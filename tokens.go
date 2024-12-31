package deepseek

import (
	"unicode"
)

// TokenEstimate represents an estimated token count
type TokenEstimate struct {
	EstimatedTokens int `json:"estimated_tokens"`
}

// EstimateTokenCount estimates the number of tokens in a text based on character type ratios
func (c *Client) EstimateTokenCount(text string) *TokenEstimate {
	var total float64
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			// Chinese character ≈ 0.6 token
			total += 0.6
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
			// English character/number/symbol ≈ 0.3 token
			total += 0.3
		}
		// Skip whitespace and other characters
	}

	// Round up to nearest integer
	estimatedTokens := int(total + 0.5)
	if estimatedTokens < 1 {
		estimatedTokens = 1
	}

	return &TokenEstimate{
		EstimatedTokens: estimatedTokens,
	}
}

// EstimateTokensFromMessages estimates the number of tokens in a list of chat messages
func (c *Client) EstimateTokensFromMessages(messages []Message) *TokenEstimate {
	var totalTokens int

	for _, msg := range messages {
		// Add tokens for role (system/user/assistant)
		totalTokens += 3 // Approximate tokens for role

		// Add tokens for content
		totalTokens += c.EstimateTokenCount(msg.Content).EstimatedTokens

		// Add tokens for function call if present
		if msg.FunctionCall != nil {
			totalTokens += 3 // Approximate tokens for function name
			totalTokens += c.EstimateTokenCount(string(msg.FunctionCall.Arguments)).EstimatedTokens
		}
	}

	return &TokenEstimate{
		EstimatedTokens: totalTokens,
	}
}
