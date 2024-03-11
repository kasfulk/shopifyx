package routes

import (
	"shopifyx/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userHandler handlers.User) {
	g := app.Group("/v1/user")
	g.Post("/register", userHandler.Register)
	g.Post("/login", userHandler.Login)
}
