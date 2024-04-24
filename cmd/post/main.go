package main

import (
	"fmt"
	"net"
	"os"
	"socio/internal/grpc/post"

	postspb "socio/internal/grpc/post/proto"
	minioRepo "socio/internal/repository/minio"
	pgRepo "socio/internal/repository/postgres"
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

	attachmentStorage, err := minioRepo.NewStaticStorage(minioClient, minioRepo.AttachmentsBucket)
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

	server := grpc.NewServer()

	postspb.RegisterPostServer(server, manager)

	fmt.Println("Post service is running on port:", port)
	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
