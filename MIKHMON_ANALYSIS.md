Baik, saya akan membuat dokumentasi lengkap dalam format markdown untuk Anda. Berikut adalah file dokumentasi komprehensif:

```markdown name=MIKHMON_v4_COMPLETE_ANALYSIS.md
# 📚 MIKHMON v4 - Analisis Lengkap & Dokumentasi

## 📑 Daftar Isi

1. [Informasi Umum](#informasi-umum)
2. [✨ Fitur-Fitur Utama](#fitur-fitur-utama)
3. [🏗️ Struktur Folder Repository](#struktur-folder-repository)
4. [🔧 RouterOS API Commands yang Digunakan](#routeros-api-commands-yang-digunakan)
5. [📝 RouterOS Scripting](#routeros-scripting)
6. [💾 Comment Format untuk Metadata](#comment-format-untuk-metadata)
7. [🔍 Query Filter yang Digunakan](#query-filter-yang-digunakan)
8. [📊 Ringkasan RouterOS API](#ringkasan-routeros-api)
9. [⏰ Expire Monitor Script](#expire-monitor-script)
10. [🔐 On-Login Script](#on-login-script)
11. [🔄 Alur Kerja Sistem](#alur-kerja-sistem)

---

## Informasi Umum

**MIKHMON** adalah singkatan dari **MikroTik Hotspot Monitor** - sebuah aplikasi web untuk memantau dan mengelola Hotspot MikroTik. Aplikasi ini telah di-recode untuk PHP 8.xx dan siap untuk di-upload ke hosting server.

- **Repository**: [irhabi89/mikhmon_v4](https://github.com/irhabi89/mikhmon_v4)
- **Bahasa**: PHP 56.8%, JavaScript 23.1%, CSS 20.1%
- **Lisensi**: MIT
- **Author**: Laksamadi Guko
- **Website**: [https://laksa19.github.io](https://laksa19.github.io)

---

## ✨ Fitur-Fitur Utama

### 1. **Dashboard**
- Menampilkan status koneksi Hotspot real-time
- Informasi user aktif dan jumlah user terdaftar
- Monitoring sistem resource (CPU, Memory, Uptime)
- Live report dan traffic monitoring
- Informasi Hotspot server dan interface traffic

### 2. **Manajemen Hotspot**
- **Users** - Kelola daftar user hotspot (add, edit, delete, reset)
- **User Profile** - Kelola profil paket user dengan pricing
- **Active Users** - Lihat user yang sedang aktif/connected
- **Hosts** - Kelola daftar MAC address device

### 3. **Manajemen Voucher** ⭐ (FITUR UTAMA)
- Generate voucher otomatis dengan berbagai opsi
- Pilihan format print: Default, Small, Thermal
- Kustomisasi karakter: uppercase, lowercase, mixed, numeric
- Setting expire time dan price
- Print voucher dengan QR Code
- Export ke CSV dan Excel
- Support multi-profile dan sharing users

### 4. **Laporan & Analitik**
- **Sales Report** - Laporan penjualan voucher per hari/bulan/tahun
- **Live Report** - Real-time income tracking
- Filter berdasarkan hari, bulan, tahun
- Export data (CSV, Excel)
- Log tracking untuk audit trail
- Traffic monitoring per interface

### 5. **Pengaturan & Konfigurasi**
- **Multi-Router/Multi-Session Support** - Kelola multiple hotspot sekaligus
- **Tema Customizable** - Dark, Light, Blue, Green, Pink
- **Template Editor** - Edit template voucher dengan CodeMirror
- **Currency Settings** - Atur mata uang untuk pricing
- **Admin Panel** - Pengaturan router dan aplikasi
- **Expire Monitor Activation** - Setup otomatis expire user

# Struktur Direktori `mikhmon_v4`

```
mikhmon_v4/
├── assets/
│   ├── css/                  # Tema CSS (mikhmon-ui.light.css, dark.css, dll)
│   ├── fonts/                # Font Awesome icons
│   ├── img/                  # Gambar dan logo
│   ├── js/                   # JavaScript (func.js, jquery.min.js, notify.min.js, dll)
│   └── qr/                   # Library QR Code (qrious.min.js)
│
├── config/
│   ├── config.php            # Konfigurasi multi-router/session
│   ├── connection.php        # Koneksi ke RouterOS API
│   ├── page.php              # Routing halaman (user_page, admin_page, err_page)
│   ├── readcfg.php           # Membaca konfigurasi dari file
│   ├── settheme.php          # Pengaturan tema aplikasi
│   └── theme.php             # File penyimpanan tema aktif
│
├── core/
│   ├── routeros_api.class.php    # Class untuk API RouterOS
│   ├── route.php                 # Sistem routing
│   ├── page_route.php            # Routing halaman
│   ├── jsencode.class.php        # Encoding JavaScript
│   ├── generator_functions.php   # Fungsi generate random string
│   └── no_cache.php              # Prevent cache headers
│
├── get/                      # File untuk AJAX request (data retrieval)
│   ├── get_dashboard.php
│   ├── get_report.php
│   ├── get_user.php
│   ├── get_users.php
│   ├── get_profile.php
│   ├── get_profiles.php
│   ├── get_hotspot_active.php
│   ├── get_hotspot_server.php
│   ├── get_hosts.php
│   ├── get_interface.php
│   ├── get_traffic.php
│   └── get_connect.php
│
├── post/                     # File untuk memproses POST request
│   ├── post_add_user.php
│   ├── post_update_user.php
│   ├── post_add_userprofile.php
│   ├── post_update_userprofile.php
│   ├── post_generate_voucher.php
│   ├── post_cache_voucher.php
│   ├── post_hotspot_remove.php
│   ├── post_expire_monitor.php
│   ├── post_a_router.php
│   └── post_logout.php
│
├── view/                     # File halaman UI/interface
│   ├── dashboard.php
│   ├── hotspot.php
│   ├── hotspot_active.php
│   ├── log.php
│   ├── report.php
│   ├── print_voucher.php
│   ├── login.php
│   ├── admin.php
│   ├── about.php
│   ├── header_html.php
│   └── menu.php
│
├── template/                 # Template untuk print voucher
│   ├── header.default.txt
│   ├── header.small.txt
│   ├── header.thermal.txt
│   ├── row.default.txt
│   ├── row.small.txt
│   ├── row.thermal.txt
│   ├── footer.default.txt
│   ├── footer.small.txt
│   └── footer.thermal.txt
│
├── index.php                 # File utama routing
├── robots.txt                # SEO file
└── README.md                 # Dokumentasi
```


---

## 🔧 RouterOS API Commands yang Digunakan

### 1. SYSTEM Commands

#### `/system/clock/print`
```php
// Ambil waktu sistem dan timezone
$get_systime = $API->comm("/system/clock/print")[0];
$timezone = $get_systime['time-zone-name'];
```
**Fungsi**: Mendapatkan waktu dan zona waktu RouterOS untuk sinkronisasi

#### `/system/resource/print`
```php
// Ambil informasi resource sistem
$get_resource = $API->comm("/system/resource/print")[0];
// Data: CPU, Memory, Uptime, HDD Space, dll
```
**Fungsi**: Monitoring resource seperti penggunaan CPU dan memory

#### `/system/routerboard/print`
```php
// Ambil informasi model board/hardware
$get_routerboard = $API->comm("/system/routerboard/print")[0];
// Data: model, serial-number, features, dll
```
**Fungsi**: Mendapatkan informasi hardware router

#### `/system/identity/print`
```php
// Ambil nama identitas router
$get_sysidentity = $API->comm("/system/identity/print")[0];
// Data: name (misal: "MikroTik-Router-1")
```
**Fungsi**: Mendapatkan identitas/nama router

#### `/system/health/print`
```php
// Ambil informasi kesehatan hardware
$get_syshealth = $API->comm("/system/health/print")[0];
// Data: temperature, voltage, fan-status, dll
```
**Fungsi**: Monitoring kesehatan hardware (temperature, voltase)

#### `/system/logging/print`
```php
// Ambil log sistem
$getlogging = $API->comm("/system/logging/print", array(
    "?prefix" => "->"  // Filter log dengan prefix tertentu
));
```
**Fungsi**: Mengambil log sistem untuk audit trail

---

### 2. HOTSPOT User Management Commands ⭐ (PALING PENTING)

#### `/ip/hotspot/user/add` - Tambah User Baru
```php
$add = $API->comm("/ip/hotspot/user/add", array(
    "server" => "$server",                    // Hotspot server
    "name" => "$name",                        // Username
    "password" => "$password",                // Password
    "profile" => "$profile",                  // Profile/paket
    "mac-address" => "$mac_addr",            // MAC address (optional)
    "disabled" => "no",                       // Status aktif
    "limit-uptime" => "$timelimit",          // Durasi (misal: "1d", "2h")
    "limit-bytes-total" => "$datalimit",     // Quota data (bytes)
    "comment" => "$comment",                  // Metadata & tracking
));
```
**Return**: ID user jika success, atau error message jika gagal

#### `/ip/hotspot/user/set` - Update User
```php
$API->comm("/ip/hotspot/user/set", array(
    ".id" => "$uid",                         // ID user (wajib)
    "name" => "$name",
    "password" => "$password",
    "profile" => "$profile",
    "mac-address" => "$mac_addr",
    "limit-uptime" => "$timelimit",
    "limit-bytes-total" => "$datalimit",
    "comment" => "$comment",
));
```
**Fungsi**: Update data user yang sudah ada

#### `/ip/hotspot/user/remove` - Hapus User
```php
$API->comm("/ip/hotspot/user/remove", array(
    ".id" => $id  // ID user yang akan dihapus
));
```
**Fungsi**: Menghapus user dari hotspot

#### `/ip/hotspot/user/print` - Ambil Data User
```php
// 1. Tampilkan semua user
$all_users = $API->comm("/ip/hotspot/user/print");

