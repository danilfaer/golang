package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/danilfaer/golang/order/internal/app"
	"github.com/danilfaer/golang/order/internal/config"
	"github.com/danilfaer/golang/platform/pkg/closer"
	"github.com/danilfaer/golang/platform/pkg/logger"
)

const (
	shutdownTimeout = 10 * time.Second
	configPath      = "deploy/compose/order/.env"
)

func main() {
	// Загружаем переменные из .env файла
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, fmt.Sprintf("failed to create app order: %v", err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, fmt.Sprintf("failed to run app order: %v", err))
		return
	}

}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, fmt.Sprintf("failed to close all: %v", err))
	}
}
