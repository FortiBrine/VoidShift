package router

import (
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

	e.StaticFS("/", echo.MustSubFS(embed.WebuiFiles, "webui"))
}
