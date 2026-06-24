## Why

KF-08 monitoring stok backend menyediakan endpoint read-only yang menggabungkan data produk, batch, dan kemasan terbuka dalam satu response terstruktur. Endpoint ini mendukung filter multi-parameter dan menjadi satu-satunya data source untuk halaman monitoring frontend.

## What Changes

- Service: query aggregasi produk + batch + kemasan terbuka dengan filter params
- Handler: `GET /api/monitoring` dengan query params: kategori_id, status_batch, status_bud, nama_produk
- Status badge logic: AMAN (>30 hari), MENDEKATI (≤30 hari), KADALUWARSA (expired)

## Capabilities

### New Capabilities

- `monitoring-endpoint`: Endpoint aggregasi stok real-time dengan filter multi-parameter

### Modified Capabilities

*(tidak ada)*

## Impact

**Files:**
- `internal/handler/monitoring_handler.go`
- `internal/service/monitoring_service.go`
- `internal/repository/batch_repository.go` — tambah FindForMonitoring dengan filter

**API Endpoints:**
- `GET /api/monitoring?kategori_id=&status_batch=&status_bud=&nama_produk=`

**Acceptance Criteria:**
- AC-08.1: Filter kombinasi menghasilkan data konsisten dengan kondisi aktual DB
- AC-08.2: Response menyertakan kemasan terbuka aktif per batch yang relevan
- AC-08.3: Status indikator batch dihitung dari expired_date (AMAN/MENDEKATI/KADALUWARSA)
