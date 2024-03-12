package webservices

import (
	"context"
	"log"

	"shopifyx/api/handlers"
	"shopifyx/api/responses"
	"shopifyx/api/routes"
	"shopifyx/configs"
	"shopifyx/db/connections"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Run() {
	app := fiber.New()

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	dbPool, err := connections.NewPgConn(config)
	if err != nil {
		log.Fatalf("failed open connection to db: %v", err)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("FAILED PING TO DB: %v", err)
	}

	deps := handlers.Dependencies{
		Cfg:    config,
		DbPool: dbPool,
	}

	// load Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// register route in another package
	routes.RouteRegister(app, deps)

	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return responses.ReturnTheResponse(c, true, int(404), "Not Found", nil)
	})

	// Here we go!
	log.Fatalln(app.Listen(":" + config.APPPort))
}
