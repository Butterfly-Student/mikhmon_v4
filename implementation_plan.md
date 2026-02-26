# Backend Gap Analysis: Mikhmon v4 Go vs MIKHMON_ANALYSIS.md & API Docs

Analisis perbandingan implementasi backend Go dengan referensi MIKHMON v4 (PHP) berdasarkan [MIKHMON_ANALYSIS.md](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/MIKHMON_ANALYSIS.md) dan [MIKHMON_v4_API_ENDPOINTS_DOCUMENTATION.md](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/MIKHMON_v4_API_ENDPOINTS_DOCUMENTATION.md).

---

## 📊 Ringkasan Gap

| Area | Status | Keterangan |
|------|--------|------------|
| RouterOS System Commands | ✅ Lengkap | Semua 5 command system sudah ada |
| Hotspot User CRUD | ✅ Lengkap | add/set/remove/print/reset-counters |
| Hotspot Profile CRUD | ✅ Lengkap | add/set/remove/print + on-login generator |
| Hotspot Active | ✅ Lengkap | print/remove/count-only |
| Hotspot Host | ✅ Lengkap | print/remove |
| Interface & Traffic | ✅ Lengkap | interface/print + monitor-traffic |
| Address Pool & Queue | ⚠️ Parsial | Hanya return nama, bukan data lengkap |
| On-Login Script | ⚠️ Ada Gap | Parse regex salah, mode "0" (no expire) kurang |
| Expire Monitor | ✅ Lengkap | Script + scheduler management sesuai |
| Voucher Generate | ⚠️ Ada Gap | Comment format berbeda dari original |
| Sales Report | ✅ Lengkap | Parse `-\|-` format sudah benar |
| System Logging | ⚠️ Ada Gap | Filter log berbeda dari original |
| User Count | ⚠️ Ada Gap | Tidak dikurangi 1 (admin user) |
| DataLimit Parsing | ⚠️ Ada Gap | UpdateUser tidak parse data_limit |
| Queue Command | ⚠️ Salah | Menggunakan `/queue/simple` bukan `/queue/tree` |
| NAT Rules | ✅ Tambahan | Fitur baru di Go, tidak ada di PHP original |
| Hotspot Server Detail | ⚠️ Parsial | Hanya return nama, bukan object lengkap |
| Test Coverage | ❌ Tidak Ada | 0 test files ditemukan |

---

## 🔍 Detail Per-Area Gap

---

### 1. On-Login Script Generator — [onlogin_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go)

#### 1a. ❌ Parse Regex Salah

