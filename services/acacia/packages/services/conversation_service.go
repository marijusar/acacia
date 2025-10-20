package services

import (
	"acacia/packages/crypto"
	"acacia/packages/db"
	"acacia/packages/llm"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrAPIKeyNotFound       = errors.New("API key not found for provider")
	ErrInvalidProvider      = errors.New("invalid provider")
)

type ConversationService struct {
	queries           *db.Queries
	providerRegistry  *llm.ProviderRegistry
	encryptionService *crypto.EncryptionService
	logger            *logrus.Logger
}

func NewConversationService(
	queries *db.Queries,
	providerRegistry *llm.ProviderRegistry,
	encryptionService *crypto.EncryptionService,
	logger *logrus.Logger,
) *ConversationService {
	return &ConversationService{
		queries:           queries,
		providerRegistry:  providerRegistry,
		encryptionService: encryptionService,
		logger:            logger,
	}
}

// ReplyToMessage handles the complete chat flow:
// 1. Saves user message to DB
// 2. Gets conversation details and history
// 3. Gets team's API key for the provider
// 4. Streams LLM response
// 5. Saves assistant response after streaming completes
func (s *ConversationService) ReplyToMessage(
	ctx context.Context,
	conversationID int64,
	userMessage string,
) (<-chan llm.StreamChunk, error) {
	// Save user message to database
	_, err := s.queries.CreateMessage(ctx, db.CreateMessageParams{
		ConversationID: conversationID,
		Role:           "user",
		Content:        userMessage,
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to save user message")
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Get conversation details (provider, model, teamID)
	conversation, err := s.queries.GetConversationByID(ctx, conversationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrConversationNotFound
		}
		s.logger.WithError(err).Error("Failed to get conversation")
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Get team's encrypted API key for this provider
	apiKeyRecord, err := s.queries.GetTeamLLMAPIKeyByTeamID(ctx, db.GetTeamLLMAPIKeyByTeamIDParams{
		TeamID:   conversation.TeamID,
		Provider: conversation.Provider,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAPIKeyNotFound
		}
		s.logger.WithError(err).Error("Failed to get API key")
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	// Decrypt the API key
	decryptedKey, err := s.encryptionService.Decrypt(apiKeyRecord.EncryptedKey)
	if err != nil {
		s.logger.WithError(err).Error("Failed to decrypt API key")
		return nil, fmt.Errorf("failed to decrypt API key: %w", err)
	}

	// Get conversation message history
	dbMessages, err := s.queries.GetMessagesByConversationID(ctx, conversationID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get conversation history")
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// Convert database messages to LLM provider format
	messages := make([]llm.Message, 0, len(dbMessages))
	for _, msg := range dbMessages {
		messages = append(messages, llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Get LLM provider instance
	provider, err := s.providerRegistry.GetProvider(conversation.Provider, decryptedKey)
	if err != nil {
		s.logger.WithError(err).WithField("provider", conversation.Provider).Error("Failed to get provider")
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// Start streaming from LLM provider
	streamChan, err := provider.StreamCompletion(ctx, messages, conversation.Model)
	if err != nil {
		s.logger.WithError(err).Error("Failed to start streaming")
		return nil, fmt.Errorf("failed to start streaming: %w", err)
	}

	// Update last used timestamp for API key
	go func() {
		if err := s.queries.UpdateLastUsedAt(context.Background(), apiKeyRecord.ID); err != nil {
			s.logger.WithError(err).Warn("Failed to update API key last used timestamp")
		}
	}()

	// Create output channel and goroutine to save assistant response after streaming
	outChan := make(chan llm.StreamChunk)
	go func() {
		defer close(outChan)

		var fullResponse string
		var streamErr error

		// Forward chunks and collect full response
		for chunk := range streamChan {
			// Forward chunk to output
			outChan <- chunk

			if chunk.Error != nil {
				streamErr = chunk.Error
				break
			}

			// Accumulate response content
			fullResponse += chunk.Content

			// If this is the last chunk and no error, save to database
			if chunk.Done && chunk.Error == nil {
				// Save assistant's response to database
				_, err := s.queries.CreateMessage(context.Background(), db.CreateMessageParams{
					ConversationID: conversationID,
					Role:           "assistant",
					Content:        fullResponse,
				})
				if err != nil {
					s.logger.WithError(err).Error("Failed to save assistant message")
					// Send error chunk
					outChan <- llm.StreamChunk{
						Content: "",
						Done:    true,
						Error:   fmt.Errorf("failed to save assistant message: %w", err),
					}
				}
			}
		}

		// If there was a streaming error and we didn't already send a final chunk
		if streamErr != nil {
			s.logger.WithError(streamErr).Error("Streaming error occurred")
		}
	}()

	return outChan, nil
}
