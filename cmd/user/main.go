package main

import (
	"fmt"
	"net"
	"os"
	"socio/internal/grpc/interceptors"
	"socio/internal/grpc/user"
	uspb "socio/internal/grpc/user/proto"
	minioRepo "socio/internal/repository/minio"
	pgRepo "socio/internal/repository/postgres"
	"socio/pkg/logger"
	customtime "socio/pkg/time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go"
	"google.golang.org/grpc"
)

var (
	DotenvPath = "../../.env"
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

	prodLogger, err := logger.NewZapLogger()
	if err != nil {
		return
	}

	defer prodLogger.Sync()

	logger := logger.NewLogger(prodLogger)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logger.UnaryLoggerInterceptor),
		grpc.ChainUnaryInterceptor(interceptors.UnaryRecoveryInterceptor),
	)

	uspb.RegisterUserServer(server, manager)

	fmt.Println("User service is running on port:", port)
	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
