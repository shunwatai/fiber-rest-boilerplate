package logging

import (
	"encoding/json"
	"golang-api-starter/internal/auth"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/helper"
	logger "golang-api-starter/internal/helper/log"
	"log"
	"slices"

	"github.com/gofiber/fiber/v2"

	customLog "golang-api-starter/internal/modules/log"
	"time"
)

var cfg = config.Cfg

/*
 * Logger() is a middleware for showing the http req & resp info
 */
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// log.Println("*********************")
		// log.Printf("I AM LOGGER....\n")
		// log.Println("*********************")

		bodyBytes := c.BodyRaw()
		// log.Printf("1reqBody: %+v, %+v \n", len(string(bodyBytes)), string(bodyBytes))
		var reqBodyJson, respBodyJson *string
		if len(string(bodyBytes)) > 0 {
			reqBodyJson = helper.ToPtr(string(bodyBytes))
		}

		reqHeader, _ := json.Marshal(c.GetReqHeaders())
		// log.Printf("reqHeader: %+v \n", string(reqHeader))

		var userId interface{}
		claims, err := auth.ParseJwt(c.Get("Authorization"))
		if err == nil {
			userId = claims["userId"]
		}
		// log.Println("JWT userId:", userId)
		// log.Println("created by:", claims["userId"], claims["username"])

		start := time.Now()
		defer func() {
			/* write to database or send to monitor service */
			ip := c.IP()
			// log.Println("from IP:", ip)
			if len(string(c.Response().Body())) > 0 {
				respBodyJson = helper.ToPtr(string(c.Response().Body()))
			}

			/* insert into logs table */
			if slices.Contains(cfg.Logging.Type, "database") {
				logData := []*customLog.Log{{
					UserId:        userId,
					IpAddress:     ip,
					HttpMethod:    c.Method(),
					Route:         c.Request().URI().String(),
					UserAgent:     string(c.Request().Header.UserAgent()),
					RequestHeader: string(reqHeader),
					RequestBody:   reqBodyJson,
					ResponseBody:  respBodyJson,
					Status:        int64(c.Response().StatusCode()),
					Duration:      time.Since(start).Milliseconds(),
					CreatedAt:     &helper.CustomDatetime{&start, helper.ToPtr(time.RFC3339)},
				}}
				log.Printf("%+v\n", logData)

				// create log to database,
				// WARN: this will slower the performance as one more database operation
				customLog.Srvc.Create(logData)
			}

			// create log to files
			if slices.Contains(cfg.Logging.Type, "zap") {
				zl := logger.ZapLog{
					Output: logger.OutputTypes{
						File:    slices.Contains(cfg.Logging.Zap.Output, "file"),
						Console: slices.Contains(cfg.Logging.Zap.Output, "console"),
					},
					Filename: &cfg.Logging.Zap.Filename,
				}
				if !zl.Output.Console && !zl.Output.File {
					return
				}
				zl.RequestLog("FIBER REQ LOG",
					"UserId", userId,
					"IpAddress", ip,
					"HttpMethod", c.Method(),
					"Route", c.Request().URI().String(),
					"UserAgent", string(c.Request().Header.UserAgent()),
					"RequestHeader", string(reqHeader),
					"RequestBody", reqBodyJson,
					"ResponseBody", respBodyJson,
					"Status", int64(c.Response().StatusCode()),
					"Duration", time.Since(start).Milliseconds(),
					"CreatedAt", &helper.CustomDatetime{&start, helper.ToPtr(time.RFC3339)},
				)
			}
		}()

		return c.Next()
	}
}
