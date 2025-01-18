package deepseek

import (
	"encoding/json"
	"testing"
)

func TestJSONExtractor_ExtractJSON(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name     string
		response *ChatCompletionResponse
		schema   json.RawMessage
		want     Person
		wantErr  bool
	}{
		{
			name: "valid direct json",
			response: &ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Content: `{"name": "John", "age": 30}`,
						},
					},
				},
			},
			want: Person{
				Name: "John",
				Age:  30,
			},
			wantErr: false,
		},
		{
			name: "valid json in code block",
			response: &ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Content: "```json\n{\"name\": \"John\", \"age\": 30}\n```",
						},
					},
				},
			},
			want: Person{
				Name: "John",
				Age:  30,
			},
			wantErr: false,
		},
		{
			name: "valid json with schema validation",
			response: &ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Content: `{"name": "John", "age": 30}`,
						},
					},
				},
			},
			schema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"name": {"type": "string"},
					"age": {"type": "number"}
				}
			}`),
			want: Person{
				Name: "John",
				Age:  30,
			},
			wantErr: false,
		},
		{
			name:     "nil response",
			response: nil,
			wantErr:  true,
		},
		{
			name: "empty choices",
			response: &ChatCompletionResponse{
				Choices: []Choice{},
			},
			wantErr: true,
		},
		{
			name: "empty content",
			response: &ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Content: "",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid json",
			response: &ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Content: "not a json",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid schema",
			response: &ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Content: `{"name": "John", "age": 30}`,
						},
					},
				},
			},
			schema:  json.RawMessage(`invalid schema`),
			wantErr: true,
		},
		{
			name: "schema type mismatch",
			response: &ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Content: `{"name": "John", "age": 30}`,
						},
					},
				},
			},
			schema:  json.RawMessage(`{"type": "array"}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je := NewJSONExtractor(tt.schema)
			var got Person
			err := je.ExtractJSON(tt.response, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got != tt.want) {
				t.Errorf("ExtractJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONExtractor_extractJSONContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "direct json",
			content: `{"name": "John"}`,
			want:    `{"name": "John"}`,
		},
		{
			name:    "json in code block with newlines",
			content: "```json\n{\"name\": \"John\"}\n```",
			want:    `{"name": "John"}`,
		},
		{
			name:    "json in code block without newlines",
			content: "```json{\"name\": \"John\"}```",
			want:    `{"name": "John"}`,
		},
		{
			name:    "json in generic code block",
			content: "```\n{\"name\": \"John\"}\n```",
			want:    `{"name": "John"}`,
		},
		{
			name:    "json with surrounding text",
			content: "Here is the JSON:\n{\"name\": \"John\"}",
			want:    `{"name": "John"}`,
		},
		{
			name:    "invalid json",
			content: "not a json",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je := NewJSONExtractor(nil)
			if got := je.extractJSONContent(tt.content); got != tt.want {
				t.Errorf("extractJSONContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
