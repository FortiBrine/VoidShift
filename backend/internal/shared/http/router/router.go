package router

import (
	"io/fs"
	"strings"

	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/embed"
	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/shared/http/handlers"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(
	e *echo.Echo,
	sessionService *session.Service,
	userService *user.Service,
	wireGuardService *wireguard.Service,
) {
	authMiddleware := auth.Middleware(sessionService, userService)

	api := e.Group("/api")
	api.GET("/health", handlers.Health)

	a := api.Group("/auth")
	a.POST("/login", auth.NewLoginHandler(sessionService, userService).Login)

	protected := api.Group("/protected")
	protected.Use(authMiddleware)
	protected.GET("/test", handlers.TestHandler)

	wgHandler := wireguard.NewHandler(wireGuardService)
	wgGroup := api.Group("/vpn/wireguard")
	wgGroup.Use(authMiddleware)
	wgGroup.GET("/networks", wgHandler.GetNetwork)
	wgGroup.POST("/networks/generate", wgHandler.GenerateNetwork)
	wgGroup.POST("/networks/:id/up", wgHandler.UpNetwork)
	wgGroup.POST("/networks/:id/down", wgHandler.DownNetwork)
	wgGroup.DELETE("/networks/:id", wgHandler.RemoveNetwork)
	wgGroup.POST("/networks/:id/peers/generate", wgHandler.GeneratePeer)

	webui := echo.MustSubFS(embed.WebuiFiles, "webui")

	e.GET("/*", func(c *echo.Context) error {
		path := strings.TrimPrefix(c.Request().URL.Path, "/")

		if strings.HasPrefix(path, "api/") {
			return echo.ErrNotFound
		}

		if path == "" {
			return c.FileFS("index.html", webui)
		}

		if _, err := fs.Stat(webui, path); err == nil {
			return c.FileFS(path, webui)
		}

		return c.FileFS("index.html", webui)
	})
}
