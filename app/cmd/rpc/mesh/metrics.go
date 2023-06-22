package mesh

import (
	"context"
	"fmt"
	"github.com/alitto/pond"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	stdPrometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	log2 "log"
	"net/http"
	"time"
)

var (
	ServiceMetricsNamespace   = "geo_mesh"
	StatTypeLabel             = "stat_type"
	NodeNameLabel             = "node_name"
	ServicerLabel             = "servicer_address"
	ChainIDLabel              = "chain_id"
	ChainNameLabel            = "chain_name"
	SessionHeightLabel        = "session_height"
	ApplicationPublicKeyLabel = "application_public_key"
	ReQueueLabel              = "requeue"
	NotifyLabel               = "is_notify"
	StatusTypeLabel           = "status_type"
	StatusCodeLabel           = "status_code"
	InstanceMoniker           = "moniker"

	runningWorkers                *stdPrometheus.GaugeVec
	idleWorkers                   *stdPrometheus.GaugeVec
	tasksSubmittedTotal           *stdPrometheus.GaugeVec
	tasksWaitingTotal             *stdPrometheus.GaugeVec
	successTasksTotal             *stdPrometheus.GaugeVec
	failedTasksTotal              *stdPrometheus.GaugeVec
	completedTasksTotal           *stdPrometheus.GaugeVec
	minWorker                     *stdPrometheus.GaugeVec
	maxWorker                     *stdPrometheus.GaugeVec
	maxCapacity                   *stdPrometheus.GaugeVec
	relayCounter                  *stdPrometheus.CounterVec
	relayTime                     *stdPrometheus.HistogramVec
	relayHandlerTime              *stdPrometheus.HistogramVec
	errCounter                    *stdPrometheus.CounterVec
	chainTime                     *stdPrometheus.HistogramVec
	optimisticSessionQueueCounter *stdPrometheus.CounterVec
)

type ServiceMetric struct {
	// relay done (between mesh and app)
	RelayCount   *stdPrometheus.CounterVec   `json:"relay_count"`
	ErrCount     *stdPrometheus.CounterVec   `json:"err_count"`
	AvgRelayTime *stdPrometheus.HistogramVec `json:"avg_relay_time"`
	// relay notified (between mesh and fullnode)
	NotifyRelayCount   *stdPrometheus.CounterVec   `json:"notify_relay_count"`
	NotifyErrCount     *stdPrometheus.CounterVec   `json:"notify_err_count"`
	NotifyAvgRelayTime *stdPrometheus.HistogramVec `json:"notify_avg_relay_time"`
}

func getErrorLabelSignature() []string {
	baseLabels := []string{InstanceMoniker, ChainIDLabel, ChainNameLabel, NotifyLabel, StatusTypeLabel, StatusCodeLabel}
	if app.GlobalMeshConfig.MetricsAttachServicerLabel {
		baseLabels = append(baseLabels, ServicerLabel)
	}
	return baseLabels
}

func getLabelSignature() []string {
	baseLabels := []string{InstanceMoniker, ChainIDLabel, ChainNameLabel, NotifyLabel}
	if app.GlobalMeshConfig.MetricsAttachServicerLabel {
		baseLabels = append(baseLabels, ServicerLabel)
	}
	return baseLabels
}

func getChainLabelSignature() []string {
	baseLabels := []string{InstanceMoniker, ChainIDLabel, ChainNameLabel, StatusCodeLabel}
	return baseLabels
}

func getSessionStorageLabelSignature() []string {
	baseLabels := []string{InstanceMoniker, ChainIDLabel, ChainNameLabel, SessionHeightLabel, ApplicationPublicKeyLabel, ReQueueLabel}
	if app.GlobalMeshConfig.MetricsAttachServicerLabel {
		baseLabels = append(baseLabels, ServicerLabel)
	}
	return baseLabels
}

func getMetricLabelSignature() []string {
	// just in case it change in the future, 1 place to update.
	return []string{InstanceMoniker, StatTypeLabel, NodeNameLabel}
}

var (
	prometheusServer *http.Server
)

