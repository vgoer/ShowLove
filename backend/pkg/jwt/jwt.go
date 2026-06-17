// Package jwt provides JWT token generation and validation using HMAC-SHA256.
package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired = errors.New("token已过期")
	ErrTokenInvalid = errors.New("token无效")
)

// Claims represents the JWT claims for access tokens.
type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// Manager handles JWT token operations.
type Manager struct {
	secret           []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

// NewManager creates a new JWT Manager.
func NewManager(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		secret:          []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

// AccessTokenTTL returns the access token time-to-live.
func (m *Manager) AccessTokenTTL() time.Duration {
	return m.accessTokenTTL
}

// RefreshTokenTTL returns the refresh token time-to-live.
func (m *Manager) RefreshTokenTTL() time.Duration {
	return m.refreshTokenTTL
}

// GenerateAccessToken creates a new access token for the given user.
func (m *Manager) GenerateAccessToken(userID, email string) (string, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenTTL)),
			Subject:   userID,
		},
		UserID: userID,
		Email:  email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken parses and validates an access token, returning the claims.
func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// GenerateRefreshToken creates a cryptographically random refresh token.
func (m *Manager) GenerateRefreshToken(userID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
