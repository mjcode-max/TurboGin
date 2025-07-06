//go:build wireinject
// +build wireinject

package wire

import (
	"TurboGin/config"
	"TurboGin/internal/controller"
	"TurboGin/internal/dao"
	"TurboGin/internal/router"
	"TurboGin/internal/service"
	"TurboGin/pkg/db"
	"TurboGin/pkg/logger"
	"TurboGin/pkg/middleware"
	"TurboGin/pkg/redis"
	"TurboGin/pkg/server"

	"github.com/google/wire"
)

var daoSet = wire.NewSet(
	dao.NewUserDAO,
)

var serviceSet = wire.NewSet(
	service.NewUserService,
)

var controllerSet = wire.NewSet(
	controller.NewContainer,
)

var routerSet = wire.NewSet(router.RegisterRoutes)

var middlewareSet = wire.NewSet(
	middleware.NewCORS,
	middleware.NewAuth,
	middleware.NewRateLimiter,
	middleware.NewRequestLog,
	middleware.NewIPAccess,
)

var systemSet = wire.NewSet(config.Load, db.NewGormDB, logger.New, redis.New, server.New)

func InitApp() (*server.Server, func(), error) {
	wire.Build(
		systemSet,
		middlewareSet,
		daoSet,
		serviceSet,
		controllerSet,
		routerSet,
	)
	return &server.Server{}, nil, nil
}
