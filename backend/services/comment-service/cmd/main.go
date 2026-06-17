package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"showlove/services/comment-service/internal/repository"
	"showlove/services/comment-service/internal/server"
	"showlove/services/comment-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[comment-service] Starting...")

	dbDSN := getEnv("DB_DSN", buildDSN("comments_db"))
	grpcPort := getEnv("GRPC_PORT", "50053")

	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalf("[comment-service] Database connection failed: %v", err)
	}

	repo := repository.NewCommentRepository(db)
	svc := service.NewCommentService(repo)

	grpcServer := grpc.NewServer()
	commentServer := server.NewCommentServer(svc)
	commentServer.RegisterGRPC(grpcServer)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("[comment-service] Listen failed: %v", err)
	}
	log.Printf("[comment-service] gRPC server on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[comment-service] Serve failed: %v", err)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func buildDSN(dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("POSTGRES_USER", "user"),
		getEnv("POSTGRES_PASSWORD", "password"),
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		dbName,
	)
}
