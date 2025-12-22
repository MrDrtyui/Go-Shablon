# Go Backend Template

A production-ready Go backend template with JWT authentication, refresh tokens, and PostgreSQL. Built following Go best practices with clean architecture and idiomatic Go patterns.

## Features

- RESTful API with Chi router v5
- JWT access tokens + secure refresh tokens
- Single-use refresh tokens with automatic rotation
- PostgreSQL 16 with Ent ORM
- Clean architecture (Handler → Service → Repository)
- Docker support for local development
- Comprehensive middleware stack
- Bcrypt password hashing
- Configuration management with YAML
- Crypto-secure token generation (crypto/rand, SHA-256)

## Project Structure

```
backend/
├── cmd/api/                  # Application entry point
├── domain/                   # Domain models and DTOs
├── internal/
│   ├── app/                  # Application initialization
│   ├── auth/                 # JWT authentication
│   ├── config/               # Configuration management
│   ├── db/                   # Database connection
│   ├── middleware/           # HTTP middleware
│   ├── refreshtoken/         # Refresh token logic
│   │   ├── token.go          # Token generation & hashing
│   │   ├── repo.go           # Data access layer
│   │   └── service.go        # Business logic
│   ├── router/               # HTTP routing
│   └── user/                 # User domain logic
│       ├── handler.go        # HTTP handlers
│       ├── service.go        # Business logic
│       ├── repo.go           # Data access
│       └── mapper.go         # DTO mappers
├── ent/                      # Ent ORM generated code
│   └── schema/               # Entity schemas
│       ├── user.go           # User entity
│       └── refreshtoken.go   # RefreshToken entity
└── config/                   # Configuration files
```

## Quick Start

### 1. Start PostgreSQL

```bash
cd docker
docker compose up -d
```

This starts PostgreSQL 16 on `localhost:5432` with credentials from `.env.postgres`.

### 2. Configure the application

Create `config/dev.yaml`:

```yaml
env: dev
http:
  port: ":9000"
database:
  url: "postgres://user:pass@localhost:5432/postgres?sslmode=disable"
jwt:
  secret: "your-secret-key-change-in-production"
  accessTtlHours: 1h        # Access token TTL (1 hour)
  refreshTtlHours: 720h     # Refresh token TTL (30 days)
```

### 3. Run the application

```bash
cd backend
CONFIG_PATH=../config/dev.yaml go run cmd/api/main.go
```

The server will start on `http://localhost:9000`

### 4. Build for production

```bash
cd backend
make build
./bin/app
```

## API Endpoints

### Authentication Endpoints

All authentication endpoints return both `access_token` and `refresh_token`.

#### 1. Register New User

Creates a new user account and returns authentication tokens.

```bash
curl -X POST http://localhost:9000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123",
    "username": "john"
  }'
```

**Response (201 Created):**
```json
{
  "user": {
    "id": 1,
    "email": "john@example.com",
    "username": "john"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "Mz8F2QxN5vR8kL3pY7dW1jH6nC9mB4sT2qA0zX..."
}
```

#### 2. Login

Authenticates existing user and returns tokens.

```bash
curl -X POST http://localhost:9000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

**Response (200 OK):**
```json
{
  "user": {
    "id": 1,
    "email": "john@example.com",
    "username": "john"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "vWzZ3fyYsTERwP678-iPrfDqiALxWJBBNCAumqw..."
}
```

#### 3. Refresh Tokens

Rotates refresh token and generates new access token. **Single-use**: the old refresh token is automatically revoked.

```bash
curl -X POST http://localhost:9000/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "vWzZ3fyYsTERwP678-iPrfDqiALxWJBBNCAumqw..."
  }'
```

**Response (200 OK):**
```json
{
  "user": {
    "id": 1,
    "email": "john@example.com",
    "username": "john"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "oTQmq2Y2sbSQjhQF7HMbULXPOvHj7KjCY9an410..."
}
```

**Note:** The new `refresh_token` is different from the old one. The old token cannot be reused.

#### 4. Logout

Revokes the refresh token, effectively logging out the user.

```bash
curl -X POST http://localhost:9000/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "oTQmq2Y2sbSQjhQF7HMbULXPOvHj7KjCY9an410..."
  }'
