package web

import (
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/helper"
	"html/template"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{s}
}

var respCode = fiber.StatusInternalServerError

func (c *Controller) HomePage(ctx *fiber.Ctx) error {
	data := fiber.Map{
		"showNavbar": true,
		"title":      "Home",
	}
	tpl := template.Must(template.ParseFiles("web/template/home.gohtml", "web/template/parts/navbar.gohtml", "web/template/base.gohtml"))

	fctx := &helper.FiberCtx{Fctx: ctx}
	respCode = fiber.StatusOK

	fctx.Fctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return tpl.ExecuteTemplate(fctx.Fctx.Response().BodyWriter(), "base.gohtml", data)
}

func (c *Controller) ErrorPage(ctx *fiber.Ctx) error {
	tpl := template.Must(template.ParseFiles("web/template/error.gohtml", "web/template/parts/navbar.gohtml", "web/template/base.gohtml"))

	fctx := &helper.FiberCtx{Fctx: ctx}
	respCode = fiber.StatusOK

	fctx.Fctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return tpl.ExecuteTemplate(fctx.Fctx.Response().BodyWriter(), "base.gohtml", fiber.Map{})
}

func (c *Controller) UnauthorisedPage(ctx *fiber.Ctx) error {
	tpl := template.Must(template.ParseFiles("web/template/unauthorised.gohtml", "web/template/parts/navbar.gohtml", "web/template/base.gohtml"))

	fctx := &helper.FiberCtx{Fctx: ctx}
	respCode = fiber.StatusOK

	fctx.Fctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return tpl.ExecuteTemplate(fctx.Fctx.Response().BodyWriter(), "base.gohtml", fiber.Map{"showNavbar": true,"errMessage":"insufficient permission"})
}
