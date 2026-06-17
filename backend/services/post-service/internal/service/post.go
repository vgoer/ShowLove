// Package service implements business logic for the post service.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"showlove/pkg/events"
	"showlove/services/post-service/internal/model"
	"showlove/services/post-service/internal/moderation"
	"showlove/services/post-service/internal/repository"
)

var (
	ErrSensitiveContent = errors.New("内容包含敏感词，请修改后重新发布")
	ErrPermissionDenied = errors.New("只能删除自己的帖子")
	ErrPostNotFound     = errors.New("帖子不存在")
)

// PostService handles post business logic.
type PostService struct {
	repo       repository.PostRepository
	filter     *moderation.Filter
	eventPub   events.Publisher
}

// NewPostService creates a new PostService.
func NewPostService(repo repository.PostRepository, filter *moderation.Filter, eventPub events.Publisher) *PostService {
	return &PostService{
		repo:     repo,
		filter:   filter,
		eventPub: eventPub,
	}
}

// CreatePostParams contains the input for creating a post.
type CreatePostParams struct {
	AuthorID       string
	AuthorNickname string
	AuthorAvatar   string
	Content        string
	MoodTag        string
	Images         []string
	VoiceURL       string
}

// CreatePost creates a new post after sensitive word check.
func (s *PostService) CreatePost(ctx context.Context, params CreatePostParams) (*model.Post, error) {
	// Check sensitive words
	if s.filter.Contains(params.Content) {
		found := s.filter.FindAll(params.Content)
		return nil, fmt.Errorf("%w: %v", ErrSensitiveContent, found)
	}

	imagesJSON, _ := json.Marshal(params.Images)

	post := &model.Post{
		AuthorID:       params.AuthorID,
		AuthorNickname: params.AuthorNickname,
		AuthorAvatar:   params.AuthorAvatar,
		Content:        params.Content,
		MoodTag:        params.MoodTag,
		Images:         string(imagesJSON),
		VoiceURL:       params.VoiceURL,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("post service: create post: %w", err)
	}

	// Publish event for AI service
	payload, _ := json.Marshal(map[string]string{
		"post_id":         post.ID,
		"content":         post.Content,
		"mood_tag":        post.MoodTag,
		"author_nickname": post.AuthorNickname,
	})
	if s.eventPub != nil {
		_ = s.eventPub.Publish(ctx, "post.created", events.Event{
			Type:    "post.created",
			Payload: payload,
		})
	}

	return post, nil
}

// GetPost retrieves a post by ID.
func (s *PostService) GetPost(ctx context.Context, postID string) (*model.Post, error) {
	post, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("post service: get post %s: %w", postID, err)
	}
	return post, nil
}

// ListPosts retrieves paginated posts.
func (s *PostService) ListPosts(ctx context.Context, sort string, page, pageSize int32) ([]*model.Post, int64, error) {
	return s.repo.List(ctx, sort, page, pageSize)
}

// DeletePost deletes a post (author only).
func (s *PostService) DeletePost(ctx context.Context, postID, userID string) error {
	err := s.repo.Delete(ctx, postID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrPostNotFound
		}
		return fmt.Errorf("post service: delete post: %w", err)
	}
	return nil
}

// SendSticker increments the sticker count for a post.
func (s *PostService) SendSticker(ctx context.Context, postID, stickerType string) (*model.Post, error) {
	post, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("post service: find post: %w", err)
	}

	switch stickerType {
	case "hug":
		post.StickerHug++
	case "cheer":
		post.StickerCheer++
	case "understand":
		post.StickerUnderstand++
	}

	if err := s.repo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("post service: update sticker: %w", err)
	}

	return post, nil
}

// ReportPost marks a post as reported.
func (s *PostService) ReportPost(ctx context.Context, postID, userID, reason string) error {
	post, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrPostNotFound
		}
		return err
	}

	post.IsReported = true
	// TODO: store report reason in a separate reports table

	return s.repo.Update(ctx, post)
}
