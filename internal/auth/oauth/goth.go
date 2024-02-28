package oauth

import (
	"golang-api-starter/internal/config"
	logger "golang-api-starter/internal/helper/logger/zap_log"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var cfg = config.Cfg

func NewGoogleOAuth() {
	// load env
	googleClientId := cfg.OAuth.OAuthGoogle.Key
	googleClientSecret := cfg.OAuth.OAuthGoogle.Secret

	// gorilla sessions
	key := ""            // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	logger.Debugf("googleClientId: %+v", googleClientId)
	logger.Debugf("googleClientSecret: %+v", googleClientSecret)
	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:7000/auth/google/callback"),
	)
	logger.Debugf("used google provider")
}
