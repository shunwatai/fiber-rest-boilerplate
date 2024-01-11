package todoDocument

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/helper"
)

type Service struct {
	repo *Repository
	ctx  *fiber.Ctx
}

func NewService(r *Repository) *Service {
	return &Service{r, nil}
}

// func (s *Service) GetIdMap(tds []*TodoDocument) map[string][]*TodoDocument {
// 	todoDocumentsMap := map[string][]*TodoDocument{}
// 	for _, td := range tds {
// 		todoDocumentsMap[td.GetId()] = append(todoDocumentsMap[td.GetId()], td)
// 	}
// 	return todoDocumentsMap
// }

func (s *Service) GetTodoIdMap(tds []*TodoDocument) map[string][]*TodoDocument {
	todoDocumentsMap := map[string][]*TodoDocument{}
	for _, td := range tds {
		todoDocumentsMap[td.GetTodoId()] = append(todoDocumentsMap[td.GetTodoId()], td)
	}
	return todoDocumentsMap
}

func (s *Service) Get(queries map[string]interface{}) ([]*TodoDocument, *helper.Pagination) {
	fmt.Printf("todoDocument service get\n")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*TodoDocument, error) {
	fmt.Printf("todoDocument service getById\n")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(todoDocuments []*TodoDocument) ([]*TodoDocument, *helper.HttpErr) {
	fmt.Printf("todoDocument service create\n")
  /*
	// use the claims for mark the "createdBy/updatedBy" in database
	claims := s.ctx.Locals("claims").(jwt.MapClaims)
	fmt.Println("req by:", claims["userId"], claims["username"])
	for _, todoDocument := range todoDocuments {
		if todoDocument.UserId == nil {
			todoDocument.UserId = claims["userId"]
		}
		if validErr := helper.ValidateStruct(*todoDocument); validErr != nil {
			return nil, &helper.HttpErr{fiber.StatusUnprocessableEntity, validErr}
		}
	}
  */

	results, err := s.repo.Create(todoDocuments)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Update(todoDocuments []*TodoDocument) ([]*TodoDocument, *helper.HttpErr) {
	fmt.Printf("todoDocument service update\n")
	results, err := s.repo.Update(todoDocuments)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Delete(ids []string) ([]*TodoDocument, error) {
	fmt.Printf("todoDocument service delete\n")
	var (
		records    = []*TodoDocument{}
		conditions = map[string]interface{}{}
	)

	cfg.LoadEnvVariables()
	if cfg.DbConf.Driver == "mongodb" {
		conditions["_id"] = ids
	} else {
		conditions["id"] = ids
	}

	records, _ = s.repo.Get(conditions)
	fmt.Printf("records: %+v\n", records)
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return records, s.repo.Delete(ids)
}
