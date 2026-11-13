# Fiber Hexagonal

REST API built with **Go Fiber v3** and **Hexagonal Architecture** (Ports & Adapters)

A clean, scalable, and maintainable REST API implementation following the hexagonal architecture pattern. This project demonstrates best practices for organizing Go code with clear separation of concerns.

## Project Structure

```
fiber/
├── cmd/
│   └── main.go                  # Composition root: builds adapters, injects them into the core
├── internal/
│   ├── core/                    # The hexagon — stdlib + uuid only (enforced by core/arch_test.go)
│   │   ├── domain/              # Entities, invariants, domain errors, ListQuery
│   │   ├── ports/
│   │   │   ├── input/           # Use-case interfaces + commands (CreateUser, UpdateUser…)
│   │   │   └── output/          # Repository interfaces the core needs
│   │   └── services/            # Use cases implementing the input ports
│   └── adapters/                # Everything that talks to the outside world
│       ├── config/              # Environment configuration
│       ├── routes/              # HTTP input adapter (Fiber)
│       │   ├── server.go        # App, middleware, health probes
│       │   └── handler/         # Handlers, request/response DTOs, query-string parser, error → status
│       └── storage/             # Output adapters
│           ├── cache/           # Redis
│           ├── orm/             # GORM connection + persistence models
│           └── repository/      # Output-port implementations, SQL translation, error translation
├── docker/                      # Docker configuration
├── script/                      # Migration + seed
└── go.mod
```

## Architecture

This project implements **Hexagonal Architecture** (Ports & Adapters). Dependencies
point inward only: adapters import the core, the core imports nothing from adapters.

```
HTTP ──▶ handler (DTO → command) ──▶ input port ──▶ service ──▶ output port ──▶ repository ──▶ Postgres
                 ◀── status code ◀── domain error ◀────────────────── translate(gorm/pg error)
```

- **Domain** (`core/domain`): entities with no `json`/`gorm`/`validate` tags. Invariants
  (`NewUser`, `Validate`) live here, so every adapter gets the same rules. Failures are
  `ErrNotFound`, `ErrConflict` and `*ValidationError` — never a driver error.
- **Input ports** (`core/ports/input`): what the application can do, expressed in commands
  (`CreateUser`, `UpdateUser` with optional fields) and typed IDs (`uuid.UUID`).
- **Output ports** (`core/ports/output`): what the core needs from storage. Listing takes a
  `domain.ListQuery` with abstract operators (`OpContains`, `OpIn`…), not SQL.
- **Services** (`core/services`): orchestrate a use case — e.g. `Update` loads, applies only the
  given fields, validates, then saves. Testable with an in-memory fake repository.
- **HTTP adapter** (`adapters/routes`): binds requests into DTOs (the client can never set `id` or
  timestamps), parses the query string, and maps domain errors to status codes. Unexpected
  errors are logged and answered with a generic 500.
- **Storage adapter** (`adapters/storage`): maps entities ↔ GORM models, domain operators ↔ SQL,
  field names ↔ columns (an allowlist), and Postgres errors ↔ domain errors.

## Features

- ✅ Hexagonal architecture, with a test that fails if the core imports an adapter
- ✅ RESTful API with Fiber v3 framework
- ✅ GORM ORM for database operations
- ✅ Redis caching support
- ✅ Query filtering and pagination
- ✅ Request validation (`validate` tags via Fiber v3 `StructValidator`)
- ✅ Request ID, CORS and rate limiting middleware
- ✅ Health probes: `/livez` (liveness) and `/readyz` (DB + Redis readiness)
- ✅ Graceful shutdown
- ✅ Docker support
- ✅ Environment-based configuration

## Getting Started

### Prerequisites

