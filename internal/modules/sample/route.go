package sample

import (
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

var (
	cfg  = config.Cfg
	Srvc = &Service{}
	ctrl = &Controller{}
)

func GetRoutes(router fiber.Router, custMiddleware *middleware.CustomMiddlewares) {
	Srvc = NewService()
	ctrl = NewController(Srvc)

	router.Get("/ping", ctrl.Ping)
	router.Get("/hallo", ctrl.HalloPage)
	router.Post("/test-email", ctrl.SendEmail)
}
