package routers

import (
	"socio/handlers"

	"github.com/gorilla/mux"
)

func NewRootRouter() (rootRouter *mux.Router) {
	rootRouter = mux.NewRouter()

	// need auth handler in post router to check if user is authenticated, will be removed when db is added
	authHandler := handlers.NewAuthHandler()
	MountAuthRouter(rootRouter, authHandler)

	MountPostsRouter(rootRouter, authHandler)
	MountStaticRouter(rootRouter)

	return
}
