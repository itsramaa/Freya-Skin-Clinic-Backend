## Why

KF-09 stock opname backend mendukung pencocokan data sistem vs fisik dengan pencatatan histori selisih sebagai jejak koreksi (audit trail). Ini adalah modul yang menyentuh paling banyak tabel dan membutuhkan DB transaction penuh untuk menjaga integritas data.

## What Changes

- Migration tabel `stok_opname` dan `detail_opname`
- Repository: create opname, get detail opname dengan item batch + kemasan terbuka, save detail, update stok batch, finalize opname
- Service: hitung selisih per item, penyesuaian stok dalam DB transaction
- Handler: `POST /api/opname`, `GET /api/opname`, `GET /api/opname/:id`, `POST /api/opname/:id/selesaikan`

## Capabilities

### New Capabilities

- `stock-opname-management`: Multi-step opname flow dengan selisih hitung otomatis, keterangan wajib, dan penyesuaian stok transaksional

### Modified Capabilities

*(tidak ada)*

## Impact

**Files:**
- `migrations/000007_create_stok_opname_table.up.sql`
- `migrations/000008_create_detail_opname_table.up.sql`
- `internal/model/opname.go`
- `internal/repository/opname_repository.go`
- `internal/service/opname_service.go`
- `internal/handler/opname_handler.go`

**Dependencies:** KF-04 (batch_stok), KF-07 (kemasan_terbuka)

**Acceptance Criteria:**
- AC-09.1: Sesi DIBATALKAN tidak mengubah stok
- AC-09.2: Keterangan wajib untuk setiap item selisih ≠ 0
- AC-09.3: Penyesuaian stok dibungkus satu DB transaction
- AC-09.4: Detail opname menyimpan stok_sistem, stok_fisik, dan selisih per item
