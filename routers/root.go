package routers

import (
	"github.com/gorilla/mux"
)

func NewRootRouter() (rootRouter *mux.Router) {
	rootRouter = mux.NewRouter()
	MountAuthRouter(rootRouter)
	MountPostsRouter(rootRouter)
	MountStaticRouter(rootRouter)

	return
}
