package todoDocument

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/helper"
	"log"
	"strconv"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{s}
}

var respCode = fiber.StatusInternalServerError

func (c *Controller) Get(ctx *fiber.Ctx) error {
	fmt.Printf("todoDocument ctrl\n")
	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	paramsMap := reqCtx.Payload.GetQueryString()
	results, pagination := c.service.Get(paramsMap)

	respCode = fiber.StatusOK
	return fctx.JsonResponse(
		respCode,
		map[string]interface{}{"data": results, "pagination": pagination},
	)
}

func (c *Controller) GetById(ctx *fiber.Ctx) error {
	fmt.Printf("todoDocument ctrl\n")
	fctx := &helper.FiberCtx{Fctx: ctx}
	id := fctx.Fctx.Params("id")
	paramsMap := map[string]interface{}{"id": id}
	results, err := c.service.GetById(paramsMap)

	if err != nil {
		respCode = fiber.StatusNotFound
		return fctx.JsonResponse(
			respCode,
			map[string]interface{}{"message": err.Error()},
		)
	}

	respCode = fiber.StatusOK
	return fctx.JsonResponse(respCode, map[string]interface{}{"data": results[0]})
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	fmt.Printf("todoDocument ctrl create\n")
	c.service.ctx = ctx
	todoDocument := &TodoDocument{}
	todoDocuments := []*TodoDocument{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	if invalidJson := reqCtx.Payload.ValidateJson(); invalidJson != nil {
		return fctx.JsonResponse(
			fiber.StatusUnprocessableEntity,
			map[string]interface{}{"message": invalidJson.Error()},
		)
	}

	todoDocumentErr, parseErr := reqCtx.Payload.ParseJsonToStruct(todoDocument, &todoDocuments)
	if parseErr != nil {
		return fctx.JsonResponse(
			fiber.StatusUnprocessableEntity,
			map[string]interface{}{"message": parseErr.Error()},
		)
	}
	if todoDocumentErr == nil {
		todoDocuments = append(todoDocuments, todoDocument)
	}
	// log.Printf("todoDocumentErr: %+v, todoDocumentsErr: %+v\n", todoDocumentErr, todoDocumentsErr)
	// for _, t := range todoDocuments {
	// 	log.Printf("todoDocuments: %+v\n", t)
	// }

	for _, todoDocument := range todoDocuments {
		if validErr := helper.ValidateStruct(*todoDocument); validErr != nil {
			return fctx.JsonResponse(
				fiber.StatusUnprocessableEntity,
				map[string]interface{}{"message": validErr.Error()},
			)
		}

		if todoDocument.Id == nil {
			continue
		} else if existing, err := c.service.GetById(map[string]interface{}{
			"id": strconv.Itoa(int(*todoDocument.Id)),
		}); err == nil && todoDocument.CreatedAt == nil {
			todoDocument.CreatedAt = existing[0].CreatedAt
		}
		// fmt.Printf("todoDocument? %+v\n", todoDocument)
	}

	// return []*TodoDocument{}
	results, httpErr := c.service.Create(todoDocuments)
	if httpErr.Err != nil {
		return fctx.JsonResponse(
			httpErr.Code,
			map[string]interface{}{"message": httpErr.Err.Error()},
		)
	}

	respCode = fiber.StatusCreated
	if todoDocumentErr == nil && len(results) > 0 {
		return fctx.JsonResponse(
			respCode,
			map[string]interface{}{"data": results[0]},
		)
	}
	return fctx.JsonResponse(
		respCode,
		map[string]interface{}{"data": results},
	)
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	fmt.Printf("todoDocument ctrl update\n")

	todoDocument := &TodoDocument{}
	todoDocuments := []*TodoDocument{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	if invalidJson := reqCtx.Payload.ValidateJson(); invalidJson != nil {
		return fctx.JsonResponse(
			fiber.StatusUnprocessableEntity,
			map[string]interface{}{"message": invalidJson.Error()},
		)
	}

	todoDocumentErr, parseErr := reqCtx.Payload.ParseJsonToStruct(todoDocument, &todoDocuments)
	if parseErr != nil {
		return fctx.JsonResponse(
			fiber.StatusUnprocessableEntity,
			map[string]interface{}{"message": parseErr.Error()},
		)
	}
	if todoDocumentErr == nil {
		todoDocuments = append(todoDocuments, todoDocument)
	}

	for _, todoDocument := range todoDocuments {
		if validErr := helper.ValidateStruct(*todoDocument); validErr != nil {
			return fctx.JsonResponse(
				fiber.StatusUnprocessableEntity,
				map[string]interface{}{"message": validErr.Error()},
			)
		}
		if todoDocument.Id == nil && todoDocument.MongoId == nil {
			return fctx.JsonResponse(
				respCode,
				map[string]interface{}{"message": "please ensure all records with id for PATCH"},
			)
		}

		conditions := map[string]interface{}{}
		conditions["id"] = todoDocument.GetId()

		existing, err := c.service.GetById(conditions)
		if len(existing) == 0 {
			respCode = fiber.StatusNotFound
			return fctx.JsonResponse(
				respCode,
				map[string]interface{}{
					"message": errors.Join(
						errors.New("cannot update non-existing records..."),
						err,
					).Error(),
				},
			)
		} else if todoDocument.CreatedAt == nil {
			todoDocument.CreatedAt = existing[0].CreatedAt
		}
	}

	results, httpErr := c.service.Update(todoDocuments)
	if httpErr.Err != nil {
		return fctx.JsonResponse(
			httpErr.Code,
			map[string]interface{}{"message": httpErr.Err.Error()},
		)
	}

	respCode = fiber.StatusOK
	if todoDocumentErr == nil && len(results) > 0 {
		return fctx.JsonResponse(
			respCode,
			map[string]interface{}{"data": results[0]},
		)
	}
	return fctx.JsonResponse(
		respCode,
		map[string]interface{}{"data": results},
	)
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	fmt.Printf("todoDocument ctrl delete\n")
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
		return fctx.JsonResponse(respCode, map[string]interface{}{"message": errors.Join(intIdsErr, strIdsErr).Error()})
	}
	if len(delIds.Ids) == 0 && len(mongoDelIds.Ids) == 0 {
		return fctx.JsonResponse(respCode, map[string]interface{}{"message": "please check the req json like the follow: {\"ids\":[]}"})
	}
	fmt.Printf("deletedIds: %+v, mongoIds: %+v\n", delIds, mongoDelIds)

	var (
		results []*TodoDocument
		err     error
	)

	if cfg.DbConf.Driver == "mongodb" {
		results, err = c.service.Delete(mongoDelIds.Ids)
	} else {
		idsString, _ := helper.ConvertNumberSliceToString(delIds.Ids)
		results, err = c.service.Delete(idsString)
	}

	if err != nil {
		log.Printf("failed to delete, err: %+v\n", err.Error())
		respCode = fiber.StatusNotFound
		return fctx.JsonResponse(
			respCode,
			map[string]interface{}{"message": err.Error()},
		)
	}

	respCode = fiber.StatusOK
	return fctx.JsonResponse(
		respCode,
		map[string]interface{}{"data": results},
	)
}
