SHELL := /bin/bash
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org

.DEFAULT_GOAL: all

LDFLAGS=-ldflags "-s -w"

.PHONY: all build check clean format help serve test tidy

all: check test build ## Default target: check, test, build

build: ## Build all executables, located under ./bin/
	@echo "Building..."
	@CGO_ENABLED=0 go build -o ./bin/example -trimpath $(LDFLAGS) ./example/...
	@CGO_ENABLED=0 go build -o ./bin/pocketbase -trimpath $(LDFLAGS) ./cmd/pocketbase/...

serve: build ## Run the pocketbase server
	@echo "Running server..."
	@./bin/pocketbase serve

clean: ## Remove all artifacts from ./bin/ and ./resources
	@rm -rf ./bin/* ./resources/*

format: ## Format go code with goimports
	@go install golang.org/x/tools/cmd/goimports@latest
	@goimports -l -w .

test: ## Run tests
	@go test -shuffle=on -race ./...

tidy: ## Run go mod tidy
	@go mod tidy

check: ## Linting and static analysis
	@if test ! -e ./bin/golangci-lint; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh; \
	fi

	@./bin/golangci-lint run -c .golangci.yml

	# TODO in 2023 check if govulncheck is a part of golangci-lint
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@govulncheck ./...

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
