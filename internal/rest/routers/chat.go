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

func MountChatRouter(rootRouter *mux.Router, pubSubRepo chat.PubSubRepository, unsentMessageAttachmentsStorage chat.UnsentMessageAttachmentsStorage, messagesRepo chat.PersonalMessagesRepository, authManager authpb.AuthClient, stickerStorage chat.StickerStorage, messageAttachmentStorage chat.MessageAttachmentStorage) {
	sanitizer := sanitizer.NewSanitizer(bluemonday.UGCPolicy())
	h := rest.NewChatServer(pubSubRepo, unsentMessageAttachmentsStorage, messagesRepo, stickerStorage, messageAttachmentStorage, sanitizer)

	csrfFreeRouter := rootRouter.PathPrefix("/chat/ws").Subrouter()
	csrfFreeRouter.HandleFunc("/", h.ServeWS).Methods("GET", "OPTIONS")
	csrfFreeRouter.Use(middleware.CreateCheckIsAuthorizedMiddleware(authManager))

	csrfRequiredRouter := rootRouter.PathPrefix("/chat").Subrouter()
	csrfRequiredRouter.Use(middleware.CreateCheckIsAuthorizedMiddleware(authManager))
	csrfRequiredRouter.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))

	csrfRequiredRouter.HandleFunc("/dialogs", h.HandleGetDialogs).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/dialogs/{receiverID:[0-9]+}/unsent-attachments/", h.HandleGetUnsentMessageAttachments).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/dialogs/{receiverID:[0-9]+}/unsent-attachments/", h.HandleCreateUnsentMessageAttachments).Methods("POST", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/dialogs/{receiverID:[0-9]+}/unsent-attachments/", h.HandleDeleteUnsentMessageAttachments).Methods("DELETE", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/dialogs/{receiverID:[0-9]+}/unsent-attachments/{fileName}", h.HandleDeleteUnsentMessageAttachment).Methods("DELETE", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/messages", h.HandleGetMessagesByDialog).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/stickers/", h.HandleGetAllStickers).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/stickers/{authorID:[0-9]+}", h.HandleGetStickersByAuthorID).Methods("GET", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/stickers/", h.HandleCreateSticker).Methods("POST", "OPTIONS")
	csrfRequiredRouter.HandleFunc("/stickers/{stickerID:[0-9]+}", h.HandleDeleteSticker).Methods("DELETE", "OPTIONS")
}
