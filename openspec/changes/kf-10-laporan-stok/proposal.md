## Why

KF-10 laporan stok backend menyediakan tiga endpoint laporan periodik: stok masuk, stok keluar, dan sisa stok terkini. Semua endpoint bersifat read-only dengan filter tanggal dan kategori/produk.

## What Changes

- Service: query laporan stok masuk (JOIN stok_masuk + batch + produk), stok keluar (JOIN stok_keluar + batch + produk), sisa stok (aggregasi dari batch_stok aktif)
- Handler: tiga GET endpoints dengan query params filter
- Tidak ada migration baru — semua data sudah ada dari tabel sebelumnya

## Capabilities

### New Capabilities

- `laporan-endpoint`: Tiga endpoint laporan stok periodik dengan filter tanggal dan kategori

### Modified Capabilities

*(tidak ada)*

## Impact

**Files:**
- `internal/handler/laporan_handler.go`
- `internal/service/laporan_service.go`
- `internal/repository/laporan_repository.go` — query kompleks JOIN multi-tabel

**API Endpoints:**
- `GET /api/laporan/stok-masuk?dari=&sampai=&kategori_id=&produk_id=`
- `GET /api/laporan/stok-keluar?dari=&sampai=&kategori_id=&produk_id=`
- `GET /api/laporan/sisa-stok?kategori_id=&produk_id=`

**Dependencies:** Semua tabel dari KF-04 s.d. KF-07

**Acceptance Criteria:**
- AC-10.1: Laporan stok masuk hanya menampilkan data dalam rentang tanggal yang diberikan
- AC-10.2: Laporan sisa stok menampilkan kondisi terkini batch aktif per produk
- AC-10.3: Semua endpoint return 401 tanpa token
