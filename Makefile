# Load environment variables from .env file
ifneq ("$(wildcard .env)","")
    include .env
    export
endif

# Database DSN for golang-migrate
DB_DSN := mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

.PHONY: dev seed docs migrate-up migrate-down migrate-fresh help

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Run the API development server with hot-reload
	@bash -c "trap 'kill 0' SIGINT SIGTERM EXIT; centrifugo --config=config.json & go run github.com/air-verse/air@latest -c .air.toml"

seed: ## Run database seeders
	go run cmd/seed/main.go

docs: ## Generate Swagger/Scalar documentation
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go

migrate-up: ## Run database migrations
	migrate -path database/migrations -database "$(DB_DSN)" up

migrate-down: ## Rollback last database migration
	migrate -path database/migrations -database "$(DB_DSN)" down 1

migrate-fresh: ## Drop all tables and run all migrations from scratch
	migrate -path database/migrations -database "$(DB_DSN)" drop -f
	migrate -path database/migrations -database "$(DB_DSN)" up

build: ## Build the application binaries
	go build -o bin/api cmd/api/main.go
	go build -o bin/seed cmd/seed/main.go

mock: ## Generate mocks using mockery
	go run github.com/vektra/mockery/v2@latest --all --dir=internal --inpackage

TEST_PKGS = $(shell go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}{{.ImportPath}}{{end}}' ./...)

test: ## Run all tests
	@if [ -n "$(TEST_PKGS)" ]; then go test -v $(TEST_PKGS); else echo "No tests found"; fi

test-unit: ## Run only unit tests
	@if [ -n "$(TEST_PKGS)" ]; then go test -v -short $(TEST_PKGS); else echo "No tests found"; fi

test-integration: ## Run only integration tests
	@if [ -n "$(TEST_PKGS)" ]; then go test -v -run Integration $(TEST_PKGS); else echo "No tests found"; fi

test-coverage: ## Run tests and show coverage
	@if [ -n "$(TEST_PKGS)" ]; then go test -v -coverprofile=coverage.out $(TEST_PKGS); go tool cover -html=coverage.out; else echo "No tests found"; fi
centrifugo: ## Run centrifugo locally
	centrifugo --config=centrifugo/config.json

