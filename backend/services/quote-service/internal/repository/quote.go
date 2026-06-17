package repository

import (
	"context"
	"errors"

	"showlove/services/quote-service/internal/model"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("语录不存在")

type QuoteRepository interface {
	Create(ctx context.Context, q *model.DailyQuote) error
	FindByDate(ctx context.Context, date string) (*model.DailyQuote, error)
	List(ctx context.Context, page, pageSize int32) ([]*model.DailyQuote, int64, error)
}

type quoteRepository struct{ db *gorm.DB }

func NewQuoteRepository(db *gorm.DB) QuoteRepository {
	return &quoteRepository{db: db}
}

func (r *quoteRepository) Create(ctx context.Context, q *model.DailyQuote) error {
	return r.db.WithContext(ctx).Create(q).Error
}

func (r *quoteRepository) FindByDate(ctx context.Context, date string) (*model.DailyQuote, error) {
	var q model.DailyQuote
	result := r.db.WithContext(ctx).Where("scheduled_date = ?", date).First(&q)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &q, result.Error
}

func (r *quoteRepository) List(ctx context.Context, page, pageSize int32) ([]*model.DailyQuote, int64, error) {
	var quotes []*model.DailyQuote
	var total int64
	query := r.db.WithContext(ctx).Model(&model.DailyQuote{})
	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Order("scheduled_date DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&quotes).Error
	return quotes, total, err
}
