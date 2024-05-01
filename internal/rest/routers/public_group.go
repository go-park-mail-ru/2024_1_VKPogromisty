package routers

import (
	authpb "socio/internal/grpc/auth/proto"
	pgpb "socio/internal/grpc/public_group/proto"
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/public_group"
	customtime "socio/pkg/time"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
)

func MountPublicGroupRouter(rootRouter *mux.Router, groupClient pgpb.PublicGroupClient, authManager authpb.AuthClient) {
	r := rootRouter.PathPrefix("/groups").Subrouter()

	h := rest.NewPublicGroupHandler(groupClient)

	r.HandleFunc("/search", h.HandleSearchByName).Methods("GET", "OPTIONS")
	r.HandleFunc("/{groupID:[0-9]+}", h.HandleGetByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/by-sub/{userID:[0-9]+}", h.HandleGetBySubscriberID).Methods("GET", "OPTIONS")
	r.HandleFunc("/{groupID:[0-9]+}/is-sub", h.HandleGetSubscriptionByPublicGroupIDAndSubscriberID).Methods("GET", "OPTIONS")
	r.HandleFunc("/{groupID:[0-9]+}/sub", h.HandleSubscribe).Methods("POST", "OPTIONS")
	r.HandleFunc("/{groupID:[0-9]+}/unsub", h.HandleUnsubscribe).Methods("POST", "OPTIONS")
	r.HandleFunc("/", h.HandleCreate).Methods("POST", "OPTIONS")
	r.HandleFunc("/{groupID:[0-9]+}", h.HandleUpdate).Methods("PUT", "OPTIONS")
	r.HandleFunc("/{groupID:[0-9]+}", h.HandleDelete).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(authManager))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