```

**Response:** `204 No Content`

### Error Responses

All endpoints return JSON error responses:

```json
{
  "error": "error message here"
}
```

**Common HTTP Status Codes:**
- `200` - Success
- `201` - Created (register)
- `204` - No Content (logout)
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid credentials/tokens)
- `500` - Internal Server Error

## Testing Scenarios

### Scenario 1: Complete Authentication Flow

```bash
# 1. Register
RESPONSE=$(curl -s -X POST http://localhost:9000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","username":"tester"}')
echo "$RESPONSE"

# Extract refresh token
REFRESH_TOKEN=$(echo "$RESPONSE" | python3 -c "import json,sys; print(json.load(sys.stdin)['refresh_token'])")

# 2. Refresh tokens
curl -X POST http://localhost:9000/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"

# 3. Logout
curl -X POST http://localhost:9000/auth/logout \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```

### Scenario 2: Error Handling

```bash
# Invalid credentials
curl -X POST http://localhost:9000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"wrongpassword"}'
# Response: {"error": "invalid email or password"}

# Missing required fields
curl -X POST http://localhost:9000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"password":"test123"}'
# Response: {"error": "email and password are required"}

# Invalid refresh token
curl -X POST http://localhost:9000/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"invalid-token"}'
# Response: {"error": "invalid or expired refresh token"}

# Duplicate email
curl -X POST http://localhost:9000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"test","username":"duplicate"}'
# Response: {"error": "...duplicate key..."}
```

### Scenario 3: Token Rotation Security

```bash
# 1. Login
RESPONSE=$(curl -s -X POST http://localhost:9000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"securepass123"}')
OLD_TOKEN=$(echo "$RESPONSE" | python3 -c "import json,sys; print(json.load(sys.stdin)['refresh_token'])")

# 2. Use refresh token (rotates it)
curl -X POST http://localhost:9000/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$OLD_TOKEN\"}"

# 3. Try to reuse old token (should fail)
curl -X POST http://localhost:9000/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$OLD_TOKEN\"}"
# Response: {"error": "invalid or expired refresh token"}
```

## Security Features

### Refresh Token Implementation

The refresh token system follows security best practices:

1. **Crypto-Secure Generation**: Uses `crypto/rand` for token generation (32 random bytes)
2. **SHA-256 Hashing**: Tokens are hashed before database storage (never stored in plaintext)
3. **Single-Use with Rotation**: Each refresh generates new tokens and revokes the old one
4. **Expiration**: Long-lived (30 days default) but explicitly revocable
5. **Logout Support**: Users can explicitly revoke tokens
6. **Database-Backed**: Tokens can be revoked server-side (no JWT refresh tokens)

### Token Flow

```
┌─────────────┐
│   Register  │ ──► Access Token (1h) + Refresh Token (30d)
│   / Login   │
└─────────────┘

┌─────────────┐
│   Refresh   │ ──► New Access Token + New Refresh Token
│             │     (Old refresh token is REVOKED)
└─────────────┘

┌─────────────┐
│   Logout    │ ──► Refresh Token REVOKED
└─────────────┘
```

### Password Security

- **Bcrypt Hashing**: Passwords hashed with `bcrypt.DefaultCost`
- **Never Exposed**: Password field excluded from all API responses
- **Secure Comparison**: Constant-time comparison via bcrypt

## Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────┐
│          HTTP Handlers              │  ← Presentation Layer
│  (Parse requests, format responses) │
└─────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────┐
│       Business Services             │  ← Business Logic Layer
│  (Validation, domain operations)    │
└─────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────┐
│         Repositories                │  ← Data Access Layer
│   (Database queries, Ent ORM)       │
└─────────────────────────────────────┘
```

### Design Patterns

- **Clean Architecture**: Separation of concerns with clear boundaries
- **Repository Pattern**: Interface-based data access for testability
- **Dependency Injection**: Loose coupling via constructor injection
- **DTO Pattern**: Domain models for request/response transformation
- **Middleware Pipeline**: Composable request processing (Chi)
- **Service Layer**: Encapsulated business logic

### Database Schema

**User Entity:**
- `id` (auto-increment)
- `email` (unique, required)
- `username` (optional)
- `password` (bcrypt hash, required)

**RefreshToken Entity:**
- `id` (auto-increment)
- `token_hash` (SHA-256, unique, indexed)
- `user_id` (foreign key to User)
- `expires_at` (timestamp)
- `created_at` (timestamp, immutable)
- `revoked` (boolean, default false)

**Relationship:** User `has many` RefreshTokens

## Development

### Generate Ent Code

After modifying schema files:

```bash
cd backend
go generate ./ent
```

### Build

```bash
cd backend
make build
```

Binary will be created at `backend/bin/app`.

### Run Tests

```bash
cd backend
make test
```

### Lint Code

```bash
cd backend
make lint
```

Requires `golangci-lint` to be installed.

### Database Migrations

Ent automatically creates/migrates schema on application startup. For production, consider using Ent's migration tools.

## Configuration

All configuration is in YAML format. Create environment-specific configs:

- `config/dev.yaml` - Development
- `config/prod.yaml` - Production
- `config/test.yaml` - Testing

**Example Production Config:**

```yaml
env: prod
http:
  port: ":8080"
database:
  url: "postgres://produser:strongpass@db.example.com:5432/proddb?sslmode=require"
jwt:
  secret: "very-strong-secret-from-env-var"
  accessTtlHours: 1h
  refreshTtlHours: 720h
```

**Environment Variable:**

```bash
export CONFIG_PATH=/path/to/config/prod.yaml
```

## Technologies

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.25+ | Programming language |
| Chi | v5.2.3 | HTTP router |
| Ent | v0.14.5 | ORM and schema management |
| PostgreSQL | 16 | Relational database |
| JWT | v5.3.0 | Access token generation |
| Bcrypt | - | Password hashing |
| Docker | - | Local database containerization |

## Middleware Stack

The application includes the following middleware (executed in order):

1. **Logger** - HTTP request logging
2. **Recoverer** - Panic recovery and error handling
3. **RequestID** - Unique request ID injection
4. **RealIP** - Client IP extraction from headers

## Best Practices Applied

- ✅ Standard library preference (`crypto/rand`, `crypto/sha256`)
- ✅ Interface-based repository pattern for testability
- ✅ Explicit error handling with custom error types
- ✅ Separation of concerns (handler/service/repository)
- ✅ Idiomatic Go naming and structure
- ✅ No global state or singletons
- ✅ Context propagation for request scoping
- ✅ Structured logging
- ✅ Graceful database connection handling
- ✅ Security-first approach (hashing, rotation, revocation)

## License

This is a template project. Use it as you see fit.

## Contributing

This is a personal template. Feel free to fork and customize for your needs.
