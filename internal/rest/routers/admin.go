package routers

import (
	uspb "socio/internal/grpc/user/proto"
	rest "socio/internal/rest/admin"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
)

func MountAdminRouter(rootRouter *mux.Router, userClient uspb.UserClient, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/admin").Subrouter()

	h := rest.NewAdminHandler(userClient)

	r.HandleFunc("/{userID}", h.HandleGetAdminByUserID).Methods("GET", "OPTIONS")
	r.HandleFunc("/all", h.HandleGetAdmins).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleCreateAdmin).Methods("POST", "OPTIONS")
	r.HandleFunc("/", h.HandleDeleteAdmin).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
	r.Use(middleware.CreateCheckAdminMiddleware(userClient))
}
