package routers

import (
	"socio/internal/rest"
	"socio/internal/rest/middleware"
	"socio/usecase/auth"
	"socio/usecase/posts"

	"github.com/gorilla/mux"
)

func MountPostsRouter(rootRouter *mux.Router, postStorage posts.PostsStorage, userStorage posts.UsersStorage, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/posts").Subrouter()
	h := rest.NewPostsHandler(postStorage, userStorage)

	r.HandleFunc("/", h.HandleListPosts).Methods("GET", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
}
