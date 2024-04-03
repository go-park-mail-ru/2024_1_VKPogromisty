package routers

import (
	rest "socio/internal/rest/chat"
	"socio/internal/rest/middleware"
	"socio/usecase/auth"
	"socio/usecase/chat"

	"github.com/gorilla/mux"
)

func MountChatRouter(rootRouter *mux.Router, pubSubRepo chat.PubSubRepository, messagesRepo chat.PersonalMessagesRepository, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/chat").Subrouter()
	h := rest.NewChatServer(pubSubRepo, messagesRepo)

	r.HandleFunc("/", h.ServeWS).Methods("GET", "OPTIONS")
	r.HandleFunc("/dialogs", h.HandleGetDialogs).Methods("GET", "OPTIONS")
	r.HandleFunc("/messages", h.HandleGetMessagesByDialog).Methods("GET", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
}
