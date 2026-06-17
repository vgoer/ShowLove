// Package router registers all API routes.
package router

import (
	"showlove/gateway/internal/handler"
	"showlove/gateway/internal/middleware"
	"showlove/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Setup configures all routes and middleware on the Gin engine.
func Setup(jwtMgr *jwt.Manager) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logging())

	// Rate limiter: 100 requests per second
	rl := middleware.NewRateLimiter(100, 100)
	r.Use(rl.Middleware())

	// Handlers
	authH := handler.NewAuthHandler()
	userH := handler.NewUserHandler()
	healthH := handler.NewHealthHandler()
	uploadH := handler.NewUploadHandler()

	// Health check (no auth)
	r.GET("/api/v1/health", healthH.Health)

	// Auth routes (no auth required)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
	}

	// Protected routes (auth required)
	api := r.Group("/api/v1")
	api.Use(middleware.Auth(jwtMgr))
	{
		// User
		api.GET("/users/me", userH.GetMe)
		api.PUT("/users/me", userH.UpdateMe)
		api.PUT("/users/me/avatar", userH.UploadAvatar)

		// Upload
		api.POST("/upload/image", uploadH.UploadImage)
	}

	return r
}
