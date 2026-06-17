package repository

import (
	"context"
	"errors"

	"showlove/services/comment-service/internal/model"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("评论不存在")

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	FindByPostID(ctx context.Context, postID string, page, pageSize int32) ([]*model.Comment, int64, error)
	Delete(ctx context.Context, id, authorID string) error
}

type commentRepository struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, c *model.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *commentRepository) FindByPostID(ctx context.Context, postID string, page, pageSize int32) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Comment{}).Where("post_id = ?", postID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").Offset(int(offset)).Limit(int(pageSize)).Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

func (r *commentRepository) Delete(ctx context.Context, id, authorID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND author_id = ?", id, authorID).Delete(&model.Comment{})
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}
