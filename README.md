# TurboGin - Go 快速开发脚手架

## 项目概述

TurboGin 是一个基于 Gin + GORM + Wire + Viper 的高性能 Go 后端开发脚手架，专为快速构建 RESTful API 服务设计。它集成了现代 Web 开发所需的常用组件，包括：

- 🚀 Gin HTTP 高性能框架
- 🗃️ GORM ORM 数据库操作
- 🔌 Wire 依赖注入
- ⚙️ Viper 配置管理
- 🔐 JWT 认证授权
- 📝 Zap 结构化日志
- 🔄 Redis 缓存支持
- ⏱️ 请求限流
- 🌐 CORS 跨域支持
- 🛡️ IP 白名单控制
- 🩺 健康检查端点

## 快速开始

### 安装要求

1. Go 1.24+ (推荐 1.24 或更高版本)
2. MySQL 8.0+ (或兼容数据库)
3. Redis (可选)

### 安装步骤

```bash
# 1. 克隆项目
git clone https://github.com/coder/TurboGin.git
cd TurboGin

# 2. 初始化项目
make init

# 3. 编辑配置文件
vi config.yaml

# 4. 运行项目
make run
```

### 测试运行

启动成功后，访问健康检查端点：
```bash
curl http://localhost:8080/api/health
```

预期响应：
```json
{"status":"healthy","version":"0.0.1"}
```

## 主要项目结构

```
TurboGin/
├── config.yaml                 # 主配置文件
├── go.mod                      # Go 模块定义
├── internal/
│   ├── controller/             # 控制器层
│   ├── dao/                    # 数据访问对象
│   ├── model/                  # 数据模型
│   ├── service/                # 业务逻辑层
│   └── wire/                   # 依赖注入配置
├── pkg/
│   ├── config/                 # 配置加载
│   ├── db/                     # 数据库连接
│   ├── logger/                 # 日志系统
│   ├── middleware/             # 中间件
│   ├── redis/                  # Redis 客户端
│   └── server/                 # HTTP 服务器
```

## 配置说明

配置文件位于项目根目录下的 `config.yaml`，支持以下配置项：

### 基础配置
```yaml
ENV: "dev"  # 运行环境: dev/test/prod
```

### 服务器配置
```yaml
SERVER:
  HOST: "0.0.0.0"
  PORT: 8080
  READ_TIMEOUT: 30s
  WRITE_TIMEOUT: 30s
  TRUSTED_PROXIES: # IP 白名单
    - "127.0.0.1"
    - "10.0.0.0/8"
```

### 数据库配置
```yaml
DATABASE:
  ENABLED: true
  DRIVER: "mysql"
  DSN: "root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  MAX_IDLE_CONNS: 10
  MAX_OPEN_CONNS: 100
```

### Redis 配置
```yaml
REDIS:
  ENABLED: false
  ADDR: "localhost:6379"
  PASSWORD: ""
  DB: 0
```

### JWT 配置
```yaml
JWT:
  ENABLED: true
  SECRET: "your-32-byte-long-secret-key-here-123456"
  EXPIRE_DURATION: 72h
```

### 日志配置
```yaml
LOG:
  LEVEL: "info"
  FORMAT: "console" # console/json
  OUTPUT: "both"    # stdout/file/both
  MAX_SIZE: 100     # MB
  MAX_BACKUPS: 7    # 保留日志文件数
```

### 中间件配置
```yaml
MIDDLEWARE:
  CORS:
    ENABLED: true
    ALLOW_ORIGINS: ["*"]
  RATE_LIMIT:
    ENABLED: true
    RPS: 100.0  # 每秒请求数
    BURST: 50   # 突发流量
```

## 功能组件

### 1. 依赖注入 (Wire)

项目使用 Wire 进行依赖注入管理：

- `wire/wire.go` - 依赖绑定声明
- `wire/wire_gen.go` - 自动生成的依赖注入代码

添加新依赖步骤：
1. 在 `wire/wire.go` 中添加新绑定
2. 运行 `make generate`
3. 更新 `wire_gen.go`

### 2. 中间件系统

