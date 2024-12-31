package main

import (
	"context"
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
	model := "deepseek-chat"

	// Count tokens in a text
	text := "Hello, how are you doing today? I hope you're having a great day!"
	fmt.Printf("Counting tokens in text: %q\n", text)

	count, err := client.CountTokens(context.Background(), model, text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total tokens: %d\n", count.TotalTokens)
	fmt.Printf("Prompt tokens: %d\n", count.Details.Prompt)
	if count.Details.Truncated {
		fmt.Println("Note: Text was truncated")
	}

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

	estimate, err := client.EstimateTokensFromMessages(context.Background(), model, messages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nEstimated total tokens: %d\n", estimate.TotalTokens)
	fmt.Printf("Estimated prompt tokens: %d\n", estimate.Details.Prompt)
	if estimate.Details.Truncated {
		fmt.Println("Note: Messages were truncated")
	}

	// Tokenize text
	fmt.Printf("\nTokenizing text: %q\n", text)
	tokenization, err := client.TokenizeText(context.Background(), model, text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nTokens:")
	for i, token := range tokenization.Tokens {
		fmt.Printf("%d. %q (ID: %d)\n", i+1, token, tokenization.IDs[i])
	}
}
