package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"acacia/packages/auth"
	"acacia/packages/db"
	"acacia/packages/httperr"
	"acacia/packages/schemas"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ConversationsController struct {
	queries   *db.Queries
	logger    *logrus.Logger
	validator *validator.Validate
}

func NewConversationsController(queries *db.Queries, logger *logrus.Logger) *ConversationsController {
	return &ConversationsController{
		queries:   queries,
		logger:    logger,
		validator: validator.New(),
	}
}

// CreateConversation creates a new conversation
func (c *ConversationsController) CreateConversation(w http.ResponseWriter, r *http.Request) error {
	userID, ok := auth.GetUserID(r)
	if !ok {
		return httperr.WithStatus(errors.New("Unauthorized"), http.StatusUnauthorized)
	}

	var req schemas.CreateConversationInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	// Validate input
	if err := c.validator.Struct(&req); err != nil {
		return httperr.WithStatus(schemas.HandleConversationValidationErrors(err), http.StatusBadRequest)
	}

	// Generate title from first 35 characters of initial message
	title := req.InitialMessage
	if len(title) > 35 {
		title = title[:35] + "..."
	}

	// Create conversation
	createParams := db.CreateConversationParams{
		UserID:   userID,
		Provider: req.Provider,
		Model:    req.Model,
		Title:    title,
	}

	conversation, err := c.queries.CreateConversation(r.Context(), createParams)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create conversation")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Build response
	response := schemas.ConversationResponse{
		ID:        conversation.ID,
		UserID:    conversation.UserID,
		Title:     conversation.Title,
		Provider:  conversation.Provider,
		Model:     conversation.Model,
		CreatedAt: conversation.CreatedAt,
		UpdatedAt: conversation.UpdatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	return nil
}

// SendMessage sends a message to an existing conversation
func (c *ConversationsController) SendMessage(w http.ResponseWriter, r *http.Request) error {
	userID, ok := auth.GetUserID(r)
	if !ok {
		return httperr.WithStatus(errors.New("Unauthorized"), http.StatusUnauthorized)
	}

	var req schemas.SendMessageInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	// Validate input
	if err := c.validator.Struct(&req); err != nil {
		return httperr.WithStatus(schemas.HandleConversationValidationErrors(err), http.StatusBadRequest)
	}

	// Verify conversation exists and belongs to user
	conversation, err := c.queries.GetConversationByID(r.Context(), req.ConversationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("Conversation not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to get conversation")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Verify user owns this conversation
	if conversation.UserID != userID {
		return httperr.WithStatus(errors.New("Forbidden: insufficient permissions"), http.StatusForbidden)
	}

	// Create user message
	messageParams := db.CreateMessageParams{
		ConversationID: req.ConversationID,
		Role:           "user",
		Content:        req.Content,
	}

	_, err = c.queries.CreateMessage(r.Context(), messageParams)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create message")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

// GetLatestConversation retrieves the latest conversation with all messages
func (c *ConversationsController) GetLatestConversation(w http.ResponseWriter, r *http.Request) error {
	userID, ok := auth.GetUserID(r)
	if !ok {
		return httperr.WithStatus(errors.New("Unauthorized"), http.StatusUnauthorized)
	}

	// Get latest conversation
	conversation, err := c.queries.GetLatestConversationByUser(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("No conversations found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to get latest conversation")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Get all messages for this conversation
	messages, err := c.queries.GetMessagesByConversationID(r.Context(), conversation.ID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get messages")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Build response
	messageResponses := make([]schemas.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		messageResponses = append(messageResponses, schemas.MessageResponse{
			ID:             msg.ID,
			ConversationID: msg.ConversationID,
			Role:           msg.Role,
			Content:        msg.Content,
			SequenceNumber: msg.SequenceNumber,
			CreatedAt:      msg.CreatedAt,
		})
	}

	response := schemas.ConversationWithMessagesResponse{
		Conversation: schemas.ConversationResponse{
			ID:        conversation.ID,
			UserID:    conversation.UserID,
			Title:     conversation.Title,
			Provider:  conversation.Provider,
			Model:     conversation.Model,
			CreatedAt: conversation.CreatedAt,
			UpdatedAt: conversation.UpdatedAt,
		},
		Messages: messageResponses,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return nil
}