#### JWT 认证
```go
// 生成令牌
token, err := auth.GenerateToken(userID, map[string]interface{}{"role": "admin"})

// 中间件使用
privateGroup.Use(auth.Middleware())
```

#### CORS 跨域
```go
// 配置示例
cors := middleware.NewCORS(config)
engine.Use(cors.Middleware())
```

#### IP 白名单
```go
ipAccess := middleware.NewIPAccess(config)
engine.Use(ipAccess.Middleware())
```

#### 请求限流
```go
rateLimiter := middleware.NewRateLimiter(config)
engine.Use(rateLimiter.Middleware())
```

#### 请求日志
```go
requestLog := middleware.NewRequestLog(logger)
engine.Use(requestLog.Middleware())
```

### 3. 日志系统

使用 Zap 实现高性能结构化日志：

```go
// 日志记录示例
log.Info("User created", 
    logger.String("username", "john"),
    logger.Int("user_id", 123))

// 错误记录
log.Error("Database error", logger.Error(err))
```

日志级别支持：
- debug
- info
- warn
- error

### 4. 数据库操作

使用 GORM 进行数据库操作：

```go
// DAO 示例
type UserDAO struct {
    IBaseDAO[model.User]
}

func NewUserDAO(db *gorm.DB) IUserDAO {
    return &UserDAO{
        IBaseDAO: NewBaseDAO[model.User](db),
    }
}

// 自定义查询
func (u UserDAO) FindByName(name string) ([]model.User, error) {
    var users []model.User
    err := u.DB().Where("name = ?", name).Find(&users).Error
    return users, err
}
```

### 5. Redis 客户端

```go
// 初始化
redisClient, err := redis.New(config, logger)

// 使用示例
ctx := context.Background()
err := redisClient.GetClient().Set(ctx, "key", "value", 10*time.Minute).Err()
```

## 添加新功能

### 添加新控制器

1. 在 `internal/controller` 创建新控制器文件
2. 实现控制器方法
```go
package controller

type ProductController struct {
    productService service.IProductService
}

func (c *ProductController) GetProduct(ctx *gin.Context) {
    // 控制器逻辑
}
```

3. 在 `internal/controller/container.go` 中添加控制器
```go
type Container struct {
    User    *UserController
    Product *ProductController // 添加新控制器
}
```

### 添加新服务

1. 在 `internal/service` 创建服务文件
```go
package service

type IProductService interface {
    GetProduct(id uint) (*model.Product, error)
}

type ProductService struct {
    productDao dao.IProductDAO
}

func NewProductService(productDao dao.IProductDAO) IProductService {
    return &ProductService{productDao: productDao}
}
```

2. 在 `wire/wire.go` 中添加服务绑定
```go
var serviceSet = wire.NewSet(
    service.NewUserService,
    service.NewProductService, // 添加新服务
)
```

### 添加新数据模型

在 `internal/model` 创建模型文件：
```go
package model

import "gorm.io/gorm"

type Product struct {
    gorm.Model
    Name  string
    Price float64
}
```

### 添加新路由

在 `router/router.go` 中注册新路由：
```go
func RegisterRoutes(ctl *controller.Container, auth *middleware.Auth) func(*gin.Engine) {
    return func(engine *gin.Engine) {
        // ...
        productGroup := engine.Group("/products")
        {
            productGroup.GET("/:id", ctl.Product.GetProduct)
        }
    }
}
```

## 健康检查

项目内置健康检查端点：

```
GET /health
```

响应示例：
```json
{
    "status": "healthy",
    "version": "0.0.1"
}
```

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目仓库
2. 创建新分支 (`git checkout -b feature/your-feature`)
3. 提交代码 (`git commit -am 'Add some feature'`)
4. 推送到分支 (`git push origin feature/your-feature`)
5. 创建 Pull Request

## 下一步计划

- [ ] 添加 proto 支持
- [ ] 集成 Prometheus 监控
- [ ] 增加分布式追踪
- [ ] 添加单元测试示例
- [ ] 支持多数据库类型

---

**TurboGin** 致力于提供简洁高效的 Go 后端开发体验，让开发者专注于业务逻辑实现。欢迎使用并提出宝贵意见！