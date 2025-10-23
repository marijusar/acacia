# Acacia MCP Server

This is the Model Context Protocol (MCP) server for the Acacia project management system. It provides LLM access to project data through a secure, HTTP-proxied architecture.

## Architecture

The MCP server acts as a **thin HTTP client** that proxies requests to the main Acacia API server, ensuring:
- **Single source of truth** for authorization (all checks happen in the API)
- **Token forwarding** - LLM passes user cookies to MCP server, which forwards them to the API
- **No database access** - MCP server has zero direct DB access, only HTTP calls
- **Reuses existing business logic** - all authorization middleware and API logic is centralized

```
User → LLM (with cookies) → MCP Server (forwards cookies) → HTTP API (auth/authz)
```

## Available Tools

The MCP server provides 4 read-only tools:

### 1. `get_user_projects`
Get list of all projects the authenticated user has access to.

**Input:** None
**Output:** Array of project summaries (ID, name, team ID, timestamps)

### 2. `get_project_details`
Get detailed information about a specific project including all columns and issues.

**Input:**
- `project_id` (integer, required): The ID of the project

**Output:** Project details with columns and nested issues

### 3. `get_issue_details`
Get detailed information about a specific issue including project and column context.

**Input:**
- `issue_id` (integer, required): The ID of the issue

**Output:** Full issue details with project name, column name, description

### 4. `search_issues`
Search and filter issues within a project.

**Input:**
- `project_id` (integer, required): The ID of the project to search in
- `column_id` (integer, optional): Filter by specific column
- `search_term` (string, optional): Search in issue name or description

**Output:** Array of matching issues

## Project Structure

```
services/acacia/
├── cmd/
│   └── mcp-server/
│       └── main.go              # MCP server entry point
├── packages/
│   ├── mcp/
│   │   ├── server.go           # MCP server setup and tool registration
│   │   ├── handlers.go         # Tool handler implementations
│   │   ├── client.go           # HTTP client for API calls
│   │   └── types.go            # Input/output type definitions
│   ├── db/                      # SHARED with API
│   ├── auth/                    # SHARED with API
│   ├── config/                  # SHARED with API
│   └── ... (other packages)
```

## Configuration

### Environment Variables

Add these to your `.env` file:

```bash
# MCP Server Configuration
MCP_PORT=8081                          # Not used with stdio transport
ACACIA_API_URL=http://localhost:8080  # Base URL of the Acacia API
```

### Building

```bash
# Build MCP server only
make build-mcp

# Build both API and MCP servers
make build

# Development mode (with hot reload)
make dev-mcp
```

### Running

#### Local Development

The MCP server uses **stdio transport** (standard for MCP servers), which means it communicates via stdin/stdout. This allows it to be invoked by MCP clients.

```bash
# Run the built binary
./bin/mcp-server

# Or in development mode
make dev-mcp
```

#### Docker

The MCP server is included in the docker-compose setup:

```bash
# Start all services (including MCP server)
docker compose up

# Start only MCP server and its dependencies
docker compose up mcp-server

# View MCP server logs
docker compose logs -f mcp-server

# Rebuild after code changes
docker compose build mcp-server
```

The MCP server will be available at `http://localhost:8083` and will automatically connect to the `acacia` API service within the Docker network.

## Authentication Flow

1. **LLM receives user request** with authentication cookies (access_token, refresh_token)
2. **LLM calls MCP tool** and passes cookies in the request context
3. **MCP server extracts cookies** from context
4. **MCP server makes HTTP request** to API with cookies in headers
5. **API authenticates user** using existing auth middleware
6. **API authorizes access** using existing authorization middleware
7. **API returns data** to MCP server
8. **MCP server formats response** and returns to LLM
9. **LLM uses data** to answer user's question

## Security

- **No direct database access** - MCP server cannot bypass API authorization
- **Cookie forwarding** - User authentication is preserved through the chain
- **API-level authorization** - All permission checks happen in the existing API layer
- **Read-only tools** - Current implementation only supports GET operations
- **HTTPS recommended** - Use HTTPS for API URL in production

## Extending

To add new tools:

1. **Define input/output types** in `packages/mcp/types.go`
2. **Implement handler** in `packages/mcp/handlers.go` using HTTP client
3. **Register tool** in `packages/mcp/server.go` using `mcp.AddTool`
4. **Rebuild** with `make build-mcp`

Example:

```go
// In types.go
type GetTeamMembersInput struct {
    TeamID int64 `json:"team_id" jsonschema:"description=The ID of the team,required"`
}

type GetTeamMembersOutput struct {
    Members []TeamMember `json:"members"`
}

// In handlers.go
func (h *ToolHandlers) GetTeamMembers(ctx context.Context, input GetTeamMembersInput) (*GetTeamMembersOutput, error) {
    cookies, err := h.extractCookies(ctx)
    if err != nil {
        return nil, err
    }

    path := fmt.Sprintf("/teams/%d/members", input.TeamID)
    respBody, err := h.apiClient.Get(path, cookies)
    // ... parse and return
}

// In server.go (registerTools function)
mcp.AddTool(s.mcpServer, &mcp.Tool{
    Name:        "get_team_members",
    Description: "Get list of team members",
}, func(ctx context.Context, req *mcp.CallToolRequest, input GetTeamMembersInput) (*mcp.CallToolResult, *GetTeamMembersOutput, error) {
    // ... implementation
})
```

## Deployment

### Development
Run both servers in separate terminals:
```bash
# Terminal 1: API server
make dev

# Terminal 2: MCP server
make dev-mcp
```

### Production
Deploy as separate binaries/containers:
- **API server**: Standard HTTP server on port 8080
- **MCP server**: Stdio-based MCP server (invoked by LLM client)

Both can share the same database and configuration.

## Troubleshooting

### "no authentication cookies found in context"
- Ensure LLM is passing cookies in the MCP request
- Check that cookies are being extracted and added to context

### "API request failed with status 401"
- Check that `ACACIA_API_URL` is correct
- Verify cookies are valid and not expired
- Ensure API server is running

### "failed to fetch X from API"
- Verify the API endpoint exists and returns expected data
- Check API logs for authorization failures
- Ensure user has permission to access the resource

## Future Enhancements

- [ ] Add write tools (create_issue, update_issue, etc.)
- [ ] Implement caching for frequently accessed data
- [ ] Add rate limiting per user
- [ ] Support for batch operations
- [ ] WebSocket transport for real-time updates
- [ ] Metrics and monitoring
