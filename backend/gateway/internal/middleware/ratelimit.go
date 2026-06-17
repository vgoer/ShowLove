package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter per IP.
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64 // tokens per second
	capacity float64 // max tokens
}

type tokenBucket struct {
	tokens   float64
	lastTime time.Time
}

// NewRateLimiter creates a new rate limiter.
// rate: requests per second allowed, capacity: burst size.
func NewRateLimiter(rate, capacity float64) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*tokenBucket),
		rate:     rate,
		capacity: capacity,
	}
}

// Middleware returns a Gin middleware that rate-limits requests by IP.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		bucket, exists := rl.buckets[ip]
		if !exists {
			bucket = &tokenBucket{
				tokens:   rl.capacity,
				lastTime: time.Now(),
			}
			rl.buckets[ip] = bucket
		}
		rl.mu.Unlock()

		// Refill tokens
		now := time.Now()
		elapsed := now.Sub(bucket.lastTime).Seconds()
		bucket.tokens += elapsed * rl.rate
		if bucket.tokens > rl.capacity {
			bucket.tokens = rl.capacity
		}
		bucket.lastTime = now

		if bucket.tokens < 1 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			return
		}

		bucket.tokens--
		c.Next()
	}
}