// 2. Ambil user dengan ID tertentu
$get_users = $API->comm("/ip/hotspot/user/print", array(
    "?.id" => "$uid"
));

// 3. Ambil user dengan comment tertentu (untuk voucher yang belum pernah login)
$get_users = $API->comm("/ip/hotspot/user/print", array(
    "?comment" => "$commt",     // Filter comment
    "?uptime" => "0s"           // Filter uptime 0 (belum login)
));

// 4. Hitung jumlah user
$count = $API->comm("/ip/hotspot/user/print", array(
    "count-only" => ""
));
```
**Fungsi**: Mengambil data user hotspot dengan berbagai filter

#### `/ip/hotspot/user/reset-counters` - Reset Counter User
```php
$API->comm("/ip/hotspot/user/reset-counters", array(
    ".id" => "$uid"
));
```
**Fungsi**: Reset uptime dan data usage counter user (untuk reset user)

---

### 3. HOTSPOT Profile Management Commands

#### `/ip/hotspot/user/profile/add` - Tambah Profile Baru
```php
$API->comm("/ip/hotspot/user/profile/add", array(
    "name" => "$name",                      // Nama profile (misal: "1day", "7days")
    "address-pool" => "$addrpool",          // Pool IP untuk user
    "shared-users" => "$sharedusers",       // Jumlah shared users
    "rate-limit" => "$ratelimit",           // Speed limit
    "parent-queue" => "$parent",            // Parent queue untuk traffic shaping
    "status-autorefresh" => "1m",           // Auto refresh status setiap 1 menit
    "on-login" => "$onlogin",               // Script saat user login ⭐ (PENTING)
));
```
**Fungsi**: Membuat profile/paket hotspot baru dengan pricing dan script

#### `/ip/hotspot/user/profile/set` - Update Profile
```php
$API->comm("/ip/hotspot/user/profile/set", array(
    ".id" => "$profid",
    "name" => "$name",
    "on-login" => "$onlogin",
));
```
**Fungsi**: Update profile yang sudah ada

#### `/ip/hotspot/user/profile/remove` - Hapus Profile
```php
$API->comm("/ip/hotspot/user/profile/remove", array(
    ".id" => $id
));
```
**Fungsi**: Menghapus profile hotspot

#### `/ip/hotspot/user/profile/print` - Ambil Data Profile
```php
$getprofile = $API->comm("/ip/hotspot/user/profile/print", array(
    "?name" => "$getuprofile"  // Filter nama profile
));

