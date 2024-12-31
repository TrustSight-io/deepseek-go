package main

import (
	"fmt"
	"log"
	"os"

	"github.com/trustsight/deepseek-go"
)

func main() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable is required")
	}

	client, err := deepseek.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	// Estimate tokens in a text
	text := "Hello, how are you doing today? I hope you're having a great day!"
	fmt.Printf("Estimating tokens in text: %q\n", text)

	estimate := client.EstimateTokenCount(text)
	fmt.Printf("Estimated tokens: %d\n", estimate.EstimatedTokens)

	// Estimate tokens in chat messages
	messages := []deepseek.Message{
		{
			Role:    deepseek.RoleSystem,
			Content: "You are a helpful assistant.",
		},
		{
			Role:    deepseek.RoleUser,
			Content: "What's the weather like today?",
		},
		{
			Role:    deepseek.RoleAssistant,
			Content: "I don't have access to real-time weather information. You would need to check a weather service or look outside for current conditions.",
		},
	}

	fmt.Printf("\nEstimating tokens in chat messages:\n")
	for _, msg := range messages {
		fmt.Printf("- %s: %q\n", msg.Role, msg.Content)
	}

	estimate = client.EstimateTokensFromMessages(messages)
	fmt.Printf("\nEstimated total tokens: %d\n", estimate.EstimatedTokens)

	// Example with Chinese text
	chineseText := "你好，世界！"
	fmt.Printf("\nEstimating tokens in Chinese text: %q\n", chineseText)
	estimate = client.EstimateTokenCount(chineseText)
	fmt.Printf("Estimated tokens: %d\n", estimate.EstimatedTokens)

	// Example with mixed text
	mixedText := "Hello 世界! How are you? 你好吗？"
	fmt.Printf("\nEstimating tokens in mixed text: %q\n", mixedText)
	estimate = client.EstimateTokenCount(mixedText)
	fmt.Printf("Estimated tokens: %d\n", estimate.EstimatedTokens)
}
