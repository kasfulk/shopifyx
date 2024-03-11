package routes

import (
	"shopifyx/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userHandler handlers.User) {
	// if you want to use the middleware, you can do it like this:
	// g.Post("/login", userHandler.Login, middleware.JwtSign(app))
	// or you can do it like this:
	// app.Group("/login", middleware.JwtSign(app), userHandler.Login)

	g := app.Group("/v1/user")
	g.Post("/register", userHandler.Register)
	g.Post("/login", userHandler.Login)
}
