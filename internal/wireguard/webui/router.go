package webui

import (
	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(
	router fiber.Router,
	wireGuardService *wireguard.Service,
	i18nService *i18n.Service,
) {
	handler := NewHandler(wireGuardService, i18nService)

	router.Get("/", handler.Networks)
	router.Get("/create-network", handler.NetworkCreatePage)
	router.Post("/create-network", handler.NetworkCreate)
	router.Get("/networks/:id", handler.NetworkDetail)
	router.Post("/networks/:id/up", handler.NetworkUp)
	router.Post("/networks/:id/down", handler.NetworkDown)
	router.Get("/networks/:id/peers/create", handler.PeerCreatePage)
	router.Post("/networks/:id/peers/create", handler.PeerCreate)
	router.Get("/peers/:peerId/config", handler.PeerConfig)
	router.Get("/peers/:peerId/qr", handler.PeerQR)
}
