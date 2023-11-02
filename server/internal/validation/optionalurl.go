package validation

import (
	"net/url"

	"github.com/go-playground/validator/v10"
)

func validateOptionalURI(fl validator.FieldLevel) bool {
	uri := fl.Field().String()

	if uri == "" {
		return true
	}

	_, err := url.ParseRequestURI(uri)
	return err == nil
}
