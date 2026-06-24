## 1. Migration & Database

- [ ] 1.1 Buat `migrations/000002_create_kategori_table.up.sql` — DDL tabel `kategori` (id UUID, kode_kategori VARCHAR(10) UNIQUE, nama_kategori VARCHAR(100) UNIQUE, created_at, updated_at) + index `idx_kategori_nama` pada `LOWER(nama_kategori)`
- [ ] 1.2 Buat `migrations/000002_create_kategori_table.down.sql` — DROP TABLE kategori
- [ ] 1.3 (Opsional) Tambah seed 5 kategori default: Skincare, Injectable, Obat, Threadlift, Facial IPL Laser
- [ ] 1.4 Jalankan `go run cmd/migrate/main.go up` — verifikasi migration berhasil

## 2. Model & DTO

- [ ] 2.1 Buat `internal/model/kategori.go` — struct: `Kategori` (id, kode_kategori, nama_kategori, created_at, updated_at), `CreateKategoriRequest` (nama_kategori), `UpdateKategoriRequest` (nama_kategori), `KategoriResponse` (+ jumlah_produk int)

## 3. Repository Layer

- [ ] 3.1 Buat `internal/repository/kategori_repository.go` — interface `KategoriRepository` dengan methods: `FindAll`, `FindByID`, `FindByNama`, `Create`, `Update`, `Delete`, `CountProdukByKategoriID`
- [ ] 3.2 Implementasi `FindAll(ctx)` — query LEFT JOIN dengan tabel produk untuk hitung jumlah_produk per kategori dalam satu query
- [ ] 3.3 Implementasi `FindByID(ctx, id)` — SELECT by UUID
- [ ] 3.4 Implementasi `FindByNama(ctx, nama)` — query `WHERE LOWER(nama_kategori) = LOWER($1)` untuk case-insensitive check
- [ ] 3.5 Implementasi `Create(ctx, kategori)` — INSERT dengan kode_kategori yang di-generate
- [ ] 3.6 Implementasi `Update(ctx, id, namaKategori)` — UPDATE nama_kategori dan updated_at
- [ ] 3.7 Implementasi `Delete(ctx, id)` — DELETE by UUID
- [ ] 3.8 Implementasi `CountProdukByKategoriID(ctx, id)` — `SELECT COUNT(*) FROM produk WHERE id_kategori = $1`

## 4. Service Layer

- [ ] 4.1 Buat `internal/service/kategori_service.go` — interface `KategoriService` dengan methods: `GetAll`, `Create`, `Update`, `Delete`
- [ ] 4.2 Implementasi `GetAll(ctx)` — call repo.FindAll, return []KategoriResponse
- [ ] 4.3 Implementasi `Create(ctx, req)` — call repo.FindByNama (duplikat check), generate kode_kategori (`KTG-` + LPAD sequence), call repo.Create
- [ ] 4.4 Implementasi `Update(ctx, id, req)` — call repo.FindByID (404 check), call repo.FindByNama (duplikat check exclude self), call repo.Update
- [ ] 4.5 Implementasi `Delete(ctx, id)` — call repo.FindByID (404 check), call repo.CountProdukByKategoriID (409 check), call repo.Delete

## 5. Handler Layer

- [ ] 5.1 Buat `internal/handler/kategori_handler.go` — struct `KategoriHandler` dengan dependency `KategoriService`
- [ ] 5.2 Implementasi `GetAll(c *fiber.Ctx) error` — call service.GetAll, return 200 response
- [ ] 5.3 Implementasi `Create(c *fiber.Ctx) error` — parse & validate body (nama_kategori wajib), call service.Create, return 201
- [ ] 5.4 Implementasi `Update(c *fiber.Ctx) error` — parse body + path param id, call service.Update, return 200
- [ ] 5.5 Implementasi `Delete(c *fiber.Ctx) error` — parse path param id, call service.Delete, return 200

## 6. Router

- [ ] 6.1 Update `internal/router/router.go` — tambah routes di protected group: `GET /api/kategori`, `POST /api/kategori`, `PUT /api/kategori/:id`, `DELETE /api/kategori/:id`; inject KategoriHandler dependency

## 7. Verifikasi

- [ ] 7.1 Verifikasi `GET /api/kategori` — response 200 dengan array kategori + jumlah_produk (AC-02.3)
- [ ] 7.2 Verifikasi `POST /api/kategori` nama baru — response 201 dengan kode_kategori ter-generate
- [ ] 7.3 Verifikasi `POST /api/kategori` nama duplikat (case-insensitive) — response 409 (AC-02.1)
- [ ] 7.4 Verifikasi `PUT /api/kategori/:id` nama valid — response 200
- [ ] 7.5 Verifikasi `PUT /api/kategori/:id` nama duplikat — response 409
- [ ] 7.6 Verifikasi `DELETE /api/kategori/:id` tanpa produk — response 200
- [ ] 7.7 Verifikasi `DELETE /api/kategori/:id` dengan produk terkait — response 409 (AC-02.2)
- [ ] 7.8 Verifikasi semua endpoint return 401 tanpa token (AC-02.4)
- [ ] 7.9 Jalankan `go build ./...` — tidak ada compile error
