package routers

import (
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/posts"
	"socio/pkg/sanitizer"
	"socio/usecase/auth"
	"socio/usecase/posts"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
)

func MountPostsRouter(rootRouter *mux.Router, postStorage posts.PostsStorage, userStorage posts.UserStorage, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/posts").Subrouter()

	sanitizer := sanitizer.NewSanitizer(bluemonday.UGCPolicy())

	h := rest.NewPostsHandler(postStorage, userStorage, sanitizer)

	r.HandleFunc("/", h.HandleGetUserPosts).Methods("GET", "OPTIONS")
	r.HandleFunc("/friends", h.HandleGetUserFriendsPosts).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleCreatePost).Methods("POST", "OPTIONS")
	r.HandleFunc("/", h.HandleUpdatePost).Methods("PUT", "OPTIONS")
	r.HandleFunc("/", h.HandleDeletePost).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
}