// getLabel - return properly formatted prometheus label
func getLabel(nodeAddress *sdk.Address, chainID string, notify bool) map[string]string {
	labels := map[string]string{
		// useful to identify different mesh instances against many writing to same prometheus like cross region.
		InstanceMoniker: app.GlobalMeshConfig.MetricsMoniker,
		ChainIDLabel:    chainID,
		NotifyLabel:     fmt.Sprintf("%v", notify),
	}

	if name, ok := ChainNameMap.Load(chainID); ok {
		labels[ChainNameLabel] = name
	} else {
		// fallback
		labels[ChainNameLabel] = chainID
	}

	if app.GlobalMeshConfig.MetricsAttachServicerLabel {
		labels[ServicerLabel] = nodeAddress.String()
	}

	return labels
}

func getErrorLabel(nodeAddress *sdk.Address, chainID string, notify bool, statusType, statusCode string) map[string]string {
	labels := getLabel(nodeAddress, chainID, notify)

	labels[StatusTypeLabel] = statusType
	labels[StatusCodeLabel] = statusCode

	return labels
}

func getChainLabel(chainID string, code int) map[string]string {
	labels := map[string]string{
		// useful to identify different mesh instances against many writing to same prometheus like cross region.
		InstanceMoniker: app.GlobalMeshConfig.MetricsMoniker,
		ChainIDLabel:    chainID,
		StatusCodeLabel: fmt.Sprintf("%d", code),
	}

	if name, ok := ChainNameMap.Load(chainID); ok {
		labels[ChainNameLabel] = name
	} else {
		// fallback
		labels[ChainNameLabel] = chainID
	}

	return labels
}

func getMetricLabel(statType, name string) map[string]string {
	return map[string]string{
		InstanceMoniker: app.GlobalMeshConfig.MetricsMoniker,
		StatTypeLabel:   statType,
		NodeNameLabel:   name,
	}
}

func getSessionStorageLabel(nodeAddress string, chainID string, applicationPubKey string, sessionHeight int64, isRequeue bool) map[string]string {
	labels := map[string]string{
		// useful to identify different mesh instances against many writing to same prometheus like cross region.
		InstanceMoniker:           app.GlobalMeshConfig.MetricsMoniker,
		ChainIDLabel:              chainID,
		ApplicationPublicKeyLabel: applicationPubKey,
		SessionHeightLabel:        fmt.Sprintf("%d", sessionHeight),
		ReQueueLabel:              fmt.Sprintf("%v", isRequeue),
	}

	if name, ok := ChainNameMap.Load(chainID); ok {
		labels[ChainNameLabel] = name
	} else {
		// fallback
		labels[ChainNameLabel] = chainID
	}

	if app.GlobalMeshConfig.MetricsAttachServicerLabel {
		labels[ServicerLabel] = nodeAddress
	}

	return labels
}

// addRelayFor - accumulate a relay on servicer and per chain counters.
func addRelayFor(chainID string, relayDuration float64, nodeAddress *sdk.Address, notify bool) {
	// add relay to accumulated count
	labels := getLabel(nodeAddress, chainID, notify)
	relayCounter.With(labels).Add(1)
	relayTime.With(labels).Observe(relayDuration)
}

// addHandlerRelayFor - accumulate a handler relay on servicer and per chain counters. (without chain)
func addHandlerRelayFor(chainID string, relayDuration float64, nodeAddress *sdk.Address, notify bool) {
	// add relay to accumulated count
	labels := getLabel(nodeAddress, chainID, notify)
	relayHandlerTime.With(labels).Observe(relayDuration)
}

// addChainTimeFor - add chain call time metric
func addChainTimeFor(chainID string, duration float64, statusCode int) {
	// add relay to accumulated count
	labels := getChainLabel(chainID, statusCode)
	chainTime.With(labels).Observe(duration)
}

// addErrorFor - accumulate a relay on servicer and per chain counters.
func addErrorFor(chainID string, nodeAddress *sdk.Address, notify bool, statusType, statusCode string) {
	// add relay to accumulated count
	labels := getErrorLabel(nodeAddress, chainID, notify, statusType, statusCode)
	errCounter.With(labels).Add(1)
}

