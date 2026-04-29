package gomind

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCreateCollectionRequestJSON verifies the JSON shape of the body
// sent to POST /v1/orgs/:org_id/collections matches the server's
// CreateCollectionRequest serializer (code/name/description, description
// omitted when empty).
func TestCreateCollectionRequestJSON(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		collName    string
		description string
		wantKeys    []string
		wantNoKeys  []string
	}{
		{
			name:        "full body",
			code:        "prod",
			collName:    "Production",
			description: "Production collection",
			wantKeys:    []string{`"code":"prod"`, `"name":"Production"`, `"description":"Production collection"`},
		},
		{
			name:       "omit empty description",
			code:       "staging",
			collName:   "Staging",
			wantKeys:   []string{`"code":"staging"`, `"name":"Staging"`},
			wantNoKeys: []string{`"description"`},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := createCollectionRequest{
				Code:        tc.code,
				Name:        tc.collName,
				Description: tc.description,
			}
			raw, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			s := string(raw)
			for _, k := range tc.wantKeys {
				if !strings.Contains(s, k) {
					t.Errorf("expected %q in body, got %s", k, s)
				}
			}
			for _, k := range tc.wantNoKeys {
				if strings.Contains(s, k) {
					t.Errorf("did not expect %q in body, got %s", k, s)
				}
			}
		})
	}
}

// TestUpdateCollectionRequestJSON verifies nil fields are omitted so
// callers can patch name OR description independently.
func TestUpdateCollectionRequestJSON(t *testing.T) {
	newName := "Renamed"
	body := updateCollectionRequest{Name: &newName}

	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(raw)

	if !strings.Contains(s, `"name":"Renamed"`) {
		t.Errorf("expected name in body, got %s", s)
	}
	if strings.Contains(s, `"description"`) {
		t.Errorf("expected nil description to be omitted, got %s", s)
	}
}

// TestMoveFactsRequestJSON verifies fact_ids is the wire key (not factIds).
func TestMoveFactsRequestJSON(t *testing.T) {
	body := moveFactsRequest{FactIDs: []string{"0x01", "0x02"}}
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(raw)
	if !strings.Contains(s, `"fact_ids":["0x01","0x02"]`) {
		t.Errorf("expected fact_ids in body, got %s", s)
	}
}

