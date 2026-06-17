// Package service implements the AI reply generation service.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"showlove/pkg/events"
	"showlove/services/ai-service/internal/prompt"
	"showlove/services/ai-service/internal/provider"
)

// AIService handles AI reply generation triggered by post creation events.
type AIService struct {
	providers map[string]provider.Provider
	active    string
	sub       events.Subscriber
}

// NewAIService creates a new AI service with configured providers.
func NewAIService(activeProvider string, providers map[string]provider.Provider, sub events.Subscriber) *AIService {
	return &AIService{
		providers: providers,
		active:    activeProvider,
		sub:       sub,
	}
}

// PostCreatedEvent represents the payload of a post.created event.
type PostCreatedEvent struct {
	PostID         string `json:"post_id"`
	Content        string `json:"content"`
	MoodTag        string `json:"mood_tag"`
	AuthorNickname string `json:"author_nickname"`
}

// Start begins listening for post.created events.
func (s *AIService) Start(ctx context.Context) error {
	return s.sub.Subscribe(ctx, "post.created", func(e events.Event) error {
		var pce PostCreatedEvent
		if err := json.Unmarshal(e.Payload, &pce); err != nil {
			return fmt.Errorf("ai service: parse event: %w", err)
		}

		reply, err := s.generateReply(ctx, pce)
		if err != nil {
			log.Printf("[ai-service] Failed to generate reply for post %s: %v", pce.PostID, err)
			return nil // Don't block event processing on AI failure
		}

		log.Printf("[ai-service] Generated reply for post %s: %s...", pce.PostID, truncate(reply, 50))
		// TODO: Call comment-service gRPC to create AI comment
		return nil
	})
}

func (s *AIService) generateReply(ctx context.Context, pce PostCreatedEvent) (string, error) {
	p, ok := s.providers[s.active]
	if !ok {
		return "", fmt.Errorf("ai service: unknown provider %q", s.active)
	}

	userPrompt := prompt.BuildUserPrompt(pce.Content, pce.MoodTag, pce.AuthorNickname)
	return p.Generate(ctx, prompt.SystemPromptZH, userPrompt)
}

// GenerateReply generates a reply for a post directly (for testing/API).
func (s *AIService) GenerateReply(ctx context.Context, content, moodTag, authorNickname string) (string, string, error) {
	p, ok := s.providers[s.active]
	if !ok {
		return "", "", fmt.Errorf("ai service: unknown provider %q", s.active)
	}

	userPrompt := prompt.BuildUserPrompt(content, moodTag, authorNickname)
	reply, err := p.Generate(ctx, prompt.SystemPromptZH, userPrompt)
	if err != nil {
		return "", "", err
	}
	return reply, p.Name(), nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
