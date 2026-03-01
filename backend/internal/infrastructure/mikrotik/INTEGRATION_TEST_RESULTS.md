# Integration Test Results

**Router:** `192.168.233.1:8728` (RouterOS 6.49.11 stable — RB750G "G-Net")
**Tanggal:** 2026-03-02 pukul 00:40 WIB
**Total waktu:** ~38s
**Hasil:** ✅ **SEMUA PASS** (75 pass, 5 skip, 0 fail)

---

## Ringkasan per File

### `integration_test.go` — Koneksi, Async, System, Listener
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_Connect | ✅ PASS | Koneksi berhasil ke router |
| TestIntegration_Async_IsAsync | ✅ PASS | Mode async aktif setelah Connect() |
| TestIntegration_RunContext_WithTimeout | ✅ PASS | RunContext + context timeout berjalan normal |
| TestIntegration_RunMany_Concurrent | ✅ PASS | 5 perintah concurrent selesai dalam ~66ms |
| TestIntegration_System_GetResource | ✅ PASS | RouterOS 6.49.11, RB750G, CPU 91%, RAM 5/32 MB, uptime 2d11h |
| TestIntegration_System_GetIdentity | ✅ PASS | Identity: "G-Net" |
| TestIntegration_System_GetHealth | ✅ PASS | Health returned (voltage/temp kosong di model ini) |
| TestIntegration_System_GetClock | ✅ PASS | Clock: mar/02/2026 00:40:25 (Asia/Jakarta) |
| TestIntegration_System_GetRouterBoardInfo | ✅ PASS | Model RB750G, serial 228E01AED081, firmware 2.23 |
| TestIntegration_Interface_GetAll | ✅ PASS | 6 interface: ether1–5, l2tp-out1 |
| TestIntegration_Pool_GetAddressPools | ✅ PASS | 4 pool: dhcp_pool0–2, e2e-test-pool |
| TestIntegration_Queue_GetAllQueues | ✅ PASS | 3 queue: TRAFIK, DIST, LOCAL |
| TestIntegration_Queue_GetAllParentQueues | ✅ PASS | 3 parent queues |
| TestIntegration_Hotspot_GetServers | ✅ PASS | 1 server: "all" |
| TestIntegration_Hotspot_GetActive | ✅ PASS | 0 sesi aktif saat test |
| TestIntegration_Hotspot_GetActiveCount | ✅ PASS | Count: 0 |
| TestIntegration_Listener_TrafficMonitor_Once | ✅ PASS | 1 sample via =once=: rx=8.1Mbps, tx=429kbps |
| TestIntegration_Listener_TrafficMonitor_Continuous | ✅ PASS | 3 sample: rx=8.1/9.8/2.0Mbps, tx=429/339/81kbps |
| TestIntegration_Listener_QueueStats | ✅ PASS | 3 sample TRAFIK: rateIn~267kbps, rateOut~5.7Mbps |
| TestIntegration_Listener_ResourceMonitor | ✅ PASS | 3 sample: cpu=3–8%, freeMem=6.0MiB |
| TestIntegration_Listener_Ping | ✅ PASS | 3x ping ke 8.8.8.8: ~26–27ms, 100% received |
| TestIntegration_RunAndListenConcurrent | ✅ PASS | RunMany + ListenArgs bersamaan — async confirmed |

---

### `hotspot_active_integration_test.go` — Active Sessions
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_HotspotActive_RemoveSession | ⏭ SKIP | Di-skip: tidak ada sesi aktif saat test |

---

### `hotspot_hosts_integration_test.go` — Hotspot Hosts
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_HotspotHosts_GetAll | ✅ PASS | 0 host aktif saat test (endpoint berfungsi normal) |
| TestIntegration_HotspotHosts_Remove | ⏭ SKIP | Di-skip: tidak ada host aktif saat test |

---

### `hotspot_listen_integration_test.go` — Hotspot Streaming (Baru)
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_HotspotActive_Listen | ⏭ SKIP | Di-skip: tidak ada sesi aktif — RouterOS tidak mengirim apapun untuk tabel kosong |
| TestIntegration_HotspotInactive_Listen | ✅ PASS | Diff diterima: 13 inactive users dalam 0.29s (13 users, 0 active) |

---

