package router

import (
	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/auth/middleware"
	"github.com/FortiBrine/VoidShift/internal/http/handlers"
	"github.com/FortiBrine/VoidShift/internal/webui"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func RegisterRoutes(
	app *fiber.App,
	authService *auth.Service,
	wireGuardService *wireguard.Service,
) {
	api := app.Group("/api")
	api.Get("/health", handlers.Health)

	auth.RegisterRoutes(app.Group("/auth"), authService)
	wireguard.RegisterRoutes(
		app.Group("/vpn/wireguard").
			Use(middleware.New()),
		wireGuardService,
	)

	app.Get("/", static.New("", static.Config{
		FS: webui.FS,
	}))

	app.Get("/*", func(c fiber.Ctx) error {
		return c.SendFile("/index.html")
	})
}
