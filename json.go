package deepseek

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONExtractor helps extract structured data from LLM responses
type JSONExtractor struct {
	// Optional JSON schema for validation
	schema json.RawMessage
}

// NewJSONExtractor creates a new JSONExtractor instance
func NewJSONExtractor(schema json.RawMessage) *JSONExtractor {
	return &JSONExtractor{
		schema: schema,
	}
}

// ExtractJSON attempts to extract and parse JSON from an LLM response
func (je *JSONExtractor) ExtractJSON(response *ChatCompletionResponse, target interface{}) error {
	if response == nil {
		return fmt.Errorf("response cannot be nil")
	}

	if len(response.Choices) == 0 {
		return fmt.Errorf("no choices in response")
	}

	content := response.Choices[0].Message.Content
	if content == "" {
		return fmt.Errorf("empty content in response")
	}

	// Try to find JSON content with or without code blocks
	jsonStr := je.extractJSONContent(content)
	if jsonStr == "" {
		return fmt.Errorf("no valid JSON content found in response")
	}

	// If schema is provided, validate the JSON against it
	if je.schema != nil {
		if err := je.validateJSON([]byte(jsonStr)); err != nil {
			return fmt.Errorf("JSON validation failed: %w", err)
		}
	}

	// Parse the JSON into the target structure
	if err := json.Unmarshal([]byte(jsonStr), target); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}

// validateJSON validates JSON content against the schema
func (je *JSONExtractor) validateJSON(data []byte) error {
	var schemaMap map[string]interface{}
	if err := json.Unmarshal(je.schema, &schemaMap); err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON data: %w", err)
	}

	// Basic type validation
	if schemaType, ok := schemaMap["type"].(string); ok {
		switch schemaType {
		case "object":
			if _, ok := jsonData.(map[string]interface{}); !ok {
				return fmt.Errorf("expected object, got %T", jsonData)
			}
		case "array":
			if _, ok := jsonData.([]interface{}); !ok {
				return fmt.Errorf("expected array, got %T", jsonData)
			}
		}
	}

	return nil
}

// extractJSONContent attempts to extract valid JSON from the content
func (je *JSONExtractor) extractJSONContent(content string) string {
	content = strings.TrimSpace(content)

	// First try to parse the entire content as JSON
	if json.Valid([]byte(content)) {
		return content
	}

	// Try to find JSON content in various formats
	patterns := []struct {
		extract func(string) string
	}{
		// Try code blocks first
		{func(s string) string {
			return je.extractBetween(s, "```json\n", "```")
		}},
		{func(s string) string {
			return je.extractBetween(s, "```json", "```")
		}},
		{func(s string) string {
			return je.extractBetween(s, "```\n", "```")
		}},
		{func(s string) string {
			return je.extractBetween(s, "```", "```")
		}},
		// Then try to find JSON objects or arrays
		{func(s string) string {
			return je.findJSONInText(s)
		}},
	}

	for _, pattern := range patterns {
		if result := pattern.extract(content); result != "" {
			return result
		}
	}

	return ""
}

// extractBetween extracts content between start and end markers
func (je *JSONExtractor) extractBetween(content, start, end string) string {
	startIdx := strings.Index(content, start)
	if startIdx == -1 {
		return ""
	}

	content = content[startIdx+len(start):]
	endIdx := strings.Index(content, end)
	if endIdx == -1 {
		return ""
	}

	result := strings.TrimSpace(content[:endIdx])
	if json.Valid([]byte(result)) {
		return result
	}
	return ""
}

// findJSONInText attempts to find valid JSON objects or arrays in text
func (je *JSONExtractor) findJSONInText(content string) string {
	// Find potential JSON start
	var start, end int
	var found bool

	// Try to find JSON object
	start = strings.Index(content, "{")
	if start != -1 {
		end = je.findMatchingBrace(content[start:])
		if end != -1 {
			found = true
			end += start
		}
	}

	// If no object found, try to find JSON array
	if !found {
		start = strings.Index(content, "[")
		if start != -1 {
			end = je.findMatchingBracket(content[start:])
			if end != -1 {
				found = true
				end += start
			}
		}
	}

	if !found {
		return ""
	}

	result := content[start : end+1]
	if json.Valid([]byte(result)) {
		return result
	}
	return ""
}

// findMatchingBrace finds the matching closing brace for a JSON object
func (je *JSONExtractor) findMatchingBrace(content string) int {
	if !strings.HasPrefix(content, "{") {
		return -1
	}

	depth := 0
	for i, char := range content {
		switch char {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// findMatchingBracket finds the matching closing bracket for a JSON array
func (je *JSONExtractor) findMatchingBracket(content string) int {
	if !strings.HasPrefix(content, "[") {
		return -1
	}

	depth := 0
	for i, char := range content {
		switch char {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
