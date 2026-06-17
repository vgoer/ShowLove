package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"showlove/services/mood-service/internal/repository"
	"showlove/services/mood-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[mood-service] Starting...")

	dbDSN := getEnv("DB_DSN", buildDSN("moods_db"))
	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalf("[mood-service] DB: %v", err)
	}

	repo := repository.NewMoodRepository(db)
	svc := service.NewMoodService(repo)
	_ = svc

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	port := getEnv("GRPC_PORT", "50054")
	lis, _ := net.Listen("tcp", ":"+port)
	log.Printf("[mood-service] gRPC on :%s", port)
	grpcServer.Serve(lis)
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
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
