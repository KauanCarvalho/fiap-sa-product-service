include .env

APP_NAME := fiap_sa_product_service
BIN_DIR := bin
DATABASE_URL := "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?charset=utf8mb4&parseTime=true"
DOCKER_COMPOSE := docker-compose
ENV_FILE := .env
GO ?= go
GOBIN := $(shell ./resolve-gobin.sh)
MIGRATE := $(GOBIN)/migrate
TEST_DATABASE_URL := "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)_test?charset=utf8mb4&parseTime=true"

ifeq ($(DB_ENV),test)
  DATABASE_URL := $(TEST_DATABASE_URL)
endif

.DEFAULT_GOAL := help
.PHONY: help deps setup-git-hooks lint check-coverage migration migrate-up migrate-down test coverage-html build-api run-api run-api-air docker-up docker-down install-tools swag

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  help            Show this help message"
	@echo "  deps            Install dependencies"
	@echo "  setup-git-hooks Install Git hooks using Lefthook"
	@echo "  lint            Run linters"
	@echo "  check-coverage  Check test coverage"
	@echo "  migration       Create a new database migration files"
	@echo "  migrate-up      Apply database migrations"
	@echo "  migrate-down    Revert one database migration"
	@echo "  test            Run tests"
	@echo "  coverage-html   Generate HTML coverage report"
	@echo "  build-api       Build the API"
	@echo "  run-api         Run the API"
	@echo "  run-api-air     Run the API with live reloading"
	@echo "  docker-up       Start Docker container(s)"
	@echo "  docker-down     Stop Docker containers"
	@echo "  install-tools   Install tools with third-party dependencies"
	@echo "  swag            Generate Swagger documentation"

deps:
	@echo "Installing dependencies..."
	$(GO) mod download

setup-git-hooks: deps
	@echo "Installing Git hooks with Lefthook..."
	$(GO) tool lefthook install

lint: deps
	@echo "Running linter..."
	$(GO) tool golangci-lint run ./... --config .golangci.yml

test:
	@echo "Running tests..."
	DB_NAME=$(DB_NAME)_test $(GO) tool godotenv -f $(ENV_FILE) $(GO) test ./... -coverprofile=coverage.out -cover -p 1

check-coverage: test
	@echo "Checking coverage..."
	$(GO) tool go-test-coverage --config=./.testcoverage.yml

coverage-html: test
	@echo "Openning coverage report..."
	$(GO) tool cover -html=coverage.out

migration:
	@echo "Creating migration files for '$(filter-out $@,$(MAKECMDGOALS))'..."
	$(MIGRATE) create -ext sql -dir ./migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@echo "Applying database migrations..."
	$(MIGRATE) -verbose -database $(DATABASE_URL) -path ./migrations up

migrate-down:
	@echo "Reverting one database migration..."
	$(MIGRATE) -verbose -database $(DATABASE_URL) -path ./migrations down 1

build-api:
	@echo "Building api..."
	$(GO) build -o $(BIN_DIR)/$(APP_NAME)_api ./cmd/api/main.go

run-api: build-api
	@echo "Running api..."
	$(BIN_DIR)/$(APP_NAME)_api

run-api-air: deps
	@echo "Running api with live reloading..."
	$(GO) tool air -c .air.api.toml

docker-up:
	@echo "Starting Docker container(s)..."
	$(DOCKER_COMPOSE) up -d $(filter-out $@,$(MAKECMDGOALS))

docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

install-tools:
	@echo "Installing tools..."
	$(GO) install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@if command -v asdf > /dev/null 2>&1; then \
		if asdf plugin-list | grep -Fq golang; then \
			asdf reshim golang; \
		fi \
	fi

swag:
	@echo "Generating Swagger documentation..."
	$(GO) tool swag init --parseDependency --parseInternal -g cmd/api/main.go -o ./swagger --ot json,go
