package routers

import (
	rest "socio/internal/rest/chat"
	"socio/internal/rest/middleware"
	"socio/pkg/sanitizer"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/chat"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
)

func MountChatRouter(rootRouter *mux.Router, pubSubRepo chat.PubSubRepository, messagesRepo chat.PersonalMessagesRepository, sessionStorage auth.SessionStorage) {
	sanitizer := sanitizer.NewSanitizer(bluemonday.UGCPolicy())
	h := rest.NewChatServer(pubSubRepo, messagesRepo, sanitizer)

	csrfFreeRouter := rootRouter.PathPrefix("/chat/ws").Subrouter()
	csrfFreeRouter.HandleFunc("/", h.ServeWS).Methods("GET", "OPTIONS")
	csrfFreeRouter.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))

	csrfRequiredRouter := rootRouter.PathPrefix("/chat").Subrouter()
	csrfRequiredRouter.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	csrfRequiredRouter.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))

	csrfRequiredRouter.HandleFunc("/dialogs", h.HandleGetDialogs).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/messages", h.HandleGetMessagesByDialog).Methods("GET", "OPTIONS")
}
