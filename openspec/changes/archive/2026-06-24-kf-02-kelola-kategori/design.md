## Context

KF-02 adalah modul CRUD data master kategori di backend. Migration dan implementasi modul ini menjadi template pola handler-service-repository yang akan diikuti oleh KF-03 s.d. KF-10.

Tabel `kategori` memiliki FK constraint dari `produk.id_kategori` — penghapusan harus dicek di level service sebelum query DELETE.

## Goals / Non-Goals

**Goals:**
- Migration tabel `kategori` dengan kode_kategori auto-generate
- Repository: FindAll (dengan count produk), FindByID, FindByNama, Create, Update, Delete, CountProdukByKategori
- Service: validasi duplikasi nama (case-insensitive), validasi no produk terkait, generate kode_kategori
- Handler: 4 endpoints REST
- Router: register ke protected group

**Non-Goals:**
- Soft delete
- Pagination/filter
- Bulk operations

## Decisions

### 1. Kode Kategori: Format KTG-XXX dengan LPAD sequence

**Keputusan:** Kode kategori di-generate dengan format `KTG-001`, `KTG-002`, dst menggunakan sequence counter dari database.

**Alasan:** Kode harus human-readable dan konsisten. Menggunakan sequence DB (bukan UUID) agar kode bisa tampil di UI dengan urutan yang mudah dibaca.

**Implementasi:** Gunakan sequence PostgreSQL atau hitung `MAX(kode_kategori)` saat insert.

### 2. Duplikasi Check: Case-insensitive dengan LOWER()

**Keputusan:** Cek duplikasi nama menggunakan `LOWER(nama_kategori) = LOWER($1)` di query.

**Alasan:** "Skincare" dan "skincare" harus dianggap sama sesuai BR-02.1. Implementasi di repository layer — bukan di service — agar atomic dengan query.

### 3. Cek Produk Terkait: Query COUNT di repository

**Keputusan:** Method `CountProdukByKategoriID(ctx, id)` yang query `SELECT COUNT(*) FROM produk WHERE id_kategori = $1`.

**Alasan:** Service layer memanggil ini sebelum delete — jika count > 0 return error tanpa query DELETE.

### 4. Response jumlah_produk: JOIN query di FindAll

**Keputusan:** `GET /api/kategori` menggunakan LEFT JOIN dengan tabel `produk` untuk count jumlah produk per kategori dalam satu query.

**Alasan:** Menghindari N+1 query — satu query untuk semua kategori + count produk sekaligus.

## API Endpoints

| Method | Path | Auth | Request | Response |
|--------|------|------|---------|---------|
| `GET` | `/api/kategori` | Bearer | — | `200: { success, data: [{ id, kode_kategori, nama_kategori, jumlah_produk }] }` |
| `POST` | `/api/kategori` | Bearer | `{ nama_kategori }` | `201: { success, data: kategori }` |
| `PUT` | `/api/kategori/:id` | Bearer | `{ nama_kategori }` | `200: { success, data: kategori }` |
| `DELETE` | `/api/kategori/:id` | Bearer | — | `200: { success, message }` |

## Database Schema

```sql
-- migrations/000002_create_kategori_table.up.sql
CREATE TABLE kategori (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode_kategori   VARCHAR(10) UNIQUE NOT NULL,
    nama_kategori   VARCHAR(100) UNIQUE NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_kategori_nama ON kategori(LOWER(nama_kategori));
```

## Architecture Flow

```
GET /api/kategori
  → JWTMiddleware
  → kategori_handler.GetAll()
      → kategori_service.GetAll(ctx)
          → kategori_repo.FindAll(ctx)
              → SELECT k.*, COUNT(p.id) as jumlah_produk FROM kategori k
                LEFT JOIN produk p ON p.id_kategori = k.id
                GROUP BY k.id ORDER BY k.kode_kategori
          → return []KategoriResponse
      → return 200 { success: true, data: [...] }

POST /api/kategori
  → JWTMiddleware
  → kategori_handler.Create()
      → Parse & validate body (nama_kategori wajib)
      → kategori_service.Create(ctx, req)
          → kategori_repo.FindByNama(ctx, namaKategori) → jika ada → return error duplikat
          → Generate kode_kategori (KTG-XXX)
          → kategori_repo.Create(ctx, kategori)
      → return 201 { success: true, data: kategori }

PUT /api/kategori/:id
  → JWTMiddleware
  → kategori_handler.Update()
      → Parse & validate body + path param id
      → kategori_service.Update(ctx, id, req)
          → kategori_repo.FindByID(ctx, id) → jika tidak ada → return 404
          → kategori_repo.FindByNama(ctx, namaKategori) → jika ada dan bukan id ini → return error duplikat
          → kategori_repo.Update(ctx, id, namaKategori)
      → return 200 { success: true, data: kategori }

DELETE /api/kategori/:id
  → JWTMiddleware
  → kategori_handler.Delete()
      → kategori_service.Delete(ctx, id)
          → kategori_repo.FindByID(ctx, id) → jika tidak ada → return 404
          → kategori_repo.CountProdukByKategoriID(ctx, id) → jika > 0 → return error produk terkait
          → kategori_repo.Delete(ctx, id)
      → return 200 { success: true, message: "Kategori berhasil dihapus." }
```

## File Structure

```
internal/
├── model/
│   └── kategori.go                  # Kategori struct, CreateKategoriRequest, UpdateKategoriRequest, KategoriResponse
├── repository/
│   └── kategori_repository.go       # Interface + implementasi (FindAll, FindByID, FindByNama, Create, Update, Delete, CountProduk)
├── service/
│   └── kategori_service.go          # Interface + implementasi business logic
├── handler/
│   └── kategori_handler.go          # GetAll, Create, Update, Delete handlers
└── router/
    └── router.go                    # UPDATED: register kategori routes

migrations/
├── 000002_create_kategori_table.up.sql
└── 000002_create_kategori_table.down.sql
```

## Risks / Trade-offs

- **[Risk] Race condition duplikasi nama** → Mitigasi: UNIQUE constraint `nama_kategori` di DB sebagai safety net — service layer check + DB constraint double protection.
- **[Risk] Kode kategori collision jika concurrent inserts** → Mitigasi: Generate kode menggunakan sequence atau MAX+1 dalam transaction. Untuk single-user system, ini cukup.

## Migration Plan

1. Jalankan migration 000002 setelah 000001 (users table)
2. Seed 5 kategori default (opsional): Skincare, Injectable, Obat, Threadlift, Facial IPL Laser
3. Verifikasi endpoints dengan Postman/curl

## Open Questions

- Apakah perlu seed 5 kategori default di migration? (memudahkan testing KF-03 langsung)
