package todo

import (
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	"log"

	"github.com/gofiber/fiber/v2"
)

// var db = &database.Postgre{
// 	Host:      "localhost",
// 	Port:      "3306",
// 	User:      "user",
// 	Pass:      "maria",
// 	TableName: "todo",
// }
var db = &database.Sqlite{
	ConnectionInfo: &database.ConnectionInfo{
		Driver:   "sqlite",
		Host:     "localhost",
		Port:     "",
		User:     "user",
		Pass:     "user",
		Database: "fiber-starter",
	},
	TableName: tableName,
}

var tableName = "todos"
var repo = NewRepository(db)
var srvc = NewService(repo)
var ctrl = NewController(srvc)

func GetRoutes(router fiber.Router) {
	r := router.Group("/todo")

	r.Get("/", func(c *fiber.Ctx) error {
		fctx := &helper.FiberCtx{Fctx: c}
		reqCtx := &helper.ReqContext{Payload: fctx}
		paramsMap := reqCtx.Payload.GetQueryString()
		results := ctrl.Get(paramsMap)

		return c.JSON(map[string]interface{}{"data": results})
	})

	r.Post("/", func(c *fiber.Ctx) error {
		todo := &Todo{}
		todos := []*Todo{}

		fctx := &helper.FiberCtx{Fctx: c}
		reqCtx := &helper.ReqContext{Payload: fctx}
		todoErr, _ := reqCtx.Payload.ParseJsonToStruct(todo, &todos)
		if todoErr == nil {
			todos = append(todos, todo)
		}
		// log.Printf("todoErr: %+v, todosErr: %+v\n", todoErr, todosErr)
		// for _, t := range todos {
		// 	log.Printf("todos: %+v\n", t)
		// }

		results := ctrl.Create(todos)

		if todoErr == nil {
			return c.JSON(map[string]interface{}{"data": results[0]})
		}
		return c.JSON(map[string]interface{}{"data": results})
	})

	r.Patch("/", func(c *fiber.Ctx) error {
		todo := &Todo{}
		todos := []*Todo{}

		fctx := &helper.FiberCtx{Fctx: c}
		reqCtx := &helper.ReqContext{Payload: fctx}
		todoErr, _ := reqCtx.Payload.ParseJsonToStruct(todo, &todos)
		if todoErr == nil {
			todos = append(todos, todo)
		}
		// log.Printf("todoErr: %+v, todosErr: %+v\n", todoErr, todosErr)
		// for _, t := range todos {
		// 	log.Printf("todos: %+v\n", t)
		// }

		results := ctrl.Update(todos)

		if todoErr == nil {
			return c.JSON(map[string]interface{}{"data": results[0]})
		}
		return c.JSON(map[string]interface{}{"data": results})
	})

	r.Delete("/", func(c *fiber.Ctx) error {
		// body := map[string]interface{}{}
		// json.Unmarshal(c.BodyRaw(), &body)
		// fmt.Printf("req body: %+v\n", body)
		delIds := struct {
			Ids []int64 `json:"ids" validate:"required,min=1,unique"`
		}{}

		fctx := &helper.FiberCtx{Fctx: c}
		reqCtx := &helper.ReqContext{Payload: fctx}
		err, _ := reqCtx.Payload.ParseJsonToStruct(&delIds, nil)
		if err != nil {
			log.Printf("failed to parse req json, %+v\n", err.Error())
			return c.JSON(map[string]interface{}{"message": err.Error()})
		}

		fmt.Printf("deletedIds: %+v\n", delIds)

		results, err := ctrl.Delete(&delIds.Ids)
		if err != nil {
			log.Printf("failed to delete, err: %+v\n", err.Error())
			return c.JSON(map[string]interface{}{"message": err.Error()})
		}

		return c.JSON(map[string]interface{}{"data": results})
	})
}
