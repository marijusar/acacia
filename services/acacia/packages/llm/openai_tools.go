package llm

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

// toolCallAccumulator is a simple struct to accumulate tool call data from streaming deltas
type toolCallAccumulator struct {
	id        string
	funcName  string
	arguments string
}

// toolCallResult holds the result of a tool call execution
type toolCallResult struct {
	toolCallID string
	result     string
	err        error
}

// StreamCompletionWithTools generates a completion with tool calling support
func (p *OpenAIProvider) StreamCompletionWithTools(
	ctx context.Context,
	messages []Message,
	model string,
) (<-chan StreamChunk, error) {
	// If no tools provided, fall back to regular streaming
	if p.tools == nil {
		return p.StreamCompletion(ctx, messages, model)
	}

	// Get tools from registry
	tools := p.tools.ListTools()
	if len(tools) == 0 {
		return p.StreamCompletion(ctx, messages, model)
	}

	// Convert tools to OpenAI format
	openaiTools := convertToolsToOpenAI(tools)

	// Create output channel
	out := make(chan StreamChunk)

	go func() {
		defer close(out)

		// Start with the original messages
		currentMessages := convertToOpenAIMessages(messages)

		// Tool calling loop - may need multiple rounds
		for {
			// Create streaming request with tools
			stream := p.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
				Messages: currentMessages,
				Model:    openai.ChatModel(model),
				Tools:    openaiTools,
			})

			var fullContent string
			var toolCallAccumulators []toolCallAccumulator

			// Process stream
			for stream.Next() {
				chunk := stream.Current()

				if len(chunk.Choices) > 0 {
					delta := chunk.Choices[0].Delta

					// Handle content
					if delta.Content != "" {
						fullContent += delta.Content
						select {
						case out <- StreamChunk{
							Content: delta.Content,
							Done:    false,
							Error:   nil,
						}:
						case <-ctx.Done():
							out <- StreamChunk{Done: true, Error: ctx.Err()}
							stream.Close()
							return
						}
					}

					// Handle tool calls - accumulate from deltas
					if len(delta.ToolCalls) > 0 {
						for _, tc := range delta.ToolCalls {
							// Find or create tool call accumulator
							if int(tc.Index) >= len(toolCallAccumulators) {
								// Expand slice
								for len(toolCallAccumulators) <= int(tc.Index) {
									toolCallAccumulators = append(toolCallAccumulators, toolCallAccumulator{})
								}
							}

							// Accumulate tool call data from delta
							if tc.ID != "" {
								toolCallAccumulators[tc.Index].id = tc.ID
							}
							if tc.Function.Name != "" {
								toolCallAccumulators[tc.Index].funcName += tc.Function.Name
							}
							if tc.Function.Arguments != "" {
								toolCallAccumulators[tc.Index].arguments += tc.Function.Arguments
							}
						}
					}
				}
			}

			stream.Close()

			// Check for errors
			if err := stream.Err(); err != nil {
				handleStreamError(out, err)
				return
			}

			// If no tool calls, we're done
			if len(toolCallAccumulators) == 0 {
				out <- StreamChunk{Content: "", Done: true, Error: nil}
				return
			}

			// Build tool call params from accumulated data
			toolCallParams := buildToolCallParams(toolCallAccumulators)

			// Add assistant message with tool calls to conversation
			currentMessages = appendAssistantMessageWithTools(currentMessages, fullContent, toolCallParams)

			// Execute each tool call and add results
			toolMessages, err := p.handleToolCalls(ctx, toolCallAccumulators)

			if err != nil {
				out <- StreamChunk{
					Content: "",
					Error:   err,
					Done:    true,
				}
			}

			currentMessages = append(currentMessages, toolMessages...)

			// Continue to next iteration to get LLM's response based on tool results
			fullContent = ""
			toolCallAccumulators = nil
		}
	}()

	return out, nil
}

func (p *OpenAIProvider) handleToolCalls(ctx context.Context, requestedToolCalls []toolCallAccumulator) ([]openai.ChatCompletionMessageParamUnion, error) {
	messages := []openai.ChatCompletionMessageParamUnion{}

	for _, acc := range requestedToolCalls {
		// Log that we're about to call the tool
		p.logger.WithField("tool_name", acc.funcName).Info("[TOOL_ORCHESTRATION] Starting tool execution")

		// Execute the tool call
		result := p.executeToolCall(ctx, acc)

		if result.err != nil {
			p.logger.WithError(result.err).WithField("tool_name", acc.funcName).Error("[TOOL_ORCHESTRATION] Tool execution failed")
			return nil, result.err
		}

		p.logger.WithField("tool_name", acc.funcName).Info("[TOOL_ORCHESTRATION] Tool execution completed successfully")

		messages = append(messages, openai.ToolMessage(result.result, result.toolCallID))
	}

	return messages, nil
}

