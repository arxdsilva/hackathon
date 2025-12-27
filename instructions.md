# Hackathon Management Microservice

A Go-based web microservice built with [Buffalo](https://gobuffalo.io/) for scheduling and managing hackathons. This service provides APIs to create, manage, and track hackathon events.

## Features

- **Event Management**: Create, update, and delete hackathon events
- **Participant Registration**: Manage participant sign-ups and team formations
- **Schedule Management**: Create and manage event schedules and timelines
- **Team Management**: Organize participants into teams
- **RESTful API**: Clean, RESTful API endpoints for all operations

## Prerequisites

- Go 1.16+ (check with `go version`)
- Buffalo CLI (`buffalo new` command available)
- PostgreSQL for database (configurable)
- Git

## Getting Started

### 1. Install Dependencies

```bash
go mod download
```

### 2. Set Up the Database

Configure your database connection in the `config/database.yml` file.

To create and migrate the database:

```bash
buffalo db create
buffalo db migrate
```

### 3. Run the Development Server

```bash
buffalo dev
```

The server will start on `http://localhost:3000` by default.

## Project Structure

```
.
├── actions/          # HTTP handlers and route logic
├── models/           # Database models
├── migrations/       # Database migrations
├── public/           # Static files
├── templates/        # HTML templates (if applicable)
├── config/           # Configuration files
├── docker-compose.yml # Docker setup (if included)
└── main.go          # Application entry point
```

## API Endpoints

### Hackathons
- `GET /hackathons` - List all hackathons
- `POST /hackathons` - Create a new hackathon
- `GET /hackathons/{id}` - Get hackathon details
- `PUT /hackathons/{id}` - Update a hackathon
- `DELETE /hackathons/{id}` - Delete a hackathon (owner only)

### Participants
- `GET /hackathons/{id}/participants` - List participants
- `POST /hackathons/{id}/participants` - Register participant
- `DELETE /hackathons/{id}/participants/{pid}` - Remove participant

### Teams
- `GET /hackathons/{id}/teams` - List teams
- `POST /hackathons/{id}/teams` - Create team
- `PUT /hackathons/{id}/teams/{tid}` - Update team

## Environment Configuration

Create a `.env` file in the project root with the following variables:

```
GO_ENV=development
PORT=3000
DATABASE_URL=postgres://user:password@localhost:5432/hackathon_db
LOG_LEVEL=debug
```

## Building for Production

```bash
buffalo build
```

This creates an executable binary in the `bin/` directory.

## Testing

Run tests with:

```bash
buffalo test
```

For specific test coverage:

```bash
go test ./... -cover
```

## Docker Deployment

If a `Dockerfile` is included:

```bash
docker build -t hackathon-service .
docker run -p 3000:3000 hackathon-service
```

Or using Docker Compose:

```bash
docker-compose up
```

## Development Workflow

1. Create a new feature branch: `git checkout -b feature/your-feature`
2. Make changes and test locally with `buffalo dev`
3. Commit with clear messages: `git commit -m "Add feature description"`
4. Push and create a pull request

## Common Commands

| Command | Purpose |
|---------|---------|
| `buffalo generate model` | Create a new data model |
| `buffalo generate action` | Create a new route handler |
| `buffalo db migrate` | Run pending migrations |
| `buffalo routes` | Display all routes |
| `buffalo dev` | Start development server with live reload |

## Troubleshooting

### Database Connection Issues
- Verify PostgreSQL/SQLite is running
- Check `DATABASE_URL` in `.env`
- Run `buffalo db create` to initialize the database

### Port Already in Use
- Change PORT in `.env` or pass `--port` flag:
  ```bash
  buffalo dev --port 3001
  ```

### Migration Errors
- Check migration files in `migrations/`
- Reset database: `buffalo db drop` then `buffalo db create`

## Contributing

1. Follow Go code style guidelines ([Effective Go](https://golang.org/doc/effective_go))
2. Write tests for new features
3. Update this README if adding new functionality
4. Keep commits atomic and well-described

## Additional Resources

- [Buffalo Documentation](https://gobuffalo.io/)
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## License

[Add your license here]