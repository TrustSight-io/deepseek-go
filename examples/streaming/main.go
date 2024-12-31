package main

import (
	"context"
	"fmt"
	"io"
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

	stream, err := client.CreateChatCompletionStream(
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
					Content: "Write a short story about a robot learning to paint.",
				},
			},
			Temperature: 0.7,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	fmt.Print("Assistant: ")
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("\nStream finished")
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(response.Choices[0].Delta.Content)
	}
}