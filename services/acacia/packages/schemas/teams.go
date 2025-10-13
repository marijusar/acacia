package schemas

import (
	"acacia/packages/db"
	"errors"

	"github.com/go-playground/validator/v10"
)

type CreateTeamInput struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type TeamResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GetUserTeamsResponse []db.Team

// HandleTeamValidationErrors converts validator errors to user-friendly messages
func HandleTeamValidationErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.New("Validation failed")
	}

	for _, e := range validationErrors {
		switch e.Field() {
		case "Name":
			if e.Tag() == "required" {
				return errors.New("Team name is required")
			}
			return errors.New("Team name must be between 1 and 255 characters")
		default:
			return errors.New("Validation failed")
		}
	}

	return errors.New("Validation failed")
}
