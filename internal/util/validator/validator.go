package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shennawardana23/graphql-pba/internal/util/exception"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validate,
	}
}

func (v *Validator) ValidateStruct(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var messages []string
		for _, e := range validationErrors {
			messages = append(messages, translateFieldError(e))
		}

		return &exception.CustomError{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: strings.Join(messages, "; "),
		}
	}
	return nil
}

func translateFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s: must be provided", field)
	case "email":
		return fmt.Sprintf("%s: must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s: must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s: must not exceed %s characters", field, e.Param())
	default:
		return fmt.Sprintf("%s: failed validation: %s", field, e.Tag())
	}
}
