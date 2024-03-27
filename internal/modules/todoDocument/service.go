package todoDocument

import (
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repo *Repository
	ctx  *fiber.Ctx
}

func NewService(r *Repository) *Service {
	return &Service{r, nil}
}

var mu sync.Mutex

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
	mu.Lock() // for avoid sqlite goroute race error
	defer mu.Unlock()
	logger.Debugf("todoDocument service get")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*TodoDocument, error) {
	logger.Debugf("todoDocument service getById")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(todoDocuments []*TodoDocument) ([]*TodoDocument, *helper.HttpErr) {
	logger.Debugf("todoDocument service create")
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
	logger.Debugf("todoDocument service update")
	results, err := s.repo.Update(todoDocuments)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Delete(ids []string) ([]*TodoDocument, error) {
	logger.Debugf("todoDocument service delete")

	getByIdsCondition := database.GetIdsMapCondition(nil, ids)
	records, _ := s.repo.Get(getByIdsCondition)
	logger.Debugf("records: %+v", records)
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return records, s.repo.Delete(ids)
}
