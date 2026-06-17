package service

import (
	"context"
	"fmt"
	"testing"

	"showlove/services/comment-service/internal/model"
	"showlove/services/comment-service/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCommentRepo struct {
	comments map[string]*model.Comment
	counter  int
}

func newMockRepo() *mockCommentRepo {
	return &mockCommentRepo{comments: make(map[string]*model.Comment)}
}

func (m *mockCommentRepo) Create(_ context.Context, c *model.Comment) error {
	m.counter++
	c.ID = fmt.Sprintf("comment-%d", m.counter)
	m.comments[c.ID] = c
	return nil
}

func (m *mockCommentRepo) FindByPostID(_ context.Context, _ string, page, pageSize int32) ([]*model.Comment, int64, error) {
	var result []*model.Comment
	for _, c := range m.comments {
		result = append(result, c)
	}
	total := int64(len(result))
	offset := (page - 1) * pageSize
	if offset >= int32(len(result)) {
		return []*model.Comment{}, total, nil
	}
	end := offset + pageSize
	if end > int32(len(result)) {
		end = int32(len(result))
	}
	return result[offset:end], total, nil
}

func (m *mockCommentRepo) Delete(_ context.Context, id, authorID string) error {
	c, ok := m.comments[id]
	if !ok || c.AuthorID != authorID {
		return repository.ErrNotFound
	}
	delete(m.comments, id)
	return nil
}

func TestCreateComment_Success(t *testing.T) {
	svc := NewCommentService(newMockRepo())
	c, err := svc.CreateComment(context.Background(), CreateCommentParams{
		PostID: "post-1", AuthorID: "user-1", AuthorNickname: "用户", Content: "暖心评论",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, c.ID)
	assert.Equal(t, "暖心评论", c.Content)
	assert.False(t, c.IsAIGenerated)
}

func TestCreateComment_AIGenerated(t *testing.T) {
	svc := NewCommentService(newMockRepo())
	c, err := svc.CreateComment(context.Background(), CreateCommentParams{
		PostID: "post-1", AuthorID: "ai-bot", AuthorNickname: "小暖", Content: "加油！一切都会好起来的", IsAIGenerated: true,
	})
	require.NoError(t, err)
	assert.True(t, c.IsAIGenerated)
}

func TestListComments(t *testing.T) {
	svc := NewCommentService(newMockRepo())
	for i := 0; i < 3; i++ {
		_, err := svc.CreateComment(context.Background(), CreateCommentParams{
			PostID: "post-1", AuthorID: "user-1", AuthorNickname: "用户", Content: "评论",
		})
		require.NoError(t, err)
	}
	comments, total, err := svc.ListComments(context.Background(), "post-1", 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, comments, 3)
}

func TestDeleteComment_Success(t *testing.T) {
	svc := NewCommentService(newMockRepo())
	c, err := svc.CreateComment(context.Background(), CreateCommentParams{
		PostID: "post-1", AuthorID: "user-1", AuthorNickname: "用户", Content: "评论",
	})
	require.NoError(t, err)

	err = svc.DeleteComment(context.Background(), c.ID, "user-1")
	require.NoError(t, err)
}

func TestDeleteComment_NotFound(t *testing.T) {
	svc := NewCommentService(newMockRepo())
	err := svc.DeleteComment(context.Background(), "nonexistent", "user-1")
	assert.ErrorIs(t, err, ErrCommentNotFound)
}

