package llm

import "context"

// Tool represents an LLM tool that can be called with arguments
type Tool interface {
	// Name returns the name of the tool
	Name() string

	// Description returns a description of what the tool does
	Description() string

	// InputSchema returns the JSON schema for the tool's input parameters
	InputSchema() map[string]any

	// Execute runs the tool with the given arguments
	// Context carries user_id and other request context from auth middleware
	Execute(ctx context.Context, args map[string]any) (any, error)
}
