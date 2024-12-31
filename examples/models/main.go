package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/trustsight-io/deepseek-go"
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
		fmt.Printf("  Object: %s\n", model.Object)
		fmt.Printf("  Owner: %s\n", model.OwnedBy)
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
	fmt.Printf("Object: %s\n", model.Object)
	fmt.Printf("Owner: %s\n", model.OwnedBy)
}
