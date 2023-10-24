package todo

import "fmt"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) Get(queries map[string]interface{}) []*Todo {
	fmt.Printf("todo service\n")
	return s.repo.Get(queries)
}

func (s *Service) Create(todos []*Todo) []*Todo {
	fmt.Printf("todo service create\n")
	return s.repo.Create(todos)
}

func (s *Service) Update(todos []*Todo) []*Todo {
	fmt.Printf("todo service update\n")
	return s.repo.Update(todos)
}

func (s *Service) Delete(ids *[]int64) ([]*Todo, error) {
	return s.repo.Delete(ids)
}
