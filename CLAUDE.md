# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based gym management system server built with gofiber framework. It provides a REST API for managing gym-related data with MariaDB/MySQL database backend.

## Development Commands

### Building and Running
- `make run` - Start the server in development mode (uses gin for hot reload)
- `make server` - Build the production binary to `bin/gym`
- `make test` - Run all tests
- `make clean` - Remove built binaries

### File Watching
- `make fswatch` - Watch controller files for changes
- `make allrun` - Start file watching and server together

### Docker Commands
- `make dockerbuild` - Build Linux binary for Docker
- `make docker` - Build Docker image (tag with `tag=` variable)
- `make dockerrun` - Run container
- `make push` - Push Docker image to registry

### Cross-Platform Building
- `make linux` - Build Linux binary

## Architecture

### Core Structure
- **main.go** - Application entry point, initializes services
- **router/** - HTTP routing and endpoint definitions
- **controllers/** - Request handlers with base Controller struct
- **models/** - Database models and connection management  
- **services/** - Background services (HTTP, cron, notifications)
- **global/** - Shared utilities (config, logging, JWT, time)

### Key Components

#### Router System
- **router/router.go** - Auto-generated REST endpoints for all models
- **router/auth.go** - JWT authentication middleware
- Routes follow pattern: `/api/{model}` with standard CRUD operations
- All API routes except `/api/jwt` require JWT authentication

#### Controller Pattern
- **controllers/controllers.go** - Base Controller struct with common functionality
- REST controllers in **controllers/rest/** - One per model
- Controllers handle Init/Close lifecycle with database connections
- Built-in pagination, file uploads, and JSON response formatting

#### Database Layer
- **models/db.go** - Connection management with retry logic
- Supports MySQL, PostgreSQL, and SQL Server
- Transaction support with isolation levels
- Connection pooling configured (100 max, 10 idle, 5min lifetime)

#### Configuration
- **global/config/config.go** - Environment-based config management
- Supports `.env.yml` files for different environments
- Environment variables override config file settings
- Database connection strings auto-generated based on type

### Services
- **services/http.go** - Main HTTP server with fiber
- **services/cron.go** - Scheduled tasks
- **services/chat.go** - WebSocket chat functionality  
- **services/notify.go** - Push notifications
- **services/fcm.go** - Firebase Cloud Messaging

### Build Tools
The project uses custom build tools (referenced in Makefile):
- `buildtool-model` - Generates model files
- `buildtool-router` - Auto-generates router.go from controllers

### Authentication
- JWT-based authentication using golang-jwt/jwt/v5
- JWT middleware in router/auth.go
- User sessions stored in Controller.Session

### Logging
- Structured logging with zerolog
- File rotation with custom lumberjack integration
- Configurable log levels and outputs (console, file, database)

### File Structure Pattern
Models follow a dual-file pattern:
- `models/{model}.go` - Main model struct
- `models/{model}/{model}.go` - Model-specific logic and methods

### Dependencies
Key external dependencies:
- gofiber/fiber/v2 - HTTP framework
- go-sql-driver/mysql - MySQL driver  
- lib/pq - PostgreSQL driver
- golang-jwt/jwt/v5 - JWT handling
- rs/zerolog - Structured logging
- spf13/viper - Configuration management

## Testing

Run tests with `make test` which executes `go test -v ./...`

## Environment Variables

Key environment variables for configuration:
- `APP_MODE` - Set to "production" or "develop"
- `PORT` - Server port
- `DB_TYPE`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASS` - Database configuration
- `LOG_LEVEL`, `LOG_CONSOLE`, `LOG_WEB`, `LOG_DATABASE`, `LOG_FILE` - Logging configuration
- `TLS_USE`, `TLS_CERT`, `TLS_KEY` - TLS configuration