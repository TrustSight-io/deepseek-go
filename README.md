# DeepSeek Go Client

A Go client library for the DeepSeek API. This library provides a simple and efficient way to interact with DeepSeek's language models.

## Features

- Chat completions with streaming support
- Function calling capabilities
- JSON mode support
- Automatic retries and error handling
- Customizable options and configurations

## Installation

```bash
go get github.com/trustsight/deepseek-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    deepseek "github.com/trustsight/deepseek-go"
)

func main() {
    client := deepseek.NewClient("your-api-key")
    
    resp, err := client.CreateChatCompletion(context.Background(), &deepseek.ChatCompletionRequest{
        Messages: []deepseek.ChatMessage{
            {
                Role:    "user",
                Content: "Hello, how are you?",
            },
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println(resp.Choices[0].Message.Content)
}
```

## Examples

Check out the [examples](./examples) directory for more detailed usage examples:

- [Basic Chat](./examples/chat)
- [Streaming](./examples/streaming)
- [Function Calling](./examples/function-calling)
- [JSON Mode](./examples/json-mode)

## Documentation

For detailed documentation and API reference, visit [pkg.go.dev](https://pkg.go.dev/github.com/trustsight/deepseek-go).

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on how to submit pull requests, report issues, and contribute to the project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
