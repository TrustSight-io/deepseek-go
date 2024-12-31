package deepseek

// ChatMessage represents a message in a chat conversation.
type ChatMessage struct {
	Role         string      `json:"role"`
	Content      string      `json:"content"`
	Name         string      `json:"name,omitempty"`
	FunctionCall interface{} `json:"function_call,omitempty"`
}

// ChatCompletionRequest represents a request to create a chat completion.
type ChatCompletionRequest struct {
	Model            string               `json:"model"`
	Messages         []ChatMessage        `json:"messages"`
	MaxTokens        int                  `json:"max_tokens,omitempty"`
	Temperature      float32              `json:"temperature,omitempty"`
	TopP             float32              `json:"top_p,omitempty"`
	N                int                  `json:"n,omitempty"`
	Stream           bool                 `json:"stream,omitempty"`
	Stop             []string             `json:"stop,omitempty"`
	PresencePenalty  float32              `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32              `json:"frequency_penalty,omitempty"`
	Functions        []FunctionDefinition `json:"functions,omitempty"`
	FunctionCall     interface{}          `json:"function_call,omitempty"`
	JSONMode         bool                 `json:"json_mode,omitempty"`
}

// ChatCompletionResponse represents a response from the chat completion API.
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice.
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Usage provides token usage information.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// FunctionDefinition represents a function that can be called by the model.
type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters"`
}

// StreamChoice represents a choice in a streaming response.
type StreamChoice struct {
	Index        int    `json:"index"`
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// Delta represents the change in a streaming response.
type Delta struct {
	Role         string `json:"role,omitempty"`
	Content      string `json:"content,omitempty"`
	FunctionCall *struct {
		Name      string `json:"name,omitempty"`
		Arguments string `json:"arguments,omitempty"`
	} `json:"function_call,omitempty"`
}

// ChatCompletionStream represents a streaming response.
type ChatCompletionStream struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}
