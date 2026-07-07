package webui

import (
	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/config"
	"github.com/FortiBrine/VoidShift/internal/middleware"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/FortiBrine/VoidShift/view/pages"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func RegisterRoutes(
	app *fiber.App,
	authService *auth.Service,
	wireGuardService *wireguard.Service,
	env config.Environment,
) {
	handler := NewHandler(authService, wireGuardService)

	app.Get("/static*", static.New("static", static.Config{
		FS: StaticFS,
	}))

	registerComponentScripts(app, env.IsDev())

	app.Get("/", handler.Home)
	app.Get("/login", handler.LoginPage)
	app.Post("/login", handler.Login)
	app.Post("/logout", handler.Logout)

	wg := app.Group("/wireguard").Use(middleware.NewAuth(middleware.AuthConfig{
		Unauthorized: func(c fiber.Ctx) error {
			return c.Redirect().To("/login")
		},
	}))
	wg.Get("/", handler.Networks)
	wg.Get("/create-network", handler.NetworkCreatePage)
	wg.Post("/create-network", handler.NetworkCreate)
	wg.Get("/networks/:id", handler.NetworkDetail)
	wg.Post("/networks/:id/up", handler.NetworkUp)
	wg.Post("/networks/:id/down", handler.NetworkDown)
	wg.Get("/networks/:id/peers/create", handler.PeerCreatePage)
	wg.Post("/networks/:id/peers/create", handler.PeerCreate)
	wg.Get("/peers/:peerId/config", handler.PeerConfig)
	wg.Get("/peers/:peerId/qr", handler.PeerQR)

	app.Use(func(c fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return Render(c, pages.NotFound(localizer(c)))
	})
}
