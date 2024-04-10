package routers

import (
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/subscriptions"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/csrf"
	"socio/usecase/subscriptions"

	"github.com/gorilla/mux"
)

func MountSubscriptionsRouter(rootRouter *mux.Router, subStorage subscriptions.SubscriptionsStorage, userStorage subscriptions.UserStorage, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/subscriptions").Subrouter()
	h := rest.NewSubscriptionsHandler(subStorage, userStorage)

	r.HandleFunc("/", h.HandleSubscription).Methods("POST", "OPTIONS")
	r.HandleFunc("/", h.HandleUnsubscription).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/subscribers", h.HandleGetSubscribers).Methods("GET", "OPTIONS")
	r.HandleFunc("/subscriptions", h.HandleGetSubscriptions).Methods("GET", "OPTIONS")
	r.HandleFunc("/friends", h.HandleGetFriends).Methods("GET", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
