package webui

import (
	"embed"

	"github.com/FortiBrine/VoidShift/internal/config"
	"github.com/FortiBrine/VoidShift/internal/i18n"
	"github.com/FortiBrine/VoidShift/view/pages"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

//go:embed static/*
var StaticFS embed.FS

func RegisterRoutes(app *fiber.App, env config.Environment) {
	app.Get("/static*", static.New("static", static.Config{
		FS: StaticFS,
	}))

	registerComponentScripts(app, env.IsDev())
}

func RegisterNotFoundHandler(app *fiber.App, i18nService *i18n.Service) {
	app.Use(func(c fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return Render(c, pages.NotFound(Localizer(i18nService, c)))
	})
}

func Render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func Localizer(
	service *i18n.Service,
	c fiber.Ctx,
) *i18n.Localizer {
	return service.FromAcceptLanguage(c.AcceptLanguage())
}
