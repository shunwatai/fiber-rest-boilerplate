package main

import (
	"golang-api-starter/internal/modules/todo"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")
	todo.GetRoutes(api)

	log.Fatal(app.Listen(":7000"))
}
