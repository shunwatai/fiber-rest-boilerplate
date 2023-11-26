package todo

import (
	"fmt"
	"golang-api-starter/internal/helper"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Get(queries map[string]interface{}) ([]*Todo, *helper.Pagination) {
	fmt.Printf("todo service\n")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*Todo, error) {
	fmt.Printf("todo service\n")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(todos []*Todo) []*Todo {
	fmt.Printf("todo service create\n")
	return s.repo.Create(todos)
}

func (s *Service) Update(todos []*Todo) []*Todo {
	fmt.Printf("todo service update\n")
	return s.repo.Update(todos)
}

func (s *Service) Delete(ids []string) ([]*Todo, error) {
	var (
		records    = []*Todo{}
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
