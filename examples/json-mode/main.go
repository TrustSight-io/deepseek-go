package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/trustsight-io/deepseek-go"
)

// Product represents a product in an e-commerce system
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	InStock     bool    `json:"in_stock"`
}

func main() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable is required")
	}

	client, err := deepseek.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	// Request product information in JSON format
	resp, err := client.CreateChatCompletion(
		context.Background(),
		&deepseek.ChatCompletionRequest{
			Model: "deepseek-chat",
			Messages: []deepseek.Message{
				{
					Role: deepseek.RoleSystem,
					Content: `You are a product information generator. 
					Generate product information in JSON format following the Product struct schema.
					Always generate valid JSON that can be parsed into the Product struct.`,
				},
				{
					Role:    deepseek.RoleUser,
					Content: "Generate a product entry for a high-end laptop computer.",
				},
			},
			JSONMode: true, // Enable JSON mode
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create a JSON extractor without schema validation
	extractor := deepseek.NewJSONExtractor(nil)
	var product Product
	if err := extractor.ExtractJSON(resp, &product); err != nil {
		log.Fatalf("Failed to extract JSON: %v", err)
	}

	// Print the formatted product information
	fmt.Println("Product Information:")
	fmt.Printf("ID: %s\n", product.ID)
	fmt.Printf("Name: %s\n", product.Name)
	fmt.Printf("Description: %s\n", product.Description)
	fmt.Printf("Price: $%.2f\n", product.Price)
	fmt.Printf("Category: %s\n", product.Category)
	fmt.Printf("In Stock: %v\n", product.InStock)

	// Print the raw JSON response
	fmt.Println("\nRaw JSON response:")
	fmt.Println(resp.Choices[0].Message.Content)
}
