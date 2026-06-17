package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"showlove/pkg/events"
	"showlove/services/notification-service/internal/repository"
	"showlove/services/notification-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[notification-service] Starting...")

	dbDSN := getEnv("DB_DSN", buildDSN("notifications_db"))
	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalf("[notification-service] DB: %v", err)
	}

	repo := repository.NewDeviceRepository(db)
	sub := events.NewNoOpPubSub()
	svc := service.NewNotificationService(repo, sub)
	_ = svc

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	port := getEnv("GRPC_PORT", "50057")
	lis, _ := net.Listen("tcp", ":"+port)
	log.Printf("[notification-service] gRPC on :%s", port)
	grpcServer.Serve(lis)
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
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
