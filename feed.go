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
