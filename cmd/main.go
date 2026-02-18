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
	"github.com/evythrossell/account-management-api/internal/adapters/logger"
	"github.com/evythrossell/account-management-api/internal/container"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("critical failure loading config: " + err.Error())
	}

	appLogger := logger.NewSimpleLogger(logger.InfoLevel)
	ctr, err := container.New(cfg, appLogger)
	if err != nil {
		appLogger.Fatal("failed to initialize container", logger.Err(err))
	}
	defer ctr.Close()

	router := httpadapter.SetupRouter(
		ctr.AccountHandler(),
		ctr.HealthHandler(),
		ctr.TransactionHandler(),
	)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("server failed to start", logger.Err(err))
		}
	}()

	appLogger.Info("server started", logger.String("port", cfg.ServerPort))
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	appLogger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("server forced to shutdown", logger.Err(err))
		os.Exit(1)
	}

	appLogger.Info("server exited gracefully")
}
