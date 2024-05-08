package routers

import (
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/subscriptions"
	customtime "socio/pkg/time"
	"socio/usecase/csrf"

	authpb "socio/internal/grpc/auth/proto"
	uspb "socio/internal/grpc/user/proto"

	"github.com/gorilla/mux"
)

func MountSubscriptionsRouter(rootRouter *mux.Router, userClient uspb.UserClient, authClient authpb.AuthClient) {
	r := rootRouter.PathPrefix("/subscriptions").Subrouter()
	h := rest.NewSubscriptionsHandler(userClient)

	r.HandleFunc("/", h.HandleSubscription).Methods("POST", "OPTIONS")
	r.HandleFunc("/", h.HandleUnsubscription).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/subscribers", h.HandleGetSubscribers).Methods("GET", "OPTIONS")
	r.HandleFunc("/subscriptions", h.HandleGetSubscriptions).Methods("GET", "OPTIONS")
	r.HandleFunc("/friends", h.HandleGetFriends).Methods("GET", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(authClient))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
