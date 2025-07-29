# Makefile for Task Management API

.PHONY: all build test lint fmt run clean coverage coverage-html

all: build

build:
	go build -o bin/taskmanager ./cmd/server

test:
	go test ./...

lint:
	go vet ./...
	staticcheck ./...

staticcheck:
	staticcheck ./...

fmt:
	gofmt -s -w .

run:
	go run ./cmd/server

clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t taskmanager:latest .

# Scan Docker image with Trivy
docker-scan:
	trivy image --exit-code 1 --severity CRITICAL,HIGH taskmanager:latest

# Run tests with coverage and output coverage.out
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Generate HTML coverage report and open in browser
coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html
