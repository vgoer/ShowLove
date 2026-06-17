package handler

import (
	"net/http"

	"showlove/gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct{}

func NewDeviceHandler() *DeviceHandler { return &DeviceHandler{} }

type RegisterDeviceRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required"` // ios, android, web
}

func (h *DeviceHandler) RegisterDevice(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if req.Platform != "ios" && req.Platform != "android" && req.Platform != "web" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "平台类型无效"})
		return
	}
	// TODO: Call notification-service gRPC
	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "设备注册成功", "data": gin.H{
		"user_id": userID, "platform": req.Platform,
	}})
}
