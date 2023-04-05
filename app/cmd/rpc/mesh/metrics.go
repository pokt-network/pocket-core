package mesh

import (
	"fmt"
	"github.com/alitto/pond"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	stdPrometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	log2 "log"
	"time"
)

var ServiceMetricsNamespace = "geo-mesh"

var (
	// pool metrics
	runningWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_workers_running",
			Help:      "Number of running worker goroutines",
		},
		[]string{"node_worker"})

	idleWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_workers_idle",
			Help:      "Number of idle worker goroutines",
		},
		[]string{"node_worker"})

	tasksSubmittedTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_submitted_total",
			Help:      "Number of tasks submitted",
		},
		[]string{"node_worker"})
	tasksWaitingTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_waiting_total",
			Help:      "Number of tasks waiting in the queue",
		},
		[]string{"node_worker"})
	successTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_successful_total",
			Help:      "Number of tasks that completed successfully",
		},
		[]string{"node_worker"})
	failedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_failed_total",
			Help:      "Number of tasks that completed with panic",
		},
		[]string{"node_worker"})
	completedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "node_tasks_completed_total",
			Help:      "Number of tasks that completed either successfully or with panic",
		},
		[]string{"worker"})

	// internal worker metrics
	internalRunningWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_workers_running",
			Help:      "Number of running worker goroutines",
		},
		[]string{"metrics_worker"})
	internalIdleWorkers = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_workers_idle",
			Help:      "Number of idle worker goroutines",
		},
		[]string{"metrics_worker"})
	internalTasksSubmittedTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_submitted_total",
			Help:      "Number of tasks submitted",
		},
		[]string{"metrics_worker"})
	internalTasksWaitingTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_waiting_total",
			Help:      "Number of tasks waiting in the queue",
		},
		[]string{"metrics_worker"})
	internalSuccessTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_successful_total",
			Help:      "Number of tasks that completed successfully",
		},
		[]string{"metrics_worker"})
	internalFailedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_failed_total",
			Help:      "Number of tasks that completed with panic",
		},
		[]string{"metrics_worker"})
	internalCompletedTasksTotal = stdPrometheus.NewGaugeVec(
		stdPrometheus.GaugeOpts{
			Namespace: ModuleName,
			Subsystem: ServiceMetricsNamespace,
			Name:      "metrics_tasks_completed_total",
			Help:      "Number of tasks that completed either successfully or with panic",
		},
		[]string{"metrics_worker"})
)

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
	nodeWorkerLabel := map[string]string{"node_worker": m.name}
	// pool metrics
	runningWorkers.With(nodeWorkerLabel).Set(float64(m.pool.RunningWorkers()))
	idleWorkers.With(nodeWorkerLabel).Set(float64(m.pool.IdleWorkers()))
	tasksSubmittedTotal.With(nodeWorkerLabel).Set(float64(m.pool.SubmittedTasks()))
	tasksWaitingTotal.With(nodeWorkerLabel).Set(float64(m.pool.WaitingTasks()))
	successTasksTotal.With(nodeWorkerLabel).Set(float64(m.pool.SuccessfulTasks()))
	failedTasksTotal.With(nodeWorkerLabel).Set(float64(m.pool.FailedTasks()))
	completedTasksTotal.With(nodeWorkerLabel).Set(float64(m.pool.CompletedTasks()))

	metricsWorkerLabel := map[string]string{"metrics_worker": m.name}
	// internal metrics
	internalRunningWorkers.With(metricsWorkerLabel).Set(float64(m.worker.RunningWorkers()))
	internalIdleWorkers.With(metricsWorkerLabel).Set(float64(m.worker.IdleWorkers()))
	internalTasksSubmittedTotal.With(metricsWorkerLabel).Set(float64(m.worker.SubmittedTasks()))
	internalTasksWaitingTotal.With(metricsWorkerLabel).Set(float64(m.worker.WaitingTasks()))
	internalSuccessTasksTotal.With(metricsWorkerLabel).Set(float64(m.worker.SuccessfulTasks()))
	internalFailedTasksTotal.With(metricsWorkerLabel).Set(float64(m.worker.FailedTasks()))
	internalCompletedTasksTotal.With(metricsWorkerLabel).Set(float64(m.worker.CompletedTasks()))
}

// AddServiceMetricErrorFor - add to prometheus metrics an error for a servicer
func (m *Metrics) AddServiceMetricErrorFor(blockchain string, address *sdk.Address) {
	m.pool.Submit(func() {
		pocketTypes.GlobalServiceMetric().AddErrorFor(blockchain, address)
	})
}

// AddServiceMetricRelayFor - add to prometheus metrics a relay for a servicer
func (m *Metrics) AddServiceMetricRelayFor(relay *pocketTypes.Relay, address *sdk.Address, relayTime time.Duration) {
	m.pool.Submit(func() {
		logger.Debug(fmt.Sprintf("adding metric for relay %s", relay.RequestHashString()))
		pocketTypes.GlobalServiceMetric().AddRelayTimingFor(
			relay.Proof.Blockchain,
			float64(relayTime.Milliseconds()),
			address,
		)
		pocketTypes.GlobalServiceMetric().AddRelayFor(relay.Proof.Blockchain, address)
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
	worker := newWorkerPool(
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
	stdPrometheus.MustRegister(
		// pool collectors
		runningWorkers,
		idleWorkers,
		tasksSubmittedTotal,
		tasksWaitingTotal,
		successTasksTotal,
		failedTasksTotal,
		completedTasksTotal,
		// internal collectors
		internalRunningWorkers,
		internalIdleWorkers,
		internalTasksSubmittedTotal,
		internalTasksWaitingTotal,
		internalSuccessTasksTotal,
		internalFailedTasksTotal,
		internalCompletedTasksTotal,
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