### `hotspot_profiles_integration_test.go` — User Profiles
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_HotspotProfiles_GetAll | ✅ PASS | 9 profile: default, Testing, Testing1, 7d, 3jam-5mbps, Testing 55, 1H, + 2 leftover inttest |
| TestIntegration_HotspotProfiles_GetByID | ✅ PASS | Lookup by ID: "default" ditemukan |
| TestIntegration_HotspotProfiles_GetByName | ✅ PASS | Lookup by name: "default" ditemukan |
| TestIntegration_HotspotProfiles_AddUpdateRemove | ✅ PASS | Add → resolve ID → update rate-limit "1M/1M" → remove |

---

### `hotspot_users_integration_test.go` — Hotspot Users
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_HotspotUsers_GetAll | ✅ PASS | 13 user: default-trial + 10 "1H" + 2 leftover inttest |
| TestIntegration_HotspotUsers_GetByProfile | ✅ PASS | Filter "1H": 10 user |
| TestIntegration_HotspotUsers_GetCount | ✅ PASS | Count returned |
| TestIntegration_HotspotUsers_GetByID | ✅ PASS | Lookup by ID: "default-trial" |
| TestIntegration_HotspotUsers_GetByName | ✅ PASS | Lookup by name: "default-trial" |
| TestIntegration_HotspotUsers_GetByComment | ✅ PASS | Filter by comment (0 result — normal) |
| TestIntegration_HotspotUsers_AddUpdateRemove | ✅ PASS | Add → resolve ID → update comment → remove |
| TestIntegration_HotspotUsers_RemoveByComment | ✅ PASS | Buat 2 user → hapus semua by comment → verified kosong |
| TestIntegration_HotspotUsers_ResetCounters | ✅ PASS | Reset counters "default-trial" berhasil |

---

### `logging_integration_test.go` — Logging
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_Logging_EnableHotspotLogging | ✅ PASS | Logging sudah terkonfigurasi / berhasil dikonfigurasi |
| TestIntegration_Logging_GetHotspotLogs_All | ✅ PASS | 0 hotspot log (belum ada aktivitas saat test) |
| TestIntegration_Logging_EnablePPPLogging | ✅ PASS | PPP logging dikonfigurasi atau sudah ada |
| TestIntegration_Logging_GetHotspotLogs_WithLimit | ✅ PASS | Limit=5 dihormati, 0 entry |

---

### `logging_listen_integration_test.go` — Log Streaming (Baru)
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_GetPPPLogs | ✅ PASS | 0 PPP log (tidak ada aktivitas PPP saat test) |
| TestIntegration_ListenHotspotLogs | ✅ PASS | Stream 5s berjalan tanpa error, 0 hotspot entries |
| TestIntegration_ListenPPPLogs | ✅ PASS | Stream 5s berjalan tanpa error, 0 PPP entries |
| TestIntegration_ListenLogs_Generic | ✅ PASS | **1012 log entries** terkumpul dalam 5s (DHCP, system, account events) |

---

### `nat_integration_test.go` — NAT Rules
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_NAT_GetRules | ✅ PASS | 8 rule NAT: srcnat masquerade + dstnat dst-nat |
| TestIntegration_NAT_GetRules_HasValidFields | ✅ PASS | Semua rule memiliki ID, Chain, Action valid |

---

### `ppp_integration_test.go` — PPP (Baru)
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_PPPSecrets_List | ✅ PASS | 0 PPP secrets (router tidak digunakan untuk PPP) |
| TestIntegration_PPPSecrets_CRUD | ✅ PASS | Add → GetByName (id=*5) → Update comment → GetByID → Remove → verify nil |
| TestIntegration_PPPSecrets_DisableEnable | ✅ PASS | Add → Disable (disabled=true) → Enable (disabled=false) → Remove |
| TestIntegration_PPPProfiles_List | ✅ PASS | 8 profile: default, testing, teteh, 10Mbps, 100Mbps, 100-RB-100, profile_test_10mbps, default-encryption |
| TestIntegration_PPPProfiles_CRUD | ✅ PASS | Add → GetByName (id=*A) → Update rateLimit 1M/1M→2M/2M → Remove → verify nil |
| TestIntegration_PPPActive_List | ✅ PASS | 0 active PPP sessions |
| TestIntegration_PPPActive_Listen | ⏭ SKIP | Di-skip: tidak ada PPP active — RouterOS tidak mengirim apapun untuk tabel kosong |
| TestIntegration_PPPInactive_Listen | ⏭ SKIP | Di-skip: tidak ada PPP secrets — tabel kosong tidak menghasilkan stream |

---

