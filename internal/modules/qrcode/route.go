package qrcode

import (
	"github.com/gofiber/fiber/v2"
)

var tableName = "qrcodes"
var Srvc = NewService()
var ctrl = NewController(Srvc)

func GetRoutes(router fiber.Router) {
	// r := router.Group("/qrcodes", jwtcheck.CheckFromHeader())
	r := router.Group("/qrcodes")
	r.Post("/from-pdf", GetQrcodeContentFromPdf)
}


// QrcodeFromPdf godoc
//
//	@Summary		Get Qrcode content from pdf
//	@Description	Get Qrcode content from pdf
//	@Tags			qrcodes
//	@Accept			json
//	@Produce		json
//	@Param			Qrcode	body	string	true	"single Qrcode request json"	SchemaExample({ "task": "take shower", "done": false })
//	@Param			Qrcodes	body	string	true	"batch Qrcode request json"	SchemaExample([{ "task": "take shower", "done": false }, { "task": "go practice", "done": false }])
//	@Security		ApiKeyAuth
//	@Router			/qrcodes [post]
func GetQrcodeContentFromPdf(c *fiber.Ctx) error {
	return ctrl.GetQrcodeContentFromPdf(c)
}

