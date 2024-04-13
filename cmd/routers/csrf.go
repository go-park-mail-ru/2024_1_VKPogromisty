package routers

import (
	rest "socio/internal/rest/csrf"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"
	"socio/usecase/auth"

	"github.com/gorilla/mux"
)

func MountCSRFRouter(rootRouter *mux.Router, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/csrf").Subrouter()

	h := rest.NewCSRFHandler(customtime.RealTimeProvider{})

	r.HandleFunc("/", h.GetCSRFToken).Methods("GET", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
}
