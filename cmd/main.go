package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-clean-architecture/config"
	"github.com/ntdat104/go-clean-architecture/infra/repo"
	"github.com/ntdat104/go-clean-architecture/internal/handler"
	"github.com/ntdat104/go-clean-architecture/internal/middleware"
	"github.com/ntdat104/go-clean-architecture/internal/service"
	"github.com/ntdat104/go-clean-architecture/pkg/logger"
)

func main() {
	config.InitConfig("config/config.yml")
	logger.InitProduction("./logs/")
	defer logger.Sync()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.ZapLoggerWithBody())

	cfg := config.GetGlobalConfig()
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler: router,
	}

	// db, err := repo.NewDB(repo.DatabaseConfig{
	// 	Driver:                  "mysql",
	// 	Url:                     "user:password@tcp(127.0.0.1:3306)/your_database_name?charset=utf8mb4&parseTime=true&loc=UTC&tls=false&readTimeout=3s&writeTimeout=3s&timeout=3s&clientFoundRows=true",
	// 	ConnMaxLifetimeInMinute: 3,
	// 	MaxOpenConns:            10,
	// 	MaxIdleConns:            1,
	// })

	db, err := repo.NewSQLiteDB(repo.SQLiteConfig{
		// Path:                  "./app.db", // or ":memory:" for tests
		Path:                  ":memory:",
		SchemaPath:            "./schema/schema.sql",
		ConnMaxLifetimeMinute: 10,
		MaxOpenConns:          10,
		MaxIdleConns:          5,
	})
	if err != nil {
		log.Fatalf("failed to new database err=%s\n", err.Error())
	}

	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	handler.NewUserHandler(router, userService)

	systemService := service.NewSystemService()
	handler.NewSystemHandler(router, systemService)

	// Run server in a goroutine
	go func() {
		log.Printf("%v started on http://%v:%v", cfg.App.Name, cfg.HTTP.Host, strconv.Itoa(cfg.HTTP.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(cfg.App.Name+" failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
