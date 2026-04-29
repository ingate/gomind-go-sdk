package gomind

import (
	"context"
	"encoding/json"
	"fmt"
)

// Remember stores a single fact in the knowledge graph.
// For additional options like normalization, use RememberWithOptions.
func (c *Client) Remember(ctx context.Context, subject, predicate, object string, context_ string) (*RememberResponse, error) {
	return c.RememberWithOptions(ctx, RememberRequest{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
		Context:   context_,
	})
}

// RememberWithOptions stores a single fact with full control over request parameters.
// Use the Normalize field to enable LLM-based normalization of abbreviations.
func (c *Client) RememberWithOptions(ctx context.Context, req RememberRequest) (*RememberResponse, error) {
	req.Collection = c.resolveCollection(req.Collection)
	respBody, err := c.post(ctx, "/v1/remember", req)
	if err != nil {
		c.logger.Error("Gomind Remember failed", "error", err)
		return nil, err
	}

	var resp APIResponse[RememberResponse]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse remember response: %w", err)
	}

	c.logger.Info("Gomind Remember success",
		"subject", req.Subject,
		"predicate", req.Predicate,
		"object", req.Object,
	)

	return &resp.Result, nil
}

// RememberMany stores multiple facts in the knowledge graph at once.
func (c *Client) RememberMany(ctx context.Context, facts []RememberRequest, source string) error {
	return c.RememberManyWithOptions(ctx, RememberManyRequest{
		Facts:  facts,
		Source: source,
	})
}

// RememberManyWithOptions stores multiple facts with full control over the
// request payload, including the optional Collection field.
func (c *Client) RememberManyWithOptions(ctx context.Context, req RememberManyRequest) error {
	req.Collection = c.resolveCollection(req.Collection)

	_, err := c.post(ctx, "/v1/remember_many", req)
	if err != nil {
		c.logger.Error("Gomind RememberMany failed", "error", err, "factCount", len(req.Facts))
		return err
	}

	c.logger.Info("Gomind RememberMany success", "factCount", len(req.Facts))
	return nil
}
