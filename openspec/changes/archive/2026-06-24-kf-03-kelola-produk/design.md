## Context

KF-03 adalah modul CRUD produk di backend. Tabel `produk` memiliki FK dari hampir semua tabel transaksi. Modul ini memperkenalkan ENUM `pola_penggunaan` yang menentukan behavior di KF-04 s.d. KF-07.

Kompleksitas utama: cek `has_transaksi` yang melihat ke beberapa tabel (stok_masuk, stok_keluar), lock perubahan `pola_penggunaan`, dan generate kode produk per-kategori.

## Goals / Non-Goals

**Goals:**
- Migration tabel `produk` dengan ENUM `pola_penggunaan`
- Repository: CRUD + CountStokAktif + CountTransaksi + HasTransaksi
- Service: validasi, generate kode produk, lock pola_penggunaan
- Handler: 4 endpoints
- Router update

**Non-Goals:**
- Soft delete
- Pagination/filter/search
- Bulk import

## Decisions

### 1. ENUM pola_penggunaan: PostgreSQL native ENUM

**Keputusan:** Gunakan PostgreSQL ENUM type `pola_penggunaan_enum` dengan values `FULL_USE` dan `PARTIAL_USE`.

**Alasan:** Native ENUM memberikan constraint di level database, lebih aman dari string arbitrary. Go dapat scan ENUM sebagai `string` dengan pgx.

### 2. Kode Produk: Format PRD-{KodeKategori}-{Sequence}

**Keputusan:** Generate kode produk dengan format `PRD-SKC-001` (contoh untuk kategori Skincare dengan kode KTG-001 → prefix SKC dari 3 huruf pertama nama).

**Alternatif:** UUID-based — ditolak karena tidak human-readable untuk tampilan UI.

**Implementasi:** Ambil 3 huruf pertama nama kategori (uppercase) + sequence count produk dalam kategori itu (`SELECT COUNT(*) FROM produk WHERE id_kategori = $1`).

### 3. Lock pola_penggunaan: Cek has_transaksi di service layer

**Keputusan:** Method `HasTransaksi(ctx, produkID)` di repository — query `SELECT EXISTS(SELECT 1 FROM stok_masuk WHERE id_produk = $1)`.

**Alasan:** Cek hanya di service layer, bukan di database constraint. Lebih fleksibel dan tidak memerlukan trigger.

### 4. has_transaksi di response GET list: Subquery atau computed

**Keputusan:** Include `has_transaksi` sebagai boolean computed dari subquery EXISTS di query FindAll.

**Alasan:** Frontend butuh info ini untuk disable field saat edit — lebih efisien di-compute di DB daripada N+1 dari service.

### 5. isi_per_kemasan: Nullable untuk FULL_USE

**Keputusan:** Kolom `isi_per_kemasan` di tabel `produk` bertipe `DECIMAL(10,3)` NULLABLE.

**Alasan:** FULL_USE tidak membutuhkan isi_per_kemasan. Validasi wajib isi hanya dilakukan di service layer untuk PARTIAL_USE.

## API Endpoints

| Method | Path | Request Body | Response |
|--------|------|-------------|---------|
| `GET` | `/api/produk` | — | `200: [{ id, kode_produk, nama_produk, id_kategori, nama_kategori, bentuk_kemasan, satuan_isi, isi_per_kemasan, pola_penggunaan, stok_kemasan, total_isi_tersedia, has_transaksi }]` |
| `POST` | `/api/produk` | `{ nama_produk, id_kategori, bentuk_kemasan, satuan_isi, isi_per_kemasan?, pola_penggunaan }` | `201: { success, data: produk }` |
| `PUT` | `/api/produk/:id` | same as POST | `200: { success, data: produk }` |
| `DELETE` | `/api/produk/:id` | — | `200: { success, message }` |

## Database Schema

