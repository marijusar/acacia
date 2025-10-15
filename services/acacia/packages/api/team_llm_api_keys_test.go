package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"acacia/packages/schemas"
	"acacia/packages/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrUpdateTeamLLMAPIKey(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should create API key successfully", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client and user
		client := testutils.CreateAuthenticatedClient(t, setup, "test@example.com", "Test User", "password123")

		// Get user ID
		var userID int64
		err := setup.DB.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
		require.NoError(t, err)

		// Create team
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, userID, "Test Team")

		// Create API key
		createReq := schemas.CreateTeamLLMAPIKeyInput{
			Provider: "anthropic",
			APIKey:   "sk-ant-api03-test-key-1234567890",
		}
		reqBody, _ := json.Marshal(createReq)

		url := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Parse response
		var apiKeyResp schemas.TeamLLMAPIKeyStatusResponse
		err = json.NewDecoder(resp.Body).Decode(&apiKeyResp)
		require.NoError(t, err)

		// Assert response data
		assert.NotZero(t, apiKeyResp.ID)
		assert.Equal(t, "anthropic", apiKeyResp.Provider)
		assert.True(t, apiKeyResp.IsActive)
		assert.NotZero(t, apiKeyResp.CreatedAt)
		assert.NotZero(t, apiKeyResp.UpdatedAt)

		// Verify API key is stored encrypted in database
		var encryptedKey string
		err = setup.DB.DB.QueryRowContext(ctx,
			"SELECT encrypted_key FROM teams_llm_api_keys WHERE id = $1",
			apiKeyResp.ID).Scan(&encryptedKey)
		require.NoError(t, err)

		// Encrypted key should not match the original
		assert.NotEqual(t, "sk-ant-api03-test-key-1234567890", encryptedKey)
		assert.NotEmpty(t, encryptedKey)
	})

	t.Run("should update existing API key", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client and user
		client := testutils.CreateAuthenticatedClient(t, setup, "update@example.com", "Update User", "password123")

		// Get user ID
		var userID int64
		err := setup.DB.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", "update@example.com").Scan(&userID)
		require.NoError(t, err)

		// Create team
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, userID, "Update Team")

		// Create initial API key
		createReq := schemas.CreateTeamLLMAPIKeyInput{
			Provider: "anthropic",
			APIKey:   "sk-ant-api03-old-key",
		}
		reqBody, _ := json.Marshal(createReq)

		url := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp1, _ := client.Do(req)
		resp1.Body.Close()
		require.Equal(t, http.StatusCreated, resp1.StatusCode)

		// Update with new API key
		updateReq := schemas.CreateTeamLLMAPIKeyInput{
			Provider: "anthropic",
			APIKey:   "sk-ant-api03-new-key",
		}
		updateBody, _ := json.Marshal(updateReq)

		req2, _ := http.NewRequest("POST", url, bytes.NewBuffer(updateBody))
		req2.Header.Set("Content-Type", "application/json")

		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusCreated, resp2.StatusCode)

		// Verify only one API key exists for this team/provider
		var count int
		err = setup.DB.DB.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM teams_llm_api_keys WHERE team_id = $1 AND provider = $2",
			teamID, "anthropic").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "Should have exactly one API key after update")
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client and user
		client := testutils.CreateAuthenticatedClient(t, setup, "json@example.com", "JSON User", "password123")

		// Get user ID
		var userID int64
		err := setup.DB.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", "json@example.com").Scan(&userID)
		require.NoError(t, err)

		// Create team
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, userID, "JSON Team")

		url := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)
		req, _ := http.NewRequest("POST", url, bytes.NewBufferString("{invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for missing provider", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client and user
		client := testutils.CreateAuthenticatedClient(t, setup, "provider@example.com", "Provider User", "password123")

		// Get user ID
		var userID int64
		err := setup.DB.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", "provider@example.com").Scan(&userID)
		require.NoError(t, err)

		// Create team
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, userID, "Provider Team")

		createReq := schemas.CreateTeamLLMAPIKeyInput{
			Provider: "",
			APIKey:   "sk-ant-api03-test-key",
		}
		reqBody, _ := json.Marshal(createReq)

		url := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		assert.Contains(t, errResp["message"], "Provider")
	})

	t.Run("should return 400 for missing API key", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client and user
		client := testutils.CreateAuthenticatedClient(t, setup, "apikey@example.com", "APIKey User", "password123")

		// Get user ID
		var userID int64
		err := setup.DB.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", "apikey@example.com").Scan(&userID)
		require.NoError(t, err)

		// Create team
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, userID, "APIKey Team")

		createReq := schemas.CreateTeamLLMAPIKeyInput{
			Provider: "anthropic",
			APIKey:   "",
		}
		reqBody, _ := json.Marshal(createReq)

		url := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		assert.Contains(t, errResp["message"], "API key")
	})

	t.Run("should return 401 for unauthenticated request", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		createReq := schemas.CreateTeamLLMAPIKeyInput{
			Provider: "anthropic",
			APIKey:   "sk-ant-api03-test-key",
		}
		reqBody, _ := json.Marshal(createReq)

		url := fmt.Sprintf("%s/teams/1/llm-api-keys", setup.Server.GetURL())
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should support multiple providers for same team", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client and user
		client := testutils.CreateAuthenticatedClient(t, setup, "multi@example.com", "Multi User", "password123")

		// Get user ID
		var userID int64
		err := setup.DB.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", "multi@example.com").Scan(&userID)
		require.NoError(t, err)

		// Create team
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, userID, "Multi Team")

		url := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)

		// Create Anthropic API key
		anthropicReq := schemas.CreateTeamLLMAPIKeyInput{
			Provider: "anthropic",
			APIKey:   "sk-ant-api03-test-key",
		}
		reqBody1, _ := json.Marshal(anthropicReq)
		req1, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody1))
		req1.Header.Set("Content-Type", "application/json")
		resp1, _ := client.Do(req1)
		resp1.Body.Close()
		require.Equal(t, http.StatusCreated, resp1.StatusCode)

		// Create OpenAI API key
		openaiReq := schemas.CreateTeamLLMAPIKeyInput{
			Provider: "openai",
			APIKey:   "sk-proj-test-key-1234567890",
		}
		reqBody2, _ := json.Marshal(openaiReq)
		req2, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody2))
		req2.Header.Set("Content-Type", "application/json")
		resp2, _ := client.Do(req2)
		resp2.Body.Close()
		require.Equal(t, http.StatusCreated, resp2.StatusCode)

		// Verify both API keys exist
		var count int
		err = setup.DB.DB.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM teams_llm_api_keys WHERE team_id = $1",
			teamID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 2, count, "Should have two API keys for different providers")
	})
}