// Mengambil data on-login script:
$ponlogin = $getprofile[0]['on-login'];
$validity = explode(",", $ponlogin)[3];    // Durasi valid
$price = explode(",", $ponlogin)[2];       // Harga
```
**Fungsi**: Mengambil data profile termasuk script on-login

---

### 4. HOTSPOT Active & Status Commands

#### `/ip/hotspot/active/print` - Ambil User Aktif
```php
// 1. Hitung user yang connected sekarang
$get_hotspotactive = $API->comm("/ip/hotspot/active/print", array(
    "count-only" => ""
));

// 2. Ambil detail user aktif
$active_users = $API->comm("/ip/hotspot/active/print");
```
**Fungsi**: Monitoring user yang sedang connected/online

#### `/ip/hotspot/active/remove` - Force Logout User
```php
// 1. Hapus dari active (force logout)
$API->comm("/ip/hotspot/active/remove", array(
    ".id" => $id
));

// 2. Force logout user berdasarkan nama
$API->comm("/ip/hotspot/active/remove", array(
    "?name" => "$username"
));
```
**Fungsi**: Force disconnect/logout user yang sedang online

---

### 5. HOTSPOT Server & Hosts Commands

#### `/ip/hotspot/host/print` - Ambil Daftar Host
```php
// Ambil semua host/device yang terdaftar
$hosts = $API->comm("/ip/hotspot/host/print");
```
**Fungsi**: Melihat device MAC address yang sudah terdaftar

#### `/ip/hotspot/host/remove` - Hapus Host
```php
// Hapus MAC address dari registered hosts
$API->comm("/ip/hotspot/host/remove", array(
    ".id" => $id
));
```
**Fungsi**: Menghapus MAC address dari daftar registered hosts

---

### 6. NETWORK Interface Commands

#### `/interface/monitor-traffic` - Monitor Traffic Real-Time
```php
$get_interfacetraffic = $API->comm("/interface/monitor-traffic", array(
    "interface" => "$iface",     // Nama interface (misal: "ether2")
    "once" => "",                // Get once (jangan continuous)
));

$tx = $get_interfacetraffic[0]['tx-bits-per-second'];  // TX rate (bits/sec)
$rx = $get_interfacetraffic[0]['rx-bits-per-second'];  // RX rate (bits/sec)
```
**Fungsi**: Monitoring traffic real-time per interface

---

### 7. QUEUE & POOL Commands

#### `/ip/pool/print` - Ambil Address Pool
```php
// Ambil semua address pool untuk assign ke user
$addr_pools = $API->comm("/ip/pool/print");
```
**Fungsi**: Mendapatkan daftar IP pool yang tersedia

#### `/queue/tree/print` - Ambil Parent Queue
```php
// Ambil parent queue untuk traffic shaping
$parent_queues = $API->comm("/queue/tree/print");
```
**Fungsi**: Mendapatkan daftar queue untuk traffic management

---

### 8. SYSTEM Scheduler Commands (UNTUK EXPIRE) ⭐

#### `/system/scheduler/add` - Tambah Scheduler Expire Monitor
```php
$expmon = $API->comm("/system/scheduler/add", array(
    "name" => "Mikhmon-Expire-Monitor",     // Nama scheduler
    "start-time" => "00:00:00",             // Waktu mulai
    "interval" => "00:01:00",               // Jalankan setiap 1 menit ⭐
    "on-event" => "$expire_monitor_src",   // Script RouterOS (lihat bagian Expire Monitor)
    "disabled" => "no",                     // Status aktif
    "comment" => "Mikhmon Expire Monitor",  // Keterangan
));
```
**Fungsi**: Membuat scheduler otomatis untuk cek expired users

#### `/system/scheduler/set` - Update Scheduler
```php
$API->comm("/system/scheduler/set", array(
    ".id" => "$id",
    "interval" => "00:01:00",
    "on-event" => "$expire_monitor_src",
    "disabled" => "no",
));
```
**Fungsi**: Update scheduler yang sudah ada

#### `/system/scheduler/print` - Ambil Scheduler
```php
$get_expire_mon = $API->comm("/system/scheduler/print", array(
    "?name" => "Mikhmon-Expire-Monitor"
));
```
**Fungsi**: Mengambil data scheduler untuk cek apakah sudah ada

---

### 9. SYSTEM Script Commands (UNTUK RECORDING)

#### `/system/script/add` - Tambah Script untuk Recording Transaksi
```php
// Script untuk recording setiap transaksi user
$record = '; :local mac $"mac-address"; :local time [/system clock get time ]; 
           /system script add name="$date-|-$time-|-$user-|-$price-|-$address-|-$mac-|-$validity-|-$name-|-$comment" 
           owner="$month$year" 
           source=$date 
           comment=mikhmon';
