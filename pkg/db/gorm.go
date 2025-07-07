package db

import (
	"fmt"
	"github.com/mjcode-max/TurboGin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	if !cfg.Database.Enabled {
		return nil, nil
	}

	// 初始化GORM
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent), // 生产环境可改为Warn
		PrepareStmt: true,                                  // 开启预编译
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// 健康检查
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

// HealthCheck 健康检查API
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
