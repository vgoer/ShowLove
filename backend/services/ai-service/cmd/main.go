package main

import (
	"log"
	"os"

	"showlove/pkg/events"
	"showlove/services/ai-service/internal/provider"
	"showlove/services/ai-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[ai-service] Starting...")

	activeProvider := getEnv("AI_PROVIDER", "openai")

	providers := map[string]provider.Provider{
		"openai": provider.NewOpenAIProvider(provider.OpenAIConfig{
			Config: provider.Config{
				APIKey:  getEnv("OPENAI_API_KEY", ""),
				BaseURL: getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
				Model:   "gpt-4o-mini",
			},
		}),
		"deepseek": provider.NewDeepSeekProvider(provider.Config{
			APIKey:  getEnv("DEEPSEEK_API_KEY", ""),
			BaseURL: "https://api.deepseek.com/v1",
			Model:   "deepseek-chat",
		}),
	}

	sub := events.NewNoOpPubSub()
	svc := service.NewAIService(activeProvider, providers, sub)

	_ = svc // Will be wired with gRPC server and NATS subscription

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	log.Printf("[ai-service] Provider: %s, waiting for events...", activeProvider)
	select {} // Block forever, processing NATS events
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
