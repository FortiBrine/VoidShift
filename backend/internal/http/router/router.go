package router

import (
	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/auth/middleware"
	"github.com/FortiBrine/VoidShift/internal/embed"
	"github.com/FortiBrine/VoidShift/internal/http/handlers"
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

	a := app.Group("/auth")
	loginHandler := auth.NewLoginHandler(authService)
	a.Post("/login", loginHandler.Login)
	a.Post("/logout", loginHandler.Logout)

	wgHandler := wireguard.NewHandler(wireGuardService)
	wgGroup := app.Group("/vpn/wireguard")
	wgGroup.Use(middleware.New())
	wgGroup.Get("/networks", wgHandler.GetNetworks)
	wgGroup.Post("/networks/generate", wgHandler.GenerateNetwork)
	wgGroup.Post("/networks/:id/up", wgHandler.UpNetwork)
	wgGroup.Post("/networks/:id/down", wgHandler.DownNetwork)
	wgGroup.Get("/networks/:id", wgHandler.GetNetwork)
	wgGroup.Delete("/networks/:id", wgHandler.RemoveNetwork)
	wgGroup.Post("/networks/:id/peers/generate", wgHandler.GeneratePeer)
	wgGroup.Get("/peers/:peerId/config", wgHandler.GetPeerConfig)
	wgGroup.Get("/peers/:peerId/config/download", wgHandler.DownloadPeerConfig)
	wgGroup.Get("/peers/:peerId/qr", wgHandler.GetPeerConfigQR)
	wgGroup.Delete("/peers/:peerId", wgHandler.RemovePeer)

	app.Get("/", static.New("", static.Config{
		FS: embed.WebuiFiles,
	}))

	app.Get("/*", func(c fiber.Ctx) error {
		return c.SendFile("/index.html")
	})
}