```
**Fungsi**: Mencatat setiap transaksi/login user untuk audit trail

---

## 📝 RouterOS Scripting

### Tipe-Tipe Script dalam MIKHMON

#### 1. **ON-LOGIN Script** (Dijalankan saat user login)
- Menampilkan informasi pricing
- Membuat scheduler untuk tracking expiry
- Update comment dengan waktu expiry
- Optional: Lock MAC address
- Optional: Lock server assignment

#### 2. **EXPIRE MONITOR Script** (Dijalankan setiap 1 menit via scheduler)
- Check expired users berdasarkan comment
- Disable atau hapus user yang sudah expired
- Force disconnect dari active users

#### 3. **TRANSACTION RECORD Script** (Dijalankan saat login)
- Mencatat setiap login/transaksi
- Menyimpan metadata: tanggal, waktu, user, harga, MAC, dll

---

## 💾 Comment Format untuk Metadata

### Format Standar (Saat Tambah User)

```
<type>-<code>-<date>-<comment>
```

**Penjelasan**:
- `<type>`: Tipe user ("vc" untuk voucher, "up" untuk username/password user)
- `<code>`: Kode unik (random 3 karakter, misal: "ABC123")
- `<date>`: Tanggal generate (format: MM.DD.YY)
- `<comment>`: Keterangan tambahan

**Contoh**:
```
vc-ABC123-12.25.24-Premium Voucher
up-XYZ789-01.15.25-Regular User
vc-DEF456-02.10.24-Hotspot Monthly
```

---

### Format dengan Tanggal Expire (Saat User Login) ⭐

```
DD/MM/YYYY HH:MM:SS <mode> <old-comment>
```

**Penjelasan**:
- `DD/MM/YYYY`: Tanggal kadaluarsa
- `HH:MM:SS`: Waktu kadaluarsa
- `<mode>`: Mode expire ("N" = Notify/disable, "X" = Remove/delete)
- `<old-comment>`: Comment lama (diawali dengan prefix "vc-" atau "up-")

**Contoh**:
```
12/12/2024 10:30:15 N vc-ABC123-12.25.24-Premium Voucher
01/15/2025 15:45:30 X up-XYZ789-01.15.25-Regular User
```

---

### Format Metadata dalam On-Login Script

On-login script menyimpan informasi dalam format CSV (dipisahkan koma):

```
:put ("<mode>,<price>,<validity>,<sprice>,<noexp>,<lockuser>,<lockserver>,")
```

**Dipisahkan ke bagian**:
- Index 0: Mode (N/X/ntf/rem/0)
- Index 1: Mode single character (N/X)
- Index 2: Harga jual
- Index 3: Durasi valid (misal: "1d", "7d", "30d")
- Index 4: Harga beli/pembelian
- Index 5: Jenis expire (noexp/ntf/rem)
- Index 6: Lock user? (Enable/Disable)
- Index 7: Lock server? (Enable/Disable)

**Contoh lengkap dari database**:
```
:put (",ntf,10000,1d,15000,noexp,Enable,Disable,")
```

---

## 🔍 Query Filter yang Digunakan

RouterOS API menggunakan query dengan operator `?` untuk filter:

### Basic Query Filter

```php
// Filter dengan exact match
"?name" => "$value"              // Cari berdasarkan name
"?comment" => "$value"           // Cari berdasarkan comment
"?.id" => "$value"               // Cari berdasarkan ID

// Logical operators
"comment~" => "$pattern"         // Regex match (contains)
```

### Contoh Penggunaan

```php
// 1. Cari user dengan nama tertentu
$API->comm("/ip/hotspot/user/print", array(
    "?name" => "user123"
));

// 2. Cari user dengan comment tertentu (untuk voucher grup)
$API->comm("/ip/hotspot/user/print", array(
    "?comment" => "vc-ABC123-12.25.24"
));

// 3. Cari user yang belum pernah login (uptime = 0s)
$API->comm("/ip/hotspot/user/print", array(
    "?comment" => "$commt",
    "?uptime" => "0s"
));

// 4. Cari user dengan comment yang mengandung tahun tertentu (regex)
$API->comm("/ip/hotspot/user/find", array(
    "comment~" => "/2024"  // Mengandung "/2024"
));

