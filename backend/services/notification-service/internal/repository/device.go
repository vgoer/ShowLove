package repository

import (
	"context"

	"showlove/services/notification-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceRepository interface {
	Upsert(ctx context.Context, d *model.DeviceToken) error
	FindByUserID(ctx context.Context, userID string) ([]*model.DeviceToken, error)
}

type deviceRepository struct{ db *gorm.DB }

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Upsert(ctx context.Context, d *model.DeviceToken) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "token"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "platform"}),
	}).Create(d).Error
}

func (r *deviceRepository) FindByUserID(ctx context.Context, userID string) ([]*model.DeviceToken, error) {
	var tokens []*model.DeviceToken
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}
