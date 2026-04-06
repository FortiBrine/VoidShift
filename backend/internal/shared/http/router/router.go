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

type Router struct {
	sessionService   *session.Service
	userService      *user.Service
	wireGuardService *wireguard.Service
}

func NewRouter(
	sessionService *session.Service,
	userService *user.Service,
	wireGuardService *wireguard.Service,
) *Router {
	return &Router{
		sessionService:   sessionService,
		userService:      userService,
		wireGuardService: wireGuardService,
	}
}

func (r *Router) Register(e *echo.Echo) {
	authMiddleware := auth.Middleware(r.sessionService, r.userService)

	api := e.Group("/api")
	api.GET("/health", handlers.Health)

	a := api.Group("/auth")
	a.POST("/login", auth.NewLoginHandler(r.sessionService, r.userService).Login)

	protected := api.Group("/protected")
	protected.Use(authMiddleware)
	protected.GET("/test", handlers.TestHandler)

	wgHandler := wireguard.NewHandler(r.wireGuardService)
	wgGroup := api.Group("/vpn/wireguard")
	wgGroup.Use(authMiddleware)
	wgGroup.GET("/networks", wgHandler.GetNetwork)
	wgGroup.POST("/networks/generate", wgHandler.GenerateNetwork)

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
