package metrics

import "github.com/prometheus/client_golang/prometheus"

type PublishWorkerMetrics struct {
	LastRunUnix prometheus.Gauge
}

func NewPublishWorkerMetrics(reg prometheus.Registerer) *PublishWorkerMetrics {
	lastRunUnix := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "project_publish_worker_last_run_unix_seconds",
		Help: "Unix timestamp of the last publish worker run.",
	})

	reg.MustRegister(lastRunUnix)

	return &PublishWorkerMetrics{
		LastRunUnix: lastRunUnix,
	}
}
