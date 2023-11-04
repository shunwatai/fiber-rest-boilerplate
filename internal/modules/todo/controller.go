package todo

import (
	"errors"
	"fmt"
	"golang-api-starter/internal/helper"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) Controller {
	return Controller{s}
}

var respCode = fiber.StatusInternalServerError

func (c *Controller) Get(ctx *fiber.Ctx) error {
	fmt.Printf("todo ctrl\n")
	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	paramsMap := reqCtx.Payload.GetQueryString()
	paramsMap["exactMatch"] = map[string]bool{
		"id": true,
	}
	results, pagination := c.service.Get(paramsMap)

	respCode = fiber.StatusOK
	return ctx.
		Status(respCode).
		JSON(map[string]interface{}{"data": results, "pagination": pagination})
}

func (c *Controller) GetById(ctx *fiber.Ctx) error {
	fmt.Printf("todo ctrl\n")
	id := ctx.Params("id")
	paramsMap := map[string]interface{}{"id": id}
	results, err := c.service.GetById(paramsMap)

	if err != nil {
		respCode = fiber.StatusNotFound
		return ctx.
			Status(respCode).
			JSON(map[string]interface{}{"msg": err.Error()})
	}
	respCode = fiber.StatusOK
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

	for _, todo := range todos {
		if todo.Id == nil {
			continue
		} else if existing, err := c.service.GetById(map[string]interface{}{
			"id": strconv.Itoa(int(*todo.Id)),
		}); err == nil && todo.CreatedAt == nil {
			todo.CreatedAt = existing[0].CreatedAt
		}
		fmt.Printf("todo? %+v\n", todo)
	}

	// return []*Todo{}
	results := c.service.Create(todos)

	respCode = fiber.StatusCreated
	if todoErr == nil && len(results) > 0 {
		return ctx.
			Status(respCode).
			JSON(map[string]interface{}{"data": results[0]})
	}
	return ctx.
		Status(respCode).
		JSON(map[string]interface{}{"data": results})
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

	for _, todo := range todos {
		existing, err := c.service.GetById(map[string]interface{}{"id": strconv.Itoa(int(*todo.Id))})
		if len(existing) == 0 {
			respCode = fiber.StatusNotFound
			return ctx.
				Status(respCode).
				JSON(map[string]interface{}{
					"message": errors.Join(
						errors.New("cannot update non-existing records..."),
						err,
					).Error(),
				})
		} else if todo.CreatedAt == nil {
			todo.CreatedAt = existing[0].CreatedAt
		}
	}

	results := c.service.Update(todos)

	respCode = fiber.StatusOK
	if todoErr == nil && len(results) > 0 {
		return ctx.
			Status(respCode).
			JSON(map[string]interface{}{"data": results[0]})
	}
	return ctx.
		Status(respCode).
		JSON(map[string]interface{}{"data": results})
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
		respCode = fiber.StatusNotFound
		return ctx.
			Status(respCode).
			JSON(map[string]interface{}{"message": err.Error()})
	}

	respCode = fiber.StatusOK
	return ctx.
		Status(respCode).
		JSON(map[string]interface{}{"data": results})
}
