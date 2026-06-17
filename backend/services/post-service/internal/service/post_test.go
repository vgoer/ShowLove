package service

import (
	"context"
	"testing"

	"showlove/pkg/events"
	"showlove/services/post-service/internal/moderation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestService() (*PostService, *mockPostRepo) {
	repo := newMockPostRepo()
	filter := moderation.NewFilter(moderation.DefaultChineseWords())
	ps := events.NewNoOpPubSub()
	svc := NewPostService(repo, filter, ps)
	return svc, repo
}

func TestCreatePost_Success(t *testing.T) {
	svc, _ := setupTestService()

	post, err := svc.CreatePost(context.Background(), CreatePostParams{
		AuthorID:       "user-1",
		AuthorNickname: "测试用户",
		Content:        "今天心情不太好，工作压力好大",
		MoodTag:        "anxious",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, post.ID)
	assert.Equal(t, "user-1", post.AuthorID)
	assert.Equal(t, "测试用户", post.AuthorNickname)
	assert.Equal(t, "anxious", post.MoodTag)
	assert.False(t, post.HasAIReply)
}

func TestCreatePost_SensitiveContent(t *testing.T) {
	svc, _ := setupTestService()

	_, err := svc.CreatePost(context.Background(), CreatePostParams{
		AuthorID:       "user-1",
		AuthorNickname: "测试用户",
		Content:        "这里涉及赌博内容",
		MoodTag:        "sad",
	})

	assert.ErrorIs(t, err, ErrSensitiveContent)
	assert.Contains(t, err.Error(), "赌博")
}

func TestGetPost_Success(t *testing.T) {
	svc, _ := setupTestService()

	created, err := svc.CreatePost(context.Background(), CreatePostParams{
		AuthorID:       "user-1",
		AuthorNickname: "测试用户",
		Content:        "一篇测试帖子",
		MoodTag:        "calm",
	})
	require.NoError(t, err)

	post, err := svc.GetPost(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, "一篇测试帖子", post.Content)
}

func TestGetPost_NotFound(t *testing.T) {
	svc, _ := setupTestService()

	_, err := svc.GetPost(context.Background(), "non-existent-id")
	assert.ErrorIs(t, err, ErrPostNotFound)
}

func TestListPosts(t *testing.T) {
	svc, _ := setupTestService()

	// Create 3 posts
	for i := 0; i < 3; i++ {
		_, err := svc.CreatePost(context.Background(), CreatePostParams{
			AuthorID:       "user-1",
			AuthorNickname: "测试用户",
			Content:        "帖子内容",
			MoodTag:        "sad",
		})
		require.NoError(t, err)
	}

	posts, total, err := svc.ListPosts(context.Background(), "latest", 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, posts, 3)
}

func TestDeletePost_Success(t *testing.T) {
	svc, _ := setupTestService()

	post, err := svc.CreatePost(context.Background(), CreatePostParams{
		AuthorID:       "user-1",
		AuthorNickname: "用户A",
		Content:        "我的帖子",
		MoodTag:        "happy",
	})
	require.NoError(t, err)

	err = svc.DeletePost(context.Background(), post.ID, "user-1")
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetPost(context.Background(), post.ID)
	assert.ErrorIs(t, err, ErrPostNotFound)
}

func TestSendSticker(t *testing.T) {
	svc, _ := setupTestService()

	post, err := svc.CreatePost(context.Background(), CreatePostParams{
		AuthorID: "user-1", AuthorNickname: "用户", Content: "需要帮助", MoodTag: "sad",
	})
	require.NoError(t, err)

	// Send hug sticker
	updated, err := svc.SendSticker(context.Background(), post.ID, "hug")
	require.NoError(t, err)
	assert.Equal(t, int32(1), updated.StickerHug)
	assert.Equal(t, int32(0), updated.StickerCheer)

	// Send cheer sticker
	updated, err = svc.SendSticker(context.Background(), post.ID, "cheer")
	require.NoError(t, err)
	assert.Equal(t, int32(1), updated.StickerCheer)
}

func TestReportPost(t *testing.T) {
	svc, _ := setupTestService()

	post, err := svc.CreatePost(context.Background(), CreatePostParams{
		AuthorID: "user-1", AuthorNickname: "用户", Content: "不当内容", MoodTag: "angry",
	})
	require.NoError(t, err)

	err = svc.ReportPost(context.Background(), post.ID, "user-2", "内容不当")
	require.NoError(t, err)

	reported, err := svc.GetPost(context.Background(), post.ID)
	require.NoError(t, err)
	assert.True(t, reported.IsReported)
}
