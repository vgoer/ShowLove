package repository

import (
	"context"
	"errors"

	"showlove/services/mood-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNotFound = errors.New("记录不存在")

type MoodRepository interface {
	Upsert(ctx context.Context, entry *model.MoodEntry) error
	FindByUserAndDateRange(ctx context.Context, userID, from, to string) ([]*model.MoodEntry, error)
}

type moodRepository struct{ db *gorm.DB }

func NewMoodRepository(db *gorm.DB) MoodRepository {
	return &moodRepository{db: db}
}

func (r *moodRepository) Upsert(ctx context.Context, entry *model.MoodEntry) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "created_at"}},
		DoUpdates: clause.AssignmentColumns([]string{"mood_level", "mood_label", "note"}),
	}).Create(entry).Error
}

func (r *moodRepository) FindByUserAndDateRange(ctx context.Context, userID, from, to string) ([]*model.MoodEntry, error) {
	var entries []*model.MoodEntry
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, from, to).
		Order("created_at ASC").Find(&entries).Error
	return entries, err
}
