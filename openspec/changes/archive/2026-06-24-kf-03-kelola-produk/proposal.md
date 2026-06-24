## Why

KF-03 adalah modul master data produk di backend. Tabel `produk` direferensikan oleh `stok_masuk`, `batch_stok`, `stok_keluar`, `kemasan_terbuka`, `stok_opname`, dan `detail_opname` — hampir semua tabel transaksi bergantung padanya. Modul ini juga memperkenalkan pola penggunaan (`FULL_USE`/`PARTIAL_USE`) yang menentukan behavior sistem di KF-04, KF-05, dan KF-07.

## What Changes

- Migration tabel `produk`
- Repository: CRUD + cek stok aktif + cek riwayat transaksi sebelum hapus
- Service: validasi kelengkapan data, generate kode produk, lock `pola_penggunaan` jika ada transaksi
- Handler: 4 endpoint REST CRUD `/api/produk`
- Router: register routes produk di protected group

## Capabilities

### New Capabilities

- `produk-management`: CRUD endpoint untuk data master produk farmasi dengan validasi, kode produk otomatis, dan proteksi integritas referensial terhadap batch/transaksi

### Modified Capabilities

*(tidak ada)*

## Impact

**Modul KF terdampak:** KF-03 (Kelola Produk)

**Files yang perlu dibuat:**
- `migrations/000003_create_produk_table.up.sql` & `.down.sql`
- `internal/model/produk.go` — struct + DTOs
- `internal/repository/produk_repository.go` — CRUD + integrity checks
- `internal/service/produk_service.go` — business logic
- `internal/handler/produk_handler.go` — 4 handlers
- Update `internal/router/router.go`

**API Endpoints:**
- `GET /api/produk` → list produk + nama kategori + stok_kemasan + total_isi_tersedia
- `POST /api/produk` → tambah produk
- `PUT /api/produk/:id` → ubah produk (lock pola_penggunaan jika ada transaksi)
- `DELETE /api/produk/:id` → hapus produk (cek stok aktif + riwayat)

**Dependencies:** KF-01 (JWT middleware), KF-02 (FK ke tabel kategori)

**Acceptance Criteria:**
- AC-03.1: POST tanpa field wajib → 400 Bad Request
- AC-03.2: DELETE produk dengan stok aktif → 409 "Produk tidak dapat dihapus karena masih memiliki stok aktif."
- AC-03.3: DELETE produk dengan riwayat transaksi → 409 "Produk tidak dapat dihapus karena memiliki riwayat transaksi."
- AC-03.4: Kode produk unik, auto-generated
- AC-03.5: PUT dengan ubah pola_penggunaan saat ada transaksi → 409 "Pola penggunaan tidak dapat diubah karena produk sudah memiliki transaksi."
