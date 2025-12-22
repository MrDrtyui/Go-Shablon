# Go Template

A production-ready Go backend template with JWT authentication and PostgreSQL. Built following Go best practices with clean architecture.

## Features

- RESTful API with Chi router
- JWT authentication
- PostgreSQL with Ent ORM
- Clean architecture (Handler -> Service -> Repository)
- Docker support
- Middleware (Logger, Recoverer, RequestID, RealIP)
- Bcrypt password hashing
- Configuration management with YAML

## Project Structure

```
backend/
├── cmd/api/              # Application entry point
├── domain/               # Domain models and DTOs
├── internal/
│   ├── app/              # Application initialization
│   ├── auth/             # JWT authentication
│   ├── config/           # Configuration management
│   ├── db/               # Database connection
│   ├── middleware/       # HTTP middleware
│   ├── router/           # HTTP routing
│   └── user/             # User domain logic
│       ├── handler.go    # HTTP handlers
│       ├── service.go    # Business logic
│       ├── repo.go       # Data access
│       └── mapper.go     # DTO mappers
├── ent/                  # Ent ORM generated code
└── config/               # Configuration files
```

## Quick Start

### 1. Start PostgreSQL
```bash
cd docker
docker compose up -d
```

### 2. Run the application
```bash
cd backend
export CONFIG_PATH=./config/dev.yaml
make run
```

The server will start on `http://localhost:9000`

## API Endpoints

### Authentication

#### Register
```bash
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "username": "johndoe"
}

Response (201):
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}

Response (200):
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

## Configuration

Edit `backend/config/dev.yaml`:

```yaml
env: "dev"

http:
  port: ":9000"

database:
  url: "postgres://user:pass@localhost:5432/postgres?sslmode=disable"

jwt:
  secret: "your-secret-key-here"
  ttlHours: 24h
```

## Development

### Build
```bash
cd backend
make build
```

### Run tests
```bash
make test
```

### Lint
```bash
make lint
```

## Architecture

This project follows Go best practices:

- **Clean Architecture**: Separation of concerns with clear layers
- **Dependency Injection**: Components are loosely coupled
- **Repository Pattern**: Data access abstraction
- **Error Handling**: Proper error types and handling
- **Security**: Bcrypt password hashing, JWT tokens
- **Middleware**: Structured request processing pipeline

## Technologies

- **Chi v5**: Lightweight HTTP router
- **Ent**: Code-first ORM for Go
- **JWT**: Token-based authentication
- **PostgreSQL**: Relational database
- **Bcrypt**: Password hashing
- **Docker**: Containerization
