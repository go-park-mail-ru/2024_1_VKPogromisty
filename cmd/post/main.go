package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"socio/internal/grpc/interceptors"
	"socio/internal/grpc/post"

	postspb "socio/internal/grpc/post/proto"
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

	attachmentStorage, err := minioRepo.NewStaticStorage(minioClient, minioRepo.PostAttachmentsBucket)
	if err != nil {
		fmt.Println(err)
		return
	}

	port := os.Getenv("GRPC_POST_SERVICE_PORT")
	lis, err := net.Listen("tcp", "0.0.0.0"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	postsStorage := pgRepo.NewPosts(db, customtime.RealTimeProvider{})
	manager := post.NewPostManager(postsStorage, attachmentStorage)

	prodLogger, err := logger.NewZapLogger(nil)
	if err != nil {
		return
	}

	defer func() {
		err = prodLogger.Sync()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	logger := logger.NewLogger(prodLogger)

	prometheus.MustRegister(
		appmetrics.PostTotalHits,
		appmetrics.PostHits,
		appmetrics.PostHitDuration,
	)

	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(MaxMessageSize),
		grpc.MaxSendMsgSize(MaxMessageSize),
		grpc.ChainUnaryInterceptor(
			logger.UnaryLoggerInterceptor,
			interceptors.PostHitMetricsInterceptor,
			interceptors.UnaryRecoveryInterceptor,
		),
	)

	metricsPort := os.Getenv("GRPC_POST_SERVICE_METRICS_PORT")

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	go func() {
		err = http.ListenAndServe(metricsPort, r)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	fmt.Println("Metrics of post service is running on port:", metricsPort)

	postspb.RegisterPostServer(server, manager)

	fmt.Println("Post service is running on port:", port)
	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
