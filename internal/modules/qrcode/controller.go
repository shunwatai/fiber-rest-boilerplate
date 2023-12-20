package qrcode

import (
	"fmt"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/helper"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) Controller {
	return Controller{s}
}

var cfg = config.Cfg
var respCode = fiber.StatusInternalServerError

func (c *Controller) GetQrcodeContentFromPdf(ctx *fiber.Ctx) error {
	fmt.Printf("qrcode ctrl GetQrcodeContentFromPdf\n")
	fctx := &helper.FiberCtx{Fctx: ctx}

	form, err := fctx.Fctx.MultipartForm()
	if err != nil { /* handle error */
		fmt.Printf("failed to get multipartForm, err: %+v\n", err.Error())
		return err
	}

	result, err := c.service.GetQrcodeContentFromPdf(form)
	if err != nil {
		fmt.Printf("GetQrcodeContentFromPdf err: %+v\n", err.Error())
		return err
	}

	respCode = fiber.StatusOK
	return fctx.JsonResponse(
		respCode,
		map[string]interface{}{"data": result},
	)
}
