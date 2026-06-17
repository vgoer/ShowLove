package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"showlove/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCORSMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, 200, w.Code)
}

func TestCORSMiddleware_OptionsRequest(t *testing.T) {
	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	jwtMgr := jwt.NewManager("test-secret-key-for-testing!!", 15*time.Minute, 7*24*time.Hour)

	r := gin.New()
	r.Use(Auth(jwtMgr))
	r.GET("/protected", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "缺少认证信息")
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtMgr := jwt.NewManager("test-secret-key-for-testing!!", 15*time.Minute, 7*24*time.Hour)

	token, err := jwtMgr.GenerateAccessToken("user-123", "test@example.com")
	assert.NoError(t, err)

	r := gin.New()
	r.Use(Auth(jwtMgr))
	r.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get(ContextKeyUserID)
		c.String(200, userID.(string))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "user-123", w.Body.String())
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	jwtMgr := jwt.NewManager("test-secret-key-for-testing!!", -1*time.Second, 7*24*time.Hour)

	token, err := jwtMgr.GenerateAccessToken("user-123", "test@example.com")
	assert.NoError(t, err)

	r := gin.New()
	r.Use(Auth(jwtMgr))
	r.GET("/protected", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	jwtMgr := jwt.NewManager("test-secret-key-for-testing!!", 15*time.Minute, 7*24*time.Hour)

	r := gin.New()
	r.Use(Auth(jwtMgr))
	r.GET("/protected", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestRateLimiter_AllowsRequest(t *testing.T) {
	rl := NewRateLimiter(100, 100)

	r := gin.New()
	r.Use(rl.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRateLimiter_BlocksExcess(t *testing.T) {
	rl := NewRateLimiter(0.01, 1) // very low rate

	r := gin.New()
	r.Use(rl.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// First request should pass
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, 200, w1.Code)

	// Second request should be rate-limited
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 429, w2.Code)
}

func TestLoggingMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(Logging())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-Id"))
}
