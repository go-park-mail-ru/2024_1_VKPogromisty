package routers

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MountMetricsRouter(rootRouter *mux.Router) {
	r := rootRouter.PathPrefix("/metrics").Subrouter()

	r.Handle("/", promhttp.Handler()).Methods("GET", "OPTIONS")
}
