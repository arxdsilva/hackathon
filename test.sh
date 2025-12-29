#!/bin/bash

# Test runner script that ensures database isolation
# This script sets up a temporary test database and runs tests

set -e

echo "🚀 Setting up isolated test database..."

# Check if Docker is available
if command -v docker &> /dev/null; then
    echo "📦 Starting test database with Docker..."

    # Start test database
    docker-compose -f docker-compose.test.yml up -d

    # Wait for database to be ready
    echo "⏳ Waiting for test database to be ready..."
    for i in {1..30}; do
        if docker-compose -f docker-compose.test.yml exec -T postgres-test pg_isready -U postgres -d hackathon_test_isolated &> /dev/null; then
            break
        fi
        sleep 2
    done

    # Run migrations
    echo "🗃️ Running migrations on test database..."
    TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5433/hackathon_test_isolated?sslmode=disable" buffalo db migrate up --env test

    # Run unit tests first (fast feedback)
    echo "🧪 Running unit tests (no database required)..."
    go test ./actions -run "TestAdminConfigUpdate_ValidationError|TestBindConfigBooleans|TestAdminConfigUpdate_NoDatabase|TestAdminConfigIndex_NoDatabase|TestAdminConfigUpdate_InvalidRole" -v

    # Run integration tests
    echo "🧪 Running integration tests..."
    go test ./actions -run ".*Integration.*" -v

    # Run any remaining tests
    echo "🧪 Running any remaining tests..."
    go test ./...

    # Clean up
    echo "🧹 Cleaning up test database..."
    docker-compose -f docker-compose.test.yml down -v

else
    echo "❌ Docker not found. Please install Docker or set up PostgreSQL manually."
    echo "See README.md for manual test database setup instructions."
    exit 1
fi

echo "✅ Tests completed successfully!"