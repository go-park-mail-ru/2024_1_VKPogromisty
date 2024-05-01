package routers

import (
	"fmt"
	"net/http"
	"os"
	authpb "socio/internal/grpc/auth/proto"
	postpb "socio/internal/grpc/post/proto"
	pgpb "socio/internal/grpc/public_group/proto"
	uspb "socio/internal/grpc/user/proto"
	pgRepo "socio/internal/repository/postgres"
	redisRepo "socio/internal/repository/redis"
	"socio/internal/rest/middleware"
	"socio/pkg/logger"
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

	personalMessageStorage := pgRepo.NewPersonalMessages(db, customtime.RealTimeProvider{})

	redisPool := redisRepo.NewPool(os.Getenv("REDIS_PROTOCOL"), os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))
	defer redisPool.Close()

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

	authClientConn, err := grpc.Dial(
		os.Getenv("GRPC_AUTH_SERVICE_HOST")+os.Getenv("GRPC_AUTH_SERVICE_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer authClientConn.Close()

	authClient := authpb.NewAuthClient(authClientConn)

	publicGroupClientConn, err := grpc.Dial(
		os.Getenv("GRPC_PUBLIC_GROUP_SERVICE_HOST")+os.Getenv("GRPC_PUBLIC_GROUP_SERVICE_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer publicGroupClientConn.Close()

	publicGroupClient := pgpb.NewPublicGroupClient(publicGroupClientConn)

	MountAuthRouter(rootRouter, authClient, userClient)
	MountCSRFRouter(rootRouter, authClient)
	MountChatRouter(rootRouter, chatPubSubRepository, personalMessageStorage, authClient)
	MountProfileRouter(rootRouter, userClient, authClient)
	MountPostsRouter(rootRouter, postClient, userClient, authClient)
	MountSubscriptionsRouter(rootRouter, userClient, authClient)
	MountPublicGroupRouter(rootRouter, publicGroupClient, postClient, authClient)

	prodLogger, err := logger.NewZapLogger()
	if err != nil {
		return
	}

	defer prodLogger.Sync()

	logger := logger.NewLogger(prodLogger)

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
