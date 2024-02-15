package logging

import (
	"encoding/json"
	"golang-api-starter/internal/auth"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/helper"
	"golang-api-starter/internal/helper/logger"
	customLog "golang-api-starter/internal/modules/log"
	"log"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
)

var cfg = config.Cfg
var zlog = logger.Zlog

/*
 * Logger() is a middleware for showing the http req & resp info
 */
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// zlog.Printf("I AM LOGGER....")

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
				if !zlog.Output.Console && !zlog.Output.File {
					return
				}
				zlog.RequestLog("FIBER REQ LOG",
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
