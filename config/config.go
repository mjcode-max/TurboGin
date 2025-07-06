package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const Version = "0.0.1"

// Config 全局配置结构体（自动生成文档）
type Config struct {
	Env        string           `mapstructure:"ENV" json:"env" yaml:"env" comment:"运行环境: dev/test/prod"`
	Version    string           `mapstructure:"VERSION" json:"version" yaml:"version"`
	Server     ServerConfig     `mapstructure:"SERVER" json:"server" yaml:"server"`
	Database   DatabaseConfig   `mapstructure:"DATABASE" json:"database" yaml:"database"`
	Redis      RedisConfig      `mapstructure:"REDIS" json:"redis" yaml:"redis"`
	JWT        AuthConfig       `mapstructure:"JWT" json:"jwt" yaml:"jwt"`
	Log        LogConfig        `mapstructure:"LOG" json:"log" yaml:"log"`
	Middleware MiddlewareConfig `mapstructure:"MIDDLEWARE" json:"middleware" yaml:"middleware"`
}

// ServerConfig HTTP服务配置
type ServerConfig struct {
	Host           string        `mapstructure:"HOST" json:"host" yaml:"host" validate:"required,hostname"`
	Port           int           `mapstructure:"PORT" json:"port" yaml:"port" validate:"required,min=1,max=65535"`
	ReadTimeout    time.Duration `mapstructure:"READ_TIMEOUT" json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"WRITE_TIMEOUT" json:"write_timeout" yaml:"write_timeout"`
	EnableSwagger  bool          `mapstructure:"ENABLE_SWAGGER" json:"enable_swagger" yaml:"enable_swagger"`
	TrustedProxies []string      `mapstructure:"TRUSTED_PROXIES" json:"trusted_proxies" yaml:"trusted_proxies"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Enabled      bool   `mapstructure:"ENABLED" json:"enabled" yaml:"enabled"`
	Driver       string `mapstructure:"DRIVER" json:"driver" yaml:"driver" validate:"oneof=mysql postgres sqlite"`
	DSN          string `mapstructure:"DSN" json:"dsn" yaml:"dsn" validate:"required_if=Enabled true"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS" json:"max_open_conns" yaml:"max_open_conns"`
	LogLevel     string `mapstructure:"LOG_LEVEL" json:"log_level" yaml:"log_level" validate:"oneof=silent error warn info"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Enabled  bool   `mapstructure:"ENABLED" json:"enabled" yaml:"enabled"`
	Addr     string `mapstructure:"ADDR" json:"addr" yaml:"addr" validate:"required_if=Enabled true"`
	Password string `mapstructure:"PASSWORD" json:"password" yaml:"password"`
	DB       int    `mapstructure:"DB" json:"db" yaml:"db"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Enabled        bool          `mapstructure:"ENABLED" json:"enabled" yaml:"enabled"`
	Secret         string        `mapstructure:"SECRET" json:"secret" yaml:"secret" validate:"required,min=32"`
	ExpireDuration time.Duration `mapstructure:"EXPIRE_DURATION" json:"expire_duration" yaml:"expire_duration"`
	Issuer         string        `mapstructure:"ISSUER" json:"issuer" yaml:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"LEVEL" json:"level" yaml:"level" validate:"oneof=debug info warn error"`
	Format     string `mapstructure:"FORMAT" json:"format" yaml:"format" validate:"oneof=json text"`
	Output     string `mapstructure:"OUTPUT" json:"output" yaml:"output" validate:"oneof=stdout file both"`
	MaxSize    int    `mapstructure:"MAX_SIZE" json:"max_size" yaml:"max_size"` // MB
	MaxBackups int    `mapstructure:"MAX_BACKUPS" json:"max_backups" yaml:"max_backups"`
	Compress   bool   `mapstructure:"COMPRESS" json:"compress" yaml:"compress"`
	MaxAge     int    `mapstructure:"MAX_AGE" json:"max_age" yaml:"max_age"`
	Dir        string `mapstructure:"DIR" json:"dir" yaml:"dir"`
	Filename   string `mapstructure:"FILENAME" json:"filename" yaml:"filename"`
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	CORS       CORSConfig      `mapstructure:"CORS" json:"cors" yaml:"cors"`
	RateLimit  RateLimitConfig `mapstructure:"RATE_LIMIT" json:"rate_limit" yaml:"rate_limit"`
	Prometheus bool            `mapstructure:"PROMETHEUS" json:"prometheus" yaml:"prometheus"`
}

