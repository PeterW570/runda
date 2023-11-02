package validation

import "github.com/go-playground/validator/v10"

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("optional_uri", validateOptionalURI)

	return v
}
