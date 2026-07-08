package webui

import (
	"github.com/FortiBrine/VoidShift/internal/auth"
	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(
	router fiber.Router,
	authService *auth.Service,
	i18nService *i18n.Service,
) {
	handler := NewHandler(authService, i18nService)

	router.Get("/login", handler.LoginPage)
	router.Post("/login", handler.Login)
	router.Post("/logout", handler.Logout)
}
