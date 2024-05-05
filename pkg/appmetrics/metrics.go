package appmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Hit = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_time_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)
)
