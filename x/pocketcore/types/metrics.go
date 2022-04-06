package types

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdPrometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tendermint/tendermint/libs/log"
	"net/http"
	"sync"
)

var (
	globalServiceMetrics *ServiceMetrics
)

const (
	ServiceMetricsKey       = "service"
	ServiceMetricsNamespace = ServiceMetricsKey
	RelayCountName          = "relay_count_for_"
	RelayCountHelp          = "the number of relays executed against: "
	ChallengeCountName      = "challenge_count_for_"
	ChallengeCountHelp      = "the number of challenges executed against: "
	ErrCountName            = "err_count_for_"
	ErrCountHelp            = "the number of errors resulting from relays executed against: "
	AvgRelayHistName        = "avg_relay_time_for_"
	AvgrelayHistHelp        = "the average relay time in ms executed against: "
	SessionsCountName       = "sessions_count_for_"
	SessionsCountHelp       = "the number of unique sessions generated for: "
	UPOKTCountName          = "tokens_earned_for_"
	UPOKTCountHelp          = "the number of tokens earned in uPOKT for : "
)

type ServiceMetrics struct {
	l               sync.Mutex
	tmLogger        log.Logger
	ServiceMetric   `json:"accumulated_service_metrics"` // total metrics
	NonNativeChains map[string]ServiceMetric             `json:"individual_service_metrics"` // metrics per chain
	prometheusSrv   *http.Server
}

type ServiceMetricsEncodable struct {
	ServiceMetric   `json:"accumulated_service_metrics"` // total metrics
	NonNativeChains []ServiceMetric                      `json:"individual_service_metrics"` // metrics per chain
}

func GlobalServiceMetric() *ServiceMetrics {
	return globalServiceMetrics
}

func InitGlobalServiceMetric(hostedBlockchains *HostedBlockchains, logger log.Logger, addr string, maxOpenConn int) {
	// create a new service metric
	serviceMetric := NewServiceMetrics(hostedBlockchains, logger)
	// set the service metrics
	globalServiceMetrics = serviceMetric
	// start metrics server
	globalServiceMetrics.prometheusSrv = globalServiceMetrics.StartPrometheusServer(addr, maxOpenConn)
}

func StopServiceMetrics() {
	if err := GlobalServiceMetric().prometheusSrv.Shutdown(context.Background()); err != nil {
		GlobalServiceMetric().tmLogger.Error("unable to shutdown service metrics server: ", err.Error())
	}
}

// startPrometheusServer starts a Prometheus HTTP server, listening for metrics
// collectors on addr.
func (sm *ServiceMetrics) StartPrometheusServer(addr string, maxOpenConn int) *http.Server {
	srv := &http.Server{
		Addr: ":" + addr,
		Handler: promhttp.InstrumentMetricHandler(
			stdPrometheus.DefaultRegisterer, promhttp.HandlerFor(
				stdPrometheus.DefaultGatherer,
				promhttp.HandlerOpts{MaxRequestsInFlight: maxOpenConn},
			),
		),
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			sm.tmLogger.Error("Prometheus HTTP server ListenAndServe", "err", err)
		}
	}()
	return srv
}

func (sm *ServiceMetrics) AddRelayFor(networkID string) {
	sm.l.Lock()
	defer sm.l.Unlock()
	// attempt to locate nn chain
	nnc, ok := sm.NonNativeChains[networkID]
	if !ok {
		sm.tmLogger.Error("unable to find corresponding networkID in service metrics: ", networkID)
		sm.NonNativeChains[networkID] = NewServiceMetricsFor(networkID)
		return
	}
	// add relay to accumulated count
	sm.RelayCount.Add(1)
	// add to individual relay count
	nnc.RelayCount.Add(1)
	// update nnc
	sm.NonNativeChains[networkID] = nnc
}

func (sm *ServiceMetrics) AddChallengeFor(networkID string) {
	sm.l.Lock()
	defer sm.l.Unlock()
	// attempt to locate nn chain
	nnc, ok := sm.NonNativeChains[networkID]
	if !ok {
		sm.tmLogger.Error("unable to find corresponding networkID in service metrics: ", networkID)
		sm.NonNativeChains[networkID] = NewServiceMetricsFor(networkID)
		return
	}
	// add to accumulated count
	sm.ChallengeCount.Add(1)
	// add to individual count
	nnc.ChallengeCount.Add(1)
	// update nnc
	sm.NonNativeChains[networkID] = nnc
}

func (sm *ServiceMetrics) AddErrorFor(networkID string) {
	sm.l.Lock()
	defer sm.l.Unlock()
	// attempt to locate nn chain
	nnc, ok := sm.NonNativeChains[networkID]
	if !ok {
		sm.tmLogger.Error("unable to find corresponding networkID in service metrics: ", networkID)
		sm.NonNativeChains[networkID] = NewServiceMetricsFor(networkID)
		return
	}
	// add to accumulated count
	sm.ErrCount.Add(1)
	// add to individual count
	nnc.ErrCount.Add(1)
	// update nnc
	sm.NonNativeChains[networkID] = nnc
}

