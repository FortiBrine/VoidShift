package wireguard

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(
	router fiber.Router,
	wireGuardService *Service,
) {
	handler := NewHandler(wireGuardService)
	router.Get("/networks", handler.GetNetworks)
	router.Post("/networks/generate", handler.GenerateNetwork)
	router.Post("/networks/:id/up", handler.UpNetwork)
	router.Post("/networks/:id/down", handler.DownNetwork)
	router.Get("/networks/:id", handler.GetNetwork)
	router.Delete("/networks/:id", handler.RemoveNetwork)
	router.Post("/networks/:id/peers/generate", handler.GeneratePeer)
	router.Get("/peers/:peerId/config", handler.GetPeerConfig)
	router.Get("/peers/:peerId/config/download", handler.DownloadPeerConfig)
	router.Get("/peers/:peerId/qr", handler.GetPeerConfigQR)
	router.Delete("/peers/:peerId", handler.RemovePeer)
}
