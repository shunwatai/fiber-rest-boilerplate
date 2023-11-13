package helper

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/strcase"
)

type IReqPayload interface {
	GetQueryString() map[string]interface{}
	ParseJsonToStruct(interface{}, interface{}) (error, error)
}

type ReqContext struct {
	Payload IReqPayload
}

type FiberCtx struct {
	Fctx *fiber.Ctx
}

func (c *FiberCtx) GetQueryString() map[string]interface{} {
	queries := c.Fctx.Queries()

	params, err := url.ParseQuery(string(c.Fctx.Request().URI().QueryString()))
	if err != nil {
		log.Printf("ParseQuery err: %+v\n", err.Error())
	}
	fmt.Printf("queries: %+v\n", queries)

	var paramsMap = make(map[string]interface{}, 0)

	for key, value := range params {
		// fmt.Printf("-->  %v = %v\n", key, value)
		snakeCase := strcase.ToSnake(key)
		if strings.Contains(snakeCase, "date") || strings.Contains(snakeCase, "_at") {
			paramsMap["withDateFilter"] = true
		}

		if len(value) == 1 {
			paramsMap[snakeCase] = value[0]
			continue
		}
		paramsMap[snakeCase] = value
	}

	// fmt.Printf("paramsMap: %+v\n", paramsMap)
	return paramsMap
}

func (c *FiberCtx) ParseJsonToStruct(single interface{}, plural interface{}) (error, error) {
	singleErr := c.Fctx.BodyParser(single)
	pluralErr := c.Fctx.BodyParser(plural)

	if pluralErr != nil {
		log.Printf("pluralErr err: %+v\n", pluralErr.Error())
	}

	if singleErr != nil {
		log.Printf("singleErr err: %+v\n", singleErr.Error())
	}

	return singleErr, pluralErr
}
