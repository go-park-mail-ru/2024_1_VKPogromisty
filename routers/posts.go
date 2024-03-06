package routers

import (
	"socio/handlers"

	"github.com/gorilla/mux"
)

func MountPostsRouter(rootRouter *mux.Router, authHandler *handlers.AuthHandler) {
	r := rootRouter.PathPrefix("/posts").Subrouter()
	h := handlers.NewPostsHandler()

	r.HandleFunc("/", h.HandleListPosts).Methods("GET", "OPTIONS")
	r.Use(authHandler.CheckIsAuthorizedMiddleware)
}
