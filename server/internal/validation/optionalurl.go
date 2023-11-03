package validation

import (
	"github.com/go-playground/validator/v10"
)

func validateOptionalURI(fl validator.FieldLevel) bool {
	uri := fl.Field().String()
	return IsURL(uri)
}
