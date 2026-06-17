package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"showlove/services/mood-service/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockMoodRepo struct {
	entries map[string]*model.MoodEntry
	counter int
}

func newMockMoodRepo() *mockMoodRepo {
	return &mockMoodRepo{entries: make(map[string]*model.MoodEntry)}
}

func (m *mockMoodRepo) Upsert(_ context.Context, e *model.MoodEntry) error {
	m.counter++
	e.ID = fmt.Sprintf("mood-%d", m.counter)
	key := e.UserID + "-" + e.CreatedAt
	m.entries[key] = e
	return nil
}

func (m *mockMoodRepo) FindByUserAndDateRange(_ context.Context, userID, from, to string) ([]*model.MoodEntry, error) {
	var result []*model.MoodEntry
	for _, e := range m.entries {
		if e.UserID == userID && e.CreatedAt >= from && e.CreatedAt <= to {
			result = append(result, e)
		}
	}
	return result, nil
}

func TestRecordMood(t *testing.T) {
	svc := NewMoodService(newMockMoodRepo())
	entry, err := svc.RecordMood(context.Background(), RecordMoodParams{
		UserID: "user-1", MoodLevel: 7, MoodLabel: "平静", Note: "还不错",
	})
	require.NoError(t, err)
	assert.Equal(t, int32(7), entry.MoodLevel)
	assert.Equal(t, "平静", entry.MoodLabel)
	assert.Equal(t, time.Now().Format("2006-01-02"), entry.CreatedAt)
}

func TestRecordMood_Overwrite(t *testing.T) {
	repo := newMockMoodRepo()
	svc := NewMoodService(repo)

	_, _ = svc.RecordMood(context.Background(), RecordMoodParams{
		UserID: "user-1", MoodLevel: 3, MoodLabel: "难过", Note: "早上",
	})
	_, _ = svc.RecordMood(context.Background(), RecordMoodParams{
		UserID: "user-1", MoodLevel: 8, MoodLabel: "开心", Note: "下午变好了",
	})

	entries, err := svc.GetMoods(context.Background(), "user-1", "2000-01-01", "2099-12-31")
	require.NoError(t, err)
	assert.Len(t, entries, 1)
	// The last upsert should have overwritten
	mood := entries[0]
	assert.Equal(t, int32(8), mood.MoodLevel)
}

func TestGetWeeklyMood(t *testing.T) {
	svc := NewMoodService(newMockMoodRepo())
	points, err := svc.GetWeeklyMood(context.Background(), "user-1")
	require.NoError(t, err)
	assert.Len(t, points, 7) // Always returns 7 days
}

