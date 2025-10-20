package llm

import (
	"context"
	"errors"
	"io"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	client *openai.Client
}

// OpenAIProviderFactory implements ProviderFactory for OpenAI
type OpenAIProviderFactory struct{}

// New creates a new OpenAI provider with the given API key
func (f *OpenAIProviderFactory) New(apiKey string) Provider {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIProvider{
		client: &client,
	}
}

// GetProviderName returns "openai"
func (p *OpenAIProvider) GetProviderName() string {
	return "openai"
}

// StreamCompletion generates a completion with streaming support
func (p *OpenAIProvider) StreamCompletion(ctx context.Context, messages []Message, model string) (<-chan StreamChunk, error) {
	// Convert our messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	for _, msg := range messages {
		switch msg.Role {
		case "user":
			openaiMessages = append(openaiMessages, openai.UserMessage(msg.Content))
		case "assistant":
			openaiMessages = append(openaiMessages, openai.AssistantMessage(msg.Content))
		case "system":
			openaiMessages = append(openaiMessages, openai.SystemMessage(msg.Content))
		}
	}

	// Create the streaming request
	stream := p.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openaiMessages,
		Model:    model,
	})

	// Create output channel
	out := make(chan StreamChunk)

	// Start goroutine to handle streaming
	go func() {
		defer close(out)
		defer stream.Close()

		var fullContent string

		for stream.Next() {
			chunk := stream.Current()

			// Check if there's content in the delta
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				content := chunk.Choices[0].Delta.Content
				fullContent += content

				select {
				case out <- StreamChunk{
					Content: content,
					Done:    false,
					Error:   nil,
				}:
				case <-ctx.Done():
					out <- StreamChunk{
						Content: "",
						Done:    true,
						Error:   ctx.Err(),
					}
					return
				}
			}
		}

		// Check for streaming errors
		if err := stream.Err(); err != nil {
			// Handle specific OpenAI errors
			var openaiErr *openai.Error
			if errors.As(err, &openaiErr) {
				switch openaiErr.StatusCode {
				case 401:
					out <- StreamChunk{
						Content: "",
						Done:    true,
						Error:   ErrInvalidAPIKey,
					}
					return
				case 429:
					out <- StreamChunk{
						Content: "",
						Done:    true,
						Error:   ErrRateLimitExceeded,
					}
					return
				}
			}

			// Check for context cancellation or EOF
			if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
				out <- StreamChunk{
					Content: "",
					Done:    true,
					Error:   nil,
				}
				return
			}

			// Generic error
			out <- StreamChunk{
				Content: "",
				Done:    true,
				Error:   err,
			}
			return
		}

		// Send final chunk indicating completion
		out <- StreamChunk{
			Content: "",
			Done:    true,
			Error:   nil,
		}
	}()

	return out, nil
}
