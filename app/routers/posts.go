package routers

import (
	"socio/internal/rest"
	"socio/internal/rest/middleware"
	"socio/usecase/auth"
	"socio/usecase/posts"

	"github.com/gorilla/mux"
)

func MountPostsRouter(rootRouter *mux.Router, postStorage posts.PostsStorage, userStorage posts.UserStorage, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/posts").Subrouter()
	h := rest.NewPostsHandler(postStorage, userStorage)

	r.HandleFunc("/", h.HandleGetUserPosts).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleCreatePost).Methods("POST", "OPTIONS")
	r.HandleFunc("/", h.HandleDeletePost).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
}
