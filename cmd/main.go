package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evythrossell/account-management-api/config"
	httpadapter "github.com/evythrossell/account-management-api/internal/adapters/http"
	logger "github.com/evythrossell/account-management-api/internal/adapters/logger"
	"github.com/evythrossell/account-management-api/internal/container"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	appLogger := logger.NewSimpleLogger(logger.InfoLevel)

	ctr, err := container.New(cfg, appLogger)
	if err != nil {
		appLogger.Fatal("failed to initialize container", logger.Err(err))
	}
	defer ctr.Close()

	router := httpadapter.SetupRouter(ctr.AccountHandler(), ctr.HealthHandler())

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("server error", logger.Err(err))
			os.Exit(1)
		}
	}()

	appLogger.Info("server started", logger.String("port", cfg.ServerPort))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("shutdown signal received")

	ctxShut, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShut); err != nil {
		appLogger.Error("shutdown failed", logger.Err(err))
		os.Exit(1)
	}

	appLogger.Info("server stopped successfully")
}
