package service

import (
	"context"
	"testing"
	"time"

	"showlove/pkg/jwt"
	"showlove/services/user-service/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestService() (*UserService, *mockUserRepo, *mockTokenRepo) {
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()
	jwtMgr := jwt.NewManager("test-secret-key-for-testing-only!!", 15*time.Minute, 7*24*time.Hour)
	svc := NewUserService(userRepo, tokenRepo, jwtMgr)
	return svc, userRepo, tokenRepo
}

func TestRegister_Success(t *testing.T) {
	svc, _, _ := setupTestService()

	result, err := svc.Register(context.Background(), RegisterParams{
		Email:    "test@example.com",
		Password: "securePass123",
		Nickname: "测试用户",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.User.ID)
	assert.Equal(t, "test@example.com", result.User.Email)
	assert.Equal(t, "测试用户", result.User.Nickname)
	assert.NotEmpty(t, result.User.Password) // hashed password is stored
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, _, _ := setupTestService()

	_, err := svc.Register(context.Background(), RegisterParams{
		Email:    "dupe@example.com",
		Password: "securePass123",
		Nickname: "First",
	})
	require.NoError(t, err)

	_, err = svc.Register(context.Background(), RegisterParams{
		Email:    "dupe@example.com",
		Password: "anotherPass456",
		Nickname: "Second",
	})

	assert.ErrorIs(t, err, ErrEmailAlreadyUsed)
}

func TestLogin_Success(t *testing.T) {
	svc, _, _ := setupTestService()

	_, err := svc.Register(context.Background(), RegisterParams{
		Email:    "login@example.com",
		Password: "securePass123",
		Nickname: "登录测试",
	})
	require.NoError(t, err)

	result, err := svc.Login(context.Background(), LoginParams{
		Email:    "login@example.com",
		Password: "securePass123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, "login@example.com", result.User.Email)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _, _ := setupTestService()

	_, err := svc.Register(context.Background(), RegisterParams{
		Email:    "wp@example.com",
		Password: "correctPass123",
		Nickname: "密码测试",
	})
	require.NoError(t, err)

	_, err = svc.Login(context.Background(), LoginParams{
		Email:    "wp@example.com",
		Password: "wrongPass456",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_WrongEmail(t *testing.T) {
	svc, _, _ := setupTestService()

	_, err := svc.Login(context.Background(), LoginParams{
		Email:    "nonexistent@example.com",
		Password: "securePass123",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestRefreshToken_Success(t *testing.T) {
	svc, _, _ := setupTestService()

	result, err := svc.Register(context.Background(), RegisterParams{
		Email:    "refresh@example.com",
		Password: "securePass123",
		Nickname: "刷新测试",
	})
	require.NoError(t, err)

	newToken, err := svc.RefreshAccessToken(context.Background(), result.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
}

func TestRefreshToken_Revoked(t *testing.T) {
	svc, _, _ := setupTestService()

	result, err := svc.Register(context.Background(), RegisterParams{
		Email:    "revoked@example.com",
		Password: "securePass123",
		Nickname: "吊销测试",
	})
	require.NoError(t, err)

	// First refresh should succeed
	_, err = svc.RefreshAccessToken(context.Background(), result.RefreshToken)
	require.NoError(t, err)

	// Second refresh with the same (now revoked) token should fail
	_, err = svc.RefreshAccessToken(context.Background(), result.RefreshToken)
	assert.ErrorIs(t, err, ErrTokenRevoked)
}

func TestGetUser_Success(t *testing.T) {
	svc, _, _ := setupTestService()

	result, err := svc.Register(context.Background(), RegisterParams{
		Email:    "get@example.com",
		Password: "securePass123",
		Nickname: "查询测试",
	})
	require.NoError(t, err)

	user, err := svc.GetUser(context.Background(), result.User.ID)
	require.NoError(t, err)
	assert.Equal(t, "get@example.com", user.Email)
}

func TestGetUser_NotFound(t *testing.T) {
	svc, _, _ := setupTestService()

	_, err := svc.GetUser(context.Background(), "non-existent-id")
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestUpdateProfile_Success(t *testing.T) {
	svc, _, _ := setupTestService()

	result, err := svc.Register(context.Background(), RegisterParams{
		Email:    "update@example.com",
		Password: "securePass123",
		Nickname: "旧昵称",
	})
	require.NoError(t, err)

	newNickname := "新昵称"
	updatedUser, err := svc.UpdateProfile(context.Background(), result.User.ID, &newNickname, nil)
	require.NoError(t, err)
	assert.Equal(t, "新昵称", updatedUser.Nickname)
}
