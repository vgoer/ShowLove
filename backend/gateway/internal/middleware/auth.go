package middleware

import (
	"net/http"
	"strings"

	"showlove/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyUserID is the Gin context key for the authenticated user ID.
	ContextKeyUserID = "user_id"
	// ContextKeyUserEmail is the Gin context key for the authenticated user email.
	ContextKeyUserEmail = "user_email"
)

// Auth returns a middleware that validates JWT Bearer tokens.
func Auth(jwtMgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少认证信息",
			})
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误，请使用 Bearer Token",
			})
			return
		}

		token := parts[1]

		claims, err := jwtMgr.ValidateToken(token)
		if err != nil {
			message := "token无效"
			if err == jwt.ErrTokenExpired {
				message = "token已过期"
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": message,
			})
			return
		}

		// Inject user info into context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUserEmail, claims.Email)

		c.Next()
	}
}

// GetUserID extracts the authenticated user ID from the Gin context.
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get(ContextKeyUserID)
	return userID.(string)
}
