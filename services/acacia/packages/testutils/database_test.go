package testutils

import (
	"context"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// TestMain sets up and tears down the global database container
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Setup happens once for all tests
	_, err := GetGlobalDatabaseContainer(ctx)
	if err != nil {
		panic("Failed to setup database container: " + err.Error())
	}

	// Run all tests
	code := m.Run()

	// Cleanup happens once after all tests
	if err := CleanupGlobalContainer(ctx); err != nil {
		panic("Failed to cleanup database container: " + err.Error())
	}

	os.Exit(code)
}

func TestDatabaseContainer(t *testing.T) {
	ctx := context.Background()

	// Get the global container (already created in TestMain)
	dbContainer, err := GetGlobalDatabaseContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to get database container: %v", err)
	}

	// Verify container is running
	if !dbContainer.Verify(ctx) {
		t.Fatal("Database container is not running")
	}

	// Test basic connection
	if err := dbContainer.DB.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Create a new test database
	testDB, err := dbContainer.CreateNewDatabase(ctx)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer testDB.Destroy(ctx)

	// Test the test database connection
	if err := testDB.DB.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Verify tables exist (from migrations)
	var count int
	err = testDB.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'issues'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query for issues table: %v", err)
	}

	if count != 1 {
		t.Fatalf("Expected issues table to exist, but got count: %d", count)
	}

	t.Log("Database container test completed successfully")
}

func TestExampleUsage(t *testing.T) {
	ctx := context.Background()

	// Get the global container (reused across all tests)
	dbContainer, err := GetGlobalDatabaseContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to get database container: %v", err)
	}

	t.Run("test with fresh database", func(t *testing.T) {
		// Create a fresh database for this test
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		if err != nil {
			t.Fatalf("Failed to create test database: %v", err)
		}
		defer testDB.Destroy(ctx)

		// Your test logic here
		// testDB.DB contains a connection to a fresh database with all migrations applied

		// For example, insert some test data
		_, err = testDB.DB.ExecContext(ctx,
			"INSERT INTO issues (name, description) VALUES ($1, $2)",
			"Test Issue", "This is a test issue")
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}

		// Verify the data
		var name, description string
		err = testDB.DB.QueryRowContext(ctx,
			"SELECT name, description FROM issues WHERE name = $1", "Test Issue").
			Scan(&name, &description)
		if err != nil {
			t.Fatalf("Failed to query test data: %v", err)
		}

		if name != "Test Issue" || description != "This is a test issue" {
			t.Fatalf("Unexpected data: name=%s, description=%s", name, description)
		}
	})

	t.Run("another test with fresh database", func(t *testing.T) {
		// Each test gets its own fresh database from the same container
		testDB, err := dbContainer.CreateNewDatabase(ctx)
		if err != nil {
			t.Fatalf("Failed to create test database: %v", err)
		}
		defer testDB.Destroy(ctx)

		// This database is clean - no data from the previous test
		var count int
		err = testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM issues").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count issues: %v", err)
		}

		if count != 0 {
			t.Fatalf("Expected empty database, but found %d issues", count)
		}
	})
}

