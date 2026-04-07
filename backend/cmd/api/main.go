package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/shared/config"
	"github.com/FortiBrine/VoidShift/internal/shared/database"
	"github.com/FortiBrine/VoidShift/internal/shared/http"
	"github.com/FortiBrine/VoidShift/internal/shared/http/middleware"
	"github.com/FortiBrine/VoidShift/internal/shared/http/router"
	"github.com/FortiBrine/VoidShift/internal/shared/http/validator"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/labstack/echo/v5"
	"golang.zx2c4.com/wireguard/wgctrl"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	databaseInitialize := func(cfg config.Config) (*gorm.DB, error) {
		if cfg.SqliteDatabasePath != "" {
			return database.NewSqliteDatabase(cfg)
		}

		if cfg.MysqlDsn != "" {
			return database.NewMysqlDatabase(cfg)
		}

		return nil, errors.New("no database configured")
	}

	db, err := databaseInitialize(cfg)

	if err != nil {
		log.Printf("failed to initialize database: %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

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

	client, err := wgctrl.New()
	if err != nil {
		log.Printf("failed to initialize WireGuard client: %v", err)
		return
	}
	defer func(client *wgctrl.Client) {
		err := client.Close()
		if err != nil {
			log.Printf("failed to close WireGuard client: %v", err)
		}
	}(client)

	wireGuardRepository := wireguard.NewGormRepository(db)
	wireGuardService := wireguard.NewService(wireGuardRepository, client)

	if err := wireGuardService.Load(); err != nil {
		log.Printf("failed to load WireGuard service: %v", err)
		return
	}

	e := echo.New()

	e.Validator = validator.NewCustomValidator()
	e.HTTPErrorHandler = http.CustomErrorHandler

	middleware.Register(e)

	r := router.NewRouter(sessionService, userService, wireGuardService)
	r.Register(e)

	startConfig := echo.StartConfig{
		Address:         cfg.HttpAddress,
		GracefulTimeout: cfg.GracefulTimeout,
	}

	if err := startConfig.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "validator", err)
	}

}
