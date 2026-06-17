package service

import (
	"context"
	"fmt"
	"time"

	"showlove/services/mood-service/internal/model"
	"showlove/services/mood-service/internal/repository"
)

type MoodService struct{ repo repository.MoodRepository }

func NewMoodService(repo repository.MoodRepository) *MoodService {
	return &MoodService{repo: repo}
}

type RecordMoodParams struct {
	UserID    string
	MoodLevel int32
	MoodLabel string
	Note      string
}

func (s *MoodService) RecordMood(ctx context.Context, params RecordMoodParams) (*model.MoodEntry, error) {
	entry := &model.MoodEntry{
		UserID:    params.UserID,
		MoodLevel: params.MoodLevel,
		MoodLabel: params.MoodLabel,
		Note:      params.Note,
		CreatedAt: time.Now().Format("2006-01-02"),
	}
	if err := s.repo.Upsert(ctx, entry); err != nil {
		return nil, fmt.Errorf("mood service: record: %w", err)
	}
	return entry, nil
}

func (s *MoodService) GetMoods(ctx context.Context, userID, from, to string) ([]*model.MoodEntry, error) {
	return s.repo.FindByUserAndDateRange(ctx, userID, from, to)
}

type MoodPoint struct {
	Date      string `json:"date"`
	MoodLevel int32  `json:"mood_level"`
	MoodLabel string `json:"mood_label"`
}

func (s *MoodService) GetWeeklyMood(ctx context.Context, userID string) ([]MoodPoint, error) {
	now := time.Now()
	from := now.AddDate(0, 0, -6).Format("2006-01-02")
	to := now.Format("2006-01-02")

	entries, err := s.repo.FindByUserAndDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	entryMap := make(map[string]*model.MoodEntry)
	for _, e := range entries {
		entryMap[e.CreatedAt] = e
	}

	points := make([]MoodPoint, 7)
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -6+i).Format("2006-01-02")
		if e, ok := entryMap[date]; ok {
			points[i] = MoodPoint{Date: date, MoodLevel: e.MoodLevel, MoodLabel: e.MoodLabel}
		} else {
			points[i] = MoodPoint{Date: date, MoodLevel: 0, MoodLabel: ""}
		}
	}
	return points, nil
}
