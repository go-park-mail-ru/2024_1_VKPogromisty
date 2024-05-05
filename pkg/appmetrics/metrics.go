package appmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	AppTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_total_hits",
			Help: "Count of hits in main app.",
		},
		[]string{},
	)
	AuthTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_total_hits",
			Help: "Count of hits in auth service.",
		},
		[]string{},
	)
	PostTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "post_total_hits",
			Help: "Count of hits in post service.",
		},
		[]string{},
	)
	PublicGroupTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "public_group_total_hits",
			Help: "Count of hits in public group service.",
		},
		[]string{},
	)
	UserTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_total_hits",
			Help: "Count of hits in user service.",
		},
		[]string{},
	)
	AppHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_methods_hits",
			Help: "Count of hits in main app.",
		},
		[]string{"method", "path", "status"},
	)
	AuthHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_methods_hits",
			Help: "Count of hits in auth service.",
		},
		[]string{"method", "status"},
	)
	PostHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "post_methods_hits",
			Help: "Count of hits in post service.",
		},
		[]string{"method", "status"},
	)
	PublicGroupHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "public_group_methods_hits",
			Help: "Count of hits in public group service.",
		},
		[]string{"method", "status"},
	)
	UserHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_methods_hits",
			Help: "Count of hits in user service.",
		},
		[]string{"method", "status"},
	)
	AppHitDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_methods_response_time_ms",
			Help: "Duration of hits in main app.",
		},
		[]string{"method", "path"},
	)
	AuthHitDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "auth_methods_response_time_ms",
			Help: "Duration of hits in auth service.",
		},
		[]string{"method"},
	)
	PostHitDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "post_methods_response_time_ms",
			Help: "Duration of hits in post service.",
		},
		[]string{"method"},
	)
	PublicGroupHitDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "public_group_methods_response_time_ms",
			Help: "Duration of hits in public group service.",
		},
		[]string{"method"},
	)
	UserHitDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "user_methods_response_time_ms",
			Help: "Duration of hits in user service.",
		},
		[]string{"method"},
	)
)
