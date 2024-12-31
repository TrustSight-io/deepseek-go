# DeepSeek Go Client

A Go client library for the DeepSeek API.

## Installation

```bash
go get github.com/trustsight/deepseek-go
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/trustsight/deepseek-go"
)

func main() {
    client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))

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
        panic(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)
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
- `models_test.go`: Model management and configuration
- `balance_test.go`: Account balance and usage operations
- `tokens_test.go`: Token counting and analysis utilities

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
