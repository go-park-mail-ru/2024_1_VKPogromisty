package routers

import (
	"fmt"
	"net/http"
	"os"
	pgRepo "socio/internal/repository/postgres"
	redisRepo "socio/internal/repository/redis"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	DotenvPath = "../.env"
)

func MountRootRouter(router *mux.Router) (err error) {
	if err = godotenv.Load(DotenvPath); err != nil {
		return
	}
	rootRouter := router.PathPrefix("/api/v1/").Subrouter()

	pgConnStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", os.Getenv("PG_USER"), os.Getenv("PG_DBNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"))
	db, err := pgRepo.NewPool(pgConnStr)
	if err != nil {
		return
	}
	defer db.Close()

	userStorage := pgRepo.NewUsers(db, customtime.RealTimeProvider{})
	postStorage := pgRepo.NewPosts(db, customtime.RealTimeProvider{})
	subStorage := pgRepo.NewSubscriptions(db, customtime.RealTimeProvider{})
	personalMessageStorage := pgRepo.NewPersonalMessages(db, customtime.RealTimeProvider{})

	redisPool := redisRepo.NewPool(os.Getenv("REDIS_PROTOCOL"), os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))
	defer redisPool.Close()

	sessionStorage := redisRepo.NewSession(redisPool)
	chatPubSubRepository := redisRepo.NewChatPubSub(redisPool)

	MountAuthRouter(rootRouter, userStorage, sessionStorage)
	MountCSRFRouter(rootRouter, sessionStorage)
	MountChatRouter(rootRouter, chatPubSubRepository, personalMessageStorage, sessionStorage)
	MountProfileRouter(rootRouter, userStorage, sessionStorage)
	MountPostsRouter(rootRouter, postStorage, userStorage, sessionStorage)
	MountSubscriptionsRouter(rootRouter, subStorage, userStorage, sessionStorage)
	MountStaticRouter(rootRouter)

	prodLogger, err := middleware.NewZapLogger()
	if err != nil {
		return
	}

	defer prodLogger.Sync()

	logger := middleware.NewLogger(prodLogger)

	rootRouter.Use(middleware.Recovery)
	rootRouter.Use(logger.LoggerMiddleware)
	rootRouter.Use(middleware.DisableCache)

	handler := cors.New(cors.Options{
		AllowedOrigins:   middleware.ALLOWED_ORIGINS,
		AllowedMethods:   middleware.ALLOWED_METHODS,
		AllowedHeaders:   middleware.ALLOWED_HEADERS,
		AllowCredentials: true,
	}).Handler(rootRouter)

	appPort := os.Getenv("APP_PORT")
	fmt.Printf("started on port %s\n", appPort)
	http.ListenAndServe(appPort, handler)

	return
}
