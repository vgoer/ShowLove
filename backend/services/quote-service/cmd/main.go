package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"showlove/services/quote-service/internal/model"
	"showlove/services/quote-service/internal/repository"
	"showlove/services/quote-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[quote-service] Starting...")

	dbDSN := getEnv("DB_DSN", buildDSN("quotes_db"))
	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalf("[quote-service] DB: %v", err)
	}

	// Auto-migrate and seed
	db.AutoMigrate(&model.DailyQuote{})
	repo := repository.NewQuoteRepository(db)

	// Seed quotes if empty
	var count int64
	db.Table("daily_quotes").Count(&count)
	if count == 0 {
		log.Println("[quote-service] Seeding 30 quotes...")
		for _, q := range service.SeedQuotes() {
			repo.Create(nil, &q) // context.Background() would be better but this is init
		}
		log.Println("[quote-service] Seeding complete")
	}

	svc := service.NewQuoteService(repo)
	_ = svc

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	port := getEnv("GRPC_PORT", "50055")
	lis, _ := net.Listen("tcp", ":"+port)
	log.Printf("[quote-service] gRPC on :%s", port)
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
