package service

import (
	"context"
	"testing"

	"showlove/services/ai-service/internal/provider"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockProvider implements Provider for testing.
type mockProvider struct{ name string }

func (m *mockProvider) Generate(_ context.Context, _, _ string) (string, error) {
	return "加油！一切都会好起来的 🌸 —— 小暖", nil
}
func (m *mockProvider) Name() string { return m.name }

func TestGenerateReply(t *testing.T) {
	providers := map[string]provider.Provider{
		"openai": &mockProvider{name: "openai"},
	}
	svc := NewAIService("openai", providers, nil)

	reply, providerName, err := svc.GenerateReply(context.Background(),
		"今天工作压力好大", "stressed", "测试用户")

	require.NoError(t, err)
	assert.Equal(t, "openai", providerName)
	assert.Contains(t, reply, "小暖")
	assert.NotEmpty(t, reply)
}

func TestGenerateReply_UnknownProvider(t *testing.T) {
	svc := NewAIService("unknown", nil, nil)
	_, _, err := svc.GenerateReply(context.Background(), "test", "sad", "user")
	assert.Error(t, err)
}
