package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/config"
	"github.com/FortiBrine/VoidShift/internal/middleware"
	"github.com/FortiBrine/VoidShift/internal/store"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/FortiBrine/VoidShift/internal/validator"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/memory/v2"
	"github.com/valyala/fasthttp"
)

type App struct {
	fiber            *fiber.App
	db               *sql.DB
	userService      *user.Service
	authService      *auth.Service
	wireGuardService *wireguard.Service
}

func NewApp(
	ctx context.Context, cfg config.Config,
	l *slog.Logger,
) (app *App, err error) {
	db, err := store.Open(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer func() {
		if err != nil {
			_ = db.Close()
		}
	}()
	if err = store.Migrate(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	sessionConfig := session.Config{
		Storage:           memory.New(),
		IdleTimeout:       30 * time.Minute,
		AbsoluteTimeout:   24 * time.Hour,
		CookieHTTPOnly:    true,
		CookiePath:        "/",
		CookieSecure:      !cfg.Environment.IsDev(),
		CookieSessionOnly: true,
		CookieSameSite:    "Lax",
		Extractor:         extractors.FromCookie("voidshift_session"),
	}

	userRepository := user.NewSqlcRepository(db)
	userService := user.NewService(userRepository)
	if err = userService.Load(ctx, cfg); err != nil {
		return nil, fmt.Errorf("failed to load user service: %w", err)
	}

	authService := auth.NewService(userService)

	wireGuardRepository := wireguard.NewSqlcRepository(db)
	wireGuardService, err := wireguard.NewService(wireGuardRepository, cfg.HostAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create wireguard service: %w", err)
	}
	defer func() {
		if err != nil {
			_ = wireGuardService.Close()
		}
	}()
	if err = wireGuardService.Load(); err != nil {
		return nil, fmt.Errorf("failed to load wireguard service: %w", err)
	}

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler:    middleware.NewCustomErrorHandler(l),
		StructValidator: validator.NewCustomValidator(),
		CaseSensitive:   true,
		ProxyHeader:     fasthttp.HeaderXForwardedFor,
	})

	fiberApp.Hooks().OnListen(func(listenData fiber.ListenData) error {
		l.Info("Server is starting",
			"host", listenData.Host,
			"port", listenData.Port,
			"tls", listenData.TLS,
		)
		return nil
	})

	middleware.Register(fiberApp, l, cfg, sessionConfig)
	RegisterRoutes(fiberApp, authService, wireGuardService, cfg.Environment.IsDev())

	app = new(App{
		fiber:            fiberApp,
		db:               db,
		userService:      userService,
		authService:      authService,
		wireGuardService: wireGuardService,
	})
	return
}

func (app *App) Start(ctx context.Context, cfg config.Config) error {
	if err := app.fiber.Listen(cfg.HttpAddress, fiber.ListenConfig{
		GracefulContext:       ctx,
		ShutdownTimeout:       cfg.GracefulTimeout,
		DisableStartupMessage: true,
	}); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start fiber app: %w", err)
	}
	return nil
}

func (app *App) Close() error {
	if app.wireGuardService != nil {
		if err := app.wireGuardService.Close(); err != nil {
			return fmt.Errorf("failed to close wireguard service: %w", err)
		}
	}
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}
	return nil
}
