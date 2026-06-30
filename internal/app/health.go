package app

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func Health(c fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(map[string]string{
		"status": "ok",
	})
}
