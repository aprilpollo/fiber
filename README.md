# Aprilpollo

A high-performance REST API built with Go, utilizing the Fiber web framework and clean architecture principles.

## 🚀 Features

- **Fast & Efficient**: Built on top of [Fiber](https://gofiber.io/) (Express-inspired web framework)
- **Clean Architecture**: Hexagonal/Ports & Adapters architecture for maintainability
- **PostgreSQL Database**: Using GORM ORM for database operations
- **JWT Authentication**: Secure token-based authentication
- **Docker Support**: Full containerization with Docker Compose
- **Request Validation**: Input validation using go-playground/validator
- **Rotating Logs**: Automatic log rotation with configurable retention
- **Health Checks**: Built-in health and version endpoints

## 📋 Prerequisites

- Go 1.24 or higher
- PostgreSQL 15
- Docker and Docker Compose (optional)

## 🏗️ Project Structure

```
.
├── cmd/                        # Application entrypoints
│   └── main.go                # Main application
├── internal/                   # Private application code
│   ├── adapter/               # Adapters for external services
│   │   ├── config/           # Configuration management
│   │   ├── handler/          # HTTP handlers (Fiber)
│   │   │   ├── fiber/        # Fiber-specific handlers
│   │   │   │   ├── middleware/  # Middleware components
│   │   │   │   └── routes/      # Route definitions
│   │   │   └── http_req/     # HTTP request utilities
│   │   └── storage/          # Database adapters
│   │       └── gorm/         # GORM implementation
│   │           ├── models/   # Database models
│   │           ├── repository/  # Data access layer
│   │           └── views/    # Database views
│   ├── core/                  # Business logic layer
│   │   ├── domain/           # Domain entities
│   │   ├── port/             # Interface definitions
│   │   └── service/          # Business services
│   └── util/                  # Utility functions
├── docker/                    # Docker configuration
│   ├── docker-compose.yml    # Service orchestration
│   └── Dockerfile            # Application container
├── scripts/                   # Database migrations & scripts
├── static/                    # Static files
├── logs/                      # Application logs
└── deploy.sh                 # Deployment script
```

## 🔧 Configuration

Create a `.env` file in the root directory:

```env
# Application
APP_NAME=aprilpollo
APP_VERSION=1.0.0
APP_MODE=development
API_PORT=8760
API_SHUTDOWN_TIMEOUT_SECONDS=30
LOG_LEVEL=info
ALLOWED_CREDENTIAL_ORIGINS=*

# JWT Configuration
JWT_SECRET_KEY=your-secret-key-here
JWT_EXPIRE_DAYS_COUNT=7
JWT_ISSUER=aprilpollo
JWT_SUBJECT=api-access
JWT_SIGNING_METHOD=HS256

# PostgreSQL Database
POSTGRE_URI=postgresql://aprilpollo:apl9921@localhost:5432/aprilpollo?sslmode=disable
POSTGRE_MAX_IDLE_CONNS=10
POSTGRE_MAX_OPEN_CONNS=100
POSTGRE_CONN_MAX_LIFETIME=0
```

## 🚦 Getting Started

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Fiber
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run PostgreSQL** (if not using Docker)
   ```bash
   # Using your local PostgreSQL installation
   createdb aprilpollo
   ```

5. **Run the application**
   ```bash
   go run cmd/main.go
   ```

### Docker Deployment

1. **Build and run with Docker Compose**
   ```bash
   cd docker
   docker-compose up -d
   ```

2. **View logs**
   ```bash
   docker-compose logs -f app
   ```

3. **Stop services**
   ```bash
   docker-compose down
   ```

## 📡 API Endpoints

### Health & Info
- `GET /health` - Health check endpoint
- `GET /version` - Get API version

### Authentication
- Authentication endpoints are defined in the auth routes module

## 🧪 Testing

```bash
go test ./...
```

## 📦 Build

### Local Build
```bash
go build -o bin/app cmd/main.go
```

### Docker Build
```bash
docker build -f docker/Dockerfile -t aprilpollo:latest .
```

## 🛠️ Tech Stack

- **[Fiber](https://gofiber.io/)** - Web framework
- **[GORM](https://gorm.io/)** - ORM library
- **[PostgreSQL](https://www.postgresql.org/)** - Database
- **[JWT](https://jwt.io/)** - Authentication
- **[Validator](https://github.com/go-playground/validator)** - Input validation
- **[godotenv](https://github.com/joho/godotenv)** - Environment configuration
- **[File Rotatelogs](https://github.com/lestrrat-go/file-rotatelogs)** - Log rotation

## 📝 Architecture

This project follows the **Hexagonal Architecture** (Ports and Adapters) pattern:

- **Core Layer**: Contains business logic, domain entities, and port interfaces
- **Adapter Layer**: Implements the ports with concrete adapters (HTTP handlers, database repositories)
- **Clean Separation**: Business logic is independent of frameworks and external services

### Dependency Flow
```
Handler → Service → Repository → Database
   ↓         ↓          ↓
Adapter   Core      Adapter
```

## 🔒 Security

- JWT-based authentication
- Environment-based configuration
- Non-root Docker user
- CORS support via middleware
- Input validation on all endpoints

## 📊 Logging

Logs are automatically rotated and stored in the `logs/` directory:
- Current log: `app.log`
- Archived logs: `app.YYYY-MM-DD.log`

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

[MIT License](https://mit-license.org)

## 👤 Author

**Aprilpollo**

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) for the amazing web framework
- [GORM](https://gorm.io/) for the powerful ORM
- The Go community for excellent libraries and tools
