package routers

import (
	authpb "socio/internal/grpc/auth/proto"
	uspb "socio/internal/grpc/user/proto"
	rest "socio/internal/rest/auth"
	customtime "socio/pkg/time"

	"github.com/gorilla/mux"
)

func MountAuthRouter(rootRouter *mux.Router, authClient authpb.AuthClient, userClient uspb.UserClient) {
	r := rootRouter.PathPrefix("/auth").Subrouter()

	h := rest.NewAuthHandler(authClient, userClient, &customtime.RealTimeProvider{})

	r.HandleFunc("/login", h.HandleLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", h.HandleRegistration).Methods("POST", "OPTIONS")
	r.HandleFunc("/logout", h.HandleLogout).Methods("DELETE", "OPTIONS")
}
