package user

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

func (s *Service) Get(queries map[string]interface{}) ([]*User, *helper.Pagination) {
	fmt.Printf("user service\n")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*User, error) {
	fmt.Printf("user service\n")

	records,_ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(users []*User) []*User {
	fmt.Printf("user service create\n")
	return s.repo.Create(users)
}

func (s *Service) Update(users []*User) []*User {
	fmt.Printf("user service update\n")
	return s.repo.Update(users)
}

func (s *Service) Delete(ids *[]int64) ([]*User, error) {
	idsString, _ := helper.ConvertNumberSliceToString(*ids)
	records,_ := s.repo.Get(map[string]interface{}{
		"id": idsString,
	})
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return s.repo.Delete(ids)
}