- Go 1.26+
- Docker & Docker Compose (optional)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd fiber
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```

4. Run the application:
```bash
go run cmd/main.go
```

### Using Docker

```bash
docker-compose up -d
```

## Query Building and Filtering

The API supports advanced filtering, sorting, and pagination through query parameters.

### Reserved Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `_limit` | Records per page (1–200, default 20) | `_limit=10` |
| `_offset` | Number of records to skip | `_offset=20` |
| `_page` | Page number (overrides `_offset`) | `_page=3` |
| `_sort` | Column to sort by | `_sort=name` |
| `_order` | Sort direction: ASC or DESC | `_order=DESC` |

A malformed reserved parameter is rejected with `400 BAD_REQUEST` rather than
silently falling back to its default — e.g. `?_limit=abc` or `?_limit=999999999`.

Only fields the repository exposes can be filtered or sorted on — anything else is a
`400 BAD_REQUEST` (`password: cannot be filtered on`):

| Resource | Fields |
|----------|--------|
| `/api/v1/users` | `id`, `name`, `email`, `created_at`, `updated_at` |
| `/api/v1/books` | `id`, `title`, `author`, `created_at`, `updated_at` |

### Filter Operators

| Operator | Suffix | SQL | Example |
|----------|--------|-----|---------|
| Equals | `_eq` | `=` | `name_eq=John` |
| Not Equals | `_ne` | `!=` | `status_ne=inactive` |
| Greater Than | `_gt` | `>` | `age_gt=18` |
| Greater Than or Equal | `_gte` | `>=` | `age_gte=18` |
| Less Than | `_lt` | `<` | `age_lt=65` |
| Less Than or Equal | `_lte` | `<=` | `price_lte=100` |
| Contains | `_contains` | `LIKE` (wildcards in the value are matched literally) | `name_contains=John` |
| In | `_in` | `IN` | `status_in=active,pending` |
| Not In | `_nin` | `NOT IN` | `status_nin=deleted,archived` |
| Is Null | `_null` | `IS NULL` | `deleted_at_null=true` |
| Is Not Null | `_notnull` | `IS NOT NULL` | `deleted_at_notnull=true` |

### Example Queries

```bash
# Page 2 of users, 10 per page
GET /api/v1/users?_limit=10&_page=2

# Users whose name contains "smith", newest first
GET /api/v1/users?name_contains=smith&_sort=created_at&_order=DESC

# Users with one of these emails
GET /api/v1/users?email_in=alice@example.com,bob@example.com

# Books by an author, sorted by title
GET /api/v1/books?author_eq=Eric%20Evans&_sort=title
```

`pagination.total` counts the rows matching the filters, not the whole table.

### Response Format

**Success Response:**
```json
{
  "code": 200,
  "message": "OK",
  "error": null,
  "payload": [
    {
      "id": "3f1c2a9e-8d4b-4c1e-9f7a-2b6d5e8c1a40",
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 100,
    "count": 10,
    "page": 1,
    "limit": 10
  }
}
```

**Error Response:**
```json
{
  "code": 400,
  "message": "BAD_REQUEST",
  "error": "_limit: must be an integer, got \"abc\"",
  "payload": null
}
```

| Status | `message` | When |
|--------|-----------|------|
| 400 | `BAD_REQUEST` | Malformed body/params, failed validation, unknown filter field |
| 404 | `NOT_FOUND` | No record with that id (GET, PUT, DELETE) |
| 409 | `CONFLICT` | Unique value already taken (e.g. user email) |
| 500 | `INTERNAL_SERVER_ERROR` | Anything unexpected — details go to the log, never the response |

`PUT` updates only the fields present in the body.

## Development

### Running Tests
```bash
go test ./...
```

Repository tests run against a real Postgres when `TEST_POSTGRES_URI` is set (they drop and
recreate the `users`/`books` tables, so point it at a throwaway database):
```bash
TEST_POSTGRES_URI="postgres://aprilpollo:aprilpollo@localhost:5432/aprilpollo_test?sslmode=disable" go test ./internal/adapters/storage/repository
```

### Code Formatting
```bash
go fmt ./...
```

## License

This project is licensed under the MIT License.
