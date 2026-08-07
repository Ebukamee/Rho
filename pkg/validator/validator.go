package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Errors(err error) []FieldError {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []FieldError{
			{
				Field:   "request",
				Message: "invalid request",
			},
		}
	}
	result := make([]FieldError, len(validationErrors))
	for _, fieldErr := range validationErrors {
		result = append(result, FieldError{
			Field:   fieldErr.Field(),
			Message: messageFor(fieldErr),
		})
	}
	return result
}

func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fe.Field(), fe.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is not valid", fe.Field())
	}
}
