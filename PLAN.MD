# Mikhmon v4 - PHP to Golang Migration Plan
 
## Context
 
Mikhmon v4 is a PHP-based MikroTik Hotspot Monitor — a web admin panel for managing MikroTik router hotspot users, vouchers, profiles, and sales reports. The current PHP codebase communicates directly with MikroTik RouterOS via a custom socket-based API client. There is no database — all hotspot data lives on the router, and app config is stored in a PHP file with custom delimiters.
 
**Why migrate?** Go provides a single-binary deployment, better performance, built-in concurrency for traffic monitoring, stronger type safety, and modern security practices. The current PHP app has no CSRF protection, uses weak XOR encryption for passwords, suppresses all errors, and has no input validation.
 
**Goal:** Produce a fully functional Go application that replicates all Mikhmon v4 features with improved security, maintainability, and deployment simplicity (single binary with embedded assets).
 
---
 
## Tech Stack
 
| Component | Choice | Rationale |
|-----------|--------|-----------|
| Web Framework | **Gin** | Most popular Go framework, great middleware ecosystem, simple API |
| RouterOS Client | **go-routeros/routeros v3** | Official Go client for MikroTik API |
| Templates | **html/template** (stdlib) | Secure by default, no external dependency |
| Session | **gorilla/sessions** | Simple cookie-based sessions, sufficient for admin panel |
| Config Storage | **JSON file** | Replace PHP config.php, simple and human-readable |
| Password Hashing | **bcrypt** (golang.org/x/crypto) | Industry standard, replaces weak XOR |
| QR Codes | **skip2/go-qrcode** | Server-side QR generation (replaces client-side QRious.js) |
| Embedded Assets | **embed** (stdlib) | Single binary with all static files |
 
---
 
## Project Structure
 
```
mikhmon_v4/
├── main.go                     # Entry point, server startup
├── go.mod / go.sum
├── config/
│   ├── config.go               # Config types and JSON read/write
│   └── mikhmon.json            # Runtime config file (routers, admin creds)
├── internal/
│   ├── routeros/
│   │   ├── client.go           # RouterOS connection wrapper (go-routeros v3)
│   │   └── pool.go             # Connection pool (per-router clients)
│   ├── handler/
│   │   ├── auth.go             # Login/logout handlers
│   │   ├── dashboard.go        # Dashboard data (system info, stats)
│   │   ├── hotspot.go          # Hotspot user CRUD, active sessions, hosts
│   │   ├── profile.go          # User profile CRUD
│   │   ├── voucher.go          # Voucher generation & printing
│   │   ├── report.go           # Sales reports & live reports
│   │   ├── admin.go            # Settings, router management, template editor
│   │   ├── traffic.go          # Real-time traffic monitoring (SSE/WebSocket)
│   │   └── system.go           # Logs, expire monitor, misc system endpoints
│   ├── middleware/
│   │   ├── auth.go             # Session validation middleware
│   │   └── csrf.go             # CSRF token middleware
│   ├── model/
│   │   ├── router.go           # Router config struct
│   │   ├── user.go             # Hotspot user struct
│   │   ├── profile.go          # Hotspot user profile struct
│   │   └── report.go           # Sales report struct
│   └── util/
│       ├── generator.go        # Voucher code generators (rand N/LC/UC/etc.)
│       └── format.go           # Data formatting helpers
├── web/
│   ├── static/                 # CSS, JS, images (embedded)
│   │   ├── css/
│   │   ├── js/
│   │   └── img/
│   └── templates/              # Go HTML templates (embedded)
│       ├── layout/
│       │   ├── base.html       # Base layout with header/footer/nav
│       │   └── menu.html       # Sidebar navigation
│       ├── login.html
│       ├── dashboard.html
│       ├── hotspot.html
│       ├── generate.html
│       ├── print_voucher.html
│       ├── report.html
│       ├── live_report.html
│       ├── log.html
│       ├── admin.html
│       ├── settings.html
│       ├── editor.html
│       └── about.html
├── voucher_templates/          # Voucher print templates (header/row/footer × default/small/thermal)
│   ├── header.default.txt
│   ├── row.default.txt
│   ├── footer.default.txt
│   ├── header.small.txt
│   ├── row.small.txt
│   ├── footer.small.txt
│   ├── header.thermal.txt
│   ├── row.thermal.txt
│   └── footer.thermal.txt
└── php_legacy/                 # Original PHP files (reference only, not deployed)
```
 
---
 
## Phased Implementation
 
