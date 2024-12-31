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
		fmt.Printf("  Context window: %d tokens\n", model.ContextWindow)
		fmt.Printf("  Max tokens: %d\n", model.MaxTokens)
		if model.PricingConfig != nil {
			fmt.Printf("  Pricing: %.4f/%s per prompt token, %.4f/%s per completion token\n",
				model.PricingConfig.PromptTokenPrice,
				model.PricingConfig.Currency,
				model.PricingConfig.CompletionTokenPrice,
				model.PricingConfig.Currency,
			)
		}
		fmt.Println()
	}

	// Get model configuration
	modelID := "deepseek-chat"
	fmt.Printf("Getting configuration for %s...\n", modelID)
	config, err := client.GetModelConfig(ctx, modelID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nModel configuration:\n")
	fmt.Printf("Context window: %d\n", config.ContextWindow)
	fmt.Printf("Max tokens: %d\n", config.MaxTokens)
	fmt.Printf("Supported features:\n")
	for feature, supported := range config.SupportedFeatures {
		fmt.Printf("- %s: %v\n", feature, supported)
	}

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
	estimate, err := client.EstimateTokensFromMessages(ctx, modelID, messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nEstimated tokens for messages: %d\n", estimate.TotalTokens)

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

	// Get usage history
	fmt.Println("\nGetting usage history...")
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7) // Last 7 days

	usage, err := client.GetUsage(ctx, &deepseek.UsageParams{
		StartTime:   startTime.Format(time.RFC3339),
		EndTime:     endTime.Format(time.RFC3339),
		Granularity: "daily",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nUsage history (%s to %s):\n", usage.StartTime, usage.EndTime)
	fmt.Printf("Total requests: %d\n", usage.Total.RequestCount)
	fmt.Printf("Total tokens: %d\n", usage.Total.TotalTokens)
	fmt.Printf("Total cost: %.4f\n", usage.Total.TotalCost)

	// Tokenize text
	text := "Hello, how are you today?"
	fmt.Printf("\nTokenizing text: %q\n", text)
	tokenization, err := client.TokenizeText(ctx, modelID, text)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nTokens:")
	for i, token := range tokenization.Tokens {
		fmt.Printf("%d. %q (ID: %d)\n", i+1, token, tokenization.IDs[i])
	}
}
