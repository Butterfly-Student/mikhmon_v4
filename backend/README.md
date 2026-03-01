# Mikhmon Backend API

Backend API untuk aplikasi Mikhmon (MikroTik Hotspot Monitor) - Sistem manajemen Hotspot dan PPP berbasis MikroTik.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-1.9+-cyan.svg)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-green.svg)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 📋 Daftar Isi

- [Fitur](#-fitur)
- [Teknologi](#-teknologi)
- [Arsitektur](#-arsitektur)
- [Struktur Project](#-struktur-project)
- [Persyaratan](#-persyaratan)
- [Instalasi](#-instalasi)
- [Konfigurasi](#-konfigurasi)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [API Endpoints](#-api-endpoints)
- [WebSocket](#-websocket)
- [Monitoring](#-monitoring)
- [Development](#-development)
- [Docker](#-docker)

## ✨ Fitur

### Autentikasi & Manajemen Admin
- JWT-based authentication
- Login/logout admin
- Manajemen session

### Manajemen Router
- CRUD router MikroTik
- Test koneksi ke router
- Multiple router support

### Hotspot Management
- **Users**: CRUD hotspot users, generate voucher
- **Profiles**: Manajemen profile hotspot (bandwidth limit, session time)
- **Active Users**: Monitoring user yang sedang aktif
- **Hosts**: Manajemen host hotspot
- **Servers**: Lihat konfigurasi hotspot server
- **Expire Monitor**: Setup script monitoring expired user

### PPP Management
- **Secrets**: CRUD PPP secrets (username/password)
- **Profiles**: Manajemen PPP profiles
- **Active**: Monitoring koneksi PPP aktif

### Voucher Management
- Generate voucher secara batch (hingga 500 voucher)
- Mode: Voucher Code (vc) atau Username/Password (up)
- Custom prefix, character set, panjang nama
- Time limit dan data limit
- Print template dengan HTML/CSS custom
- Cache vouchers

### Monitoring Real-time (WebSocket)
- **Resource Monitor**: CPU, Memory, Uptime, Disk usage
- **Traffic Monitor**: Traffic per interface real-time
- **Queue Monitor**: Simple queue monitoring
- **Ping Monitor**: Ping test ke host tertentu
- **Log Monitor**: Monitoring log MikroTik (Hotspot, PPP)
- **User Monitor**: Hotspot & PPP active/inactive users

### Reports
- Sales report berdasarkan owner atau hari
- Export report ke CSV
- Summary statistics

### Network Management
- **Interfaces**: List dan monitoring interface
- **NAT Rules**: Lihat konfigurasi NAT/Firewall
- **Queues**: Simple queue management
- **Address Pools**: IP pool management

### System Information
- Router resources
- Health check
- Router identity
- RouterBoard info
- System clock
- Dashboard data (summary)

## 🛠 Teknologi

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL 15+ dengan GORM
- **Cache**: Redis (opsional, untuk pub-sub)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **WebSocket**: Gorilla WebSocket
- **Logging**: Uber Zap
- **Configuration**: Viper
- **RouterOS API**: Custom client (pkg/routeros)

## 🏗 Arsitektur

Project menggunakan **Clean Architecture / Layered Architecture**:

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                        │
│              (HTTP Handlers, WebSocket, Middleware)          │
├─────────────────────────────────────────────────────────────┤
│                     Use Case Layer                           │
│              (Business Logic, Orchestration)                 │
├─────────────────────────────────────────────────────────────┤
│                   Domain Layer                               │
│         (Entities, DTOs, Repository Interfaces)              │
├─────────────────────────────────────────────────────────────┤
│                Infrastructure Layer                          │
│  (Database, RouterOS Client, Auth, Cache, External Services) │
└─────────────────────────────────────────────────────────────┘
```

### Alur Request

```
HTTP Request → Router → Handler → UseCase → Repository/Service → Database/MikroTik
                                                ↓
HTTP Response ← Handler ← UseCase ← Repository/Service ←
```

## 📁 Struktur Project

```
.
├── cmd/
│   ├── api/                    # Entry point aplikasi API
│   │   └── main.go
│   └── test/                   # Utility untuk testing
│       └── main.go
│
├── internal/
│   ├── domain/
│   │   ├── dto/                # Data Transfer Objects
│   │   │   ├── auth.go
│   │   │   ├── dashboard.go
│   │   │   ├── hotspot.go
│   │   │   ├── interface.go
│   │   │   ├── nat.go
│   │   │   ├── ping.go
│   │   │   ├── ppp.go
│   │   │   ├── queue.go
│   │   │   ├── report.go
│   │   │   ├── system.go
│   │   │   └── voucher.go
│   │   ├── entity/             # Domain entities
│   │   │   ├── admin.go        # Admin user
│   │   │   ├── router.go       # MikroTik router
│   │   │   └── setting.go      # App settings & print templates
│   │   ├── repository/         # Repository interfaces
│   │   │   └── interfaces.go
│   │   └── service/            # Service interfaces
│   │       └── mikrotik_service.go
│   │
│   ├── usecase/                # Business logic
│   │   ├── auth_usecase.go
│   │   ├── router_usecase.go
│   │   └── mikrotik/
│   │       ├── helpers.go
│   │       ├── hotspot_usecase.go
│   │       ├── interface_usecase.go
│   │       ├── log_usecase.go
│   │       ├── nat_usecase.go
│   │       ├── pool_usecase.go
│   │       ├── ppp_active_usecase.go
│   │       ├── ppp_profile_usecase.go
│   │       ├── ppp_secret_usecase.go
│   │       ├── queue_usecase.go
│   │       ├── report_usecase.go
│   │       ├── system_usecase.go
│   │       └── voucher_usecase.go
│   │
│   └── infrastructure/
│       ├── auth/               # JWT implementation
│       │   └── jwt.go
│       ├── cache/              # Redis cache
│       │   ├── interface.go
│       │   └── redis.go
│       ├── config/             # Configuration management
│       │   └── config.go
│       ├── database/           # PostgreSQL connection & migrations
│       │   ├── postgres.go
│       │   └── seed.go
│       ├── http/
│       │   ├── handler/        # HTTP handlers
│       │   │   ├── base.go
│       │   │   ├── auth_handler.go
│       │   │   ├── router_handler.go
│       │   │   └── mikrotik/
│       │   │       ├── hotspot_handler.go
│       │   │       ├── interface_handler.go
│       │   │       ├── log_handler.go
│       │   │       ├── nat_handler.go
│       │   │       ├── pool_handler.go
│       │   │       ├── ppp_active_handler.go
│       │   │       ├── ppp_profile_handler.go
│       │   │       ├── ppp_secret_handler.go
│       │   │       ├── queue_handler.go
│       │   │       ├── report_handler.go
│       │   │       ├── system_handler.go
│       │   │       ├── voucher_handler.go
│       │   │       └── ws/     # WebSocket handlers
│       │   │           ├── hotspot_active_monitor_handler.go
│       │   │           ├── hotspot_inactive_monitor_handler.go
│       │   │           ├── log_monitor_handler.go
│       │   │           ├── ping_handler.go
│       │   │           ├── ppp_active_monitor_handler.go
│       │   │           ├── ppp_inactive_monitor_handler.go
│       │   │           ├── queue_monitor_handler.go
│       │   │           ├── resource_monitor_handler.go
│       │   │           └── traffic_monitor_handler.go
│       │   ├── middleware/     # HTTP middleware
│       │   │   ├── auth.go
│       │   │   ├── cors.go
│       │   │   └── zap_logger.go
│       │   └── router.go       # Route definitions
│       ├── logger/             # Zap logger setup
│       │   └── logger.go
│       ├── mikrotik/           # RouterOS API implementation
│       │   ├── client.go       # RouterOS client
│       │   ├── manager.go      # Connection manager
│       │   ├── expire_monitor.go
│       │   ├── helpers.go
│       │   ├── hotspot_active.go
│       │   ├── hotspot_hosts.go
│       │   ├── hotspot_profiles.go
│       │   ├── hotspot_servers.go
│       │   ├── hotspot_service.go
│       │   ├── hotspot_users.go
│       │   ├── interfaces.go
│       │   ├── logging.go
│       │   ├── nat.go
│       │   ├── onlogin_generator.go
│       │   ├── ping.go
│       │   ├── pool.go
│       │   ├── ppp_active.go
│       │   ├── ppp_profile.go
│       │   ├── ppp_secret.go
│       │   ├── queue.go
│       │   ├── reports.go
│       │   ├── system.go
│       │   ├── voucher_generator.go
│       │   └── test/           # Integration tests
│       └── repository/postgres/# Repository implementations
│           ├── admin_repository.go
│           ├── router_repository.go
│           └── setting_repository.go
│
├── pkg/
│   ├── pubsub/                 # Redis pub-sub utility
│   │   └── pubsub.go
│   └── routeros/               # Custom RouterOS API client
│       ├── async.go
│       ├── chan_reply.go
│       ├── client.go
│       ├── client_test.go
│       ├── error.go
│       ├── listen.go
│       ├── logger.go
│       ├── proto/              # Protocol implementation
│       │   ├── io_context.go
│       │   ├── reader.go
│       │   ├── sentence.go
│       │   ├── writer.go
│       │   └── *_test.go
│       ├── reply.go
│       └── run.go
│
├── .env                        # Environment variables
├── .env.example                # Environment template
├── docker-compose.yml          # Docker services
├── Dockerfile                  # API container
├── Makefile                    # Build commands
├── go.mod
├── go.sum
└── README.md
```

## 📋 Persyaratan

- Go 1.21 atau lebih tinggi
- PostgreSQL 15 atau lebih tinggi
- Redis 7 (opsional, untuk pub-sub)
- MikroTik RouterOS dengan API enabled

## 🚀 Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/irhabi89/mikhmon.git
cd mikhmon
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Database dengan Docker

```bash
# Jalankan PostgreSQL dan Redis
docker-compose up -d postgres redis
```

### 4. Setup Environment Variables

```bash
cp .env.example .env
# Edit .env sesuai konfigurasi Anda
```

### 5. Jalankan Aplikasi

```bash
# Development
make run

# atau

go run ./cmd/api/main.go
```

Aplikasi akan berjalan di `http://localhost:8080`

## ⚙️ Konfigurasi

Konfigurasi dapat dilakukan melalui environment variables atau file `.env`:

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `SERVER_PORT` | 8080 | Port server HTTP |
| `SERVER_ENVIRONMENT` | development | Environment (development/production) |
| `DATABASE_HOST` | localhost | Host PostgreSQL |
| `DATABASE_PORT` | 5432 | Port PostgreSQL |
| `DATABASE_USER` | mikhmon | Username PostgreSQL |
| `DATABASE_PASSWORD` | mikhmon | Password PostgreSQL |
| `DATABASE_NAME` | mikhmon | Nama database |
| `DATABASE_SSLMODE` | disable | SSL mode (disable/require) |
| `REDIS_HOST` | localhost | Host Redis |
| `REDIS_PORT` | 6379 | Port Redis |
| `REDIS_PASSWORD` | - | Password Redis (opsional) |
| `REDIS_DB` | 0 | Database Redis |
| `JWT_SECRET` | your-secret-key | Secret key untuk JWT |
| `JWT_EXPIRY` | 24h | JWT expiry time |
| `INTERNAL_WS_KEY` | mikhmon-ws-internal-key | Key untuk WebSocket internal |

## ▶️ Menjalankan Aplikasi

### Development Mode

```bash
# Auto-reload dengan air (jika terinstall)
air

# Atau tanpa auto-reload
make run
```

### Production Mode

```bash
# Build binary
make build

# Run binary
./bin/api
```

### Docker

```bash
# Build image
make docker-build

# Run dengan docker-compose (database saja)
make docker-run

# Stop
make docker-stop
```

## 🔌 API Endpoints

### Autentikasi

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/auth/login` | Login admin |
| GET | `/api/v1/auth/me` | Get current user |
| POST | `/api/v1/auth/logout` | Logout |

### Router Management

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/routers` | Create router |
| GET | `/api/v1/routers` | List routers |
| GET | `/api/v1/routers/:id` | Get router detail |
| PUT | `/api/v1/routers/:id` | Update router |
| DELETE | `/api/v1/routers/:id` | Delete router |
| POST | `/api/v1/routers/test-connection` | Test MikroTik connection |

### Hotspot (per router)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/v1/mikrotik/:router_id/hotspot/active/count` | Count active users |
| GET | `/api/v1/mikrotik/:router_id/hotspot/profiles` | List profiles |
| POST | `/api/v1/mikrotik/:router_id/hotspot/profiles` | Create profile |
| GET | `/api/v1/mikrotik/:router_id/hotspot/profiles/:name` | Get profile by name |
| PUT | `/api/v1/mikrotik/:router_id/hotspot/profiles/:name` | Update profile |
| DELETE | `/api/v1/mikrotik/:router_id/hotspot/profiles/:name` | Delete profile |
| GET | `/api/v1/mikrotik/:router_id/hotspot/users` | List users |
| POST | `/api/v1/mikrotik/:router_id/hotspot/users` | Create user |
| GET | `/api/v1/mikrotik/:router_id/hotspot/users/:id` | Get user detail |
| PUT | `/api/v1/mikrotik/:router_id/hotspot/users/:id` | Update user |
| DELETE | `/api/v1/mikrotik/:router_id/hotspot/users/:id` | Delete user |
| GET | `/api/v1/mikrotik/:router_id/hotspot/active` | List active users |
| DELETE | `/api/v1/mikrotik/:router_id/hotspot/active/:id` | Disconnect active user |
| GET | `/api/v1/mikrotik/:router_id/hotspot/hosts` | List hosts |
| DELETE | `/api/v1/mikrotik/:router_id/hotspot/hosts/:id` | Delete host |
| GET | `/api/v1/mikrotik/:router_id/hotspot/servers` | List servers |
| POST | `/api/v1/mikrotik/:router_id/hotspot/expire-monitor` | Setup expire monitor |
| GET | `/api/v1/mikrotik/:router_id/hotspot/expire-monitor/script` | Get expire monitor script |

### PPP (per router)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/v1/mikrotik/:router_id/ppp/secrets` | List secrets |
| POST | `/api/v1/mikrotik/:router_id/ppp/secrets` | Create secret |
| GET | `/api/v1/mikrotik/:router_id/ppp/secrets/:id` | Get secret |
| PUT | `/api/v1/mikrotik/:router_id/ppp/secrets/:id` | Update secret |
| DELETE | `/api/v1/mikrotik/:router_id/ppp/secrets/:id` | Delete secret |
| PATCH | `/api/v1/mikrotik/:router_id/ppp/secrets/:id/disable` | Disable secret |
| PATCH | `/api/v1/mikrotik/:router_id/ppp/secrets/:id/enable` | Enable secret |
| GET | `/api/v1/mikrotik/:router_id/ppp/profiles` | List profiles |
| POST | `/api/v1/mikrotik/:router_id/ppp/profiles` | Create profile |
| GET | `/api/v1/mikrotik/:router_id/ppp/profiles/:id` | Get profile |
| PUT | `/api/v1/mikrotik/:router_id/ppp/profiles/:id` | Update profile |
| DELETE | `/api/v1/mikrotik/:router_id/ppp/profiles/:id` | Delete profile |
| GET | `/api/v1/mikrotik/:router_id/ppp/active` | List active connections |
| DELETE | `/api/v1/mikrotik/:router_id/ppp/active/:id` | Disconnect active |

### Voucher (per router)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/mikrotik/:router_id/vouchers/generate` | Generate vouchers |
| GET | `/api/v1/mikrotik/:router_id/vouchers` | List vouchers |
| POST | `/api/v1/mikrotik/:router_id/vouchers/cache` | Cache vouchers |
| DELETE | `/api/v1/mikrotik/:router_id/vouchers` | Delete vouchers |

### Network

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/v1/mikrotik/:router_id/interfaces` | List interfaces |
| GET | `/api/v1/mikrotik/:router_id/interfaces/:name/traffic` | Get traffic stats |
| GET | `/api/v1/mikrotik/:router_id/nat` | List NAT rules |
| GET | `/api/v1/mikrotik/:router_id/queues` | List queues |
| GET | `/api/v1/mikrotik/:router_id/queues/parents` | List parent queues |
| GET | `/api/v1/mikrotik/:router_id/pools` | List address pools |

### System & Reports

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/v1/mikrotik/:router_id/system/resources` | System resources |
| GET | `/api/v1/mikrotik/:router_id/system/health` | Health check |
| GET | `/api/v1/mikrotik/:router_id/system/identity` | Router identity |
| GET | `/api/v1/mikrotik/:router_id/system/routerboard` | RouterBoard info |
| GET | `/api/v1/mikrotik/:router_id/system/clock` | System clock |
| GET | `/api/v1/mikrotik/:router_id/system/dashboard` | Dashboard data |
| GET | `/api/v1/mikrotik/:router_id/system/status` | Router status |
| GET | `/api/v1/mikrotik/:router_id/logs` | System logs |
| GET | `/api/v1/mikrotik/:router_id/logs/hotspot` | Hotspot logs |
| GET | `/api/v1/mikrotik/:router_id/logs/ppp` | PPP logs |
| GET | `/api/v1/mikrotik/:router_id/reports/sales` | Sales report |
| GET | `/api/v1/mikrotik/:router_id/reports/summary` | Report summary |
| GET | `/api/v1/mikrotik/:router_id/reports/export` | Export CSV |

### Health Check

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/health` | Health check endpoint |

## 🔌 WebSocket

WebSocket endpoints untuk monitoring real-time:

| Endpoint | Deskripsi |
|----------|-----------|
| `WS /api/v1/ws/mikrotik/monitor/resource/:router_id` | Resource monitoring |
| `WS /api/v1/ws/mikrotik/monitor/interface/:router_id` | Traffic monitoring |
| `WS /api/v1/ws/mikrotik/monitor/queue/:router_id` | Queue monitoring |
| `WS /api/v1/ws/mikrotik/monitor/ping/:router_id` | Ping monitoring |
| `WS /api/v1/ws/mikrotik/monitor/logs/:router_id` | System log monitoring |
| `WS /api/v1/ws/mikrotik/monitor/hotspot-logs/:router_id` | Hotspot log monitoring |
| `WS /api/v1/ws/mikrotik/monitor/ppp-logs/:router_id` | PPP log monitoring |
| `WS /api/v1/ws/mikrotik/monitor/ppp-active/:router_id` | PPP active users |
| `WS /api/v1/ws/mikrotik/monitor/ppp-inactive/:router_id` | PPP inactive users |
| `WS /api/v1/ws/mikrotik/monitor/hotspot-active/:router_id` | Hotspot active users |
| `WS /api/v1/ws/mikrotik/monitor/hotspot-inactive/:router_id` | Hotspot inactive users |

### Contoh Koneksi WebSocket

```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws/mikrotik/monitor/resource/1');

ws.onopen = () => {
  console.log('Connected to resource monitor');
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Resource data:', data);
  // { cpu_load: 15, memory_usage: 45, uptime: "1d2h3m", ... }
};

ws.onclose = () => {
  console.log('Disconnected');
};
```

## 📊 Monitoring

### Default Credentials

Setelah pertama kali menjalankan aplikasi, default admin user akan dibuat otomatis:

- **Username**: `admin`
- **Password**: `admin123`

> ⚠️ **IMPORTANT**: Segera ubah default password setelah login pertama!

### Default Router

Default router juga akan dibuat untuk testing:

- **Name**: MikroTik-1
- **Host**: 192.168.233.1
- **Port**: 8728
- **Username**: admin
- **Password**: r00t

Update konfigurasi router sesuai dengan MikroTik Anda.

## 🧪 Development

### Commands Makefile

```bash
make build          # Build binary
make run            # Run development server
make test           # Run tests
make clean          # Clean build artifacts
make fmt            # Format code
make lint           # Run linter
make deps           # Download dependencies
make docker-build   # Build Docker image
make docker-run     # Run with docker-compose
make docker-stop    # Stop docker-compose
```

### Testing

```bash
# Run all tests
go test -v ./...

# Run specific package test
go test -v ./internal/infrastructure/mikrotik/...

# Run integration tests (requires MikroTik connection)
go test -v ./internal/infrastructure/mikrotik/test/...
```

### Database Migrations

Migration otomatis dilakukan saat aplikasi start. Struktur tabel:

- `admin_users` - Administrator accounts
- `routers` - MikroTik router configurations
- `settings` - Application settings
- `print_templates` - Voucher print templates

## 🐳 Docker

### Development dengan Docker Compose

```bash
# Start services (PostgreSQL + Redis)
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f postgres
docker-compose logs -f redis
```

### Build Production Image

```bash
# Build image
docker build -t mikhmon-api:latest .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_HOST=host.docker.internal \
  -e DATABASE_PORT=5432 \
  -e DATABASE_USER=mikhmon \
  -e DATABASE_PASSWORD=mikhmon \
  -e DATABASE_NAME=mikhmon \
  mikhmon-api:latest
```

## 🔒 Keamanan

1. **Ganti Default Credentials**: Segera ubah username dan password default
2. **JWT Secret**: Gunakan secret key yang kuat untuk production
3. **HTTPS**: Gunakan HTTPS di production
4. **Firewall**: Batasi akses API dan MikroTik API port
5. **Redis**: Password protect Redis jika digunakan di production

## 📝 Lisensi

[MIT License](LICENSE)

## 🤝 Kontribusi

Kontribusi sangat diterima! Silakan buat Pull Request atau Issue.

## 📧 Support

Untuk pertanyaan atau bantuan, silakan buat issue di repository.

---

**Catatan**: Project ini dalam tahap pengembangan aktif. Fitur dan API dapat berubah sewaktu-waktu.
