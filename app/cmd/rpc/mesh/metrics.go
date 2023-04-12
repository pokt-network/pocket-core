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
	"github.com/puzpuzpuz/xsync"
	"github.com/robfig/cron/v3"
	log2 "log"
	"net/http"
	"time"
)

var (
	ServiceMetricsNamespace = "geo_mesh"
	NodeWorkerLabel         = "node_worker"
	MetricsWorkerLabel      = "metrics_worker"
	ServicerLabel           = "validator_address"
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

var (
	prometheusServer *http.Server

	// total metrics per servicer
	servicerMetrics = &ServiceMetric{}

	// metric per chain
	perChainMetrics = xsync.NewMapOf[*ServiceMetric]()

	// pool metrics
	runningWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_workers_running",
			Help:      "Number of running worker goroutines",
		},
		[]string{NodeWorkerLabel})

	idleWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_workers_idle",
			Help:      "Number of idle worker goroutines",
		},
		[]string{NodeWorkerLabel})

	tasksSubmittedTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_submitted_total",
			Help:      "Number of tasks submitted",
		},
		[]string{NodeWorkerLabel})
	tasksWaitingTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_waiting_total",
			Help:      "Number of tasks waiting in the queue",
		},
		[]string{NodeWorkerLabel})
	successTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_successful_total",
			Help:      "Number of tasks that completed successfully",
		},
		[]string{NodeWorkerLabel})
	failedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_failed_total",
			Help:      "Number of tasks that completed with panic",
		},
		[]string{NodeWorkerLabel})
	completedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_completed_total",
			Help:      "Number of tasks that completed either successfully or with panic",
		},
		[]string{NodeWorkerLabel})
	minWorker = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_min_workers",
			Help:      "Number min workers of node pool",
		},
		[]string{NodeWorkerLabel})
	maxWorker = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_max_workers",
			Help:      "Number max workers of node pool",
		},
		[]string{NodeWorkerLabel})
	maxCapacity = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_max_capacity",
			Help:      "Number max capacity of node pool",
		},
		[]string{NodeWorkerLabel})

	// internal worker metrics
	internalRunningWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_workers_running",
			Help:      "Number of running worker goroutines",
		},
		[]string{MetricsWorkerLabel})
	internalIdleWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_workers_idle",
			Help:      "Number of idle worker goroutines",
		},
		[]string{MetricsWorkerLabel})
	internalTasksSubmittedTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_submitted_total",
			Help:      "Number of tasks submitted",
		},
		[]string{MetricsWorkerLabel})
	internalTasksWaitingTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_waiting_total",
			Help:      "Number of tasks waiting in the queue",
		},
		[]string{MetricsWorkerLabel})
	internalSuccessTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_successful_total",
			Help:      "Number of tasks that completed successfully",
		},
		[]string{MetricsWorkerLabel})
	internalFailedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_failed_total",
			Help:      "Number of tasks that completed with panic",
		},
		[]string{MetricsWorkerLabel})
	internalCompletedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_completed_total",
			Help:      "Number of tasks that completed either successfully or with panic",
		},
		[]string{MetricsWorkerLabel})
	internalMinWorker = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_min_workers",
			Help:      "Number min workers of metrics pool",
		},
		[]string{MetricsWorkerLabel})
	internalMaxWorker = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_max_workers",
			Help:      "Number max workers of metrics pool",
		},
		[]string{MetricsWorkerLabel})
	internalMaxCapacity = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_max_capacity",
			Help:      "Number max capacity of metrics pool",
		},
		[]string{MetricsWorkerLabel})
)

// getServicerLabel - return properly formatted prometheus label
func getServicerLabel(nodeAddress *sdk.Address) map[string]string {
	return map[string]string{ServicerLabel: nodeAddress.String()}
}

