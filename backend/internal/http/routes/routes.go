package routes

import (
	"github.com/FortiBrine/VoidShift/internal/config"
	"github.com/FortiBrine/VoidShift/internal/embed"
	"github.com/FortiBrine/VoidShift/internal/http/handlers"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Echo, cfg config.Config) {
	api := e.Group("/api")
	api.GET("/health", handlers.Health)

	e.StaticFS("/", echo.MustSubFS(embed.WebuiFiles, "webui"))
}
