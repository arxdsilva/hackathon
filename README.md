# Hackathon Management Platform

[![CI](https://github.com/arxdsilva/hackathon/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/arxdsilva/hackathon/actions/workflows/unit-tests.yml)
[![codecov](https://codecov.io/gh/arxdsilva/hackathon/branch/main/graph/badge.svg)](https://codecov.io/gh/arxdsilva/hackathon)

A production-ready platform for running internal company hackathons. Built with Go and Buffalo framework for high performance and reliability.

## Why This Platform?

**🚀 Launch Internal Hackathons in Minutes**
Create and manage hackathon events for your organization with a single deployment. No complex configuration required - just set up, invite employees, and start innovating.

**👥 Built for Company Teams**
Seamless team formation across departments with automatic membership tracking. Employees can create projects, join teams, and collaborate with built-in presentation management.

**🔒 Enterprise-Grade Security**
Complete audit logging tracks every action. Role-based access control, CSRF protection, and bcrypt password hashing with per-user salts keep your internal data secure.

**📊 Powerful Admin Dashboard**
Monitor everything from a central dashboard. Manage employees, track projects, view audit logs, and configure platform settings - all from one place.

**🎯 Self-Hosted & Open Source**
Keep your hackathon data within your organization. Deploy on your own infrastructure with full control over your intellectual property and project submissions.

**💻 Developer Friendly**
Clean repository pattern, comprehensive Makefile with 25+ commands, hot reload for development, and extensive documentation make customization straightforward.

## Table of Contents

- [Features](#features)
  - [User Management & Authentication](#user-management--authentication)
  - [Hackathon Management](#hackathon-management)
  - [Project & Team Management](#project--team-management)
  - [File Management](#file-management)
  - [Admin Dashboard](#admin-dashboard)
  - [Security & Audit](#security--audit)
  - [User Interface & Experience](#user-interface--experience)
  - [Technical Features](#technical-features)
- [Screenshots](#screenshots)
- [Quick Start](#quick-start)
  - [Prerequisites](#prerequisites)
  - [Initial Setup](#initial-setup)
  - [Development](#development)
  - [Makefile Commands](#makefile-commands)
- [Docker Deployment](#docker-deployment)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [CI/CD](#cicd)
- [Contributing](#contributing)
- [License](#license)

## Features

### User Management & Authentication
- **User Registration** - Create accounts with email validation and password policy enforcement
- **Secure Authentication** - bcrypt password hashing with per-user salts and session management
- **User Profiles** - Personal profiles with name, email, company/team, and role information
- **Profile Editing** - Users can update their personal information and change passwords
- **Password Reset** - Forced password reset functionality for new accounts
- **Role-Based Access** - Two-tier system with Owner (admin) and Hacker (participant) roles

### Hackathon Management
- **Hackathon Creation** - Owners can create hackathons with title, description, dates, and status
- **Schedule Management** - Add rich-text schedules with HTML formatting support (nullable field support)
- **Status Tracking** - Four status types: upcoming, active, completed, and hidden
- **Hackathon Listing** - Browse all hackathons with filtering and pagination
- **Detailed Views** - Individual hackathon pages with statistics, timeline, and project listings
- **Owner Controls** - Edit and delete hackathons (owner-only)

### Project & Team Management
- **Project Creation** - Users can create one project per hackathon with name, description, and links
- **Project Images** - Upload and manage project images stored in database
- **Team Formation** - Project membership system allowing users to join teams
- **Join/Leave Projects** - Users can join projects and owners can manage memberships
- **Unique Constraints** - Prevents users from creating multiple projects per hackathon
- **Team Member Display** - Shows all team members with roles (owner/member) and join timestamps
- **Presentation Opt-In** - Projects can toggle presentation status with order tracking

### File Management
- **File Uploads** - Upload files associated with hackathons and projects
- **Database Storage** - Files stored as binary data in PostgreSQL with metadata
- **File Download** - Secure file retrieval with proper content-type headers
- **File Management** - View, list, and delete uploaded files
- **Access Control** - File operations restricted to authenticated users

### Admin Dashboard
- **Comprehensive Overview** - Statistics on users, hackathons, projects, and audit events
- **User Management** - View, create, edit, and delete user accounts
- **Role Assignment** - Promote/demote users between owner and hacker roles
- **Account Protection** - Prevents deletion of owner accounts to maintain system access
- **Force Password Reset** - Require users to change passwords on next login
- **Hackathon Overview** - Admin view of all hackathons across the platform
- **Project Monitoring** - View all projects with filtering and search
- **Audit Logs** - Complete activity log with user actions, timestamps, and IP addresses
- **Presentations Dashboard** - Track all projects opting to present with order management
- **System Configuration** - Password policies and platform settings management

### Security & Audit
- **Comprehensive Audit Logging** - All user actions logged with timestamps, IP addresses, and user agents
- **Password Policies** - Configurable minimum length, uppercase, numbers, and special character requirements
- **CSRF Protection** - Built-in Cross-Site Request Forgery protection on all forms
- **Session Security** - Secure session handling with proper authentication checks
- **Owner Protection** - Prevents accidental lockouts by protecting admin accounts
- **Input Validation** - Server-side validation for all forms and data entry
- **SQL Injection Protection** - Parameterized queries via Pop ORM

### User Interface & Experience
- **Responsive Design** - Mobile-friendly layout that works on all screen sizes
- **Modern UI** - Professional interface with sticky navigation and clean styling
- **Bootstrap 5** - Clean, consistent styling with Bootstrap components
- **Font Awesome Icons** - Rich iconography throughout the interface
- **Flash Messages** - User feedback for all actions (success, error, info)
- **Pagination** - Efficient browsing of large datasets
- **Custom Branding** - Professional favicon with "H" logo and custom styling
- **Inline Forms** - Seamless inline editing and interactions

### Technical Features
- **Buffalo Framework** - Built on Go Buffalo for high performance
- **PostgreSQL Database** - Robust relational database with UUID support
- **Pop ORM** - Database migrations and model management
- **Docker Support** - Full containerization with docker-compose
- **Asset Pipeline** - Webpack for CSS/JS bundling and optimization
- **RESTful API** - Clean URL structure and HTTP methods
- **Repository Pattern** - Clean architecture with repository interfaces
- **Comprehensive Makefile** - Developer-friendly commands for all operations

## Screenshots

Here are some screenshots of the hackathon management platform:

### Home Page
![Home Page](docs/screenshots/home_page.png)
*Landing page with authentication options and platform overview*

### Hackathons Listing
![Hackathons Page](docs/screenshots/hackathons_page.png)
*Browse all available hackathons with status indicators and filtering*

### Hackathon Details
![Hackathon Page](docs/screenshots/hackathon_page.png)
*Individual hackathon view with projects, schedule, and team information*

### Admin Dashboard
![Admin Dashboard](docs/screenshots/admin_dashboard.png)
*Comprehensive admin panel with user management, audit logs, and system configuration*

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