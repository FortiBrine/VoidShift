package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FortiBrine/VoidShift/internal/config"
	"github.com/FortiBrine/VoidShift/internal/database"
	"github.com/FortiBrine/VoidShift/internal/http/middleware"
	"github.com/FortiBrine/VoidShift/internal/http/routes"
	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	cfg := config.Load()
	db, err := database.NewSqliteDatabase(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err != nil {
		log.Printf("failed to initialize database: %v", err)
		return
	}

	sessionRepository := session.NewGormRepository(db)
	sessionService := session.NewService(sessionRepository, 5*24*time.Hour)

	if err := sessionService.Load(); err != nil {
		log.Printf("failed to load session service: %v", err)
		return
	}

	userRepository := user.NewGormRepository(db)
	userService := user.NewService(userRepository)

	if err := userService.Load(ctx, cfg); err != nil {
		log.Printf("failed to load user service: %v", err)
		return
	}

	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	middleware.Register(e)

	router := routes.NewRouter(sessionService, userService)
	router.Register(e)

	startConfig := echo.StartConfig{
		Address:         cfg.HttpAddress,
		GracefulTimeout: cfg.GracefulTimeout,
	}

	if err := startConfig.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
