package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/config"
	"github.com/FortiBrine/VoidShift/internal/database"
	apphttp "github.com/FortiBrine/VoidShift/internal/http"
	"github.com/FortiBrine/VoidShift/internal/http/middleware"
	"github.com/FortiBrine/VoidShift/internal/http/router"
	"github.com/FortiBrine/VoidShift/internal/http/validator"
	"github.com/FortiBrine/VoidShift/internal/user"
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
) (*App, error) {
	db, err := database.Open(cfg)
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(ctx, db); err != nil {
		return nil, err
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
	if err := userService.Load(ctx, cfg); err != nil {
		return nil, err
	}

	authService := auth.NewService(userService)

	wireGuardRepository := wireguard.NewSqlcRepository(db)
	wireGuardService, err := wireguard.NewService(wireGuardRepository, cfg.HostAddress)
	if err != nil {
		return nil, err
	}
	if err := wireGuardService.Load(); err != nil {
		_ = wireGuardService.Close()
		return nil, err
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:    apphttp.NewCustomErrorHandler(l),
		StructValidator: validator.NewCustomValidator(),
		CaseSensitive:   true,
		ProxyHeader:     fasthttp.HeaderXForwardedFor,
	})

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		l.Info("Server is starting",
			"host", listenData.Host,
			"port", listenData.Port,
			"tls", listenData.TLS,
		)
		return nil
	})

	middleware.Register(app, l, cfg, sessionConfig)
	router.RegisterRoutes(app, authService, wireGuardService)

	return &App{
		fiber:            app,
		db:               db,
		userService:      userService,
		authService:      authService,
		wireGuardService: wireGuardService,
	}, nil
}

func (app *App) Start(ctx context.Context, cfg config.Config) error {
	if err := app.fiber.Listen(cfg.HttpAddress, fiber.ListenConfig{
		GracefulContext:       ctx,
		ShutdownTimeout:       cfg.GracefulTimeout,
		DisableStartupMessage: true,
	}); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (app *App) Close() error {
	if app.wireGuardService != nil {
		if err := app.wireGuardService.Close(); err != nil {
			return err
		}
	}
	if app.db != nil {
		return app.db.Close()
	}
	return nil
}
