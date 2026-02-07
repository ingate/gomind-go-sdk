# Gomind Go SDK

Go client library for the Gomind knowledge graph API.

## Installation

```bash
go get github.com/ingate/gomind-go-sdk
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    gomind "github.com/ingate/gomind-go-sdk"
)

func main() {
    // Create a new client (uses https://api.gominddb.com by default)
    client, err := gomind.NewClient("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Remember a fact
    resp, err := client.Remember(ctx, "John", "works_at", "Acme Corp", "employment")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created fact: %s\n", resp.ID)

    // Recall facts
    recallResp, err := client.Recall(ctx, "John", 10)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d facts\n", len(recallResp.Facts))

    // Format facts for LLM context
    contextStr := gomind.FormatFactsAsContext(recallResp.Facts)
    fmt.Println(contextStr)
}
```

## Options

```go
// Custom base URL (default is https://api.gominddb.com)
client, _ := gomind.NewClient(apiKey, gomind.WithBaseURL("https://api.gomind.example.com"))

// Custom HTTP client
client, _ := gomind.NewClient(apiKey,
    gomind.WithBaseURL(baseURL),
    gomind.WithHTTPClient(customHTTPClient),
)

// Custom timeout
client, _ := gomind.NewClient(apiKey,
    gomind.WithBaseURL(baseURL),
    gomind.WithTimeout(60 * time.Second),
)

// Custom logger (implements gomind.Logger interface)
client, _ := gomind.NewClient(apiKey,
    gomind.WithBaseURL(baseURL),
    gomind.WithLogger(myLogger),
)
```

## API Methods

### Memory Operations

- `Remember(ctx, subject, predicate, object, context)` - Store a single fact
- `RememberMany(ctx, facts, source)` - Store multiple facts
- `Recall(ctx, query, limit)` - Search for facts
- `RecallConnections(ctx, entity, depth)` - Get connected entities
- `Forget(ctx, subject, predicate, object)` - Remove a specific fact
- `ForgetEntity(ctx, entity)` - Remove all facts about an entity

### Feed Operations

- `Feed(ctx, content, source)` - Extract facts from content
- `FeedMessages(ctx, messages, source)` - Extract facts from conversation
- `FeedAsync(ctx, content, source)` - Async fact extraction
- `GetJobStatus(ctx, jobID)` - Check async job status

### Utilities

- `GetSystemPrompt(ctx)` - Get recommended LLM system prompt
- `FormatFactsAsContext(facts)` - Format facts for LLM context (TOON format)
- `Encode(v)` - Encode any value to TOON format
- `EncodeTabular(name, rows, fields...)` - Encode tabular data to TOON

## License

MIT
