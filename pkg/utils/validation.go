package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ParseValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// Just return the first error
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s is required", e.Field())
			case "email":
				return fmt.Sprintf("%s is not a valid email", e.Field())
			case "min":
				return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
			case "max":
				return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
			default:
				return fmt.Sprintf("%s is invalid", e.Field())
			}
		}
	}
	return err.Error()
}
