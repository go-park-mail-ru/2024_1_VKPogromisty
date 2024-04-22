package main

import (
	"fmt"
	"net"
	"os"
	"socio/internal/grpc/user"
	uspb "socio/internal/grpc/user/proto"
	pgRepo "socio/internal/repository/postgres"
	customtime "socio/pkg/time"

	"github.com/joho/godotenv"
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

	port := os.Getenv("GRPC_USER_SERVICE_PORT")
	lis, err := net.Listen("tcp", "0.0.0.0"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	userStorage := pgRepo.NewUsers(db, customtime.RealTimeProvider{})
	manager := user.NewUserManager(userStorage)

	server := grpc.NewServer()

	uspb.RegisterUserServer(server, manager)

	fmt.Println("User service is running on port:", port)
	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
