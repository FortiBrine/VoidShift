package middleware

import (
	"log/slog"

	"github.com/FortiBrine/VoidShift/internal/config"
	"github.com/FortiBrine/VoidShift/internal/logger"
	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Register(
	app *fiber.App,
	l *slog.Logger,
	cfg config.Config,
	sessionConfig session.Config,
) {
	app.Use(requestid.New())
	app.Use(logger.NewRequestLogger(logger.Config{
		Logger: l,
	}))
	app.Use(recoverer.New())

	app.Use(session.New(sessionConfig))
}
