package service

import (
	"context"
	"errors"

	"showlove/services/user-service/internal/model"
	"showlove/services/user-service/internal/repository"
)

// mockUserRepo implements UserRepository for testing.
type mockUserRepo struct {
	users map[string]*model.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*model.User)}
}

func (m *mockUserRepo) Create(_ context.Context, user *model.User) error {
	for _, u := range m.users {
		if u.Email == user.Email {
			return repository.ErrAlreadyExists
		}
	}
	if user.ID == "" {
		user.ID = "generated-id-" + user.Email
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) FindByID(_ context.Context, id string) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (m *mockUserRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepo) Update(_ context.Context, user *model.User) error {
	m.users[user.ID] = user
	return nil
}

// mockTokenRepo implements TokenRepository for testing.
type mockTokenRepo struct {
	tokens map[string]*model.RefreshToken
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{tokens: make(map[string]*model.RefreshToken)}
}

func (m *mockTokenRepo) Create(_ context.Context, token *model.RefreshToken) error {
	if token.ID == "" {
		token.ID = "rt-" + token.UserID
	}
	m.tokens[token.Token] = token
	return nil
}

func (m *mockTokenRepo) FindByToken(_ context.Context, token string) (*model.RefreshToken, error) {
	rt, ok := m.tokens[token]
	if !ok || rt.Revoked {
		return nil, errors.New("token not found")
	}
	return rt, nil
}

func (m *mockTokenRepo) Revoke(_ context.Context, tokenID string) error {
	for _, t := range m.tokens {
		if t.ID == tokenID {
			t.Revoked = true
			return nil
		}
	}
	return nil
}

func (m *mockTokenRepo) RevokeAllForUser(_ context.Context, userID string) error {
	for _, t := range m.tokens {
		if t.UserID == userID {
			t.Revoked = true
		}
	}
	return nil
}
