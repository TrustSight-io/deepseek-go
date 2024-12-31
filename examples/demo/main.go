package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/trustsight/deepseek-go"
)

func main() {
	// Create a client with custom configuration
	client, err := deepseek.NewClient(
		os.Getenv("DEEPSEEK_API_KEY"),
		deepseek.WithHTTPClient(&http.Client{
			Timeout: time.Minute,
		}),
		deepseek.WithMaxRetries(2),
		deepseek.WithDebug(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// List available models
	fmt.Println("Listing available models...")
	models, err := client.ListModels(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nAvailable models:")
	for _, model := range models.Data {
		fmt.Printf("- %s:\n", model.ID)
		fmt.Printf("  Object: %s\n", model.Object)
		fmt.Printf("  Owner: %s\n", model.OwnedBy)
		fmt.Println()
	}

	modelID := "deepseek-chat"

	// Create a chat completion
	fmt.Println("\nCreating chat completion...")
	messages := []deepseek.Message{
		{
			Role:    deepseek.RoleSystem,
			Content: "You are a helpful assistant.",
		},
		{
			Role:    deepseek.RoleUser,
			Content: "What's the weather like today?",
		},
	}

	// First, estimate tokens
	estimate := client.EstimateTokensFromMessages(messages)
	fmt.Printf("\nEstimated tokens for messages: %d\n", estimate.EstimatedTokens)

	// Then create completion
	resp, err := client.CreateChatCompletion(
		ctx,
		&deepseek.ChatCompletionRequest{
			Model:       modelID,
			Messages:    messages,
			Temperature: 0.7,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nAssistant: %s\n", resp.Choices[0].Message.Content)

	// Create a streaming chat completion
	fmt.Println("\nCreating streaming chat completion...")
	stream, err := client.CreateChatCompletionStream(
		ctx,
		&deepseek.ChatCompletionRequest{
			Model: modelID,
			Messages: []deepseek.Message{
				{
					Role:    deepseek.RoleUser,
					Content: "Tell me a story about a brave knight.",
				},
			},
			Stream: true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	fmt.Print("\nStreaming response: ")
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Print(response.Choices[0].Delta.Content)
	}
	fmt.Println()

	// Check account balance
	fmt.Println("\nChecking account balance...")
	balance, err := client.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nAccount Status: %v\n", balance.IsAvailable)
	for _, info := range balance.BalanceInfos {
		fmt.Printf("\nBalance Info for %s:\n", info.Currency)
		fmt.Printf("  Total Balance: %s\n", info.TotalBalance)
		fmt.Printf("  Granted Balance: %s\n", info.GrantedBalance)
		fmt.Printf("  Topped Up Balance: %s\n", info.ToppedUpBalance)
	}

}
