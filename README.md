# URL Shortener Service

A production-ready URL shortening microservice built with Go, featuring RESTful API design, PostgreSQL integration, authentication middleware, and comprehensive logging for enterprise-grade applications.

## 🎯 Overview

This project provides a robust, scalable solution for URL shortening with enterprise-level features:
- **Microservice Architecture**: Clean, modular design following Go best practices
- **RESTful API**: Well-designed endpoints with proper HTTP semantics
- **Database Integration**: PostgreSQL with connection pooling and migrations
- **Authentication**: Basic authentication for secure URL management
- **Comprehensive Logging**: Structured logging with different environment configurations
- **Production Ready**: Configuration management, error handling, and monitoring

## 🛠️ Technology Stack

- **Language**: Go 1.24+
- **Web Framework**: Chi v5 - Lightweight, idiomatic HTTP router
- **Database**: PostgreSQL with connection pooling
- **Authentication**: Basic Auth middleware
- **Logging**: Structured logging with slog
- **Configuration**: Environment-based configuration with cleanenv
- **Validation**: Request validation with validator
- **Testing**: Comprehensive test suite with testify

## 🚀 Features

### Core Functionality
- **URL Shortening**: Generate short aliases for long URLs
- **URL Redirection**: Seamless redirect from short to original URL
- **URL Management**: Create and delete shortened URLs
- **Authentication**: Secure access to URL management endpoints
- **Database Persistence**: Reliable storage with PostgreSQL
- **Error Handling**: Comprehensive error management and logging

### API Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| `POST` | `/url` | Create shortened URL | Required |
| `DELETE` | `/url/{alias}` | Delete shortened URL | Required |
| `GET` | `/{alias}` | Redirect to original URL | None |

### Architecture Highlights
- **Clean Architecture**: Separation of concerns with layered design
- **Dependency Injection**: Testable and maintainable code structure
- **Middleware Pipeline**: Request ID, logging, recovery, and authentication
- **Configuration Management**: Environment-based configuration
- **Structured Logging**: JSON logging for production environments

## 📋 Prerequisites

- Go 1.24 or higher installed
- PostgreSQL server running
- Database created for the application
- Git for cloning the repository

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd url_shortener
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**
   ```sql
   CREATE DATABASE url_shortener;
   ```

4. **Run database migrations**
   ```bash
   # Apply migrations from migrations/ directory
   # Check migrations/ directory for available migration files
   ```

5. **Configure environment variables**
   Create a `.env` file in the project root:
   ```env
   # Database Configuration
   STORAGE_PATH=postgres://user:password@localhost/url_shortener?sslmode=disable
   
   # Server Configuration
   ADDRESS=:8080
   TIMEOUT=10s
   IDLE_TIMEOUT=60s
   
   # Authentication
   HTTP_SERVER_USER=admin
   HTTP_SERVER_PASSWORD=your_secure_password
   
   # Environment
   ENV=local
   ```

## 🎮 Usage

### Starting the Service

```bash
# Development mode
go run cmd/url_shortner/main.go

# Production build
go build -o url-shortener cmd/url_shortner/main.go
./url-shortener
```

### API Usage Examples

#### 1. Create a Shortened URL
```bash
curl -X POST http://localhost:8080/url \
  -H "Content-Type: application/json" \
  -u admin:your_secure_password \
  -d '{"url": "https://www.example.com/very/long/url"}'
```

**Response:**
```json
{
  "alias": "abc123",
  "url": "https://www.example.com/very/long/url"
}
```

#### 2. Redirect to Original URL
```bash
curl -L http://localhost:8080/abc123
```
This will redirect to the original URL.

#### 3. Delete a Shortened URL
```bash
curl -X DELETE http://localhost:8080/url/abc123 \
  -u admin:your_secure_password
```

**Response:**
```json
{
  "status": "ok"
}
```

## 🏗️ Project Structure

```
url_shortener/
├── cmd/
│   └── url_shortner/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── http-server/
│   │   └── handlers/
│   │       ├── delete.go       # Delete URL handler
│   │       ├── redirect.go     # Redirect handler
│   │       └── save.go         # Save URL handler
│   ├── lib/
│   │   └── sl.go               # Logging utilities
│   └── storage/
│       └── postgres/
│           └── postgres.go      # Database operations
├── migrations/                 # Database migration files
├── test/                       # Test files
├── config/                     # Configuration files
├── go.mod                      # Go module file
├── go.sum                      # Go module checksums
├── Taskfile.yaml               # Task definitions
├── .env                        # Environment variables (create this)
└── README.md                   # This file
```

## 🧩 Core Components

### Configuration Management
```go
type Config struct {
    Env         string        `yaml:"env" env:"ENV" env-default:"local"`
    StoragePath string        `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
    Address     string        `yaml:"address" env:"ADDRESS" env-default:":8080"`
    Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"10s"`
    IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
    User        string        `yaml:"user" env:"HTTP_SERVER_USER" env-required:"true"`
    Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD" env-required:"true"`
}
```

### Key Systems

1. **HTTP Server**
   - Chi router with middleware pipeline
   - Request ID generation and logging
   - Panic recovery and error handling
   - Graceful shutdown support

2. **Authentication**
   - Basic authentication middleware
   - Configurable user credentials
   - Secure credential storage

3. **Database Layer**
   - PostgreSQL connection pooling
   - Transaction management
   - Query optimization
   - Connection health checks

4. **Logging System**
   - Structured logging with slog
   - Environment-specific log formats
   - Request tracing and correlation IDs

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `STORAGE_PATH` | PostgreSQL connection string | Yes | - |
| `ADDRESS` | Server bind address | No | `:8080` |
| `TIMEOUT` | Request timeout | No | `10s` |
| `IDLE_TIMEOUT` | Idle timeout | No | `60s` |
| `HTTP_SERVER_USER` | Basic auth username | Yes | - |
| `HTTP_SERVER_PASSWORD` | Basic auth password | Yes | - |
| `ENV` | Environment (local/dev/prod) | No | `local` |

### Logging Configuration

- **Local**: Text format with debug level
- **Development**: JSON format with debug level
- **Production**: JSON format with info level

## 🎯 Performance Considerations

- **Connection Pooling**: Optimized database connections
- **Request Timeout**: Configurable timeouts for reliability
- **Idle Timeout**: Resource cleanup for idle connections
- **Structured Logging**: Efficient logging for production
- **Middleware Pipeline**: Optimized request processing

## 🧪 Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./test/...
```

### Test Structure
- Unit tests for individual components
- Integration tests for API endpoints
- Database testing with test containers
- Mock implementations for external dependencies

## 🚧 Development Notes

### Development Workflow
```bash
# Install task runner (if using Taskfile)
go install github.com/go-task/task/v3/cmd/task@latest

# Run development tasks
task dev          # Start development server
task test         # Run tests
task build        # Build application
task lint         # Run linter
```

### Extending the Service
- Add URL analytics and click tracking
- Implement custom alias generation
- Add rate limiting middleware
- Create admin dashboard
- Add URL expiration functionality
- Implement caching layer

### Production Considerations
- Set up monitoring and alerting
- Configure backup strategies
- Implement health check endpoints
- Add metrics collection
- Set up log aggregation
- Configure load balancing

### Security Enhancements
- Implement JWT authentication
- Add CORS configuration
- Implement rate limiting
- Add input sanitization
- Set up HTTPS/TLS
- Add security headers

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📞 Contact

For questions or support regarding this service, please open an issue in the repository.

---

**Built with ❤️ using Go and PostgreSQL**
