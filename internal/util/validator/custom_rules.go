package validator

import (
	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidations adds custom validation rules to the validator
func RegisterCustomValidations(v *validator.Validate) {
	// Register custom validations
	_ = v.RegisterValidation("customtag", customValidation)
	_ = v.RegisterValidation("phone", validatePhone)
	// Add more custom validations as needed
}

// Example custom validation
func customValidation(fl validator.FieldLevel) bool {
	// Your custom validation logic
	return true
}

// Phone number validation
func validatePhone(fl validator.FieldLevel) bool {
	// Phone validation logic
	return true
}