### Phase 1: Foundation (Core infrastructure)
 
**Files to create:**
- `main.go` — Gin server setup, route registration, embed directives, graceful shutdown
- `go.mod` — Module init with dependencies
- `config/config.go` — Config structs (`AdminConfig`, `RouterConfig`, `AppConfig`), JSON load/save, password hashing with bcrypt
- `config/mikhmon.json` — Default config (empty routers list, default admin)
- `internal/routeros/client.go` — Wrapper around go-routeros v3: `Connect()`, `Run()`, `RunWithFilter()`, helper to parse reply sentences into `map[string]string` slices
- `internal/routeros/pool.go` — Map of session-name → `*routeros.Client`, lazy connect, auto-reconnect on failure
- `internal/middleware/auth.go` — Gin middleware: check gorilla session for "mikhmon" key, redirect to login if missing
- `internal/middleware/csrf.go` — CSRF token generation and validation middleware
- `internal/model/router.go` — `RouterConfig` struct with fields: Name, Host, Username, Password, HotspotName, DNSName, Currency, Phone, Email, InfoLP, IdleTimeout, SessionName, ReportFlag, Token
 
**Verification:**
- `go build` compiles successfully
- Server starts and serves a "hello world" page
- Can connect to a MikroTik router via go-routeros and print system identity
 
---
 
### Phase 2: Authentication & Admin
 
**Files to create:**
- `internal/handler/auth.go` — Login page (GET), login handler (POST: validate credentials, set session), logout handler (POST: destroy session)
- `web/templates/login.html` — Login form (port from `view/login.php`)
- `web/templates/layout/base.html` — Base template with header, nav, footer slots
- `web/templates/layout/menu.html` — Sidebar navigation
- `internal/handler/admin.go` — Router CRUD: add router config, remove router, save settings, theme management, template editor (save voucher templates)
- `web/templates/admin.html` — Admin settings page
- `web/templates/settings.html` — Router configuration page
- `web/templates/editor.html` — Template editor page
 
**Key mapping from PHP:**
| PHP | Go Handler |
|-----|------------|
| `view/login.php` + index.php login logic | `handler/auth.go` → `HandleLogin()`, `HandleLogout()` |
| `post/post_a_router.php` | `handler/admin.go` → `HandleAddRouter()`, `HandleRemoveRouter()`, `HandleSaveRouter()` |
| `post/post_template.php` | `handler/admin.go` → `HandleSaveTemplate()` |
| `view/admin.php`, `view/admin/settings.php`, `view/admin/editor.php` | Templates: `admin.html`, `settings.html`, `editor.html` |
 
**Verification:**
- Login/logout flow works
- Can add/remove/edit router configurations
- Session persists across requests
- CSRF tokens validate correctly
 
---
 
### Phase 3: Dashboard & System Info
 
**Files to create:**
- `internal/handler/dashboard.go` — Fetch system info from RouterOS:
  - `/system/clock/print`, `/system/resource/print`, `/system/routerboard/print`
  - `/system/identity/print`, `/system/health/print`
  - `/ip/hotspot/user/print` (count), `/ip/hotspot/active/print` (count)
  - Returns JSON for AJAX or renders template
- `web/templates/dashboard.html` — Dashboard page with system stats, hotspot counts, income summary
- `internal/handler/system.go` — Log viewing (`/log/print`), expire monitor management (`/system/scheduler/print`, `/add`, `/set`)
 
**Key mapping from PHP:**
| PHP | Go Handler |
|-----|------------|
| `get/get_dashboard.php` | `handler/dashboard.go` → `HandleGetDashboard()` |
| `get/get_connect.php` | `handler/dashboard.go` → `HandleTestConnection()` |
| `get/get_expire_mon.php` | `handler/system.go` → `HandleGetExpireMonitor()` |
| `post/post_expire_monitor.php` | `handler/system.go` → `HandleSetExpireMonitor()` |
| `view/log.php` + `get/get_log.php` (if exists) | `handler/system.go` → `HandleGetLogs()` |
 
**Verification:**
- Dashboard shows live system resource info from router
- Connection test endpoint works
- Expire monitor can be created/enabled/disabled
 
---
 
### Phase 4: Hotspot User Management (Core Feature)
 