func (sm *ServiceMetrics) AddRelayTimingFor(networkID string, relayTime float64) {
	sm.l.Lock()
	defer sm.l.Unlock()
	// attempt to locate nn chain
	nnc, ok := sm.NonNativeChains[networkID]
	if !ok {
		sm.tmLogger.Error("unable to find corresponding networkID in service metrics: ", networkID)
		sm.NonNativeChains[networkID] = NewServiceMetricsFor(networkID)
		return
	}
	// add to accumulated hist
	sm.AverageRelayTime.Observe(relayTime)
	// add to individual hist
	nnc.AverageRelayTime.Observe(relayTime)
	// update nnc
	sm.NonNativeChains[networkID] = nnc
}

func (sm *ServiceMetrics) AddSessionFor(networkID string) {
	sm.l.Lock()
	defer sm.l.Unlock()
	// attempt to locate nn chain
	nnc, ok := sm.NonNativeChains[networkID]
	if !ok {
		sm.tmLogger.Error("unable to find corresponding networkID in service metrics: ", networkID)
		sm.NonNativeChains[networkID] = NewServiceMetricsFor(networkID)

		return
	}
	// add to accumulated count
	sm.TotalSessions.Add(1)
	// add to individual count
	nnc.TotalSessions.Add(1)
	// update nnc
	sm.NonNativeChains[networkID] = nnc
}

func (sm *ServiceMetrics) AddUPOKTEarnedFor(networkID string, upoktEarned float64) {
	sm.l.Lock()
	defer sm.l.Unlock()
	// attempt to locate nn chain
	nnc, ok := sm.NonNativeChains[networkID]
	if !ok {
		sm.tmLogger.Error("unable to find corresponding networkID in service metrics: ", networkID)
		sm.NonNativeChains[networkID] = NewServiceMetricsFor(networkID)
		return
	}
	// add to accumulated count
	sm.UPOKTEarned.Add(1)
	// add to individual count
	nnc.UPOKTEarned.Add(1)
	// update nnc
	sm.NonNativeChains[networkID] = nnc
}

func KeyForServiceMetrics() []byte {
	return []byte(ServiceMetricsKey)
}

func NewServiceMetrics(hostedBlockchains *HostedBlockchains, logger log.Logger) *ServiceMetrics {
	serviceMetrics := ServiceMetrics{
		ServiceMetric:   NewServiceMetricsFor("all"),
		NonNativeChains: make(map[string]ServiceMetric),
	}
	if hostedBlockchains != nil {
		for _, hb := range hostedBlockchains.M {
			serviceMetrics.NonNativeChains[hb.ID] = NewServiceMetricsFor(hb.ID)
		}
	}
	// add the logger
	serviceMetrics.tmLogger = logger
	// return the metrics
	return &serviceMetrics
}

type ServiceMetric struct {
	RelayCount       metrics.Counter   `json:"relay_count"`
	ChallengeCount   metrics.Counter   `json:"challenge_count"`
	ErrCount         metrics.Counter   `json:"err_count"`
	AverageRelayTime metrics.Histogram `json:"avg_relay_time"`
	TotalSessions    metrics.Counter   `json:"total_sessions"`
	UPOKTEarned      metrics.Counter   `json:"upokt_earned"`
}

func NewServiceMetricsFor(networkID string) ServiceMetric {
	//labels := make([]string, 1)
	// relay counter metric
	relayCounter := prometheus.NewCounterFrom(stdPrometheus.CounterOpts{
		Namespace: ModuleName,
		Subsystem: ServiceMetricsNamespace,
		Name:      RelayCountName + networkID,
		Help:      RelayCountHelp + networkID,
	}, nil)
	// challenge counter metric
	challengeCounter := prometheus.NewCounterFrom(stdPrometheus.CounterOpts{
		Namespace: ModuleName,
		Subsystem: ServiceMetricsNamespace,
		Name:      ChallengeCountName + networkID,
		Help:      ChallengeCountHelp + networkID,
	}, nil)
	// err counter metric
	errCounter := prometheus.NewCounterFrom(stdPrometheus.CounterOpts{
		Namespace: ModuleName,
		Subsystem: ServiceMetricsNamespace,
		Name:      ErrCountName + networkID,
		Help:      ErrCountHelp + networkID,
	}, nil)
	// Avg relay time histogram metric
	avgRelayTime := prometheus.NewHistogramFrom(stdPrometheus.HistogramOpts{
		Namespace:   ModuleName,
		Subsystem:   ServiceMetricsNamespace,
		Name:        AvgRelayHistName + networkID,
		Help:        AvgrelayHistHelp + networkID,
		ConstLabels: nil,
		Buckets:     stdPrometheus.LinearBuckets(1, 20, 20),
	}, nil)
	// session counter metric
	totalSessions := prometheus.NewCounterFrom(stdPrometheus.CounterOpts{
		Namespace: ModuleName,
		Subsystem: ServiceMetricsNamespace,
		Name:      SessionsCountName + networkID,
		Help:      SessionsCountHelp + networkID,
	}, nil)
	// tokens earned metric
	uPOKTEarned := prometheus.NewCounterFrom(stdPrometheus.CounterOpts{
		Namespace: ModuleName,
		Subsystem: ServiceMetricsNamespace,
		Name:      UPOKTCountName + networkID,
		Help:      UPOKTCountHelp + networkID,
	}, nil)
	return ServiceMetric{
		RelayCount:       relayCounter,
		ChallengeCount:   challengeCounter,
		ErrCount:         errCounter,
		AverageRelayTime: avgRelayTime,
		TotalSessions:    totalSessions,
		UPOKTEarned:      uPOKTEarned,
	}
}
