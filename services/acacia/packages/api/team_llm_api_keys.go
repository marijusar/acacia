package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"acacia/packages/crypto"
	"acacia/packages/db"
	"acacia/packages/httperr"
	"acacia/packages/schemas"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type TeamLLMAPIKeysController struct {
	queries           *db.Queries
	logger            *logrus.Logger
	validator         *validator.Validate
	encryptionService *crypto.EncryptionService
}

func NewTeamLLMAPIKeysController(
	queries *db.Queries,
	logger *logrus.Logger,
	encryptionService *crypto.EncryptionService,
) *TeamLLMAPIKeysController {
	return &TeamLLMAPIKeysController{
		queries:           queries,
		logger:            logger,
		validator:         validator.New(),
		encryptionService: encryptionService,
	}
}

// CreateOrUpdateAPIKey creates or updates an LLM API key for the team
func (c *TeamLLMAPIKeysController) CreateOrUpdateAPIKey(w http.ResponseWriter, r *http.Request) error {
	teamIDStr := chi.URLParam(r, "id")
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid team ID"), http.StatusBadRequest)
	}

	var req schemas.CreateTeamLLMAPIKeyInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httperr.WithStatus(errors.New("Invalid JSON"), http.StatusBadRequest)
	}

	if err := c.validator.Struct(&req); err != nil {
		return httperr.WithStatus(schemas.HandleTeamLLMAPIKeyValidationErrors(err), http.StatusBadRequest)
	}

	// Encrypt the API key
	encryptedKey, err := c.encryptionService.Encrypt(req.APIKey)
	if err != nil {
		c.logger.WithError(err).Error("Failed to encrypt API key")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Check if API key already exists for this team and provider
	existingKey, err := c.queries.GetTeamLLMAPIKeyByTeamID(r.Context(), db.GetTeamLLMAPIKeyByTeamIDParams{
		TeamID:   teamID,
		Provider: req.Provider,
	})

	// Handle unexpected errors
	if err != nil && err != sql.ErrNoRows {
		c.logger.WithError(err).Error("Failed to check existing API key")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	var apiKey db.TeamsLlmApiKey

	// Create new API key if it doesn't exist
	if err == sql.ErrNoRows {
		apiKey, err = c.queries.CreateTeamLLMAPIKey(r.Context(), db.CreateTeamLLMAPIKeyParams{
			TeamID:       teamID,
			Provider:     req.Provider,
			EncryptedKey: encryptedKey,
		})
		if err != nil {
			c.logger.WithError(err).Error("Failed to create API key")
			return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
		}
	} else {
		// Update existing API key
		apiKey, err = c.queries.UpdateTeamLLMAPIKey(r.Context(), db.UpdateTeamLLMAPIKeyParams{
			ID:           existingKey.ID,
			EncryptedKey: encryptedKey,
		})
		if err != nil {
			c.logger.WithError(err).Error("Failed to update API key")
			return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
		}
	}

	response := schemas.TeamLLMAPIKeyStatusResponse{
		ID:        apiKey.ID,
		Provider:  apiKey.Provider,
		IsActive:  apiKey.IsActive.Bool,
		CreatedAt: apiKey.CreatedAt,
		UpdatedAt: apiKey.UpdatedAt,
	}

	if apiKey.LastUsedAt.Valid {
		response.LastUsedAt = &apiKey.LastUsedAt.Time
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	return nil
}

// GetAPIKeys returns all LLM API keys for the team (without the actual keys)
func (c *TeamLLMAPIKeysController) GetAPIKeys(w http.ResponseWriter, r *http.Request) error {
	teamIDStr := chi.URLParam(r, "id")
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid team ID"), http.StatusBadRequest)
	}

	// Get all API keys for this team
	apiKeys, err := c.queries.GetAllTeamLLMAPIKeys(r.Context(), teamID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get API keys")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	// Build response without exposing actual API keys
	response := make(schemas.TeamLLMAPIKeysListResponse, 0, len(apiKeys))
	for _, key := range apiKeys {
		item := schemas.TeamLLMAPIKeyStatusResponse{
			ID:        key.ID,
			Provider:  key.Provider,
			IsActive:  key.IsActive.Bool,
			CreatedAt: key.CreatedAt,
			UpdatedAt: key.UpdatedAt,
		}
		if key.LastUsedAt.Valid {
			item.LastUsedAt = &key.LastUsedAt.Time
		}
		response = append(response, item)
	}

	json.NewEncoder(w).Encode(response)
	return nil
}

// DeleteAPIKey deletes an LLM API key
func (c *TeamLLMAPIKeysController) DeleteAPIKey(w http.ResponseWriter, r *http.Request) error {
	apiKeyIDStr := chi.URLParam(r, "keyId")
	apiKeyID, err := strconv.ParseInt(apiKeyIDStr, 10, 64)
	if err != nil {
		return httperr.WithStatus(errors.New("Invalid API key ID"), http.StatusBadRequest)
	}

	err = c.queries.DeleteTeamLLMAPIKey(r.Context(), apiKeyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return httperr.WithStatus(errors.New("API key not found"), http.StatusNotFound)
		}
		c.logger.WithError(err).Error("Failed to delete API key")
		return httperr.WithStatus(errors.New("Internal server error"), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
