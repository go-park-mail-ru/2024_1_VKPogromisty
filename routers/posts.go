package routers

import (
	"socio/handlers"

	"github.com/gorilla/mux"
)

func MountPostsRouter(rootRouter *mux.Router, authHandler *handlers.AuthHandler) {
	r := rootRouter.PathPrefix("/posts").Subrouter()
	h := handlers.NewPostsHandler()

	r.Use(authHandler.CheckIsAuthorized)
	r.HandleFunc("/", h.HandleListPosts).Methods("GET")
}
