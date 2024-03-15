package routers

import (
	repository "socio/internal/repository/map"
	"socio/internal/rest"
	"socio/internal/rest/middleware"

	"github.com/gorilla/mux"
)

func MountPostsRouter(rootRouter *mux.Router, postStorage *repository.Posts, userStorage *repository.Users, sessionStorage *repository.Sessions) {
	r := rootRouter.PathPrefix("/posts").Subrouter()
	h := rest.NewPostsHandler(postStorage, userStorage)

	r.HandleFunc("/", h.HandleListPosts).Methods("GET", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
}