File: [onlogin_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go#L230)

```go
// Current (SALAH)
putPattern := regexp.MustCompile(`:put \(,\"([^\"]*)\",([^,]*),([^,]*),([^,]*),,([^,]*),([^,]*),\"\)`)
```

Format output sebenarnya dari [buildHeader()](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go#71-97) (line 87):
```go
`:put (",%s,%.0f,%s,%.0f,,%s,%s,"); :local mode "%s"; {`
```

Outputnya: `:put (",ntf,10000,1d,15000,,Enable,Disable,");`

Regex saat ini menggunakan escaped quote di dalam capture group yang **tidak match** dengan format yang di-generate. Harus difix agar bisa parse on-login script yang sudah ada di MikroTik.

#### 1b. ⚠️ Mode "0" (Tanpa Expire) Belum Complete

Dari MIKHMON_ANALYSIS, ada mode `"0"` yang berarti **tanpa expire**. Di [buildExpirationLogic()](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go#98-139) sudah return `""` untuk mode ini, tapi [buildFooter()](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go#164-180) juga return `""`, sehingga jika `LockUser=Enable` akan menambah lock script tanpa wrapper `{..}`.

#### 1c. ⚠️ Missing Closing Brace pada Expiration Script

Di [buildExpirationLogic()](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go#98-139), script berakhir dengan:
```
/sys sch remove [find where name="$user"];
```
Tapi tidak ada closing `}` untuk block `do={...}`. Block `do={` di line 112 tidak pernah ditutup di dalam expiration logic, bergantung pada Footer. Ini perlu diverifikasi apakah assembly-nya benar.

---

### 2. Voucher Generator — [voucher_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/voucher_generator.go)

#### 2a. ⚠️ Comment Format Tidak Sesuai Original

Dari MIKHMON_ANALYSIS, format comment:
```
vc-<gencode>-<date>-<comment>
```

Di [voucher_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/voucher_generator.go) line 26:
```go
comment := fmt.Sprintf("%s-%s-%s", req.Mode, time.Now().Format("01.02.06"), req.Comment)
```
Output: `vc-01.02.06-Daily Voucher` — **tanpa gencode**.

Tapi di [voucher_usecase.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/voucher_usecase.go) line 56, comment sudah benar:
```go
comment := fmt.Sprintf("%s-%s-%s-%s", req.Mode, gencode, time.Now().Format("01.02.06"), strings.TrimSpace(req.Comment))
```
Output: `vc-538-01.02.06-Daily Voucher` — **sesuai original**.

**Gap**: `VoucherGenerator.GenerateBatch()` menghasilkan comment yang salah, tetapi `VoucherUseCase.GenerateVouchers()` meng-override-nya. Ini duplikasi logika yang membingungkan.

---

### 3. Queue Command Salah — [helpers.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/helpers.go)

#### ❌ Menggunakan `/queue/simple` Bukan `/queue/tree`

File: [helpers.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/helpers.go#L38)

```go
// Current (SALAH)
reply, err := client.RunContext(ctx, "/queue/simple/print")
```

Dari MIKHMON_ANALYSIS:
```php
// Original (BENAR)
$parent_queues = $API->comm("/queue/tree/print");
```

Harus diganti ke `/queue/tree/print`.

---

### 4. Address Pool — Hanya Return Nama

File: [helpers.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/helpers.go#L10)

Backend hanya return `[]string` (daftar nama pool). Original PHP mengembalikan object lengkap termasuk `ranges`. Untuk keperluan UI dropdown, ini cukup. Tapi jika frontend butuh detail, perlu enrichment.

---

### 5. Hotspot Server — Hanya Return Nama

File: [helpers.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/helpers.go#L54)

API docs menunjukkan hotspot server mengembalikan object lengkap ([name](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/voucher_generator.go#85-96), `address-pool`, `profile`, `interface`). Backend hanya return `[]string`. Perlu ditingkatkan jika frontend butuh data server lengkap.

---

### 6. User Count Tidak Dikurangi 1 (Admin)

File: [hotspot_users.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/hotspot_users.go#L309)

Dari MIKHMON_ANALYSIS & API docs:
> `hotspot_users`: Total user hotspot terdaftar (**dikurangi 1 admin**)

Backend saat ini mengembalikan count mentah tanpa pengurangan. Ini bisa menyebabkan jumlah user di dashboard **lebih 1** dari seharusnya.

---

### 7. UpdateUser Tidak Parse DataLimit

File: [hotspot_service.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/hotspot_service.go#L151)

```go
LimitBytesTotal: 0, // TODO: parse from req.DataLimit
```

Ada TODO yang belum selesai. [UpdateUser](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/hotspot_service.go#137-158) tidak pernah mengirim `limit-bytes-total` ke MikroTik, sehingga data limit **tidak bisa diupdate** via API.

---

### 8. System Logging Query Berbeda

File: [logging.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/logging.go#L22)

```go
// Current
client.RunContext(ctx, "/log/print", "?topics=hotspot,info,debug")
```

Original PHP:
```php
$API->comm("/system/logging/print", array("?prefix" => "->"));
```

Ini berbeda secara fundamental. Original menggunakan `/system/logging/print` dengan filter prefix, sedangkan Go menggunakan `/log/print` dengan filter topics. Fungsinya berbeda: satu mengambil konfigurasi logging, satu mengambil log entries. Logika Go mungkin lebih tepat untuk shows log, tapi perlu diverifikasi filter topics-nya benar.

---

### 9. On-Login Script — `$exp 7 16` vs `$exp 7 15`

File: [onlogin_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go#L120)

```routeros
// Go version:
:local t [:pic $exp 7 16];
```

Di MIKHMON_ANALYSIS:
```routeros
// PHP version:
:local t [:pic $exp 7 15];
```

Perbedaan index `15` vs `16` bisa menyebabkan parsing waktu expire **off-by-one character**. Perlu diverifikasi mana yang benar.

---

## 📋 Improvement Plan (Prioritas)

### Priority 1 — Bug Fixes (Critical) 🔴

| # | Item | File | Impact |
|---|------|------|--------|
| 1 | Fix on-login Parse regex agar match dengan format yang di-generate | [onlogin_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go) | Profile yang sudah ada tidak bisa di-parse |
| 2 | Fix queue command dari `/queue/simple` ke `/queue/tree` | [helpers.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/helpers.go) | Parent queue salah |
| 3 | Implement DataLimit parsing di UpdateUser | [hotspot_service.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/hotspot_service.go) | Data limit tidak bisa diupdate |
| 4 | Verifikasi `[:pic $exp 7 16]` vs `[:pic $exp 7 15]` | [onlogin_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go) | Expire time bisa salah |

### Priority 2 — Feature Gaps (Important) 🟡

| # | Item | File | Impact |
|---|------|------|--------|
| 5 | Kurangi 1 dari user count (admin) | [hotspot_users.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/hotspot_users.go) / [dashboard_usecase.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/dashboard_usecase.go) | User count lebih 1 |
| 6 | Perkaya return Hotspot Servers (object, bukan string) | [helpers.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/helpers.go) + new DTO | Frontend terbatas info |
| 7 | Perkaya return Address Pools (object dengan ranges) | [helpers.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/helpers.go) + new DTO | Frontend terbatas info |
| 8 | Bersihkan duplikasi comment formatting di VoucherGenerator | [voucher_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/voucher_generator.go) + [voucher_usecase.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/voucher_usecase.go) | Maintenance risk |
| 9 | Fix mode "0" (no expire) agar lock scripts wrapped benar | [onlogin_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go) | Edge case error |

### Priority 3 — Enhancements (Nice to Have) 🟢

| # | Item | File | Impact |
|---|------|------|--------|
| 10 | Tambah unit tests untuk onlogin_generator | Baru: `onlogin_generator_test.go` | Quality assurance |
| 11 | Tambah unit tests untuk voucher_generator | Baru: `voucher_generator_test.go` | Quality assurance |
| 12 | Tambah unit tests untuk reports parser | Baru: `reports_test.go` | Quality assurance |
| 13 | Review logging query (`/log/print` vs original) | [logging.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/logging.go) | Log accuracy |

---

## Verification Plan

### Automated Tests

Saat ini **tidak ada test files** (`0 files *_test.go`). Setelah implementasi fix, perlu:

```bash
# Run all tests
cd /home/butterfly_student/code/Mikrotik/mikhmon_v4/backend && go test ./...

# Run specific package tests 
cd /home/butterfly_student/code/Mikrotik/mikhmon_v4/backend && go test ./internal/infrastructure/mikrotik/...
```

### Manual Verification

> [!IMPORTANT]
> Karena project ini terhubung langsung ke RouterOS device, verifikasi penuh **membutuhkan router MikroTik yang aktif**. Mohon informasikan apakah ada router test yang bisa digunakan, atau apakah verifikasi hanya bisa dilakukan via code review dan unit tests.

Untuk verifikasi tanpa router:
1. **Build check**: `cd /home/butterfly_student/code/Mikrotik/mikhmon_v4/backend && go build ./...`
2. **Code review** per-fix terhadap referensi MIKHMON_ANALYSIS.md
3. **Unit tests** dengan mock MikroTik client untuk fungsi pure-logic (parser, generator, helpers)
