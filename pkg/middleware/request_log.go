package middleware

import (
	"TurboGin/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

// RequestLog 请求日志中间件
type RequestLog struct {
	log *logger.Logger
}

// NewRequestLog 构造函数
func NewRequestLog(log *logger.Logger) *RequestLog {
	return &RequestLog{log: log}
}

// Middleware 生成Gin中间件
func (r *RequestLog) Middleware() gin.HandlerFunc {
	if r == nil {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		fields := []logger.Field{
			logger.Int("status", c.Writer.Status()),
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", query),
			logger.String("ip", c.ClientIP()),
			logger.String("user-agent", c.Request.UserAgent()),
			logger.Duration("latency", latency),
			logger.String("time", end.Format(time.RFC3339)),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, logger.Any("errors", c.Errors.Errors()))
		}

		r.log.Info("HTTP Request", fields...)
	}
}
