package todo

import (
	"fmt"
	"golang-api-starter/internal/helper"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Get(queries map[string]interface{}) ([]*Todo, *helper.Pagination) {
	fmt.Printf("todo service get\n")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*Todo, error) {
	fmt.Printf("todo service getById\n")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(todos []*Todo) ([]*Todo, *helper.HttpErr) {
	fmt.Printf("todo service create\n")
	results, err := s.repo.Create(todos)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Update(todos []*Todo) ([]*Todo, *helper.HttpErr) {
	fmt.Printf("todo service update\n")
	results, err := s.repo.Update(todos)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Delete(ids []string) ([]*Todo, error) {
	records,_ := s.repo.Get(map[string]interface{}{
		"id": ids,
	})
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return s.repo.Delete(ids)
}
