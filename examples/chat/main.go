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

	resp, err := client.CreateChatCompletion(
		context.Background(),
		&deepseek.ChatCompletionRequest{
			Model: "deepseek-chat",
			Messages: []deepseek.ChatMessage{
				{
					Role:    "system",
					Content: "You are a helpful assistant.",
				},
				{
					Role:    "user",
					Content: "What is the capital of France?",
				},
			},
			Temperature: 0.7,
			MaxTokens:   100,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
}
