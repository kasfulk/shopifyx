package routes

import (
	"shopifyx/api/handlers"
	"shopifyx/db/functions"

	"github.com/gofiber/fiber/v2"
)

func RouteRegister(app *fiber.App, deps handlers.Dependencies) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	userHandler := handlers.User{
		Database: functions.NewUser(deps.DbPool, deps.Cfg),
	}

	UserRoutes(app, userHandler)

	productHandler := handlers.Product{
		Database:     functions.NewProductFn(deps.DbPool),
		UserDatabase: functions.NewUser(deps.DbPool, deps.Cfg),
		BankDatabase: functions.NewBank(deps.DbPool),
	}

	ProductRoutes(app, productHandler)

	imageUploaderHandler := handlers.ImageUploader{
		Uploader: functions.NewImageUploader(deps.Cfg),
	}

	ImageRoutes(app, imageUploaderHandler)

	bankAccountHandler := handlers.BankHandler{
		Bank: *functions.NewBank(deps.DbPool),
	}

	BankRoutes(app, bankAccountHandler)
}
