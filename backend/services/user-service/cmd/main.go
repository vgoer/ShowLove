// Package main is the entry point for the user service.
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"showlove/pkg/jwt"
	"showlove/services/user-service/internal/repository"
	"showlove/services/user-service/internal/server"
	"showlove/services/user-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[user-service] Starting...")

	// Configuration from environment
	dbDSN := getEnv("DB_DSN", "postgres://showlove:showlove123@localhost:5432/users_db?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-in-production")
	grpcPort := getEnv("GRPC_PORT", "50051")

	// Database
	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("[user-service] Failed to connect to database: %v", err)
	}
	log.Println("[user-service] Database connected")

	// Repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// JWT Manager
	jwtMgr := jwt.NewManager(jwtSecret, 15*time.Minute, 7*24*time.Hour)

	// Service
	svc := service.NewUserService(userRepo, tokenRepo, jwtMgr)

	// gRPC Server
	grpcServer := grpc.NewServer()
	userServer := server.NewUserServer(db, svc)
	userServer.RegisterGRPC(grpcServer)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("[user-service] Failed to listen: %v", err)
	}

	log.Printf("[user-service] gRPC server listening on :%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[user-service] Failed to serve: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
