## Why

KF-04 adalah modul pencatatan penerimaan stok di backend. Ini menghasilkan data `stok_masuk` dan `batch_stok` — dua tabel yang menjadi fondasi seluruh sistem transaksi. Business rule kunci: merge batch jika kombinasi produk+expiredDate sudah ada, buat batch baru jika belum.

## What Changes

- Migration tabel `stok_masuk` dan `batch_stok`
- Repository: create stok masuk, find/create batch by produk+expiredDate, update batch stok
- Service: logika merge vs create batch, hitung total_isi_masuk
- Handler: `POST /api/stok-masuk`, `GET /api/stok-masuk`
- Background worker stub: update status batch KADALUWARSA (diimplementasi di KF-06)

## Capabilities

### New Capabilities

- `stok-masuk-management`: Pencatatan penerimaan stok dengan logika merge/create batch otomatis berdasarkan produk + expired date

### Modified Capabilities

*(tidak ada)*

## Impact

**Modul KF terdampak:** KF-04 (Kelola Stok Masuk)

**Files yang perlu dibuat:**
- `migrations/000004_create_stok_masuk_table.up.sql`
- `migrations/000005_create_batch_stok_table.up.sql`
- `internal/model/stok_masuk.go`, `internal/model/batch_stok.go`
- `internal/repository/stok_masuk_repository.go`, `internal/repository/batch_repository.go`
- `internal/service/stok_masuk_service.go`
- `internal/handler/stok_masuk_handler.go`

**API Endpoints:**
- `GET /api/stok-masuk` → list riwayat penerimaan dengan info produk & batch
- `POST /api/stok-masuk` → catat penerimaan baru (merge/create batch otomatis)

**Dependencies:** KF-01 (JWT), KF-02 (FK kategori via produk), KF-03 (FK produk)

**Acceptance Criteria:**
- AC-04.1: expired_date ≤ tanggal_penerimaan → 400 Bad Request
- AC-04.2: Dua POST dengan produk+expired_date sama → satu batch terakumulasi (bukan dua batch)
- AC-04.3: jumlah_kemasan ≤ 0 → 400 Bad Request
