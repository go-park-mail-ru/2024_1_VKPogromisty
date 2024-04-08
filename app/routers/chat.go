package routers

import (
	rest "socio/internal/rest/chat"
	"socio/internal/rest/middleware"
	"socio/pkg/sanitizer"
	"socio/usecase/auth"
	"socio/usecase/chat"

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
	csrfRequiredRouter.Use(middleware.CSRFMiddleware)
	csrfRequiredRouter.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))

	csrfRequiredRouter.HandleFunc("/dialogs", h.HandleGetDialogs).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/messages", h.HandleGetMessagesByDialog).Methods("GET", "OPTIONS")
}
