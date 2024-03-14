package oauth

import (
	"golang-api-starter/internal/config"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

var cfg = config.Cfg

var sessionCfg = session.Config{
	Expiration:   720 * time.Hour, // 30 days
	KeyLookup:    "cookie:_gothic_session",
	CookieDomain: "",
	CookiePath:   "/",
	CookieSecure: false,
	// CookieSecure:   os.Getenv("ENVIRONMENT") == "production",
	CookieHTTPOnly: true, // Should always be enabled
	CookieSameSite: "Lax",
	KeyGenerator:   utils.UUIDv4,
}

// create session handler
var sessions = session.New(sessionCfg)

func NewGoogleOAuth() {
	// load env
	googleClientId := cfg.OAuth.OAuthGoogle.Key
	googleClientSecret := cfg.OAuth.OAuthGoogle.Secret
	googleCallbackUrl := cfg.OAuth.OAuthGoogle.CallbackUrl

	goth_fiber.SessionStore = sessions

	// logger.Debugf("googleClientId: %+v", googleClientId)
	// logger.Debugf("googleClientSecret: %+v", googleClientSecret)
	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, googleCallbackUrl),
	)
	// logger.Debugf("used google provider")
}
