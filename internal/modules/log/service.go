package log

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

func (s *Service) Get(queries map[string]interface{}) ([]*Log, *helper.Pagination) {
	fmt.Printf("log service get\n")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*Log, error) {
	fmt.Printf("log service getById\n")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(logs []*Log) ([]*Log, *helper.HttpErr) {
	fmt.Printf("log service create\n")
  /*
	// use the claims for mark the "createdBy/updatedBy" in database
	claims := s.ctx.Locals("claims").(jwt.MapClaims)
	fmt.Println("req by:", claims["userId"], claims["username"])
	for _, log := range logs {
		if log.UserId == nil {
			log.UserId = claims["userId"]
		}
		if validErr := helper.ValidateStruct(*log); validErr != nil {
			return nil, &helper.HttpErr{fiber.StatusUnprocessableEntity, validErr}
		}
	}
  */

	results, err := s.repo.Create(logs)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Update(logs []*Log) ([]*Log, *helper.HttpErr) {
	fmt.Printf("log service update\n")
	results, err := s.repo.Update(logs)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Delete(ids []string) ([]*Log, error) {
	fmt.Printf("log service delete\n")

	getByIdsCondition := helper.GetIdsMapCondition(nil, ids)
	records, _ := s.repo.Get(getByIdsCondition)
	fmt.Printf("records: %+v\n", records)
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return records, s.repo.Delete(ids)
}
