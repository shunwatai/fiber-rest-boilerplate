package todo

import (
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	logger "golang-api-starter/internal/helper/logger/zap_log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	repo *Repository
	ctx  *fiber.Ctx
}

func NewService(r *Repository) *Service {
	return &Service{r, nil}
}

func (s *Service) Get(queries map[string]interface{}) ([]*Todo, *helper.Pagination) {
	logger.Debugf("todo service get")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*Todo, error) {
	logger.Debugf("todo service getById")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(todos []*Todo) ([]*Todo, *helper.HttpErr) {
	logger.Debugf("todo service create")

	// use the claims for mark the "createdBy/updatedBy" in database
	claims := s.ctx.Locals("claims").(jwt.MapClaims)
	logger.Debugf("req by userId: %+v, username: %+v", claims["userId"], claims["username"])
	for _, todo := range todos {
		if todo.UserId == nil {
			todo.UserId = claims["userId"]
		}
		if validErr := helper.ValidateStruct(*todo); validErr != nil {
			return nil, &helper.HttpErr{fiber.StatusUnprocessableEntity, validErr}
		}
	}

	results, err := s.repo.Create(todos)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Update(todos []*Todo) ([]*Todo, *helper.HttpErr) {
	logger.Debugf("todo service update")
	results, err := s.repo.Update(todos)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Delete(ids []string) ([]*Todo, error) {
	logger.Debugf("todo service delete")

	getByIdsCondition := database.GetIdsMapCondition(nil, ids)
	records, _ := s.repo.Get(getByIdsCondition)
	logger.Debugf("records: %+v\n", records)
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return records, s.repo.Delete(ids)
}
