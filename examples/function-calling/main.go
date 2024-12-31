package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/trustsight/deepseek-go"
)

// WeatherInfo represents weather information for a location
type WeatherInfo struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
	Condition   string  `json:"condition"`
}

// getCurrentWeather simulates getting weather data
func getCurrentWeather(location string) WeatherInfo {
	// In a real application, this would make an API call to a weather service
	return WeatherInfo{
		Location:    location,
		Temperature: 22.5,
		Unit:        "celsius",
		Condition:   "sunny",
	}
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

	// Define the function that the model can call
	weatherFunction := deepseek.FunctionDefinition{
		Name:        "get_current_weather",
		Description: "Get the current weather in a given location",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "The city and state, e.g., San Francisco, CA",
				},
			},
			"required": []string{"location"},
		},
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		&deepseek.ChatCompletionRequest{
			Model: "deepseek-chat",
			Messages: []deepseek.ChatMessage{
				{
					Role:    "user",
					Content: "What's the weather like in Paris?",
				},
			},
			Functions: []deepseek.FunctionDefinition{weatherFunction},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Handle the function call
	if resp.Choices[0].Message.FunctionCall != nil {
		functionCall := resp.Choices[0].Message.FunctionCall.(map[string]interface{})
		var args struct {
			Location string `json:"location"`
		}
		if err := json.Unmarshal([]byte(functionCall["arguments"].(string)), &args); err != nil {
			log.Fatal(err)
		}

		// Get the weather data
		weather := getCurrentWeather(args.Location)

		// Send the function result back to continue the conversation
		resp, err = client.CreateChatCompletion(
			context.Background(),
			&deepseek.ChatCompletionRequest{
				Model: "deepseek-chat",
				Messages: []deepseek.ChatMessage{
					{
						Role:    "user",
						Content: "What's the weather like in Paris?",
					},
					{
						Role:         "assistant",
						Content:      "",
						FunctionCall: functionCall,
					},
					{
						Role:    "function",
						Name:    "get_current_weather",
						Content: fmt.Sprintf("The current weather in %s is %.1f°%s and %s", weather.Location, weather.Temperature, weather.Unit, weather.Condition),
					},
				},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
}
