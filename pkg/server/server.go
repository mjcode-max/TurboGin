package server

import (
	"context"
	"fmt"
	"github.com/mjcode-max/TurboGin/config"
	"github.com/mjcode-max/TurboGin/internal/controller"
	"github.com/mjcode-max/TurboGin/pkg/logger"
	"github.com/mjcode-max/TurboGin/pkg/middleware"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server represents the main server structure
type Server struct {
	engine *gin.Engine
	http   *http.Server
	db     *gorm.DB
	cfg    *config.Config
	log    *logger.Logger

	// Middleware
	middlewares struct {
		auth      *middleware.Auth
		cors      *middleware.CORS
		rateLimit *middleware.RateLimiter
		allowed   *middleware.IPAccess
	}

	// Controllers
	controllers *controller.Container
}

// New creates a new Server instance (dependency injection entry point)
func New(
	cfg *config.Config,
	db *gorm.DB,
	log *logger.Logger,
	auth *middleware.Auth,
	cors *middleware.CORS,
	rateLimit *middleware.RateLimiter,
	allowed *middleware.IPAccess,
	controllers *controller.Container,
	registerRoutes func(*gin.Engine),
) *Server {
	s := &Server{
		db:          db,
		cfg:         cfg,
		log:         log,
		controllers: controllers,
	}

	s.middlewares.auth = auth
	s.middlewares.cors = cors
	s.middlewares.rateLimit = rateLimit
	s.middlewares.allowed = allowed

	s.initializeEngine(registerRoutes)
	s.configureHTTPServer()

	return s
}

// initializeEngine sets up the Gin engine with middleware and routes
func (s *Server) initializeEngine(registerRoutes func(*gin.Engine)) {
	setGinMode(s.cfg.Env)
	s.engine = gin.Default()

	s.engine.Use(s.middlewares.allowed.Middleware())

	// Apply middleware
	if s.cfg.Middleware.CORS.Enabled {
		s.engine.Use(s.middlewares.cors.Middleware())
	}

	if s.cfg.Middleware.RateLimit.Enabled {
		s.engine.Use(s.middlewares.rateLimit.Middleware())
	}

	// Configure JSON prefix
	s.engine.SecureJsonPrefix("api")

	// Register routes
	registerRoutes(s.engine)

	// Register health check endpoint
	s.engine.GET("/health", s.healthCheck)
}

// configureHTTPServer sets up the HTTP server configuration
func (s *Server) configureHTTPServer() {
	s.http = &http.Server{
		Addr:           s.resolveAddress(),
		Handler:        s.engine,
		ReadTimeout:    s.cfg.Server.ReadTimeout,
		WriteTimeout:   s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}
}

// Run starts the server with graceful shutdown
func (s *Server) Run() error {
	// Start HTTP server in a goroutine
	go func() {
		s.log.Info("Server starting",
			zap.String("address", s.http.Addr),
			zap.String("env", s.cfg.Env),
		)

		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal("Server start failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Begin graceful shutdown
	return s.shutdown()
}

// shutdown performs graceful server shutdown
func (s *Server) shutdown() error {
	s.log.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	// Clean up other resources
	s.cleanup()
	return nil
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	if err := s.db.Exec("SELECT 1").Error; err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": s.cfg.Version,
	})
}

// cleanup releases server resources
func (s *Server) cleanup() {
	if sqlDB, err := s.db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	s.log.Info("Server resources released")
}

// -------------------------- Helper Methods --------------------------

// setGinMode sets Gin mode based on environment
func setGinMode(env string) {
	switch env {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

// resolveAddress constructs the server address from config
func (s *Server) resolveAddress() string {
	return s.cfg.Server.Host + ":" + strconv.Itoa(s.cfg.Server.Port)
}
