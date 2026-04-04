package router

import (
	"io/fs"
	"strings"

	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/embed"
	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/shared/http/handlers"
	"github.com/FortiBrine/VoidShift/internal/user"
	"github.com/labstack/echo/v5"
)

type Router struct {
	sessionService *session.Service
	userService    *user.Service
}

func NewRouter(
	sessionService *session.Service,
	userService *user.Service,
) *Router {
	return &Router{
		sessionService: sessionService,
		userService:    userService,
	}
}

func (r *Router) Register(e *echo.Echo) {
	api := e.Group("/api")
	api.GET("/health", handlers.Health)

	a := api.Group("/auth")
	a.POST("/login", auth.NewLoginHandler(r.sessionService, r.userService).Login)

	protected := api.Group("/protected")
	protected.Use(auth.Middleware(r.sessionService, r.userService))
	protected.GET("/test", handlers.TestHandler)

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
