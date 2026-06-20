.PHONY: all build test lint clean run docker

APP_NAME := teldrive
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -ldflags "-w -s -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

all: lint test build

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/teldrive

test:
	go test -v -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ coverage.out

run: build
	./bin/$(APP_NAME) run

docker:
	docker build -f docker/Dockerfile -t $(APP_NAME):$(VERSION) .

dev:
	air -c .air.toml

.DEFAULT_GOAL := all
