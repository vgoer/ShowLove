// Package main is the entry point for the post service.
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"showlove/pkg/events"
	"showlove/services/post-service/internal/moderation"
	"showlove/services/post-service/internal/repository"
	"showlove/services/post-service/internal/server"
	"showlove/services/post-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[post-service] Starting...")

	dbDSN := getEnv("DB_DSN", buildDSN("posts_db"))
	grpcPort := getEnv("GRPC_PORT", "50052")

	// Database
	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("[post-service] Database connection failed: %v", err)
	}
	log.Println("[post-service] Database connected")

	// Repository
	postRepo := repository.NewPostRepository(db)

	// Sensitive word filter
	filter := moderation.NewFilter(moderation.DefaultChineseWords())

	// Event publisher (no-op for now, full NATS integration later)
	eventPub := events.NewNoOpPubSub()

	// Service
	svc := service.NewPostService(postRepo, filter, eventPub)

	// gRPC Server
	grpcServer := grpc.NewServer()
	postServer := server.NewPostServer(svc)
	postServer.RegisterGRPC(grpcServer)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("[post-service] Listen failed: %v", err)
	}

	log.Printf("[post-service] gRPC server on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[post-service] Serve failed: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
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
