package routers

import (
	authpb "socio/internal/grpc/auth/proto"
	rest "socio/internal/rest/csrf"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"

	"github.com/gorilla/mux"
)

func MountCSRFRouter(rootRouter *mux.Router, authManager authpb.AuthClient) {
	r := rootRouter.PathPrefix("/csrf").Subrouter()

	h := rest.NewCSRFHandler(customtime.RealTimeProvider{})

	r.HandleFunc("/", h.GetCSRFToken).Methods("GET", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(authManager))
}
