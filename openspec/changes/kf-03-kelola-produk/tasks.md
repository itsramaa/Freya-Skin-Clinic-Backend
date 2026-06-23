## 1. Migration & Database

- [ ] 1.1 Buat `migrations/000003_create_produk_table.up.sql` — CREATE TYPE `pola_penggunaan_enum` AS ENUM ('FULL_USE', 'PARTIAL_USE'); CREATE TABLE produk (id UUID, kode_produk VARCHAR(20) UNIQUE, nama_produk VARCHAR(200), id_kategori UUID FK→kategori, bentuk_kemasan VARCHAR(50), satuan_isi VARCHAR(20), isi_per_kemasan DECIMAL(10,3) NULLABLE, pola_penggunaan pola_penggunaan_enum, created_at, updated_at)
- [ ] 1.2 Buat `migrations/000003_create_produk_table.down.sql` — DROP TABLE produk; DROP TYPE pola_penggunaan_enum
- [ ] 1.3 Tambah index: `idx_produk_id_kategori` pada kolom id_kategori, `idx_produk_nama` pada nama_produk
- [ ] 1.4 Jalankan `go run cmd/migrate/main.go up` — verifikasi migration dan ENUM terbuat

## 2. Model & DTO

- [ ] 2.1 Buat `internal/model/produk.go` — struct: `Produk` (semua kolom DB), `CreateProdukRequest` (nama_produk, id_kategori, bentuk_kemasan, satuan_isi, isi_per_kemasan *float64, pola_penggunaan), `UpdateProdukRequest` (same), `ProdukResponse` (+ nama_kategori string, stok_kemasan int, total_isi_tersedia float64, has_transaksi bool)

## 3. Repository Layer

- [ ] 3.1 Buat `internal/repository/produk_repository.go` — interface `ProdukRepository` dengan methods: `FindAll`, `FindByID`, `Create`, `Update`, `Delete`, `HasTransaksi`, `CountStokAktif`, `CountTransaksi`
- [ ] 3.2 Implementasi `FindAll(ctx)` — query dengan JOIN kategori + LEFT JOIN batch_stok + subquery EXISTS(stok_masuk) untuk has_transaksi; GROUP BY produk.id
- [ ] 3.3 Implementasi `FindByID(ctx, id)` — SELECT by UUID
- [ ] 3.4 Implementasi `Create(ctx, produk)` — INSERT dengan kode_produk yang sudah di-generate
- [ ] 3.5 Implementasi `Update(ctx, id, req)` — UPDATE semua field yang dikirim
- [ ] 3.6 Implementasi `Delete(ctx, id)` — DELETE by UUID
- [ ] 3.7 Implementasi `HasTransaksi(ctx, id)` — `SELECT EXISTS(SELECT 1 FROM stok_masuk WHERE id_produk = $1)`
- [ ] 3.8 Implementasi `CountStokAktif(ctx, id)` — `SELECT COUNT(*) FROM batch_stok WHERE id_produk = $1 AND status = 'AKTIF'`
- [ ] 3.9 Implementasi `CountTransaksi(ctx, id)` — `SELECT COUNT(*) FROM stok_masuk WHERE id_produk = $1`

## 4. Service Layer

- [ ] 4.1 Buat `internal/service/produk_service.go` — interface `ProdukService` dengan methods: `GetAll`, `Create`, `Update`, `Delete`
- [ ] 4.2 Implementasi `GetAll(ctx)` — call repo.FindAll, return []ProdukResponse
- [ ] 4.3 Implementasi `Create(ctx, req)` — validasi field wajib; validasi isi_per_kemasan wajib jika PARTIAL_USE; cek kategori exists; generate kode_produk (PRD-{3-char-prefix}-{seq}); call repo.Create
- [ ] 4.4 Implementasi `Update(ctx, id, req)` — FindByID (404); HasTransaksi + cek pola_penggunaan berubah (409); validasi isi_per_kemasan jika PARTIAL_USE; call repo.Update
- [ ] 4.5 Implementasi `Delete(ctx, id)` — FindByID (404); CountStokAktif (409); CountTransaksi (409); call repo.Delete

## 5. Handler Layer

- [ ] 5.1 Buat `internal/handler/produk_handler.go` — struct `ProdukHandler` dengan dependency `ProdukService` dan `KategoriService` (untuk validasi id_kategori)
- [ ] 5.2 Implementasi `GetAll(c *fiber.Ctx) error` — call service.GetAll, return 200
- [ ] 5.3 Implementasi `Create(c *fiber.Ctx) error` — parse body, validasi wajib, call service.Create, return 201
- [ ] 5.4 Implementasi `Update(c *fiber.Ctx) error` — parse body + path param, call service.Update, return 200
- [ ] 5.5 Implementasi `Delete(c *fiber.Ctx) error` — parse path param, call service.Delete, return 200

## 6. Router

- [ ] 6.1 Update `internal/router/router.go` — tambah routes di protected group: `GET /api/produk`, `POST /api/produk`, `PUT /api/produk/:id`, `DELETE /api/produk/:id`; inject ProdukHandler

## 7. Verifikasi

- [ ] 7.1 Verifikasi `GET /api/produk` — response 200 dengan field has_transaksi, stok_kemasan, total_isi_tersedia
- [ ] 7.2 Verifikasi `POST /api/produk` FULL_USE tanpa isi_per_kemasan — response 201 (AC-03.1)
- [ ] 7.3 Verifikasi `POST /api/produk` PARTIAL_USE tanpa isi_per_kemasan — response 400 (AC-03.1)
- [ ] 7.4 Verifikasi `PUT /api/produk/:id` ganti pola_penggunaan saat has_transaksi=true — response 409 (AC-03.5)
- [ ] 7.5 Verifikasi `DELETE /api/produk/:id` produk dengan stok aktif — response 409 (AC-03.2)
- [ ] 7.6 Verifikasi `DELETE /api/produk/:id` produk dengan riwayat transaksi — response 409 (AC-03.3)
- [ ] 7.7 Verifikasi kode_produk unik dan ter-generate otomatis (AC-03.4)
- [ ] 7.8 Verifikasi semua endpoint return 401 tanpa token
- [ ] 7.9 Jalankan `go build ./...` — tidak ada compile error
