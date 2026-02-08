# BE Modular Golang

[![Go Version](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A scalable backend API built with **Go** following **clean architecture** and **modular design principles**. This project demonstrates enterprise-level backend development with microservices-ready structure.

## 🏗️ Architecture

This project follows a **modular monolith architecture** with clear separation of concerns:

```
Backend/
├── cmd/                    # Application entry points
│   └── api/               # Main API server
├── internal/              # Private application code
│   ├── common/           # Shared infrastructure
│   │   ├── auth/        # Authentication & authorization
│   │   ├── configs/     # Configuration management
│   │   ├── database/    # Database connections
│   │   ├── events/      # Event handling
│   │   ├── infrastructure/ # Infrastructure services
│   │   │   ├── kafka/   # Kafka integration
│   │   │   ├── redis/   # Redis client
│   │   │   └── ...
│   │   ├── logs/        # Logging utilities
│   │   └── server/      # HTTP server setup
│   └── modules/         # Business domain modules
│       └── [module]/    # Each module has its own structure
│           ├── delivery/      # HTTP handlers (controllers)
│           ├── domain/        # Business entities & interfaces
│           ├── infrastructure/ # Data access layer
│           └── usecase/       # Business logic
├── pkg/                  # Public libraries
├── configs/             # Configuration files
├── migrations/          # Database migrations
├── deployments/         # Docker & deployment configs
├── docs/               # API documentation (Swagger)
└── api/                # Protocol buffers definitions
```

## ✨ Features

- 🏛️ **Clean Architecture**: Clear separation between layers
- 🔌 **Modular Design**: Easy to add/remove business modules
- 🔐 **Authentication**: JWT & OAuth2 with Keycloak integration
- 📊 **Multiple Databases**: PostgreSQL, Redis, Cassandra, Neo4j
- 📨 **Event-Driven**: Apache Kafka for async messaging
- 🔍 **Full-Text Search**: Elasticsearch integration
- 🔄 **Dependency Injection**: Using Google Wire
- 📝 **API Documentation**: Auto-generated Swagger/OpenAPI docs
- 🐳 **Docker Support**: Ready for containerization
- 🔥 **Live Reload**: Air for development hot-reload
- 📦 **gRPC Support**: Protocol Buffers integration

## 🛠️ Tech Stack

### Core
- **Language**: Go 1.25.5
- **Framework**: Gin Web Framework
- **DI**: Google Wire

### Databases
- **PostgreSQL**: Primary relational database
- **Redis**: Caching & session storage
- **Cassandra**: Wide-column store
- **Neo4j**: Graph database
- **Elasticsearch**: Search & analytics

### Message Queue
- **Apache Kafka**: Event streaming platform

### Authentication
- **Keycloak**: Identity and access management
- **JWT**: JSON Web Tokens
- **OAuth2/OIDC**: Standard protocols

### DevOps
- **Docker**: Containerization
- **Docker Compose**: Local development
- **Air**: Live reload for development
- **Makefile**: Build automation

## 🚀 Getting Started

### Prerequisites

- Go 1.25.5 or higher
- Docker & Docker Compose
- Make (optional)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/NguyenNhatquang522004/BE-modular-golang.git
cd BE-modular-golang
```

2. **Copy environment variables**
```bash
cp configs/config.env.example configs/config.env
# Edit configs/config.env with your settings
```

3. **Start dependencies with Docker**
```bash
docker-compose -f deployments/docker-compose.yml up -d
```

4. **Install dependencies**
```bash
go mod download
```

5. **Generate Wire dependencies**
```bash
cd cmd/api
wire
# or using make
make wire
```

6. **Run database migrations**
```bash
make migrate-up
# or manually
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up
```

7. **Run the application**
```bash
# Development mode with live reload
make dev
# or
air

# Production mode
make run
# or
go run cmd/api/main.go
```

The API will be available at `http://localhost:8081`

## 📋 Available Commands

```bash
# Run application
make run              # Run in production mode
make dev              # Run with hot-reload (Air)

# Build
make build            # Build binary

# Database migrations
make migrate-up       # Apply all migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration file

# Code generation
make wire             # Generate Wire dependency injection
make swagger          # Generate Swagger documentation
make proto            # Generate protobuf files

# Testing
make test             # Run tests
make test-coverage    # Run tests with coverage

# Docker
make docker-build     # Build Docker image
make docker-up        # Start all services
make docker-down      # Stop all services

# Code quality
make lint             # Run linter
make fmt              # Format code
```

## 📚 API Documentation

After starting the server, visit:
- **Swagger UI**: `http://localhost:8081/swagger/index.html`
- **API Docs**: Available in `/docs` directory

## 🧩 Module Structure

Each business module follows this structure:

```
module-name/
├── delivery/           # HTTP handlers
│   ├── http/          # REST endpoints
│   └── grpc/          # gRPC services (optional)
├── domain/            # Business entities & interfaces
│   ├── entity.go      # Domain models
│   └── repository.go  # Repository interfaces
├── infrastructure/    # External dependencies
│   ├── persistence/   # Database implementations
│   └── external/      # Third-party services
└── usecase/           # Business logic
    └── service.go     # Use case implementations
```

## 🔐 Environment Variables

Key environment variables in `configs/config.env`:

```env
# Server
APP_PORT=8081
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=dbname
DB_USER=user
DB_PASSWORD=password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Kafka
KAFKA_BROKERS=localhost:9092

# Keycloak
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=your-realm
KEYCLOAK_CLIENT_ID=your-client-id
KEYCLOAK_CLIENT_SECRET=your-secret

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific module
go test ./internal/modules/[module-name]/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go standard formatting (use `gofmt`)
- Write meaningful commit messages
- Add tests for new features
- Update documentation as needed

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**NguyenNhatquang522004**

- GitHub: [@NguyenNhatquang522004](https://github.com/NguyenNhatquang522004)

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [Google Wire](https://github.com/google/wire)
- [Swaggo](https://github.com/swaggo/swag)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- Clean Architecture principles by Robert C. Martin

## 📞 Support

For questions or issues, please:
- Open an issue on GitHub
- Contact: [Your Email]

---

⭐ If you find this project helpful, please give it a star!
