package webservices

import (
	"log"
	"time"

	"shopifyx/api/responses"
	"shopifyx/api/routes"
	"shopifyx/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Run() {
	app := fiber.New()
	loc, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		app.Use(func(c *fiber.Ctx) error {
			return responses.ReturnTheResponse(c, true, int(500), "Can not init the timezone", nil)
		})
	}
	time.Local = loc // -> this is setting the global timezone

	config := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// load Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// register route in another package
	routes.RouteRegister(app)

	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return responses.ReturnTheResponse(c, true, int(404), "Not Found", nil)
	})

	// Here we go!
	log.Fatalln(app.Listen(config.Server.Host + ":" + config.Server.Port))
}
