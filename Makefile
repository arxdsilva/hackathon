# Hackathon Management Platform Makefile
# A comprehensive build and development tool for the Buffalo application

.PHONY: help setup install deps build run dev test test-unit test-integration db-setup db-migrate db-reset db-seed assets assets-dev assets-build docker-build docker-run docker-stop clean lint fmt vet mod-tidy

# Default target
help: ## Show this help message
	@echo "Hackathon Management Platform - Development Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make setup          # Initial project setup"
	@echo "  make dev            # Start development server"
	@echo "  make test           # Run all tests"
	@echo "  make test-unit      # Run unit tests only (fast)"
	@echo "  make db-migrate     # Run database migrations"

# Setup and Installation
setup: deps db-setup assets-dev ## Complete project setup (dependencies, database, assets)
	@echo "✅ Project setup complete!"

install: ## Install Go and Node.js dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	yarn install

deps: ## Download Go dependencies
	@echo "📦 Downloading Go dependencies..."
	go mod download

# Development
dev: ## Start development server with hot reload
	@echo "🚀 Starting development server..."
	buffalo dev

run: ## Run the application (production mode)
	@echo "🚀 Starting application..."
	buffalo task

# Building
build: assets-build ## Build the application for production
	@echo "🔨 Building application..."
	buffalo build -o bin/app

# Testing
test: test-unit test-integration ## Run all tests (unit + integration)

test-unit: ## Run unit tests only (no database required)
	@echo "🧪 Running unit tests..."
	go test ./actions -run "TestAdminConfigUpdate_ValidationError|TestBindConfigBooleans|TestAdminConfigUpdate_NoDatabase|TestAdminConfigIndex_NoDatabase|TestAdminConfigUpdate_InvalidRole" -v

test-integration: ## Run integration tests (requires database)
	@echo "🧪 Running integration tests..."
	go test ./actions -run ".*Integration.*" -v

test-full: ## Run full test suite with database setup (requires Docker)
	@echo "🧪 Running full test suite..."
	./test.sh

# Database Operations
db-setup: ## Set up development database
	@echo "🗄️ Setting up database..."
	@echo "Make sure PostgreSQL is running on port 5432"
	@echo "Creating database: hackathon_development"
	createdb hackathon_development || echo "Database may already exist"
	psql -d hackathon_development -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";" || echo "Extension may already exist"

db-migrate: ## Run database migrations
	@echo "🗃️ Running database migrations..."
	buffalo db migrate up

db-reset: ## Reset database (drop and recreate)
	@echo "💥 Resetting database..."
	buffalo db reset

db-seed: ## Seed database with initial data
	@echo "🌱 Seeding database..."
	buffalo task db:seed

db-docker-up: ## Start database with Docker Compose
	@echo "🐳 Starting database with Docker..."
	docker-compose up -d postgres

db-docker-down: ## Stop database Docker container
	@echo "🐳 Stopping database..."
	docker-compose down

# Assets
assets: assets-build ## Build production assets

assets-dev: ## Build development assets with watch mode
	@echo "🎨 Building development assets..."
	yarn run dev

assets-build: ## Build production assets
	@echo "🎨 Building production assets..."
	yarn run build

# Docker Operations
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t hackathon .

docker-run: ## Run application in Docker
	@echo "🐳 Running application in Docker..."
	docker run -p 3000:3000 -e GO_ENV=production hackathon

docker-stop: ## Stop all Docker containers
	@echo "🐳 Stopping Docker containers..."
	docker-compose down

docker-dev: ## Start full development environment with Docker
	@echo "🐳 Starting development environment..."
	docker-compose up -d

# Code Quality
lint: ## Run linter
	@echo "🔍 Running linter..."
	golangci-lint run

fmt: ## Format Go code
	@echo "📝 Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

mod-tidy: ## Clean up Go modules
	@echo "🧹 Tidying Go modules..."
	go mod tidy

# Cleaning
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -rf tmp/
	rm -rf public/assets/
	rm -rf node_modules/.cache/

clean-all: clean ## Clean everything including dependencies
	@echo "🧹 Deep cleaning..."
	rm -rf node_modules/
	go clean -cache
	go clean -testcache
	go clean -modcache

# Utility
version: ## Show application version
	@echo "🏷️ Application version:"
	@grep -E "^version:" config/buffalo-app.toml || echo "Version not found in config"

status: ## Show development environment status
	@echo "📊 Development Environment Status:"
	@echo "Go version: $$(go version)"
	@echo "Node version: $$(node --version)"
	@echo "Yarn version: $$(yarn --version)"
	@echo "Database: $$(pg_isready -h localhost -p 5432 >/dev/null 2>&1 && echo '✅ Connected' || echo '❌ Not connected')"
	@echo "Docker: $$(docker --version 2>/dev/null | head -1 || echo '❌ Not installed')"