package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"showlove/services/user-service/internal/model"
)

// TokenRepository defines the data access interface for refresh tokens.
type TokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, tokenID string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}

// tokenRepository implements TokenRepository with GORM.
type tokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository creates a new TokenRepository.
func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	result := r.db.WithContext(ctx).
		Where("token = ? AND revoked = false AND expires_at > ?", token, time.Now()).
		First(&rt)
	if result.Error != nil {
		return nil, ErrNotFound
	}
	return &rt, nil
}

func (r *tokenRepository) Revoke(ctx context.Context, tokenID string) error {
	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked", true).Error
}

func (r *tokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}
