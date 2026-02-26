# Test MikroTik Client

File ini untuk melakukan test langsung ke MikroTik router tanpa melalui HTTP API.

## Cara Penggunaan

1. **Konfigurasi Router**
   
   Edit file `cmd/test_mikrotik/main.go` dan ubah konstanti di bagian atas:
   
   ```go
   const (
       RouterHost     = "192.168.88.1"  // Ganti dengan IP router Anda
       RouterPort     = 8728              // Port API MikroTik
       RouterUser     = "admin"             // Username MikroTik
       RouterPassword = "your_password"      // Ganti dengan password
   )
   ```

2. **Jalankan Test**
   
   ```bash
   cd /home/butterfly_student/code/Mikrotik/mikhmon_v4/backend
   go run cmd/test_mikrotik/main.go
   ```

## Test yang Dilakukan

File ini akan memanggil semua fungsi di `/internal/infrastructure/mikrotik/`:

### Hotspot Operations
- **GetHotspotServers** - Mengambil daftar server hotspot
- **GetHotspotUsers** - Mengambil daftar user hotspot
- **GetHotspotUsersCount** - Mengambil jumlah user hotspot
- **GetUserProfiles** - Mengambil daftar profil user
- **GetUserProfileByID** - Mengambil profil by ID
- **GetUserProfileByName** - Mengambil profil by nama
- **GetHotspotActive** - Mengambil sesi aktif
- **GetHotspotActiveCount** - Mengambil jumlah sesi aktif
- **GetHotspotHosts** - Mengambil daftar host hotspot

### System Operations
- **GetSystemResource** - Mengambil info resource (CPU, RAM, HDD)
- **GetSystemHealth** - Mengambil info kesehatan (Voltage, Temperature)
- **GetSystemIdentity** - Mengambil identity router
- **GetRouterBoardInfo** - Mengambil info routerboard
- **GetSystemClock** - Mengambil waktu sistem

### Network Operations
- **GetInterfaces** - Mengambil daftar interface
- **StartTrafficMonitorListen** - Monitoring traffic interface (3 detik)
- **GetNATRules** - Mengambil aturan NAT

### Queue Operations
- **GetAllQueues** - Mengambil semua simple queues
- **GetAllParentQueues** - Mengambil parent queues

### Other Operations
- **GetAddressPools** - Mengambil IP address pools
- **GetHotspotLogs** - Mengambil log hotspot (5 terakhir)
- **GetSalesReports** - Mengambil laporan penjualan
- **GetSalesReportsByDay** - Mengambil laporan penjualan per hari
- **StartPingListen** - Testing ping ke 8.8.8.8 (3x)
- **EnsureExpireMonitor** - Setup expire monitor

### Generators
- **GenerateBatch** - Generate 5 voucher codes
- **Generate (On-Login)** - Generate on-login script
- **Parse (On-Login)** - Parse on-login script

## Contoh Output

```
==========================================
MikroTik Client Direct Testing
==========================================
Router: 192.168.88.1:8728

--- Hotspot Users ---
Total Users: 15
First User: {Name: user1 Profile: default ...}

--- Hotspot Profiles ---
Total Profiles: 3
First Profile: Name=default ExpireMode=remc Price=5000

--- System Info ---
Resource: CPU=15% FreeMemory=512MB
Identity: MikroTik-Router
RouterBoard: Model=RB750Gr3 Serial=XXXXXX

--- Ping (Streaming) ---
Ping #1: 8.8.8.8 -> 24.50ms, Received: true
Ping #2: 8.8.8.8 -> 25.20ms, Received: true
Ping #3: 8.8.8.8 -> 23.80ms, Received: true

--- Voucher Generator ---
Generated 5 vouchers
First Voucher: Vq7x3k2m / Vq7x3k2m
```
