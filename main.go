package main

import (
	"fmt"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/modules/todo"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config := config.Config{}
	config.LoadEnvVariables()
	config.WatchConfig()
	fmt.Printf("server config: %+v\n", config.ServerConf)
	fmt.Printf("db config: %+v\n", config.DbConf)

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")
	todo.GetRoutes(api)

	log.Fatal(app.Listen(fmt.Sprintf(":%s",config.ServerConf.Port)))
}