### `reports_integration_test.go` — Sales Reports
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_Reports_GetAll | ✅ PASS | 5 report (3 dari sesi sebelumnya + 2 dari run ini) |
| TestIntegration_Reports_GetByOwner | ✅ PASS | Filter by owner mengembalikan 1 report yang sesuai |
| TestIntegration_Reports_GetByDay | ✅ PASS | Filter "mar/01/2026": 3 report |
| TestIntegration_Reports_AddAndVerify | ✅ PASS | Add → retrieve by owner → harga dan username cocok |

---

### `expire_monitor_integration_test.go` — Expire Monitor
| Test | Status | Keterangan |
|------|--------|-----------|
| TestIntegration_ExpireMonitor_Ensure | ✅ PASS | Status "existing" — scheduler sudah ada dan aktif |
| TestIntegration_ExpireMonitor_IdempotentCall | ✅ PASS | Panggil 2x → keduanya "existing", tidak ada duplikat |
| TestIntegration_ExpireMonitor_ScriptContent | ✅ PASS | Script 1513 char; ExpireMode=remc, Validity=30d, Price=5000, SellingPrice=5500 |

---

### Unit Tests (tidak perlu router)
| Test | Status | Keterangan |
|------|--------|-----------|
| TestOnLoginGenerator_Parse (4 subtest) | ✅ PASS | Parse ntf/remc/rem/mode0 |
| TestOnLoginGenerator_Generate_ContainsPicIndex | ✅ PASS | |
| TestOnLoginGenerator_Generate_ExpirationBlock | ✅ PASS | |
| TestOnLoginGenerator_Mode0_WithLock | ✅ PASS | |
| TestOnLoginGenerator_Mode0_NoLock | ✅ PASS | |
| TestParseRate (11 subtest) | ✅ PASS | bps/kbps/Mbps/Gbps parsing |
| TestSplitRateValue (5 subtest) | ✅ PASS | |
| TestSplitSlashValue (5 subtest) | ✅ PASS | |
| TestParseQueueStatsSentence | ✅ PASS | |
| TestVoucherGenerator_* (5 test) | ✅ PASS | |
| TestParseDataLimit (10 subtest) | ✅ PASS | |

---

## Temuan Penting dari Real Router

| Temuan | Detail |
|--------|--------|
| RouterOS version | 6.49.11 (stable) — tidak support `ret` dari `/add`, ID harus di-lookup by name setelah add |
| Interface aktif | ether1, ether3, ether5, l2tp-out1 |
| Hotspot server | "all" (1 server) |
| Address pool | dhcp_pool0, dhcp_pool1, dhcp_pool2, e2e-test-pool |
| Simple queue | TRAFIK, DIST, LOCAL (throughput: rateIn~267kbps, rateOut~5.7Mbps) |
| Traffic aktif | ether1 rx ~8–9 Mbps / tx ~80–430 kbps saat test |
| Ping ke 8.8.8.8 | ~26–27ms, 100% received |
| Expire Monitor | Scheduler "Mikhmon-Expire-Monitor" sudah ada di router |
| PPP profiles | 8 profile tersedia (testing, teteh, 10Mbps, 100Mbps, dll), 0 active session |
| Log stream volume | `/log/print =follow=` menghasilkan **1012 entries dalam 5s** (sangat aktif) |
| RouterOS follow empty table | Tidak mengirim apapun untuk tabel kosong — test yang perlu data aktif di-skip |

---

## Catatan

- Test dengan **SKIP** (5 test) bukan kegagalan — di-skip karena kondisi router saat test:
  - `HotspotActive_RemoveSession`, `HotspotActive_Listen` — tidak ada sesi aktif hotspot
  - `HotspotHosts_Remove` — tidak ada hotspot host
  - `PPPActive_Listen`, `PPPInactive_Listen` — router tidak dikonfigurasi sebagai PPP server (0 secrets, 0 active)
- RouterOS 6.x tidak mengembalikan ID dari `/add` melalui `reply.Re`. ID harus diambil via lookup by name — sudah ditangani di semua test CRUD.
- **ListenLogs deadlock fix**: Log stream sangat aktif (1000+ entries/5s). Implementasi `ListenLogs` menggunakan non-blocking send ke resultChan (drop if full) agar `listenReply.Chan()` selalu terbaca dan asyncLoop RouterOS tidak stall saat `Cancel()` dipanggil.
- Test report meninggalkan data di `/system/script` router — efek samping yang disengaja.
