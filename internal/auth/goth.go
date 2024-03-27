package auth

import (
	"golang-api-starter/internal/config"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

var cfg = config.Cfg

func getSessionCfg() session.Config {
	return session.Config{
		Expiration:     720 * time.Hour, // 30 days
		KeyLookup:      "cookie:_gothic_session",
		CookieDomain:   "",
		CookiePath:     "/",
		CookieSecure:   cfg.ServerConf.Env == "prod",
		CookieHTTPOnly: true, // Should always be enabled
		CookieSameSite: "Lax",
		KeyGenerator:   utils.UUIDv4,
	}
}

func NewOAuth() {
	// create session handler
	var sessions = session.New(getSessionCfg())

	// load env
	googleClientId := cfg.OAuth.OAuthGoogle.Key
	googleClientSecret := cfg.OAuth.OAuthGoogle.Secret
	googleCallbackUrl := cfg.OAuth.OAuthGoogle.CallbackUrl
	// logger.Debugf("key: %+v, secret: %+v, callback: %+v", googleClientId, googleClientSecret, googleCallbackUrl)

	githubClientId := cfg.OAuth.OAuthGithub.Key
	githubClientSecret := cfg.OAuth.OAuthGithub.Secret
	githubCallbackUrl := cfg.OAuth.OAuthGithub.CallbackUrl
	// logger.Debugf("key: %+v, secret: %+v, callback: %+v", githubClientId, githubClientSecret, githubCallbackUrl)

	goth_fiber.SessionStore = sessions

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, googleCallbackUrl),
		github.New(githubClientId, githubClientSecret, githubCallbackUrl),
	)
}
