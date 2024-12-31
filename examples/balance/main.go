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
}
