package routers

import (
	"socio/handlers"

	"github.com/gorilla/mux"
)

func MountStaticRouter(rootRouter *mux.Router) {
	r := rootRouter.PathPrefix("/static").Subrouter()
	handler := handlers.StaticHandler{}

	r.HandleFunc("/{fileName}", handler.HandleServeStatic).Methods("GET")
}
