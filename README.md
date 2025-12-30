# Hackathon Management Platform

[![CI](https://github.com/arxdsilva/hackathon/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/arxdsilva/hackathon/actions/workflows/unit-tests.yml)
[![codecov](https://codecov.io/gh/arxdsilva/hackathon/branch/main/graph/badge.svg)](https://codecov.io/gh/arxdsilva/hackathon)

This is a comprehensive web application built with Buffalo for managing hackathons. It provides a complete solution for creating, organizing, and running hackathon events with advanced features including audit logging, secure admin controls, and responsive design. The platform supports user registration, project submissions, team formation, file management, and detailed scheduling, all wrapped in a modern dark theme interface.

## Features

- User registration and authentication
- Hackathon creation and management
- Project submissions and team formation
- File upload and management
- Schedule management
- Admin dashboard with comprehensive user management
- Dark theme UI with sticky navigation
- Docker containerization
- **Comprehensive audit logging** - All user actions are logged with timestamps, IP addresses, and user agents
- **Owner account protection** - Prevents deletion of owner/admin accounts for security
- **Responsive layout** - CSS Grid-based admin panels with sticky sidebar navigation
- **Custom branding** - Professional favicon and visual identity
- **Team member display** - Project pages show all team members with their roles and join dates
- **Presentation management** - Projects can opt-in to presentations with toggle functionality and admin oversight
- **Admin presentations dashboard** - Comprehensive view of all presenting projects across hackathons with presentation order tracking

## Security Features

- **Audit Logging**: Complete tracking of all user actions with detailed logs including timestamps, IP addresses, and user agents
- **Owner Protection**: Owner/admin accounts cannot be deleted, preventing accidental lockouts
- **CSRF Protection**: Built-in Cross-Site Request Forgery protection
- **Session Management**: Secure session handling with proper authentication

## New Features (December 28, 2025)

### Team Member Display
Project pages now display comprehensive team information including:
- **Team member list** with names, emails, and roles (owner/member)
- **Join timestamps** showing when each member joined the project
- **Visual indicators** distinguishing project owners from regular members
- **Responsive design** that works well on all screen sizes

### Presentation Management System
A complete presentation workflow for hackathon projects featuring:
- **Presentation toggle** - Project owners can opt their projects in/out of presentations
- **Presentation status tracking** - Clear visual indicators showing presentation eligibility
- **Audit logging** - All presentation status changes are logged for transparency
- **Presentation order** - Automatic timestamp-based ordering for presentation scheduling

### Admin Presentations Dashboard
Enhanced admin panel with dedicated presentations oversight:
- **Presenting projects statistics** - Count of projects opting to present across all hackathons
- **Comprehensive project listing** - All presenting projects with detailed information
- **Presentation order display** - Projects ordered by when they opted to present
- **Cross-hackathon visibility** - Admin can see presentations from all hackathons they manage
- **Project details** - Names, descriptions, repository links, and team information
- **Visual project cards** - Clean card-based layout with project images and metadata

### Hackathon Presentation Display
Hackathon overview pages now feature:
- **Presenting projects section** - Dedicated area showing all projects that will present
- **Presentation order** - Projects displayed in the order they opted to present
- **Project thumbnails** - Visual representation with project images
- **Team member counts** - Quick overview of team sizes
- **Direct links** - Easy navigation to individual project pages

## Screenshots

Here are some screenshots of the hackathon management platform:

### Logged Homepage
![Logged Homepage](docs/screenshots/logged_homepage.png)

### Hackathon Overview
![Hackathon Overview](docs/screenshots/hackathon_overview.png)

### Admin Panel
![Admin Panel](docs/screenshots/admin_panel.png)

*Note: The above screenshots show the platform before the latest feature additions (team member display, presentation management, and enhanced admin dashboard). The current version includes these new features as documented above.*

## Quick Start

### Prerequisites

- **Go** (version 1.25 or later)
- **Node.js** (version 18 or later)
- **Yarn** (version 1.x)
- **PostgreSQL** (version 15 or later)
- **Docker** and **Docker Compose** (for containerized deployment)

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/arxdsilva/hackathon.git
cd hackathon

# Complete project setup (recommended)
make setup

# Or set up manually:
make install          # Install dependencies
make db-docker-up     # Start PostgreSQL with Docker
make db-setup         # Set up database
make db-migrate       # Run migrations
make assets-dev       # Build assets
```

### Development

```bash
# Start development server
make dev

# Run unit tests (fast, no database)
make test-unit

# Run full test suite
make test-full

# View all available commands
make help
```

### Makefile Commands

This project includes a comprehensive Makefile with common development commands:

```bash
# Development
make dev              # Start development server
make build            # Build for production
make run              # Run production build

# Testing
make test-unit        # Run unit tests (fast, no DB)
make test             # Run all tests

# Database
make db-setup         # Set up development database
make db-migrate       # Run migrations
make db-reset         # Reset database
make db-docker-up     # Start DB with Docker

# Assets
make assets-dev       # Build assets with watch mode
make assets-build     # Build production assets

# Code Quality
make fmt              # Format Go code
make vet              # Run go vet
make lint             # Run linter

# Docker
make docker-build     # Build Docker image
make docker-dev       # Start full dev environment

# Utilities
make status           # Show environment status
make clean            # Clean build artifacts
make help             # Show all commands
```

Install Go dependencies:
```bash
go mod download
```

Install Node.js dependencies:
```bash
yarn install
```

### 3. Database Setup

Create a PostgreSQL database for the application. You can use the provided Docker setup or set up PostgreSQL locally.

#### Option A: Using Docker (Recommended)

```bash
# Start PostgreSQL in Docker
docker run --name hackathon-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -p 5432:5432 -d postgres:15-alpine

# Initialize the database with the schema
docker exec -i hackathon-postgres psql -U postgres -d postgres < init-db.sql
```

#### Option B: Local PostgreSQL

Create a database named `hackathon_development` and run the initialization script:

```sql
CREATE DATABASE hackathon_development;
\c hackathon_development;
\i init-db.sql;
```

### 4. Environment Configuration

The application uses the following environment variables (defaults are provided for development):

- `GO_ENV=development`
- `DATABASE_URL=postgres://postgres:postgres@localhost:5432/hackathon_development?sslmode=disable`
- `PORT=3000`
- `LOG_LEVEL=debug`

You can override these by creating a `.env` file or setting them in your shell.

### 5. Run Database Migrations

```bash
# Run migrations to set up the database schema
buffalo db migrate up
```

### 6. Start the Development Server

```bash
buffalo dev
```

The application will be available at [http://127.0.0.1:3000](http://127.0.0.1:3000).

## Docker Deployment

For production deployment or isolated development environment, use Docker Compose:

### Build and Run with Docker Compose

```bash
# Build the application
docker-compose build

# Start all services (PostgreSQL + App)
docker-compose up -d

# View logs
docker-compose logs -f app
```

The application will be available at [http://localhost:3000](http://localhost:3000).

### Docker Commands

```bash
# Stop services
docker-compose down

# Rebuild after code changes
docker-compose build --no-cache

# View logs for specific service
docker-compose logs postgres
docker-compose logs app

# Access database directly
docker-compose exec postgres psql -U postgres -d hackathon_development
```

## Available Commands

### Buffalo Commands

- `buffalo dev` - Start development server with hot reload
- `buffalo build` - Build the application binary
- `buffalo db migrate up` - Run database migrations
- `buffalo db migrate down` - Rollback migrations
- `buffalo db migrate status` - Check migration status
- `buffalo routes` - List all application routes

### Asset Management

- `yarn build` - Build production assets
- `yarn dev` - Watch and rebuild assets during development

## Project Structure

```
├── actions/              # Buffalo actions (controllers)
├── models/               # Database models with audit logging
├── templates/            # Plush templates with responsive layouts
├── assets/               # CSS, JS, and image assets
├── migrations/           # Database migrations including audit_logs
├── public/               # Static files including custom favicon
├── docs/                 # Documentation and screenshots
│   └── screenshots/      # Application screenshots
├── grifts/               # Buffalo tasks
├── config/               # Application configuration
├── docker-compose.yml    # Docker services configuration
└── init-db.sql          # Database initialization script
```

## Testing

The application includes comprehensive unit tests to ensure code quality and functionality.

### Test Types

- **Unit Tests**: Pure business logic tests that don't require database connections
- **Mock-Based Tests**: Repository interface tests using generated mocks
- **Validation Tests**: Business logic validation without external dependencies

### Test Coverage

[![codecov](https://codecov.io/gh/arxdsilva/hackathon/branch/main/graph/badge.svg)](https://codecov.io/gh/arxdsilva/hackathon)

### Running Unit Tests (No Database Required)

Unit tests test business logic in isolation and can be run without any database setup:

```bash
# Run all unit tests (no database required)
go test ./actions -v

# Run tests with coverage
go test ./actions -cover

# Generate coverage report
go test ./actions -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run all tests
go test ./...
```

### Test Organization

Tests are organized by functionality in corresponding `<filename>_test.go` files:

- `auth_test.go` - Authentication tests
- `users_test.go` - User management tests
- `profile_test.go` - Profile display tests
- `admin_test.go` - Admin functionality tests
- `files_test.go` - File management tests
- `pages_test.go` - Public pages tests
- `project_memberships_test.go` - Project membership tests
- `company_configurations_unit_test.go` - Configuration validation tests
- `repository_mock_example_test.go` - Mock usage examples

### Test Execution

Unit tests run completely independently without any external dependencies and provide fast feedback during development.

## CI/CD

This project uses GitHub Actions for continuous integration. The CI pipeline automatically runs on every push and pull request to the `main` and `develop` branches.

### GitHub Actions Workflows

- **Unit Tests**: Runs isolated unit tests that don't require a database connection

### Local Testing Options

- **Unit Tests**: Run instantly without any setup (`go test ./actions -run "Test..."`)
- **CI**: Handles unit tests automatically on every push/PR

### Workflow Triggers

The CI pipeline runs automatically when:
- Code is pushed to `main` or `develop` branches
- Pull requests are opened against `main` or `develop` branches

### Local Testing vs CI

- **Unit tests** run the same locally and in CI (no external dependencies)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `buffalo test`
5. Submit a pull request

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.

<!-- presentation: resource -->