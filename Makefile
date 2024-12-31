.PHONY: lint test coverage build

lint:
	golangci-lint run

test:
	go test -v -race ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

build:
	go build .