```sql
-- migrations/000003_create_produk_table.up.sql
CREATE TYPE pola_penggunaan_enum AS ENUM ('FULL_USE', 'PARTIAL_USE');

CREATE TABLE produk (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode_produk      VARCHAR(20) UNIQUE NOT NULL,
    nama_produk      VARCHAR(200) NOT NULL,
    id_kategori      UUID NOT NULL,
    bentuk_kemasan   VARCHAR(50) NOT NULL,
    satuan_isi       VARCHAR(20) NOT NULL,
    isi_per_kemasan  DECIMAL(10,3),
    pola_penggunaan  pola_penggunaan_enum NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_produk_id_kategori FOREIGN KEY (id_kategori)
        REFERENCES kategori(id) ON DELETE RESTRICT
);

CREATE INDEX idx_produk_id_kategori ON produk(id_kategori);
CREATE INDEX idx_produk_nama ON produk(nama_produk);
```

## Architecture Flow

```
GET /api/produk
  → JWTMiddleware → produk_handler.GetAll()
      → produk_service.GetAll(ctx)
          → produk_repo.FindAll(ctx)
              → SELECT p.*, k.nama_kategori,
                  COALESCE(SUM(b.jumlah_kemasan) FILTER (WHERE b.status='AKTIF'), 0) as stok_kemasan,
                  COALESCE(SUM(b.sisa_isi) FILTER (WHERE b.status='AKTIF'), 0) as total_isi_tersedia,
                  EXISTS(SELECT 1 FROM stok_masuk sm WHERE sm.id_produk = p.id) as has_transaksi
                FROM produk p JOIN kategori k ON p.id_kategori = k.id
                LEFT JOIN batch_stok b ON b.id_produk = p.id
                GROUP BY p.id, k.nama_kategori

POST /api/produk
  → JWTMiddleware → produk_handler.Create()
      → Validate body (nama, id_kategori, bentuk_kemasan, satuan_isi, pola_penggunaan wajib;
                       isi_per_kemasan wajib jika PARTIAL_USE)
      → produk_service.Create(ctx, req)
          → kategori_repo.FindByID(ctx, id_kategori) → 404 check
          → Generate kode_produk (PRD-{PREFIX}-{SEQ})
          → produk_repo.Create(ctx, produk)
      → return 201

PUT /api/produk/:id
  → JWTMiddleware → produk_handler.Update()
      → produk_service.Update(ctx, id, req)
          → produk_repo.FindByID → 404 check
          → produk_repo.HasTransaksi(ctx, id) → jika true dan pola berubah → return 409
          → produk_repo.Update(ctx, id, req)
      → return 200

DELETE /api/produk/:id
  → JWTMiddleware → produk_handler.Delete()
      → produk_service.Delete(ctx, id)
          → produk_repo.FindByID → 404 check
          → produk_repo.CountStokAktif(ctx, id) → jika > 0 → 409 stok aktif
          → produk_repo.CountTransaksi(ctx, id) → jika > 0 → 409 riwayat transaksi
          → produk_repo.Delete(ctx, id)
      → return 200
```

## File Structure

```
internal/
├── model/produk.go
├── repository/produk_repository.go   # FindAll, FindByID, Create, Update, Delete, HasTransaksi, CountStokAktif, CountTransaksi
├── service/produk_service.go
├── handler/produk_handler.go
└── router/router.go                  # UPDATED

migrations/
├── 000003_create_produk_table.up.sql
└── 000003_create_produk_table.down.sql
```

## Risks / Trade-offs

- **[Risk] Query FindAll kompleks (multiple JOIN + subquery)** → Mitigasi: stok_kemasan dan total_isi_tersedia dihitung dari batch_stok — query tetap satu round-trip. Index pada `id_produk` di batch_stok penting untuk performa.
- **[Risk] Kode produk collision concurrent insert** → Mitigasi: UNIQUE constraint di DB sebagai safety net. Single-user system — race condition sangat jarang terjadi.

## Migration Plan

1. Jalankan migration 000003 setelah 000002
2. Verifikasi ENUM `pola_penggunaan_enum` terbuat
3. Verifikasi FK constraint ke tabel kategori aktif

## Open Questions

- Format 3-huruf prefix kode produk: otomatis dari nama kategori atau dikonfigurasi manual per kategori?
