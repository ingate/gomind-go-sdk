package gomind

import (
	"context"
	"encoding/json"
	"fmt"
)

// MindRequest is the request body for the mind endpoint.
type MindRequest struct {
	Prompt       string            `json:"prompt"`
	Context      map[string]string `json:"context,omitempty"`
	OutputSchema map[string]any    `json:"output_schema"`
}

// MindResponse is the response from the mind endpoint.
type MindResponse struct {
	Result map[string]any `json:"result"`
	Meta   MindMeta       `json:"meta"`
}

// MindMeta contains metadata about the mind request execution.
type MindMeta struct {
	TokensUsed int `json:"tokens_used"`
	LatencyMs  int `json:"latency_ms"`
}

// Mind makes a natural language request that is fulfilled by an internal LLM agent.
// The agent uses Gomind's MCP server to query and manipulate the knowledge graph automatically.
func (c *Client) Mind(ctx context.Context, prompt string, context_ map[string]string, outputSchema map[string]any) (*MindResponse, error) {
	req := MindRequest{
		Prompt:       prompt,
		Context:      context_,
		OutputSchema: outputSchema,
	}

	respBody, err := c.post(ctx, "/v1/mind", req)
	if err != nil {
		c.logger.Error("Gomind Mind failed", "error", err)
		return nil, err
	}

	var resp APIResponse[MindResponse]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse mind response: %w", err)
	}

	c.logger.Info("Gomind Mind success",
		"tokensUsed", resp.Result.Meta.TokensUsed,
		"latencyMs", resp.Result.Meta.LatencyMs,
	)

	return &resp.Result, nil
}
