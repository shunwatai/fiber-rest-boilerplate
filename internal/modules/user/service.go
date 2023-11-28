package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-api-starter/internal/helper"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
	fmt.Printf("user service get\n")
	users, pagination := s.repo.Get(queries)

	return users, pagination
}

func (s *Service) GetById(queries map[string]interface{}) ([]*User, error) {
	fmt.Printf("user service getById\n")

	records, _ := s.repo.Get(queries)
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

func (s *Service) Update(users []*User) ([]*User, *helper.HttpErr) {
	fmt.Printf("user service update\n")

	userIds := []string{}
	for _, user := range users {
		userIds = append(userIds, strconv.Itoa(int(*user.Id)))
	}

	// create map by existing user from DB
	userIdMap := map[string]User{}
	existings, _ := s.repo.Get(map[string]interface{}{"id": userIds})
	for _, user := range existings {
		userIdMap[strconv.Itoa(int(*user.Id))] = *user
	}

	// check reqJson for non-existing ids
	// also reuse the map storing the req's user which use for create the "update data"
	nonExistIds := []int64{}
	for _, reqUser := range users {
		_, ok := userIdMap[strconv.Itoa(int(*reqUser.Id))]
		if !ok {
			nonExistIds = append(nonExistIds, *reqUser.Id)
		}
		userIdMap[strconv.Itoa(int(*reqUser.Id))] = *reqUser
	}

	if len(nonExistIds) > 0 || len(existings) == 0 {
		respCode = fiber.StatusNotFound
		notFoundMsg := fmt.Sprintf("cannot update non-existing id(s): %+v", nonExistIds)
		httpErr := &helper.HttpErr{
			Code: fiber.StatusNotFound,
			Err:  errors.New(notFoundMsg),
		}
		return nil, httpErr
	}

	// combining the req user that match with the existing user for update
	for _, originalUser := range existings {
		user := userIdMap[strconv.Itoa(int(*originalUser.Id))] // get the req user
		if user.CreatedAt == nil {
			user.CreatedAt = originalUser.CreatedAt
		}
		if user.Password == nil {
			user.Password = originalUser.Password
		} else {
			hashUserPassword(user.Password)
		}
		newUserBytes, _ := json.Marshal(user)       // convert req user into []byte
		json.Unmarshal(newUserBytes, &originalUser) // unmarshal the req user into its original db record
	}

	return s.repo.Update(existings), nil
}

func (s *Service) Delete(ids *[]int64) ([]*User, error) {
	idsString, _ := helper.ConvertNumberSliceToString(*ids)
	records, _ := s.repo.Get(map[string]interface{}{
		"id": idsString,
	})
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return s.repo.Delete(ids)
}
