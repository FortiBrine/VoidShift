package middleware

import (
	"github.com/FortiBrine/VoidShift/internal/shared/logger"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Register(
	e *echo.Echo,
) {
	e.Use(logger.RequestLogger(logger.RequestLoggerConfig{
		Logger: e.Logger,
		Skipper: func(c *echo.Context) bool {
			switch c.Request().URL.Path {
			case "/health":
				return true
			default:
				return false
			}
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
}
