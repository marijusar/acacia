package schemas

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type CreateConversationInput struct {
	Provider       string `json:"provider" validate:"required"`
	Model          string `json:"model" validate:"required"`
	InitialMessage string `json:"initial_message" validate:"required,min=1,max=10000"`
}

type SendMessageInput struct {
	ConversationID int64  `json:"conversation_id" validate:"required"`
	Content        string `json:"content" validate:"required,min=1,max=10000"`
}

type MessageResponse struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	SequenceNumber int32     `json:"sequence_number"`
	CreatedAt      time.Time `json:"created_at"`
}

type ConversationResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Provider  string    `json:"provider"`
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type ConversationWithMessagesResponse struct {
	Conversation ConversationResponse `json:"conversation"`
	Messages     []MessageResponse    `json:"messages"`
}

func HandleConversationValidationErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.New("Validation failed")
	}

	for _, e := range validationErrors {
		switch e.Field() {
		case "Content":
			if e.Tag() == "required" {
				return errors.New("Message content is required")
			}
			return errors.New("Message content must be between 1 and 10000 characters")
		case "InitialMessage":
			if e.Tag() == "required" {
				return errors.New("Initial message is required")
			}
			return errors.New("Initial message must be between 1 and 10000 characters")
		case "Provider":
			return errors.New("Provider is required")
		case "Model":
			return errors.New("Model is required")
		case "ConversationID":
			return errors.New("Conversation ID is required")
		default:
			return errors.New("Validation failed")
		}
	}
	return errors.New("Validation failed")
}
