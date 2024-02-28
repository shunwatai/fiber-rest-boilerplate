package oauth

import (
	"golang-api-starter/internal/helper"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

var respCode = fiber.StatusInternalServerError

func OAuthGetAuth(ctx *fiber.Ctx) error {
	respCode = fiber.StatusOK
	fctx := &helper.FiberCtx{Fctx: ctx}
	return fctx.JsonResponse(
		respCode,
		fiber.Map{"data": "OAuthGetAuth"},
	)
}

func OAuthGetUser(ctx *fiber.Ctx) error {
	respCode = fiber.StatusOK
	fctx := &helper.FiberCtx{Fctx: ctx}
	// try to get the user without re-authenticating
	if gothUser, err := goth_fiber.CompleteUserAuth(ctx); err == nil {
		// t, _ := template.New("foo").Parse(userTemplate)
		// t.Execute(res, gothUser)
		return fctx.JsonResponse(
			respCode,
			fiber.Map{"data": gothUser},
		)
	} else {
		goth_fiber.BeginAuthHandler(ctx)
	}
	return nil
}

func OAuthLogout(ctx *fiber.Ctx) error {
	respCode = fiber.StatusOK
	fctx := &helper.FiberCtx{Fctx: ctx}
	return fctx.JsonResponse(
		respCode,
		fiber.Map{"data": "OAuthLogout"},
	)
}

func OAuthProviderPage(ctx *fiber.Ctx) error {
	tpl := template.Must(template.ParseGlob("web/template/oauth/sign-in.tmpl"))

	fctx := &helper.FiberCtx{Fctx: ctx}
	respCode = fiber.StatusOK

	fctx.Fctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return tpl.Execute(fctx.Fctx.Response().BodyWriter(), struct{}{})
}
