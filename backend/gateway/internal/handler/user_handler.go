package handler

import (
	"net/http"

	"showlove/gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user profile requests.
type UserHandler struct{}

// NewUserHandler creates a new UserHandler.
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetMe returns the current user's profile.
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// TODO: Call user-service gRPC GetUser
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id": userID,
		},
	})
}

// UpdateMeRequest represents the profile update request body.
type UpdateMeRequest struct {
	Nickname *string `json:"nickname"`
	Bio      *string `json:"bio"`
}

// UpdateMe updates the current user's profile.
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// TODO: Call user-service gRPC UpdateProfile
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data": gin.H{
			"id": userID,
		},
	})
}

// UploadAvatar handles avatar upload.
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := middleware.GetUserID(c)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择头像文件"})
		return
	}

	// TODO: Upload to storage and call user-service gRPC UploadAvatar
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "上传成功",
		"data": gin.H{
			"id":        userID,
			"file_name": file.Filename,
			"file_size": file.Size,
		},
	})
}
