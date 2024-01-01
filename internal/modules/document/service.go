package document

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

func (s *Service) Get(queries map[string]interface{}) ([]*Document, *helper.Pagination) {
	fmt.Printf("document service get\n")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*Document, error) {
	fmt.Printf("document service getById\n")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(documents []*Document) ([]*Document, *helper.HttpErr) {
	fmt.Printf("document service create\n")
  /*
	// use the claims for mark the "createdBy/updatedBy" in database
	claims := s.ctx.Locals("claims").(jwt.MapClaims)
	fmt.Println("req by:", claims["userId"], claims["username"])
	for _, document := range documents {
		if document.UserId == nil {
			document.UserId = claims["userId"]
		}
		if validErr := helper.ValidateStruct(*document); validErr != nil {
			return nil, &helper.HttpErr{fiber.StatusUnprocessableEntity, validErr}
		}
	}
  */

	results, err := s.repo.Create(documents)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Update(documents []*Document) ([]*Document, *helper.HttpErr) {
	fmt.Printf("document service update\n")
	results, err := s.repo.Update(documents)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Delete(ids []string) ([]*Document, error) {
	fmt.Printf("document service delete\n")
	var (
		records    = []*Document{}
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
