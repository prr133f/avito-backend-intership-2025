package metrics

import "github.com/prometheus/client_golang/prometheus"

func InitMetrics() {
	prometheus.MustRegister(RequestsTotal)
}
