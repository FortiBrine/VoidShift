package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type Config struct {
	Next func(c fiber.Ctx) bool
}

var ConfigDefault = Config{
	Next: nil,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	return cfg
}

func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	return func(c fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		sess := session.FromContext(c)
		if sess.Get("username") == nil {
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.Next()
	}
}
