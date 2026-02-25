# Mikhmon v4 Go Backend

Clean Architecture implementation of Mikhmon hotspot management system backend in Go.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Mikhmon v4 Go Backend                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐ │
│  │   Handler    │  │   Handler    │  │   Handler    │  │     Handler      │ │
│  │    Auth      │  │   Router     │  │   Hotspot    │  │    Dashboard     │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └────────┬─────────┘ │
│         │                 │                 │                   │           │
│  ┌──────┴─────────────────┴─────────────────┴───────────────────┴─────────┐ │
│  │                              Use Case Layer                             │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐  │ │
│  │  │   Auth   │ │  Router  │ │ Hotspot  │ │ Voucher  │ │    Report    │  │ │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────────┘  │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                    │                                        │
│  ┌─────────────────────────────────┼─────────────────────────────────────┐ │
│  │                         Domain Layer                                   │ │
│  │  ┌──────────────┐  ┌───────────┴────────┐  ┌──────────────────────┐  │ │
│  │  │   Entity     │  │  Repository (Port)  │  │  Service Interface    │  │ │
│  │  │  - AdminUser │  │  - AdminUserRepo    │  │  - MikroTikClient     │  │ │
│  │  │  - Router    │  │  - RouterRepo       │  │  - HotspotOps         │  │ │
│  │  │  - Setting   │  │  - SettingRepo      │  │  - ReportOps          │  │ │
│  │  │  - Template  │  │  - TemplateRepo     │  │                       │  │ │
│  │  └──────────────┘  └────────────────────┘  └──────────────────────┘  │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                    │                                        │
│  ┌─────────────────────────────────┼─────────────────────────────────────┐ │
│  │                    Infrastructure Layer                                │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │ │
│  │  │PostgreSQL│ │  Redis   │ │MikroTik  │ │   JWT    │ │  Gin HTTP  │  │ │
│  │  │  (GORM)  │ │  Cache   │ │   API    │ │   Auth   │ │  Router    │  │ │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └────────────┘  │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Key Design Decisions

### 1. No Database for Hotspot Data
Unlike traditional approaches, this implementation stores ALL hotspot data (users, profiles, active sessions) in MikroTik itself via the API. PostgreSQL is only used for:
- `admin_users` - Mikhmon login credentials
- `routers` - Router connection configurations
- `settings` - Application settings
- `print_templates` - Voucher print templates

### 2. On-Login Script Generator
The most critical component is the on-login script generator (`internal/infrastructure/mikrotik/onlogin_generator.go`). It generates RouterOS scripts that:
- Calculate expiration dates using MikroTik scheduler
- Lock MAC addresses
- Lock to specific servers
- Record sales data to `/system/script` for reporting

### 3. Sales Reporting via MikroTik Scripts
Reports are stored as named scripts in MikroTik's `/system/script` with format:
```
Name: "date-|-time-|-user-|-price-|-ip-|-mac-|-validity-|-profile-|-comment"
Owner: "jan2024" (month/year for grouping)
Source: date
Comment: "mikhmon"
```

## Project Structure

```
cmd/
└── api/
    └── main.go                 # Application entry point

internal/
├── domain/
│   ├── entity/                 # Database entities (4 tables only)
│   │   ├── admin.go
│   │   ├── router.go
│   │   ├── setting.go
│   │   └── print_template.go
│   ├── dto/                    # Data Transfer Objects
│   │   ├── auth.go
│   │   ├── hotspot.go          # User, Profile, Active DTOs
│   │   ├── voucher.go
│   │   ├── report.go
│   │   └── dashboard.go
│   ├── repository/             # Repository interfaces
│   │   └── interfaces.go
│   └── service/                # Service interfaces
│       └── mikrotik_service.go
│
├── usecase/                    # Business logic
│   ├── auth_usecase.go
│   ├── router_usecase.go
│   ├── hotspot_usecase.go
│   ├── voucher_usecase.go
│   ├── report_usecase.go
│   └── dashboard_usecase.go
│
└── infrastructure/
    ├── auth/
    │   └── jwt.go              # JWT implementation
    ├── cache/
    │   └── redis.go            # Redis client
    ├── config/
    │   └── config.go           # Configuration management
    ├── database/
    │   └── postgres.go         # Database connection
    ├── http/
    │   ├── handler/            # HTTP handlers
    │   │   ├── auth_handler.go
    │   │   ├── router_handler.go
    │   │   ├── hotspot_handler.go
    │   │   ├── voucher_handler.go
    │   │   ├── report_handler.go
    │   │   └── dashboard_handler.go
    │   ├── middleware/         # Auth & CORS middleware
    │   │   ├── auth.go
    │   │   └── cors.go
    │   └── router.go           # Route setup
    ├── mikrotik/               # MikroTik API implementation
    │   ├── client.go           # API client wrapper
    │   ├── hotspot_users.go    # User operations
    │   ├── hotspot_profiles.go # Profile operations
    │   ├── hotspot_active.go   # Active sessions
    │   ├── reports.go          # Sales reporting
    │   ├── system.go           # System resources
    │   ├── onlogin_generator.go # Script generator
    │   └── voucher_generator.go # Voucher code generation
    └── repository/postgres/    # Repository implementations
        ├── admin_repository.go
        ├── router_repository.go
        ├── setting_repository.go
        └── template_repository.go
```

