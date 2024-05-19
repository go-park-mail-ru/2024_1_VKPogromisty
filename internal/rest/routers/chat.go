package routers

import (
	authpb "socio/internal/grpc/auth/proto"
	rest "socio/internal/rest/chat"
	"socio/internal/rest/middleware"
	"socio/pkg/sanitizer"
	customtime "socio/pkg/time"
	"socio/usecase/chat"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
)

func MountChatRouter(rootRouter *mux.Router, pubSubRepo chat.PubSubRepository, messagesRepo chat.PersonalMessagesRepository, authManager authpb.AuthClient, stickerStorage chat.StickerStorage) {
	sanitizer := sanitizer.NewSanitizer(bluemonday.UGCPolicy())
	h := rest.NewChatServer(pubSubRepo, messagesRepo, stickerStorage, sanitizer)

	csrfFreeRouter := rootRouter.PathPrefix("/chat/ws").Subrouter()
	csrfFreeRouter.HandleFunc("/", h.ServeWS).Methods("GET", "OPTIONS")
	csrfFreeRouter.Use(middleware.CreateCheckIsAuthorizedMiddleware(authManager))

	csrfRequiredRouter := rootRouter.PathPrefix("/chat").Subrouter()
	csrfRequiredRouter.Use(middleware.CreateCheckIsAuthorizedMiddleware(authManager))
	csrfRequiredRouter.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))

	csrfRequiredRouter.HandleFunc("/dialogs", h.HandleGetDialogs).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/messages", h.HandleGetMessagesByDialog).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/stickers/", h.HandleGetAllStickers).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/stickers/{authorID}", h.HandleGetStickersByAuthorID).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/stickers/", h.HandleCreateSticker).Methods("POST", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/stickers/{stickerID}", h.HandleDeleteSticker).Methods("DELETE", "OPTIONS")
}
