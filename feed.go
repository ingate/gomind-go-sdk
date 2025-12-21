package gomind

import (
	"context"
	"encoding/json"
	"fmt"
)

// Feed ingests raw content and automatically extracts facts using LLM.
func (c *Client) Feed(ctx context.Context, content string, source string) (*FeedResponse, error) {
	req := FeedRequest{
		Content: content,
		Source:  source,
	}

	respBody, err := c.post(ctx, "/v1/feed", req)
	if err != nil {
		c.logger.Error("Gomind Feed failed", "error", err)
		return nil, err
	}

	var resp FeedResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse feed response: %w", err)
	}

	c.logger.Info("Gomind Feed success",
		"status", resp.Status,
		"factsExtracted", resp.FactsExtracted,
	)

	return &resp, nil
}

// FeedMessages ingests conversation messages and automatically extracts facts.
func (c *Client) FeedMessages(ctx context.Context, messages []FeedMessage, source string) (*FeedResponse, error) {
	req := FeedRequest{
		Messages: messages,
		Source:   source,
	}

	respBody, err := c.post(ctx, "/v1/feed", req)
	if err != nil {
		c.logger.Error("Gomind FeedMessages failed", "error", err)
		return nil, err
	}

	var resp FeedResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse feed response: %w", err)
	}

	c.logger.Info("Gomind FeedMessages success",
		"status", resp.Status,
		"factsExtracted", resp.FactsExtracted,
		"messageCount", len(messages),
	)

	return &resp, nil
}

// FeedAsync ingests content asynchronously and returns a job ID.
// Use GetJobStatus to check the status of the job.
func (c *Client) FeedAsync(ctx context.Context, content string, source string) (string, error) {
	resp, err := c.Feed(ctx, content, source)
	if err != nil {
		return "", err
	}

	if resp.JobID != "" {
		return resp.JobID, nil
	}

	// If no job ID, it was processed synchronously
	return "", nil
}

// GetJobStatus checks the status of an async feed job.
func (c *Client) GetJobStatus(ctx context.Context, jobID string) (*JobStatusResponse, error) {
	respBody, err := c.get(ctx, fmt.Sprintf("/v1/jobs/%s", jobID))
	if err != nil {
		c.logger.Error("Gomind GetJobStatus failed", "error", err, "jobID", jobID)
		return nil, err
	}

	var resp APIResponse[JobStatusResponse]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse job status response: %w", err)
	}

	return &resp.Data, nil
}

// GetSystemPrompt fetches the recommended system prompt for LLM integration.
func (c *Client) GetSystemPrompt(ctx context.Context) (string, error) {
	respBody, err := c.get(ctx, "/v1/system-prompt")
	if err != nil {
		c.logger.Error("Gomind GetSystemPrompt failed", "error", err)
		return "", err
	}

	var resp APIResponse[SystemPromptResponse]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("failed to parse system prompt response: %w", err)
	}

	return resp.Data.Prompt, nil
}
