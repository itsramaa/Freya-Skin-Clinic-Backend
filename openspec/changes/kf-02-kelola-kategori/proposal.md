## Why

KF-02 adalah modul master data kategori di backend. Kategori adalah dependency data master pertama — tabel `kategori` direferensikan oleh tabel `produk` (FK), sehingga harus ada sebelum KF-03 dapat diimplementasi. Endpoint CRUD kategori juga menjadi template pola handler-service-repository yang akan diikuti oleh semua modul berikutnya.

## What Changes

- Migration tabel `kategori`
- Repository layer: CRUD + cek relasi produk sebelum hapus
- Service layer: validasi duplikasi nama (case-insensitive), validasi ketergantungan produk
- Handler layer: 4 endpoint REST CRUD `/api/kategori`
- Router: daftarkan routes kategori di protected group

## Capabilities

### New Capabilities

- `kategori-management`: CRUD endpoint untuk data master kategori farmasi dengan validasi duplikasi dan proteksi integritas referensial

### Modified Capabilities

*(tidak ada)*

## Impact

**Modul KF terdampak:** KF-02 (Kelola Kategori)

**Files yang perlu dibuat:**
- `migrations/000002_create_kategori_table.up.sql` & `.down.sql`
- `internal/model/kategori.go` — struct Kategori, CreateKategoriRequest, UpdateKategoriRequest, KategoriResponse
- `internal/repository/kategori_repository.go` — interface + implementasi CRUD
- `internal/service/kategori_service.go` — interface + implementasi business logic
- `internal/handler/kategori_handler.go` — 4 handler functions
- Update `internal/router/router.go` — register kategori routes

**API Endpoints:**
- `GET /api/kategori` → list semua kategori dengan jumlah produk terkait
- `POST /api/kategori` → tambah kategori baru
- `PUT /api/kategori/:id` → ubah nama kategori
- `DELETE /api/kategori/:id` → hapus kategori (gagal jika ada produk terkait)

**Dependencies:** KF-01 (middleware JWT harus sudah ada)

**Acceptance Criteria:**
- AC-02.1: POST kategori dengan nama yang sama (case-insensitive) → 409 Conflict "Nama kategori sudah terdaftar dalam sistem."
- AC-02.2: DELETE kategori yang masih punya produk terkait → 409 Conflict "Kategori tidak dapat dihapus karena masih memiliki produk terkait."
- AC-02.3: GET /api/kategori mengembalikan list dengan field jumlah_produk yang akurat
- AC-02.4: Semua endpoint kategori return 401 tanpa Bearer token
