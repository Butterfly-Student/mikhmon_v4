# Analisis Backend vs Mikhmon (berdasarkan kode saat ini)

Tanggal analisis: 25 Februari 2026

## Catatan penting tentang referensi
File referensi yang Anda minta untuk dibandingkan saat ini kosong (0 byte):
- `MIKHMON_ANALYSIS.md`
- `MIKHMON_v4_API_ENDPOINTS_DOCUMENTATION.md`

Akibatnya, analisis di bawah **hanya** berdasarkan implementasi backend di repo ini, bukan perbandingan langsung dengan dokumen Mikhmon. Jika Anda mengisi kedua file tersebut, saya bisa ulangi analisis dengan pemetaan yang presisi.

## Ringkasan cepat
- Penyimpanan data di MikroTik sudah ada untuk: `hotspot users`, `hotspot user profiles`, dan `sales report` di `/system/script`.
- On-login script generator sudah dibuat, tetapi terdapat indikasi bug sintaks dan parsing sehingga metadata Mikhmon kemungkinan tidak terdeteksi/berjalan semestinya.
- Mekanisme **expiration enforcement** (hapus/limit user ketika masa habis) **belum terimplementasi** di backend.
- Modul **reporting** sisi usecase/handler **belum jalan** (placeholder), walaupun parser data report sudah ada.

## 1. Penyimpanan data di MikroTik
**Terimplementasi:**
- Hotspot users disimpan di `/ip/hotspot/user` melalui `AddHotspotUser`, `UpdateHotspotUser`, `RemoveHotspotUser`. (lihat `backend/internal/infrastructure/mikrotik/hotspot_users.go`)
- Hotspot user profiles disimpan di `/ip/hotspot/user/profile` melalui `AddUserProfile`, `UpdateUserProfile`, `RemoveUserProfile`. (lihat `backend/internal/infrastructure/mikrotik/hotspot_profiles.go`)
- Sales report disimpan sebagai entri `/system/script` dengan `comment=mikhmon` dan `owner` berupa bulan/tahun. (lihat `backend/internal/infrastructure/mikrotik/reports.go`)

**Catatan:**
- Penyimpanan report mengandalkan format nama script `date-| -time-| -user-| -price-| -address-| -mac-| -validity-| -profile-| -comment`.

## 2. On-login script
**Terimplementasi:**
- Generator on-login script ada di `backend/internal/infrastructure/mikrotik/onlogin_generator.go`.
- Script melakukan:
  - Output metadata `:put` (expire mode, price, validity, selling price, lock options).
  - Kalkulasi expiration dan update `comment` user berdasarkan `scheduler` dan `next-run`.
  - Pencatatan transaksi ke `/system/script` untuk mode `remc` dan `ntfc`.
  - Lock MAC / lock server bila diaktifkan.

**Potensi masalah serius (kemungkinan bug):**
- Di `buildExpirationLogic` terdapat banyak pemakaian `:pic` yang seharusnya `:pick` (RouterOS). Ini akan membuat script gagal dijalankan di RouterOS.
- Format `:put` yang dihasilkan:
  - Output: `:put (",remc,5000,30d,5500,,Enable,Disable,");`
  - Regex parser di `Parse()` mengharapkan pola `:put (,"...",...` (koma sebelum kutip). Pola ini **tidak cocok** dengan output aktual. Akibatnya `ExpireMode`, `Validity`, `Price`, dll kemungkinan tidak ter-parse.

## 3. Expiration users
**Yang ada sekarang:**
- On-login script menghitung tanggal expire, lalu menulisnya ke `comment` user dalam format `"<date/time> <mode>"`.
- Mode `rem` / `ntf` dibedakan sebagai `X` / `N` pada `comment`.

**Yang belum ada:**
- Tidak ditemukan mekanisme **monitoring expiration** yang secara otomatis:
  - menghapus user saat expire (mode `rem`), atau
  - membatasi user (mode `ntf`),
  - ataupun scheduler global yang melakukan cleanup.
- Fungsi `GenerateExpiredAction()` ada, tetapi **tidak digunakan** di mana pun.

**Kesimpulan:**
- Expiration hanya disimpan sebagai metadata di `comment`, tetapi **tidak dieksekusi** untuk disable/hapus user.

## 4. Reporting
**Yang ada:**
- Parser laporan di `/system/script` sudah ada (`parseSalesReport`).
- Bisa membaca report berdasarkan `owner` (bulan) atau `source` (hari).

**Yang belum ada:**
- Usecase `ReportUseCase.GetSalesReport()` dan `GetSalesReportByDay()` **masih placeholder** (selalu return kosong).
- Handler laporan akan selalu mengembalikan data kosong.
- Export CSV belum diimplementasikan.

**Kesimpulan:**
- Reporting **belum berfungsi** walaupun fondasi parsing MikroTik sudah ada.

## 5. Lainnya yang relevan
- Update user data limit belum diproses: di `hotspot_service.UpdateUser`, `LimitBytesTotal` selalu `0` (TODO parse dari `req.DataLimit`).
- Voucher `comment` di `VoucherUseCase.GenerateVouchers()` selalu diawali `vc-...` meskipun mode bisa `up`. Ini bisa membuat on-login menganggap semua voucher sebagai mode `vc`.

## Rekomendasi prioritas perbaikan
1. Perbaiki sintaks on-login script (`:pic` -> `:pick`) dan sesuaikan regex parse header agar sesuai output.
2. Implementasikan expiration enforcement (scheduler global atau job backend) yang mengeksekusi `GenerateExpiredAction()`.
3. Implementasikan `ReportUseCase.GetSalesReport()` dan `GetSalesReportByDay()` agar laporan berfungsi.
4. Implementasikan parse `DataLimit` di update user.
5. Sesuaikan prefix comment voucher sesuai mode (`vc` vs `up`).

## File yang dianalisis
- `backend/internal/infrastructure/mikrotik/onlogin_generator.go`
- `backend/internal/infrastructure/mikrotik/hotspot_profiles.go`
- `backend/internal/infrastructure/mikrotik/hotspot_users.go`
- `backend/internal/infrastructure/mikrotik/reports.go`
- `backend/internal/usecase/report_usecase.go`
- `backend/internal/usecase/voucher_usecase.go`
- `backend/internal/infrastructure/http/handler/report_handler.go`

