package server

import (
	"fmt"
	"golang-api-starter/internal/auth/oauth"
	"golang-api-starter/internal/config"
	zlog "golang-api-starter/internal/helper/logger/zap_log"
	"golang-api-starter/internal/middleware/logging"
	"golang-api-starter/internal/modules/document"
	"golang-api-starter/internal/modules/log"
	"golang-api-starter/internal/modules/qrcode"
	"golang-api-starter/internal/modules/todo"
	"golang-api-starter/internal/modules/todoDocument"
	"golang-api-starter/internal/modules/user"
	lg "log"
	"strings"

	_ "golang-api-starter/docs" // docs is generated by Swag CLI, you have to import it.

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger" // swagger handler
)

type Fiber struct {
	App *fiber.App
}

var cfg = config.Cfg

func (f *Fiber) GetApp() {
	cfg.LoadEnvVariables()
	zlog.NewZlog()
	f.App = fiber.New(fiber.Config{
		// Prefork:       true,
		CaseSensitive: true,
		StrictRouting: false,
		ServerHeader:  "Fiber",
		BodyLimit:     500 << 20, // 500Mb
	})
}

func (f *Fiber) LoadMiddlewares() {
	f.App.Use(logger.New())
	f.App.Use(recover.New())
	f.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	oauth.NewGoogleOAuth()
}

func (f *Fiber) LoadSwagger() {
	/* for swagger web */
	serverUrl := fmt.Sprintf("http://%s/swagger/doc.json", fmt.Sprintf("%s:%s", cfg.ServerConf.Host, cfg.ServerConf.Port))
	f.App.Get("/swagger/*", swagger.HandlerDefault)
	f.App.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         fmt.Sprintf("http://%s/doc.json", serverUrl),
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		// OAuth: &swagger.OAuthConfig{
		// 	AppName:  "OAuth Provider",
		// 	ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		// },
		// Ability to change OAuth2 redirect uri location
		// OAuth2RedirectUrl: fmt.Sprintf("http://%s:8080/swagger/oauth2-redirect.html", serverUrl),
	}))
}

func (f *Fiber) LoadAllRoutes() {
	api := f.App.Group("/api", logging.Logger())
	document.GetRoutes(api)
	log.GetRoutes(api)
	qrcode.GetRoutes(api)
	todo.GetRoutes(api)
	todoDocument.GetRoutes(api)
	user.GetRoutes(api)

	// a custom 404 handler instead of default "Cannot GET /page-not-found"
	// ref: https://github.com/gofiber/fiber/issues/748#issuecomment-687503079
	f.App.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "Resource Not Found",
		})
	})
}

func (f *Fiber) Start() {
	cfg.WatchConfig()

	fmt.Println(strings.Repeat("*", 50))
	fmt.Printf("server env: %+v\n", cfg.ServerConf.Env)
	fmt.Printf("using DB: %+v\n", cfg.DbConf.Driver)
	fmt.Println(strings.Repeat("*", 50))

	lg.Fatal(f.App.Listen(fmt.Sprintf(":%s", cfg.ServerConf.Port)))
}

var Api = &Fiber{}