## API Endpoints

### Authentication
```
POST /api/v1/auth/login        # Login with username/password
GET  /api/v1/auth/me           # Get current user info
```

### Routers (Protected)
```
GET    /api/v1/routers              # List all routers
POST   /api/v1/routers              # Create router
GET    /api/v1/routers/:id          # Get router details
PUT    /api/v1/routers/:id          # Update router
DELETE /api/v1/routers/:id          # Delete router
POST   /api/v1/routers/:id/test     # Test connection
```

### Hotspot (Protected)
```
GET    /api/v1/hotspot/:router_id/users           # List users
POST   /api/v1/hotspot/:router_id/users           # Create user
GET    /api/v1/hotspot/:router_id/users/:id       # Get user
PUT    /api/v1/hotspot/:router_id/users/:id       # Update user
DELETE /api/v1/hotspot/:router_id/users/:id       # Delete user

GET    /api/v1/hotspot/:router_id/profiles        # List profiles
POST   /api/v1/hotspot/:router_id/profiles        # Create profile
PUT    /api/v1/hotspot/:router_id/profiles/:id    # Update profile
DELETE /api/v1/hotspot/:router_id/profiles/:id    # Delete profile

GET    /api/v1/hotspot/:router_id/active         # List active sessions
```

### Vouchers (Protected)
```
POST   /api/v1/vouchers/:router_id/generate       # Generate vouchers
GET    /api/v1/vouchers/:router_id?comment=xxx    # Get vouchers by comment
DELETE /api/v1/vouchers/:router_id?comment=xxx    # Delete vouchers
```

### Reports (Protected)
```
GET /api/v1/reports/:router_id/sales    # Get sales report
GET /api/v1/reports/:router_id/summary  # Get summary stats
GET /api/v1/reports/:router_id/export   # Export to CSV
```

### Dashboard (Protected)
```
GET /api/v1/dashboard/:router_id           # Dashboard data
GET /api/v1/dashboard/:router_id/resources # System resources
GET /api/v1/dashboard/:router_id/status    # Router status
```

## Configuration

Environment variables or config file:

```yaml
server:
  port: 8080
  environment: development

database:
  host: localhost
  port: 5432
  user: mikhmon
  password: mikhmon
  name: mikhmon
  ssl_mode: disable

redis:
  host: localhost
  port: 6379
  db: 0

jwt:
  secret: your-secret-key
  expiry: 24h
```

## Running the Application

### Development
```bash
# Start dependencies
docker-compose up -d postgres redis

# Run migrations
go run ./cmd/migrate/main.go

# Start server
go run ./cmd/api/main.go
```

### Production
```bash
# Build
go build -o bin/api ./cmd/api/main.go

# Run
./bin/api
```

### Docker
```bash
# Build image
docker build -t mikhmon-api:latest .

# Run with docker-compose
docker-compose up -d
```

## Development Status

### Completed ✅
- [x] Domain layer (entities, DTOs, interfaces)
- [x] Infrastructure layer (MikroTik API client, Redis, PostgreSQL)
- [x] On-login script generator (critical component)
- [x] Use cases (business logic)
- [x] HTTP handlers
- [x] Middleware (JWT auth, CORS)
- [x] Router setup
- [x] Docker support

### Remaining Tasks 🚧
- [ ] Fix type mismatches in use cases (string vs uint IDs)
- [ ] Add missing DTO types (VoucherGenerateRequest, VoucherBatchResult)
- [ ] Align repository interfaces with implementations
- [ ] Add comprehensive error handling
- [ ] Add logging middleware
- [ ] Add rate limiting
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Add Swagger documentation
- [ ] Add monitoring/metrics

## Key Files to Review

1. `internal/infrastructure/mikrotik/onlogin_generator.go` - The heart of Mikhmon logic
2. `internal/infrastructure/mikrotik/client.go` - MikroTik API wrapper
3. `internal/usecase/hotspot_usecase.go` - Business logic
4. `internal/domain/dto/hotspot.go` - Data structures

## Dependencies

- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/go-routeros/routeros/v3` - MikroTik API client
- `gorm.io/gorm` + `gorm.io/driver/postgres` - ORM
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `github.com/spf13/viper` - Configuration management
- `golang.org/x/crypto/bcrypt` - Password hashing

## License

MIT License - Same as original Mikhmon project
