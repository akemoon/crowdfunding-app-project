package metrics

import "github.com/prometheus/client_golang/prometheus"

type FinishWorkerMetrics struct {
	LastRunUnix prometheus.Gauge
}

func NewFinishWorkerMetrics(reg prometheus.Registerer) *FinishWorkerMetrics {
	lastRunUnix := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "project_finish_worker_last_run_unix_seconds",
		Help: "Unix timestamp of the last finish worker run.",
	})

	reg.MustRegister(lastRunUnix)

	return &FinishWorkerMetrics{
		LastRunUnix: lastRunUnix,
	}
}
