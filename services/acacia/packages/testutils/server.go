package testutils

import (
	"acacia/packages/config"
	"acacia/packages/db"
	"acacia/packages/schemas"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/guregu/null"
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
	// Generate a 32-byte encryption key for testing
	encryptionKey := []byte("test-encryption-key-32-bytes!!!!")

	env := &config.Environment{
		Port:          "8080",
		DatabaseURL:   "test",
		Env:           "test",
		JWTSecret:     "test-secret-key-for-testing-only",
		EncryptionKey: encryptionKey,
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

// authenticatedTransport wraps http.RoundTripper to add cookies to each request
type authenticatedTransport struct {
	cookies []*http.Cookie
	base    http.RoundTripper
}

func (t *authenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add cookies to the request
	for _, cookie := range t.cookies {
		req.AddCookie(cookie)
	}
	return t.base.RoundTrip(req)
}

// CreateAuthenticatedClient registers and logs in a test user, returning an HTTP client that automatically
// includes authentication cookies in all requests
func CreateAuthenticatedClient(t *testing.T, setup *IntegrationTestSetup, email, name, password string) *http.Client {
	// Register the user
	registerReq := schemas.RegisterUserInput{
		Email:    email,
		Name:     name,
		Password: password,
	}
	reqBody, err := json.Marshal(registerReq)
	require.NoError(t, err)

	registerURL := fmt.Sprintf("%s/users/register", setup.Server.GetURL())
	registerResp, err := http.Post(registerURL, "application/json", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	defer registerResp.Body.Close()
	require.Equal(t, http.StatusCreated, registerResp.StatusCode)

	// Login to get authentication cookies
	loginReq := schemas.LoginUserInput{
		Email:    email,
		Password: password,
	}
	loginBody, err := json.Marshal(loginReq)
	require.NoError(t, err)

	loginURL := fmt.Sprintf("%s/users/login", setup.Server.GetURL())
	loginResp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(loginBody))
	require.NoError(t, err)
	defer loginResp.Body.Close()
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	// Extract cookies from response
	cookies := loginResp.Cookies()

	// Create client with custom transport that adds cookies to every request
	return &http.Client{
		Transport: &authenticatedTransport{
			cookies: cookies,
			base:    http.DefaultTransport,
		},
	}
}

// CreateTeamAndAddUser creates a team and adds the specified user to it
// Returns the team ID
func CreateTeamAndAddUser(t *testing.T, ctx context.Context, setup *IntegrationTestSetup, userID int64, teamName string) int64 {
	// Create team
	team, err := setup.Queries.CreateTeam(ctx, db.CreateTeamParams{
		Name:        teamName,
		Description: null.StringFrom("Test team"),
	})
	require.NoError(t, err)

	// Add user to team
	_, err = setup.Queries.AddTeamMember(ctx, db.AddTeamMemberParams{
		TeamID: team.ID,
		UserID: userID,
	})
	require.NoError(t, err)

	return team.ID
}
