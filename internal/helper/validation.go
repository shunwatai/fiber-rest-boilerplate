package helper

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang-api-starter/internal/config"
)

var cfg = config.Cfg

func ValidateStruct(strct interface{}) error {
	var invalidErrs []error
	validate := validator.New(validator.WithRequiredStructEnabled())

	// validates that an enum is within the interval
	err := validate.RegisterValidation("id_custom_validation", func(fl validator.FieldLevel) bool {
		// fmt.Printf("what is is? %+v\n", fl.Field().Interface())
		cfg.LoadEnvVariables()
		if cfg.DbConf.Driver == "mongodb" {
			_, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}
		} else {
			_, ok := fl.Field().Interface().(int64)
			if !ok {
				return false
			}
		}
		return true
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := validate.Struct(strct); err != nil {
		// fmt.Printf("validate err: %+v\n", err)
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Error())
			invalidErrs = append(invalidErrs, validationError)
		}

		return errors.Join(invalidErrs...)
	}
	return nil
}
