package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"socio/internal/grpc/interceptors"
	"socio/internal/grpc/user"
	uspb "socio/internal/grpc/user/proto"
	minioRepo "socio/internal/repository/minio"
	pgRepo "socio/internal/repository/postgres"
	"socio/pkg/appmetrics"
	"socio/pkg/logger"
	customtime "socio/pkg/time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	DotenvPath     = "../../.env"
	MaxMessageSize = 1024 * 1024 * 100
)

func main() {
	if err := godotenv.Load(DotenvPath); err != nil {
		fmt.Println(err)
		return
	}

	pgConnStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", os.Getenv("PG_USER"), os.Getenv("PG_DBNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"))
	db, err := pgRepo.NewPool(pgConnStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	minioClient, err := minio.New(os.Getenv("MINIO_HOST"), os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), false)
	if err != nil {
		fmt.Println(err)
		return
	}

	avatarStorage, err := minioRepo.NewStaticStorage(minioClient, minioRepo.UserAvatarsBucket)
	if err != nil {
		fmt.Println(err)
		return
	}

	port := os.Getenv("GRPC_USER_SERVICE_PORT")
	lis, err := net.Listen("tcp", "0.0.0.0"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	userStorage := pgRepo.NewUsers(db, customtime.RealTimeProvider{})
	subsciptionsStorage := pgRepo.NewSubscriptions(db, customtime.RealTimeProvider{})

	manager := user.NewUserManager(userStorage, subsciptionsStorage, avatarStorage)

	prodLogger, err := logger.NewZapLogger(nil)
	if err != nil {
		return
	}

	defer prodLogger.Sync()

	logger := logger.NewLogger(prodLogger)

	prometheus.MustRegister(
		appmetrics.AuthTotalHits,
		appmetrics.AuthHits,
		appmetrics.UserHitDuration,
	)

	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(MaxMessageSize),
		grpc.MaxSendMsgSize(MaxMessageSize),
		grpc.ChainUnaryInterceptor(
			logger.UnaryLoggerInterceptor,
			interceptors.UserHitMetricsInterceptor,
			interceptors.UnaryRecoveryInterceptor,
		),
	)

	metricsPort := os.Getenv("GRPC_USER_SERVICE_METRICS_PORT")

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(metricsPort, r)
	fmt.Println("Metrics of user service is running on port:", metricsPort)

	uspb.RegisterUserServer(server, manager)

	fmt.Println("User service is running on port:", port)
	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
