package routers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	pgRepo "socio/internal/repository/postgres"
	redisRepo "socio/internal/repository/redis"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"
	"socio/utils"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func MountRootRouter() (err error) {
	if err = godotenv.Load("../.env"); err != nil {
		return
	}
	rootRouter := mux.NewRouter().PathPrefix("/api/v1/").Subrouter()

	pgConnStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", os.Getenv("PG_USER"), os.Getenv("PG_DBNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"))
	db, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		return
	}
	defer db.Close()

	userStorage := pgRepo.NewUsers(db, customtime.RealTimeProvider{})
	postStorage := pgRepo.NewPosts(db, customtime.RealTimeProvider{})
	subStorage := pgRepo.NewSubscriptions(db, customtime.RealTimeProvider{})

	sessionConn, err := redis.Dial(os.Getenv("REDIS_PROTOCOL"), os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"), redis.DialPassword(os.Getenv("REDIS_PASSWORD")))
	if err != nil {
		return
	}
	defer sessionConn.Close()
	sessionStorage := redisRepo.NewSession(sessionConn)

	MountAuthRouter(rootRouter, userStorage, sessionStorage)
	MountPostsRouter(rootRouter, postStorage, userStorage, sessionStorage)
	MountSubscriptionsRouter(rootRouter, subStorage, userStorage, sessionStorage)
	MountStaticRouter(rootRouter)

	rootRouter.Use(middleware.SetUpCORS)
	rootRouter.Use(middleware.DisableCache)

	fmt.Printf("started on port %s\n", utils.PORT)
	http.ListenAndServe(utils.PORT, rootRouter)

	return
}
