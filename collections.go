package gomind

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Collection is a server-side partition of an org's knowledge graph.
// Each collection has a human-friendly Code (immutable after creation)
// and a display Name / Description. CreatedAt/UpdatedAt are unix seconds.
// FactCount is populated only when GetCollection is called with
// ?include=fact_count.
type Collection struct {
	ID          string `json:"id"`
	OrgID       string `json:"org_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	FactCount   *int   `json:"fact_count,omitempty"`
}

// DeleteSummary is returned by DeleteCollection and reports the cross-DB
// cleanup totals (Postgres row + Dgraph facts/entities).
type DeleteSummary struct {
	DeletedFacts    int `json:"deleted_facts"`
	DeletedEntities int `json:"deleted_entities"`
}

// MoveSummary is returned by MoveFactsToCollection.
type MoveSummary struct {
	MovedFacts int `json:"moved_facts"`
}

// createCollectionRequest is the wire body for POST /collections.
type createCollectionRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// updateCollectionRequest is the wire body for PATCH /collections/:id.
// Code is intentionally absent — collection codes are immutable.
type updateCollectionRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// moveFactsRequest is the wire body for POST /collections/:id/move.
type moveFactsRequest struct {
	FactIDs []string `json:"fact_ids"`
}

// ListCollections returns all collections for the given org.
func (c *Client) ListCollections(ctx context.Context, orgID string) ([]Collection, error) {
	if strings.TrimSpace(orgID) == "" {
		return nil, fmt.Errorf("org id is required")
	}

	endpoint := fmt.Sprintf("/v1/orgs/%s/collections/", url.PathEscape(orgID))
	respBody, err := c.get(ctx, endpoint)
	if err != nil {
		c.logger.Error("Gomind ListCollections failed", "error", err, "orgID", orgID)
		return nil, err
	}

	var resp APIResponse[[]Collection]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse list collections response: %w", err)
	}
	return resp.Result, nil
}

// CreateCollection creates a new collection. The code must match
// ^[a-z0-9][a-z0-9-]{0,62}$ and not be a reserved alias
// (default/none/null/nil/undefined).
func (c *Client) CreateCollection(ctx context.Context, orgID, code, name, description string) (*Collection, error) {
	if strings.TrimSpace(orgID) == "" {
		return nil, fmt.Errorf("org id is required")
	}

	body := createCollectionRequest{
		Code:        code,
		Name:        name,
		Description: description,
	}

	endpoint := fmt.Sprintf("/v1/orgs/%s/collections/", url.PathEscape(orgID))
	respBody, err := c.post(ctx, endpoint, body)
	if err != nil {
		c.logger.Error("Gomind CreateCollection failed", "error", err, "orgID", orgID, "code", code)
		return nil, err
	}

	var resp APIResponse[Collection]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse create collection response: %w", err)
	}
	return &resp.Result, nil
}

// GetCollection fetches a single collection by its internal ID.
// Returns a Collection with FactCount populated.
func (c *Client) GetCollection(ctx context.Context, orgID, id string) (*Collection, error) {
	if strings.TrimSpace(orgID) == "" {
		return nil, fmt.Errorf("org id is required")
	}
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("collection id is required")
	}

	endpoint := fmt.Sprintf("/v1/orgs/%s/collections/%s?include=fact_count",
		url.PathEscape(orgID), url.PathEscape(id))
	respBody, err := c.get(ctx, endpoint)
	if err != nil {
		c.logger.Error("Gomind GetCollection failed", "error", err, "orgID", orgID, "id", id)
		return nil, err
	}

	var resp APIResponse[Collection]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse get collection response: %w", err)
	}
	return &resp.Result, nil
}

// UpdateCollection patches the name and/or description. Pass nil to leave
// a field unchanged. The code is immutable and cannot be updated.
func (c *Client) UpdateCollection(ctx context.Context, orgID, id string, name, description *string) (*Collection, error) {
	if strings.TrimSpace(orgID) == "" {
		return nil, fmt.Errorf("org id is required")
	}
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("collection id is required")
	}

	body := updateCollectionRequest{
		Name:        name,
		Description: description,
	}

	endpoint := fmt.Sprintf("/v1/orgs/%s/collections/%s",
		url.PathEscape(orgID), url.PathEscape(id))
	respBody, err := c.patch(ctx, endpoint, body)
	if err != nil {
		c.logger.Error("Gomind UpdateCollection failed", "error", err, "orgID", orgID, "id", id)
		return nil, err
	}

	var resp APIResponse[Collection]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse update collection response: %w", err)
	}
	return &resp.Result, nil
}

// DeleteCollection hard-deletes a collection and all of its facts and
// entities. The returned summary reports the totals removed from the
// knowledge graph.
func (c *Client) DeleteCollection(ctx context.Context, orgID, id string) (*DeleteSummary, error) {
	if strings.TrimSpace(orgID) == "" {
		return nil, fmt.Errorf("org id is required")
	}
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("collection id is required")
	}

	endpoint := fmt.Sprintf("/v1/orgs/%s/collections/%s",
		url.PathEscape(orgID), url.PathEscape(id))
	respBody, err := c.delete(ctx, endpoint)
	if err != nil {
		c.logger.Error("Gomind DeleteCollection failed", "error", err, "orgID", orgID, "id", id)
		return nil, err
	}

	var resp APIResponse[DeleteSummary]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse delete collection response: %w", err)
	}
	return &resp.Result, nil
}

// MoveFactsToCollection moves the given facts into the target collection.
// The target collection ID is the *internal* collection ID, not the code.
// Fails with 409 if any of the facts' referenced entities are shared with
// facts outside the move set (v1 rejects rather than clones).
func (c *Client) MoveFactsToCollection(ctx context.Context, orgID, targetID string, factIDs []string) (*MoveSummary, error) {
	if strings.TrimSpace(orgID) == "" {
		return nil, fmt.Errorf("org id is required")
	}
	if strings.TrimSpace(targetID) == "" {
		return nil, fmt.Errorf("target collection id is required")
	}
	if len(factIDs) == 0 {
		return nil, fmt.Errorf("at least one fact id is required")
	}

	body := moveFactsRequest{FactIDs: factIDs}

	endpoint := fmt.Sprintf("/v1/orgs/%s/collections/%s/move",
		url.PathEscape(orgID), url.PathEscape(targetID))
	respBody, err := c.post(ctx, endpoint, body)
	if err != nil {
		c.logger.Error("Gomind MoveFactsToCollection failed",
			"error", err, "orgID", orgID, "targetID", targetID, "count", len(factIDs))
		return nil, err
	}

	var resp APIResponse[MoveSummary]
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse move facts response: %w", err)
	}
	return &resp.Result, nil
}
