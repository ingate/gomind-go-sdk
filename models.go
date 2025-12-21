package gomind

// APIResponse is the generic response wrapper from Gomind API
type APIResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

// Entity represents an entity in the knowledge graph
type Entity struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	EntityType string `json:"entity_type,omitempty"`
}

// Fact represents a single fact (triplet) in the knowledge graph
type Fact struct {
	ID        string  `json:"id,omitempty"`
	Subject   *Entity `json:"subject,omitempty"`
	Predicate string  `json:"predicate"`
	Object    *Entity `json:"object,omitempty"`
	Value     string  `json:"value,omitempty"`
	Context   string  `json:"context,omitempty"`
	Source    string  `json:"source,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
}

// RememberRequest is the request body for the remember endpoint
type RememberRequest struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
	Context   string `json:"context,omitempty"`
}

// RememberResponse is the response from the remember endpoint
type RememberResponse struct {
	ID        string  `json:"id"`
	Subject   *Entity `json:"subject"`
	Predicate string  `json:"predicate"`
	Object    *Entity `json:"object"`
	CreatedAt string  `json:"created_at"`
}

// RememberManyRequest is the request body for the remember_many endpoint
type RememberManyRequest struct {
	Facts  []RememberRequest `json:"facts"`
	Source string            `json:"source,omitempty"`
}

// RecallRequest is the request body for the recall endpoint
type RecallRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// RecallResponse is the response from the recall endpoint
type RecallResponse struct {
	Facts []Fact `json:"facts"`
}

// RecallConnectionsRequest is the request body for the recall_connections endpoint
type RecallConnectionsRequest struct {
	Entity string `json:"entity"`
	Depth  int    `json:"depth,omitempty"`
}

// ForgetRequest is the request body for the forget endpoint
type ForgetRequest struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
}

// ForgetEntityRequest is the request body for the forget_entity endpoint
type ForgetEntityRequest struct {
	Entity string `json:"entity"`
}

// FeedMessage represents a message in a conversation for the feed endpoint
type FeedMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// FeedRequest is the request body for the feed endpoint
type FeedRequest struct {
	Content  string        `json:"content,omitempty"`
	Messages []FeedMessage `json:"messages,omitempty"`
	Source   string        `json:"source,omitempty"`
}

// FeedResponse is the response from the feed endpoint (sync mode)
type FeedResponse struct {
	Status          string `json:"status"`
	FactsExtracted  int    `json:"facts_extracted,omitempty"`
	EntitiesCreated int    `json:"entities_created,omitempty"`
	Facts           []Fact `json:"facts,omitempty"`
	JobID           string `json:"job_id,omitempty"`
}

// JobStatusResponse is the response from the jobs endpoint
type JobStatusResponse struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	FactsExtracted  int    `json:"facts_extracted,omitempty"`
	EntitiesCreated int    `json:"entities_created,omitempty"`
	Facts           []Fact `json:"facts,omitempty"`
}

// SystemPromptResponse is the response from the system-prompt endpoint
type SystemPromptResponse struct {
	Prompt string `json:"prompt"`
}