**Files to create:**
- `internal/handler/hotspot.go` — All hotspot user operations:
  - List users (with server filter): `/ip/hotspot/user/print`
  - Get single user: `/ip/hotspot/user/print` with `.id` or `name` filter
  - Add user: `/ip/hotspot/user/add`
  - Update user: `/ip/hotspot/user/set` + `/ip/hotspot/user/reset-counters`
  - Remove user(s): `/ip/hotspot/user/remove`
  - List active sessions: `/ip/hotspot/active/print`
  - Remove active session: `/ip/hotspot/active/remove`
  - List hosts: `/ip/hotspot/host/print`
  - Remove host: `/ip/hotspot/host/remove`
  - List hotspot servers: `/ip/hotspot/print`
- `internal/model/user.go` — `HotspotUser` struct with JSON tags
- `web/templates/hotspot.html` — Hotspot management page (users, profiles, active, hosts tabs)
 
**Key mapping from PHP:**
| PHP | Go Handler |
|-----|------------|
| `get/get_users.php` | `HandleGetUsers()` |
| `get/get_user.php` | `HandleGetUser()` |
| `get/get_hotspot_active.php` | `HandleGetActiveSessions()` |
| `get/get_hotspot_server.php` | `HandleGetHotspotServers()` |
| `get/get_hosts.php` | `HandleGetHosts()` |
| `get/get_tot_users.php` | `HandleGetTotalUsers()` |
| `post/post_add_user.php` | `HandleAddUser()` |
| `post/post_update_user.php` | `HandleUpdateUser()` |
| `post/post_hotspot_remove.php` | `HandleRemoveHotspotItem()` |
 
**Verification:**
- Can list, add, edit, delete hotspot users
- Can view and kick active sessions
- Can view and remove hosts
- Server filtering works
 
---
 
### Phase 5: User Profiles
 
**Files to create:**
- `internal/handler/profile.go` — Profile CRUD:
  - List profiles: `/ip/hotspot/user/profile/print`
  - Get profile: with `.id` or `name` filter
  - Add profile: `/ip/hotspot/user/profile/add` with complex on-login script generation
  - Update profile: `/ip/hotspot/user/profile/set`
  - Remove profile: `/ip/hotspot/user/profile/remove`
  - List related data: address pools (`/ip/pool/print`), parent queues (`/queue/simple/print`), NAT rules (`/ip/firewall/nat/print`), interfaces (`/interface/print`)
- `internal/model/profile.go` — `UserProfile` struct
 
**Key mapping from PHP:**
| PHP | Go Handler |
|-----|------------|
| `get/get_profiles.php` | `HandleGetProfiles()` |
| `get/get_profile.php` | `HandleGetProfile()` |
| `get/get_addr_pool.php` | `HandleGetAddressPools()` |
| `get/get_parent_queue.php` | `HandleGetQueues()` |
| `get/get_nat.php` | `HandleGetNATRules()` |
| `get/get_interface.php` | `HandleGetInterfaces()` |
| `post/post_add_userprofile.php` | `HandleAddProfile()` |
| `post/post_update_userprofile.php` | `HandleUpdateProfile()` |
 
**Critical:** The on-login script generation (from `post_add_userprofile.php`) is complex — it generates RouterOS scripts for expiration monitoring, MAC locking, server locking, and sales recording. This logic must be carefully ported to Go string building.
 
**Verification:**
- Can create profiles with full on-login script
- Can update and delete profiles
- Address pools, queues, NAT, interfaces load correctly
 
---
 
### Phase 6: Voucher Generation & Printing
 
**Files to create:**
- `internal/handler/voucher.go` — Voucher operations:
  - Generate vouchers in bulk: loop `/ip/hotspot/user/add` with generated codes
  - Cache vouchers: `/ip/hotspot/user/print` filtered by comment
  - Print vouchers: render voucher template with variable substitution + QR codes
- `internal/util/generator.go` — Port all 7 random generators from `core/generator_functions.php`:
  - `RandN()`, `RandLC()`, `RandUC()`, `RandULC()`, `RandNLC()`, `RandNUC()`, `RandNULC()`
- `web/templates/generate.html` — Voucher generation form
- `web/templates/print_voucher.html` — Voucher print view with QR codes
 
**Template variable substitution** (in voucher templates):
- `%username%`, `%password%`, `%hotspotName%`, `%dnsName%`, `%validity%`
- `%limitUptime%`, `%limitBytesTotal%`, `%price%`, `%qrCode%`, `%#%`
 
**QR Code:** Use `skip2/go-qrcode` to generate QR as base64 PNG data URI server-side (replaces client-side QRious.js canvas rendering).
 
