package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"acacia/packages/db"
	"acacia/packages/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteProjectStatusColumn(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should delete project status column successfully", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "column1@example.com", "Test User", "password123")

		// Get user and create team
		user, err := setup.Queries.GetUserByEmail(ctx, "column1@example.com")
		require.NoError(t, err)
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, user.ID, "Test Team")

		// Create a test project first
		project, err := setup.Queries.CreateProject(ctx, db.CreateProjectParams{
			Name:   "Test Project",
			TeamID: teamID,
		})
		require.NoError(t, err)

		// Create two columns (need at least 2 to delete one)
		column1, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 1",
		})
		require.NoError(t, err)

		column2, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 2",
		})
		require.NoError(t, err)

		// Delete the first column
		deleteURL := fmt.Sprintf("%s/project-columns/%d", setup.Server.GetURL(), column1.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify column was deleted by trying to get it
		_, err = setup.Queries.GetProjectStatusColumnByID(ctx, column1.ID)
		assert.Error(t, err) // Should not be found

		// Verify other column still exists
		remainingColumn, err := setup.Queries.GetProjectStatusColumnByID(ctx, column2.ID)
		require.NoError(t, err)
		assert.Equal(t, column2.ID, remainingColumn.ID)
	})

	t.Run("should return 400 when trying to delete the last column", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "column2@example.com", "Test User", "password123")

		// Get user and create team
		user, err := setup.Queries.GetUserByEmail(ctx, "column2@example.com")
		require.NoError(t, err)
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, user.ID, "Test Team")

		// Create a test project
		project, err := setup.Queries.CreateProject(ctx, db.CreateProjectParams{
			Name:   "Test Project",
			TeamID: teamID,
		})
		require.NoError(t, err)

		// Create only one column
		column, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Only Column",
		})
		require.NoError(t, err)

		// Try to delete the only column
		deleteURL := fmt.Sprintf("%s/project-columns/%d", setup.Server.GetURL(), column.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify column still exists
		existingColumn, err := setup.Queries.GetProjectStatusColumnByID(ctx, column.ID)
		require.NoError(t, err)
		assert.Equal(t, column.ID, existingColumn.ID)
	})

	t.Run("should return 404 for non-existent column", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "column3@example.com", "Test User", "password123")

		// Try to delete non-existent column
		deleteURL := fmt.Sprintf("%s/project-columns/999999", setup.Server.GetURL())
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid column ID", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "column4@example.com", "Test User", "password123")

		// Try to delete column with invalid ID
		deleteURL := fmt.Sprintf("%s/project-columns/invalid", setup.Server.GetURL())
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should reorder columns correctly after deletion", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create authenticated client
		client := testutils.CreateAuthenticatedClient(t, setup, "column5@example.com", "Test User", "password123")

		// Get user and create team
		user, err := setup.Queries.GetUserByEmail(ctx, "column5@example.com")
		require.NoError(t, err)
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, user.ID, "Test Team")

		// Create a test project
		project, err := setup.Queries.CreateProject(ctx, db.CreateProjectParams{
			Name:   "Test Project",
			TeamID: teamID,
		})
		require.NoError(t, err)

		// Create three columns
		column1, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 1",
		})
		require.NoError(t, err)

		column2, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 2",
		})
		require.NoError(t, err)

		column3, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 3",
		})
		require.NoError(t, err)

		// Delete the middle column (column2)
		deleteURL := fmt.Sprintf("%s/project-columns/%d", setup.Server.GetURL(), column2.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify columns were reordered correctly
		remainingColumns, err := setup.Queries.GetProjectStatusColumnsByProjectID(ctx, int32(project.ID))
		require.NoError(t, err)
		assert.Len(t, remainingColumns, 2)

		// Column 1 should still be at position 0
		assert.Equal(t, column1.ID, remainingColumns[0].ID)
		assert.Equal(t, int16(0), remainingColumns[0].PositionIndex)

		// Column 3 should now be at position 1 (shifted from position 2)
		assert.Equal(t, column3.ID, remainingColumns[1].ID)
		assert.Equal(t, int16(1), remainingColumns[1].PositionIndex)
	})

	t.Run("should return 403 when trying to delete column from project user is not member of", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create first user and their team with a project and columns
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

		// Create two columns (need at least 2)
		column1, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 1",
		})
		require.NoError(t, err)

		_, err = setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 2",
		})
		require.NoError(t, err)

		// Create second user (not part of team1)
		client2 := testutils.CreateAuthenticatedClient(t, setup, "user2@example.com", "User 2", "password123")

		// User 2 tries to delete Team 1's column (should fail)
		deleteURL := fmt.Sprintf("%s/project-columns/%d", setup.Server.GetURL(), column1.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		resp, err := client2.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 403 Forbidden
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify column still exists
		existingColumn, err := setup.Queries.GetProjectStatusColumnByID(ctx, column1.ID)
		require.NoError(t, err)
		assert.Equal(t, column1.ID, existingColumn.ID)
	})
}

