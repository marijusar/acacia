package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"acacia/packages/db"
	"acacia/packages/schemas"
	"acacia/packages/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProject(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should create project successfully", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Prepare request body
		createReq := schemas.CreateProjectInput{
			Name: "Test Project",
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make HTTP request to the server

		url := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
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
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Make HTTP request with invalid JSON
		url := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
		resp, err := http.Post(url, "application/json", bytes.NewBufferString("{invalid json"))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for empty name", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Prepare request body with empty name
		createReq := schemas.CreateProjectInput{
			Name: "",
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make HTTP request
		url := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for name that is too long", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Prepare request body with name that exceeds max length (256 characters)
		longName := strings.Repeat("a", 256)
		createReq := schemas.CreateProjectInput{
			Name: longName,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make HTTP request
		resp, err := http.Post(setup.Server.GetURL()+"/projects", "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetProjectByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should get project by ID successfully", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// First create a project to retrieve
		createReq := schemas.CreateProjectInput{
			Name: "Test Project for Get",
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Create project
		createURL := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
		createResp, err := http.Post(createURL, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createdProject db.Project
		err = json.NewDecoder(createResp.Body).Decode(&createdProject)
		require.NoError(t, err)

		// Now get the project by ID
		getURL := fmt.Sprintf("%s/projects/%d", setup.Server.GetURL(), createdProject.ID)
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
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Try to get non-existent project
		getURL := fmt.Sprintf("%s/projects/999999", setup.Server.GetURL())
		getResp, err := http.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
	})

	t.Run("should return 400 for invalid project ID", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Try to get project with invalid ID
		getURL := fmt.Sprintf("%s/projects/invalid", setup.Server.GetURL())
		getResp, err := http.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, getResp.StatusCode)
	})
}

func TestGetProjects(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should get all projects successfully", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create some test projects
		projects := []string{"Project 1", "Project 2", "Project 3"}
		createdProjects := make([]db.Project, 0, len(projects))

		for _, projectName := range projects {
			createReq := schemas.CreateProjectInput{
				Name: projectName,
			}
			reqBody, err := json.Marshal(createReq)
			require.NoError(t, err)

			// Create project
			createURL := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
			createResp, err := http.Post(createURL, "application/json", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			defer createResp.Body.Close()

			var createdProject db.Project
			err = json.NewDecoder(createResp.Body).Decode(&createdProject)
			require.NoError(t, err)
			createdProjects = append(createdProjects, createdProject)
		}

		// Now get all projects
		getURL := fmt.Sprintf("%s/projects", setup.Server.GetURL())
		getResp, err := http.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		// Parse response body
		var retrievedProjects []db.Project
		err = json.NewDecoder(getResp.Body).Decode(&retrievedProjects)
		require.NoError(t, err)

		// Assert response body
		assert.Len(t, retrievedProjects, 3)

		// Check that all created projects are in the response
		for _, created := range createdProjects {
			found := false
			for _, retrieved := range retrievedProjects {
				if created.ID == retrieved.ID {
					assert.Equal(t, created.Name, retrieved.Name)
					assert.Equal(t, created.CreatedAt, retrieved.CreatedAt)
					assert.Equal(t, created.UpdatedAt, retrieved.UpdatedAt)
					found = true
					break
				}
			}
			assert.True(t, found, "Created project with ID %d should be found in response", created.ID)
		}
	})

	t.Run("should return empty array when no projects exist", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Get all projects (should be empty)
		getURL := fmt.Sprintf("%s/projects", setup.Server.GetURL())
		getResp, err := http.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		// Parse response body
		var retrievedProjects []db.Project
		err = json.NewDecoder(getResp.Body).Decode(&retrievedProjects)
		require.NoError(t, err)

		// Assert response body is empty array
		assert.Len(t, retrievedProjects, 0)
	})
}
