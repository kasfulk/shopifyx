package routes

import (
	"shopifyx/api/handlers"
	"shopifyx/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App, h handlers.Product) {
	g := app.Group("/v1/product")
	g.Get("", middleware.OptionalJWTAuth(), h.GetProducts)
	g.Get("/:id", h.GetProductDetail)
	g.Post("/:id/buy", middleware.JWTAuth(), h.BuyProduct)
	g.Post("/:id/stock", middleware.JWTAuth(), h.UpdateStock)
	g.Post("", middleware.JWTAuth(), h.AddProduct)
	g.Patch("/:id", middleware.JWTAuth(), h.UpdateProduct)
	g.Delete("/:id", middleware.JWTAuth(), h.DeleteProduct)
}
