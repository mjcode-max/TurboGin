package middleware

import (
	"TurboGin/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// IPAccess IP访问控制中间件
type IPAccess struct {
	allowedIPs []string // 允许的IP列表
}

// NewIPAccess 构造函数
func NewIPAccess(cfg *config.Config) *IPAccess {
	return &IPAccess{
		allowedIPs: cfg.Server.TrustedProxies,
	}
}

// Middleware 生成Gin中间件
func (ia *IPAccess) Middleware() gin.HandlerFunc {
	if ia == nil {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		clientIP := getClientIP(c)

		// 检查IP是否在允许列表中
		if !ia.isIPAllowed(clientIP) {
			abortWithError(c, http.StatusForbidden, "Access denied for your IP address")
			return
		}

		c.Next()
	}
}

// isIPAllowed 检查IP是否被允许访问
func (ia *IPAccess) isIPAllowed(ip string) bool {
	for _, allowedIP := range ia.allowedIPs {
		if ip == allowedIP {
			return true
		}
	}
	return false
}

// getClientIP 从请求中获取客户端IP
func getClientIP(c *gin.Context) string {
	// 尝试从X-Forwarded-For获取(如果有代理)
	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// 如果没有代理，直接使用RemoteAddr
	return strings.Split(c.Request.RemoteAddr, ":")[0]
}
