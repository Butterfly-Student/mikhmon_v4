# Test Queue Monitor MikroTik

Test sederhana untuk monitoring queue MikroTik menggunakan `GetAllParentQueues` dan `StartQueueStatsListen`.

## Cara Penggunaan

1. **Konfigurasi Router**

   Edit file `cmd/test_queue_monitor/main.go` dan ubah konstanti di bagian atas:

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
   go run cmd/test_queue_monitor/main.go
   ```

## Output yang Diharapkan

```
==========================================
MikroTik Queue Monitor Test
==========================================
Router: 192.168.88.1:8728

--- Get All Parent Queues ---
Total Parent Queues: 3
Parent Queues: [queue_parent_1 queue_parent_2 queue_parent_3]
First Queue: queue_parent_1

--- Start Queue Stats Monitor (5 seconds) ---
Monitoring queue: queue_parent_1
Monitoring started (5 seconds)...

[0.00s] Queue: queue_parent_1 | In: 1024000 bytes (1 MB) | Out: 512000 bytes (500 KB) | Rate In: 100000 bps | Rate Out: 50000 bps
[1.00s] Queue: queue_parent_1 | In: 2048000 bytes (2 MB) | Out: 1024000 bytes (1 MB) | Rate In: 100000 bps | Rate Out: 50000 bps
[2.00s] Queue: queue_parent_1 | In: 3072000 bytes (3 MB) | Out: 1536000 bytes (1.5 MB) | Rate In: 100000 bps | Rate Out: 50000 bps
[3.00s] Queue: queue_parent_1 | In: 4096000 bytes (4 MB) | Out: 2048000 bytes (2 MB) | Rate In: 100000 bps | Rate Out: 50000 bps
[4.00s] Queue: queue_parent_1 | In: 5120000 bytes (5 MB) | Out: 2560000 bytes (2.5 MB) | Rate In: 100000 bps | Rate Out: 50000 bps

--- Monitoring Stopped ---
```

## Fungsi yang Diuji

| Fungsi | Deskripsi |
|---------|-----------|
| `GetAllParentQueues` | Mengambil daftar parent queues |
| `StartQueueStatsListen` | Streaming statistik queue (Bytes In/Out, Rate In/Out) |
| `DefaultQueueStatsConfig` | Membuat konfigurasi default untuk monitoring |

## Fitur

- ✅ Mengambil semua parent queues dari MikroTik
- ✅ Monitoring stats queue pertama yang ditemukan
- ✅ Menampilkan statistik real-time setiap detik selama 5 detik
- ✅ Format bytes yang mudah dibaca (B, KB, MB)
- ✅ Auto cleanup koneksi dengan `defer cancel()`
