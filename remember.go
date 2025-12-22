package gomind

import (
	"context"
	"encoding/json"
	"fmt"
)

// Remember stores a single fact in the knowledge graph.
func (c *Client) Remember(ctx context.Context, subject, predicate, object string, context_ string) (*RememberResponse, error) {
	req := RememberRequest{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
		Context:   context_,
	}

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
		"subject", subject,
		"predicate", predicate,
		"object", object,
	)

	return &resp.Result, nil
}

// RememberMany stores multiple facts in the knowledge graph at once.
func (c *Client) RememberMany(ctx context.Context, facts []RememberRequest, source string) error {
	req := RememberManyRequest{
		Facts:  facts,
		Source: source,
	}

	_, err := c.post(ctx, "/v1/remember_many", req)
	if err != nil {
		c.logger.Error("Gomind RememberMany failed", "error", err, "factCount", len(facts))
		return err
	}

	c.logger.Info("Gomind RememberMany success", "factCount", len(facts))
	return nil
}
