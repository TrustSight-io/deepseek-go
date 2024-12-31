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

	// List available models
	fmt.Println("Listing available models:")
	models, err := client.ListModels(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, model := range models.Data {
		fmt.Printf("- %s:\n", model.ID)
		fmt.Printf("  Context Window: %d\n", model.ContextWindow)
		fmt.Printf("  Max Tokens: %d\n", model.MaxTokens)
		if model.PricingConfig != nil {
			fmt.Printf("  Pricing:\n")
			fmt.Printf("    Prompt Token Price: %f %s\n", model.PricingConfig.PromptTokenPrice, model.PricingConfig.Currency)
			fmt.Printf("    Completion Token Price: %f %s\n", model.PricingConfig.CompletionTokenPrice, model.PricingConfig.Currency)
		}
		fmt.Printf("  Capabilities:\n")
		for capability, enabled := range model.Capabilities {
			fmt.Printf("    %s: %v\n", capability, enabled)
		}
		fmt.Println()
	}

	// Get specific model details
	modelID := "deepseek-chat"
	fmt.Printf("\nGetting details for model %s:\n", modelID)
	model, err := client.GetModel(context.Background(), modelID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Model ID: %s\n", model.ID)
	fmt.Printf("Created: %d\n", model.Created)
	fmt.Printf("Owner: %s\n", model.OwnedBy)
	fmt.Printf("Context Window: %d\n", model.ContextWindow)
	fmt.Printf("Max Tokens: %d\n", model.MaxTokens)

	// Get model configuration
	fmt.Printf("\nGetting configuration for model %s:\n", modelID)
	config, err := client.GetModelConfig(context.Background(), modelID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Context Window: %d\n", config.ContextWindow)
	fmt.Printf("Max Tokens: %d\n", config.MaxTokens)
	fmt.Printf("\nSupported Features:\n")
	for feature, supported := range config.SupportedFeatures {
		fmt.Printf("- %s: %v\n", feature, supported)
	}
	fmt.Printf("\nDefault Parameters:\n")
	for param, value := range config.DefaultParameters {
		fmt.Printf("- %s: %v\n", param, value)
	}
}
