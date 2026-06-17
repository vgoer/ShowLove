package service

import (
	"context"
	"errors"
	"fmt"

	"showlove/services/comment-service/internal/model"
	"showlove/services/comment-service/internal/repository"
)

var (
	ErrCommentNotFound = errors.New("评论不存在")
	ErrPermissionDenied = errors.New("只能删除自己的评论")
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

type CreateCommentParams struct {
	PostID         string
	AuthorID       string
	AuthorNickname string
	AuthorAvatar   string
	Content        string
	IsAIGenerated  bool
}

func (s *CommentService) CreateComment(ctx context.Context, params CreateCommentParams) (*model.Comment, error) {
	comment := &model.Comment{
		PostID:         params.PostID,
		AuthorID:       params.AuthorID,
		AuthorNickname: params.AuthorNickname,
		AuthorAvatar:   params.AuthorAvatar,
		Content:        params.Content,
		IsAIGenerated:  params.IsAIGenerated,
	}
	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("comment service: create: %w", err)
	}
	return comment, nil
}

func (s *CommentService) ListComments(ctx context.Context, postID string, page, pageSize int32) ([]*model.Comment, int64, error) {
	return s.repo.FindByPostID(ctx, postID, page, pageSize)
}

func (s *CommentService) DeleteComment(ctx context.Context, commentID, userID string) error {
	err := s.repo.Delete(ctx, commentID, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrCommentNotFound
	}
	return err
}
