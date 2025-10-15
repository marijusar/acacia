package schemas

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type CreateTeamLLMAPIKeyInput struct {
	Provider string `json:"provider" validate:"required,min=1,max=50"`
	APIKey   string `json:"api_key" validate:"required,min=1"`
}

type UpdateTeamLLMAPIKeyInput struct {
	APIKey string `json:"api_key" validate:"required,min=1"`
}

type TeamLLMAPIKeyStatusResponse struct {
	ID         int64      `json:"id"`
	Provider   string     `json:"provider"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

type TeamLLMAPIKeysListResponse []TeamLLMAPIKeyStatusResponse

// HandleTeamLLMAPIKeyValidationErrors converts validator errors to user-friendly messages
func HandleTeamLLMAPIKeyValidationErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.New("Validation failed")
	}

	for _, e := range validationErrors {
		switch e.Field() {
		case "Provider":
			if e.Tag() == "required" {
				return errors.New("Provider is required")
			}
			return errors.New("Provider must be between 1 and 50 characters")
		case "APIKey":
			if e.Tag() == "required" {
				return errors.New("API key is required")
			}
			return errors.New("API key cannot be empty")
		default:
			return errors.New("Validation failed")
		}
	}

	return errors.New("Validation failed")
}
