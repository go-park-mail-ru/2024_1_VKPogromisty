package routers

import (
	csatpb "socio/internal/grpc/csat/proto"
	userpb "socio/internal/grpc/user/proto"
	rest "socio/internal/rest/csat"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
)

func MountCSATRouter(rootRouter *mux.Router, CSATClient csatpb.CSATClient, userClient userpb.UserClient, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/csat").Subrouter()

	h := rest.NewCSATHandler(CSATClient)

	r.HandleFunc("/pool/", h.CreatePool).Methods("POST", "OPTIONS")
	r.HandleFunc("/pool/", h.UpdatePool).Methods("PUT", "OPTIONS")
	r.HandleFunc("/pool/", h.DeletePool).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/question/", h.CreateQuestion).Methods("POST", "OPTIONS")
	r.HandleFunc("/question/", h.UpdateQuestion).Methods("PUT", "OPTIONS")
	r.HandleFunc("/question/", h.DeleteQuestion).Methods("DELETE", "OPTIONS")

	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
	r.Use(middleware.CreateCheckAdminMiddleware(userClient))
}
