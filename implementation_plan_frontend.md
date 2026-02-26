# Frontend Alignment dengan Backend Changes

Menyesuaikan frontend TypeScript/React dengan perubahan backend Go yang sudah dilakukan pada sesi sebelumnya.

## Perubahan Backend yang Perlu Disesuaikan

| Backend Change | File | Frontend Impact |
|---|---|---|
| [DashboardData](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go#137-150) struct punya `health`, `routerBoard`, `interfaces`, `hotspotLogs` | [dto/dashboard.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go) | [DashboardData](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go#137-150) type di frontend belum punya field ini |
| `stats` sebagai object `{totalUsers, activeUsers}` | [dto/dashboard.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go) | ✅ Sudah ditangani dengan `??` fallback |
| [SystemHealth](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go#26-33) punya `voltage`, `temperature` sebagai string | [dto/dashboard.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go) | [SystemResources](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend/src/types/index.ts#171-187) type pakai `voltage?: number\|string` — perlu pastikan konsisten |
| [LogEntry](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go#123-129) struct: [id](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/voucher_usecase.go#154-175), `time`, `topics`, `message` | [dto/dashboard.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go) | Type belum ada di frontend |
| [Interface](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go#64-82) struct: lengkap dengan `rxByte`, `txByte`, dll. | [dto/dashboard.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go) | Type belum ada di frontend |
| [GetHotspotLogs](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/logging.go#10-42) → `GET /dashboard/:id/logs` | [dashboard_handler.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/http/handler/dashboard_handler.go) | Belum ada di `dashboardApi` |
| [GetInterfaces](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/http/handler/dashboard_handler.go#155-175) → `GET /dashboard/:id/interfaces` | [dashboard_handler.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/http/handler/dashboard_handler.go) | Belum ada di `dashboardApi` |
| [GetNATRules](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/http/handler/dashboard_handler.go#203-223) → `GET /dashboard/:id/nat` | [dashboard_handler.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/http/handler/dashboard_handler.go) | Belum ada di `dashboardApi` |
| [UpdateUser](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/http/handler/hotspot_handler.go#120-148) sekarang parse `dataLimit` via [ParseDataLimit()](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/voucher_generator.go#160-188) | [hotspot_service.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/hotspot_service.go) | Form UsersPage belum kirim `timeLimit`/`dataLimit` |
| [AddUser](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/hotspot_usecase.go#59-74) di [hotspot_usecase.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/hotspot_usecase.go) masih `LimitBytesTotal: 0` | [hotspot_usecase.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/hotspot_usecase.go) | Backend bug — akan difix bersamaan |
| User count dikurangi 1 (admin) | [hotspot_users.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/hotspot_users.go) | ✅ Tidak perlu perubahan frontend |

---

## Proposed Changes

### Backend (Fix Kecil Sisa)

#### [MODIFY] [hotspot_usecase.go](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/hotspot_usecase.go)

[AddUser](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/hotspot_usecase.go#59-74) masih hardcode `LimitBytesTotal: 0`. Perlu parse `req.DataLimit` sama seperti [UpdateUser](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/http/handler/hotspot_handler.go#120-148).

---

### Frontend

#### [MODIFY] [types/index.ts](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend/src/types/index.ts)

- Tambah [LogEntry](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go#123-129) interface (sesuai `dto.LogEntry`: [id](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/usecase/voucher_usecase.go#154-175), `time`, `topics`, `message`)
- Tambah `NetworkInterface` interface (sesuai `dto.Interface`: [name](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/infrastructure/mikrotik/voucher_generator.go#84-95), `type`, `running`, `disabled`, `rxByte`, `txByte`, ...)
- Tambah [RouterBoardInfo](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go#52-61) interface (sesuai `dto.RouterBoardInfo`)
- Tambah [SystemHealth](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go#26-33) interface dengan `voltage: string`, `temperature: string`, `fanSpeed?: string`
- Update [DashboardData](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/dashboard.go#137-150) agar include `health?`, `routerBoard?`, `interfaces?`, `hotspotLogs?`, `connectionError?`
- Update [SystemResources](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend/src/types/index.ts#171-187): field `voltage` dan `temperature` ubah ke `string` untuk konsisten dengan backend

#### [MODIFY] [dashboard.ts](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend/src/api/dashboard.ts)

Tambah 3 fungsi API baru:
- `getLogs(routerId, limit?)` → `GET /dashboard/:id/logs?limit=N`
- `getInterfaces(routerId)` → `GET /dashboard/:id/interfaces`
- `getNATRules(routerId)` → `GET /dashboard/:id/nat`

#### [MODIFY] [hotspot.ts](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend/src/api/hotspot.ts)

- [createUser](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend/src/api/hotspot.ts#45-57): tambah `timeLimit` dan `dataLimit` di payload (sudah ada di [AddUserRequest](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/hotspot.go#96-107) DTO)
- [updateUser](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend/src/api/hotspot.ts#58-70): tambah `timeLimit` dan `dataLimit` di payload (sudah ada di [UpdateUserRequest](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/backend/internal/domain/dto/hotspot.go#122-134) DTO)
- Tambah fungsi `getUsersCount(routerId)` → `GET /hotspot/:id/users/count` (sudah ada di backend, belum di frontend)
- `deleteActiveUser` dan `deleteHost`: tambah fungsi yang belum ada

#### [MODIFY] [UsersPage.tsx](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend/src/pages/hotspot/UsersPage.tsx)

- Tambah field `timeLimit` dan `dataLimit` di form schema dan modal
- Saat edit user, populate field `timeLimit` (`limitUptime`) dan `dataLimit` dari user data
- Kirim `timeLimit` dan `dataLimit` saat create/update

#### [MODIFY] [DashboardPage.tsx](file:///home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend/src/pages/DashboardPage.tsx)

- Gunakan field `health` dari `dashboardData` untuk Voltage dan Temperature (sekarang masih dari `resources`, padahal backend mengirimkannya dalam `health` object terpisah)
- Fallback field `routerBoard.model` untuk System Info

---

## Verification Plan

### Automated — Build Check

```bash
cd /home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend
npm run build 2>&1
```

Harus selesai tanpa TypeScript error.

### Automated — Type Check

```bash
cd /home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend
npx tsc --noEmit 2>&1
```

Harus 0 error.

### Manual — Dev Server

```bash
cd /home/butterfly_student/code/Mikrotik/mikhmon_v4/frontend
npm run dev
```

Buka browser ke `http://localhost:5173` dan verifikasi:
1. Dashboard menampilkan data resource (CPU, Memory, Storage)
2. Dashboard menampilkan Voltage dan Temperature dari health data
3. System Info menampilkan model dari routerBoard
4. User management dapat membuat/update user dengan time limit dan data limit
