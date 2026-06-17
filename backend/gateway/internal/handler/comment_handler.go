package handler

import (
	"net/http"
	"strconv"

	"showlove/gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

// CommentHandler handles comment-related requests.
type CommentHandler struct{}

// NewCommentHandler creates a new CommentHandler.
func NewCommentHandler() *CommentHandler {
	return &CommentHandler{}
}

// CreateCommentRequest represents the create comment request body.
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateComment handles creating a comment on a post.
func (h *CommentHandler) CreateComment(c *gin.Context) {
	postID := c.Param("id")
	userID := middleware.GetUserID(c)

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// TODO: Call comment-service gRPC CreateComment
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "评论成功",
		"data": gin.H{
			"id":        "temp-comment-id",
			"post_id":   postID,
			"author_id": userID,
			"content":   req.Content,
		},
	})
}

// ListComments handles listing comments for a post.
func (h *CommentHandler) ListComments(c *gin.Context) {
	postID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// TODO: Call comment-service gRPC ListComments
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"comments":  []interface{}{},
			"post_id":   postID,
			"page":      page,
			"page_size": pageSize,
			"total":     0,
		},
	})
}

// DeleteComment handles deleting a comment (author only).
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	userID := middleware.GetUserID(c)

	// TODO: Call comment-service gRPC DeleteComment
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "评论已删除",
		"data":    gin.H{"id": commentID, "deleted_by": userID},
	})
}
