package schemas

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

// HandleValidationErrors converts validator errors to user-friendly messages
func HandleUserValidationErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.New("Validation failed")
	}

	for _, e := range validationErrors {
		switch e.Field() {
		case "Email":
			return errors.New("Invalid email format")
		case "Name":
			if e.Tag() == "required" {
				return errors.New("Name is required")
			}
			return errors.New("Name must be between 1 and 100 characters")
		case "Password":
			if e.Tag() == "required" {
				return errors.New("Password is required")
			} else if e.Tag() == "min" {
				return errors.New("Password must be at least 6 characters")
			} else if e.Tag() == "max" {
				return errors.New("Password must be at most 50 characters")
			}
		default:
			return errors.New("Validation failed")
		}
	}

	return errors.New("Validation failed")
}
