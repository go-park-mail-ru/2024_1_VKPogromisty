package routers

import (
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/subscriptions"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/csrf"

	uspb "socio/internal/grpc/user/proto"

	"github.com/gorilla/mux"
)

func MountSubscriptionsRouter(rootRouter *mux.Router, userService uspb.UserClient, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/subscriptions").Subrouter()
	h := rest.NewSubscriptionsHandler(userService)

	r.HandleFunc("/", h.HandleSubscription).Methods("POST", "OPTIONS")
	r.HandleFunc("/", h.HandleUnsubscription).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/subscribers", h.HandleGetSubscribers).Methods("GET", "OPTIONS")
	r.HandleFunc("/subscriptions", h.HandleGetSubscriptions).Methods("GET", "OPTIONS")
	r.HandleFunc("/friends", h.HandleGetFriends).Methods("GET", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
