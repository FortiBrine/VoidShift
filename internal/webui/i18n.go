package webui

import (
	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/gofiber/fiber/v3"
)

func localizer(c fiber.Ctx) *i18n.Localizer {
	return i18n.FromAcceptLanguage(c.AcceptLanguage())
}