type CORSConfig struct {
	Enabled          bool     `mapstructure:"ENABLED" json:"enabled" yaml:"enabled"`
	AllowOrigins     []string `mapstructure:"ALLOW_ORIGINS" json:"allow_origins" yaml:"allow_origins"`
	AllowMethods     []string `mapstructure:"ALLOW_METHODS" json:"allow_methods" yaml:"allow_methods"`
	AllowHeaders     []string `mapstructure:"ALLOW_HEADERS" json:"allow_headers" yaml:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"EXPOSE_HEADERS" json:"expose_headers" yaml:"expose_headers"`
	AllowCredentials bool     `mapstructure:"ALLOW_CREDENTIALS" json:"allow_credentials" yaml:"allow_credentials"`
}

type RateLimitConfig struct {
	Enabled bool    `mapstructure:"ENABLED" json:"enabled" yaml:"enabled"`
	RPS     float64 `mapstructure:"RPS" json:"rps" yaml:"rps"`
	Burst   int     `mapstructure:"BURST" json:"burst" yaml:"burst"`
}

// Load 加载配置
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// 环境变量配置
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	setDefaults(v)

	// 读取配置
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 解析到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 配置校验
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("配置校验失败: %w", err)
	}

	cfg.Version = Version
	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// 服务器默认值
	v.SetDefault("ENV", "dev")
	v.SetDefault("SERVER.HOST", "0.0.0.0")
	v.SetDefault("SERVER.PORT", 8080)
	v.SetDefault("SERVER.READ_TIMEOUT", 30*time.Second)
	v.SetDefault("SERVER.WRITE_TIMEOUT", 30*time.Second)
	v.SetDefault("SERVER.ENABLE_SWAGGER", true)

	// 数据库默认值
	v.SetDefault("DATABASE.ENABLED", true)
	v.SetDefault("DATABASE.DRIVER", "mysql")
	v.SetDefault("DATABASE.MAX_IDLE_CONNS", 10)
	v.SetDefault("DATABASE.MAX_OPEN_CONNS", 100)
	v.SetDefault("DATABASE.LOG_LEVEL", "warn")

	// Redis默认值
	v.SetDefault("REDIS.ENABLED", false)
	v.SetDefault("REDIS.DB", 0)

	// JWT默认值
	v.SetDefault("JWT.EXPIRE_DURATION", 72*time.Hour)
	v.SetDefault("JWT.ISSUER", "myapp")

	// 日志默认值
	v.SetDefault("LOG.LEVEL", "info")
	v.SetDefault("LOG.FORMAT", "json")
	v.SetDefault("LOG.OUTPUT", "stdout")
	v.SetDefault("LOG.MAX_SIZE", 100)
	v.SetDefault("LOG.MAX_BACKUPS", 7)

	// 中间件默认值
	v.SetDefault("MIDDLEWARE.CORS.ENABLED", true)
	v.SetDefault("MIDDLEWARE.CORS.ALLOW_METHODS", []string{"GET", "POST", "PUT", "DELETE"})
	v.SetDefault("MIDDLEWARE.RATE_LIMIT.ENABLED", true)
	v.SetDefault("MIDDLEWARE.RATE_LIMIT.RPS", 100)
	v.SetDefault("MIDDLEWARE.RATE_LIMIT.BURST", 50)
}

func validateConfig(cfg *Config) error {
	// 校验DSN格式
	if cfg.Database.Enabled {
		if _, err := url.Parse(cfg.Database.DSN); err != nil {
			return fmt.Errorf("数据库DSN格式错误: %w", err)
		}
	}

	// 校验JWT密钥长度
	if len(cfg.JWT.Secret) < 32 {
		return fmt.Errorf("JWT密钥长度至少32位")
	}

	return nil
}
