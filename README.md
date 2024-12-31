# DeepSeek Go Client

A Go client library for the DeepSeek API.

## Installation

```bash
go get github.com/trustsight-io/deepseek-go
```

## Usage

### Chat Completion

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/trustsight-io/deepseek-go"
)

func main() {
    client, err := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.CreateChatCompletion(
        context.Background(),
        &deepseek.ChatCompletionRequest{
            Model: "deepseek-chat",
            Messages: []deepseek.Message{
                {
                    Role:    deepseek.RoleUser,
                    Content: "Hello!",
                },
            },
        },
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)
}
```

### Streaming Chat Completion

```go
stream, err := client.CreateChatCompletionStream(
    context.Background(),
    &deepseek.ChatCompletionRequest{
        Model: "deepseek-chat",
        Messages: []deepseek.Message{
            {
                Role:    deepseek.RoleUser,
                Content: "Tell me a story.",
            },
        },
        Stream: true,
    },
)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for {
    response, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(response.Choices[0].Delta.Content)
}
```

### Token Estimation

```go
// Estimate tokens in text
text := "Hello 世界!"
estimate := client.EstimateTokenCount(text)
fmt.Printf("Estimated tokens: %d\n", estimate.EstimatedTokens)

// Estimate tokens in messages
messages := []deepseek.Message{
    {
        Role:    deepseek.RoleSystem,
        Content: "You are a helpful assistant.",
    },
    {
        Role:    deepseek.RoleUser,
        Content: "Hello!",
    },
}
estimate = client.EstimateTokensFromMessages(messages)
fmt.Printf("Estimated total tokens: %d\n", estimate.EstimatedTokens)
```

### List Available Models

```go
models, err := client.ListModels(context.Background())
if err != nil {
    log.Fatal(err)
}

for _, model := range models.Data {
    fmt.Printf("- %s:\n", model.ID)
    fmt.Printf("  Object: %s\n", model.Object)
    fmt.Printf("  Owner: %s\n", model.OwnedBy)
}
```

### Check Account Balance

```go
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
```

## Running Tests

### Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Add your DeepSeek API key to `.env`:
   ```
   DEEPSEEK_API_KEY=your_api_key_here
   ```

   You can get your API key from the [DeepSeek Platform](https://api.deepseek.com).

3. (Optional) Configure test timeout:
   ```
   # Default is 30s, increase for slower connections
   TEST_TIMEOUT=1m
   ```

### Test Organization

The tests are organized into several files:
- `client_test.go`: Client configuration and error handling
- `chat_test.go`: Chat completion functionality (including streaming)
- `models_test.go`: Model listing and retrieval
- `balance_test.go`: Account balance operations
- `tokens_test.go`: Token estimation utilities

### Running Tests

1. Run all tests (requires API key):
   ```bash
   go test -v ./...
   ```

2. Run tests in short mode (skips API calls):
   ```bash
   go test -v -short ./...
   ```

3. Run tests with race detection:
   ```bash
   go test -v -race ./...
   ```

4. Run tests with coverage:
   ```bash
   go test -v -coverprofile=coverage.txt -covermode=atomic ./...
   ```

   View coverage in browser:
   ```bash
   go tool cover -html=coverage.txt
   ```

5. Run specific test:
   ```bash
   # Example: Run only chat completion tests
   go test -v -run TestCreateChatCompletion ./...
   ```

### Test Environment Variables

- `DEEPSEEK_API_KEY`: Your DeepSeek API key (required for API tests)
- `TEST_TIMEOUT`: Test timeout duration (default: 30s)

### Common Issues

1. "invalid character 'A' looking for beginning of value":
   - This usually means the API returned HTML instead of JSON
   - Check if your API key is valid
   - Verify you're using the correct API base URL

2. "context deadline exceeded":
   - Increase test timeout in `.env`: `TEST_TIMEOUT=60s`
   - Check your internet connection

3. "unexpected end of JSON input":
   - The API response was truncated
   - Could indicate network issues
   - Try increasing the timeout

4. Tests are skipped:
   - Make sure `DEEPSEEK_API_KEY` is set in `.env`
   - For quick testing, use `-short` flag to skip API tests

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

___

<a href='https://ko-fi.com/O5O31892PV' target='_blank'><img height='36' style='border:0px;height:36px;' src='https://storage.ko-fi.com/cdn/kofi6.png?v=6' border='0' alt='Buy Me a Coffee at ko-fi.com' /></a>

Created by [pocok](pocok.dev) @ [TrustSight](https://www.trustsight.io)
