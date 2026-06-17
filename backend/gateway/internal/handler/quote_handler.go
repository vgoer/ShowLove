package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type QuoteHandler struct{}

func NewQuoteHandler() *QuoteHandler { return &QuoteHandler{} }

func (h *QuoteHandler) GetTodayQuote(c *gin.Context) {
	// TODO: Call quote-service gRPC
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"text_zh": "万物皆有裂痕，那是光照进来的地方。",
		"text_en": "There is a crack in everything. That's how the light gets in.",
		"author":  "Leonard Cohen",
	}})
}
