package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	// Get current balance
	fmt.Println("Getting current balance:")
	balance, err := client.GetBalance(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account Status: %v\n", balance.IsAvailable)
	for _, info := range balance.BalanceInfos {
		fmt.Printf("\nBalance Info for %s:\n", info.Currency)
		fmt.Printf("  Total Balance: %s\n", info.TotalBalance)
		fmt.Printf("  Granted Balance: %s\n", info.GrantedBalance)
		fmt.Printf("  Topped Up Balance: %s\n", info.ToppedUpBalance)
	}

	// Get usage history for the last 7 days
	fmt.Println("\nGetting usage history for the last 7 days:")
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	usage, err := client.GetUsage(context.Background(), &deepseek.UsageParams{
		StartTime:   startTime.Format(time.RFC3339),
		EndTime:     endTime.Format(time.RFC3339),
		Granularity: "daily",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nUsage from %s to %s:\n", usage.StartTime, usage.EndTime)
	for _, record := range usage.Data {
		fmt.Printf("\nTimestamp: %s\n", record.Timestamp)
		fmt.Printf("  Requests: %d\n", record.RequestCount)
		fmt.Printf("  Prompt Tokens: %d\n", record.PromptTokens)
		fmt.Printf("  Completion Tokens: %d\n", record.CompletionTokens)
		fmt.Printf("  Total Tokens: %d\n", record.TotalTokens)
		fmt.Printf("  Cost: %.4f\n", record.Cost)
	}

	fmt.Printf("\nTotals:\n")
	fmt.Printf("Total Requests: %d\n", usage.Total.RequestCount)
	fmt.Printf("Total Prompt Tokens: %d\n", usage.Total.PromptTokens)
	fmt.Printf("Total Completion Tokens: %d\n", usage.Total.CompletionTokens)
	fmt.Printf("Total Tokens: %d\n", usage.Total.TotalTokens)
	fmt.Printf("Total Cost: %.4f\n", usage.Total.TotalCost)
}
