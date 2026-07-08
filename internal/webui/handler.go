package webui

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func IsAuthenticated(c fiber.Ctx) bool {
	return session.FromContext(c).Get("username") != nil
}
