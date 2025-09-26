package testutils

import (
	"acacia/packages/config"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
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
	server := config.NewServer(d, l)
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
