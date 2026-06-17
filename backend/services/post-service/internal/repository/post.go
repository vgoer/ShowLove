package repository

import (
	"context"
	"errors"

	"showlove/services/post-service/internal/model"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("帖子不存在")
)

// PostRepository defines the data access interface for posts.
type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	FindByID(ctx context.Context, id string) (*model.Post, error)
	List(ctx context.Context, sort string, page, pageSize int32) ([]*model.Post, int64, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id, authorID string) error
}

// postRepository implements PostRepository with GORM.
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new PostRepository.
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) FindByID(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&post)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

func (r *postRepository) List(ctx context.Context, sort string, page, pageSize int32) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Post{}).Where("is_hidden = false")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	orderClause := "created_at DESC"
	if sort == "most_helped" {
		orderClause = "comment_count DESC"
	}

	if err := query.Order(orderClause).Offset(int(offset)).Limit(int(pageSize)).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) Update(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *postRepository) Delete(ctx context.Context, id, authorID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND author_id = ?", id, authorID).Delete(&model.Post{})
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}
