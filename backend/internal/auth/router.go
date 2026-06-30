package auth

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(
	router fiber.Router,
	authService *Service,
) {
	handler := NewHandler(authService)
	router.Post("/login", handler.Login)
	router.Post("/logout", handler.Logout)
}
