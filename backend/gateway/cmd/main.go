// Package main is the entry point for the API Gateway.
package main

import (
	"log"
	"os"
	"time"

	"showlove/gateway/internal/router"
	"showlove/pkg/jwt"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[gateway] Starting...")

	// Configuration
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-in-production")
	gatewayPort := getEnv("GATEWAY_PORT", "8080")

	// JWT Manager
	jwtMgr := jwt.NewManager(jwtSecret, 15*time.Minute, 7*24*time.Hour)

	// Setup router
	r := router.Setup(jwtMgr)

	log.Printf("[gateway] HTTP server listening on :%s", gatewayPort)
	if err := r.Run(":" + gatewayPort); err != nil {
		log.Fatalf("[gateway] Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