// addQueueCountFor - accumulate session storage queue counts for optimistic sessions
func addQueueCountFor(session *Session, nodeAddress string, isRequeue bool) {
	labels := getSessionStorageLabel(nodeAddress, session.Chain, session.AppPublicKey, session.BlockHeight, isRequeue)
	optimisticSessionQueueCounter.With(labels).Add(1)
}

type Metrics struct {
	name string
	// internal worker to queue the metrics tasks
	worker *pond.WorkerPool
	// worker pool that will be tracked
	pool *pond.WorkerPool

	cron *cron.Cron
}

// report - add values to all the collectors
func (m *Metrics) report() {
	nodeWorkerLabel := getMetricLabel("node", m.name)
	// pool metrics
	runningWorkers.With(nodeWorkerLabel).Set(float64(m.pool.RunningWorkers()))
	idleWorkers.With(nodeWorkerLabel).Set(float64(m.pool.IdleWorkers()))
	tasksSubmittedTotal.With(nodeWorkerLabel).Set(float64(m.pool.SubmittedTasks()))
	tasksWaitingTotal.With(nodeWorkerLabel).Set(float64(m.pool.WaitingTasks()))
	successTasksTotal.With(nodeWorkerLabel).Set(float64(m.pool.SuccessfulTasks()))
	failedTasksTotal.With(nodeWorkerLabel).Set(float64(m.pool.FailedTasks()))
	completedTasksTotal.With(nodeWorkerLabel).Set(float64(m.pool.CompletedTasks()))
	minWorker.With(nodeWorkerLabel).Set(float64(m.pool.MinWorkers()))
	maxWorker.With(nodeWorkerLabel).Set(float64(m.pool.MaxWorkers()))
	maxCapacity.With(nodeWorkerLabel).Set(float64(m.pool.MaxCapacity()))

	metricsWorkerLabel := getMetricLabel("metric", m.name)
	// internal metrics
	// pool metrics
	runningWorkers.With(metricsWorkerLabel).Set(float64(m.worker.RunningWorkers()))
	idleWorkers.With(metricsWorkerLabel).Set(float64(m.worker.IdleWorkers()))
	tasksSubmittedTotal.With(metricsWorkerLabel).Set(float64(m.worker.SubmittedTasks()))
	tasksWaitingTotal.With(metricsWorkerLabel).Set(float64(m.worker.WaitingTasks()))
	successTasksTotal.With(metricsWorkerLabel).Set(float64(m.worker.SuccessfulTasks()))
	failedTasksTotal.With(metricsWorkerLabel).Set(float64(m.worker.FailedTasks()))
	completedTasksTotal.With(metricsWorkerLabel).Set(float64(m.worker.CompletedTasks()))
	minWorker.With(metricsWorkerLabel).Set(float64(m.worker.MinWorkers()))
	maxWorker.With(metricsWorkerLabel).Set(float64(m.worker.MaxWorkers()))
	maxCapacity.With(metricsWorkerLabel).Set(float64(m.worker.MaxCapacity()))
}

// AddServiceMetricErrorFor - add to prometheus metrics an error for a servicer
func (m *Metrics) AddServiceMetricErrorFor(blockchain string, address *sdk.Address, notify bool, statusType, statusCode string) {
	m.worker.Submit(func() {
		addErrorFor(blockchain, address, notify, statusType, statusCode)
	})
}

// AddServiceMetricRelayFor - add to prometheus metrics a relay for a servicer
func (m *Metrics) AddServiceMetricRelayFor(relay *pocketTypes.Relay, address *sdk.Address, relayTime time.Duration, notify bool) {
	m.worker.Submit(func() {
		logger.Debug(fmt.Sprintf("adding metric for relay %s", relay.RequestHashString()))
		addRelayFor(
			relay.Proof.Blockchain,
			float64(relayTime.Milliseconds()),
			address,
			notify,
		)
	})
}

