// Package service implements the business logic for the user service.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"showlove/pkg/jwt"
	"showlove/services/user-service/internal/model"
	"showlove/services/user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("邮箱或密码错误")
	ErrEmailAlreadyUsed   = errors.New("该邮箱已被注册")
	ErrTokenRevoked       = errors.New("token已被吊销")
)

// UserService handles user business logic.
type UserService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	jwtMgr    *jwt.Manager
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, jwtMgr *jwt.Manager) *UserService {
	return &UserService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtMgr:    jwtMgr,
	}
}

// RegisterParams contains the registration input.
type RegisterParams struct {
	Email    string
	Password string
	Nickname string
}

// RegisterResult contains the registration output.
type RegisterResult struct {
	User         *model.User
	AccessToken  string
	RefreshToken string
}

// Register creates a new user account.
func (s *UserService) Register(ctx context.Context, params RegisterParams) (*RegisterResult, error) {
	// Check if email already exists
	_, err := s.userRepo.FindByEmail(ctx, params.Email)
	if err == nil {
		return nil, ErrEmailAlreadyUsed
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("user service: check email: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("user service: hash password: %w", err)
	}

	user := &model.User{
		Email:    params.Email,
		Password: string(hashedPassword),
		Nickname: params.Nickname,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return nil, ErrEmailAlreadyUsed
		}
		return nil, fmt.Errorf("user service: create user: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtMgr.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("user service: generate access token: %w", err)
	}

	refreshToken, err := s.jwtMgr.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("user service: generate refresh token: %w", err)
	}

	// Store refresh token
	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.jwtMgr.RefreshTokenTTL()),
	}
	if err := s.tokenRepo.Create(ctx, rt); err != nil {
		return nil, fmt.Errorf("user service: store refresh token: %w", err)
	}

	return &RegisterResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// LoginParams contains the login input.
type LoginParams struct {
	Email    string
	Password string
}

// Login authenticates a user and returns tokens.
func (s *UserService) Login(ctx context.Context, params LoginParams) (*RegisterResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("user service: find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.jwtMgr.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("user service: generate access token: %w", err)
	}

	refreshToken, err := s.jwtMgr.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("user service: generate refresh token: %w", err)
	}

	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.jwtMgr.RefreshTokenTTL()),
	}
	if err := s.tokenRepo.Create(ctx, rt); err != nil {
		return nil, fmt.Errorf("user service: store refresh token: %w", err)
	}

	return &RegisterResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshAccessToken validates a refresh token and issues a new access token.
func (s *UserService) RefreshAccessToken(ctx context.Context, tokenStr string) (string, error) {
	rt, err := s.tokenRepo.FindByToken(ctx, tokenStr)
	if err != nil {
		return "", ErrTokenRevoked
	}

	user, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return "", fmt.Errorf("user service: find user: %w", err)
	}

	// Revoke the old refresh token (rotation)
	if err := s.tokenRepo.Revoke(ctx, rt.ID); err != nil {
		return "", fmt.Errorf("user service: revoke token: %w", err)
	}

	accessToken, err := s.jwtMgr.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", fmt.Errorf("user service: generate access token: %w", err)
	}

	return accessToken, nil
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("user service: get user %s: %w", userID, err)
	}
	return user, nil
}

// UpdateProfile updates a user's profile fields.
func (s *UserService) UpdateProfile(ctx context.Context, userID string, nickname, bio *string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if nickname != nil {
		user.Nickname = *nickname
	}
	if bio != nil {
		user.Bio = *bio
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("user service: update user: %w", err)
	}

	return user, nil
}

// UpdateAvatar updates a user's avatar URL.
func (s *UserService) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.AvatarURL = avatarURL
	return s.userRepo.Update(ctx, user)
}
