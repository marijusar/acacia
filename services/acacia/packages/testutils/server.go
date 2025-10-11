package testutils

import (
	"acacia/packages/config"
	"acacia/packages/db"
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type PortManager struct {
	mu    *sync.Mutex
	min   int
	max   int
	ports map[int]bool
}

func NewPortManager(min int, max int) *PortManager {
	p := make(map[int]bool)

	for i := min; i <= max; i++ {
		p[i] = false
	}

	return &PortManager{
		mu:    &sync.Mutex{},
		min:   min,
		max:   max,
		ports: p,
	}
}

func (p *PortManager) GetPort() (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := p.min; i <= p.max; i++ {
		if p.ports[i] == false {
			p.ports[i] = true
			return i, nil
		}

	}
	return -1, errors.New("Error while retrieving empty port. All ports taken.")
}

func (p *PortManager) ReleasePort(port int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ports[port] == false {
		return errors.New("Error while returning port to the pool. Port is already returned")
	}

	p.ports[port] = false

	return nil

}

var (
	portManager     *PortManager
	portManagerOnce sync.Once
)

// TestServer wraps a config.Server with port management for testing
type TestServer struct {
	*config.Server
	port   int
	logger *logrus.Logger
}

// initPortPool initializes the port pool with 1000 ports (9000-9999)
// This range avoids common Unix process ports and well-known service ports
func initPortPool() {
	portManagerOnce.Do(func() {
		portManager = NewPortManager(9000, 9999)
	})
}

// getPort retrieves a port from the pool
func getPort() (int, error) {
	if portManager == nil {
		initPortPool()
	}
	return portManager.GetPort()

}

// returnPort returns a port to the pool
func returnPort(port int) error {
	return portManager.ReleasePort(port)
}

// NewTestServer creates a new test server with automatic port allocation
// Accepts the same parameters as config.NewServer: db.Queries and logrus.Logger
func NewTestServer(d *config.Database, l *logrus.Logger) (*TestServer, error) {
	// Create test environment
	env := &config.Environment{
		Port:        "8080",
		DatabaseURL: "test",
		Env:         "test",
		JWTSecret:   "test-secret-key-for-testing-only",
	}

	server := config.NewServer(d, l, env)
	port, err := getPort()

	if err != nil {
		return nil, err
	}

	fmt.Printf("Taking port %d\n", port)

	return &TestServer{
		Server: server,
		port:   port,
	}, nil
}

// StartServer starts the HTTP server on the allocated port
func (ts *TestServer) StartServer() error {
	// We need to access the router from the Server struct
	// Since it's not exported, we'll use the existing ListenAndServe method pattern
	// but run it in a goroutine with our allocated port

	go func() {
		ts.Server.ListenAndServe(strconv.Itoa(ts.port))
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	return nil
}

// GetPort returns the allocated port number
func (ts *TestServer) GetPort() int {
	return ts.port
}

// GetURL returns the full URL to the server
func (ts *TestServer) GetURL() string {
	return fmt.Sprintf("http://localhost:%d", ts.port)
}

// Close returns the port to the pool
// Note: The original config.Server.ListenAndServe doesn't provide graceful shutdown,
// so this mainly handles port cleanup. For production tests, consider implementing
// a wrapper with proper graceful shutdown if needed.
func (ts *TestServer) Close() {
	// Return port to pool
	ts.Server.Close()
	returnPort(ts.port)
}

// IntegrationTestSetup encapsulates the common setup for integration tests
type IntegrationTestSetup struct {
	Server  *TestServer
	DB      *TestDatabase
	Queries *db.Queries
	Cleanup func()
}

// WithIntegrationTestSetup sets up a fresh database and server for integration tests
// It handles all the boilerplate: creating test DB, setting up queries, starting server, etc.
// Returns a setup struct with server, DB, queries, and a cleanup function that must be called
func WithIntegrationTestSetup(ctx context.Context, t *testing.T) *IntegrationTestSetup {
	// Get the global database container
	dbContainer, err := GetGlobalDatabaseContainer(ctx)
	require.NoError(t, err)

	// Create a fresh test database
	testDB, err := dbContainer.CreateNewDatabase(ctx)
	require.NoError(t, err)

	// Set up queries and server
	queries := db.New(testDB.DB)
	database := &config.Database{
		Queries: queries,
		Conn:    testDB.DB,
	}
	logger := logrus.New()
	server, err := NewTestServer(database, logger)
	require.NoError(t, err)

	// Start server
	err = server.StartServer()
	require.NoError(t, err)

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Create cleanup function
	cleanup := func() {
		server.Close()
		testDB.Destroy(ctx)
	}

	return &IntegrationTestSetup{
		Server:  server,
		DB:      testDB,
		Queries: queries,
		Cleanup: cleanup,
	}
}
