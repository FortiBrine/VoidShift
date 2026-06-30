package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type AuthConfig struct {
	Next func(c fiber.Ctx) bool
}

var ConfigDefault = AuthConfig{
	Next: nil,
}

func configDefault(config ...AuthConfig) AuthConfig {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	return cfg
}

func NewAuth(config ...AuthConfig) fiber.Handler {
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
