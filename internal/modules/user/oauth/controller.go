package oauth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
	"golang-api-starter/internal/helper"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"html/template"
)

var respCode = fiber.StatusInternalServerError

func OAuthGetAuth(ctx *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		logger.Errorf("goth_fiber.CompleteUserAuth err: %+v", err)
		return err
	}

	provider, _ := goth_fiber.GetProviderName(ctx)
	goth_fiber.StoreInSession(provider, user.AccessToken, ctx)

	// logger.Debugf("authed user: %+v", user)
	respCode = fiber.StatusOK
	fctx := &helper.FiberCtx{Fctx: ctx}
	return fctx.JsonResponse(
		respCode,
		fiber.Map{"data": user},
	)
}

func OAuthLogin(ctx *fiber.Ctx) error {
	return goth_fiber.BeginAuthHandler(ctx)
}

func OAuthLogout(ctx *fiber.Ctx) error {
	respCode = fiber.StatusOK
	fctx := &helper.FiberCtx{Fctx: ctx}
	if err := goth_fiber.Logout(ctx); err!=nil{
		logger.Errorf("failed to logout...")
	}
	return fctx.JsonResponse(
		respCode,
		fiber.Map{"data": "OAuthLogout"},
	)
}

func OAuthProviderPage(ctx *fiber.Ctx) error {
	tpl := template.Must(template.ParseGlob("web/template/oauth/sign-in.gohtml"))

	fctx := &helper.FiberCtx{Fctx: ctx}
	respCode = fiber.StatusOK

	fctx.Fctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return tpl.Execute(fctx.Fctx.Response().BodyWriter(), struct{}{})
}
