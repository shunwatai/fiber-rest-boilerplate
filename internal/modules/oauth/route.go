package oauth

import (
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/config"
)

var (
	cfg       = config.Cfg
	tableName = "users"
	Srvc      = &Service{}
	ctrl      = &Controller{}
)

func GetRoutes(router fiber.Router) {
	Srvc = NewService()
	ctrl = NewController(Srvc)

	// oauth for google
	oauthRoute := router.Group("/oauth")
	oauthRoute.Get("/:provider/callback", OAuthGetAuth)
	oauthRoute.Get("/:provider/login", OAuthLogin)
	oauthRoute.Get("/logout", OAuthLogout)
	oauthRoute.Get("/sign-in", OAuthProviderPage)

}

func OAuthGetAuth(c *fiber.Ctx) error {
	return ctrl.OAuthGetAuth(c)
}
func OAuthLogin(c *fiber.Ctx) error {
	return ctrl.OAuthLogin(c)
}
func OAuthLogout(c *fiber.Ctx) error {
	return ctrl.OAuthLogout(c)
}
func OAuthProviderPage(c *fiber.Ctx) error {
	return ctrl.OAuthProviderPage(c)
}
