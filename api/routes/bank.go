package routes

import (
	"shopifyx/api/handlers"
	"shopifyx/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func BankRoutes(app *fiber.App, h handlers.BankHandler) {
	g := app.Group("/v1/bank").Use(middleware.JWTAuth())

	g.Post("/account", h.Create)
}