// AddServiceHandlerMetricRelayFor - add metrics of the handler of a relay (without chain execution)
func (m *Metrics) AddServiceHandlerMetricRelayFor(relay *pocketTypes.Relay, address *sdk.Address, relayTime time.Duration, notify bool) {
	m.worker.Submit(func() {
		logger.Debug(fmt.Sprintf("adding handler metric for relay %s", relay.RequestHashString()))
		addHandlerRelayFor(
			relay.Proof.Blockchain,
			float64(relayTime.Milliseconds()),
			address,
			notify,
		)
	})
}

// AddChainMetricFor - add metrics for chain call time
func (m *Metrics) AddChainMetricFor(chain string, duration time.Duration, statusCode int) {
	m.worker.Submit(func() {
		addChainTimeFor(chain, float64(duration.Milliseconds()), statusCode)
	})
}

// AddSessionStorageMetricQueueFor - add metrics for optimistic sessions queue/requeue jobs
func (m *Metrics) AddSessionStorageMetricQueueFor(session *Session, address string, isRequeue bool) {
	m.worker.Submit(func() {
		logger.Debug("adding metric for optimistic session")
		addQueueCountFor(session, address, isRequeue)
	})
}

// Start - start metrics cron
func (m *Metrics) Start() {
	m.cron.Start()
}

// Stop - stop metrics crons and internal worker pool
func (m *Metrics) Stop() {
	// stop cron job
	m.cron.Stop()
	// stop worker and wait it dispatch queue tasks because they are not persisted
	m.worker.StopAndWait()
}

// NewWorkerPoolMetrics - create a metric instance with the needed worker and crons
func NewWorkerPoolMetrics(name string, pool *pond.WorkerPool) *Metrics {
	worker := NewWorkerPool(
		name,
		app.GlobalMeshConfig.MetricsWorkerStrategy,
		app.GlobalMeshConfig.MetricsMaxWorkers,
		app.GlobalMeshConfig.MetricsMaxWorkersCapacity,
		app.GlobalMeshConfig.MetricsWorkersIdleTimeout,
	)

	metrics := &Metrics{
		name:   name,
		worker: worker,
		pool:   pool,
		cron:   cron.New(),
	}

	// schedule the metrics report using cron
	_, err := metrics.cron.AddFunc(fmt.Sprintf("@every %ds", app.GlobalMeshConfig.MetricsReportInterval), func() {
		metrics.report()
	})

	if err != nil {
		log2.Fatal(err)
	}

	return metrics
}

