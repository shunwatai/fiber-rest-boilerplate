package web

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

	router.Get("/home", custMiddleware.JwtChecker.CheckJwt(), ctrl.HomePage)
	router.Get("/error", ctrl.ErrorPage)
}
