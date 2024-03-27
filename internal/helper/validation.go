package helper

import (
	"errors"
	"golang-api-starter/internal/config"
	logger "golang-api-starter/internal/helper/logger/zap_log"

	"github.com/go-playground/validator/v10"
)

var cfg = config.Cfg

func ValidateStruct(strct interface{}) error {
	var invalidErrs []error
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.RegisterValidation("id_custom_validation", func(fl validator.FieldLevel) bool {
		// fmt.Printf("what is is? %+v, db: %+v\n", fl.Field().Interface(),cfg.DbConf.Driver)
		if cfg.DbConf.Driver == "mongodb" {
			_, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}
		} else {
			_, isFloat64 := fl.Field().Interface().(float64)
			_, isInt64 := fl.Field().Interface().(int64)
			_, isFlexInt := fl.Field().Interface().(FlexInt)

			if !isFloat64 && !isInt64 && !isFlexInt {
				return false
			}

		}
		return true
	})
	if err != nil {
		logger.Errorf("RegisterValidation err: %+v\n", err)
		return err
	}

	if err := validate.Struct(strct); err != nil {
		// logger.Errorf("validate err: %+v\n", err)
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			logger.Debugf("validate.Struct err: %+v\n", err)
			invalidErrs = append(invalidErrs, validationError)
		}

		return errors.Join(invalidErrs...)
	}
	return nil
}