// RegisterMetrics - register to prom all the collectors
func RegisterMetrics() {
	// pool metrics
	runningWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "workers_running",
			Help:      "Number of running worker goroutines",
		},
		getMetricLabelSignature(),
	)

	idleWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "workers_idle",
			Help:      "Number of idle worker goroutines",
		},
		getMetricLabelSignature(),
	)

	tasksSubmittedTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "tasks_submitted_total",
			Help:      "Number of tasks submitted",
		},
		getMetricLabelSignature(),
	)
	tasksWaitingTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "tasks_waiting_total",
			Help:      "Number of tasks waiting in the queue",
		},
		getMetricLabelSignature(),
	)
	successTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "tasks_successful_total",
			Help:      "Number of tasks that completed successfully",
		},
		getMetricLabelSignature(),
	)
	failedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "tasks_failed_total",
			Help:      "Number of tasks that completed with panic",
		},
		getMetricLabelSignature(),
	)
	completedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "tasks_completed_total",
			Help:      "Number of tasks that completed either successfully or with panic",
		},
		getMetricLabelSignature(),
	)
	minWorker = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "min_workers",
			Help:      "Number min workers of node pool",
		},
		getMetricLabelSignature(),
	)
	maxWorker = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "max_workers",
			Help:      "Number max workers of node pool",
		},
		getMetricLabelSignature(),
	)
	maxCapacity = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "max_capacity",
			Help:      "Number max capacity of node pool",
		},
		getMetricLabelSignature(),
	)

	// relay counter metric
	relayCounter = stdPrometheus.NewCounterVec(
		stdPrometheus.CounterOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "relay_count",
			Help:      "Number of relays executed",
		},
		getLabelSignature(),
	)
	// Avg relay time histogram metric (include chain execution)
	relayTime = stdPrometheus.NewHistogramVec(
		stdPrometheus.HistogramOpts{
			Namespace:   ModuleName,
			Subsystem:   ServiceMetricsNamespace,
			Name:        "relay_time",
			Help:        "Relay duration in milliseconds (with chain call)",
			ConstLabels: nil,
			Buckets:     stdPrometheus.LinearBuckets(1, 20, 20),
		},
		getLabelSignature(),
	)
	// Avg relay time histogram metric (without chain execution)
	relayHandlerTime = stdPrometheus.NewHistogramVec(
		stdPrometheus.HistogramOpts{
			Namespace:   ModuleName,
			Subsystem:   ServiceMetricsNamespace,
			Name:        "relay_handler_time",
			Help:        "Relay duration in milliseconds (without chain call)",
			ConstLabels: nil,
			Buckets:     stdPrometheus.LinearBuckets(1, 20, 20),
		},
		getLabelSignature(),
	)
	// err counter metric
	errCounter = stdPrometheus.NewCounterVec(
		stdPrometheus.CounterOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "error_count",
			Help:      "Number of errors resulting from relays (mesh or chain)",
		},
		getErrorLabelSignature(),
	)

	// Avg chain execute time histogram metric
	chainTime = stdPrometheus.NewHistogramVec(
		stdPrometheus.HistogramOpts{
			Namespace:   ModuleName,
			Subsystem:   ServiceMetricsNamespace,
			Name:        "chain_time",
			Help:        "Chain call duration in milliseconds",
			ConstLabels: nil,
			Buckets:     stdPrometheus.LinearBuckets(1, 20, 20),
		},
		getChainLabelSignature(),
	)

	// optimistic session queue counter
	optimisticSessionQueueCounter = stdPrometheus.NewCounterVec(
		stdPrometheus.CounterOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "optimistic_session_queue_count",
			Help:      "Number of optimistic sessions delivered to queue to be validated",
		},
		getSessionStorageLabelSignature(),
	)

	stdPrometheus.MustRegister(
		// pool collectors
		runningWorkers,
		idleWorkers,
		tasksSubmittedTotal,
		tasksWaitingTotal,
		successTasksTotal,
		failedTasksTotal,
		completedTasksTotal,
		minWorker,
		maxWorker,
		maxCapacity,
		// servicer
		relayCounter,
		relayTime,
		relayHandlerTime,
		errCounter,
		// chain
		chainTime,
		// session storage
		optimisticSessionQueueCounter,
	)
}

// UnregisterMetrics - unregister from prom all the collectors
func UnregisterMetrics() {
	// unregister node collectors
	stdPrometheus.Unregister(runningWorkers)
	stdPrometheus.Unregister(idleWorkers)
	stdPrometheus.Unregister(tasksSubmittedTotal)
	stdPrometheus.Unregister(tasksWaitingTotal)
	stdPrometheus.Unregister(successTasksTotal)
	stdPrometheus.Unregister(failedTasksTotal)
	stdPrometheus.Unregister(completedTasksTotal)
}

// StartPrometheusServer - starts a Prometheus HTTP server, listening for metrics
// collectors on addr.
func StartPrometheusServer() *http.Server {
	// Register node and metrics collectors
	RegisterMetrics()

	prometheusServer = &http.Server{
		Addr: ":" + app.GlobalMeshConfig.PrometheusAddr,
		Handler: promhttp.InstrumentMetricHandler(
			stdPrometheus.DefaultRegisterer, promhttp.HandlerFor(
				stdPrometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					MaxRequestsInFlight: app.GlobalMeshConfig.PrometheusMaxOpenfiles,
				},
			),
		),
	}

	go func() {
		if err := prometheusServer.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			logger.Error("Prometheus HTTP server ListenAndServe", "err", err)
		}
	}()

	return prometheusServer
}

// StopPrometheusServer - stop prometheus server instance.
func StopPrometheusServer() {
	// Unregister node and metrics collectors
	UnregisterMetrics()
	// stop receiving new requests
	logger.Info("stopping prometheus http server...")
	if prometheusServer != nil {
		if err := prometheusServer.Shutdown(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("prometheus http server shutdown error: %s", err.Error()))
		}
	}
	logger.Info("prometheus http server stopped!")
}
