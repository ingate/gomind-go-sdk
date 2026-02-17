// Package tools provides LLM tool definitions for Gomind.
// It supports multiple providers through converter functions.
package tools

// ParamType represents a JSON Schema type.
type ParamType string

const (
	TypeString  ParamType = "string"
	TypeInteger ParamType = "integer"
	TypeArray   ParamType = "array"
	TypeObject  ParamType = "object"
)

// Param defines a tool parameter.
type Param struct {
	Name        string
	Type        ParamType
	Description string
	Required    bool
	Items       *Param            // For array types
	Properties  map[string]*Param // For nested objects
}

// Definition defines a tool that can be converted to any provider format.
type Definition struct {
	Name        string
	Description string
	Parameters  []*Param
}

// Definitions returns all Gomind tool definitions.
func Definitions() []Definition {
	return []Definition{
		rememberDef(),
		rememberManyDef(),
		recallDef(),
		recallConnectionsDef(),
		feedDef(),
		forgetDef(),
		forgetEntityDef(),
		mindDef(),
	}
}

func rememberDef() Definition {
	return Definition{
		Name:        "remember",
		Description: "Store a fact or relationship in memory",
		Parameters: []*Param{
			{Name: "subject", Type: TypeString, Description: "The entity the fact is about", Required: true},
			{Name: "predicate", Type: TypeString, Description: "The relationship type (e.g., works_at, likes, has_skill)", Required: true},
			{Name: "object", Type: TypeString, Description: "The related entity or value", Required: true},
			{Name: "context", Type: TypeString, Description: "Optional context or source for the fact"},
		},
	}
}

func rememberManyDef() Definition {
	return Definition{
		Name:        "remember_many",
		Description: "Store multiple facts or relationships in memory at once",
		Parameters: []*Param{
			{
				Name:        "facts",
				Type:        TypeArray,
				Description: "Array of facts to store",
				Required:    true,
				Items: &Param{
					Type: TypeObject,
					Properties: map[string]*Param{
						"subject":   {Name: "subject", Type: TypeString, Description: "The entity the fact is about", Required: true},
						"predicate": {Name: "predicate", Type: TypeString, Description: "The relationship type", Required: true},
						"object":    {Name: "object", Type: TypeString, Description: "The related entity or value", Required: true},
					},
				},
			},
			{Name: "source", Type: TypeString, Description: "Source identifier for traceability"},
		},
	}
}

func recallDef() Definition {
	return Definition{
		Name:        "recall",
		Description: "Search memory for facts about an entity or topic",
		Parameters: []*Param{
			{Name: "query", Type: TypeString, Description: "Entity or topic to search for", Required: true},
			{Name: "predicate", Type: TypeString, Description: "Filter by a single relationship type (e.g. works_at)"},
			{Name: "predicates", Type: TypeArray, Description: "Filter by multiple relationship types (e.g. [question_text, answer_text]). Results match any listed predicate.", Items: &Param{Type: TypeString}},
			{Name: "limit", Type: TypeInteger, Description: "Maximum number of results to return (default: 10)"},
		},
	}
}

func recallConnectionsDef() Definition {
	return Definition{
		Name:        "recall_connections",
		Description: "Get all entities connected to a specific entity in the knowledge graph",
		Parameters: []*Param{
			{Name: "entity", Type: TypeString, Description: "The entity to find connections for", Required: true},
			{Name: "depth", Type: TypeInteger, Description: "How many levels of connections to traverse (default: 2)"},
		},
	}
}

func feedDef() Definition {
	return Definition{
		Name:        "feed",
		Description: "Ingest raw text content (e.g., meeting notes, documents) and automatically extract structured facts",
		Parameters: []*Param{
			{Name: "content", Type: TypeString, Description: "Raw text content to extract facts from", Required: true},
			{Name: "source", Type: TypeString, Description: "Source identifier for traceability"},
		},
	}
}

func forgetDef() Definition {
	return Definition{
		Name:        "forget",
		Description: "Remove a specific fact from memory",
		Parameters: []*Param{
			{Name: "subject", Type: TypeString, Description: "The entity the fact is about", Required: true},
			{Name: "predicate", Type: TypeString, Description: "The relationship type", Required: true},
			{Name: "object", Type: TypeString, Description: "The related entity or value", Required: true},
		},
	}
}

func forgetEntityDef() Definition {
	return Definition{
		Name:        "forget_entity",
		Description: "Remove all facts about a specific entity from memory",
		Parameters: []*Param{
			{Name: "entity", Type: TypeString, Description: "The entity to remove all facts about", Required: true},
		},
	}
}

func mindDef() Definition {
	return Definition{
		Name:        "mind",
		Description: "Make a natural language request that is fulfilled by an internal LLM agent using the knowledge graph",
		Parameters: []*Param{
			{Name: "prompt", Type: TypeString, Description: "Natural language request with optional {{placeholder}} syntax for variable interpolation", Required: true},
			{
				Name:        "context",
				Type:        TypeObject,
				Description: "Key-value pairs to interpolate into the prompt placeholders",
			},
			{
				Name:        "output_schema",
				Type:        TypeObject,
				Description: "Expected response structure with field definitions (type, items, description)",
				Required:    true,
			},
		},
	}
}
