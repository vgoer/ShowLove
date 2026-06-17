//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"showlove/pkg/events"
	"showlove/services/post-service/internal/model"
	"showlove/services/post-service/internal/moderation"
	"showlove/services/post-service/internal/repository"
	"showlove/services/post-service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupPostDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "postgres://showlove:showlove123@localhost:5432/posts_db?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
	}
	db.AutoMigrate(&model.Post{})
	t.Cleanup(func() {
		db.Exec("DELETE FROM posts")
	})
	return db
}

func TestIntegration_PostService_CreateAndList(t *testing.T) {
	db := setupPostDB(t)
	ctx := context.Background()

	repo := repository.NewPostRepository(db)
	filter := moderation.NewFilter(moderation.DefaultChineseWords())
	pub := events.NewNoOpPubSub()
	svc := service.NewPostService(repo, filter, pub)

	// Create posts
	for i := 0; i < 3; i++ {
		_, err := svc.CreatePost(ctx, service.CreatePostParams{
			AuthorID:       "user-1",
			AuthorNickname: "集成测试",
			Content:        "这是一篇集成测试帖子",
			MoodTag:        "calm",
		})
		require.NoError(t, err)
	}

	// List
	posts, total, err := svc.ListPosts(ctx, "latest", 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, posts, 3)
}

func TestIntegration_PostService_SensitiveContentRejected(t *testing.T) {
	db := setupPostDB(t)
	ctx := context.Background()

	repo := repository.NewPostRepository(db)
	filter := moderation.NewFilter(moderation.DefaultChineseWords())
	svc := service.NewPostService(repo, filter, events.NewNoOpPubSub())

	_, err := svc.CreatePost(ctx, service.CreatePostParams{
		AuthorID:       "user-1",
		AuthorNickname: "测试",
		Content:        "涉及赌博的违规内容",
		MoodTag:        "angry",
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, service.ErrSensitiveContent)
}
