# Fiber Hexagonal

REST API built with **Go Fiber** and **Hexagonal Architecture** (Ports & Adapters)

A clean, scalable, and maintainable REST API implementation following the hexagonal architecture pattern. This project demonstrates best practices for organizing Go code with clear separation of concerns.

## Project Structure

```
fiber/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── adapters/            # External interfaces (HTTP, Database, Cache)
│   │   ├── config/          # Configuration management
│   │   ├── routes/          # HTTP routes and handlers
│   │   ├── storage/         # Database and cache adapters
│   │   │   ├── cache/       # Redis cache
│   │   │   ├── orm/         # GORM database
│   │   │   └── repository/  # Data access layer
│   ├── core/                # Business logic (isolated from frameworks)
│   │   ├── domain/          # Entity definitions
│   │   ├── ports/           # Interfaces (contracts)
│   │   │   ├── input/       # Input ports (use cases)
│   │   │   └── output/      # Output ports (repositories)
│   │   └── services/        # Business services
│   └── pkg/                 # Shared utilities
│       └── query/           # Query building and filtering
├── docker/                  # Docker configuration
├── script/                  # Utility scripts
├── go.mod                   # Go module definition
└── README.md               # This file
```

## Architecture

This project implements **Hexagonal Architecture** (also known as Ports & Adapters):

- **Core**: Contains business logic and domain models (independent of frameworks)
- **Ports**: Define interfaces for communication between layers
  - **Input Ports**: Use cases/services
  - **Output Ports**: Repository interfaces
- **Adapters**: Implement concrete functionality
  - **Incoming Adapters**: HTTP handlers
  - **Outgoing Adapters**: Database, cache, external services

## Features

- ✅ Clean hexagonal architecture
- ✅ RESTful API with Fiber framework
- ✅ GORM ORM for database operations
- ✅ Redis caching support
- ✅ Query filtering and pagination
- ✅ Docker support
- ✅ Environment-based configuration

## Getting Started

### Prerequisites

- Go 1.18+
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
| `_limit` | Number of records per page | `_limit=10` |
| `_offset` | Number of records to skip | `_offset=20` |
| `_page` | Page number (alternative to offset) | `_page=3` |
| `_sort` | Column to sort by | `_sort=name` |
| `_order` | Sort direction: ASC or DESC | `_order=DESC` |

### Filter Operators

| Operator | Suffix | SQL | Example |
|----------|--------|-----|---------|
| Equals | `_eq` | `=` | `name_eq=John` |
| Not Equals | `_ne` | `!=` | `status_ne=inactive` |
| Greater Than | `_gt` | `>` | `age_gt=18` |
| Greater Than or Equal | `_gte` | `>=` | `age_gte=18` |
| Less Than | `_lt` | `<` | `age_lt=65` |
| Less Than or Equal | `_lte` | `<=` | `price_lte=100` |
| Contains | `_contains` | `LIKE` | `name_contains=John` |
| In | `_in` | `IN` | `status_in=active,pending` |
| Not In | `_nin` | `NOT IN` | `status_nin=deleted,archived` |
| Is Null | `_null` | `IS NULL` | `deleted_at_null=true` |
| Is Not Null | `_notnull` | `IS NOT NULL` | `deleted_at_notnull=true` |

### Example Queries

#### Basic Pagination
```bash
# Get first 10 users
GET /api/v1/users?_limit=10

# Get users on page 2 (10 per page)
GET /api/v1/users?_limit=10&_page=2

# Get users with offset
GET /api/v1/users?_limit=10&_offset=20
```

#### Filtering
```bash
# Get active users
GET /api/v1/users?status_eq=active

# Get users older than 18
GET /api/v1/users?age_gt=18

# Get users with names containing "John"
GET /api/v1/users?name_contains=John

# Get users with status either 'active' or 'pending'
GET /api/v1/users?status_in=active,pending

# Get users without a deleted_at timestamp
GET /api/v1/users?deleted_at_null=true
```

#### Sorting
```bash
# Sort by name ascending
GET /api/v1/users?_sort=name&_order=ASC

# Sort by created date descending
GET /api/v1/users?_sort=created_at&_order=DESC
```

#### Combined Queries
```bash
# Get active users aged 18-65, sorted by name, page 1
GET /api/v1/users?status_eq=active&age_gte=18&age_lte=65&_sort=name&_order=ASC&_limit=20&_page=1

# Search books with price less than 50, sorted by price descending
GET /api/v1/books?price_lt=50&_sort=price&_order=DESC&_limit=10

# Get non-deleted users with names containing "smith"
GET /api/v1/users?name_contains=smith&deleted_at_null=true&_limit=20
```

### Response Format

**Success Response:**
```json
{
  "code": 200,
  "message": "OK",
  "error": null,
  "payload": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
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
  "error": "Invalid filter parameter",
  "payload": null
}
```

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

## License

This project is licensed under the MIT License.
