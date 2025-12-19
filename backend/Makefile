.PHONY: run build test lint

run:
	go run ./cmd/api/main.go

build:
	go build -o bin/app ./cmd/api/main.go

test:
	go test ./...

lint:
	golangci-lint run
