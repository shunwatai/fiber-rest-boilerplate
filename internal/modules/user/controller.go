package user

import (
	"errors"
	"fmt"
	"golang-api-starter/internal/auth"
	"golang-api-starter/internal/helper"
	"log"
	"time"
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

/* helper func for Login & Refresh funcs below */
func SetRefreshTokenInCookie(result map[string]interface{}, c *fiber.Ctx) {
	cfg.LoadEnvVariables()
	env := cfg.ServerConf.Env
	refreshToken := result["refreshToken"].(string)
	cookie := &fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 720), // 30 days
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
	}
	if env == "local" {
		cookie.Secure = false
	}

	c.Cookie(cookie)
	delete(result, "refreshToken")
}

func (c *Controller) Get(ctx *fiber.Ctx) error {
	fmt.Printf("user ctrl\n")
	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	paramsMap := reqCtx.Payload.GetQueryString()
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
			"id": user.GetId(),
		}); err == nil && user.CreatedAt == nil {
			user.CreatedAt = existing[0].CreatedAt
		}
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

	for _, user := range users {
		if user.Id == nil && user.MongoId == nil {
			return ctx.
				Status(respCode).
				JSON(map[string]interface{}{"message": "please ensure all records with id for PATCH"})
		}

		cfg.LoadEnvVariables()
		conditions := map[string]interface{}{}
		conditions["id"] = user.GetId()

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
		} else if user.CreatedAt == nil {
			user.CreatedAt = existing[0].CreatedAt
		}
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
		results []*User
		err     error
	)

	cfg.LoadEnvVariables()
	if cfg.DbConf.Driver == "mongodb" {
		results, err = c.service.Delete(mongoDelIds.Ids)
	} else {
		idsString, _ := helper.ConvertNumberSliceToString(delIds.Ids)
		results, err = c.service.Delete(idsString)
	}
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

func (c *Controller) Login(ctx *fiber.Ctx) error {
	fmt.Printf("user ctrl create\n")
	user := &User{}
	users := []*User{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	if userErr, _ := reqCtx.Payload.ParseJsonToStruct(user, &users); userErr != nil {
		log.Printf("userErr: %+v\n", userErr)
	}
	// log.Printf("login req: %+v\n", user)

	result, httpErr := c.service.Login(user)
	if httpErr != nil {
		return ctx.
			Status(httpErr.Code).
			JSON(map[string]interface{}{"message": httpErr.Err.Error()})
	}

	SetRefreshTokenInCookie(result, ctx)
	respCode = fiber.StatusOK
	return ctx.
		Status(respCode).
		JSON(map[string]interface{}{"data": result})
}

func (c *Controller) Refresh(ctx *fiber.Ctx) error {
	// Read cookie
	cookie := ctx.Cookies("refreshToken")

	refreshToken := "Bearer " + cookie
	fmt.Printf("%s\n", refreshToken)

	claims, err := auth.ParseJwt(refreshToken)
	if claims["tokenType"] != "refreshToken" || err != nil {
		respCode = fiber.StatusExpectationFailed
		return ctx.
			Status(respCode).
			JSON(map[string]interface{}{"message": "Invalid Token type... please try to login again"})
	}

	result := map[string]interface{}{}
	cfg.LoadEnvVariables()
	if cfg.DbConf.Driver == "mongodb" {
		userId := claims["userId"].(string)
		result, _ = c.service.Refresh(&User{MongoId: &userId})
	} else {
		userId := int64(claims["userId"].(float64))
		result, _ = c.service.Refresh(&User{Id: &userId})
	}
	SetRefreshTokenInCookie(result, ctx)

	respCode = fiber.StatusOK
	return ctx.
		Status(respCode).
		JSON(map[string]interface{}{"data": result})
}
