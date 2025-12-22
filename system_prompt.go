package gomind

import (
	"context"
	"encoding/json"
	"fmt"
)

// SystemPrompt fetches the recommended system prompt for LLM integration.
func (c *Client) SystemPrompt(ctx context.Context) (*SystemPromptResponse, error) {
	respBody, err := c.get(ctx, "/v1/system-prompt")
	if err != nil {
		c.logger.Error("Gomind SystemPrompt failed", "error", err)
		return nil, err
	}

	var resp APIResponse[SystemPromptResponse]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse system-prompt response: %w", err)
	}

	c.logger.Info("Gomind SystemPrompt success")

	return &resp.Result, nil
}