func TestGetTeamLLMAPIKeys(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should return all API keys for team", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client and user
		client := testutils.CreateAuthenticatedClient(t, setup, "get@example.com", "Get User", "password123")

		// Get user ID
		var userID int64
		err := setup.DB.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", "get@example.com").Scan(&userID)
		require.NoError(t, err)

		// Create team
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, userID, "Get Team")

		// Create multiple API keys
		createURL := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)

		providers := []string{"anthropic", "openai", "cohere"}
		for _, provider := range providers {
			createReq := schemas.CreateTeamLLMAPIKeyInput{
				Provider: provider,
				APIKey:   fmt.Sprintf("sk-%s-test-key", provider),
			}
			reqBody, _ := json.Marshal(createReq)
			req, _ := http.NewRequest("POST", createURL, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := client.Do(req)
			resp.Body.Close()
		}

		// Get all API keys
		getURL := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)
		req, _ := http.NewRequest("GET", getURL, nil)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var apiKeys schemas.TeamLLMAPIKeysListResponse
		err = json.NewDecoder(resp.Body).Decode(&apiKeys)
		require.NoError(t, err)

		// Assert response
		assert.Len(t, apiKeys, 3, "Should return all 3 API keys")

		// Verify providers
		providerMap := make(map[string]bool)
		for _, key := range apiKeys {
			providerMap[key.Provider] = true
			assert.NotZero(t, key.ID)
			assert.True(t, key.IsActive)
			assert.NotZero(t, key.CreatedAt)
		}

		assert.True(t, providerMap["anthropic"])
		assert.True(t, providerMap["openai"])
		assert.True(t, providerMap["cohere"])
	})

	t.Run("should return empty array for team with no API keys", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client and user
		client := testutils.CreateAuthenticatedClient(t, setup, "empty@example.com", "Empty User", "password123")

		// Get user ID
		var userID int64
		err := setup.DB.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", "empty@example.com").Scan(&userID)
		require.NoError(t, err)

		// Create team
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, userID, "Empty Team")

		url := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)
		req, _ := http.NewRequest("GET", url, nil)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var apiKeys schemas.TeamLLMAPIKeysListResponse
		err = json.NewDecoder(resp.Body).Decode(&apiKeys)
		require.NoError(t, err)

		assert.Len(t, apiKeys, 0)
	})

	t.Run("should return 401 for unauthenticated request", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		url := fmt.Sprintf("%s/teams/1/llm-api-keys", setup.Server.GetURL())
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestDeleteTeamLLMAPIKey(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should delete API key successfully", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client and user
		client := testutils.CreateAuthenticatedClient(t, setup, "delete@example.com", "Delete User", "password123")

		// Get user ID
		var userID int64
		err := setup.DB.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", "delete@example.com").Scan(&userID)
		require.NoError(t, err)

		// Create team
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, userID, "Delete Team")

		// Create API key
		createReq := schemas.CreateTeamLLMAPIKeyInput{
			Provider: "anthropic",
			APIKey:   "sk-ant-api03-test-key",
		}
		reqBody, _ := json.Marshal(createReq)

		createURL := fmt.Sprintf("%s/teams/%d/llm-api-keys", setup.Server.GetURL(), teamID)
		req, _ := http.NewRequest("POST", createURL, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		createResp, _ := client.Do(req)
		var apiKeyResp schemas.TeamLLMAPIKeyStatusResponse
		json.NewDecoder(createResp.Body).Decode(&apiKeyResp)
		createResp.Body.Close()

		// Delete API key
		deleteURL := fmt.Sprintf("%s/teams/%d/llm-api-keys/%d", setup.Server.GetURL(), teamID, apiKeyResp.ID)
		deleteReq, _ := http.NewRequest("DELETE", deleteURL, nil)

		deleteResp, err := client.Do(deleteReq)
		require.NoError(t, err)
		defer deleteResp.Body.Close()

		assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

		// Verify API key is deleted from database
		var count int
		err = setup.DB.DB.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM teams_llm_api_keys WHERE id = $1",
			apiKeyResp.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "API key should be deleted")
	})

	t.Run("should return 401 for unauthenticated request", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		url := fmt.Sprintf("%s/teams/1/llm-api-keys/1", setup.Server.GetURL())
		req, _ := http.NewRequest("DELETE", url, nil)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
