package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"acacia/packages/db"
	"acacia/packages/schemas"
	"acacia/packages/testutils"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueAuthorization(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("should return 403 when trying to create issue in column from team user is not member of", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create first user and their team with a project and column
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

		// Create a column
		column, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "To Do",
		})
		require.NoError(t, err)

		// Create second user (not part of team1)
		client2 := testutils.CreateAuthenticatedClient(t, setup, "user2@example.com", "User 2", "password123")

		// User 2 tries to create an issue in Team 1's column (should fail)
		description := "This should not be created"
		createReq := schemas.CreateIssueInput{
			Name:        "Unauthorized Issue",
			Description: &description,
			ColumnId:    column.ID,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/issues", setup.Server.GetURL())
		resp, err := client2.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 403 Forbidden
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 403 when trying to get issue from team user is not member of", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create first user and their team with a project, column, and issue
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

		// Create a column
		column, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "To Do",
		})
		require.NoError(t, err)

		// Create an issue
		issue, err := setup.Queries.CreateIssue(ctx, db.CreateIssueParams{
			Name:        "Team 1 Issue",
			ColumnID:    column.ID,
			Description: null.StringFrom("This is a team 1 issue"),
		})
		require.NoError(t, err)

		// Create second user (not part of team1)
		client2 := testutils.CreateAuthenticatedClient(t, setup, "user2@example.com", "User 2", "password123")

		// User 2 tries to get Team 1's issue (should fail)
		getURL := fmt.Sprintf("%s/issues/%d", setup.Server.GetURL(), issue.ID)
		getResp, err := client2.Get(getURL)
		require.NoError(t, err)
		defer getResp.Body.Close()

		// Assert response status is 403 Forbidden
		assert.Equal(t, http.StatusForbidden, getResp.StatusCode)
	})

	t.Run("should return 403 when trying to update issue from team user is not member of", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create first user and their team with a project, column, and issue
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

		// Create a column
		column, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "To Do",
		})
		require.NoError(t, err)

		// Create an issue
		issue, err := setup.Queries.CreateIssue(ctx, db.CreateIssueParams{
			Name:        "Team 1 Issue",
			ColumnID:    column.ID,
			Description: null.StringFrom("This is a team 1 issue"),
		})
		require.NoError(t, err)

		// Create second user (not part of team1)
		client2 := testutils.CreateAuthenticatedClient(t, setup, "user2@example.com", "User 2", "password123")

		// User 2 tries to update Team 1's issue (should fail)
		updateReq := schemas.UpdateIssueInput{
			ID:          issue.ID,
			Name:        "Hacked Title",
			Description: "Hacked Description",
			ColumnId:    column.ID,
		}
		reqBody, err := json.Marshal(updateReq)
		require.NoError(t, err)

		updateURL := fmt.Sprintf("%s/issues", setup.Server.GetURL())
		req, err := http.NewRequest("PUT", updateURL, bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client2.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 403 Forbidden
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify issue was not updated
		existingIssue, err := setup.Queries.GetIssueByID(ctx, issue.ID)
		require.NoError(t, err)
		assert.Equal(t, "Team 1 Issue", existingIssue.Name)
	})

	t.Run("should return 403 when trying to delete issue from team user is not member of", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create first user and their team with a project, column, and issue
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

		// Create a column
		column, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "To Do",
		})
		require.NoError(t, err)

		// Create an issue
		issue, err := setup.Queries.CreateIssue(ctx, db.CreateIssueParams{
			Name:        "Team 1 Issue",
			ColumnID:    column.ID,
			Description: null.StringFrom("This is a team 1 issue"),
		})
		require.NoError(t, err)

		// Create second user (not part of team1)
		client2 := testutils.CreateAuthenticatedClient(t, setup, "user2@example.com", "User 2", "password123")

		// User 2 tries to delete Team 1's issue (should fail)
		deleteURL := fmt.Sprintf("%s/issues/%d", setup.Server.GetURL(), issue.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		resp, err := client2.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 403 Forbidden
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify issue still exists
		existingIssue, err := setup.Queries.GetIssueByID(ctx, issue.ID)
		require.NoError(t, err)
		assert.Equal(t, issue.ID, existingIssue.ID)
	})

	t.Run("should allow team member to create issue successfully", func(t *testing.T) {
		t.Parallel()
		setup := testutils.WithIntegrationTestSetup(ctx, t)
		defer setup.Cleanup()

		// Create user and team with a project and column
		client := testutils.CreateAuthenticatedClient(t, setup, "user1@example.com", "User 1", "password123")
		user, err := setup.Queries.GetUserByEmail(ctx, "user1@example.com")
		require.NoError(t, err)
		teamID := testutils.CreateTeamAndAddUser(t, ctx, setup, user.ID, "Team 1")

		// Create a project
		project, err := setup.Queries.CreateProject(ctx, db.CreateProjectParams{
			Name:   "Team 1 Project",
			TeamID: teamID,
		})
		require.NoError(t, err)

		// Create a column
		column, err := setup.Queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "To Do",
		})
		require.NoError(t, err)

		// User creates an issue in their team's column (should succeed)
		description := "This is my issue"
		createReq := schemas.CreateIssueInput{
			Name:        "My Issue",
			Description: &description,
			ColumnId:    column.ID,
		}
		reqBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		url := fmt.Sprintf("%s/issues", setup.Server.GetURL())
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 201 Created
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Parse response body
		var issue db.Issue
		err = json.NewDecoder(resp.Body).Decode(&issue)
		require.NoError(t, err)

		// Assert issue was created correctly
		assert.NotZero(t, issue.ID)
		assert.Equal(t, "My Issue", issue.Name)
		assert.Equal(t, "This is my issue", issue.Description.String)
		assert.Equal(t, column.ID, issue.ColumnID)
	})
}
