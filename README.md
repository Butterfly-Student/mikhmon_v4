# mikhmon_v4 (Go Edition)

Repository ini sudah difinalisasi ke runtime Golang.

## Stack runtime aktif
- Go
- Gin
- go-routeros v3
- HTML template embedded

## Service dependencies (PostgreSQL + Redis)
File yang disediakan:
- `.env`
- `docker-compose.yaml`

Menjalankan dependency:

```bash
docker compose up -d
```

## Validasi migrasi total
Jalankan:

```bash
./scripts/verify_no_php.sh
```

Jika sukses, output akan menampilkan:

```text
OK: no PHP files found
```
