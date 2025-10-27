# mini-trh-backend

## Project Structure

```
mini-trh-backend/
├── internal/                          # Internal shared utilities
│   ├── consts/
│   │   └── consts.go                 # Application constants
│   ├── logger/
│   │   └── logger.go                 # Zap-based logging setup
│   └── utils/
│       ├── hash.go                   # Password hashing utilities
│       └── response.go               # HTTP response helpers
│
└── pkg/                               # Core business logic
    ├── domain/
    │   └── entities/                  # Domain models
    │       ├── deployment.go         # Deployment entity
    │       ├── integration.go        # Integration entity
    │       ├── node.go               # Node entity
    │       └── user.go               # User entity
    │
    ├── infrastructure/
    │   └── postgres/                  # Database layer
    │       ├── database.go           # Database connection
    │       ├── repositories/         # Data access layer
    │       │   ├── deployment_repository.go
    │       │   ├── integration_repository.go
    │       │   ├── node_repository.go
    │       │   ├── repository.go     # Repository factory
    │       │   └── user_repository.go
    │       └── schemas/              # Database schemas
    │           ├── deployment_schema.go
    │           ├── integration_schema.go
    │           ├── migrate.go        # Migration runner
    │           ├── node_schema.go
    │           └── user_schema.go
    │
    └── services/                      # Business logic layer
        ├── auth_service.go           # Authentication service
        ├── integration_service.go    # Integration management
        ├── jwt_service.go            # JWT token service
        ├── node_deployment_service.go # Node deployment logic
        └── taskmanager/
            └── task_manager.go       # Async task queue
```

## Architecture

The project follows a **layered architecture** with clear separation of concerns:

```
┌─────────────────────────────────────┐
│   API Layer (To be implemented)     │  ← HTTP handlers & routes
├─────────────────────────────────────┤
│   Service Layer                     │  ← Business logic
│   - Authentication & JWT            │
│   - Node deployment                 │
│   - Integration management          │
│   - Async task processing           │
├─────────────────────────────────────┤
│   Repository Layer                  │  ← Data access
│   - CRUD operations                 │
│   - Database queries                │
├─────────────────────────────────────┤
│   Domain Layer                      │  ← Data models
│   - Entities & domain logic         │
├─────────────────────────────────────┤
│   Infrastructure Layer              │  ← External dependencies
│   - PostgreSQL connection           │
│   - GORM ORM                        │
└─────────────────────────────────────┘
```

### Key Design Patterns

- **Repository Pattern**: Abstracts data access logic
- **Dependency Injection**: Services receive dependencies via constructors
- **Domain-Driven Design**: Clear domain entities with business logic
- **Async Processing**: Worker pool-based task queue for long-running operations

## Directory Details

### `/internal` - Internal Utilities

Private packages that can only be used within this project.

- **consts/**: Application-wide constants (user roles, status values)
- **logger/**: Structured logging using Uber Zap
- **utils/**: Common utilities (password hashing, HTTP responses)

### `/pkg/domain/entities` - Domain Layer

Business domain models representing core concepts:

- **User**: Authentication and authorization
- **Node**: Blockchain node configuration and state
- **Deployment**: Node deployment process tracking
- **Integration**: Additional features (monitoring, backup, explorer)

### `/pkg/infrastructure/postgres` - Data Access Layer

PostgreSQL database interaction:

- **database.go**: Connection initialization and pooling
- **repositories/**: CRUD operations for each entity
- **schemas/**: GORM table schemas and migrations

### `/pkg/services` - Business Logic Layer

Core application logic:

- **AuthService**: User registration and login
- **JWTService**: Token generation and validation
- **NodeDeploymentService**: Node creation and deployment orchestration
- **IntegrationService**: Integration installation and management
- **TaskManager**: Asynchronous task queue with worker pool

## Tech Stack

- **Language**: Go 1.x
- **Web Framework**: Gin (planned)
- **ORM**: GORM
- **Database**: PostgreSQL with JSONB support
- **Authentication**: JWT with bcrypt password hashing
- **Logging**: Uber Zap
- **Other**: Google UUID, golang.org/x/crypto

## Database Schema

```
users (1) ──────→ (N) nodes
                    ├─ (1) ──────→ (N) deployments
                    └─ (1) ──────→ (N) integrations
```

### Tables

- **users**: User authentication and roles
- **nodes**: Blockchain nodes with JSONB config
- **deployments**: Multi-step deployment tracking
- **integrations**: Additional features per node

## Getting Started

### Prerequisites

- Go 1.x
- PostgreSQL 12+

### Environment Variables

Create a `.env` file based on `.env.example`:

```env
PORT=8000

POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=trh_thomas_db
POSTGRES_HOST=localhost
POSTGRES_PORT=5433

JWT_SECRET=your-super-secret-jwt-key-change-in-production

DEFAULT_ADMIN_EMAIL=admin@thomas.com
DEFAULT_ADMIN_PASSWORD=admin123
```

### Installation

```bash
# Install dependencies
go mod download

# Run database migrations
# (Migrations run automatically on service start)

# Run the application
go run cmd/main.go
```

## Current Implementation Status

- ✅ Domain entities
- ✅ Repository layer with PostgreSQL
- ✅ Service layer with business logic
- ✅ Authentication and JWT
- ✅ Async task processing
- ⏳ HTTP API handlers (planned)

## License

[Add license information]
