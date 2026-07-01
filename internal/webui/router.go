package webui

import (
	"github.com/FortiBrine/VoidShift/view/pages"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func RegisterRoutes(
	router fiber.Router,
) {
	router.Get("/static*", static.New("static", static.Config{
		FS: StaticFS,
	}))

	router.Get("/", func(c fiber.Ctx) error {
		return Render(c, pages.Home())
	})

	router.Use(NotFoundMiddleware)
}
