package handler

import (
	"net/http"

	"showlove/gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

type MoodHandler struct{}

func NewMoodHandler() *MoodHandler { return &MoodHandler{} }

type RecordMoodRequest struct {
	MoodLevel int32  `json:"mood_level" binding:"required"`
	MoodLabel string `json:"mood_label" binding:"required"`
	Note      string `json:"note"`
}

func (h *MoodHandler) RecordMood(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req RecordMoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	// TODO: Call mood-service gRPC
	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "记录成功", "data": gin.H{
		"user_id": userID, "mood_level": req.MoodLevel, "mood_label": req.MoodLabel,
	}})
}

func (h *MoodHandler) GetMoods(c *gin.Context) {
	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")
	// TODO: Call mood-service gRPC
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"entries": []interface{}{}, "from": from, "to": to}})
}

func (h *MoodHandler) GetWeeklyMood(c *gin.Context) {
	// TODO: Call mood-service gRPC
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"points": []interface{}{}}})
}
