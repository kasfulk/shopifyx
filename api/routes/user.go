package routes

import (
	"shopifyx/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	handlers := handlers.User{}
	g := app.Group("/user")
	g.Post("/register", handlers.Register)
	g.Post("/login", handlers.Login)
}
