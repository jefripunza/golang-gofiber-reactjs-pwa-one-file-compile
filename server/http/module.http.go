package http

import (
	"backend/server/module"

	"github.com/gofiber/fiber/v2"
)

func Module(app *fiber.App) {
	api := app.Group("/api")

	// --------------------------
	// --------------------------

	Example := module.Example{}
	Example.Route(api)

	// --------------------------
	// --------------------------

}
