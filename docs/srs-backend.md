# SRS — Backend Specification
## Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic

---

## Daftar Isi

1. [Tujuan & Cakupan](#1-tujuan--cakupan)
2. [Justifikasi Teknologi](#2-justifikasi-teknologi)
3. [Struktur Proyek](#3-struktur-proyek)
4. [Arsitektur Lapisan (Handler → Service → Repository)](#4-arsitektur-lapisan-handler--service--repository)
5. [Middleware Autentikasi](#5-middleware-autentikasi)
6. [Modul Handler](#6-modul-handler)
7. [Modul Service (Logika Bisnis)](#7-modul-service-logika-bisnis)
8. [Repository Layer](#8-repository-layer)
9. [Background Worker](#9-background-worker)
10. [Penanganan Error & Format Respons](#10-penanganan-error--format-respons)
11. [Konvensi Kode](#11-konvensi-kode)

---

## 1. Tujuan & Cakupan

Dokumen ini merinci spesifikasi lapisan logika bisnis (backend) Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic, yang dibangun menggunakan **Go dengan Fiber** sebagai web framework sesuai arsitektur yang ditetapkan pada `srs-overview.md` § 8.2. Backend berfungsi sebagai pusat pemrosesan seluruh logika bisnis, autentikasi, serta pengelolaan data yang dikomunikasikan ke frontend melalui REST API (lihat `srs-api.md`) dan ke basis data melalui repository layer (lihat `srs-database.md`).

---

## 2. Justifikasi Teknologi

Lapisan logika bisnis dibangun menggunakan **Go dengan Fiber** sebagai web framework. Fiber menyediakan routing dan middleware untuk penanganan request API, termasuk middleware autentikasi berbasis token yang mendukung pembatasan akses sesuai hak pengguna yang telah ditentukan.

Pemilihan Go didasarkan pada kebutuhan sistem akan proses yang berjalan di luar alur utama penanganan permintaan pengguna, yaitu:

- **Pemantauan status batch** — pengecekan batch yang telah melewati `expiredDate` untuk diperbarui statusnya menjadi `KADALUWARSA` secara otomatis (KF-06, KNF-07).
- **Pemantauan status BUD** — pengecekan kemasan terbuka yang telah melewati batas BUD (28 hari) untuk dinonaktifkan secara otomatis (KF-07, KNF-07).

Go menyediakan model konkurensi berbasis **goroutine dan channel** yang memungkinkan kedua proses pemantauan tersebut berjalan sebagai background worker secara terpisah dari proses penanganan permintaan utama, sehingga validasi data dan respons terhadap pengguna tidak terhambat oleh proses pemantauan yang berjalan di latar belakang. Selain itu, pemilihan Go konsisten dengan kebutuhan integritas data transaksi stok yang memerlukan pengelolaan DB transaction secara eksplisit (KNF-04).

---

## 3. Struktur Proyek

```
backend/
├── main.go                        ← Entry point: init app, routes, background worker
├── go.mod
├── go.sum
├── config/
│   └── config.go                  ← Konfigurasi env (DB DSN, JWT secret, port)
├── middleware/
│   └── auth.go                    ← Middleware JWT: validasi token tiap request
├── handler/
│   ├── auth_handler.go            ← POST /api/auth/login, PUT /api/auth/password
│   ├── kategori_handler.go        ← CRUD /api/kategori
│   ├── produk_handler.go          ← CRUD /api/produk
│   ├── stok_masuk_handler.go      ← POST /api/stok-masuk, GET list & detail
│   ├── stok_keluar_handler.go     ← POST /api/stok-keluar, GET preview-batch
│   ├── monitoring_handler.go      ← GET /api/monitoring
│   ├── opname_handler.go          ← POST /api/opname, GET list & detail
│   └── laporan_handler.go         ← GET /api/laporan (stok-masuk, stok-keluar, sisa-stok)
├── service/
│   ├── auth_service.go
│   ├── kategori_service.go
│   ├── produk_service.go
│   ├── stok_masuk_service.go
│   ├── stok_keluar_service.go
│   ├── monitoring_service.go
│   ├── opname_service.go
│   ├── laporan_service.go
│   └── worker_service.go          ← Logika pembaruan status batch & BUD
├── repository/
│   ├── user_repository.go
│   ├── kategori_repository.go
│   ├── produk_repository.go
│   ├── stok_masuk_repository.go
│   ├── batch_repository.go
│   ├── stok_keluar_repository.go
│   ├── kemasan_terbuka_repository.go
│   ├── opname_repository.go
│   └── laporan_repository.go
├── model/
│   ├── user.go
│   ├── kategori.go
│   ├── produk.go
│   ├── stok_masuk.go
│   ├── batch_stok.go
│   ├── stok_keluar.go
│   ├── kemasan_terbuka.go
│   ├── stok_opname.go
│   └── detail_opname.go
└── util/
    ├── jwt.go                     ← Generate & parse JWT token
    ├── hash.go                    ← Bcrypt hash & verify password
    └── response.go                ← Helper format JSON respons standar
```

---

## 4. Arsitektur Lapisan (Handler → Service → Repository)

Backend mengikuti pemisahan tiga lapisan yang tegas, sesuai KNF-05 (Maintainability):

```
Request HTTP
     │
     ▼
┌─────────────┐
│   Handler   │  ← Parsing request, validasi format input, panggil service, kembalikan respons
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Service   │  ← Seluruh logika bisnis: validasi aturan bisnis, orkestrasi DB transaction,
│             │    penerapan FEFO, pengelolaan BUD, kalkulasi stok
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Repository  │  ← Akses database: query SQL via sqlx/pgx, tidak ada logika bisnis di sini
└──────┬──────┘
       │
       ▼
  PostgreSQL
```

**Prinsip:**
- Handler **tidak** memanggil repository secara langsung; seluruh koordinasi ada di service.
- Repository **tidak** mengandung logika bisnis; hanya menerima parameter dan mengembalikan data/error.
- DB transaction dikelola di level service menggunakan `sql.Tx` yang diteruskan ke repository.

---

## 5. Middleware Autentikasi

### 5.1 Deskripsi

Middleware autentikasi (`middleware/auth.go`) memvalidasi token JWT pada setiap request ke endpoint yang dilindungi sebelum request mencapai handler, sesuai KNF-02 (Keamanan Akses).

### 5.2 Alur Validasi

```
Request masuk
     │
     ▼
Ambil header Authorization: Bearer <token>
     │
     ├── Tidak ada token ──────────────► HTTP 401 Unauthorized
     │
     ▼
Parse & validasi JWT (secret, expiry, signature)
     │
     ├── Token tidak valid / expired ──► HTTP 401 Unauthorized
     │
     ▼
Simpan payload (idUser) ke context Fiber
     │
     ▼
Lanjut ke handler
```

### 5.3 Endpoint yang Tidak Diproteksi

Hanya satu endpoint yang dapat diakses tanpa token:

| Endpoint | Alasan |
|---|---|
| `POST /api/auth/login` | Entry point autentikasi |

Seluruh endpoint lainnya (termasuk `PUT /api/auth/password`) **wajib** diproteksi middleware ini.

### 5.4 Konfigurasi Token

| Parameter | Nilai |
|---|---|
| Algoritma | HS256 |
| Masa berlaku (expiry) | Dikonfigurasi via env `JWT_EXPIRY` (rekomendasi: 8 jam sesuai jam operasional klinik) |
| Claim wajib | `sub` (idUser), `exp`, `iat` |
| Secret | Dikonfigurasi via env `JWT_SECRET` (tidak di-hardcode) |

---

## 6. Modul Handler

Handler bertanggung jawab atas:
1. Parsing body/query parameter dari request.
2. Validasi format input dasar (field wajib tidak kosong, tipe data sesuai).
3. Memanggil service yang sesuai.
4. Mengembalikan respons JSON terstandar (lihat § 10).

Handler **tidak** mengandung logika bisnis — validasi aturan bisnis (duplikasi, ketersediaan stok, FEFO) sepenuhnya ada di service.

### 6.1 Auth Handler (`auth_handler.go`)

| Fungsi | Method | Path | Terkait KF |
|---|---|---|---|
| `Login` | POST | `/api/auth/login` | KF-01 |
| `GantiPassword` | PUT | `/api/auth/password` | KF-01 |

**Login:**
- Membaca `username` dan `password` dari body.
- Memanggil `authService.Login(username, password)`.
- Jika berhasil: kembalikan token JWT + flag `isDefaultPassword`.
- Jika gagal: kembalikan HTTP 401 dengan pesan `"Kredensial tidak valid."`.

**GantiPassword:**
- Membaca `passwordBaru` dari body; `idUser` diambil dari context (token).
- Memanggil `authService.GantiPassword(idUser, passwordBaru)`.

### 6.2 Kategori Handler (`kategori_handler.go`)

| Fungsi | Method | Path | Terkait KF |
|---|---|---|---|
| `GetAllKategori` | GET | `/api/kategori` | KF-02 |
| `CreateKategori` | POST | `/api/kategori` | KF-02 |
| `UpdateKategori` | PUT | `/api/kategori/:id` | KF-02 |
| `DeleteKategori` | DELETE | `/api/kategori/:id` | KF-02 |

### 6.3 Produk Handler (`produk_handler.go`)

| Fungsi | Method | Path | Terkait KF |
|---|---|---|---|
| `GetAllProduk` | GET | `/api/produk` | KF-03 |
| `GetProdukByID` | GET | `/api/produk/:id` | KF-03 |
| `CreateProduk` | POST | `/api/produk` | KF-03 |
| `UpdateProduk` | PUT | `/api/produk/:id` | KF-03 |
| `DeleteProduk` | DELETE | `/api/produk/:id` | KF-03 |

### 6.4 Stok Masuk Handler (`stok_masuk_handler.go`)

| Fungsi | Method | Path | Terkait KF |
|---|---|---|---|
| `GetAllStokMasuk` | GET | `/api/stok-masuk` | KF-04 |
| `CreateStokMasuk` | POST | `/api/stok-masuk` | KF-04 |

### 6.5 Stok Keluar Handler (`stok_keluar_handler.go`)

| Fungsi | Method | Path | Terkait KF |
|---|---|---|---|
| `GetAllStokKeluar` | GET | `/api/stok-keluar` | KF-05 |
| `PreviewBatchFEFO` | GET | `/api/stok-keluar/preview-batch` | KF-06 |
| `CreateStokKeluar` | POST | `/api/stok-keluar` | KF-05, KF-06, KF-07 |

### 6.6 Monitoring Handler (`monitoring_handler.go`)

| Fungsi | Method | Path | Terkait KF |
|---|---|---|---|
| `GetMonitoringStok` | GET | `/api/monitoring` | KF-08 |
| `GetDetailProduk` | GET | `/api/monitoring/:idProduk` | KF-08 |

### 6.7 Opname Handler (`opname_handler.go`)

| Fungsi | Method | Path | Terkait KF |
|---|---|---|---|
| `GetAllOpname` | GET | `/api/opname` | KF-09 |
| `GetDetailOpname` | GET | `/api/opname/:id` | KF-09 |
| `MulaiOpname` | POST | `/api/opname` | KF-09 |
| `SelesaiOpname` | PUT | `/api/opname/:id/selesai` | KF-09 |

### 6.8 Laporan Handler (`laporan_handler.go`)

| Fungsi | Method | Path | Terkait KF |
|---|---|---|---|
| `GetLaporanStokMasuk` | GET | `/api/laporan/stok-masuk` | KF-10 |
| `GetLaporanStokKeluar` | GET | `/api/laporan/stok-keluar` | KF-10 |
| `GetLaporanSisaStok` | GET | `/api/laporan/sisa-stok` | KF-10 |

---

## 7. Modul Service (Logika Bisnis)

### 7.1 Auth Service (`auth_service.go`)

**`Login(username, password string) (token string, isDefaultPassword bool, err error)`**
1. Ambil user dari repository berdasarkan `username`.
2. Verifikasi `password` terhadap hash menggunakan bcrypt (`util/hash.go`).
3. Jika valid: generate JWT token (`util/jwt.go`), kembalikan token dan flag `isDefaultPassword`.
4. Jika tidak valid: kembalikan error `ErrKredensialTidakValid`.

**`GantiPassword(idUser int, passwordBaru string) error`**
1. Hash `passwordBaru` menggunakan bcrypt.
2. Update `password` dan set `isDefaultPassword = false` di tabel `users`.

### 7.2 Kategori Service (`kategori_service.go`)

**`TambahKategori(namaKategori string) error`**
1. Periksa duplikasi `namaKategori` (case-insensitive) via repository.
2. Jika duplikat: kembalikan error `ErrDuplikasiKategori`.
3. Simpan kategori baru; kode kategori di-generate otomatis (format: `KAT-xxx`).

**`UpdateKategori(idKategori int, namaKategori string) error`**
1. Periksa duplikasi — kecuali terhadap `idKategori` yang sedang diubah.
2. Simpan perubahan.

**`HapusKategori(idKategori int) error`**
1. Hitung produk terkait via repository.
2. Jika jumlah > 0: kembalikan error `ErrKategoriMemilikiProduk`.
3. Hapus kategori.

### 7.3 Produk Service (`produk_service.go`)

**`TambahProduk(data ProdukInput) error`**
1. Validasi data (`isiPerKemasan > 0`, `polaPenggunaan` valid).
2. Generate kode produk otomatis (format: `PRD-xxxxx`).
3. Simpan produk.

**`HapusProduk(idProduk int) error`**
1. Periksa stok aktif (`batch_stok` dengan `statusBatch = AKTIF` dan `stokKemasan > 0`).
2. Periksa riwayat transaksi (`stok_masuk` atau `stok_keluar` terkait).
3. Jika salah satu terpenuhi: kembalikan error yang sesuai.
4. Hapus produk.

### 7.4 Stok Masuk Service (`stok_masuk_service.go`)

**`SimpanStokMasuk(data StokMasukInput, idUser int) error`**

Validasi backend (KNF-04):
- `jumlahKemasan > 0`
- `expiredDate > tanggalPenerimaan`
- `tanggalPenerimaan ≤ hari ini`
- Produk harus ada di data master

Logika batch (sesuai catatan class diagram):
1. Cari batch existing untuk `idProduk` + `expiredDate` yang sama.
2. Jika batch ditemukan: tambahkan `jumlahKemasan` ke `stokKemasan` batch tersebut dan perbarui `totalIsiTersedia`.
3. Jika tidak ditemukan: buat batch baru, generate `kodeBatch` otomatis (format: `BCH-{idProduk}-{yyyyMMdd}-{seq}`), set `statusBatch = AKTIF`.
4. Simpan record `stok_masuk`.
5. Seluruh langkah 1–4 dibungkus dalam **DB transaction**.

### 7.5 Stok Keluar Service (`stok_keluar_service.go`)

**`GetBatchPrioritasFEFO(idProduk int) (BatchPreview, error)`**
1. Query batch berstatus `AKTIF` milik produk, urutkan `expiredDate ASC`.
2. Ambil batch dengan `expiredDate` paling dekat.
3. Jika produk `PARTIAL_USE`: sertakan info `kemasanTerbuka` aktif pada batch tersebut (jika ada).

**`SimpanStokKeluar(data StokKeluarInput, idUser int) error`**

Seluruh proses dibungkus dalam **DB transaction** (KNF-04).

*Alur Full Use:*
1. Dapatkan batch prioritas FEFO.
2. Validasi `stokKemasan >= jumlahKemasanDipakai`.
3. Kurangi `stokKemasan` pada batch; perbarui `totalIsiTersedia` produk.
4. Update `statusBatch` jika `stokKemasan = 0` → `HABIS`.
5. Simpan record `stok_keluar`.

*Alur Partial Use:*
1. Dapatkan batch prioritas FEFO.
2. Periksa kemasan terbuka aktif pada batch tersebut.
3. Jika **tidak ada kemasan terbuka aktif**:
   - Buka kemasan baru: kurangi `stokKemasan` batch sebesar 1.
   - Hitung `isiAwal = isiPerKemasan` dari data master produk.
   - Tetapkan BUD: `bud = tanggalPenggunaan + 28 hari` (KF-07).
   - Buat record `kemasan_terbuka` baru dengan `statusBUD = AKTIF`.
   - Kurangi `isiTersisa` sebesar `jumlahIsiDipakai`.
4. Jika **ada kemasan terbuka aktif**:
   - Periksa `statusBUD`: jika `KADALUWARSA`, nonaktifkan kemasan terbuka tersebut dan lanjutkan seperti poin 3 (buka kemasan baru).
   - Jika `statusBUD = AKTIF`: kurangi `isiTersisa` sebesar `jumlahIsiDipakai`.
   - Jika `isiTersisa = 0` setelah pengurangan: set `statusBUD = KADALUWARSA` (kemasan habis).
5. Perbarui `totalIsiTersedia` produk.
6. Simpan record `stok_keluar`.

### 7.6 Monitoring Service (`monitoring_service.go`)

**`GetMonitoringStok(filter MonitoringFilter) ([]MonitoringProduk, error)`**
1. Query produk dengan join ke `batch_stok` dan `kemasan_terbuka`.
2. Hitung agregat: total `stokKemasan`, total `totalIsiTersedia` per produk.
3. Tentukan indikator status expired date per batch (AMAN / MENDEKATI / KADALUWARSA) berdasarkan selisih hari dari tanggal hari ini ke `expiredDate`.
4. Terapkan filter: kategori, status expired, status BUD, nama produk.

### 7.7 Opname Service (`opname_service.go`)

**`MulaiOpname(idUser int) (StokOpnameDetail, error)`**
1. Buat sesi `stok_opname` baru.
2. Kumpulkan seluruh data yang perlu diperiksa: semua `batch_stok` aktif per produk + semua `kemasan_terbuka` aktif untuk produk PARTIAL_USE.
3. Kembalikan daftar item opname.

**`SelesaiOpname(idOpname int, items []ItemOpname) error`**

Seluruh proses dibungkus dalam **DB transaction**:
1. Hitung `selisih = stokFisik - stokSistem` per item.
2. Jika ada selisih: validasi `keterangan` tidak kosong (KF-09, BR-09.2).
3. Lakukan penyesuaian stok: update `stokKemasan` pada `batch_stok` atau `isiTersisa` pada `kemasan_terbuka` sesuai `stokFisik`.
4. Simpan baris `detail_opname` untuk setiap item.
5. Update `statusOpname = SELESAI`.

### 7.8 Worker Service (`worker_service.go`)

Lihat § 9 (Background Worker) untuk detail lengkap.

### 7.9 Laporan Service (`laporan_service.go`)

**`GetLaporanStokMasuk(params LaporanParams) ([]LaporanStokMasukRow, error)`**
- Query `stok_masuk` JOIN `produk` JOIN `batch_stok` dengan filter `tanggalPenerimaan` antara `periodeAwal` dan `periodeAkhir`.
- Filter kategori bersifat opsional.

**`GetLaporanStokKeluar(params LaporanParams) ([]LaporanStokKeluarRow, error)`**
- Query `stok_keluar` JOIN `produk` JOIN `batch_stok` dengan filter `tanggalPenggunaan` antara `periodeAwal` dan `periodeAkhir`.

**`GetLaporanSisaStok(params LaporanParams) ([]LaporanSisaStokRow, error)`**
- Query posisi stok terkini per produk dan batch pada akhir `periodeAkhir`.

---

## 8. Repository Layer

Repository hanya bertanggung jawab atas operasi baca/tulis ke database. Tidak ada logika bisnis di lapisan ini. Setiap fungsi repository menerima opsional parameter `*sql.Tx` untuk mendukung DB transaction yang dikelola oleh service.

### 8.1 Contoh Pola Repository

```go
// batch_repository.go

// Mencari batch existing berdasarkan produk dan expired date
func (r *BatchRepository) FindByProdukAndExpiredDate(
    tx *sql.Tx, idProduk int, expiredDate time.Time,
) (*BatchStok, error)

// Membuat batch baru
func (r *BatchRepository) Create(tx *sql.Tx, batch *BatchStok) (int, error)

// Menambah stok ke batch existing
func (r *BatchRepository) TambahStok(tx *sql.Tx, idBatch int, jumlahKemasan int, totalIsiTambahan float64) error

// Mendapatkan batch AKTIF dengan FEFO (urut expiredDate ASC)
func (r *BatchRepository) GetBatchAktifFEFO(tx *sql.Tx, idProduk int) ([]BatchStok, error)

// Memperbarui status batch
func (r *BatchRepository) UpdateStatus(tx *sql.Tx, idBatch int, status string) error

// Memperbarui stok kemasan batch
func (r *BatchRepository) UpdateStokKemasan(tx *sql.Tx, idBatch int, jumlah int) error
```

### 8.2 Tabel Repository dan Fungsi Utama

| Repository | Tabel Utama | Fungsi Kunci |
|---|---|---|
| `UserRepository` | `users` | `FindByUsername`, `UpdatePassword` |
| `KategoriRepository` | `kategori` | `FindAll`, `FindByNama`, `Create`, `Update`, `Delete`, `CountProdukTerkait` |
| `ProdukRepository` | `produk` | `FindAll`, `FindByID`, `Create`, `Update`, `Delete`, `HasStokAktif`, `HasRiwayatTransaksi` |
| `StokMasukRepository` | `stok_masuk` | `Create`, `FindAll` |
| `BatchRepository` | `batch_stok` | `FindByProdukAndExpiredDate`, `Create`, `TambahStok`, `GetBatchAktifFEFO`, `UpdateStatus`, `UpdateStokKemasan` |
| `StokKeluarRepository` | `stok_keluar` | `Create`, `FindAll` |
| `KemasanTerbukaRepository` | `kemasan_terbuka` | `FindAktifByBatch`, `Create`, `UpdateIsiTersisa`, `Nonaktifkan` |
| `OpnameRepository` | `stok_opname`, `detail_opname` | `CreateSesi`, `CreateDetail`, `UpdateStatus`, `FindAll`, `FindDetailByOpname` |
| `LaporanRepository` | (multi-tabel) | `GetStokMasukByPeriode`, `GetStokKeluarByPeriode`, `GetSisaStokByPeriode` |

---

## 9. Background Worker

### 9.1 Deskripsi

Background worker adalah goroutine yang berjalan secara periodik dan independen dari alur request-response HTTP, sesuai KNF-07 (Availability Background Worker). Worker diinisialisasi saat startup aplikasi di `main.go` dan tidak memerlukan trigger manual atau cron eksternal.

### 9.2 Worker 1 — Pembaruan Status Batch

**File:** `service/worker_service.go` — fungsi `StartBatchStatusWorker()`

**Tugas:** Memeriksa seluruh batch berstatus `AKTIF` atau `HABIS` yang `expiredDate`-nya telah terlewati (`expiredDate < NOW()`), lalu memperbarui `statusBatch` menjadi `KADALUWARSA`.

**Alur:**

```
Goroutine dimulai saat aplikasi start
     │
     └─► Loop tak terbatas:
              │
              ▼
         Jalankan query:
         UPDATE batch_stok
         SET status_batch = 'KADALUWARSA'
         WHERE expired_date < NOW()
           AND status_batch IN ('AKTIF', 'HABIS')
              │
              ├── Sukses → log jumlah baris yang diperbarui
              ├── Error  → log error, TIDAK panic, TIDAK stop goroutine
              │
              ▼
         time.Sleep(interval)   ← default: 1 jam (dikonfigurasi via env WORKER_BATCH_INTERVAL)
```

**Acceptance Criteria:** AC-NFR07.1 — Status batch yang melewati `expiredDate` berubah menjadi `KADALUWARSA` dalam satu siklus worker tanpa interaksi user.

### 9.3 Worker 2 — Pembaruan Status BUD

**File:** `service/worker_service.go` — fungsi `StartBUDStatusWorker()`

**Tugas:** Memeriksa seluruh `kemasan_terbuka` berstatus `AKTIF` yang `bud`-nya telah terlewati (`bud < NOW()`), lalu memperbarui `statusBUD` menjadi `KADALUWARSA`.

**Alur:**

```
Goroutine dimulai saat aplikasi start
     │
     └─► Loop tak terbatas:
              │
              ▼
         Jalankan query:
         UPDATE kemasan_terbuka
         SET status_bud = 'KADALUWARSA'
         WHERE bud < NOW()
           AND status_bud = 'AKTIF'
              │
              ├── Sukses → log jumlah baris yang diperbarui
              ├── Error  → log error, TIDAK panic, TIDAK stop goroutine
              │
              ▼
         time.Sleep(interval)   ← default: 1 jam (dikonfigurasi via env WORKER_BUD_INTERVAL)
```

**Acceptance Criteria:** AC-NFR07.2 — Status `KemasanTerbuka` yang BUD-nya lewat berubah menjadi `KADALUWARSA` dalam satu siklus worker tanpa interaksi user.

### 9.4 Fault Tolerance

- Worker **tidak** menggunakan `panic`/`recover` untuk menghentikan proses — setiap error siklus di-log dan worker melanjutkan ke siklus berikutnya.
- Restart aplikasi backend otomatis menjalankan ulang worker (goroutine di-spawn di `main.go`), sesuai AC-NFR07.3.
- Worker berjalan di goroutine terpisah — kegagalan worker tidak memengaruhi penanganan request HTTP.

### 9.5 Inisialisasi di `main.go`

```go
func main() {
    // ... inisialisasi DB, config, router ...

    // Jalankan background worker sebagai goroutine terpisah
    go workerService.StartBatchStatusWorker()
    go workerService.StartBUDStatusWorker()

    // Mulai HTTP server (blocking)
    app.Listen(":" + config.Port)
}
```

---

## 10. Penanganan Error & Format Respons

### 10.1 Format Respons Sukses

```json
{
  "success": true,
  "message": "Pesan sukses sesuai KF",
  "data": { ... }
}
```

Untuk respons tanpa data (aksi delete, update):
```json
{
  "success": true,
  "message": "Kategori berhasil dihapus."
}
```

### 10.2 Format Respons Error

```json
{
  "success": false,
  "message": "Pesan error sesuai KF",
  "errors": [ ... ]   // opsional, untuk validasi field
}
```

### 10.3 HTTP Status Code

| Kondisi | Status Code |
|---|---|
| Operasi berhasil (baca/buat/ubah/hapus) | `200 OK` |
| Resource baru berhasil dibuat | `201 Created` |
| Request tidak valid (validasi gagal) | `400 Bad Request` |
| Token tidak ada / tidak valid / expired | `401 Unauthorized` |
| Resource tidak ditemukan | `404 Not Found` |
| Konflik (duplikasi, constraint) | `409 Conflict` |
| Error server internal | `500 Internal Server Error` |

### 10.4 Error Sentinel

Definisi error bisnis di `service/` menggunakan error sentinel untuk memudahkan handler memetakan ke status code:

| Error Sentinel | HTTP Code | Pesan |
|---|---|---|
| `ErrKredensialTidakValid` | 401 | `"Kredensial tidak valid."` |
| `ErrDuplikasiKategori` | 409 | `"Nama kategori sudah terdaftar dalam sistem."` |
| `ErrKategoriMemilikiProduk` | 409 | `"Kategori tidak dapat dihapus karena masih memiliki produk terkait."` |
| `ErrProdukMemilikiStokAktif` | 409 | `"Produk tidak dapat dihapus karena masih memiliki stok aktif."` |
| `ErrProdukMemilikiRiwayat` | 409 | `"Produk tidak dapat dihapus karena memiliki riwayat transaksi."` |
| `ErrStokTidakCukup` | 400 | `"Stok tidak mencukupi untuk transaksi ini."` |
| `ErrParameterLaporanTidakValid` | 400 | `"Tanggal akhir tidak boleh lebih kecil dari tanggal awal."` |

---

## 11. Konvensi Kode

### 11.1 Penamaan

- Nama fungsi: `CamelCase` (Go idiom).
- Nama tabel/kolom: `snake_case` (sesuai `srs-database.md`).
- Nama struct model: sesuai nama kelas pada class diagram (§ 4.3.4 PDF).

### 11.2 DB Transaction

Semua operasi yang mengubah lebih dari satu tabel **wajib** menggunakan `sql.Tx`, sesuai KNF-04:
- `SimpanStokMasuk` → mengubah `stok_masuk` + `batch_stok`.
- `SimpanStokKeluar` → mengubah `stok_keluar` + `batch_stok` + `kemasan_terbuka`.
- `SelesaiOpname` → mengubah `stok_opname` + `detail_opname` + `batch_stok`/`kemasan_terbuka`.

### 11.3 Logging

- Gunakan structured logging (misal: `log/slog` Go standard library).
- Level `INFO` untuk siklus worker sukses; level `ERROR` untuk kegagalan DB.
- Password dan token **tidak pernah** masuk ke log (KNF-02, AC-NFR02.2).

### 11.4 Konfigurasi Environment

| Variabel | Deskripsi |
|---|---|
| `DATABASE_URL` | DSN PostgreSQL |
| `JWT_SECRET` | Secret key JWT |
| `JWT_EXPIRY` | Durasi token (misal: `8h`) |
| `PORT` | Port HTTP server (default: `8080`) |
| `WORKER_BATCH_INTERVAL` | Interval worker batch (default: `1h`) |
| `WORKER_BUD_INTERVAL` | Interval worker BUD (default: `1h`) |
