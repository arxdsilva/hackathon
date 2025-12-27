# Hackathon Management Platform

This is a web application built with Buffalo for managing hackathons. It allows users to create and participate in hackathons, manage projects, teams, and schedules.

## Features

- User registration and authentication
- Hackathon creation and management
- Project submissions and team formation
- File upload and management
- Schedule management
- Admin dashboard
- Dark theme UI
- Docker containerization

## Prerequisites

Before running this application, make sure you have the following installed:

- **Go** (version 1.23 or later)
- **Node.js** (version 18 or later)
- **Yarn** (version 1.x)
- **PostgreSQL** (version 15 or later)
- **Docker** and **Docker Compose** (for containerized deployment)

## Local Development Setup

### 1. Clone the Repository

```bash
git clone https://github.com/arxdsilva/hackathon.git
cd hackathon
```

### 2. Install Dependencies

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
├── actions/          # Buffalo actions (controllers)
├── models/           # Database models
├── templates/        # Plush templates
├── assets/           # CSS, JS, and image assets
├── migrations/       # Database migrations
├── public/           # Static files
├── grifts/           # Buffalo tasks
├── config/           # Application configuration
└── docker-compose.yml # Docker services configuration
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `buffalo test`
5. Submit a pull request

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.
