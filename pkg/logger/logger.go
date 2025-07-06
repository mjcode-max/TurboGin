package logger

import (
	"TurboGin/config"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 确保所有zap方法可用
var (
	Any        = zap.Any
	Bool       = zap.Bool
	Float64    = zap.Float64
	Time       = zap.Time
	Strings    = zap.Strings
	Ints       = zap.Ints
	Errors     = zap.Errors
	NamedError = zap.NamedError
	Object     = zap.Object
	Reflect    = zap.Reflect
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// Logger 封装zap.Logger并提供更友好的API
type Logger struct {
	*zap.Logger
	cfg *config.LogConfig
}

// New 构造函数（线程安全）
func New(cfg *config.Config) (*Logger, error) {
	var initErr error
	once.Do(func() {
		core, err := buildCore(&cfg.Log)
		if err != nil {
			initErr = err
			return
		}

		// 构建Logger
		globalLogger = zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	})

	if initErr != nil {
		return nil, fmt.Errorf("failed to init logger: %w", initErr)
	}

	return &Logger{
		Logger: globalLogger,
		cfg:    &cfg.Log,
	}, nil
}

// buildCore 构建日志核心
func buildCore(cfg *config.LogConfig) (zapcore.Core, error) {
	// 解析日志级别
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// 编码器配置（生产环境优化）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder, // 终端彩色输出
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 选择编码器
	var encoder zapcore.Encoder
	switch cfg.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return nil, fmt.Errorf("unsupported log format: %s", cfg.Format)
	}

	// 多输出源
	cores := make([]zapcore.Core, 0, 2)

	// 控制台输出
	if cfg.Output == "stdout" || cfg.Output == "both" {
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			level,
		))
	}

	// 文件输出
	if cfg.Output == "file" || cfg.Output == "both" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(cfg.Dir, cfg.Filename),
			MaxSize:    cfg.MaxSize,    // MB
			MaxBackups: cfg.MaxBackups, // 保留旧日志文件数
			MaxAge:     cfg.MaxAge,     // 保留天数
			Compress:   cfg.Compress,   // 是否压缩
		})

		cores = append(cores, zapcore.NewCore(
			encoder,
			fileWriter,
			level,
		))
	}

	// 组合多个Core
	return zapcore.NewTee(cores...), nil
}

// Sync 刷新缓冲区的日志
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// WithFields 结构化日志
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		cfg:    l.cfg,
	}
}

// ==================== 辅助函数 ====================

// Field 快捷生成字段
type Field = zap.Field

func String(key, val string) Field                 { return zap.String(key, val) }
func Int(key string, val int) Field                { return zap.Int(key, val) }
func Error(err error) Field                        { return zap.Error(err) }
func Duration(key string, val time.Duration) Field { return zap.Duration(key, val) }
