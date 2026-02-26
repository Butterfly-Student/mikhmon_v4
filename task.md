# Backend Gap Analysis — Implementation Tasks

## Priority 1 — Bug Fixes (Critical) 🔴
- [ ] Fix on-login Parse regex di [onlogin_generator.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/onlogin_generator.go)
- [ ] Fix queue command `/queue/simple` → `/queue/tree` di [helpers.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/helpers.go)
- [ ] Implement DataLimit parsing di [hotspot_service.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/hotspot_service.go) UpdateUser
- [ ] Verifikasi `[:pic $exp 7 16]` vs `[:pic $exp 7 15]` di on-login script

## Priority 2 — Feature Gaps (Important) 🟡
- [ ] Kurangi 1 dari user count (admin) di dashboard
- [ ] Perkaya Hotspot Servers return (object, bukan string)
- [ ] Perkaya Address Pools return (object dengan ranges)
- [ ] Bersihkan duplikasi comment di VoucherGenerator
- [ ] Fix mode "0" (no expire) lock scripts wrapping

## Priority 3 — Enhancements 🟢  
- [ ] Unit tests untuk onlogin_generator
- [ ] Unit tests untuk voucher_generator
- [ ] Unit tests untuk reports parser
- [ ] Review logging query

## Verification
- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
