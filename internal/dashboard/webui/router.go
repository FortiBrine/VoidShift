package webui

import (
	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/FortiBrine/VoidShift/internal/middleware"
	"github.com/FortiBrine/VoidShift/internal/wireguard"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(
	router fiber.Router,
	wireGuardService *wireguard.Service,
	i18nService *i18n.Service,
	hostAddress string,
) {
	handler := NewHandler(wireGuardService, i18nService, hostAddress)

	router.Get("/", middleware.NewAuth(middleware.AuthConfig{
		Unauthorized: func(c fiber.Ctx) error {
			return c.Redirect().To("/login")
		},
	}), handler.Home)
}
