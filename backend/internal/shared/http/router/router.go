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
	userService *user.Service,
	sessionService *session.Service,
	authService *auth.Service,
	wireGuardService *wireguard.Service,
) {
	authMiddleware := auth.Middleware(sessionService, userService)

	api := e.Group("/api")
	api.GET("/health", handlers.Health)

	a := api.Group("/auth")
	loginHandler := auth.NewLoginHandler(authService)
	a.POST("/login", loginHandler.Login)
	a.POST("/logout", loginHandler.Logout)

	wgHandler := wireguard.NewHandler(wireGuardService)
	wgGroup := api.Group("/vpn/wireguard")
	wgGroup.Use(authMiddleware)
	wgGroup.GET("/networks", wgHandler.GetNetworks)
	wgGroup.POST("/networks/generate", wgHandler.GenerateNetwork)
	wgGroup.POST("/networks/:id/up", wgHandler.UpNetwork)
	wgGroup.POST("/networks/:id/down", wgHandler.DownNetwork)
	wgGroup.GET("/networks/:id", wgHandler.GetNetwork)
	wgGroup.DELETE("/networks/:id", wgHandler.RemoveNetwork)
	wgGroup.POST("/networks/:id/peers/generate", wgHandler.GeneratePeer)
	wgGroup.GET("/peers/:peerId/config", wgHandler.GetPeerConfig)
	wgGroup.GET("/peers/:peerId/config/download", wgHandler.DownloadPeerConfig)
	wgGroup.GET("/peers/:peerId/qr", wgHandler.GetPeerConfigQR)
	wgGroup.DELETE("/peers/:peerId", wgHandler.RemovePeer)

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