**Key mapping from PHP:**
| PHP | Go Handler |
|-----|------------|
| `post/post_generate_voucher.php` | `HandleGenerateVouchers()` |
| `post/post_cache_voucher.php` | `HandleCacheVouchers()` |
| `view/generate.php` | Template: `generate.html` |
| `view/print_voucher.php` | `HandlePrintVouchers()` + template |
 
**Verification:**
- Can generate N vouchers with various character types
- Voucher print page renders with correct data
- QR codes contain correct login URL
- Voucher templates (default/small/thermal) work
 
---
 
### Phase 7: Reports & Traffic Monitoring
 
**Files to create:**
- `internal/handler/report.go` — Sales reports:
  - Fetch reports from RouterOS system scripts: `/system/script/print` (parse owner format `MMYYYYcomment`)
  - Filter by date range, profile, etc.
  - Live report: daily/monthly income aggregation
- `internal/handler/traffic.go` — Real-time traffic:
  - Use go-routeros `Listen()` for `/interface/monitor-traffic`
  - Serve via Server-Sent Events (SSE) or WebSocket for live updates
- `internal/model/report.go` — `SalesReport` struct
- `web/templates/report.html` — Report page with filtering
- `web/templates/live_report.html` — Live income dashboard
 
**Key mapping from PHP:**
| PHP | Go Handler |
|-----|------------|
| `get/get_report.php` | `HandleGetReports()` |
| `get/get_traffic.php` | `HandleGetTraffic()` (SSE endpoint) |
| `view/report.php` | Template: `report.html` |
| `view/live_report.php` | Template: `live_report.html` |
 
**Verification:**
- Sales reports load and filter correctly
- Live traffic updates stream in real-time
- Report date filtering works
 
---
 
### Phase 8: Frontend & Static Assets
 
**Tasks:**
- Port all CSS/JS from `assets/` to `web/static/`
- Update JavaScript to work with new Go API endpoints (same JSON format)
- Ensure all AJAX calls include CSRF token
- Test all pages render correctly
- Add `//go:embed web/static` and `//go:embed web/templates` directives
- Embed voucher templates with `//go:embed voucher_templates`
 
**Verification:**
- All pages render identically to PHP version
- JavaScript interactions work (AJAX user/profile CRUD)
- Theme switching works
- Mobile responsiveness preserved
 
---
 
### Phase 9: Config Migration & Polish
 
**Tasks:**
- Write a one-time migration tool (in `cmd/migrate/main.go`) to convert `config/config.php` → `config/mikhmon.json`
  - Parse PHP delimiter format
  - Decrypt XOR passwords and re-hash with bcrypt
- Add graceful shutdown handling
- Add structured logging (slog stdlib)
- Add request rate limiting middleware
- Final security review: input validation on all POST handlers
 
**Verification:**
- Migration tool successfully converts existing PHP config
- Application works end-to-end with migrated config
- All features tested against real MikroTik router
 
---
 
## Route Mapping Summary
 
