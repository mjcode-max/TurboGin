package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mjcode-max/TurboGin/config"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

// RateLimiter 基于IP的限流器
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rps      float64
	burst    int
	enabled  bool
}

// NewRateLimiter 构造函数
func NewRateLimiter(cfg *config.Config) *RateLimiter {
	if !cfg.Middleware.RateLimit.Enabled {
		return nil
	}
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      cfg.Middleware.RateLimit.RPS,
		burst:    cfg.Middleware.RateLimit.Burst,
		enabled:  true,
	}
}

// Middleware 生成Gin中间件
func (r *RateLimiter) Middleware() gin.HandlerFunc {
	if r == nil || !r.enabled {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		limiter := r.getLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
				"code":  http.StatusTooManyRequests,
			})
			return
		}
		c.Next()
	}
}

// getLimiter 获取或创建限流器
func (r *RateLimiter) getLimiter(ip string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	if limiter, exists := r.limiters[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(r.rps), r.burst)
	r.limiters[ip] = limiter
	return limiter
}
