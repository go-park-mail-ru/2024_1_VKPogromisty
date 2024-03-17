package routers

import (
	"fmt"
	"os"
	mapRepo "socio/internal/repository/map"
	redisRepo "socio/internal/repository/redis"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

func NewRootRouter() (rootRouter *mux.Router, err error) {
	rootRouter = mux.NewRouter().PathPrefix("/api/v1/").Subrouter()

	userStorage := mapRepo.NewUsers(customtime.RealTimeProvider{}, &sync.Map{})
	fmt.Println("here ", os.Getenv("REDIS_PROTOCOL"))
	sessionConn, err := redis.Dial(os.Getenv("REDIS_PROTOCOL"), os.Getenv("REDIS_URL"), redis.DialPassword(os.Getenv("REDIS_PASSWORD")))
	if err != nil {
		return
	}
	sessionStorage := redisRepo.NewSession(sessionConn)

	postStorage := mapRepo.NewPosts(customtime.RealTimeProvider{}, &sync.Map{})

	MountAuthRouter(rootRouter, userStorage, sessionStorage)
	MountPostsRouter(rootRouter, postStorage, userStorage, sessionStorage)
	MountStaticRouter(rootRouter)

	rootRouter.Use(middleware.SetUpCORS)
	rootRouter.Use(middleware.DisableCache)

	return
}
