package todo

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/helper"
	"log"
	"strconv"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) Controller {
	return Controller{s}
}

var cfg = config.Cfg
var respCode = fiber.StatusInternalServerError

func (c *Controller) Get(ctx *fiber.Ctx) error {
	fmt.Printf("todo ctrl\n")
	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	paramsMap := reqCtx.Payload.GetQueryString()
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
			JSON(map[string]interface{}{"message": err.Error()})
	}
	respCode = fiber.StatusOK
	return ctx.JSON(map[string]interface{}{"data": results[0]})
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	fmt.Printf("todo ctrl create\n")
	c.service.ctx = ctx
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
	results, httpErr := c.service.Create(todos)
	if httpErr.Err != nil {
		return ctx.
			Status(httpErr.Code).
			JSON(map[string]interface{}{"message": httpErr.Err.Error()})
	}

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
		if todo.Id == nil && todo.MongoId == nil {
			return ctx.
				Status(respCode).
				JSON(map[string]interface{}{"message": "please ensure all records with id for PATCH"})
		}

		cfg.LoadEnvVariables()
		conditions := map[string]interface{}{}
		conditions["id"] = todo.GetId()

		existing, err := c.service.GetById(conditions)
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

	results, httpErr := c.service.Update(todos)
	if httpErr.Err != nil {
		return ctx.
			Status(httpErr.Code).
			JSON(map[string]interface{}{"message": httpErr.Err.Error()})
	}

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
	fmt.Printf("todo ctrl delete\n")
	// body := map[string]interface{}{}
	// json.Unmarshal(c.BodyRaw(), &body)
	// fmt.Printf("req body: %+v\n", body)
	delIds := struct {
		Ids []int64 `json:"ids" validate:"required,unique"`
	}{}

	mongoDelIds := struct {
		Ids []string `json:"ids" validate:"required,unique"`
	}{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	intIdsErr, strIdsErr := reqCtx.Payload.ParseJsonToStruct(&delIds, &mongoDelIds)
	if intIdsErr != nil && strIdsErr != nil {
		log.Printf("failed to parse req json, %+v\n", errors.Join(intIdsErr, strIdsErr).Error())
		return ctx.JSON(map[string]interface{}{"message": errors.Join(intIdsErr, strIdsErr).Error()})
	}
	fmt.Printf("deletedIds: %+v, mongoIds: %+v\n", delIds, mongoDelIds)

	var (
		results []*Todo
		err     error
	)

	cfg.LoadEnvVariables()
	if cfg.DbConf.Driver == "mongodb" {
		results, err = c.service.Delete(mongoDelIds.Ids)
	} else {
		idsString, _ := helper.ConvertNumberSliceToString(delIds.Ids)
		results, err = c.service.Delete(idsString)
	}

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
