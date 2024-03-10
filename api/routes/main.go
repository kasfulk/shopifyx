package routes

import (
	"shopifyx/configs"

	"github.com/gofiber/fiber/v2"
)

func RouteRegister(app *fiber.App) {
	config := configs.LoadConfig()
	ver := app.Group("/" + config.Server.Version)

	ver.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
}
