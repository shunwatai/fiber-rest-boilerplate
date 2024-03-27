package sample

import (
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/config"
)

var (
	cfg  = config.Cfg
	Srvc = &Service{}
	ctrl = &Controller{}
)

func GetRoutes(router fiber.Router) {
	Srvc = NewService()
	ctrl = NewController(Srvc)

	router.Get("/hallo", ctrl.HalloPage)
	router.Get("/ping", ctrl.Ping)
	router.Post("/test-email", ctrl.SendEmail)
}
