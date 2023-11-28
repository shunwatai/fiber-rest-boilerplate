package user

import (
	"fmt"
	"golang-api-starter/internal/helper"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func sanitise(users Users) {
	for _, u := range users {
		u.Password = nil
	}
}

func NewController(s *Service) Controller {
	return Controller{s}
}

var respCode = fiber.StatusInternalServerError

func (c *Controller) Get(ctx *fiber.Ctx) error {
	fmt.Printf("user ctrl\n")
	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	paramsMap := reqCtx.Payload.GetQueryString()
	paramsMap["exactMatch"] = map[string]bool{
		"id": true,
	}
	results, pagination := c.service.Get(paramsMap)
	sanitise(results)

	respCode = fiber.StatusOK
	return ctx.
		Status(respCode).
		JSON(map[string]interface{}{"data": results, "pagination": pagination})
}

func (c *Controller) GetById(ctx *fiber.Ctx) error {
	fmt.Printf("user ctrl\n")
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
	fmt.Printf("user ctrl create\n")
	user := &User{}
	users := []*User{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	userErr, _ := reqCtx.Payload.ParseJsonToStruct(user, &users)
	if userErr == nil {
		users = append(users, user)
	}
	// log.Printf("userErr: %+v, usersErr: %+v\n", userErr, usersErr)
	// for _, t := range users {
	// 	log.Printf("users: %+v\n", t)
	// }

	for _, user := range users {
		if user.Id == nil {
			continue
		} else if existing, err := c.service.GetById(map[string]interface{}{
			"id": strconv.Itoa(int(*user.Id)),
		}); err == nil && user.CreatedAt == nil {
			user.CreatedAt = existing[0].CreatedAt
		}
		fmt.Printf("user? %+v\n", user)
	}

	results, httpErr := c.service.Create(users)
	sanitise(results)
	if httpErr.Err != nil {
		return ctx.
			Status(httpErr.Code).
			JSON(map[string]interface{}{"message": httpErr.Err.Error()})
	}

	respCode = fiber.StatusCreated
	if userErr == nil && len(results) > 0 {
		return ctx.
			Status(respCode).
			JSON(map[string]interface{}{"data": results[0]})
	}
	return ctx.
		Status(respCode).
		JSON(map[string]interface{}{"data": results})
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	fmt.Printf("user ctrl update\n")

	user := &User{}
	users := []*User{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	userErr, _ := reqCtx.Payload.ParseJsonToStruct(user, &users)
	if userErr == nil {
		users = append(users, user)
	}
	// log.Printf("userErr: %+v, usersErr: %+v\n", userErr, usersErr)
	// for _, t := range users {
	// 	log.Printf("users: %+v\n", t)
	// }

	userIds := []string{}
	for _, user := range users {
		if user.Id == nil {
			return ctx.
				Status(respCode).
				JSON(map[string]interface{}{"message": "please ensure all records with key 'id' for PATCH"})
		}
		userIds = append(userIds, strconv.Itoa(int(*user.Id)))
	}

	results, httpErr := c.service.Update(users)
	if httpErr.Err != nil {
		return ctx.
			Status(httpErr.Code).
			JSON(map[string]interface{}{"message": httpErr.Err.Error()})
	}
	sanitise(results)

	respCode = fiber.StatusOK
	if userErr == nil && len(results) > 0 {
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
	sanitise(results)
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
