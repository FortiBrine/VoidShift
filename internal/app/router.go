package app

import (
	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/middleware"
	"github.com/FortiBrine/VoidShift/internal/webui"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(
	app *fiber.App,
	authService *auth.Service,
	wireGuardService *wireguard.Service,
) {
	api := app.Group("/api")
	api.Get("/health", Health)

	auth.RegisterRoutes(api.Group("/auth"), authService)
	wireguard.RegisterRoutes(
		api.Group("/vpn/wireguard").
			Use(middleware.NewAuth()),
		wireGuardService,
	)

	webui.RegisterRoutes(app, authService, wireGuardService)
}
