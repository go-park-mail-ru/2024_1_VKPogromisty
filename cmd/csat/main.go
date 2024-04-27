package main

import (
	"fmt"
	"net"
	"os"

	"socio/internal/grpc/csat"
	csatpb "socio/internal/grpc/csat/proto"
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

	port := os.Getenv("GRPC_CSAT_SERVICE_PORT")
	lis, err := net.Listen("tcp", "0.0.0.0"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	CSATStorage := pgRepo.NewCSAT(db, customtime.RealTimeProvider{})
	manager := csat.NewCSATManager(CSATStorage)

	server := grpc.NewServer()

	csatpb.RegisterCSATServer(server, manager)

	fmt.Println("CSAT service is running on port:", port)
	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