// addRelayFor - accumulate a relay on servicer and per chain counters.
func addRelayFor(networkID string, relayTime float64, nodeAddress *sdk.Address, notify bool) {
	// add relay to accumulated count
	labels := getServicerLabel(nodeAddress)

	var relayCounter *stdPrometheus.CounterVec
	var avgRelay *stdPrometheus.HistogramVec
	var chainRelayCounter *stdPrometheus.CounterVec
	var chainAvgRelay *stdPrometheus.HistogramVec

	var chainMetric *ServiceMetric

	if cm, ok := perChainMetrics.Load(networkID); !ok {
		logger.Error("unable to find corresponding networkID in service metrics: ", networkID)
		// path on the fly
		chainMetric = NewServiceMetricsFor(networkID)
		perChainMetrics.Store(networkID, chainMetric)
	} else {
		chainMetric = cm
	}

	if !notify {
		relayCounter = servicerMetrics.RelayCount
		avgRelay = servicerMetrics.AvgRelayTime
		chainRelayCounter = chainMetric.RelayCount
		chainAvgRelay = chainMetric.AvgRelayTime
	} else {
		relayCounter = servicerMetrics.NotifyRelayCount
		avgRelay = servicerMetrics.NotifyAvgRelayTime
		chainRelayCounter = chainMetric.NotifyRelayCount
		chainAvgRelay = chainMetric.NotifyAvgRelayTime
	}

	relayCounter.With(labels).Add(1)
	avgRelay.With(labels).Observe(relayTime)
	chainRelayCounter.With(labels).Add(1)
	chainAvgRelay.With(labels).Observe(relayTime)
}

// addErrorFor - accumulate a relay on servicer and per chain counters.
func addErrorFor(networkID string, nodeAddress *sdk.Address, notify bool) {
	// add relay to accumulated count
	labels := getServicerLabel(nodeAddress)

	var errCounter *stdPrometheus.CounterVec
	var chainErrCounter *stdPrometheus.CounterVec

	var chainMetric *ServiceMetric

	if cm, ok := perChainMetrics.Load(networkID); !ok {
		logger.Error("unable to find corresponding networkID in service metrics: ", networkID)
		// path on the fly
		chainMetric = NewServiceMetricsFor(networkID)
		perChainMetrics.Store(networkID, chainMetric)
	} else {
		chainMetric = cm
	}

	if !notify {
		errCounter = servicerMetrics.ErrCount
		chainErrCounter = chainMetric.ErrCount
	} else {
		errCounter = servicerMetrics.NotifyErrCount
		chainErrCounter = chainMetric.NotifyErrCount
	}

	errCounter.With(labels).Add(1)
	chainErrCounter.With(labels).Add(1)
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
	nodeWorkerLabel := map[string]string{NodeWorkerLabel: m.name}
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

	metricsWorkerLabel := map[string]string{MetricsWorkerLabel: m.name}
	// internal metrics
	internalRunningWorkers.With(metricsWorkerLabel).Set(float64(m.worker.RunningWorkers()))
	internalIdleWorkers.With(metricsWorkerLabel).Set(float64(m.worker.IdleWorkers()))
	internalTasksSubmittedTotal.With(metricsWorkerLabel).Set(float64(m.worker.SubmittedTasks()))
	internalTasksWaitingTotal.With(metricsWorkerLabel).Set(float64(m.worker.WaitingTasks()))
	internalSuccessTasksTotal.With(metricsWorkerLabel).Set(float64(m.worker.SuccessfulTasks()))
	internalFailedTasksTotal.With(metricsWorkerLabel).Set(float64(m.worker.FailedTasks()))
	internalCompletedTasksTotal.With(metricsWorkerLabel).Set(float64(m.worker.CompletedTasks()))
	internalMinWorker.With(metricsWorkerLabel).Set(float64(m.worker.MinWorkers()))
	internalMaxWorker.With(metricsWorkerLabel).Set(float64(m.worker.MaxWorkers()))
	internalMaxCapacity.With(metricsWorkerLabel).Set(float64(m.worker.MaxCapacity()))
}

