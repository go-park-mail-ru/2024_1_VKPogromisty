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
	r := rootRouter.PathPrefix("/chat").Subrouter()

	sanitizer := sanitizer.NewSanitizer(bluemonday.UGCPolicy())

	h := rest.NewChatServer(pubSubRepo, messagesRepo, sanitizer)

	r.HandleFunc("/", h.ServeWS).Methods("GET", "OPTIONS")
	r.HandleFunc("/dialogs", h.HandleGetDialogs).Methods("GET", "OPTIONS")
	r.HandleFunc("/messages", h.HandleGetMessagesByDialog).Methods("GET", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	r.Use(middleware.CSRFMiddleware)
}
