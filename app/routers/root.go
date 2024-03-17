package routers

import (
	repository "socio/internal/repository/map"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"
	"sync"

	"github.com/gorilla/mux"
)

func NewRootRouter() (rootRouter *mux.Router) {
	rootRouter = mux.NewRouter().PathPrefix("/api/v1/").Subrouter()

	userStorage := repository.NewUsers(customtime.RealTimeProvider{}, &sync.Map{})
	sessionStorage := repository.NewSessions(&sync.Map{})
	postStorage := repository.NewPosts(customtime.RealTimeProvider{}, &sync.Map{})

	MountAuthRouter(rootRouter, userStorage, sessionStorage)
	MountPostsRouter(rootRouter, postStorage, userStorage, sessionStorage)
	MountStaticRouter(rootRouter)

	rootRouter.Use(middleware.SetUpCORS)
	rootRouter.Use(middleware.DisableCache)

	return
}
