## Why

KF-07 (BUD) adalah backend-only concern — pengelolaan kemasan terbuka untuk produk Partial Use. BUD ditetapkan otomatis 28 hari sejak kemasan dibuka, dipantau oleh background worker, dan dikecek saat stok keluar Partial Use diproses.

## What Changes

- Logika BUD di `stok_keluar_service`: cek kemasan terbuka aktif, tetapkan BUD 28 hari, nonaktifkan jika expired
- Background worker goroutine: setiap 1 jam cek kemasan_terbuka WHERE bud < NOW() AND status = AKTIF → update status ke KADALUWARSA
- Repository: FindKemasanTerbukaAktifByBatch, CreateKemasanTerbuka, UpdateKemasanTerbuka, UpdateStatusExpiredBUD

## Capabilities

### New Capabilities

- `bud-management`: Pengelolaan kemasan terbuka dengan BUD otomatis 28 hari dan background worker monitoring

### Modified Capabilities

- `stok-keluar-management`: MODIFIED — integrasi BUD logic untuk Partial Use flow

## Impact

**Files terdampak:**
- `internal/repository/kemasan_terbuka_repository.go` — CRUD + FindAktifByBatch
- `internal/service/stok_keluar_service.go` — BUD logic di Partial Use flow
- `internal/service/worker_service.go` — goroutine update status BUD expired
- `migrations/000006_create_kemasan_terbuka_table.up.sql`

**Dependencies:** KF-04 (batch_stok), KF-05 (stok_keluar)

**Acceptance Criteria:**
- AC-07.1: Produk Full Use tidak pernah menghasilkan baris kemasan_terbuka
- AC-07.2: BUD = tanggal_dibuka + 28 hari (selalu)
- AC-07.3: Kemasan terbuka BUD expired → status KADALUWARSA via worker atau interaksi user