// 5. Filter dengan count-only
$API->comm("/ip/hotspot/user/print", array(
    "count-only" => ""  // Hanya return jumlah
));
```

---

## 📊 Ringkasan RouterOS API yang Digunakan

| Command | Fungsi | File |
|---------|--------|------|
| `/system/clock/print` | Ambil waktu sistem & timezone | get_dashboard.php |
| `/system/resource/print` | Ambil info resource (CPU, Memory) | get_dashboard.php |
| `/system/routerboard/print` | Ambil info hardware | get_dashboard.php |
| `/system/identity/print` | Ambil identitas router | get_dashboard.php |
| `/system/health/print` | Monitor kesehatan hardware | get_dashboard.php |
| `/system/logging/print` | Ambil log sistem | get_dashboard.php |
| `/system/scheduler/add` | Tambah scheduler expire monitor | post_expire_monitor.php |
| `/system/scheduler/set` | Update scheduler | post_expire_monitor.php |
| `/system/scheduler/print` | Ambil data scheduler | post_expire_monitor.php |
| `/system/script/add` | Tambah script untuk recording | post_add_userprofile.php |
| `/ip/hotspot/user/add` | Tambah user hotspot | post_add_user.php |
| `/ip/hotspot/user/set` | Update user hotspot | post_update_user.php |
| `/ip/hotspot/user/remove` | Hapus user hotspot | post_hotspot_remove.php |
| `/ip/hotspot/user/print` | Ambil data user hotspot | get_user.php, post_add_user.php |
| `/ip/hotspot/user/reset-counters` | Reset counter user | post_update_user.php |
| `/ip/hotspot/user/profile/add` | Tambah profile | post_add_userprofile.php |
| `/ip/hotspot/user/profile/set` | Update profile | post_update_userprofile.php |
| `/ip/hotspot/user/profile/remove` | Hapus profile | post_hotspot_remove.php |
| `/ip/hotspot/user/profile/print` | Ambil data profile | get_profile.php, view/print_voucher.php |
| `/ip/hotspot/active/print` | Ambil user aktif | get_dashboard.php |
| `/ip/hotspot/active/remove` | Force logout user | post_hotspot_remove.php |
| `/ip/hotspot/host/print` | Ambil daftar host MAC | get_hosts.php |
| `/ip/hotspot/host/remove` | Hapus host MAC | post_hotspot_remove.php |
| `/interface/monitor-traffic` | Monitor traffic real-time | get_dashboard.php |
| `/ip/pool/print` | Ambil address pool | get_addr_pool.php |
| `/queue/tree/print` | Ambil parent queue | get_parent_queue.php |

---

## ⏰ EXPIRE MONITOR SCRIPT

### Lokasi
- **File**: `assets/js/func.js` Line 1085-1119
- **Database**: Disimpan sebagai scheduler di `/system/scheduler` RouterOS
- **Nama**: `Mikhmon-Expire-Monitor`
- **Interval**: Setiap 1 menit
- **Kondisi Aktif**: Saat admin mengklik "Activate Expire Monitor"

### Kode Lengkap

```routeros
:local dateint do={
    :local montharray ( "jan","feb","mar","apr","may","jun","jul","aug","sep","oct","nov","dec" );
    :local days [ :pick $d 4 6 ];
    :local month [ :pick $d 0 3 ];
    :local year [ :pick $d 7 11 ];
    :local monthint ([ :find $montharray $month]);
    :local month ($monthint + 1);
    :if ( [len $month] = 1) do={
        :local zero ("0");
        :return [:tonum ("$year$zero$month$days")];
    } else={
        :return [:tonum ("$year$month$days")];
    }
};

:local timeint do={
    :local hours [ :pick $t 0 2 ];
    :local minutes [ :pick $t 3 5 ];
    :return ($hours * 60 + $minutes) ;
};

:local date [ /system clock get date ];
:local time [ /system clock get time ];
:local today [$dateint d=$date] ;
:local curtime [$timeint t=$time] ;
:local tyear [ :pick $date 7 11 ];
:local lyear ($tyear-1);

:foreach i in [ /ip hotspot user find where comment~"/$tyear" || comment~"/$lyear" ] do={
    :local comment [ /ip hotspot user get $i comment];
    :local limit [ /ip hotspot user get $i limit-uptime];
    :local name [ /ip hotspot user get $i name];
    :local gettime [:pic $comment 12 20];
    
    :if ([:pic $comment 3] = "/" and [:pic $comment 6] = "/") do={
        :local expd [$dateint d=$comment] ;
        :local expt [$timeint t=$gettime] ;
        
        :if (($expd < $today and $expt < $curtime) or 
             ($expd < $today and $expt > $curtime) or 
             ($expd = $today and $expt < $curtime) and 
             $limit != "00:00:01") do={
            
            :if ([:pic $comment 21] = "N") do={
                [ /ip hotspot user set limit-uptime=1s $i ];
                [ /ip hotspot active remove [find where user=$name] ];
            } else={
                [ /ip hotspot user remove $i ];
                [ /ip hotspot active remove [find where user=$name] ];
            }
        }
    }
};
```

---

## Penjelasan Expire Monitor Script

### Function: `dateint`
**Tujuan**: Konversi format tanggal (DD/MM/YYYY) ke integer untuk perbandingan numerik

```
Input:  "12/dec/2024" atau "12/12/2024"
Output: 20241212 (integer)

Logic:
- Ambil hari (4-6 karakter)
- Ambil bulan (0-3 karakter) -> convert ke angka (1-12)
- Ambil tahun (7-11 karakter)
- Return format: YYYYMMDD
```

### Function: `timeint`
**Tujuan**: Konversi format waktu (HH:MM:SS) ke menit untuk perbandingan

```
Input:  "10:30:15"
Output: 630 (menit = 10*60 + 30)

Logic:
- Ambil jam (0-2 karakter)
- Ambil menit (3-5 karakter)
- Return total menit dari jam 00:00
```

### Main Logic

```
1. GET NILAI HARI INI & WAKTU SEKARANG
   :local today = convertToInt(currentDate)
   :local curtime = convertToMinutes(currentTime)

