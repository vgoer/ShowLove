package service

import (
	"context"
	"errors"
	"fmt"

	"showlove/services/post-service/internal/model"
	"showlove/services/post-service/internal/repository"
)

// mockPostRepo implements PostRepository for testing.
type mockPostRepo struct {
	posts   map[string]*model.Post
	counter int
}

func newMockPostRepo() *mockPostRepo {
	return &mockPostRepo{posts: make(map[string]*model.Post)}
}

func (m *mockPostRepo) Create(_ context.Context, post *model.Post) error {
	m.counter++
	if post.ID == "" {
		post.ID = fmt.Sprintf("post-%d", m.counter)
	}
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepo) FindByID(_ context.Context, id string) (*model.Post, error) {
	post, ok := m.posts[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return post, nil
}

func (m *mockPostRepo) List(_ context.Context, _ string, page, pageSize int32) ([]*model.Post, int64, error) {
	var result []*model.Post
	for _, p := range m.posts {
		if !p.IsHidden {
			result = append(result, p)
		}
	}
	total := int64(len(result))
	offset := (page - 1) * pageSize
	if offset >= int32(len(result)) {
		return []*model.Post{}, total, nil
	}
	end := offset + pageSize
	if end > int32(len(result)) {
		end = int32(len(result))
	}
	return result[offset:end], total, nil
}

func (m *mockPostRepo) Update(_ context.Context, post *model.Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepo) Delete(_ context.Context, id, authorID string) error {
	post, ok := m.posts[id]
	if !ok {
		return repository.ErrNotFound
	}
	if post.AuthorID != authorID {
		return errors.New("permission denied")
	}
	delete(m.posts, id)
	return nil
}
