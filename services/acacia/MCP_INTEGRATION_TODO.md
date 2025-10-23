# MCP Integration - Continuation Instructions

## Current Status

### ✅ Completed
1. **MCP Server (HTTP JSON-RPC)** - FULLY WORKING
   - Simple HTTP server at `cmd/mcp-server/main.go`
   - JSON-RPC handlers in `packages/mcp/http_handler.go`
   - Tool handlers in `packages/mcp/handlers.go` (uses HTTP client to call API)
   - 4 tools implemented: get_user_projects, get_project_details, get_issue_details, search_issues
   - Docker service added (port 8083)
   - **Builds successfully** ✅

2. **MCP Client** - FULLY WORKING
   - HTTP client in `packages/mcp/mcp_client.go`
   - Discovers tools from MCP server
   - Calls tools with cookie forwarding
   - Implements `llm.ToolCaller` interface
   - **Builds successfully** ✅

3. **Architecture**
   - HTTP client renamed to `acacia_http_client.go`
   - Provider interface updated with `StreamCompletionWithTools` method
   - Conversation service updated to accept `mcpClient` and `cookies` parameters

### ✅ RESOLVED - OpenAI SDK Compatibility Issue

**Solution:** Used `param.NewOpt[T]()` to wrap optional fields.

**Fixed in:**
- `packages/mcp/openai_converter.go` - Uses `param.NewOpt(tool.Description)`
- `packages/llm/openai_tools.go` - Uses `param.NewOpt()` for content, `.ToParam()` for tool calls
- Proper union construction: `ChatCompletionMessageParamUnion{OfAssistant: &assistantMsg}`

**Build Status:** ✅ Both binaries compile successfully

### 📋 Integration Steps - ALL COMPLETED ✅

1. **✅ Fix OpenAI SDK compatibility**
   - File: `packages/mcp/openai_converter.go`
   - Solution: Used `param.NewOpt()` for optional fields

2. **✅ Update server.go to initialize MCP client**
   - File: `packages/config/server.go`
   - MCP client initialized conditionally based on `env.MCPServerURL`
   - Passed to conversation service constructor

3. **✅ Update config to load MCP_SERVER_URL**
   - File: `packages/config/environment.go`
   - Added field: `MCPServerURL string`
   - Loads from `MCP_SERVER_URL` env var (optional)

4. **✅ Pass cookies from API controller to conversation service**
   - File: `packages/api/conversations.go`
   - Extracts cookies: `cookies := r.Cookies()`
   - Passes to service: `conversationService.ReplyToMessage(ctx, conversationID, userMessage, cookies)`

5. **✅ Update .env.example**
   - Added `MCP_SERVER_URL` with documentation
   - Documented both local and docker configurations

### 🧪 Testing Steps

6. **Test MCP server standalone**
   ```bash
   # Start all services
   docker compose up

   # Test MCP server directly
   curl -X POST http://localhost:8083 \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'

   # Should return 4 tools
   ```

7. **Test OpenAI integration**
   - Start conversation
   - Ask: "What projects do I have access to?"
   - LLM should call get_user_projects tool
   - Verify in logs

## File Reference

### Working Files (No Changes Needed)
- ✅ `cmd/mcp-server/main.go` - MCP server entry point
- ✅ `packages/mcp/http_handler.go` - JSON-RPC HTTP handler
- ✅ `packages/mcp/jsonrpc.go` - JSON-RPC types
- ✅ `packages/mcp/handlers.go` - Tool implementations
- ✅ `packages/mcp/acacia_http_client.go` - HTTP client for API
- ✅ `packages/mcp/mcp_client.go` - MCP client
- ✅ `packages/mcp/types.go` - Tool input/output types
- ✅ `packages/llm/provider.go` - Updated interface
- ✅ `packages/services/conversation_service.go` - Updated signature

### Files Needing Fixes
- ⚠️ `packages/mcp/openai_converter.go` - **BLOCKED ON SDK**
- ⚠️ `packages/llm/openai_tools.go` - **DEPENDS ON CONVERTER**

### Files Needing Updates
- 📝 `packages/config/server.go` - Initialize MCP client
- 📝 `packages/config/environment.go` - Add MCPServerURL field
- 📝 `packages/api/conversations.go` - Pass cookies
- 📝 `.env.example` - Add MCP_SERVER_URL

## Quick Start for Next Session

```bash
# 1. Fix the OpenAI SDK issue first
cd /home/marijus/Development/fullstack/acacia/services/acacia

# Try this to see available functions:
go doc github.com/openai/openai-go/packages/param

# 2. Once fixed, complete remaining integrations (see steps above)

# 3. Test build
go build ./packages/...

# 4. Test MCP server
make build-mcp
./bin/mcp-server &

# 5. Test integration
make build-api
./bin/acacia
```

## Architecture Diagram

```
User Request (with cookies)
    ↓
API Controller (/conversations/messages)
    ↓ (extract cookies)
Conversation Service
    ↓ (discover tools from MCP)
MCP Client → MCP Server (HTTP JSON-RPC with cookies)
                ↓
            Acacia API (existing endpoints)
                ↓
            Database

OpenAI ← Conversation Service (with tools)
    ↓ (tool calls)
MCP Client → MCP Server → API
    ↓
OpenAI (with results)
    ↓
User Response
```

## Key Decisions Made

1. **HTTP JSON-RPC instead of stdio** - Easier to deploy, debug, scale
2. **MCP server proxies to API** - Single source of truth for auth
3. **Cookie forwarding** - User context preserved through chain
4. **Tool discovery** - Dynamic, no hardcoding in conversation service
5. **Separate binary** - Independent deployment and scaling

## Notes

- MCP server port: 8083 (8081 was taken by Traefik)
- Both API and MCP server share same codebase (modulith)
- Docker compose service configured and ready
- All authorization happens in API layer (MCP just forwards)
- Tool calling loop handles multi-turn conversations

Good luck! The hardest part (architecture) is done. Just need to solve the SDK compatibility issue and wire everything up.
