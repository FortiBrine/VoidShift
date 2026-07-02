package webui

import (
	"path"
	"strings"

	"github.com/FortiBrine/VoidShift/view/components/ui"
	"github.com/gofiber/fiber/v3"
)

func registerComponentScripts(app *fiber.App, isDevelopment bool) {
	app.Get("/templui/js/:file", func(c fiber.Ctx) error {
		fileName := path.Base(c.Params("file"))
		component := strings.TrimSuffix(strings.TrimSuffix(fileName, ".min.js"), ".js")

		file, err := ui.ScriptFiles.ReadFile(path.Join(component, fileName))
		if err != nil {
			return fiber.ErrNotFound
		}

		c.Set(fiber.HeaderContentType, "text/javascript; charset=utf-8")
		if isDevelopment {
			c.Set(fiber.HeaderCacheControl, "no-store")
		} else {
			c.Set(fiber.HeaderCacheControl, "public, max-age=31536000")
		}
		return c.Send(file)
	})
}
