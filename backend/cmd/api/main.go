package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/FortiBrine/VoidShift/internal/config"
	"github.com/FortiBrine/VoidShift/internal/http/middleware"
	"github.com/FortiBrine/VoidShift/internal/http/routes"
	"github.com/labstack/echo/v5"
)

func main() {
	cfg := config.Load()

	e := echo.New()

	middleware.Register(e)
	routes.Register(e, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	startConfig := echo.StartConfig{
		Address: ":8080",
	}

	if err := startConfig.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
