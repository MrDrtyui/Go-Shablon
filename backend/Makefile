.PHONY: run build test lint swagger swagger-clean dev

run:
	go run ./cmd/api/main.go

build:
	go build -o bin/app ./cmd/api/main.go

test:
	go test ./...

lint:
	golangci-lint run

swagger:
	swag init -g cmd/api/main.go -o docs

swagger-clean:
	rm -rf docs/

dev: swagger
	go run ./cmd/api/main.go
