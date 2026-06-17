// Package provider defines the AI model provider interface.
package provider

import "context"

// Provider is the interface for AI model backends.
type Provider interface {
	// Generate sends a prompt to the AI model and returns the generated text.
	Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error)
	// Name returns the provider's identifier.
	Name() string
}

// Config holds common configuration for AI providers.
type Config struct {
	APIKey  string
	BaseURL string
	Model   string
}
