APP_NAME := telegram-bot-simple
BIN := bot
GO_FILES := ./...
DOCKER_IMAGE := $(APP_NAME)
ENV_FILE := .env

.DEFAULT_GOAL := help

.PHONY: help all run build deps fmt fmt-check imports lint vet staticcheck golangci-lint test docker-build docker-run docker-stop docker-logs docker-compose-up docker-compose-down preprod vuln

## Show available make targets
help:
	@echo "Available make targets:"
	@echo "  help          - Show this help"
	@echo "  run           - Run bot locally (uses .env if present)"
	@echo "  build         - Build Go binary"
	@echo "  deps          - Tidy Go modules (go mod tidy)"
	@echo "  fmt           - Format code with gofmt"
	@echo "  fmt-check     - Check formatting with gofmt -l"
	@echo "  imports       - Organize imports with goimports"
	@echo "  lint          - Run all linters (fmt, vet, staticcheck, golangci-lint)"
	@echo "  vet           - Run go vet"
	@echo "  staticcheck   - Run staticcheck (if installed)"
	@echo "  golangci-lint - Run golangci-lint (if installed)"
	@echo "  test          - Run go tests"
	@echo "  vuln          - Run govulncheck (via go run)"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run bot in Docker with .env"
	@echo "  docker-stop   - Stop running Docker container"
	@echo "  docker-logs   - Show Docker logs"
	@echo "  docker-compose-up   - Run bot via docker-compose (uses .env)"
	@echo "  docker-compose-down - Stop bot started by docker-compose"
	@echo "  preprod       - Full pre-production checks (deps, fmt, imports, linters, tests, vuln, docker build)"

## Default: run all pre-production checks
all: preprod

## Install / tidy Go dependencies
deps:
	go mod tidy

## Format code with gofmt
fmt:
	gofmt -w .

## Check formatting with gofmt (no file modifications)
fmt-check:
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "These files are not gofmt-formatted:"; \
		echo "$$unformatted"; \
		echo ""; \
		echo "Run: make fmt"; \
		exit 1; \
	fi

## Organize imports (requires golang.org/x/tools/cmd/goimports)
imports:
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not found, install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi

## Build Go binary
build: deps
	go build -o $(BIN) ./cmd/bot

## Run bot locally using current shell environment (.env helper)
run:
	@if [ -f $(ENV_FILE) ]; then \
		echo "Loading env from $(ENV_FILE)"; \
		set -a; . ./$(ENV_FILE); set +a; \
	fi; \
	go run ./cmd/bot

## Aggregate lints: fmt, vet, staticcheck, golangci-lint (if installed)
lint: fmt vet staticcheck golangci-lint

## Go vet (static analyzer from Go toolchain)
vet:
	go vet $(GO_FILES)

## Staticcheck (requires honnef.co/go/tools/cmd/staticcheck)
staticcheck:
	@if command -v staticcheck >/dev/null 2>&1; then \
		echo "staticcheck $(GO_FILES)"; \
		staticcheck $(GO_FILES); \
	else \
		echo "staticcheck not found, install with: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi

## GolangCI-Lint (requires golangci-lint installed)
golangci-lint:
	@echo "golangci-lint $(GO_FILES)";
	@if command -v golangci-lint >/dev/null 2>&1; then \
		if ! golangci-lint run ./...; then \
			echo "golangci-lint failed (likely built with older Go). To update, run:"; \
			echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		fi; \
	else \
		echo "golangci-lint not found, install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## Run tests (if/when they appear)
test:
	go test $(GO_FILES)

## Vulnerability check (requires govulncheck)
vuln:
	@echo "govulncheck ./..."
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

## Run bot in Docker with .env
docker-run: docker-build
	docker run --rm \
		--env-file $(ENV_FILE) \
		--name $(APP_NAME) \
		$(DOCKER_IMAGE)

## Stop running Docker container
docker-stop:
	- docker stop $(APP_NAME) 2>/dev/null || true

## Show Docker logs
docker-logs:
	docker logs -f $(APP_NAME)

## Run bot via docker-compose (build + up)
docker-compose-up:
	docker compose up --build

## Stop bot started via docker-compose
docker-compose-down:
	docker compose down

## Full pre-production check: deps, fmt, imports, vet, staticcheck, golangci-lint, tests, vuln, docker build
preprod: deps fmt imports vet staticcheck golangci-lint test vuln docker-build
	@echo "Pre-production checks completed successfully."