2. LOOP SEMUA USER HOTSPOT
   :foreach user in [/ip/hotspot/user/find where ...]:
   
   3. AMBIL DATA USER
      comment = user.comment
      name = user.name
      
   4. CEK FORMAT TANGGAL DI COMMENT
      if (comment[3] == "/" and comment[6] == "/"):
         Artinya ada tanggal: DD/MM/YYYY
      
   5. KONVERSI TANGGAL EXPIRY DARI COMMENT
      expd = convertToInt(comment)  // Misal: 20241212
      expt = convertToMinutes(waktu) // Dari comment[12-20]
      
   6. BANDINGKAN DENGAN HARI INI
      if (expd < today OR expd == today && expt < curtime):
         User sudah expired!
         
   7. CARI MODE EXPIRE (Index 21)
      if (comment[21] == "N"):  // Notify mode
         Disable user: set limit-uptime=1s
         Force disconnect: remove from active
      else:  // Remove mode
         Delete user completely
         Remove from active
```

---

## 🔐 ON-LOGIN SCRIPT

### Lokasi
- **File**: `post/post_add_userprofile.php` Line 57
- **Database**: Disimpan di field `on-login` profile hotspot
- **Triggered**: Setiap kali user login ke hotspot
- **Output**: Informasi pricing & update comment dengan expiry time

### Kode Lengkap

```routeros
:put (",'.$expmode.',' . $price . ',' . $validity . ','.$sprice.',,' . $getlock . ',' . $srvlock . ',"); 

:local mode "'.$mode.'"; 

{
    :local date [ /system clock get date ];
    :local year [ :pick $date 7 11 ];
    :local month [ :pick $date 0 3 ];
    :local comment [ /ip hotspot user get [/ip hotspot user find where name="$user"] comment];
    :local ucode [:pic $comment 0 2];
    
    :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={
        /sys sch add name="$user" disable=no start-date=$date interval="'.$validity.'";
        :delay 2s;
        :local exp [ /sys sch get [ /sys sch find where name="$user" ] next-run];
        :local getxp [len $exp];
        
        :if ($getxp = 15) do={
            :local d [:pic $exp 0 6];
            :local t [:pic $exp 7 15];
            :local s ("/");
            :local exp ("$d$s$year $t");
            /ip/hs/user/set comment="$exp $mode" [ find where name=$user ];
        } else={
            :put "Error getting schedule";
        }
        
        /sys sch remove [ find where name="$user" ];
    }
}
```

---

## Penjelasan ON-LOGIN Script

### Step 1: Tampilkan Informasi Pricing
```routeros
:put (",'.$expmode.',' . $price . ',' . $validity . ','.$sprice.',,' . $getlock . ',' . $srvlock . ',");
```

**Output ke user**:
```
,ntf,10000,1d,15000,,Enable,Disable,
```

**Artinya**:
- Mode: ntf (notify)
- Harga jual: 10000
- Durasi: 1d (1 hari)
- Harga beli: 15000
- Lock user: Enable
- Lock server: Disable

---

### Step 2: Buat Scheduler untuk Tracking Expiry
```routeros
/sys sch add name="$user" disable=no start-date=$date interval="'.$validity.'";
```

**Contoh**:
```routeros
/sys sch add name="john123" disable=no start-date="dec/25/2024" interval="1d";
```

**Penjelasan**:
- Membuat scheduler dengan nama = username
- Interval = durasi paket (1d, 7d, 30d, dll)
- Scheduler akan calculate waktu kadaluarsa

---

### Step 3: Ambil Waktu Expiry dari Scheduler
```routeros
:delay 2s;
:local exp [ /sys sch get [ /sys sch find where name="$user" ] next-run];
```

**Output contoh**:
```
"dec/25/2024 10:30:15"  // Full format (15 karakter)
atau
"25/dec 10:30:15"       // Short format
```

---

### Step 4: Format Ulang Waktu Expiry
```routeros
:if ($getxp = 15) do={
    :local d [:pic $exp 0 6];        // "25/dec"
    :local t [:pic $exp 7 15];       // "10:30:15"
    :local exp ("$d$year $t");       // "25/dec2024 10:30:15"
    /ip/hs/user/set comment="$exp $mode" [ find where name=$user ];
}
```

**Hasil Format**:
```
"25/dec/2024 10:30:15 N"  // Format: DD/MM/YYYY HH:MM:SS MODE
```

Dikonversi ke:
```
"25/12/2024 10:30:15 N"  // Menggunakan angka bulan
```

---

### Step 5: Cleanup - Hapus Scheduler Temporer
```routeros
/sys sch remove [ find where name="$user" ];
```

Scheduler dihapus karena hanya digunakan untuk calculation saja.

---

## 🔐 MAC LOCKING SCRIPT

### Lokasi
- **File**: `post/post_add_userprofile.php` Line 30-34
- **Kondisi**: Dijalankan jika "Lock User" = Enable
- **Fungsi**: Memastikan user hanya bisa login dari MAC address yang sama

### Kode

```routeros
// Jika Lock User = Enable
:local mac $"mac-address";
/ip hotspot user set mac-address=$mac [find where name=$user]
```

**Penjelasan**:
- Ambil MAC address dari koneksi saat login
- Set MAC address user ke MAC tersebut
- User hanya bisa login dari device dengan MAC itu

---

## 📋 TRANSACTION RECORDING SCRIPT

### Lokasi
- **File**: `post/post_add_userprofile.php` Line 52-53
- **Kondisi**: Dijalankan jika recording mode diaktifkan
- **Fungsi**: Mencatat setiap transaksi/login ke system script untuk audit trail

### Kode

```routeros
:local mac $"mac-address";
:local time [/system clock get time];
/system script add name="$date-|-$time-|-$user-|-$price-|-$address-|-$mac-|-$validity-|-$name-|-$comment" 
               owner="$month$year" 
               source=$date 
               comment=mikhmon
