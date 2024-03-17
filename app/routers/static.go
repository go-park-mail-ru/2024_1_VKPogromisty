package routers

import (
	"socio/internal/rest"

	"github.com/gorilla/mux"
)

func MountStaticRouter(rootRouter *mux.Router) {
	r := rootRouter.PathPrefix("/static").Subrouter()
	handler := rest.StaticHandler{}

	r.HandleFunc("/{fileName}", handler.HandleServeStatic).Methods("GET", "OPTIONS")
}
