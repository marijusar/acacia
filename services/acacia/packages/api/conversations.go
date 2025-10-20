package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"acacia/packages/auth"
	"acacia/packages/db"
	"acacia/packages/httperr"
	"acacia/packages/schemas"
	"acacia/packages/services"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ConversationsController struct {
	queries             *db.Queries
	logger              *logrus.Logger
	validator           *validator.Validate
	conversationService *services.ConversationService
}

func NewConversationsController(
	queries *db.Queries,
	logger *logrus.Logger,
	conversationService *services.ConversationService,
) *ConversationsController {
	return &ConversationsController{
		queries:             queries,
		logger:              logger,
		validator:           validator.New(),
		conversationService: conversationService,
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

	teamID, err := c.queries.GetTeamIDByProject(context.Background(), req.ProjectID)

	if err != nil {
		return httperr.WithStatus(errors.New("Cannot find project's team"), http.StatusInternalServerError)
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
		TeamID:   teamID,
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

// SendMessage sends a message to an existing conversation and streams the LLM response
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

	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get streaming channel from conversation service
	streamChan, err := c.conversationService.ReplyToMessage(r.Context(), req.ConversationID, req.Content)
	if err != nil {
		c.logger.WithError(err).Error("Failed to start reply stream")
		// Send error as SSE event
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return nil
	}

	// Stream chunks to client
	flusher, ok := w.(http.Flusher)
	if !ok {
		return httperr.WithStatus(errors.New("Streaming not supported"), http.StatusInternalServerError)
	}

	for chunk := range streamChan {
		if chunk.Error != nil {
			c.logger.WithError(chunk.Error).Error("Error in stream")
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", chunk.Error.Error())
			flusher.Flush()
			break
		}

		if chunk.Content != "" {
			// Send content chunk
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", chunk.Content)
			flusher.Flush()
		}

		if chunk.Done {
			// Send done event
			fmt.Fprintf(w, "event: done\ndata: \n\n")
			flusher.Flush()
			break
		}
	}

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
