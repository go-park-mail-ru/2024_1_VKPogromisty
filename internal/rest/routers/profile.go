package routers

import (
	uspb "socio/internal/grpc/user/proto"
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/profile"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
)

func MountProfileRouter(rootRouter *mux.Router, userClient uspb.UserClient, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/profile").Subrouter()

	h := rest.NewProfileHandler(userClient)

	r.HandleFunc("/{userID}", h.HandleGetProfile).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleGetProfile).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleUpdateProfile).Methods("PUT", "OPTIONS")
	r.HandleFunc("/", h.HandleDeleteProfile).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
