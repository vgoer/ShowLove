package service

import (
	"context"
	"testing"
	"time"

	"showlove/services/quote-service/internal/model"
	"showlove/services/quote-service/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockQuoteRepo struct {
	quotes map[string]*model.DailyQuote
}

func newMockQuoteRepo() *mockQuoteRepo {
	return &mockQuoteRepo{quotes: make(map[string]*model.DailyQuote)}
}

func (m *mockQuoteRepo) Create(_ context.Context, q *model.DailyQuote) error {
	m.quotes[q.ScheduledDate] = q
	return nil
}

func (m *mockQuoteRepo) FindByDate(_ context.Context, date string) (*model.DailyQuote, error) {
	q, ok := m.quotes[date]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return q, nil
}

func (m *mockQuoteRepo) List(_ context.Context, _, _ int32) ([]*model.DailyQuote, int64, error) {
	var result []*model.DailyQuote
	for _, q := range m.quotes {
		result = append(result, q)
	}
	return result, int64(len(result)), nil
}

func TestGetTodayQuote(t *testing.T) {
	repo := newMockQuoteRepo()
	svc := NewQuoteService(repo)

	today := time.Now().Format("2006-01-02")
	_, _ = svc.CreateQuote(context.Background(), "你好世界", "Hello World", "测试", today)

	quote, err := svc.GetTodayQuote(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "你好世界", quote.TextZH)
}

func TestSeedQuotes(t *testing.T) {
	quotes := SeedQuotes()
	assert.Len(t, quotes, 30)
	// All quotes should have content
	for _, q := range quotes {
		assert.NotEmpty(t, q.TextZH)
		assert.NotEmpty(t, q.TextEN)
		assert.NotEmpty(t, q.ScheduledDate)
	}
}
