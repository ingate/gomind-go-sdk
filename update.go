package gomind

import (
	"context"
	"encoding/json"
	"fmt"
)

// Update rewrites a single existing fact in place.
func (c *Client) Update(ctx context.Context, req UpdateRequest) (*RememberResponse, error) {
	return c.UpdateWithOptions(ctx, req)
}

// UpdateWithOptions rewrites a fact with full control over the request payload,
// including the optional Collection field.
func (c *Client) UpdateWithOptions(ctx context.Context, req UpdateRequest) (*RememberResponse, error) {
	req.Collection = c.resolveCollection(req.Collection)
	respBody, err := c.post(ctx, "/v1/update", req)
	if err != nil {
		c.logger.Error("Gomind Update failed",
			"error", err,
			"subject", req.Subject,
			"predicate", req.Predicate,
			"object", req.Object,
		)
		return nil, err
	}

	var resp APIResponse[RememberResponse]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	c.logger.Info("Gomind Update success",
		"subject", req.Subject,
		"predicate", req.Predicate,
		"object", req.Object,
	)

	return &resp.Result, nil
}
