package llm

import (
	"context"
	"errors"
)

var (
	ErrProviderNotSupported = errors.New("provider not supported")
	ErrAPIKeyNotFound       = errors.New("API key not found for provider")
	ErrInvalidAPIKey        = errors.New("invalid API key")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
)

// Message represents a chat message in a provider-agnostic format
type Message struct {
	Role    string `json:"role"`    // "user", "assistant", or "system"
	Content string `json:"content"` // The message content
}

// StreamChunk represents a chunk of streamed response
type StreamChunk struct {
	Content string // The text content of this chunk
	Done    bool   // Whether this is the final chunk
	Error   error  // Any error that occurred
}

// Provider defines the interface that all LLM providers must implement
type Provider interface {
	// StreamCompletion generates a completion with streaming support
	// Returns a channel that will receive chunks of the response
	StreamCompletion(ctx context.Context, messages []Message, model string) (<-chan StreamChunk, error)

	// GetProviderName returns the name of the provider (e.g., "openai", "anthropic")
	GetProviderName() string
}

// ProviderFactory creates provider instances with an API key
type ProviderFactory interface {
	New(apiKey string) Provider
}

// ProviderRegistry manages available LLM provider factories
type ProviderRegistry struct {
	factories map[string]ProviderFactory
}

// NewProviderRegistry creates a new provider registry with all provider factories
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		factories: map[string]ProviderFactory{
			"openai": &OpenAIProviderFactory{},
		},
	}
}

// GetProvider creates a provider instance with the given API key
func (r *ProviderRegistry) GetProvider(name string, apiKey string) (Provider, error) {
	factory, ok := r.factories[name]
	if !ok {
		return nil, ErrProviderNotSupported
	}
	return factory.New(apiKey), nil
}
