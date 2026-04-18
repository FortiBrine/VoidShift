//go:build linux

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FortiBrine/VoidShift/internal/app"
	"github.com/FortiBrine/VoidShift/internal/shared/config"
	"github.com/labstack/echo/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("error loading config: %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	application, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Printf("failed to initialize app: %v", err)
		return
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("failed to close app: %v", err)
		}
	}()

	startConfig := echo.StartConfig{
		Address:         cfg.HttpAddress,
		GracefulTimeout: cfg.GracefulTimeout,
	}

	if err := startConfig.Start(ctx, application.Echo); err != nil {
		application.Echo.Logger.Error("failed to start server", "error", err)
	}
}
