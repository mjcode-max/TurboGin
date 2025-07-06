package middleware

import (
	"TurboGin/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// Auth JWT认证中间件
type Auth struct {
	cfg *config.AuthConfig
}

// NewAuth 构造函数
func NewAuth(cfg *config.Config) *Auth {
	if !cfg.JWT.Enabled {
		return nil
	}
	return &Auth{cfg: &cfg.JWT}
}

// Middleware 生成Gin中间件
func (a *Auth) Middleware() gin.HandlerFunc {
	if a == nil {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			abortWithError(c, http.StatusUnauthorized, "Authorization header required")
			return
		}

		claims, err := a.parseToken(tokenString)
		if err != nil {
			abortWithError(c, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// 存储claims到上下文
		for k, v := range claims {
			c.Set(k, v)
		}
		c.Next()
	}
}

// GenerateToken 生成JWT令牌 (供Service层调用)
func (a *Auth) GenerateToken(userID uint, extraClaims map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(a.cfg.ExpireDuration).Unix(),
	}
	for k, v := range extraClaims {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.cfg.Secret))
}

// parseToken 解析JWT令牌
func (a *Auth) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.cfg.Secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// extractToken 从请求头提取Token
func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:6]) == "BEARER" {
		return bearerToken[7:]
	}
	return bearerToken
}

// abortWithError 统一错误响应
func abortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"error":   message,
		"code":    code,
		"request": c.FullPath(),
	})
}
