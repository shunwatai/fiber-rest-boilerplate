package document

import (
	"errors"
	"fmt"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/helper"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/maps"
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
	fmt.Printf("document ctrl\n")
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
	fmt.Printf("document ctrl\n")
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
	fmt.Printf("document ctrl create\n")
	c.service.ctx = ctx
	fctx := &helper.FiberCtx{Fctx: ctx}

	form, err := fctx.Fctx.MultipartForm()
	if err != nil { /* handle error */
		fmt.Printf("failed to get multipartForm, err: %+v\n", err.Error())
		return fctx.JsonResponse(
			respCode,
			map[string]interface{}{"message": err.Error()},
		)
	}

	results, httpErr := c.service.Create(form)
	if httpErr.Err != nil {
		fmt.Printf("document upload failed err: %+v\n", httpErr.Err)
		return fctx.JsonResponse(
			respCode,
			map[string]interface{}{"message": httpErr.Error()},
		)
	}

	respCode = fiber.StatusCreated
	if len(results) == 1 {
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
	fmt.Printf("document ctrl update\n")

	document := &Document{}
	documents := []*Document{}

	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	if invalidJson := reqCtx.Payload.ValidateJson(); invalidJson != nil {
		return fctx.JsonResponse(
			fiber.StatusUnprocessableEntity,
			map[string]interface{}{"message": invalidJson.Error()},
		)
	}

	documentErr, parseErr := reqCtx.Payload.ParseJsonToStruct(document, &documents)
	if parseErr != nil {
		return fctx.JsonResponse(
			fiber.StatusUnprocessableEntity,
			map[string]interface{}{"message": parseErr.Error()},
		)
	}
	if documentErr == nil {
		documents = append(documents, document)
	}

	for _, document := range documents {
		if validErr := helper.ValidateStruct(*document); validErr != nil {
			return fctx.JsonResponse(
				fiber.StatusUnprocessableEntity,
				map[string]interface{}{"message": validErr.Error()},
			)
		}
		if document.Id == nil && document.MongoId == nil {
			return fctx.JsonResponse(
				respCode,
				map[string]interface{}{"message": "please ensure all records with id for PATCH"},
			)
		}

		cfg.LoadEnvVariables()
		conditions := map[string]interface{}{}
		conditions["id"] = document.GetId()

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
		} else if document.CreatedAt == nil {
			document.CreatedAt = existing[0].CreatedAt
		}
	}

	results, httpErr := c.service.Update(documents)
	if httpErr.Err != nil {
		return fctx.JsonResponse(
			httpErr.Code,
			map[string]interface{}{"message": httpErr.Err.Error()},
		)
	}

	respCode = fiber.StatusOK
	if documentErr == nil && len(results) > 0 {
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
	fmt.Printf("document ctrl delete\n")
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
	fmt.Printf("deletedIds: %+v, mongoIds: %+v\n", delIds, mongoDelIds)

	var (
		results []*Document
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

func (c *Controller) GetDocument(ctx *fiber.Ctx) error {
	fmt.Printf("GetDocument ctrl\n")
	fctx := &helper.FiberCtx{Fctx: ctx}
	reqCtx := &helper.ReqContext{Payload: fctx}
	id := fctx.Fctx.Params("id")
	paramsMap := reqCtx.Payload.GetQueryString()
	maps.Copy(paramsMap, map[string]interface{}{"id": id})
	fileBuffer, fileType, fileName, err := c.service.GetDocument(paramsMap)

	if err != nil {
		respCode = fiber.StatusNotFound
		return fctx.JsonResponse(
			respCode,
			map[string]interface{}{"message": err.Error()},
		)
	}

	respCode = fiber.StatusOK
	fctx.Fctx.Response().Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	fctx.Fctx.Response().Header.Set("Content-Type", fileType)
	_, err = fctx.Fctx.Write(fileBuffer)
	return err
}
