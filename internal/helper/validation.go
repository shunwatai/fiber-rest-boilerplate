package helper

import (
	"golang-api-starter/internal/config"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var cfg = config.Cfg
var Validate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	Validate.RegisterValidation("id_custom_validation", func(fl validator.FieldLevel) bool {
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

	Validate.RegisterValidation("alphanumspace", func(fl validator.FieldLevel) bool {
		var regexAlphaNumSpace = regexp.MustCompile("^[ \\p{L}\\p{N}]+$")
		return regexAlphaNumSpace.MatchString(fl.Field().String())
	})
}
