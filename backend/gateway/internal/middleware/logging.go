package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// ContextKeyTraceID is the Gin context key for the request trace ID.
	ContextKeyTraceID = "trace_id"
)

// Logging returns a middleware that logs each request with trace ID, method, path, and latency.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set(ContextKeyTraceID, traceID)
		c.Header("X-Request-Id", traceID)

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		log.Printf("[gateway] %s | %s %s | %d | %v | %s",
			traceID[:8], method, path, status, latency, clientIP)
	}
}
