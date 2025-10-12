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

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "test@example.com", "Test User", "password123")

		// Get user and create team
		user, err := setup.Queries.GetUserByEmail(ctx, "test@example.com")
		require.NoError(t, err)
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, user.ID, "Test Team")

		// Prepare request body
		createReq := schemas.CreateProjectInput{
			Name:   "Test Project",
			TeamID: teamID,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make authenticated HTTP request to the server
		url := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
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
		assert.Equal(t, teamID, project.TeamID)
		assert.NotZero(t, project.CreatedAt)
		assert.NotZero(t, project.UpdatedAt)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "test2@example.com", "Test User", "password123")

		// Make HTTP request with invalid JSON
		url := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
		resp, err := client.Post(url, "application/json", bytes.NewBufferString("{invalid json"))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for empty name", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "test3@example.com", "Test User", "password123")

		// Get user and create team
		user, err := setup.Queries.GetUserByEmail(ctx, "test3@example.com")
		require.NoError(t, err)
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, user.ID, "Test Team")

		// Prepare request body with empty name
		createReq := schemas.CreateProjectInput{
			Name:   "",
			TeamID: teamID,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make HTTP request
		url := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for name that is too long", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "test4@example.com", "Test User", "password123")

		// Get user and create team
		user, err := setup.Queries.GetUserByEmail(ctx, "test4@example.com")
		require.NoError(t, err)
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, user.ID, "Test Team")

		// Prepare request body with name that exceeds max length (256 characters)
		longName := strings.Repeat("a", 256)
		createReq := schemas.CreateProjectInput{
			Name:   longName,
			TeamID: teamID,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Make HTTP request
		resp, err := client.Post(setup.Server.GetURL()+"/projects", "application/json", bytes.NewBuffer(reqBody))
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

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "test5@example.com", "Test User", "password123")

		// Get user and create team
		user, err := setup.Queries.GetUserByEmail(ctx, "test5@example.com")
		require.NoError(t, err)
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, user.ID, "Test Team")

		// First create a project to retrieve
		createReq := schemas.CreateProjectInput{
			Name:   "Test Project for Get",
			TeamID: teamID,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		// Create project
		createURL := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
		createResp, err := client.Post(createURL, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createdProject db.Project
		err = json.NewDecoder(createResp.Body).Decode(&createdProject)
		require.NoError(t, err)

		// Now get the project by ID
		getURL := fmt.Sprintf("%s/projects/%d", setup.Server.GetURL(), createdProject.ID)
		getResp, err := client.Get(getURL)
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

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "test6@example.com", "Test User", "password123")

		// Try to get non-existent project
		getURL := fmt.Sprintf("%s/projects/999999", setup.Server.GetURL())
		getResp, err := client.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
	})

	t.Run("should return 400 for invalid project ID", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "test7@example.com", "Test User", "password123")

		// Try to get project with invalid ID
		getURL := fmt.Sprintf("%s/projects/invalid", setup.Server.GetURL())
		getResp, err := client.Get(getURL)
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

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "test8@example.com", "Test User", "password123")

		// Get user and create team
		user, err := setup.Queries.GetUserByEmail(ctx, "test8@example.com")
		require.NoError(t, err)
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, user.ID, "Test Team")

		// Create some test projects
		projects := []string{"Project 1", "Project 2", "Project 3"}
		createdProjects := make([]db.Project, 0, len(projects))

		for _, projectName := range projects {
			createReq := schemas.CreateProjectInput{
				Name:   projectName,
				TeamID: teamID,
			}
			reqBody, err := json.Marshal(createReq)
			require.NoError(t, err)

			// Create project
			createURL := fmt.Sprintf("%s%s", setup.Server.GetURL(), "/projects")
			createResp, err := client.Post(createURL, "application/json", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			defer createResp.Body.Close()

			var createdProject db.Project
			err = json.NewDecoder(createResp.Body).Decode(&createdProject)
			require.NoError(t, err)
			createdProjects = append(createdProjects, createdProject)
		}

		// Now get all projects
		getURL := fmt.Sprintf("%s/projects", setup.Server.GetURL())
		getResp, err := client.Get(getURL)
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

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "test9@example.com", "Test User", "password123")

		// Get all projects (should be empty)
		getURL := fmt.Sprintf("%s/projects", setup.Server.GetURL())
		getResp, err := client.Get(getURL)
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

	t.Run("should return 403 when trying to create project for team user is not member of", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create first user and their team
		_ = testutils.CreateAuthenticatedClient(t, setup, "user1@example.com", "User 1", "password123")
		user1, err := setup.Queries.GetUserByEmail(ctx, "user1@example.com")
		require.NoError(t, err)
		team1ID := testutils.CreateTeamAndAddUser(t, ctx, setup, user1.ID, "Team 1")

		// Create second user (not part of team1)
		client2 := testutils.CreateAuthenticatedClient(t, setup, "user2@example.com", "User 2", "password123")

		// User 2 tries to create a project for Team 1 (should fail)
		createReq := schemas.CreateProjectInput{
			Name:   "Unauthorized Project",
			TeamID: team1ID,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/projects", setup.Server.GetURL())
		resp, err := client2.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 403 Forbidden
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 403 when trying to get project from team user is not member of", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create first user and their team with a project
		_ = testutils.CreateAuthenticatedClient(t, setup, "user1@example.com", "User 1", "password123")
		user1, err := setup.Queries.GetUserByEmail(ctx, "user1@example.com")
		require.NoError(t, err)
		team1ID := testutils.CreateTeamAndAddUser(t, ctx, setup, user1.ID, "Team 1")

		// Create a project for team 1
		project, err := setup.Queries.CreateProject(ctx, db.CreateProjectParams{
			Name:   "Team 1 Project",
			TeamID: team1ID,
		})
		require.NoError(t, err)

		// Create second user (not part of team1)
		client2 := testutils.CreateAuthenticatedClient(t, setup, "user2@example.com", "User 2", "password123")

		// User 2 tries to get Team 1's project (should fail)
		getURL := fmt.Sprintf("%s/projects/%d", setup.Server.GetURL(), project.ID)
		getResp, err := client2.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status is 403 Forbidden
		assert.Equal(t, http.StatusForbidden, getResp.StatusCode)
	})

	t.Run("should return 403 when trying to update project from team user is not member of", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create first user and their team with a project
		_ = testutils.CreateAuthenticatedClient(t, setup, "user1@example.com", "User 1", "password123")
		user1, err := setup.Queries.GetUserByEmail(ctx, "user1@example.com")
		require.NoError(t, err)
		team1ID := testutils.CreateTeamAndAddUser(t, ctx, setup, user1.ID, "Team 1")

		// Create a project for team 1
		project, err := setup.Queries.CreateProject(ctx, db.CreateProjectParams{
			Name:   "Team 1 Project",
			TeamID: team1ID,
		})
		require.NoError(t, err)

		// Create second user (not part of team1)
		client2 := testutils.CreateAuthenticatedClient(t, setup, "user2@example.com", "User 2", "password123")

		// User 2 tries to update Team 1's project (should fail)
		updateReq := schemas.UpdateProjectInput{
			Name: "Hacked Project Name",
		}
		reqBody, err := json.Marshal(updateReq)
		require.NoError(t, err)

		updateURL := fmt.Sprintf("%s/projects/%d", setup.Server.GetURL(), project.ID)
		req, err := http.NewRequest("PUT", updateURL, bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client2.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 403 Forbidden
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 403 when trying to delete project from team user is not member of", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create first user and their team with a project
		_ = testutils.CreateAuthenticatedClient(t, setup, "user1@example.com", "User 1", "password123")
		user1, err := setup.Queries.GetUserByEmail(ctx, "user1@example.com")
		require.NoError(t, err)
		team1ID := testutils.CreateTeamAndAddUser(t, ctx, setup, user1.ID, "Team 1")

		// Create a project for team 1
		project, err := setup.Queries.CreateProject(ctx, db.CreateProjectParams{
			Name:   "Team 1 Project",
			TeamID: team1ID,
		})
		require.NoError(t, err)

		// Create second user (not part of team1)
		client2 := testutils.CreateAuthenticatedClient(t, setup, "user2@example.com", "User 2", "password123")

		// User 2 tries to delete Team 1's project (should fail)
		deleteURL := fmt.Sprintf("%s/projects/%d", setup.Server.GetURL(), project.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		resp, err := client2.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 403 Forbidden
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify project still exists
		existingProject, err := setup.Queries.GetProjectByID(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, project.ID, existingProject.ID)
	})
}
