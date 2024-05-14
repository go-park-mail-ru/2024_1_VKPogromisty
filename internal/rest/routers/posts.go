package routers

import (
	authpb "socio/internal/grpc/auth/proto"
	post "socio/internal/grpc/post/proto"
	pgpb "socio/internal/grpc/public_group/proto"
	user "socio/internal/grpc/user/proto"
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/posts"
	customtime "socio/pkg/time"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
)

func MountPostsRouter(rootRouter *mux.Router, postsClient post.PostClient, userClient user.UserClient, publicGroupClient pgpb.PublicGroupClient, authManager authpb.AuthClient) {
	r := rootRouter.PathPrefix("/posts").Subrouter()

	h := rest.NewPostsHandler(postsClient, userClient, publicGroupClient)

	r.HandleFunc("/{postID:[0-9]+}", h.HandleGetPostByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleGetUserPosts).Methods("GET", "OPTIONS")
	r.HandleFunc("/friends", h.HandleGetUserFriendsPosts).Methods("GET", "OPTIONS")
	r.HandleFunc("/groups", h.HandleGetGroupPostsBySubscriptions).Methods("GET", "OPTIONS")
	r.HandleFunc("/all", h.HandleGetPostsByGroupSubIDsAndUserSubIDs).Methods("GET", "OPTIONS")
	r.HandleFunc("/new", h.HandleGetNewPosts).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleCreatePost).Methods("POST", "OPTIONS")
	r.HandleFunc("/", h.HandleUpdatePost).Methods("PUT", "OPTIONS")
	r.HandleFunc("/", h.HandleDeletePost).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/liked", h.HandleGetLikedPosts).Methods("GET", "OPTIONS")
	r.HandleFunc("/like", h.HandleLikePost).Methods("POST", "OPTIONS")
	r.HandleFunc("/unlike", h.HandleUnlikePost).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/{postID:[0-9]+}/comments", h.HandleGetCommentsByPostID).Methods("GET", "OPTIONS")
	r.HandleFunc("/comments", h.HandleCreateComment).Methods("POST", "OPTIONS")
	r.HandleFunc("/comments", h.HandleUpdateComment).Methods("PUT", "OPTIONS")
	r.HandleFunc("/comments", h.HandleDeleteComment).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/comments/like", h.HandleLikeComment).Methods("POST", "OPTIONS")
	r.HandleFunc("/comments/unlike", h.HandleUnlikeComment).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(authManager))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
