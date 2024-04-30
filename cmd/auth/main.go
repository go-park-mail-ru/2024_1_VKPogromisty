package main

import (
	"fmt"
	"net"
	"os"
	"socio/internal/grpc/auth"
	"socio/internal/grpc/interceptors"
	"socio/pkg/logger"

	authpb "socio/internal/grpc/auth/proto"
	uspb "socio/internal/grpc/user/proto"
	redisRepo "socio/internal/repository/redis"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	DotenvPath = "../../.env"
)

func main() {
	if err := godotenv.Load(DotenvPath); err != nil {
		fmt.Println(err)
		return
	}

	port := os.Getenv("GRPC_AUTH_SERVICE_PORT")
	lis, err := net.Listen("tcp", "0.0.0.0"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	redisPool := redisRepo.NewPool(os.Getenv("REDIS_PROTOCOL"), os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))
	defer redisPool.Close()

	sessionStorage := redisRepo.NewSession(redisPool)

	userClientConn, err := grpc.Dial(
		os.Getenv("GRPC_USER_SERVICE_HOST")+os.Getenv("GRPC_USER_SERVICE_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer userClientConn.Close()

	userClient := uspb.NewUserClient(userClientConn)

	manager := auth.NewAuthManager(userClient, sessionStorage)

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

	authpb.RegisterAuthServer(server, manager)

	fmt.Println("Auth service is running on port:", port)
	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
