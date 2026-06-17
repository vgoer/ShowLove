package handler

import (
	"net/http"
	"strconv"

	"showlove/gateway/internal/middleware"
	"showlove/pkg/validator"

	"github.com/gin-gonic/gin"
)

// PostHandler handles post-related requests.
type PostHandler struct{}

// NewPostHandler creates a new PostHandler.
func NewPostHandler() *PostHandler {
	return &PostHandler{}
}

// CreatePostRequest represents the create post request body.
type CreatePostRequest struct {
	Content string   `json:"content" binding:"required"`
	MoodTag string   `json:"mood_tag" binding:"required"`
	Images  []string `json:"images"`
}

// CreatePost handles post creation.
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := validator.ValidateContent(req.Content, 5000); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// TODO: Call post-service gRPC CreatePost
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "发布成功",
		"data": gin.H{
			"id":            "temp-id",
			"author_id":     userID,
			"content":       req.Content,
			"mood_tag":      req.MoodTag,
			"images":        req.Images,
			"comment_count": 0,
		},
	})
}

// ListPosts handles post listing with sorting and pagination.
func (h *PostHandler) ListPosts(c *gin.Context) {
	sort := c.DefaultQuery("sort", "latest")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// TODO: Call post-service gRPC ListPosts
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"posts":     []interface{}{},
			"page":      page,
			"page_size": pageSize,
			"total":     0,
			"sort":      sort,
		},
	})
}

// GetPost handles retrieving a single post.
func (h *PostHandler) GetPost(c *gin.Context) {
	postID := c.Param("id")

	// TODO: Call post-service gRPC GetPost
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id": postID,
		},
	})
}

// DeletePost handles post deletion (author only).
func (h *PostHandler) DeletePost(c *gin.Context) {
	postID := c.Param("id")
	userID := middleware.GetUserID(c)

	// TODO: Call post-service gRPC DeletePost
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    gin.H{"id": postID, "deleted_by": userID},
	})
}

// SendStickerRequest represents the send sticker request body.
type SendStickerRequest struct {
	StickerType string `json:"sticker_type" binding:"required"` // hug, cheer, understand
}

// SendSticker handles sending a sticker to a post.
func (h *PostHandler) SendSticker(c *gin.Context) {
	postID := c.Param("id")

	var req SendStickerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// TODO: Call post-service gRPC SendSticker
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "贴纸发送成功",
		"data":    gin.H{"post_id": postID, "sticker_type": req.StickerType},
	})
}

// ReportPostRequest represents the report post request body.
type ReportPostRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// ReportPost handles reporting a post.
func (h *PostHandler) ReportPost(c *gin.Context) {
	postID := c.Param("id")

	var req ReportPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// TODO: Call post-service gRPC ReportPost
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "举报已提交",
		"data":    gin.H{"post_id": postID, "reason": req.Reason},
	})
}
