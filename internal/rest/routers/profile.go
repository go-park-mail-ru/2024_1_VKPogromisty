package routers

import (
	authpb "socio/internal/grpc/auth/proto"
	uspb "socio/internal/grpc/user/proto"
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/profile"
	customtime "socio/pkg/time"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
)

func MountProfileRouter(rootRouter *mux.Router, userClient uspb.UserClient, authClient authpb.AuthClient) {
	r := rootRouter.PathPrefix("/profile").Subrouter()

	h := rest.NewProfileHandler(userClient)

	r.HandleFunc("/search", h.HandleSearchByName).Methods("GET", "OPTIONS")
	r.HandleFunc("/{userID:[0-9]+}", h.HandleGetProfile).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleGetProfile).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleUpdateProfile).Methods("PUT", "OPTIONS")
	r.HandleFunc("/", h.HandleDeleteProfile).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(authClient))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
