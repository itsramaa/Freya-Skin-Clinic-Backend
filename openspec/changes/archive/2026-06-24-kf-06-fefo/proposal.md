## Why

KF-06 (FEFO) dan KF-07 (BUD) adalah backend-only concerns — tidak ada halaman UI tersendiri, keduanya adalah bagian dari KF-05. FEFO adalah algoritma pemilihan batch otomatis, BUD adalah pengelolaan kemasan terbuka + background worker.

## What Changes

- Implementasi `getBatchPrioritasFEFO()` di service layer stok keluar
- Algoritma: SELECT batch WHERE produk=X AND status=AKTIF ORDER BY expired_date ASC LIMIT 1
- Mendukung pemotongan lintas-batch jika satu batch tidak mencukupi
- Background worker goroutine: setiap 1 jam cek batch expired_date < NOW() → update status ke KADALUWARSA

## Capabilities

### New Capabilities

- `fefo-algorithm`: Logika FEFO otomatis terintegrasi di stok_keluar_service

### Modified Capabilities

- `stok-keluar-management`: MODIFIED — integrasi FEFO algorithm untuk pemilihan batch

## Impact

**Files terdampak:**
- `internal/service/stok_keluar_service.go` — tambah FEFO logic
- `internal/repository/batch_repository.go` — tambah FindBatchPrioritasFEFO
- `internal/service/worker_service.go` — background worker update batch status
- `cmd/api/main.go` — start goroutine worker saat startup

**Dependencies:** KF-04 (batch_stok), KF-05 (stok keluar service)

**Acceptance Criteria:**
- AC-06.1: Sistem selalu memilih batch dengan expired_date terdekat di antara batch AKTIF
- AC-06.2: Batch HABIS atau KADALUWARSA tidak masuk kandidat FEFO
- AC-06.3: Background worker mengupdate status batch KADALUWARSA dalam interval ≤ 1 jam
