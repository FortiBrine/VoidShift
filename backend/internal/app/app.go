package app

import (
	"context"
	"time"

	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/shared/config"
	"github.com/FortiBrine/VoidShift/internal/shared/database"
	"github.com/FortiBrine/VoidShift/internal/shared/http"
	"github.com/FortiBrine/VoidShift/internal/shared/http/middleware"
	"github.com/FortiBrine/VoidShift/internal/shared/http/router"
	"github.com/FortiBrine/VoidShift/internal/shared/http/validator"
	"github.com/FortiBrine/VoidShift/internal/shared/logger"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/labstack/echo/v5"
	"golang.zx2c4.com/wireguard/wgctrl"
	"gorm.io/gorm"
)

type App struct {
	Echo             *echo.Echo
	DB               *gorm.DB
	SessionService   *session.Service
	UserService      *user.Service
	AuthService      *auth.Service
	WireGuardService *wireguard.Service
	WireGuardClient  *wgctrl.Client
}

func NewApp(ctx context.Context, cfg config.Config) (*App, error) {
	l := logger.New(cfg.Environment)

	db, err := database.Open(cfg, l)
	if err != nil {
		return nil, err
	}

	sessionRepository := session.NewGormRepository(db)
	sessionService := session.NewService(sessionRepository, 5*24*time.Hour)
	if err := sessionService.Load(); err != nil {
		return nil, err
	}

	userRepository := user.NewGormRepository(db)
	userService := user.NewService(userRepository)
	if err := userService.Load(ctx, cfg); err != nil {
		return nil, err
	}

	authService := auth.NewService(sessionService, userService)

	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}

	wireGuardRepository := wireguard.NewGormRepository(db)
	wireGuardService := wireguard.NewService(wireGuardRepository, client)
	if err := wireGuardService.Load(); err != nil {
		_ = client.Close()
		return nil, err
	}

	e := echo.New()
	e.Logger = l
	e.Validator = validator.NewCustomValidator()
	e.HTTPErrorHandler = http.CustomErrorHandler
	middleware.Register(e)
	router.RegisterRoutes(e, userService, sessionService, authService, wireGuardService)

	return &App{
		Echo:             e,
		DB:               db,
		SessionService:   sessionService,
		UserService:      userService,
		AuthService:      authService,
		WireGuardService: wireGuardService,
		WireGuardClient:  client,
	}, nil
}

func (a *App) Close() error {
	if a.WireGuardClient != nil {
		return a.WireGuardClient.Close()
	}
	return nil
}
