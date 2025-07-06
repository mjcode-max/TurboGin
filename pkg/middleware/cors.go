package middleware

import (
	"TurboGin/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// CORS 跨域中间件
type CORS struct {
	cfg *config.CORSConfig
}

// NewCORS 构造函数
func NewCORS(cfg *config.Config) *CORS {
	if !cfg.Middleware.CORS.Enabled {
		return nil
	}
	return &CORS{cfg: &cfg.Middleware.CORS}
}

// Middleware 生成Gin中间件
func (c *CORS) Middleware() gin.HandlerFunc {
	if c == nil {
		return func(ctx *gin.Context) { ctx.Next() }
	}

	config := cors.Config{
		AllowOrigins:     c.cfg.AllowOrigins,
		AllowMethods:     c.cfg.AllowMethods,
		AllowHeaders:     c.cfg.AllowHeaders,
		ExposeHeaders:    c.cfg.ExposeHeaders,
		AllowCredentials: c.cfg.AllowCredentials,
		MaxAge:           12 * time.Hour,
	}

	// 动态允许Origin
	if len(c.cfg.AllowOrigins) == 1 && c.cfg.AllowOrigins[0] == "*" {
		config.AllowOriginFunc = func(origin string) bool {
			return true
		}
	}

	return cors.New(config)
}
