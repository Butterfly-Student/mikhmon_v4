# 📡 MIKHMON v4 - API Endpoints Documentation

**Comprehensive API Reference dengan Request & Response**

---

## 📋 Table of Contents

1. [API Base URL & Authentication](#api-base-url--authentication)
2. [GET Endpoints (Data Retrieval)](#get-endpoints-data-retrieval)
3. [POST Endpoints (Data Modification)](#post-endpoints-data-modification)
4. [Response Format & Error Handling](#response-format--error-handling)
5. [Complete Request-Response Examples](#complete-request-response-examples)

---

## 🔐 API Base URL & Authentication

### Base URL

```
http://localhost/mikhmon_v4/
```

### Session-Based Authentication

MIKHMON menggunakan session-based authentication, bukan API keys.

```php
// Session dimulai saat login
session_start();
$_SESSION["mikhmon"] = true;    // Session flag
$_SESSION["m_user"] = "user01"; // Current user session

// Untuk multi-router support:
// $_SESSION[$m_user] = session_name
// $data[$m_user] = array(
//     1 => "session!IP",
//     2 => "session@|@username",
//     3 => "session#|#password"
// )
```

### Required Headers

```http
Content-Type: application/x-www-form-urlencoded
Content-Type: application/json
Cookie: PHPSESSID=<session_id>
```

---

## 📥 GET Endpoints (Data Retrieval)

### 1. Get System Resource Information

**Endpoint**: `/index.php?<SESSION>/dashboard&get_sys_resource`

**Method**: `GET`

**Description**: Mengambil informasi resource sistem RouterOS (CPU, Memory, Uptime, Health)

**Request Parameters**: Tidak ada

**Request Example**:

```http
GET /index.php?user01/dashboard&get_sys_resource HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
{
  "systime": {
    "time": "Feb/25/2026 14:45:00",
    "time-zone-name": "Asia/Jakarta",
    "gmt-offset": "+07:00"
  },
  "resource": {
    "uptime": "42d23h15m30s",
    "cpu-load": "45",
    "free-memory": "2097152",
    "total-memory": "4194304",
    "cpu-count": "2"
  },
  "syshealth": {
    "temperature": "65C",
    "voltage": "12V"
  },
  "model": "RB951G-2HnD",
  "identity": "MikroTik"
}
```

**File**: `get/get_dashboard.php`

---

### 2. Get Hotspot Information

**Endpoint**: `/index.php?<SESSION>/dashboard&get_hotspotinfo`

**Method**: `GET`

**Description**: Mengambil informasi user hotspot (total user, user aktif)

**Request Parameters**: Tidak ada

**Request Example**:

```http
GET /index.php?user01/dashboard&get_hotspotinfo HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
{
  "hotspot_users": 150,
  "hotspot_active": 45
}
```

**Penjelasan**:
- `hotspot_users`: Total user hotspot terdaftar (dikurangi 1 admin)
- `hotspot_active`: Jumlah user yang sedang aktif/connected

**File**: `get/get_dashboard.php`

---

### 3. Get Traffic Real-Time

**Endpoint**: `/index.php?<SESSION>/dashboard&get_traffic`

**Method**: `GET`

**Query Parameters**:
- `iface` (string, required): Nama interface untuk monitoring

**Request Example**:

```http
GET /index.php?user01/dashboard&get_traffic&iface=ether1 HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
{
  "tx": 5242880,
  "rx": 10485760
}
```

**Penjelasan**:
- `tx`: Transmit bits per second (downstream)
- `rx`: Receive bits per second (upstream)
- Nilai dalam bits/second (bukan bytes)

**Conversions**:
```
Kbps = bits/second / 1000
Mbps = bits/second / 1000000
Gbps = bits/second / 1000000000
```

**File**: `get/get_dashboard.php`

---

### 4. Get System Log

**Endpoint**: `/index.php?<SESSION>/dashboard&get_log`

**Method**: `GET`

**Query Parameters**:
- `f` (string, optional): "true" untuk force refresh, "false" untuk cache

**Request Example**:

```http
GET /index.php?user01/dashboard&get_log&f=true HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
[
  {
    ".id": "*1",
    "time": "Feb/25/2026 14:45:00",
    "topics": "system,info",
    "message": "System startup",
    "account": "admin"
  },
  {
    ".id": "*2",
    "time": "Feb/25/2026 14:46:15",
    "topics": "hotspot,info",
    "message": "User admin created",
    "account": "admin"
  }
]
```

**File**: `get/get_dashboard.php`

---

### 5. Get Hotspot Active Users

**Endpoint**: `/index.php?<SESSION>/hotspot&get_hotspot_active`

**Method**: `GET`

**Description**: Mengambil daftar user yang sedang aktif/connected ke hotspot

**Request Parameters**: Tidak ada

**Request Example**:

```http
GET /index.php?user01/hotspot&get_hotspot_active HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
[
  {
    ".id": "*1",
    "name": "user123",
    "address": "192.168.1.100",
    "mac-address": "00:11:22:33:44:55",
    "login-time": "Feb/25/2026 14:30:00",
    "uptime": "15m30s",
    "bytes-in": "1048576",
    "bytes-out": "2097152"
  },
  {
    ".id": "*2",
    "name": "user456",
    "address": "192.168.1.101",
    "mac-address": "66:77:88:99:AA:BB",
    "login-time": "Feb/25/2026 14:35:00",
    "uptime": "10m",
    "bytes-in": "512000",
    "bytes-out": "1024000"
  }
]
```

**File**: `get/get_hotspot_active.php`

---

### 6. Get Hotspot Server List

**Endpoint**: `/index.php?<SESSION>/hotspot&get_hotspot_server`

**Method**: `GET`

**Query Parameters**:
- `f` (string, optional): "true" untuk force refresh

**Request Example**:

```http
GET /index.php?user01/hotspot&get_hotspot_server&f=false HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
[
  {
    ".id": "*1",
    "name": "Hotspot-Server-1",
    "address-pool": "Pool-Main",
    "profile": "default",
    "interface": "ether2"
  },
  {
    ".id": "*2",
    "name": "Hotspot-Server-2",
    "address-pool": "Pool-Guest",
    "profile": "guest",
    "interface": "ether3"
  }
]
```

**File**: `get/get_hotspot_server.php`

---

### 7. Get Hotspot Hosts/MAC List

**Endpoint**: `/index.php?<SESSION>/hotspot&get_hosts`

**Method**: `GET`

**Description**: Mengambil daftar host yang ter-register di hotspot

**Request Parameters**: Tidak ada

**Request Example**:

```http
GET /index.php?user01/hotspot&get_hosts HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
[
  {
    ".id": "*1",
    "mac-address": "00:11:22:33:44:55",
    "comment": "Office-Device-1",
    "address": "192.168.1.50"
  },
  {
    ".id": "*2",
    "mac-address": "66:77:88:99:AA:BB",
    "comment": "Guest-Laptop",
    "address": "192.168.1.51"
  }
]
```

**File**: `get/get_hosts.php`

---

### 8. Get Sales Report

**Endpoint**: `/index.php?<SESSION>/report&get_report`

**Method**: `GET`

**Query Parameters**:
- `day` (string, required): Format "Mon/dd/YYYY" (e.g., "Feb/25/2026")
- `f` (string, optional): "true" untuk force refresh

**Request Example**:

```http
GET /index.php?user01/report&get_report&day=Feb/25/2026&f=false HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
[
  {
    ".id": "3",
    "name": "up-ABC123-02.25.26-OfficeStaff",
    "owner": "0202",
    "source": "Feb/25/2026",
    "comment": "mikhmon"
  },
  {
    ".id": "4",
    "name": "vc-XYZ789-02.25.26-Premium1Day",
    "owner": "0202",
    "source": "Feb/25/2026",
    "comment": "mikhmon"
  }
]
```

**Penjelasan**:
- Data disimpan sebagai `/system/script` dengan nama = transaction record
- Format name: `<type>-<code>-<date>-<description>`
- `owner`: Bulan+Tahun (MMYY format)
- `source`: Tanggal transaks

i (MMM/DD/YYYY format)

**File**: `get/get_report.php`

---

### 9. Get User Total Count

**Endpoint**: `/index.php?<SESSION>/users&get_tot_users`

**Method**: `GET`

**Query Parameters**:
- `name` (string, optional): Username (untuk compatibility)

**Request Example**:

```http
GET /index.php?user01/users&get_tot_users&name=admin HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
{
  "users": 149
}
```

**Penjelasan**: Total user dikurangi 1 (admin user)

**File**: `get/get_tot_users.php`

---

### 10. Get Network Interfaces

**Endpoint**: `/index.php?<SESSION>/settings&get_interface`

**Method**: `GET`

**Description**: Mengambil daftar interface network untuk configuration

**Request Parameters**: Tidak ada

**Request Example**:

```http
GET /index.php?user01/settings&get_interface HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
[
  {
    ".id": "*1",
    "name": "ether1",
    "type": "ether",
    "mtu": "1500",
    "mac-address": "00:11:22:33:44:55",
    "running": "true"
  },
  {
    ".id": "*2",
    "name": "ether2",
    "type": "ether",
    "mtu": "1500",
    "mac-address": "66:77:88:99:AA:BB",
    "running": "true"
  }
]
```

**File**: `get/get_interface.php`

---

### 11. Get Address Pool List

**Endpoint**: `/index.php?<SESSION>/profile&get_addr_pool`

**Method**: `GET`

**Query Parameters**:
- `f` (string, optional): "true" untuk force refresh

**Request Example**:

```http
GET /index.php?user01/profile&get_addr_pool&f=false HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
[
  {
    ".id": "*1",
    "name": "Pool-Main",
    "ranges": "192.168.1.100-192.168.1.200"
  },
  {
    ".id": "*2",
    "name": "Pool-Guest",
    "ranges": "192.168.2.100-192.168.2.200"
  }
]
```

**File**: `get/get_addr_pool.php`

---

### 12. Get Parent Queue List

**Endpoint**: `/index.php?<SESSION>/profile&get_parent_queue`

**Method**: `GET`

**Query Parameters**:
- `f` (string, optional): "true" untuk force refresh

**Request Example**:

```http
GET /index.php?user01/profile&get_parent_queue&f=true HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
[
  {
    ".id": "*1",
    "name": "QueueMain",
    "max-limit": "10M",
    "burst-limit": "20M"
  },
  {
    ".id": "*2",
    "name": "QueueGuest",
    "max-limit": "5M",
    "burst-limit": "10M"
  }
]
```

**File**: `get/get_parent_queue.php`

---

### 13. Get NAT Rules

**Endpoint**: `/index.php?<SESSION>/settings&get_nat`

**Method**: `GET`

**Description**: Mengambil daftar NAT rules untuk konfigurasi

**Request Parameters**: Tidak ada

**Request Example**:

```http
GET /index.php?user01/settings&get_nat HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```json
[
  {
    ".id": "*1",
    "chain": "srcnat",
    "action": "masquerade",
    "out-interface": "ether1",
    "comment": "Masquerade-LAN"
  },
  {
    ".id": "*2",
    "chain": "srcnat",
    "action": "masquerade",
    "out-interface": "ether2",
    "comment": "Masquerade-Hotspot"
  }
]
```

**File**: `get/get_nat.php`

---

### 14. Test Connection

**Endpoint**: `/index.php?<SESSION>/connect`

**Method**: `GET`

**Description**: Test koneksi ke RouterOS device

**Request Parameters**: Tidak ada

**Request Example**:

```http
GET /index.php?user01/connect HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

**Response** (Success - Status 200):

```
Connected
```

**Possible Responses**:
- `Connected` - Koneksi berhasil
- `Invalid username or password` - Kredensial salah
- `Error` - Koneksi gagal ke RouterOS

**File**: `get/get_connect.php`

---

## 📤 POST Endpoints (Data Modification)

### 1. Add Hotspot User

**Endpoint**: `/post/post_add_user.php`

**Method**: `POST`

**Description**: Membuat user hotspot baru (username/password atau voucher)

**Request Content-Type**: `application/x-www-form-urlencoded`

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sessname` | string | Yes | Session name (format: "?user01") |
| `server` | string | Yes | Nama hotspot server |
| `name` | string | Yes | Username |
| `password` | string | Yes | Password |
| `profile` | string | Yes | Nama profile |
| `macaddr` | string | No | MAC address (00:00:00:00:00:00 jika kosong) |
| `timelimit` | string | No | Time limit (1d, 2h, 30m, dll) |
| `datalimit` | string | No | Data limit (1G, 500M, 1024000 bytes) |
| `comment` | string | No | Keterangan tambahan |

**Request Example**:

```http
POST /post/post_add_user.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&server=Hotspot-Server-1&name=testuser&password=testpass123&profile=1day-profile&macaddr=&timelimit=1d&datalimit=1G&comment=Test+User
```

**Response** (Success - Status 200):

```json
{
  "message": "success",
  "data": {
    ".id": "*5",
    "name": "testuser",
    "password": "testpass123",
    "profile": "1day-profile",
    "server": "Hotspot-Server-1",
    "mac-address": "00:00:00:00:00:00",
    "disabled": "false",
    "limit-uptime": "1d",
    "limit-bytes-total": "1073741824",
    "comment": "up-testuser-02.25.26-Test User",
    "uptime": "0s",
    "bytes-in": "0",
    "bytes-out": "0"
  }
}
```

**Response** (Error - Status 200):

```json
{
  "message": "error",
  "data": {
    "error": "duplicate name"
  }
}
```

**Penjelasan Comment Format**:
- Jika `name == password` → Voucher type: `"vc-..."`
- Jika `name != password` → User type: `"up-..."`

**File**: `post/post_add_user.php`

---

### 2. Update Hotspot User

**Endpoint**: `/post/post_update_user.php`

**Method**: `POST`

**Description**: Update/edit user hotspot yang sudah ada

**Request Content-Type**: `application/x-www-form-urlencoded`

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sessname` | string | Yes | Session name (format: "?user01") |
| `uid` | string | Yes | User ID (.id dari /ip/hotspot/user) |
| `server` | string | Yes | Nama hotspot server |
| `name` | string | Yes | Username |
| `password` | string | Yes | Password |
| `profile` | string | Yes | Nama profile |
| `macaddr` | string | No | MAC address |
| `timelimit` | string | No | Time limit |
| `datalimit` | string | No | Data limit |
| `comment` | string | No | Keterangan |
| `expdate` | string | No | Expiry date (DD/MM/YYYY format) |
| `ucode` | string | No | User code/voucher code |
| `reset` | string | No | "yes" untuk reset counter |

**Request Example**:

```http
POST /post/post_update_user.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&uid=*5&server=Hotspot-Server-1&name=testuser&password=newpass123&profile=1day-profile&macaddr=&timelimit=2d&datalimit=2G&comment=Updated+Test&expdate=&ucode=&reset=no
```

**Response** (Success - Status 200):

```json
{
  "message": "success",
  "data": {
    ".id": "*5",
    "name": "testuser",
    "password": "newpass123",
    "profile": "1day-profile",
    "limit-uptime": "2d",
    "limit-bytes-total": "2147483648",
    "comment": "up-testuser-02.25.26-Updated Test"
  }
}
```

**Reset Counter Example**:

```http
POST /post/post_update_user.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&uid=*5&reset=yes&...
```

Effect: User's uptime dan bytes-in/bytes-out di-reset ke 0

**File**: `post/post_update_user.php`

---

### 3. Generate Vouchers

**Endpoint**: `/post/post_generate_voucher.php`

**Method**: `POST`

**Description**: Generate multiple vouchers sekaligus dengan setting yang sama

**Request Content-Type**: `application/x-www-form-urlencoded`

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sessname` | string | Yes | Session name |
| `qty` | integer | Yes | Jumlah voucher (1-50 per batch) |
| `server` | string | Yes | Hotspot server name |
| `user` | string | Yes | Type: "up" (user) atau "vc" (voucher) |
| `userl` | integer | Yes | Length username (4-16 karakter) |
| `prefix` | string | No | Prefix untuk username |
| `char` | string | Yes | Character type: "lower", "upper", "upplow", "mix", "num" |
| `profile` | string | Yes | Profile name |
| `timelimit` | string | No | Time limit (1d, 7d, 30d, dll) |
| `datalimit` | string | No | Data limit (100M, 1G, dll) |
| `gcomment` | string | No | Comment untuk group voucher |
| `gencode` | integer | Yes | Generation code (random ID untuk batch) |

**Request Example**:

```http
POST /post/post_generate_voucher.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&qty=50&server=Hotspot-Server-1&user=vc&userl=8&prefix=&char=mix&profile=1day-profile&timelimit=1d&datalimit=1G&gcomment=Daily+Voucher&gencode=538
```

**Response** (Success - Status 200):

```json
{
  "message": "success",
  "data": {
    "count": "50",
    "comment": "vc-538-02.25.26-Daily Voucher"
  }
}
```

**Penjelasan**:
- Script akan generate 50 user dengan username random
- Format: `vc-<gencode>-<date>-<comment>`
- Response memberikan total count dan comment yang disimpan
- Client-side akan melakukan batch request jika qty > 50

**Character Types**:
- `lower`: Lowercase letters only (abc...xyz)
- `upper`: Uppercase letters only (ABC...XYZ)
- `upplow`: Mixed uppercase+lowercase (AaBbCc...)
- `mix`: Mixed with numbers (Aa1Bb2Cc3...)
- `num`: Numbers only (0123456789)

**Data Limit Formats**:
- `100M` = 100 * 1048576 bytes
- `1G` = 1 * 1073741824 bytes
- `1024000` = 1024000 bytes

**File**: `post/post_generate_voucher.php`

---

### 4. Cache Generated Vouchers

**Endpoint**: `/post/post_cache_voucher.php`

**Method**: `POST`

**Description**: Cache voucher yang sudah di-generate untuk printing

**Request Content-Type**: `application/x-www-form-urlencoded`

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sessname` | string | Yes | Session name |
| `qty` | integer | Yes | Jumlah |
| `user` | string | Yes | "up" atau "vc" |
| `gcomment` | string | Yes | Comment yang sama dengan generate |
| `gencode` | integer | Yes | Generation code yang sama |

**Request Example**:

```http
POST /post/post_cache_voucher.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&qty=50&user=vc&gcomment=Daily+Voucher&gencode=538
```

**Response** (Success - Status 200):

```json
{
  "message": "success",
  "data": {
    "count": "50",
    "comment": "vc-538-02.25.26-Daily Voucher"
  }
}
```

**Response** (Error - Status 200):

```json
{
  "message": "error",
  "data": {
    "error": "not found"
  }
}
```

**Penjelasan**:
- Cache tersimpan di `$_SESSION[$m_user.$comment]`
- Digunakan sebelum print untuk ambil data dari session

**File**: `post/post_cache_voucher.php`

---

### 5. Add User Profile

**Endpoint**: `/post/post_add_userprofile.php`

**Method**: `POST`

**Description**: Membuat profile/paket baru dengan on-login script

**Request Content-Type**: `application/x-www-form-urlencoded`

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sessname` | string | Yes | Session name |
| `name` | string | Yes | Profile name |
| `addresspool` | string | Yes | Address pool name |
| `sharedusers` | integer | Yes | Jumlah shared users |
| `ratelimit` | string | Yes | Rate limit (1M, 10M, dll) |
| `parentqueue` | string | Yes | Parent queue name |
| `expmode` | string | Yes | Expire mode: "ntf" (notify) atau "rem" (remove) |
| `validity` | string | Yes | Validity period (1d, 7d, 30d, etc) |
| `price` | decimal | No | Selling price |
| `sellingprice` | decimal | No | Purchase price |
| `lockuser` | string | Yes | "Enable" atau "Disable" MAC lock |
| `lockserver` | string | Yes | "Enable" atau "Disable" server lock |

**Request Example**:

```http
POST /post/post_add_userprofile.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&name=Premium-1Day&addresspool=Pool-Main&sharedusers=1&ratelimit=10M&parentqueue=QueueMain&expmode=rem&validity=1d&price=25000&sellingprice=20000&lockuser=Enable&lockserver=Disable
```

**Response** (Success - Status 200):

```json
{
  "message": "success",
  "data": {
    ".id": "*3",
    "name": "Premium-1Day",
    "address-pool": "Pool-Main",
    "shared-users": "1",
    "rate-limit": "10M",
    "parent-queue": "QueueMain",
    "on-login": ":put (\"rem,25000,1d,20000,,Enable,Disable,\"); :local mode \"X\"; {...on-login script...}"
  }
}
```

**on-login Script Content**:

Format yang disimpan di `on-login`:
```
:put ("mode,price,validity,sprice,,lockuser,lockserver,");
:local mode "N/X";
{script untuk handling expiry}
```

Parsing di client:
```javascript
let parts = onlogin.split(",");
let mode = parts[1];        // "ntf" atau "rem"
let price = parts[2];       // Harga jual
let validity = parts[3];    // Durasi
let sprice = parts[4];      // Harga beli
let lockuser = parts[6];    // Lock MAC
let lockserver = parts[7];  // Lock server
```

**File**: `post/post_add_userprofile.php`

---

### 6. Remove Hotspot Resources

**Endpoint**: `/post/post_hotspot_remove.php`

**Method**: `POST`

**Description**: Hapus user, profile, active session, atau host

**Request Content-Type**: `application/x-www-form-urlencoded`

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sessname` | string | Yes | Session name |
| `where` | string | Yes | "user_", "profile_", "active_", atau "host_" |
| `id` | string | Yes | ID yang akan dihapus |

**Request Examples**:

**Delete User**:
```http
POST /post/post_hotspot_remove.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&where=user_&id=*5
```

**Delete Profile**:
```http
POST /post/post_hotspot_remove.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&where=profile_&id=*3
```

**Force Logout Active User**:
```http
POST /post/post_hotspot_remove.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&where=active_&id=*1
```

**Delete Host/MAC**:
```http
POST /post/post_hotspot_remove.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&where=host_&id=*2
```

**Response** (Success - Status 200):

```json
{
  "message": "success"
}
```

**Response** (Error - Status 200):

```json
{
  "message": "error"
}
```

**File**: `post/post_hotspot_remove.php`

---

### 7. Setup Expire Monitor

**Endpoint**: `/post/post_expire_monitor.php`

**Method**: `POST`

**Description**: Setup scheduler untuk auto-expire user yang sudah kadaluarsa

**Request Content-Type**: `application/x-www-form-urlencoded`

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sessname` | string | Yes | Session name |
| `expmon` | string | Yes | RouterOS expire monitor script |

**Request Example**:

```http
POST /post/post_expire_monitor.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&expmon=:local+dateint+do={...}
```

Script terlalu panjang, lihat di [EXPIRE MONITOR SCRIPT Section](#expire-monitor-script)

**Response** (Success - Status 200):

```json
{
  "message": "success"
}
```

**Possible Responses**:
- `"success"` - Scheduler created/updated
- `"Mikhmon-Expire-Monitor"` - Already exists dan aktif

**Effects di RouterOS**:
1. Jika scheduler belum ada: Create `/system/scheduler` baru dengan nama "Mikhmon-Expire-Monitor"
2. Jika scheduler ada tapi disabled: Enable kembali
3. Interval: 1 menit (00:01:00)

**File**: `post/post_expire_monitor.php`

---

### 8. Logout

**Endpoint**: `/post/post_logout.php`

**Method**: `POST`

**Description**: Logout user dan destroy session

**Request Content-Type**: `application/x-www-form-urlencoded`

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `logout` | string | Yes | Logout message/confirmation |

**Request Example**:

```http
POST /post/post_logout.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

logout=Logout+Success
```

**Response** (Status 200):

```
Logout Success
```

**Effect**:
- Session destroyed
- Redirect ke login page
- All cached data cleared

**File**: `post/post_logout.php`

---

### 9. Add Router/Session

**Endpoint**: `/post/post_a_router.php`

**Method**: `POST`

**Description**: Manage multi-router configuration

**Request Content-Type**: `application/x-www-form-urlencoded`

**Request Parameters** (untuk Add):

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `do` | string | Yes | "add", "remove", atau "save" |
| `router_` | string | Yes | Router identifier |
| `session` | string | Yes (save) | Session name (format: "user01") |
| `ipmik` | string | Yes (save) | RouterOS IP (encoded) |
| `usermik` | string | Yes (save) | RouterOS username (encoded) |
| `passmik` | string | Yes (save) | RouterOS password (encoded) |
| `hotspotname` | string | Yes (save) | Hotspot name |
| `currency` | string | No (save) | Currency symbol |
| `phone` | string | No (save) | Admin phone |

**Add New Router Example**:

```http
POST /post/post_a_router.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

do=add&router_=sess_
```

**Response** (Success - Status 200):

```json
{
  "message": "Success",
  "sesname": "session123"
}
```

**Save Router Config Example**:

```http
POST /post/post_a_router.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

do=save&router_=session123&session=user01&ipmik=<encoded_ip>&usermik=<encoded_user>&passmik=<encoded_pass>&hotspotname=MyHotspot&currency=Rp&phone=628123456789
```

**Response** (Success - Status 200):

```json
{
  "message": "Success",
  "sess": "session123"
}
```

**Remove Router Example**:

```http
POST /post/post_a_router.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

do=remove&router_=session123
```

**Response** (Success - Status 200):

```json
{
  "message": "Success"
}
```

**Encoding Function** (JavaScript):
```javascript
function blah(value) {
    let encoded = btoa(btoa(value));  // Double base64
    let xored = "";
    for (let i = 0; i < encoded.length; i++) {
        xored += String.fromCharCode(encoded.charCodeAt(i) ^ 10);
    }
    return btoa(xored);
}
```

**File**: `post/post_a_router.php`

---

## 📋 Response Format & Error Handling

### Standard Response Format

**Success Response**:
```json
{
  "message": "success",
  "data": {
    // Response data
  }
}
```

**Error Response**:
```json
{
  "message": "error",
  "data": {
    "error": "Error description"
  }
}
```

### Common HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success (atau error dalam JSON) |
| 403 | Forbidden (direct access without index.php) |
| 500 | Server error |

### Common Error Messages

| Error | Cause |
|-------|-------|
| `"not found"` | Data tidak ditemukan |
| `"duplicate name"` | Username/name sudah ada |
| `"invalid parameter"` | Parameter tidak valid |
| `"Gagal terhubung ke MikroTik"` | Connection failed |
| `"No session"` | Session tidak aktif |
| `"Permission denied"` | Tidak ada akses |

### Exception Handling

RouterOS API dapat return trap/error:

```php
if(!empty($response['!trap'][0]['message'])){
    $error = $response['!trap'][0]['message'];
    // Handle error
}
```

Common RouterOS Traps:
- `"duplicate name"` - Item sudah ada
- `"no such command"` - Command tidak dikenali
- `"invalid value for key"` - Value tidak valid
- `"connection refused"` - Koneksi ditolak

---

## 📌 Complete Request-Response Examples

### Complete Flow: Generate & Print Voucher

**Step 1: Generate 50 vouchers**

```http
POST /post/post_generate_voucher.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&qty=50&server=Hotspot-Server-1&user=vc&userl=8&prefix=&char=mix&profile=1day-profile&timelimit=1d&datalimit=1G&gcomment=Daily+100K&gencode=538
```

Response:
```json
{
  "message": "success",
  "data": {
    "count": "50",
    "comment": "vc-538-02.25.26-Daily 100K"
  }
}
```

**Step 2: Cache vouchers**

```http
POST /post/post_cache_voucher.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&qty=50&user=vc&gcomment=Daily+100K&gencode=538
```

Response:
```json
{
  "message": "success",
  "data": {
    "count": "50",
    "comment": "vc-538-02.25.26-Daily 100K"
  }
}
```

**Step 3: Get data for print**

```http
GET /index.php?user01/print_voucher HTTP/1.1
Host: localhost
Cookie: PHPSESSID=abc123xyz789
```

System retrieves from:
- `$_SESSION[$m_user.$comment]` - Cached user data
- Generate print view dengan username/password
- Show QR code untuk setiap voucher

---

### Complete Flow: Create User Profile & Add User

**Step 1: Create Profile**

```http
POST /post/post_add_userprofile.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&name=Premium-Daily&addresspool=Pool-Main&sharedusers=1&ratelimit=10M&parentqueue=QueueMain&expmode=rem&validity=1d&price=25000&sellingprice=20000&lockuser=Enable&lockserver=Disable
```

Response:
```json
{
  "message": "success",
  "data": {
    ".id": "*3",
    "name": "Premium-Daily",
    "on-login": ":put (\"rem,25000,1d,20000,,Enable,Disable,\"); ..."
  }
}
```

**Step 2: Add User dengan Profile**

```http
POST /post/post_add_user.php HTTP/1.1
Host: localhost
Content-Type: application/x-www-form-urlencoded
Cookie: PHPSESSID=abc123xyz789

sessname=?user01&server=Hotspot-Server-1&name=john123&password=john123&profile=Premium-Daily&macaddr=&timelimit=1d&datalimit=1G&comment=John+Doe
```

Response:
```json
{
  "message": "success",
  "data": {
    ".id": "*10",
    "name": "john123",
    "password": "john123",
    "profile": "Premium-Daily",
    "comment": "vc-john123-02.25.26-John Doe",
    "limit-uptime": "1d",
    "limit-bytes-total": "1073741824"
  }
}
```

---

## 🔐 Security Considerations

### Session Management

```php
// Session harus aktif untuk semua POST requests
if(!isset($_SESSION["mikhmon"])){
    die("No session");
}
```

### Input Validation

```php
// Sanitize input untuk mencegah injection
$name = preg_replace('/[^a-zA-Z0-9\-_]/', '', $_POST['name']);
```

### Encoding Sensitive Data

```php
// Router credentials di-encode 3x:
// 1. Base64 double
// 2. XOR dengan 10
// 3. Base64 lagi
```

---

**API Documentation Version**: 1.0  
**Last Updated**: 2026-02-25  
**Repository**: https://github.com/irhabi89/mikhmon_v4