// Helper functions

func convertToolsToOpenAI(tools []Tool) []openai.ChatCompletionToolParam {
	openaiTools := make([]openai.ChatCompletionToolParam, len(tools))

	for i, tool := range tools {
		openaiTools[i] = openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        tool.Name(),
				Description: param.NewOpt(tool.Description()),
				Parameters:  openai.FunctionParameters(tool.InputSchema()),
			},
			// Type will default to "function"
		}
	}

	return openaiTools
}

func convertToOpenAIMessages(messages []Message) []openai.ChatCompletionMessageParamUnion {
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
	return openaiMessages
}

func handleStreamError(out chan<- StreamChunk, err error) {
	var openaiErr *openai.Error
	if errors.As(err, &openaiErr) {
		switch openaiErr.StatusCode {
		case 401:
			out <- StreamChunk{Content: "", Done: true, Error: ErrInvalidAPIKey}
			return
		case 429:
			out <- StreamChunk{Content: "", Done: true, Error: ErrRateLimitExceeded}
			return
		}
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
		out <- StreamChunk{Content: "", Done: true, Error: nil}
		return
	}

	out <- StreamChunk{Content: "", Done: true, Error: err}
}

// executeToolCall executes a single tool call and returns the result
func (p *OpenAIProvider) executeToolCall(
	ctx context.Context,
	acc toolCallAccumulator,
) toolCallResult {
	toolName := acc.funcName

	p.logger.WithField("tool_name", toolName).Info("[TOOL_EXECUTION] Parsing tool arguments")

	// Parse arguments
	var args map[string]any
	if err := json.Unmarshal([]byte(acc.arguments), &args); err != nil {
		p.logger.WithError(err).WithField("tool_name", toolName).Error("[TOOL_EXECUTION] Failed to parse tool arguments")
		return toolCallResult{
			toolCallID: acc.id,
			err:        err,
		}
	}

	p.logger.WithField("tool_name", toolName).Info("[TOOL_EXECUTION] Getting tool from registry")

	// Get tool from registry
	tool, ok := p.tools.GetTool(toolName)
	if !ok {
		err := errors.New("unknown tool: " + toolName)
		p.logger.WithField("tool_name", toolName).Error("[TOOL_EXECUTION] Unknown tool")
		return toolCallResult{
			toolCallID: acc.id,
			err:        err,
		}
	}

	p.logger.WithField("tool_name", toolName).Info("[TOOL_EXECUTION] Executing tool")

	// Call tool directly with context (has user_id from auth middleware!)
	result, err := tool.Execute(ctx, args)
	if err != nil {
		p.logger.WithError(err).WithField("tool_name", toolName).Error("[TOOL_EXECUTION] Tool execution failed")
		return toolCallResult{
			toolCallID: acc.id,
			err:        err,
		}
	}

	// Convert result to string
	resultStr, ok := result.(string)
	if !ok {
		resultJSON, _ := json.Marshal(result)
		resultStr = string(resultJSON)
	}

	p.logger.WithField("tool_name", toolName).Info("[TOOL_EXECUTION] Tool execution completed successfully")

	return toolCallResult{
		toolCallID: acc.id,
		result:     resultStr,
		err:        nil,
	}
}

// buildToolCallParams constructs OpenAI tool call parameters from accumulated tool call data
func buildToolCallParams(accumulators []toolCallAccumulator) []openai.ChatCompletionMessageToolCallParam {
	toolCallParams := make([]openai.ChatCompletionMessageToolCallParam, len(accumulators))
	for i, acc := range accumulators {
		toolCallParams[i] = openai.ChatCompletionMessageToolCallParam{
			ID: acc.id,
			Function: openai.ChatCompletionMessageToolCallFunctionParam{
				Name:      acc.funcName,
				Arguments: acc.arguments,
			},
			// Type field will default to "function"
		}
	}
	return toolCallParams
}

// appendAssistantMessageWithTools adds an assistant message with tool calls to the conversation
func appendAssistantMessageWithTools(
	messages []openai.ChatCompletionMessageParamUnion,
	content string,
	toolCalls []openai.ChatCompletionMessageToolCallParam,
) []openai.ChatCompletionMessageParamUnion {
	// Create assistant message with tool calls
	assistantMsg := openai.ChatCompletionAssistantMessageParam{
		ToolCalls: toolCalls,
	}

	// Only add content if not empty
	if content != "" {
		assistantMsg.Content = openai.ChatCompletionAssistantMessageParamContentUnion{
			OfString: param.NewOpt(content),
		}
	}

	// Wrap in union type and append
	assistantMsgUnion := openai.ChatCompletionMessageParamUnion{
		OfAssistant: &assistantMsg,
	}

	return append(messages, assistantMsgUnion)
}
