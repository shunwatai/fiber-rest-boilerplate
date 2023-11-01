package server

import (
	"fmt"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/modules/todo"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Fiber struct {
	App *fiber.App
}

func (f *Fiber) GetApp() {
	f.App = fiber.New()
}

func (f *Fiber) LoadMiddlewares() {
	f.App.Use(logger.New())
	f.App.Use(recover.New())
}

func (f *Fiber) LoadAllRoutes() {
	api := f.App.Group("/api")
	todo.GetRoutes(api)
}

func (f *Fiber) Start() {
	config := config.Cfg
	config.LoadEnvVariables()
	config.WatchConfig()

	fmt.Printf("server config: %+v\n", config.ServerConf)
	fmt.Printf("db config: %+v\n", config.DbConf)

	log.Fatal(f.App.Listen(fmt.Sprintf(":%s", config.ServerConf.Port)))
}

var Api = &Fiber{}