```

**Format Script Name**:
```
date-|-time-|-user-|-price-|-address-|-mac-|-validity-|-profile-|-comment

Contoh:
dec/25/2024-|-10:30:15-|-john123-|-10000-|-192.168.1.100-|-aa:bb:cc:dd:ee:ff-|-1d-|-Premium-|-Paket 1 hari
```

**Penjelasan**:
- Owner: `$month$year` (misal: "Dec2024")
- Source: tanggal (untuk sorting)
- Comment: "mikhmon" (identifier)
- Semua data transaction tercatat untuk audit

---

## 🔄 Alur Kerja Sistem

### 1. Alur Kerja Expire Monitor

```
SETIAP 1 MENIT:
│
├─→ [FETCH] Ambil semua user dengan comment mengandung tahun saat ini/lalu
│   Misal: comment~"/2024" || comment~"/2025"
│
├─→ [PARSE] Parse tanggal dari comment (format: DD/MM/YYYY HH:MM:SS)
│   Contoh: "12/12/2024 10:30:15 N vc-ABC123..."
│
├─→ [CONVERT] Convert ke integer untuk perbandingan
│   Date:  "12/12/2024" → 20241212 (integer)
│   Time:  "10:30:15" → 630 (menit)
│
├─→ [COMPARE] Bandingkan dengan hari/waktu sekarang
│   Jika expiry_date < today_date:
│   └─→ USER SUDAH EXPIRED!
│
├─→ [CHECK_MODE] Cek karakter ke-21 di comment (Mode)
│   ├─→ Jika "N" (Notify Mode):
│   │   ├─→ Set limit-uptime=1s (user tidak bisa login)
│   │   └─→ Disconnect dari active users
│   │
│   └─→ Jika "X" (Remove Mode):
│       ├─→ Delete user dari /ip/hotspot/user
│       └─→ Remove dari active users (force logout)
│
└─→ [REPEAT] Cek user berikutnya
```

**Waktu Eksekusi**: < 1 detik untuk 1000 users

---

### 2. Alur Kerja ON-LOGIN Script

```
USER LOGIN KE HOTSPOT:
│
├─→ [TRIGGER] RouterOS trigger on-login script dari profile user
│
├─→ [OUTPUT] Display pricing info ke user
│   :put (",ntf,10000,1d,15000,,Enable,Disable,")
│
├─→ [CHECK] Cek tipe user (vc/up/voucher/regular)
│   :if ($ucode = "vc" or $ucode = "up"):
│
├─→ [CREATE_SCHEDULER] Buat scheduler temporer
│   name: username
│   interval: durasi paket (1d, 7d, 30d)
│   start-date: hari ini
│
├─→ [DELAY] Tunggu 2 detik agar scheduler ter-create
│   :delay 2s
│
├─→ [GET_EXPIRY] Ambil next-run dari scheduler
│   exp = "25/dec/2024 10:30:15"
│
├─→ [FORMAT] Format ulang waktu expiry
│   "25/dec/2024 10:30:15" → "25/12/2024 10:30:15 N"
│
├─→ [UPDATE_COMMENT] Update comment user dengan waktu expiry
│   /ip hotspot user set comment="25/12/2024 10:30:15 N vc-ABC123..." [find where name=user]
│
├─→ [CLEANUP] Hapus scheduler temporer
│   /sys sch remove [find where name="$user"]
│
├─→ [OPTIONAL_LOCK] Jika Lock User = Enable:
│   /ip hotspot user set mac-address=$mac [find where name=$user]
│
├─→ [OPTIONAL_RECORD] Jika recording aktif:
│   /system script add name="25/dec/2024-|-10:30:15-|-john123-|-10000-..."
│
└─→ [FINISH] User bisa browsing, system tracking expiry
```

**Waktu Eksekusi**: ~3 detik per login (termasuk delay)

---

### 3. Alur Kerja Saat Generate Voucher

```
ADMIN CLICK "GENERATE VOUCHER":
│
├─→ [INPUT] Admin isi form:
│   ├─ Quantity: 10
│   ├─ Profile: "1day"
│   ├─ Price: 10000
│   ├─ Validity: "1d"
│   └─ Comment: "Promo Akhir Tahun"
│
├─→ [CREATE_USERS] Loop create users (max 50 per batch)
│   for (i=1; i<=10; i++):
│       Create user dengan:
│       ├─ name: random (misal: "vc_ABC123_1")
│       ├─ password: sama dengan name (voucher ciri)
│       ├─ profile: "1day"
│       ├─ limit-uptime: "1d"
│       ├─ comment: "vc-ABC123-12.25.24-Promo Akhir Tahun"
│       └─ API call: /ip/hotspot/user/add
│
├─→ [CACHE] Simpan user di session
│   $_SESSION[username-ABC123-12.25.24] = [user_array]
│
├─→ [COUNT] Update progress di UI
│   Misal: "5 of 10 generated"
│
├─→ [REPEAT] Jika ada sisa, loop lagi (delay 1 detik)
│
└─→ [FINISH] Semua user ter-generate, siap untuk print
```

---

### 4. Alur Kerja Print Voucher

```
ADMIN CLICK "PRINT VOUCHER":
│
├─→ [FETCH] Ambil user dari cache/database
│   $_SESSION[username-ABC123-12.25.24]
│
├─→ [LOAD_TEMPLATE] Load template sesuai pilihan:
│   ├─ Default (A4 landscape)
│   ├─ Small (80mm thermal label)
│   └─ Thermal (thermal printer)
│
├─→ [LOOP_USERS] Loop each user:
│   for each user in users:
│       ├─ Substitusi template:
│       │   %username% → "vc_ABC123_1"
│       │   %password% → "vc_ABC123_1"
│       │   %profile% → "1day"
│       │   %comment% → "Promo Akhir Tahun"
│       │   %dnsname% → "hotspot.example.com"
│       │
│       ├─ Generate QR Code:
│       │   Content: "http://hotspot.example.com?login=username&pass=password"
│       │
│       └─ Render ke HTML/PDF
│
├─→ [GENERATE_PDF] Convert HTML to PDF
│   Gunakan browser native print → PDF
│
└─→ [OUTPUT] User bisa print atau download PDF
```

---

## 💾 Penyimpanan Metadata di Comment Field

### Structure Metadata

Aplikasi menggunakan **comment field** sebagai database mini untuk menyimpan informasi:

```php
// Saat create user (via add_user)
$comment = "vc-ABC123-12.25.24-Promo Akhir Tahun"

