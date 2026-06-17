package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles file upload requests.
type UploadHandler struct{}

// NewUploadHandler creates a new UploadHandler.
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage handles image file upload.
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择文件"})
		return
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的图片格式，仅支持 JPEG/PNG/GIF/WebP"})
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "图片大小不能超过 10MB"})
		return
	}

	// TODO: Upload to storage (MinIO)
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "上传成功",
		"data": gin.H{
			"file_name": file.Filename,
			"file_size": file.Size,
			"url":       "http://localhost:9000/showlove/images/" + file.Filename,
		},
	})
}
