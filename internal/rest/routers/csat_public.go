package routers

import (
	csatpb "socio/internal/grpc/csat/proto"
	rest "socio/internal/rest/csat"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
)

func MountCSATPublicRouter(rootRouter *mux.Router, CSATClient csatpb.CSATClient, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/csat-public").Subrouter()

	h := rest.NewCSATHandler(CSATClient)

	r.HandleFunc("/pools/", h.GetPools).Methods("GET", "OPTIONS")
	r.HandleFunc("/pools/{poolID:[0-9]+}/questions", h.GetQuestionsByPoolID).Methods("GET", "OPTIONS")
	r.HandleFunc("/pools/{poolID:[0-9]+}/unanswered", h.GetUnansweredQuestionsByPoolID).Methods("GET", "OPTIONS")
	r.HandleFunc("/questions/reply", h.CreateReply).Methods("POST", "OPTIONS")
	r.HandleFunc("/pools/{poolID:[0-9]+}/stats", h.GetStatsByPoolID).Methods("GET", "OPTIONS")

	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
