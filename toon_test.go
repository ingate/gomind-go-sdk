package gomind

import (
	"testing"
)

func TestEncodeTabular(t *testing.T) {
	tests := []struct {
		name     string
		arrName  string
		rows     []map[string]string
		fields   []string
		expected string
	}{
		{
			name:     "empty rows",
			arrName:  "items",
			rows:     []map[string]string{},
			fields:   []string{"a", "b"},
			expected: "items[0]{a,b}:",
		},
		{
			name:    "single row",
			arrName: "facts",
			rows: []map[string]string{
				{"subject": "Alice", "predicate": "likes", "object": "Bob"},
			},
			fields:   []string{"subject", "predicate", "object"},
			expected: "facts[1]{subject,predicate,object}:\n  Alice,likes,Bob",
		},
		{
			name:    "multiple rows",
			arrName: "memory",
			rows: []map[string]string{
				{"subject": "Alice", "predicate": "likes", "object": "Bob"},
				{"subject": "Bob", "predicate": "works_at", "object": "Acme"},
			},
			fields:   []string{"subject", "predicate", "object"},
			expected: "memory[2]{subject,predicate,object}:\n  Alice,likes,Bob\n  Bob,works_at,Acme",
		},
		{
			name:    "values with commas need quoting",
			arrName: "data",
			rows: []map[string]string{
				{"name": "Smith, John", "role": "admin"},
			},
			fields:   []string{"name", "role"},
			expected: "data[1]{name,role}:\n  \"Smith, John\",admin",
		},
		{
			name:    "values with quotes need escaping",
			arrName: "data",
			rows: []map[string]string{
				{"name": `Say "Hello"`, "value": "test"},
			},
			fields:   []string{"name", "value"},
			expected: "data[1]{name,value}:\n  \"Say \"\"Hello\"\"\",test",
		},
		{
			name:    "empty values",
			arrName: "data",
			rows: []map[string]string{
				{"a": "", "b": "value"},
			},
			fields:   []string{"a", "b"},
			expected: "data[1]{a,b}:\n  ,value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeTabular(tt.arrName, tt.rows, tt.fields...)
			if result != tt.expected {
				t.Errorf("EncodeTabular() =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestEscapeValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"", ""},
		{"hello world", "hello world"},
		{"has,comma", `"has,comma"`},
		{`has"quote`, `"has""quote"`},
		{"has\nnewline", `"has
newline"`},
		{"has:colon", `"has:colon"`},
		{"has[bracket", `"has[bracket"`},
		{"has{brace", `"has{brace"`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeValue(tt.input)
			if result != tt.expected {
				t.Errorf("escapeValue(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"nil", nil, "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Encode(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatFactsAsContext(t *testing.T) {
	tests := []struct {
		name     string
		facts    []Fact
		expected string
	}{
		{
			name:     "empty facts",
			facts:    []Fact{},
			expected: "",
		},
		{
			name: "single fact with entity object",
			facts: []Fact{
				{
					Subject:   &Entity{Name: "Alice"},
					Predicate: "likes",
					Object:    &Entity{Name: "Bob"},
				},
			},
			expected: "memory[1]{subject,predicate,object}:\n  Alice,likes,Bob",
		},
		{
			name: "single fact with value",
			facts: []Fact{
				{
					Subject:   &Entity{Name: "Alice"},
					Predicate: "prefers",
					Value:     "dark mode",
				},
			},
			expected: "memory[1]{subject,predicate,object}:\n  Alice,prefers,dark mode",
		},
		{
			name: "multiple facts",
			facts: []Fact{
				{
					Subject:   &Entity{Name: "Alice"},
					Predicate: "works_at",
					Object:    &Entity{Name: "Acme Corp"},
				},
				{
					Subject:   &Entity{Name: "Bob"},
					Predicate: "likes",
					Value:     "coffee",
				},
			},
			expected: "memory[2]{subject,predicate,object}:\n  Alice,works_at,Acme Corp\n  Bob,likes,coffee",
		},
		{
			name: "fact with nil subject filtered",
			facts: []Fact{
				{
					Subject:   nil,
					Predicate: "likes",
					Object:    &Entity{Name: "Bob"},
				},
			},
			expected: "",
		},
		{
			name: "fact with nil object and empty value filtered",
			facts: []Fact{
				{
					Subject:   &Entity{Name: "Alice"},
					Predicate: "likes",
					Object:    nil,
					Value:     "",
				},
			},
			expected: "",
		},
		{
			name: "fact with special characters",
			facts: []Fact{
				{
					Subject:   &Entity{Name: "John Smith, Jr."},
					Predicate: "title",
					Value:     "CEO",
				},
			},
			expected: "memory[1]{subject,predicate,object}:\n  \"John Smith, Jr.\",title,CEO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFactsAsContext(tt.facts)
			if result != tt.expected {
				t.Errorf("FormatFactsAsContext() =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestEncodeTabularAuto(t *testing.T) {
	type User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email,omitempty"`
		IsActive bool   `json:"active"`
		internal string // unexported, should be ignored
	}

	tests := []struct {
		name     string
		arrName  string
		items    []User
		fields   []string
		expected string
	}{
		{
			name:     "auto detect all fields",
			arrName:  "users",
			items:    []User{{ID: 1, Name: "Alice", Email: "alice@example.com", IsActive: true}},
			fields:   nil,
			expected: "users[1]{id,name,email,active}:\n  1,Alice,alice@example.com,true",
		},
		{
			name:     "specify subset of fields",
			arrName:  "users",
			items:    []User{{ID: 1, Name: "Alice", Email: "alice@example.com", IsActive: true}},
			fields:   []string{"name", "active"},
			expected: "users[1]{name,active}:\n  Alice,true",
		},
		{
			name:     "multiple items",
			arrName:  "users",
			items:    []User{{ID: 1, Name: "Alice", IsActive: true}, {ID: 2, Name: "Bob", IsActive: false}},
			fields:   []string{"id", "name", "active"},
			expected: "users[2]{id,name,active}:\n  1,Alice,true\n  2,Bob,false",
		},
		{
			name:     "empty slice with fields",
			arrName:  "users",
			items:    []User{},
			fields:   []string{"id", "name"},
			expected: "users[0]{id,name}:",
		},
		{
			name:     "empty slice no fields",
			arrName:  "users",
			items:    []User{},
			fields:   nil,
			expected: "users[0]{}:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeTabularAuto(tt.arrName, tt.items, tt.fields...)
			if result != tt.expected {
				t.Errorf("EncodeTabularAuto() =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestEncodeTabularAutoWithNestedStruct(t *testing.T) {
	type Address struct {
		Name string
		City string
	}

	type Person struct {
		ID      int      `json:"id"`
		Name    string   `json:"name"`
		Address *Address `json:"address"`
	}

	people := []Person{
		{ID: 1, Name: "Alice", Address: &Address{Name: "Home", City: "NYC"}},
		{ID: 2, Name: "Bob", Address: nil},
	}

	// Nested struct with Name field should extract Name
	result := EncodeTabularAuto("people", people, "id", "name", "address")
	expected := "people[2]{id,name,address}:\n  1,Alice,Home\n  2,Bob,"

	if result != expected {
		t.Errorf("EncodeTabularAuto() =\n%q\nwant:\n%q", result, expected)
	}
}
