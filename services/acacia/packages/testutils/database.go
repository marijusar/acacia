package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	globalContainer *DatabaseContainer
	containerOnce   sync.Once
)

// DatabaseContainer wraps the PostgreSQL testcontainer with utility functions
type DatabaseContainer struct {
	Container     *postgres.PostgresContainer
	ConnectionURL string
	DB            *sql.DB
}

// TestDatabase represents a template-based test database
type TestDatabase struct {
	Name          string
	ConnectionURL string
	DB            *sql.DB
	container     *DatabaseContainer
}

// GetGlobalDatabaseContainer returns the singleton database container instance
// Creates it once and reuses for all tests
func GetGlobalDatabaseContainer(ctx context.Context) (*DatabaseContainer, error) {
	var setupErr error

	containerOnce.Do(func() {
		globalContainer, setupErr = setupDatabaseContainer(ctx)
	})

	return globalContainer, setupErr
}

// setupDatabaseContainer starts a PostgreSQL container and runs migrations (internal)
func setupDatabaseContainer(ctx context.Context) (*DatabaseContainer, error) {
	// Start PostgreSQL container
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:17-alpine"),
		postgres.WithDatabase("acacia"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("root"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Get connection string
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Wait for database to be ready
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbContainer := &DatabaseContainer{
		Container:     container,
		ConnectionURL: connStr,
		DB:            db,
	}

	// Run migrations
	if err := dbContainer.RunMigrations(ctx); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return dbContainer, nil
}

// RunMigrations executes all SQL files in the schema directory
func (dc *DatabaseContainer) RunMigrations(ctx context.Context) error {
	schemaDir := "../../schema"

	// Read all SQL files from schema directory
	err := filepath.WalkDir(schemaDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".up.sql") {
			return nil
		}

		// Read SQL file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		// Execute SQL
		if _, err := dc.DB.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", path, err)
		}

		return nil
	})

	return err
}

// CreateNewDatabase creates a new test database using the template
func (dc *DatabaseContainer) CreateNewDatabase(ctx context.Context) (*TestDatabase, error) {
	// Generate unique database name
	dbName := "test_" + strings.ReplaceAll(uuid.New().String(), "-", "_")

	// Create database using template
	query := fmt.Sprintf(`CREATE DATABASE %s TEMPLATE acacia`, pq.QuoteIdentifier(dbName))
	if _, err := dc.DB.ExecContext(ctx, query); err != nil {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}

	// Create connection string for new database
	connStr := strings.Replace(dc.ConnectionURL, "acacia", dbName, 1)

	// Connect to new database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	return &TestDatabase{
		Name:          dbName,
		ConnectionURL: connStr,
		DB:            db,
		container:     dc,
	}, nil
}

// Destroy drops the test database and closes connections
func (td *TestDatabase) Destroy(ctx context.Context) error {
	// Close database connection
	if td.DB != nil {
		td.DB.Close()
	}

	// Drop database
	query := fmt.Sprintf(`DROP DATABASE %s`, pq.QuoteIdentifier(td.Name))
	if _, err := td.container.DB.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to drop test database: %w", err)
	}

	return nil
}

// Close cleans up the database container
func (dc *DatabaseContainer) Close(ctx context.Context) error {
	if dc.DB != nil {
		dc.DB.Close()
	}
	if dc.Container != nil {
		return dc.Container.Terminate(ctx)
	}
	return nil
}

// CleanupGlobalContainer cleans up the global container
// Should be called when all tests are finished (e.g., in TestMain)
func CleanupGlobalContainer(ctx context.Context) error {
	if globalContainer != nil {
		return globalContainer.Close(ctx)
	}
	return nil
}

// Verify checks if the container is still running
func (dc *DatabaseContainer) Verify(ctx context.Context) bool {
	if dc.Container == nil {
		return false
	}

	state, err := dc.Container.State(ctx)
	if err != nil {
		return false
	}

	return state.Running
}
