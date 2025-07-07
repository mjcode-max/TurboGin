//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/mjcode-max/TurboGin/config"
	"github.com/mjcode-max/TurboGin/internal/controller"
	"github.com/mjcode-max/TurboGin/internal/dao"
	"github.com/mjcode-max/TurboGin/internal/router"
	"github.com/mjcode-max/TurboGin/internal/service"
	"github.com/mjcode-max/TurboGin/pkg/db"
	"github.com/mjcode-max/TurboGin/pkg/logger"
	"github.com/mjcode-max/TurboGin/pkg/middleware"
	"github.com/mjcode-max/TurboGin/pkg/redis"
	"github.com/mjcode-max/TurboGin/pkg/server"

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
