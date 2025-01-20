package helper

import (
	"golang-api-starter/internal/config"
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var cfg = config.Cfg
var Validate = validator.New(validator.WithRequiredStructEnabled())

func isStrongPassword(userInput string) bool {
	var (
		hasMinLen   = false
		hasUpper    = false
		hasLower    = false
		hasNumber   = false
		hasSpecial  = false
		secureLevel = 0
	)
	if len(userInput) >= 4 {
		hasMinLen = true
	}
	for _, char := range userInput {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	for _, check := range []bool{hasUpper, hasLower, hasNumber, hasSpecial} {
		if check {
			secureLevel++
		}
	}

	return hasMinLen && secureLevel >= 3
	// return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

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

	// custom validator allow space in alphanum
	Validate.RegisterValidation("alphanumspace", func(fl validator.FieldLevel) bool {
		var regexAlphaNumSpace = regexp.MustCompile("^[ \\p{L}\\p{N}]+$")
		return regexAlphaNumSpace.MatchString(fl.Field().String())
	})

	// custom validator for requiring password len=4 with BIGsmallSymbol
	Validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return isStrongPassword(fl.Field().String())
	})
}
