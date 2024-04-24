package routers

import (
	post "socio/internal/grpc/post/proto"
	user "socio/internal/grpc/user/proto"
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/posts"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/csrf"

	"github.com/gorilla/mux"
)

func MountPostsRouter(rootRouter *mux.Router, postsClient post.PostClient, userClient user.UserClient, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/posts").Subrouter()

	h := rest.NewPostsHandler(postsClient, userClient)

	r.HandleFunc("/{postID:[0-9]+}", h.HandleGetPostByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleGetUserPosts).Methods("GET", "OPTIONS")
	r.HandleFunc("/friends", h.HandleGetUserFriendsPosts).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleCreatePost).Methods("POST", "OPTIONS")
	r.HandleFunc("/", h.HandleUpdatePost).Methods("PUT", "OPTIONS")
	r.HandleFunc("/", h.HandleDeletePost).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/like", h.HandleLikePost).Methods("POST", "OPTIONS")
	r.HandleFunc("/unlike", h.HandleUnlikePost).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
