package routes

import (
	"shopifyx/api/handlers"
	"shopifyx/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App, h handlers.Product) {
	g := app.Group("/v1/product").Use(middleware.JWTAuth())
	g.Post("/:id/buy", h.BuyProduct)
	g.Post("/:id/stock", h.UpdateStock)
	g.Post("", h.AddProduct)
	g.Patch("/:id", h.UpdateProduct)
}
