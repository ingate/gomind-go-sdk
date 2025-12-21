package gomind

import (
	"context"
	"encoding/json"
	"fmt"
)

// Recall searches for facts in the knowledge graph.
func (c *Client) Recall(ctx context.Context, query string, limit int) (*RecallResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	req := RecallRequest{
		Query: query,
		Limit: limit,
	}

	respBody, err := c.post(ctx, "/v1/recall", req)
	if err != nil {
		c.logger.Error("Gomind Recall failed", "error", err, "query", query)
		return nil, err
	}

	var resp APIResponse[RecallResponse]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse recall response: %w", err)
	}

	c.logger.Info("Gomind Recall success",
		"query", query,
		"factsFound", len(resp.Data.Facts),
	)

	return &resp.Data, nil
}

// RecallConnections gets all entities connected to a specific entity.
func (c *Client) RecallConnections(ctx context.Context, entity string, depth int) (*RecallResponse, error) {
	if depth <= 0 {
		depth = 2
	}

	req := RecallConnectionsRequest{
		Entity: entity,
		Depth:  depth,
	}

	respBody, err := c.post(ctx, "/v1/recall_connections", req)
	if err != nil {
		c.logger.Error("Gomind RecallConnections failed", "error", err, "entity", entity)
		return nil, err
	}

	var resp APIResponse[RecallResponse]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse recall_connections response: %w", err)
	}

	c.logger.Info("Gomind RecallConnections success",
		"entity", entity,
		"factsFound", len(resp.Data.Facts),
	)

	return &resp.Data, nil
}

// factRow represents a flattened fact for TOON encoding.
type factRow struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
}

// FormatFactsAsContext formats recalled facts as a context string for the LLM.
func FormatFactsAsContext(facts []Fact) string {
	if len(facts) == 0 {
		return ""
	}

	// Build rows from valid facts
	rows := make([]factRow, 0, len(facts))
	for _, fact := range facts {
		subject := getEntityName(fact.Subject)
		object := getObjectValue(fact)
		if subject != "" && object != "" {
			rows = append(rows, factRow{
				Subject:   subject,
				Predicate: fact.Predicate,
				Object:    object,
			})
		}
	}

	if len(rows) == 0 {
		return ""
	}

	return EncodeTabularAuto("memory", rows)
}

// getEntityName safely extracts the name from an Entity pointer.
func getEntityName(e *Entity) string {
	if e != nil {
		return e.Name
	}
	return ""
}

// getObjectValue gets the object value from a Fact, preferring Object.Name over Value.
func getObjectValue(fact Fact) string {
	if fact.Object != nil && fact.Object.Name != "" {
		return fact.Object.Name
	}
	return fact.Value
}
