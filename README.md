# UShort - URL Shortening Service

A high-performance URL shortening service built with Go, featuring rate limiting, authentication, and monitoring capabilities.

## Features

- URL shortening with base62 encoding
- JWT-based authentication
- Rate limiting using token bucket algorithm
- Redis caching for improved performance
- PostgreSQL database with optimized indexing
- Swagger API documentation
- Monitoring with Grafana
- Docker multi-container setup
- CI/CD pipeline ready

## Tech Stack

- Backend: Go with Gin framework
- Database: PostgreSQL
- Cache: Redis
- Container: Docker
- Monitoring: Grafana
- API Documentation: Swagger
- Testing: Go testing framework

## Project Structure

```
.
├── cmd/            # Application entry points
├── internal/       # Private application code
├── pkg/           # Public library code
├── docker/        # Docker configuration
├── docs/          # Documentation
├── scripts/       # Build and deployment scripts
└── tests/         # Test files
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- PostgreSQL
- Redis

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/ushort.git
cd ushort
```

2. Copy environment file and configure
```bash
cp .env.example .env
```

3. Run with Docker Compose
```bash
docker-compose up -d
```

### API Documentation

Once the service is running, access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## Development

### Running Tests
```bash
make test
```

### Building
```bash
make build
```

## License

MIT License 