```go
// Auth (no middleware)
POST /login           → auth.HandleLogin
GET  /login           → auth.HandleLoginPage
POST /logout          → auth.HandleLogout
 
// All routes below require auth middleware
// Dashboard
GET  /:session/dashboard          → dashboard.HandleDashboardPage
GET  /:session/api/dashboard      → dashboard.HandleGetDashboard (JSON)
GET  /:session/api/connect        → dashboard.HandleTestConnection (JSON)
 
// Hotspot Users
GET  /:session/hotspot            → hotspot.HandleHotspotPage
GET  /:session/api/users          → hotspot.HandleGetUsers (JSON)
GET  /:session/api/user/:id       → hotspot.HandleGetUser (JSON)
POST /:session/api/user           → hotspot.HandleAddUser (JSON)
PUT  /:session/api/user/:id       → hotspot.HandleUpdateUser (JSON)
DELETE /:session/api/user/:id     → hotspot.HandleRemoveUser (JSON)
 
// Active Sessions
GET  /:session/api/active         → hotspot.HandleGetActiveSessions (JSON)
DELETE /:session/api/active/:id   → hotspot.HandleRemoveActiveSession (JSON)
 
// Hosts
GET  /:session/api/hosts          → hotspot.HandleGetHosts (JSON)
DELETE /:session/api/host/:id     → hotspot.HandleRemoveHost (JSON)
 
// Hotspot Servers
GET  /:session/api/servers        → hotspot.HandleGetHotspotServers (JSON)
 
// Profiles
GET  /:session/api/profiles       → profile.HandleGetProfiles (JSON)
GET  /:session/api/profile/:id    → profile.HandleGetProfile (JSON)
POST /:session/api/profile        → profile.HandleAddProfile (JSON)
PUT  /:session/api/profile/:id    → profile.HandleUpdateProfile (JSON)
DELETE /:session/api/profile/:id  → profile.HandleRemoveProfile (JSON)
 
// Supporting data for profiles
GET  /:session/api/pools          → profile.HandleGetAddressPools (JSON)
GET  /:session/api/queues         → profile.HandleGetQueues (JSON)
GET  /:session/api/nat            → profile.HandleGetNATRules (JSON)
GET  /:session/api/interfaces     → profile.HandleGetInterfaces (JSON)
 
// Vouchers
GET  /:session/generate           → voucher.HandleGeneratePage
POST /:session/api/generate       → voucher.HandleGenerateVouchers (JSON)
POST /:session/api/cache-voucher  → voucher.HandleCacheVouchers (JSON)
GET  /:session/print-voucher      → voucher.HandlePrintVouchers (HTML)
 
// Reports
GET  /:session/report             → report.HandleReportPage
GET  /:session/api/report         → report.HandleGetReports (JSON)
GET  /:session/live-report        → report.HandleLiveReportPage
 
// Traffic
GET  /:session/api/traffic        → traffic.HandleGetTraffic (SSE)
 
// System
GET  /:session/log                → system.HandleLogPage
GET  /:session/api/logs           → system.HandleGetLogs (JSON)
GET  /:session/api/expire-monitor → system.HandleGetExpireMonitor (JSON)
POST /:session/api/expire-monitor → system.HandleSetExpireMonitor (JSON)
 
// Admin
GET  /admin/settings              → admin.HandleSettingsPage
POST /admin/api/router            → admin.HandleAddRouter (JSON)
PUT  /admin/api/router/:name      → admin.HandleSaveRouter (JSON)
DELETE /admin/api/router/:name    → admin.HandleRemoveRouter (JSON)
POST /admin/api/template          → admin.HandleSaveTemplate (JSON)
GET  /admin/template-editor       → admin.HandleEditorPage
```
 
---
 
## Key Implementation Notes
 
### RouterOS Client Wrapper (`internal/routeros/client.go`)
```go
// Core wrapper functions needed:
func Connect(host, user, password string) (*routeros.Client, error)
func RunCommand(client *routeros.Client, command string, args ...string) ([]map[string]string, error)
// RunCommand sends command, parses reply.Re sentences into []map[string]string
// Each sentence's Map() gives key-value pairs
```
 
### On-Login Script Generation (Phase 5 - Critical)
The most complex piece to port is the on-login script builder from `post/post_add_userprofile.php`. This generates RouterOS scripting language code that runs when a user connects. It must be ported as a Go string builder that produces identical RouterOS scripts.
 
### Session-Scoped Router Connections
Each browser session can switch between multiple routers. The connection pool (`internal/routeros/pool.go`) should:
- Keep connections alive per router session name
- Reconnect on connection failure
- Close idle connections after timeout
 
### Config File Format (`config/mikhmon.json`)
```json
{
  "admin": {
    "username": "admin",
    "password_hash": "$2a$10$..."
  },
  "routers": [
    {
      "session_name": "router1",
      "host": "192.168.1.1:8728",
      "username": "admin",
      "password": "encrypted_or_plain",
      "hotspot_name": "hotspot1",
      "dns_name": "hotspot.example.com",
      "currency": "Rp",
      "phone": "",
      "email": "",
      "info_lp": "",
      "idle_timeout": "30",
      "report_flag": "",
      "token": ""
    }
  ],
  "server": {
    "port": 8080,
    "session_secret": "random-32-byte-key"
  }
}
```
 
---
 
## Verification Plan (End-to-End)
 
1. **Build:** `go build -o mikhmon` produces single binary
2. **Start:** `./mikhmon` starts server on configured port
3. **Login:** Access web UI, login with admin credentials
4. **Add Router:** Configure a MikroTik router connection
5. **Dashboard:** View system resources, hotspot stats
6. **Users:** List, add, edit, delete hotspot users
7. **Profiles:** Create profile with on-login script, verify script on router
8. **Vouchers:** Generate 10 vouchers, print with QR codes
9. **Reports:** View sales reports, filter by date
10. **Traffic:** View live traffic monitoring
11. **Logs:** View system logs from router
12. **Expire Monitor:** Enable/disable scheduler
13. **Template Editor:** Edit and save voucher template
14. **Multi-Router:** Switch between router sessions
15. **Migration:** Run migration tool on existing PHP config, verify all routers imported