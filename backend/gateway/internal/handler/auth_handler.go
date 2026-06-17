package handler

import (
	"net/http"

	"showlove/pkg/validator"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests.
// In production, this would call the user-service gRPC client.
type AuthHandler struct{}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// RegisterRequest represents the registration request body.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

// Register handles user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if err := validator.ValidatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if err := validator.ValidateNickname(req.Nickname); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// TODO: Call user-service gRPC Register
	c.JSON(http.StatusNotImplemented, gin.H{"code": 501, "message": "user-service gRPC client pending"})
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// TODO: Call user-service gRPC Login
	c.JSON(http.StatusNotImplemented, gin.H{"code": 501, "message": "user-service gRPC client pending"})
}

// RefreshRequest represents the token refresh request body.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh handles token refresh.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// TODO: Call user-service gRPC RefreshToken
	c.JSON(http.StatusNotImplemented, gin.H{"code": 501, "message": "user-service gRPC client pending"})
}
