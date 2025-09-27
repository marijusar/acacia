package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"acacia/packages/config"
	"acacia/packages/db"
	"acacia/packages/testutils"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteProjectStatusColumn(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Get the global database container
	dbContainer, err := testutils.GetGlobalDatabaseContainer(ctx)
	require.NoError(t, err)

	t.Run("should delete project status column successfully", func(t *testing.T) {
		t.Parallel()
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		database := &config.Database{
			Queries: queries,
			Conn:    testDB.DB,
		}
		logger := logrus.New()
		server, err := testutils.NewTestServer(database, logger)

		require.NoError(t, err)

		// Start server in goroutine
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Create a test project first
		project, err := queries.CreateProject(ctx, "Test Project")
		require.NoError(t, err)

		// Create two columns (need at least 2 to delete one)
		column1, err := queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 1",
		})
		require.NoError(t, err)

		column2, err := queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 2",
		})
		require.NoError(t, err)

		// Delete the first column
		deleteURL := fmt.Sprintf("%s/project-columns/%d", server.GetURL(), column1.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify column was deleted by trying to get it
		_, err = queries.GetProjectStatusColumnByID(ctx, column1.ID)
		assert.Error(t, err) // Should not be found

		// Verify other column still exists
		remainingColumn, err := queries.GetProjectStatusColumnByID(ctx, column2.ID)
		require.NoError(t, err)
		assert.Equal(t, column2.ID, remainingColumn.ID)
	})

	t.Run("should return 400 when trying to delete the last column", func(t *testing.T) {
		t.Parallel()
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		database := &config.Database{
			Queries: queries,
			Conn:    testDB.DB,
		}
		logger := logrus.New()
		server, err := testutils.NewTestServer(database, logger)

		require.NoError(t, err)
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Create a test project
		project, err := queries.CreateProject(ctx, "Test Project")
		require.NoError(t, err)

		// Create only one column
		column, err := queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Only Column",
		})
		require.NoError(t, err)

		// Try to delete the only column
		deleteURL := fmt.Sprintf("%s/project-columns/%d", server.GetURL(), column.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify column still exists
		existingColumn, err := queries.GetProjectStatusColumnByID(ctx, column.ID)
		require.NoError(t, err)
		assert.Equal(t, column.ID, existingColumn.ID)
	})

	t.Run("should return 404 for non-existent column", func(t *testing.T) {
		t.Parallel()
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		database := &config.Database{
			Queries: queries,
			Conn:    testDB.DB,
		}
		logger := logrus.New()
		server, err := testutils.NewTestServer(database, logger)

		require.NoError(t, err)
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Try to delete non-existent column
		deleteURL := fmt.Sprintf("%s/project-columns/999999", server.GetURL())
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid column ID", func(t *testing.T) {
		t.Parallel()
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		database := &config.Database{
			Queries: queries,
			Conn:    testDB.DB,
		}
		logger := logrus.New()
		server, err := testutils.NewTestServer(database, logger)

		require.NoError(t, err)
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Try to delete column with invalid ID
		deleteURL := fmt.Sprintf("%s/project-columns/invalid", server.GetURL())
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status is 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should reorder columns correctly after deletion", func(t *testing.T) {
		t.Parallel()
		// Create a fresh test database
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		require.NoError(t, err)
		defer testDB.Destroy(ctx)

		// Set up queries and server
		queries := db.New(testDB.DB)
		database := &config.Database{
			Queries: queries,
			Conn:    testDB.DB,
		}
		logger := logrus.New()
		server, err := testutils.NewTestServer(database, logger)

		require.NoError(t, err)
		err = server.StartServer()
		require.NoError(t, err)
		defer server.Close()

		// Wait for server to start
		time.Sleep(100 * time.Millisecond)

		// Create a test project
		project, err := queries.CreateProject(ctx, "Test Project")
		require.NoError(t, err)

		// Create three columns
		column1, err := queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 1",
		})
		require.NoError(t, err)

		column2, err := queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 2",
		})
		require.NoError(t, err)

		column3, err := queries.CreateProjectStatusColumn(ctx, db.CreateProjectStatusColumnParams{
			ProjectID: int32(project.ID),
			Name:      "Column 3",
		})
		require.NoError(t, err)

		// Delete the middle column (column2)
		deleteURL := fmt.Sprintf("%s/project-columns/%d", server.GetURL(), column2.ID)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify columns were reordered correctly
		remainingColumns, err := queries.GetProjectStatusColumnsByProjectID(ctx, int32(project.ID))
		require.NoError(t, err)
		assert.Len(t, remainingColumns, 2)

		// Column 1 should still be at position 0
		assert.Equal(t, column1.ID, remainingColumns[0].ID)
		assert.Equal(t, int16(0), remainingColumns[0].PositionIndex)

		// Column 3 should now be at position 1 (shifted from position 2)
		assert.Equal(t, column3.ID, remainingColumns[1].ID)
		assert.Equal(t, int16(1), remainingColumns[1].PositionIndex)
	})
}

