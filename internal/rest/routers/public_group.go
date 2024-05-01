package routers

import (
	authpb "socio/internal/grpc/auth/proto"
	postpb "socio/internal/grpc/post/proto"
	pgpb "socio/internal/grpc/public_group/proto"
	uspb "socio/internal/grpc/user/proto"
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/public_group"
	customtime "socio/pkg/time"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
)

func MountPublicGroupRouter(rootRouter *mux.Router, groupClient pgpb.PublicGroupClient, postClient postpb.PostClient, userClient uspb.UserClient, authManager authpb.AuthClient) {
	publicRouter := rootRouter.PathPrefix("/groups").Subrouter()

	h := rest.NewPublicGroupHandler(groupClient, postClient, userClient)

	publicRouter.HandleFunc("/search", h.HandleSearchByName).Methods("GET", "OPTIONS")
	publicRouter.HandleFunc("/{groupID:[0-9]+}", h.HandleGetByID).Methods("GET", "OPTIONS")
	publicRouter.HandleFunc("/by-sub/{userID:[0-9]+}", h.HandleGetBySubscriberID).Methods("GET", "OPTIONS")
	publicRouter.HandleFunc("/{groupID:[0-9]+}/is-sub", h.HandleGetSubscriptionByPublicGroupIDAndSubscriberID).Methods("GET", "OPTIONS")
	publicRouter.HandleFunc("/{groupID:[0-9]+}/sub", h.HandleSubscribe).Methods("POST", "OPTIONS")
	publicRouter.HandleFunc("/{groupID:[0-9]+}/unsub", h.HandleUnsubscribe).Methods("POST", "OPTIONS")
	publicRouter.HandleFunc("/", h.HandleCreate).Methods("POST", "OPTIONS")
	publicRouter.HandleFunc("/{groupID:[0-9]+}/posts/", h.HandleGetGroupPosts).Methods("GET", "OPTIONS")
	publicRouter.Use(middleware.CreateCheckIsAuthorizedMiddleware(authManager))
	publicRouter.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))

	adminRouter := publicRouter.PathPrefix("/groups").Subrouter()

	adminRouter.HandleFunc("/{groupID:[0-9]+}/admins/", h.HandleGetAdminsByPublicGroupID).Methods("GET", "OPTIONS")
	adminRouter.HandleFunc("/{groupID:[0-9]+}/admins/check", h.HandleCheckIfUserIsAdmin).Methods("GET", "OPTIONS")
	adminRouter.HandleFunc("/{groupID:[0-9]+}/admins/", h.HandleCreatePublicGroupAdmin).Methods("POST", "OPTIONS")
	adminRouter.HandleFunc("/{groupID:[0-9]+}/admins/", h.HandleDeletePublicGroupAdmin).Methods("DELETE", "OPTIONS")
	adminRouter.HandleFunc("/{groupID:[0-9]+}", h.HandleUpdate).Methods("PUT", "OPTIONS")
	adminRouter.HandleFunc("/{groupID:[0-9]+}", h.HandleDelete).Methods("DELETE", "OPTIONS")
	adminRouter.HandleFunc("/{groupID:[0-9]+}/posts/", h.HandleCreateGroupPost).Methods("POST", "OPTIONS")
	adminRouter.Use(middleware.CreateCheckIsAuthorizedMiddleware(authManager))
	adminRouter.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
	adminRouter.Use(middleware.CreateCheckPublicGroupAdminMiddleware(userClient))
}
