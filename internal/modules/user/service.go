package user

import (
	"fmt"
	"golang-api-starter/internal/helper"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r}
}

func hashUserPassword(pwd *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*pwd), bcrypt.MinCost)
	if err != nil {
		return err
	}

	*pwd = string(hash)
	return nil
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

func (s *Service) Create(users []*User) ([]*User, error) {
	fmt.Printf("user service create\n")
	newUserNames := []string{}
	for _, user := range users {
		newUserNames = append(newUserNames, user.Name)
		// hash plain password
		if err := hashUserPassword(user.Password); err != nil {
			return nil, fmt.Errorf(err.Error())
		}
	}

	// check if duplicated by "name"
	existingUsers, _ := s.repo.Get(map[string]interface{}{"name": newUserNames})
	if len(existingUsers) > 0 {
		errMsg := fmt.Sprintf("user service create error: provided user name(s) %+v already exists.\n", newUserNames)
		log.Printf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	return s.repo.Create(users), nil
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
