package middleware

import (
	"github.com/FortiBrine/VoidShift/internal/shared/logger"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Register(
	e *echo.Echo,
) {
	e.Use(logger.RequestLogger(e.Logger))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
}