// TestCreateCollectionEndToEnd spins up a test server and verifies the
// client posts to the correct URL with the expected body shape.
func TestCreateCollectionEndToEnd(t *testing.T) {
	var capturedPath string
	var capturedBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"OK","result":{"id":"col_123","org_id":"org_abc","code":"prod","name":"Production","description":"","created_at":1700000000,"updated_at":0}}`))
	}))
	defer srv.Close()

	client, err := NewClient("test-key", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	col, err := client.CreateCollection(context.Background(), "org_abc", "prod", "Production", "")
	if err != nil {
		t.Fatalf("CreateCollection: %v", err)
	}

	if capturedPath != "/v1/orgs/org_abc/collections/" {
		t.Errorf("unexpected path %q", capturedPath)
	}
	if !strings.Contains(string(capturedBody), `"code":"prod"`) {
		t.Errorf("body missing code, got %s", capturedBody)
	}
	if col.ID != "col_123" || col.Code != "prod" {
		t.Errorf("unexpected parsed collection: %+v", col)
	}
}

// TestWithCollectionInjection verifies the client default collection
// flows into request bodies that leave Collection nil, and that an
// explicit per-request override (including DefaultBucket()) wins.
//
// The DefaultBucket() case pins the bug-16 fix: a *string pointer to
// the empty string must travel on the wire as `"collection":""` so
// the server interprets it as an explicit selection of the default
// bucket. Pre-fix the SDK collapsed empty strings into "use client
// default", making the override unreachable.
func TestWithCollectionInjection(t *testing.T) {
	var capturedBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"OK","result":{"subject":"s","predicate":"p","object":"o"}}`))
	}))
	defer srv.Close()

	client, err := NewClient("test-key", WithBaseURL(srv.URL), WithCollection("prod"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// nil per-request collection inherits the client default.
	if _, err := client.Remember(context.Background(), "s", "p", "o", ""); err != nil {
		t.Fatalf("Remember: %v", err)
	}
	if !strings.Contains(string(capturedBody), `"collection":"prod"`) {
		t.Errorf("expected collection=prod in body, got %s", capturedBody)
	}

	// Per-request override wins over client default.
	if _, err := client.RememberWithOptions(context.Background(), RememberRequest{
		Subject: "s", Predicate: "p", Object: "o", Collection: CollectionScope("staging"),
	}); err != nil {
		t.Fatalf("RememberWithOptions: %v", err)
	}
	if !strings.Contains(string(capturedBody), `"collection":"staging"`) {
		t.Errorf("expected per-request collection=staging to win, got %s", capturedBody)
	}

	// DefaultBucket() forces the server default bucket even when the
	// client has WithCollection("prod") configured.
	if _, err := client.RememberWithOptions(context.Background(), RememberRequest{
		Subject: "s", Predicate: "p", Object: "o", Collection: DefaultBucket(),
	}); err != nil {
		t.Fatalf("RememberWithOptions(DefaultBucket): %v", err)
	}
	if !strings.Contains(string(capturedBody), `"collection":""`) {
		t.Errorf("expected explicit collection=\"\" on the wire, got %s", capturedBody)
	}
	if strings.Contains(string(capturedBody), `"collection":"prod"`) {
		t.Errorf("DefaultBucket() must not fall back to client default, got %s", capturedBody)
	}
}

// TestResolveCollection_NoClientDefault confirms that with no client
// default and a nil per-request override, the field is omitted from
// the wire body (server-side query/header scope can still apply).
func TestResolveCollection_NoClientDefault(t *testing.T) {
	var capturedBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"OK","result":{"subject":"s","predicate":"p","object":"o"}}`))
	}))
	defer srv.Close()

	client, err := NewClient("test-key", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	if _, err := client.Remember(context.Background(), "s", "p", "o", ""); err != nil {
		t.Fatalf("Remember: %v", err)
	}
	if strings.Contains(string(capturedBody), `"collection"`) {
		t.Errorf("expected collection field to be omitted, got %s", capturedBody)
	}
}

// TestResolveCollection_ExplicitEmptyAcrossEndpoints sweeps all
// memory operations to confirm DefaultBucket() forces the wire field
// to "" on each one (bug 16 verification requirement: SDK must allow
// explicit empty for every endpoint).
func TestResolveCollection_ExplicitEmptyAcrossEndpoints(t *testing.T) {
	var capturedBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		// Generic OK response, ignored by tests below.
		_, _ = w.Write([]byte(`{"status":"OK","result":{}}`))
	}))
	defer srv.Close()

	client, err := NewClient("test-key", WithBaseURL(srv.URL), WithCollection("prod"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	ctx := context.Background()

	// Each invocation overwrites capturedBody, so assert per-call.
	check := func(t *testing.T, label string) {
		t.Helper()
		if !strings.Contains(string(capturedBody), `"collection":""`) {
			t.Errorf("%s: expected collection=\"\" on the wire, got %s", label, capturedBody)
		}
	}

	if _, err := client.RememberWithOptions(ctx, RememberRequest{
		Subject: "s", Predicate: "p", Object: "o", Collection: DefaultBucket(),
	}); err != nil {
		t.Fatalf("Remember: %v", err)
	}
	check(t, "remember")

	if err := client.RememberManyWithOptions(ctx, RememberManyRequest{
		Facts:      []RememberRequest{{Subject: "s", Predicate: "p", Object: "o"}},
		Collection: DefaultBucket(),
	}); err != nil {
		t.Fatalf("RememberMany: %v", err)
	}
	check(t, "remember_many")

	if _, err := client.RecallWithOptions(ctx, RecallRequest{
		Query: "x", Collection: DefaultBucket(),
	}); err != nil {
		t.Fatalf("Recall: %v", err)
	}
	check(t, "recall")

	if _, err := client.RecallConnectionsWithOptions(ctx, RecallConnectionsRequest{
		Entity: "e", Collection: DefaultBucket(),
	}); err != nil {
		t.Fatalf("RecallConnections: %v", err)
	}
	check(t, "recall_connections")

	if err := client.ForgetWithOptions(ctx, ForgetRequest{
		Subject: "s", Predicate: "p", Object: "o", Collection: DefaultBucket(),
	}); err != nil {
		t.Fatalf("Forget: %v", err)
	}
	check(t, "forget")

	if err := client.ForgetEntityWithOptions(ctx, ForgetEntityRequest{
		Entity: "e", Collection: DefaultBucket(),
	}); err != nil {
		t.Fatalf("ForgetEntity: %v", err)
	}
	check(t, "forget_entity")

	if _, err := client.FeedWithOptions(ctx, FeedRequest{
		Content: "x", Collection: DefaultBucket(),
	}); err != nil {
		t.Fatalf("Feed: %v", err)
	}
	check(t, "feed")
}
