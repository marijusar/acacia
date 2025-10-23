package llm

// ToolRegistry manages available LLM tools
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry creates a new tool registry with the given tools
func NewToolRegistry(tools []Tool) *ToolRegistry {
	registry := &ToolRegistry{
		tools: make(map[string]Tool),
	}

	for _, tool := range tools {
		registry.tools[tool.Name()] = tool
	}

	return registry
}

// GetTool retrieves a tool by name
func (r *ToolRegistry) GetTool(name string) (Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

// ListTools returns all registered tools
func (r *ToolRegistry) ListTools() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}
