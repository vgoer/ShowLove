package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidate_Success(t *testing.T) {
	secret := "test-secret-key-32-bytes-long!!"
	mgr := NewManager(secret, 15*time.Minute, 7*24*time.Hour)

	userID := "user-123"
	email := "test@example.com"

	token, err := mgr.GenerateAccessToken(userID, email)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := mgr.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestGenerateAndValidate_ExpiredToken(t *testing.T) {
	secret := "test-secret-key-32-bytes-long!!"
	mgr := NewManager(secret, -1*time.Second, 7*24*time.Hour) // negative TTL = instant expiry

	userID := "user-123"
	email := "test@example.com"

	token, err := mgr.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	_, err = mgr.ValidateToken(token)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestGenerateAndValidate_InvalidToken(t *testing.T) {
	secret := "test-secret-key-32-bytes-long!!"
	mgr := NewManager(secret, 15*time.Minute, 7*24*time.Hour)

	_, err := mgr.ValidateToken("invalid.token.here")
	assert.Error(t, err)
}

func TestGenerateAndValidate_WrongSecret(t *testing.T) {
	mgr1 := NewManager("secret-key-aaaaaaaaaaaaaaaaaaaaaa", 15*time.Minute, 7*24*time.Hour)
	mgr2 := NewManager("secret-key-bbbbbbbbbbbbbbbbbbbbbb", 15*time.Minute, 7*24*time.Hour)

	token, err := mgr1.GenerateAccessToken("user-123", "test@example.com")
	require.NoError(t, err)

	_, err = mgr2.ValidateToken(token)
	assert.Error(t, err)
}

func TestGenerateRefreshToken(t *testing.T) {
	secret := "test-secret-key-32-bytes-long!!"
	mgr := NewManager(secret, 15*time.Minute, 7*24*time.Hour)

	token, err := mgr.GenerateRefreshToken("user-123")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	assert.True(t, len(token) > 32) // should be a decently long random string
}

func TestRefreshToken_HasCorrectTTL(t *testing.T) {
	secret := "test-secret-key-32-bytes-long!!"
	mgr := NewManager(secret, 15*time.Minute, 7*24*time.Hour)

	assert.Equal(t, 15*time.Minute, mgr.AccessTokenTTL())
	assert.Equal(t, 7*24*time.Hour, mgr.RefreshTokenTTL())
}
