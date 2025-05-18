# User Management System

A user management system built with Go, featuring data import capabilities using WorkerPool cuncurrency pattern and RESTful API endpoint to get user along their addresses.

## Features

- Bulk user data import with concurrent processing
- User and address management
- RESTful API endpoints
- PostgreSQL database integration
- Docker support for easy deployment
- Configurable worker pool for data processing
- Comprehensive test coverage (unit, integration, and handler tests)
- Automated test database setup

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- PostgreSQL 14 (if running without Docker)

## Project Structure

```
.
├── api/            # API handlers and routes (presenters would be written here too)
├── cmd/            # Application entry points
├── config/         # Configuration management
├── internal/       # Internal packages (user, address operations, validations specific to domain can be placed here, we can define domain models and seperate them from entities in this folder)
├── pkg/            # Shared packages
│   ├── load/      # Data loading utilities
│   └── storage/   # Database operations
├── service/        # Business logic layer
├── test/          # Test files and utilities
│   ├── integration/  # Integration tests
│   ├── mocks/       # Generated mocks
│   └── config.yaml  # Test configuration
└── data/          # Data files
```

## Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd sika
   ```

2. Start the PostgreSQL database using Docker:
   ```bash
   docker-compose up -d
   ```

3. Configure the application:
   - Copy `config.yaml.example` to `config.yaml`
   - Update the configuration as needed

4. Run the application:
   ```bash
   go run cmd/api/main.go --config config.yaml --file data/users_data.json

   #or

   make run
   ```

## Testing

The project includes comprehensive test coverage:

### Unit Tests
- Service layer tests with mocked dependencies
- Worker pool concurrency tests
- Error handling and edge cases

### Integration Tests
- Database operations
- End-to-end user import flow
- Concurrent processing with real database
- Data consistency checks

### Handler Tests
- API endpoint validation
- Request/response handling
- Error scenarios
- Status code verification

Run tests using:
```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Clean up test database
make clean-test-db
```

## Data Import

The application supports bulk data import from JSON files(`data/users_data.json`). The import process:
- Uses a worker pool for concurrent processing
- Handles user and address data
- Supports batch operations on create addresses for better performance
- Includes error handling and logging

To import data:
```bash
go run cmd/api/main.go -file path/to/users_data.json
```

To force re-import data:
1. Delete the `.data_imported` file
2. Run the application again

## API Endpoints

- `GET /users/:id` - Get user by ID
- More endpoints to be documented...

## Configuration

The application can be configured using `config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: user
  pass: pass
  dbname: sika-db
```

## Future Improvements

1. API Enhancements:
   - Add pagination for user listing
   - Implement filtering and sorting

2. Performance Optimizations:
   - Add caching layer (Redis)

3. Monitoring and Observability:
   - Implement structured logging using logrus and sentry or loki

4. Security Enhancements:
   - Add authentication
   - Implement rate limiting
