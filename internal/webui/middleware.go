package webui

import (
	"github.com/FortiBrine/VoidShift/view/pages"
	"github.com/gofiber/fiber/v3"
)

func NotFoundMiddleware(c fiber.Ctx) error {
	c.Status(fiber.StatusNotFound)
	return Render(c, pages.NotFound(localizer(c)))
}
