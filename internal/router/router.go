package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mjcode-max/TurboGin/internal/controller"
	"github.com/mjcode-max/TurboGin/pkg/middleware"
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
