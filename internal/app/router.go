package app

import (
	"github.com/FortiBrine/VoidShift/internal/auth"
	authwebui "github.com/FortiBrine/VoidShift/internal/auth/webui"
	"github.com/FortiBrine/VoidShift/internal/config"
	dashboardwebui "github.com/FortiBrine/VoidShift/internal/dashboard/webui"
	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/FortiBrine/VoidShift/internal/middleware"
	"github.com/FortiBrine/VoidShift/internal/webui"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	wireguardwebui "github.com/FortiBrine/VoidShift/internal/wireguard/webui"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(
	app *fiber.App,
	authService *auth.Service,
	wireGuardService *wireguard.Service,
	i18nService *i18n.Service,
	hostAddress string,
	env config.Environment,
) {
	api := app.Group("/api")
	api.Get("/health", Health)

	auth.RegisterRoutes(api.Group("/auth"), authService)
	wireguard.RegisterRoutes(
		api.Group("/vpn/wireguard").
			Use(middleware.NewAuth()),
		wireGuardService,
	)
	api.Use(func(c fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	webui.RegisterRoutes(app, env)
	authwebui.RegisterRoutes(app, authService, i18nService)
	dashboardwebui.RegisterRoutes(app, wireGuardService, i18nService, hostAddress)

	wg := app.Group("/wireguard").Use(middleware.NewAuth(middleware.AuthConfig{
		Unauthorized: func(c fiber.Ctx) error {
			return c.Redirect().To("/login")
		},
	}))
	wireguardwebui.RegisterRoutes(wg, wireGuardService, i18nService)

	webui.RegisterNotFoundHandler(app, i18nService)
}