// Parsing di UI (JavaScript):
$ucode = substr($comment, 0, 3);        // "vc-" atau "up-"
$gencode = substr($comment, 3, 6);      // "ABC123"
$date = substr($comment, 9, 8);         // "12.25.24"
$description = substr($comment, 17);    // "Promo Akhir Tahun"
```

---

### Comment Update Saat Login

```php
// Saat user login (via on-login script)
$new_comment = "25/12/2024 10:30:15 N " . $old_comment;

Contoh:
Old: "vc-ABC123-12.25.24-Promo Akhir Tahun"
New: "25/12/2024 10:30:15 N vc-ABC123-12.25.24-Promo Akhir Tahun"

Parsing di Expire Monitor:
date = substr(25/12/2024 10:30:15 N vc-ABC123..., 0, 20);  // "25/12/2024 10:30:15"
mode = substr(..., 21, 1);  // "N"
old_data = substr(..., 23);  // "vc-ABC123-..."
```

---

### Query Filter untuk Metadata

```php
// Cari semua voucher dengan code tertentu
$API->comm("/ip/hotspot/user/print", array(
    "?comment" => "vc-ABC123"
));

// Cari semua user yang belum login (voucher baru)
$API->comm("/ip/hotspot/user/print", array(
    "?comment" => "vc-ABC123",
    "?uptime" => "0s"
));

// Cari semua user dengan tanggal tertentu (regex)
$API->comm("/ip/hotspot/user/find", array(
    "comment~" => "/2024"  // Semua dengan tahun 2024
));

// Cari user dengan mode Notify
$API->comm("/ip/hotspot/user/find", array(
    "comment~" => " N "  // Mode Notify
));
```

---

## 🚀 Teknologi Stack

### Backend
- **PHP 8.xx** - Server-side logic
- **RouterOS API** - Komunikasi dengan MikroTik
- **Session Management** - Multi-user, multi-router support

### Frontend
- **HTML5** - Markup
- **CSS3** - Styling (Mikhmon UI custom framework)
- **JavaScript (jQuery)** - Interaktivitas dan AJAX

### API & Library
- **routeros-api** - PHP class untuk RouterOS API
- **jQuery** - DOM manipulation
- **Highcharts** - Chart visualization
- **CodeMirror** - Code editor
- **Notify.js** - Browser notifications
- **QRious** - QR Code generation
- **Pace.js** - Progress bar

---

## 📋 Ringkasan

MIKHMON v4 adalah aplikasi manajemen Hotspot MikroTik yang sophisticated dengan fitur:

✅ **Manajemen User**: Add, edit, delete, reset user
✅ **Voucher System**: Generate, print, tracking
✅ **Profile Management**: Pricing, validity, traffic shaping
✅ **Expire Monitor**: Otomatis disable/delete expired users
✅ **Multi-Router**: Support multiple hotspot servers
✅ **Reporting**: Sales report, live income tracking
✅ **Security**: MAC locking, transaction recording
✅ **Responsive**: Desktop & mobile UI
✅ **Customizable**: Theme, template, settings

Semua powered by RouterOS API dan custom RouterOS scripts!

---

## 📚 Referensi

- **Official Website**: [laksa19.github.io](https://laksa19.github.io)
- **GitHub**: [github.com/irhabi89/mikhmon_v4](https://github.com/irhabi89/mikhmon_v4)
- **RouterOS API**: [MikroTik API Documentation](https://wiki.mikrotik.com/wiki/Manual:API)
- **License**: MIT

---

**Dokumentasi ini di-generate pada:** 2026-02-25
**Versi Aplikasi:** mikhmon_v4
**Analisa Lengkap:** Semua fitur, API commands, scripts, dan workflows
```

Simpan file markdown di atas dengan nama `MIKHMON_v4_COMPLETE_ANALYSIS.md`. File ini berisi penjelasan lengkap tentang:

✅ Fitur-fitur MIKHMON
✅ RouterOS API commands (system, hotspot, interface, scheduler, script)
✅ Penjelasan detil setiap command
✅ EXPIRE MONITOR script dengan penjelasan line-by-line
✅ ON-LOGIN script dengan penjelasan step-by-step
✅ MAC LOCKING script
✅ TRANSACTION RECORDING script
✅ Comment format & metadata
✅ Query filters
✅ Alur kerja lengkap (expire monitor, on-login, voucher generation, print)
✅ Penyimpanan metadata di database
✅ Ringkasan teknologi stack