// AddServiceMetricErrorFor - add to prometheus metrics an error for a servicer
func (m *Metrics) AddServiceMetricErrorFor(blockchain string, address *sdk.Address, notify bool) {
	m.worker.Submit(func() {
		addErrorFor(blockchain, address, notify)
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

// NewServiceMetricsFor - Create prometheus vector for received network id
func NewServiceMetricsFor(networkID string) *ServiceMetric {
	// relay counter metric
	relayCounter := stdPrometheus.NewCounterVec(
		stdPrometheus.CounterOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      pocketTypes.RelayCountName + networkID,
			Help:      pocketTypes.RelayCountHelp + networkID,
		},
		[]string{ServicerLabel},
	)
	// err counter metric
	errCounter := stdPrometheus.NewCounterVec(
		stdPrometheus.CounterOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      pocketTypes.ErrCountName + networkID,
			Help:      pocketTypes.ErrCountName + networkID,
		},
		[]string{ServicerLabel},
	)
	// Avg relay time histogram metric
	avgRelayTime := stdPrometheus.NewHistogramVec(
		stdPrometheus.HistogramOpts{
			Namespace:   ModuleName,
			Subsystem:   ServiceMetricsNamespace,
			Name:        pocketTypes.AvgRelayHistName + networkID,
			Help:        pocketTypes.AvgrelayHistHelp + networkID,
			ConstLabels: nil,
			Buckets:     stdPrometheus.LinearBuckets(1, 20, 20),
		},
		[]string{ServicerLabel},
	)

	// relay counter metric
	notifyRelayCounter := stdPrometheus.NewCounterVec(
		stdPrometheus.CounterOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "notify_" + pocketTypes.RelayCountName + networkID,
			Help:      "the number of relays notified to fullnode against: " + networkID,
		},
		[]string{ServicerLabel},
	)
	// err counter metric
	notifyErrCounter := stdPrometheus.NewCounterVec(
		stdPrometheus.CounterOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "notify_" + pocketTypes.ErrCountName + networkID,
			Help:      "the number of errors resulting from notify relays executed against: " + networkID,
		},
		[]string{ServicerLabel},
	)
	// Avg relay time histogram metric
	notifyAvgRelayTime := stdPrometheus.NewHistogramVec(
		stdPrometheus.HistogramOpts{
			Namespace:   ModuleName,
			Subsystem:   ServiceMetricsNamespace,
			Name:        "notify_" + pocketTypes.AvgRelayHistName + networkID,
			Help:        "the average notify relay time in ms executed against: " + networkID,
			ConstLabels: nil,
			Buckets:     stdPrometheus.LinearBuckets(1, 20, 20),
		},
		[]string{ServicerLabel},
	)

	return &ServiceMetric{
		RelayCount:         relayCounter,
		ErrCount:           errCounter,
		AvgRelayTime:       avgRelayTime,
		NotifyRelayCount:   notifyRelayCounter,
		NotifyErrCount:     notifyErrCounter,
		NotifyAvgRelayTime: notifyAvgRelayTime,
	}
}

// RegisterMetricsForChains - Register vector for each loaded chain, this will not register a chain 2 times.
func RegisterMetricsForChains() {
	// read current chains
	currentChains := GetChains()

	currentChains.L.Lock()
	for _, chain := range currentChains.M {
		if _, ok := perChainMetrics.Load(chain.ID); !ok {
			chainMetrics := NewServiceMetricsFor(chain.ID)
			perChainMetrics.Store(chain.ID, chainMetrics)
			stdPrometheus.MustRegister(
				chainMetrics.RelayCount,
				chainMetrics.AvgRelayTime,
				chainMetrics.ErrCount,
				chainMetrics.NotifyRelayCount,
				chainMetrics.NotifyErrCount,
				chainMetrics.NotifyAvgRelayTime,
			)
		}
	}
	currentChains.L.Unlock()
}

// RegisterMetrics - register to prom all the collectors
func RegisterMetrics() {
	servicerMetrics = NewServiceMetricsFor("all")
	stdPrometheus.MustRegister(
		servicerMetrics.RelayCount,
		servicerMetrics.AvgRelayTime,
		servicerMetrics.ErrCount,
		servicerMetrics.NotifyRelayCount,
		servicerMetrics.NotifyErrCount,
		servicerMetrics.NotifyAvgRelayTime,
	)

	RegisterMetricsForChains()

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
		// internal collectors
		internalRunningWorkers,
		internalIdleWorkers,
		internalTasksSubmittedTotal,
		internalTasksWaitingTotal,
		internalSuccessTasksTotal,
		internalFailedTasksTotal,
		internalCompletedTasksTotal,
		internalMinWorker,
		internalMaxWorker,
		internalMaxCapacity,
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

	// unregister node collectors
	stdPrometheus.Unregister(internalRunningWorkers)
	stdPrometheus.Unregister(internalIdleWorkers)
	stdPrometheus.Unregister(internalTasksSubmittedTotal)
	stdPrometheus.Unregister(internalTasksWaitingTotal)
	stdPrometheus.Unregister(internalSuccessTasksTotal)
	stdPrometheus.Unregister(internalFailedTasksTotal)
	stdPrometheus.Unregister(internalCompletedTasksTotal)
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
