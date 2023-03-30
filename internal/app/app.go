package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"techodom/internal/config"
	"techodom/internal/repository"
	"techodom/internal/transport/http"
	"techodom/internal/usecase"
	"techodom/pkg/cache/inmemory"
	"techodom/pkg/postgres"
	"time"
)

func Run() {
	godotenv.Load()
	cfg := config.New()

	var (
		logger *zap.Logger
		err    error
	)

	if cfg.App.Mode == "debug" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}

	e := echo.New()

	pool, err := postgres.NewConnectionPool(cfg.DBConfig.Host, cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Dbname, cfg.DBConfig.Port)
	if err != nil {
		logger.Error("error while connecting to database", zap.Error(err))
		return
	}
	defer pool.Close()

	if pool == nil {
		logger.Error("pool is nil")
		return
	}

	cache := inmemory.NewCache(time.Duration(cfg.Cache.TTL) * time.Second)
	redRepo := repository.NewRedirectRepository(pool, logger, cache)
	adminRepo := repository.NewAdmin(pool, logger)

	redUse := usecase.NewUserRedirect(redRepo, logger)
	adminUse := usecase.NewAdmin(adminRepo, logger)

	http.NewUserRedirect(redUse, logger).Register(e)
	http.NewAdminRedirect(adminUse, logger, cfg.AdminMW.Key).Register(e)
	// graceful shutdown
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.App.Port)); err != nil {
			logger.Error("error while starting server", zap.Error(err))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-quit

	if err := e.Shutdown(nil); err != nil {
		logger.Error("error while shutting down server", zap.Error(err))
	}

	logger.Info("server stopped")
}
