package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"acacia/packages/db"
	"acacia/packages/schemas"
	"acacia/packages/testutils"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProject(t *testing.T) {
	ctx := context.Background()

	// Get the global database container
	dbContainer, err := testutils.GetGlobalDatabaseContainer(ctx)
	require.NoError(t, err)

	t.Run("should create project successfully", func(t *testing.T) {
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		logger := logrus.New()
		server, err := testutils.NewTestServer(queries, logger)

		require.NoError(t, err)

		// Start server in goroutine
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Prepare request body
		createReq := schemas.CreateProjectInput{
			Name: "Test Project",
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make HTTP request to the server

		url := fmt.Sprintf("%s%s", server.GetURL(), "/projects")
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Parse response body
		var project db.Project
		err = json.NewDecoder(resp.Body).Decode(&project)
		require.NoError(t, err)

		// Assert response body
		assert.NotZero(t, project.ID)
		assert.Equal(t, "Test Project", project.Name)
		assert.NotZero(t, project.CreatedAt)
		assert.NotZero(t, project.UpdatedAt)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		logger := logrus.New()
		server, err := testutils.NewTestServer(queries, logger)

		require.NoError(t, err)
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Make HTTP request with invalid JSON
		url := fmt.Sprintf("%s%s", server.GetURL(), "/projects")
		resp, err := http.Post(url, "application/json", bytes.NewBufferString("{invalid json"))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for empty name", func(t *testing.T) {
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		logger := logrus.New()
		server, err := testutils.NewTestServer(queries, logger)

		require.NoError(t, err)
		// Start server in goroutine
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Prepare request body with empty name
		createReq := schemas.CreateProjectInput{
			Name: "",
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make HTTP request
		url := fmt.Sprintf("%s%s", server.GetURL(), "/projects")
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for name that is too long", func(t *testing.T) {
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		logger := logrus.New()
		server, err := testutils.NewTestServer(queries, logger)

		require.NoError(t, err)
		// Start server in goroutine
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Prepare request body with name that exceeds max length (256 characters)
		longName := strings.Repeat("a", 256)
		createReq := schemas.CreateProjectInput{
			Name: longName,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make HTTP request
		resp, err := http.Post(server.GetURL()+"/projects", "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetProjectByID(t *testing.T) {
	ctx := context.Background()

	// Get the global database container
	dbContainer, err := testutils.GetGlobalDatabaseContainer(ctx)
	require.NoError(t, err)

	t.Run("should get project by ID successfully", func(t *testing.T) {
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		logger := logrus.New()
		server, err := testutils.NewTestServer(queries, logger)

		require.NoError(t, err)

		// Start server in goroutine
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// First create a project to retrieve
		createReq := schemas.CreateProjectInput{
			Name: "Test Project for Get",
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Create project
		createURL := fmt.Sprintf("%s%s", server.GetURL(), "/projects")
		createResp, err := http.Post(createURL, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createdProject db.Project
		err = json.NewDecoder(createResp.Body).Decode(&createdProject)
		require.NoError(t, err)

		// Now get the project by ID
		getURL := fmt.Sprintf("%s/projects/%d", server.GetURL(), createdProject.ID)
		getResp, err := http.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		// Parse response body
		var retrievedProject db.Project
		err = json.NewDecoder(getResp.Body).Decode(&retrievedProject)
		require.NoError(t, err)

		// Assert response body
		assert.Equal(t, createdProject.ID, retrievedProject.ID)
		assert.Equal(t, "Test Project for Get", retrievedProject.Name)
		assert.Equal(t, createdProject.CreatedAt, retrievedProject.CreatedAt)
		assert.Equal(t, createdProject.UpdatedAt, retrievedProject.UpdatedAt)
	})

	t.Run("should return 404 for non-existent project", func(t *testing.T) {
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		logger := logrus.New()
		server, err := testutils.NewTestServer(queries, logger)

		require.NoError(t, err)
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Try to get non-existent project
		getURL := fmt.Sprintf("%s/projects/999999", server.GetURL())
		getResp, err := http.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
	})

	t.Run("should return 400 for invalid project ID", func(t *testing.T) {
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		logger := logrus.New()
		server, err := testutils.NewTestServer(queries, logger)

		require.NoError(t, err)
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Try to get project with invalid ID
		getURL := fmt.Sprintf("%s/projects/invalid", server.GetURL())
		getResp, err := http.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, getResp.StatusCode)
	})
}
