package qrcode

import (
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/helper"
	"golang-api-starter/internal/helper/logger/zap_log"

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
	logger.Debugf("qrcode ctrl GetQrcodeContentFromPdf\n")
	fctx := &helper.FiberCtx{Fctx: ctx}

	form, err := fctx.Fctx.MultipartForm()
	if err != nil { /* handle error */
		logger.Debugf("failed to get multipartForm, err: %+v\n", err.Error())
		return err
	}

	result, err := c.service.GetQrcodeContentFromPdf(form)
	if err != nil {
		logger.Debugf("GetQrcodeContentFromPdf err: %+v\n", err.Error())
		return err
	}

	respCode = fiber.StatusOK
	return fctx.JsonResponse(
		respCode,
		map[string]interface{}{"data": result},
	)
}
