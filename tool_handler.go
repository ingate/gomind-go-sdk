package gomind

import (
	"context"
	"encoding/json"
	"fmt"
)

// HandleToolCall executes a Gomind tool call and returns the result.
// It routes the tool call to the appropriate API method based on the tool name.
func (c *Client) HandleToolCall(ctx context.Context, name string, arguments string) (any, error) {
	switch name {
	case "remember":
		var req RememberRequest
		if err := json.Unmarshal([]byte(arguments), &req); err != nil {
			return nil, fmt.Errorf("failed to parse remember arguments: %w", err)
		}
		return c.Remember(ctx, req.Subject, req.Predicate, req.Object, req.Context)

	case "remember_many":
		var req RememberManyRequest
		if err := json.Unmarshal([]byte(arguments), &req); err != nil {
			return nil, fmt.Errorf("failed to parse remember_many arguments: %w", err)
		}
		if err := c.RememberMany(ctx, req.Facts, req.Source); err != nil {
			return nil, err
		}
		return map[string]string{"status": "OK"}, nil

	case "recall":
		var req RecallRequest
		if err := json.Unmarshal([]byte(arguments), &req); err != nil {
			return nil, fmt.Errorf("failed to parse recall arguments: %w", err)
		}
		return c.Recall(ctx, req.Query, req.Limit)

	case "recall_connections":
		var req RecallConnectionsRequest
		if err := json.Unmarshal([]byte(arguments), &req); err != nil {
			return nil, fmt.Errorf("failed to parse recall_connections arguments: %w", err)
		}
		return c.RecallConnections(ctx, req.Entity, req.Depth)

	case "feed":
		var req FeedRequest
		if err := json.Unmarshal([]byte(arguments), &req); err != nil {
			return nil, fmt.Errorf("failed to parse feed arguments: %w", err)
		}
		return c.Feed(ctx, req.Content, req.Source)

	case "forget":
		var req ForgetRequest
		if err := json.Unmarshal([]byte(arguments), &req); err != nil {
			return nil, fmt.Errorf("failed to parse forget arguments: %w", err)
		}
		if err := c.Forget(ctx, req.Subject, req.Predicate, req.Object); err != nil {
			return nil, err
		}
		return map[string]string{"status": "OK"}, nil

	case "forget_entity":
		var req ForgetEntityRequest
		if err := json.Unmarshal([]byte(arguments), &req); err != nil {
			return nil, fmt.Errorf("failed to parse forget_entity arguments: %w", err)
		}
		if err := c.ForgetEntity(ctx, req.Entity); err != nil {
			return nil, err
		}
		return map[string]string{"status": "OK"}, nil

	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// HandleToolCallJSON is a convenience method that returns the tool call result as a JSON string.
// This is useful for directly returning the result to the LLM.
func (c *Client) HandleToolCallJSON(ctx context.Context, name string, arguments string) (string, error) {
	result, err := c.HandleToolCall(ctx, name, arguments)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(jsonBytes), nil
}
