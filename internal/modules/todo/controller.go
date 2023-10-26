package todo

import (
	"fmt"
	"golang-api-starter/internal/helper"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) Controller {
	return Controller{s}
}

func (c *Controller) Get(ctx *fiber.Ctx) error {
	fmt.Printf("todo ctrl\n")
	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	paramsMap := reqCtx.Payload.GetQueryString()
	results := c.service.Get(paramsMap)

	return ctx.JSON(map[string]interface{}{"data": results})
}

func (c *Controller) GetById(ctx *fiber.Ctx) error {
	fmt.Printf("todo ctrl\n")
	id := ctx.Params("id")
	paramsMap := map[string]interface{}{"id": id}
	results := c.service.Get(paramsMap)

	if len(results) == 0 {
		return ctx.JSON(map[string]interface{}{"msg": fmt.Sprintf("record with id: %s not found", id)})
	}
	return ctx.JSON(map[string]interface{}{"data": results[0]})
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	fmt.Printf("todo ctrl create\n")
	todo := &Todo{}
	todos := []*Todo{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	todoErr, _ := reqCtx.Payload.ParseJsonToStruct(todo, &todos)
	if todoErr == nil {
		todos = append(todos, todo)
	}
	// log.Printf("todoErr: %+v, todosErr: %+v\n", todoErr, todosErr)
	// for _, t := range todos {
	// 	log.Printf("todos: %+v\n", t)
	// }

	t := time.Now()
	for _, todo := range todos {
		// t := time.Now().Format("2006-01-02 15:04:05")
		if todo.CreatedAt == nil {
			todo.CreatedAt = &t
		}
		if todo.UpdatedAt == nil {
			todo.UpdatedAt = &t
		}
	}
	// return []*Todo{}
	results := c.service.Create(todos)

	if todoErr == nil && len(results) > 0 {
		return ctx.JSON(map[string]interface{}{"data": results[0]})
	}
	return ctx.JSON(map[string]interface{}{"data": results})
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	fmt.Printf("todo ctrl update\n")

	todo := &Todo{}
	todos := []*Todo{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	todoErr, _ := reqCtx.Payload.ParseJsonToStruct(todo, &todos)
	if todoErr == nil {
		todos = append(todos, todo)
	}
	// log.Printf("todoErr: %+v, todosErr: %+v\n", todoErr, todosErr)
	// for _, t := range todos {
	// 	log.Printf("todos: %+v\n", t)
	// }

	t := time.Now()
	for _, todo := range todos {
		if todo.Id == nil {
			todo.CreatedAt = &t
		}
		todo.UpdatedAt = &t
	}

	results := c.service.Update(todos)

	if todoErr == nil && len(results) > 0 {
		return ctx.JSON(map[string]interface{}{"data": results[0]})
	}
	return ctx.JSON(map[string]interface{}{"data": results})
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	// body := map[string]interface{}{}
	// json.Unmarshal(c.BodyRaw(), &body)
	// fmt.Printf("req body: %+v\n", body)
	delIds := struct {
		Ids []int64 `json:"ids" validate:"required,min=1,unique"`
	}{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	err, _ := reqCtx.Payload.ParseJsonToStruct(&delIds, nil)
	if err != nil {
		log.Printf("failed to parse req json, %+v\n", err.Error())
		return ctx.JSON(map[string]interface{}{"message": err.Error()})
	}

	fmt.Printf("deletedIds: %+v\n", delIds)

	results, err := c.service.Delete(&delIds.Ids)
	if err != nil {
		log.Printf("failed to delete, err: %+v\n", err.Error())
		return ctx.JSON(map[string]interface{}{"message": err.Error()})
	}

	return ctx.JSON(map[string]interface{}{"data": results})
}
