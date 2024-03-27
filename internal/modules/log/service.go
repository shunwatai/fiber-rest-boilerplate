package log

import (
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	logger "golang-api-starter/internal/helper/logger/zap_log"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repo *Repository
	ctx  *fiber.Ctx
}

func NewService(r *Repository) *Service {
	return &Service{r, nil}
}

func (s *Service) Get(queries map[string]interface{}) ([]*Log, *helper.Pagination) {
	logger.Debugf("log service get\n")
	return s.repo.Get(queries)
}

func (s *Service) GetById(queries map[string]interface{}) ([]*Log, error) {
	logger.Debugf("log service getById\n")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(logs []*Log) ([]*Log, *helper.HttpErr) {
	logger.Debugf("log service create\n")
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
	logger.Debugf("log service update")
	results, err := s.repo.Update(logs)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Delete(ids []string) ([]*Log, error) {
	logger.Debugf("log service delete")

	getByIdsCondition := database.GetIdsMapCondition(nil, ids)
	records, _ := s.repo.Get(getByIdsCondition)
	logger.Debugf("records: %+v", records)
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return records, s.repo.Delete(ids)
}
