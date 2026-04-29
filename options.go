package gomind

import (
	"net/http"
	"time"
)

// Option is a functional option for configuring the Client.
type Option func(*Client)

// Logger interface for optional logging.
type Logger interface {
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

// noopLogger is a no-op logger implementation.
type noopLogger struct{}

func (n *noopLogger) Info(msg string, keysAndValues ...any)  {}
func (n *noopLogger) Error(msg string, keysAndValues ...any) {}

// WithBaseURL sets the API base URL. If not set, defaults to https://api.gominddb.com.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL != "" {
			c.baseURL = baseURL
		}
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithCollection sets a default collection code applied to every memory
// operation when the per-request Collection field is nil. Reserved
// aliases ("default", "none", "null", "nil", "undefined") and the empty
// string all map to the default bucket server-side.
//
// Per-request override semantics (plan §9 body precedence):
//
//   - Collection == nil               → fall back to this client default.
//   - Collection == &"" (DefaultBucket) → force the server default bucket
//     even if WithCollection set a non-empty default. Use this to
//     escape the configured scope for a single request.
//   - Collection == &"code"           → scope to that collection.
//
// The wire field is always omitted when Collection is nil (no client
// default) so server-side query/header scope can still apply.
func WithCollection(code string) Option {
	return func(c *Client) {
		c.collection = code
	}
}

// CollectionScope wraps a non-nil collection code as a *string so it can
// be assigned to RememberRequest.Collection (and the other request
// types). Equivalent to taking the address of a local variable.
func CollectionScope(code string) *string {
	return &code
}

// DefaultBucket returns a *string that forces the server default
// bucket for a single request, overriding any client-level
// WithCollection default. The wire body carries `"collection": ""`
// which the server interprets as an explicit selection of the
// default bucket (plan §9 body precedence).
func DefaultBucket() *string {
	s := ""
	return &s
}
