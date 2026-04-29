package gomind

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the Gomind API client.
type Client struct {
	baseURL    string
	apiKey     string
	collection string
	httpClient *http.Client
	logger     Logger
}

// NewClient creates a new Gomind client.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	c := &Client{
		baseURL: "https://api.gominddb.com",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: &noopLogger{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// doRequest performs an HTTP request to the Gomind API.
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body any) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// post performs a POST request to the Gomind API.
func (c *Client) post(ctx context.Context, endpoint string, body any) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPost, endpoint, body)
}

// resolveCollection applies plan §9 body precedence to merge the
// per-request override with the client default set via WithCollection.
//
//   - reqCol != nil               → body wins, even when *reqCol == "".
//     This is what lets callers escape the
//     client default and force the server's
//     default bucket on a per-request basis
//     (DefaultBucket()).
//   - reqCol == nil, c.collection == "" → leave the field absent so
//     server-side query/header scope can apply.
//   - reqCol == nil, c.collection != "" → inject the client default.
//
// Pre-fix this returned a plain string and treated "" the same as
// "absent", which made the explicit-empty override unreachable from
// the SDK (bug 16).
func (c *Client) resolveCollection(reqCol *string) *string {
	if reqCol != nil {
		return reqCol
	}
	if c.collection == "" {
		return nil
	}
	return &c.collection
}

// get performs a GET request to the Gomind API.
func (c *Client) get(ctx context.Context, endpoint string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, endpoint, nil)
}

// patch performs a PATCH request to the Gomind API.
func (c *Client) patch(ctx context.Context, endpoint string, body any) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPatch, endpoint, body)
}

// delete performs a DELETE request to the Gomind API.
func (c *Client) delete(ctx context.Context, endpoint string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodDelete, endpoint, nil)
}
