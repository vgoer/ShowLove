package main

import (
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

	dbDSN := getEnv("DB_DSN", "postgres://showlove:showlove123@localhost:5432/notifications_db?sslmode=disable")
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
