package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-clean-architecture/application/handler"
	"github.com/ntdat104/go-clean-architecture/application/middleware"
	"github.com/ntdat104/go-clean-architecture/application/service"
	"github.com/ntdat104/go-clean-architecture/config"
	"github.com/ntdat104/go-clean-architecture/infra/repository"
	"github.com/ntdat104/go-clean-architecture/pkg/logger"
	"go.uber.org/zap"
)

const ServiceName = "go-clean-architecture"

// Constants for application settings
const (
	// DefaultShutdownTimeout is the default timeout for graceful shutdown
	DefaultShutdownTimeout = 5 * time.Second
	// DefaultMetricsAddr is the default address for the metrics server
	DefaultMetricsAddr = ":9090"
)

func main() {
	fmt.Println("Starting " + ServiceName)

	// Initialize configuration
	config.Init("./config", "config")
	fmt.Println("Configuration initialized")

	// Initialize logging
	logger.Init()
	logger.Logger.Info("Application starting",
		zap.String("service", ServiceName),
		zap.String("env", string(config.GlobalConfig.Env)))

	// Initialize metrics collection system
	middleware.InitializeMetrics()
	logger.Logger.Info("Metrics collection system initialized")

	repository.InitializeRepositories(
		repository.WithMySQLite(),
		repository.WithMySQL(),
		repository.WithMyPostgreSQL(),
		repository.WithRedis(),
	)

	// Start metrics server in a separate goroutine if enabled
	if config.GlobalConfig.MetricsServer != nil && config.GlobalConfig.MetricsServer.Enabled {
		metricsAddr := config.GlobalConfig.MetricsServer.Addr
		if metricsAddr == "" {
			metricsAddr = DefaultMetricsAddr
		}
		go func() {
			if err := middleware.StartMetricsServer(metricsAddr); err != nil {
				logger.Logger.Error("Failed to start metrics server", zap.Error(err))
			}
		}()
		logger.Logger.Info("Metrics server started", zap.String("address", metricsAddr))
	} else {
		logger.Logger.Info("Metrics server is disabled")
	}

	// Create context and cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.ZapLoggerWithBody())

	systemService := service.NewSystemService()
	handler.NewSystemHandler(router, systemService)

	srv := &http.Server{
		Addr:    config.GlobalConfig.HTTPServer.Addr,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("%v started on http://%v%v", config.GlobalConfig.App.Name, "localhost", config.GlobalConfig.HTTPServer.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(config.GlobalConfig.App.Name+" failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
