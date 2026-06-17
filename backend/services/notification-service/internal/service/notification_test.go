package service

import (
	"context"
	"fmt"
	"testing"

	"showlove/services/notification-service/internal/model"
	"showlove/services/notification-service/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDeviceRepo struct {
	tokens map[string]*model.DeviceToken
	count  int
}

func newMockDeviceRepo() *mockDeviceRepo {
	return &mockDeviceRepo{tokens: make(map[string]*model.DeviceToken)}
}

func (m *mockDeviceRepo) Upsert(_ context.Context, d *model.DeviceToken) error {
	m.count++
	d.ID = fmt.Sprintf("dev-%d", m.count)
	m.tokens[d.Token] = d // same token overwrites
	return nil
}

func (m *mockDeviceRepo) FindByUserID(_ context.Context, userID string) ([]*model.DeviceToken, error) {
	var result []*model.DeviceToken
	for _, d := range m.tokens {
		if d.UserID == userID {
			result = append(result, d)
		}
	}
	return result, nil
}

func TestRegisterDevice(t *testing.T) {
	repo := newMockDeviceRepo()
	svc := NewNotificationService(repo, nil)

	dt, err := svc.RegisterDevice(context.Background(), "user-1", "fcm-token-abc123", "android")
	require.NoError(t, err)
	assert.Equal(t, "user-1", dt.UserID)
	assert.Equal(t, "android", dt.Platform)
	assert.NotEmpty(t, dt.ID)
}

func TestRegisterDevice_Upsert(t *testing.T) {
	repo := newMockDeviceRepo()
	svc := NewNotificationService(repo, nil)

	// Register same token for different users — token map key overwrites
	_, _ = svc.RegisterDevice(context.Background(), "user-1", "shared-token", "ios")
	_, _ = svc.RegisterDevice(context.Background(), "user-2", "shared-token", "android")

	// The last upsert wins for the same token
	assert.Equal(t, 2, repo.count)
}

func TestSendPush(t *testing.T) {
	repo := newMockDeviceRepo()
	svc := NewNotificationService(repo, nil)

	_, _ = svc.RegisterDevice(context.Background(), "user-1", "token-aaa", "ios")
	_, _ = svc.RegisterDevice(context.Background(), "user-1", "token-bbb", "web")

	// Should not error
	err := svc.SendPush(context.Background(), "user-1", "新评论", "有人评论了你的帖子")
	assert.NoError(t, err)
}

// Ensure mock implements interface
var _ repository.DeviceRepository = (*mockDeviceRepo)(nil)
