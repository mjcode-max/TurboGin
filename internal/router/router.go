package router

import (
	"TurboGin/internal/controller"
	"TurboGin/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(ctl *controller.Container, auth *middleware.Auth) func(*gin.Engine) {
	return func(engine *gin.Engine) {
		// ==================== 公共路由 ====================
		publicGroup := engine.Group("/v1")
		{
			publicGroup.POST("/register", ctl.User.CreateUser)
		}

		// ==================== 受保护路由 ====================
		privateGroup := engine.Group("/v1")
		privateGroup.Use(auth.Middleware())
		{
			userRoutes := privateGroup.Group("/users")
			{
				userRoutes.GET("/:id", ctl.User.GetUser)
			}
		}
	}
}
