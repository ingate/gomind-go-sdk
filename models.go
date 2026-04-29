package gomind

// APIResponse is the generic response wrapper from Gomind API
type APIResponse[T any] struct {
	Status string `json:"status"`
	Result T      `json:"result"`
}

// Entity represents an entity in the knowledge graph
type Entity struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	EntityType string `json:"entity_type,omitempty"`
}

// Fact represents a single fact (triplet) in the knowledge graph
// This matches the API's FactOutput format
type Fact struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object,omitempty"`
	Value     string `json:"value,omitempty"`
	Context   string `json:"context,omitempty"`
	Source    string `json:"source,omitempty"`
}

// RememberRequest is the request body for the remember endpoint.
//
// Collection is a *string so the SDK can express three intents
// distinctly (plan §9 body precedence):
//
//   - nil       → fall back to the client default set via WithCollection.
//   - &""       → explicitly select the server default bucket, even if a
//     non-empty client default is configured. Use the
//     DefaultBucket() helper to construct this value.
//   - &"code"   → scope to the named collection. Use CollectionScope("code").
type RememberRequest struct {
	Subject    string  `json:"subject"`
	Predicate  string  `json:"predicate"`
	Object     string  `json:"object"`
	Context    string  `json:"context,omitempty"`
	Normalize  bool    `json:"normalize,omitempty"`
	Collection *string `json:"collection,omitempty"`
}

// RememberResponse is the response from the remember endpoint
// Matches API's FactOutput format
type RememberResponse struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object,omitempty"`
	Value     string `json:"value,omitempty"`
	Context   string `json:"context,omitempty"`
	Source    string `json:"source,omitempty"`
}

// RememberManyRequest is the request body for the remember_many endpoint.
// See RememberRequest.Collection for the *string semantics.
type RememberManyRequest struct {
	Facts      []RememberRequest `json:"facts"`
	Source     string            `json:"source,omitempty"`
	Collection *string           `json:"collection,omitempty"`
}

// RecallRequest is the request body for the recall endpoint.
// See RememberRequest.Collection for the *string semantics.
type RecallRequest struct {
	Query      string   `json:"query,omitempty"`
	Predicate  string   `json:"predicate,omitempty"`
	Predicates []string `json:"predicates,omitempty"`
	EntityType string   `json:"entity_type,omitempty"`
	RelatedTo  string   `json:"related_to,omitempty"`
	Depth      int      `json:"depth,omitempty"`
	FuzzyMatch bool     `json:"fuzzy_match,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Collection *string  `json:"collection,omitempty"`
}

// RecallResponse is the response from the recall endpoint
type RecallResponse struct {
	Facts       []Fact   `json:"facts"`
	Count       int      `json:"count"`
	SearchMode  string   `json:"search_mode,omitempty"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// RecallConnectionsRequest is the request body for the recall_connections endpoint.
// See RememberRequest.Collection for the *string semantics.
type RecallConnectionsRequest struct {
	Entity     string  `json:"entity"`
	Depth      int     `json:"depth,omitempty"`
	Collection *string `json:"collection,omitempty"`
}

// ForgetRequest is the request body for the forget endpoint.
// See RememberRequest.Collection for the *string semantics.
type ForgetRequest struct {
	Subject    string  `json:"subject"`
	Predicate  string  `json:"predicate"`
	Object     string  `json:"object"`
	Collection *string `json:"collection,omitempty"`
}

// ForgetEntityRequest is the request body for the forget_entity endpoint.
// See RememberRequest.Collection for the *string semantics.
type ForgetEntityRequest struct {
	Entity     string  `json:"entity"`
	Collection *string `json:"collection,omitempty"`
}

// FeedMessage represents a message in a conversation for the feed endpoint
type FeedMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// FeedRequest is the request body for the feed endpoint.
// See RememberRequest.Collection for the *string semantics.
type FeedRequest struct {
	Content    string        `json:"content,omitempty"`
	Messages   []FeedMessage `json:"messages,omitempty"`
	Source     string        `json:"source,omitempty"`
	Collection *string       `json:"collection,omitempty"`
}

// FeedResponse is the response from the feed endpoint (sync mode)
type FeedResponse struct {
	Status          string `json:"status"`
	FactsExtracted  int    `json:"facts_extracted,omitempty"`
	FactsCreated    int    `json:"facts_created,omitempty"`
	EntitiesCreated int    `json:"entities_created,omitempty"`
	Facts           []Fact `json:"facts,omitempty"`
	JobID           string `json:"job_id,omitempty"`
}

// SystemPromptResponse is the response from the system-prompt endpoint
type SystemPromptResponse struct {
	Prompt string `json:"prompt"`
}
