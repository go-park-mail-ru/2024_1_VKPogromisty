package routers

import (
	"fmt"
	"net/http"
	"os"
	postpb "socio/internal/grpc/post/proto"
	uspb "socio/internal/grpc/user/proto"
	pgRepo "socio/internal/repository/postgres"
	redisRepo "socio/internal/repository/redis"
	"socio/internal/rest/middleware"
	customtime "socio/pkg/time"

	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	DotenvPath    = "../../.env"
	AppPortEnv    = "APP_PORT"
	PgUserEnv     = "PG_USER"
	PgDBNameEnv   = "PG_DBNAME"
	PgPasswordEnv = "PG_PASSWORD"
	PgHostEnv     = "PG_HOST"
	PgPortEnv     = "PG_PORT"
)

func MountRootRouter(router *mux.Router) (err error) {
	if err = godotenv.Load(DotenvPath); err != nil {
		return
	}
	rootRouter := router.PathPrefix("/api/v1/").Subrouter()

	pgConnStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv(PgUserEnv),
		os.Getenv(PgDBNameEnv),
		os.Getenv(PgPasswordEnv),
		os.Getenv(PgHostEnv),
		os.Getenv(PgPortEnv),
	)
	db, err := pgRepo.NewPool(pgConnStr)
	if err != nil {
		return
	}
	defer db.Close()

	userStorage := pgRepo.NewUsers(db, customtime.RealTimeProvider{})
	personalMessageStorage := pgRepo.NewPersonalMessages(db, customtime.RealTimeProvider{})

	redisPool := redisRepo.NewPool(os.Getenv("REDIS_PROTOCOL"), os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))
	defer redisPool.Close()

	sessionStorage := redisRepo.NewSession(redisPool)
	chatPubSubRepository := redisRepo.NewChatPubSub(redisPool)

	userClientConn, err := grpc.Dial(
		os.Getenv("GRPC_USER_SERVICE_HOST")+os.Getenv("GRPC_USER_SERVICE_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer userClientConn.Close()

	userClient := uspb.NewUserClient(userClientConn)

	postClientConn, err := grpc.Dial(
		os.Getenv("GRPC_POST_SERVICE_HOST")+os.Getenv("GRPC_POST_SERVICE_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer postClientConn.Close()

	postClient := postpb.NewPostClient(postClientConn)

	MountAuthRouter(rootRouter, userStorage, sessionStorage)
	MountCSRFRouter(rootRouter, sessionStorage)
	MountChatRouter(rootRouter, chatPubSubRepository, personalMessageStorage, sessionStorage)
	MountProfileRouter(rootRouter, userClient, sessionStorage)
	MountPostsRouter(rootRouter, postClient, userClient, sessionStorage)
	MountSubscriptionsRouter(rootRouter, userClient, sessionStorage)
	MountStaticRouter(rootRouter)
	MountAdminRouter(rootRouter, userClient, sessionStorage)

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

	appPort := os.Getenv(AppPortEnv)
	fmt.Printf("started on port %s\n", appPort)
	http.ListenAndServe(appPort, handler)

	return
}
