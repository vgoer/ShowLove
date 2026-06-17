//go:build integration
// +build integration

// Package integration provides integration tests that require a running database.
// Run with: go test -tags=integration ./test/integration/... -v
package integration

import (
	"context"
	"testing"
	"time"

	"showlove/pkg/jwt"
	"showlove/services/user-service/internal/model"
	"showlove/services/user-service/internal/repository"
	"showlove/services/user-service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory test database connection.
// In a real environment, this would use testcontainers-go.
// For CI, set DB_DSN to a test database.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// Use environment variable or skip
	dsn := "postgres://showlove:showlove123@localhost:5432/users_db?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
	}
	// Auto-migrate
	db.AutoMigrate(&model.User{}, &model.RefreshToken{})
	t.Cleanup(func() {
		db.Exec("DELETE FROM refresh_tokens")
		db.Exec("DELETE FROM users")
	})
	return db
}

func TestIntegration_UserService_RegisterAndLogin(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	jwtMgr := jwt.NewManager("integration-test-secret-key!!!", 15*time.Minute, 7*24*time.Hour)
	svc := service.NewUserService(userRepo, tokenRepo, jwtMgr)

	// Register
	result, err := svc.Register(ctx, service.RegisterParams{
		Email:    "integration@example.com",
		Password: "TestPass123",
		Nickname: "集成测试用户",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.User.ID)
	assert.NotEmpty(t, result.AccessToken)

	// Login
	loginResult, err := svc.Login(ctx, service.LoginParams{
		Email:    "integration@example.com",
		Password: "TestPass123",
	})
	require.NoError(t, err)
	assert.Equal(t, result.User.ID, loginResult.User.ID)

	// Get user
	user, err := svc.GetUser(ctx, result.User.ID)
	require.NoError(t, err)
	assert.Equal(t, "integration@example.com", user.Email)
}

func TestIntegration_UserService_RefreshToken(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	jwtMgr := jwt.NewManager("integration-test-secret-key!!!", 15*time.Minute, 7*24*time.Hour)
	svc := service.NewUserService(userRepo, tokenRepo, jwtMgr)

	// Register to get tokens
	result, err := svc.Register(ctx, service.RegisterParams{
		Email:    "refresh-int@example.com",
		Password: "TestPass123",
		Nickname: "刷新测试",
	})
	require.NoError(t, err)

	// Refresh
	newToken, err := svc.RefreshAccessToken(ctx, result.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newToken)
}
