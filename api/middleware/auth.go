package middleware

import (
	"shopifyx/configs"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func JwtSign(app *fiber.App) fiber.Handler {
	config, err := configs.LoadConfig()
	if err != nil {
		// Handling error when loading config, you might want to return an error here
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
	}

	// Initialize JWT middleware
	jwt := jwtware.New(jwtware.Config{
		SigningKey: []byte(config.JWTSecret),
	})

	// Return the middleware handler function
	return func(c *fiber.Ctx) error {
		// Apply JWT middleware to the context
		err := jwt(c)
		if err != nil {
			// If JWT validation fails, set response status to 403 Forbidden
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}
		// If JWT validation succeeds, continue to the next handler
		return c.Next()
	}
}
