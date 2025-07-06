package controller

import (
	"TurboGin/internal/service"
)

// Container 集中管理所有控制器
type Container struct {
	User *UserController

	// 添加其他控制器...
}

// NewContainer 构造函数（依赖所有需要的Service）
func NewContainer(
	userService service.IUserService,

	// 其他Service...
) *Container {
	return &Container{
		User: NewUserController(userService),
	}
}
