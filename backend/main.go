package main

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

//go:embed public/*
var embeddedFiles embed.FS

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// API
	api := e.Group("/api")
	api.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"ok": true,
		})
	})

	e.StaticFS("/", echo.MustSubFS(embeddedFiles, "public"))

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
