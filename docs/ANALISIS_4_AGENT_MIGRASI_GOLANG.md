# Analisis Repository Mikhmon v4 dengan 4 Agent + Blueprint Migrasi Golang

## Ringkasan Eksekutif

Analisis ini memecah pekerjaan menjadi 4 agent agar migrasi dari PHP ke Go terstruktur, terukur, dan minim regresi.

- **Agent-1 (Arsitektur):** menilai struktur kode saat ini dan mendesain clean architecture target.
- **Agent-2 (Domain & Integrasi MikroTik):** memetakan fitur inti hotspot/profile/report dan strategi komunikasi `go-routeros v3` berbasis `ListenArgs` (tanpa pooling untuk data real-time).
- **Agent-3 (Data & Security):** memigrasikan konfigurasi router + akun admin ke PostgreSQL, caching ke Redis, dan hardening keamanan.
- **Agent-4 (Delivery & Quality):** memastikan testability, observability, CI, dan rollout bertahap.

---

## Agent-1 — Arsitektur & Boundary

### Temuan
- Repository sudah memiliki implementasi Go berbasis handler monolitik (Gin + handler langsung).
- Folder legacy PHP (`get/`, `post/`, `view/`) masih ada sehingga boundary domain belum sepenuhnya tegas.
- Dependency antar layer belum sepenuhnya mengikuti dependency rule clean architecture.

### Rekomendasi
Gunakan layer berikut:

1. `domain/` (entity + interface repository)
2. `usecase/` (business rules)
3. `infrastructure/` (postgres, redis, routeros)
4. `delivery/http` (Gin handlers)
5. `bootstrap/` (wiring dependency)

---

## Agent-2 — Domain Hotspot & MikroTik

### Temuan
- Core value Mikhmon: operasi user hotspot, profile, voucher, report live.
- Data operasional hotspot tetap di RouterOS (tidak dipindah ke PostgreSQL).
- Realtime lebih tepat menggunakan `ListenArgs` dibanding pooling saat stream event/traffic.

### Rekomendasi
- Buat `RouterOSListener` abstraction untuk stream event dan pembacaan metrik.
- Simpan hanya data administratif di database:
  - akun admin
  - konfigurasi router
- Endpoint operasional tetap mengeksekusi command RouterOS secara on-demand.

---

## Agent-3 — Data Platform, Cache, Security

### Temuan
- Kebutuhan masa depan: konfigurasi terpusat multi-router, login aman, audit yang rapi.
- Caching dibutuhkan untuk data non-kritis (mis. metadata router/session state), bukan untuk kebenaran utama data hotspot.

### Rekomendasi
- **PostgreSQL + GORM** untuk `users` dan `routers`.
- **Redis** untuk cache session/context singkat.
- Password hash `bcrypt`/`argon2`, secret via env.
- CSRF + secure cookie + rate limiting login.

---

## Agent-4 — Delivery, Observability, dan Risiko

### Temuan
- Sudah ada basis Go, sehingga strategi terbaik adalah migrasi bertahap, bukan rewrite sekali jalan.

### Rekomendasi fase delivery
1. **Foundation:** bootstrap infra (Zap, GORM, Redis, migration).
2. **Auth & Router Config:** login + CRUD router full via PostgreSQL.
3. **Hotspot Core:** user/profile/report via usecase + gateway RouterOS.
4. **Compat & Cutover:** endpoint compatibility, nonaktifkan PHP secara bertahap.

### KPI migrasi
- p95 response time endpoint utama
- login success rate
- error rate RouterOS command
- parity fitur vs PHP legacy

---

## Target Arsitektur Teknis (Modern Go)

```text
cmd/modern
  -> bootstrap
     -> config
     -> logger (zap)
     -> postgres (gorm)
     -> redis
     -> repositories
     -> usecases
     -> handlers
```

Prinsip:
- inward dependency (outer layer boleh tahu inner layer, bukan sebaliknya)
- interface-first di domain/usecase
- transport-agnostic usecase

---

## Strategi Migrasi dari Kondisi Repo Saat Ini

1. Pertahankan service Go existing agar operasional tetap jalan.
2. Tambahkan fondasi clean architecture dalam namespace baru (`internal/clean/...`) agar tidak mengganggu flow lama.
3. Migrasikan modul per modul:
   - auth/admin config dulu
   - hotspot/profile
   - report/traffic
4. Setelah parity tercapai, lakukan cutover route dan arsipkan kode PHP.

Dokumen ini dipasangkan dengan implementasi awal clean architecture pada commit yang sama sebagai baseline teknis.
