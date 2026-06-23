# SRS — REST API Specification
## Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic

---

## Daftar Isi

1. [Tujuan & Cakupan](#1-tujuan--cakupan)
2. [Konvensi Umum API](#2-konvensi-umum-api)
3. [Format Respons Standar](#3-format-respons-standar)
4. [Autentikasi & Middleware](#4-autentikasi--middleware)
5. [API-01 — Auth](#5-api-01--auth)
6. [API-02 — Kategori](#6-api-02--kategori)
7. [API-03 — Produk](#7-api-03--produk)
8. [API-04 — Stok Masuk](#8-api-04--stok-masuk)
9. [API-05 — Stok Keluar](#9-api-05--stok-keluar)
10. [API-06 — Monitoring Stok](#10-api-06--monitoring-stok)
11. [API-07 — Stock Opname](#11-api-07--stock-opname)
12. [API-08 — Laporan Stok](#12-api-08--laporan-stok)
13. [Kode Error & Pesan Standar](#13-kode-error--pesan-standar)
14. [Matriks Traceability API](#14-matriks-traceability-api)

---

## 1. Tujuan & Cakupan

Dokumen ini mendefinisikan kontrak REST API lengkap antara lapisan antarmuka (frontend React) dan lapisan logika bisnis (backend Go + Fiber) Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic. Setiap endpoint dirinci mencakup method HTTP, path, parameter request, struktur respons sukses, kemungkinan respons error, serta keterkaitan dengan kebutuhan fungsional (KF) yang telah ditetapkan pada `srs-fr.md`.

API ini merupakan satu-satunya jalur komunikasi antara frontend dan backend — tidak ada akses langsung dari frontend ke basis data. Seluruh kalkulasi stok, penerapan FEFO, dan pengelolaan BUD dikembalikan dalam respons API sebagai nilai siap tampil, sesuai KNF-04 (Keakuratan dan Integritas Data).

---

## 2. Konvensi Umum API

### 2.1 Base URL

```
/api
```

Seluruh endpoint diawali dengan prefix `/api`. Dalam lingkungan produksi, server backend berjalan di balik reverse proxy yang menyajikan frontend dan API pada domain yang sama untuk menghindari konfigurasi CORS yang kompleks.

### 2.2 Format Data

- Seluruh request dan respons menggunakan format **JSON** dengan header `Content-Type: application/json`.
- Nilai tanggal menggunakan format **ISO 8601**: `YYYY-MM-DD` untuk tipe `DATE`, dan `YYYY-MM-DDTHH:mm:ssZ` untuk tipe `TIMESTAMP`.
- Nilai desimal (jumlah isi, isi tersisa) menggunakan tipe `number` (float) pada JSON.

### 2.3 HTTP Method

| Method | Semantik |
|---|---|
| `GET` | Membaca data (tidak mengubah state) |
| `POST` | Membuat resource baru |
| `PUT` | Memperbarui resource yang ada (replace atau partial) |
| `DELETE` | Menghapus resource |

### 2.4 Autentikasi

Seluruh endpoint kecuali `POST /api/auth/login` **wajib** menyertakan header:

```
Authorization: Bearer <token_jwt>
```

Token diperoleh dari respons login dan disimpan di `authStore` (sisi frontend). Detail validasi token dijabarkan pada `srs-backend.md` § Middleware Autentikasi.

### 2.5 Paginasi (Query Parameter Umum)

Endpoint yang mengembalikan daftar mendukung parameter paginasi opsional:

| Parameter | Tipe | Default | Keterangan |
|---|---|---|---|
| `page` | `integer` | `1` | Halaman yang diminta |
| `limit` | `integer` | `20` | Jumlah item per halaman |

---

## 3. Format Respons Standar

### 3.1 Respons Sukses (dengan Data)

```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": { ... }
}
```

Untuk respons daftar:

```json
{
  "success": true,
  "message": "Data berhasil diambil.",
  "data": {
    "items": [ ... ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 85,
      "totalPages": 5
    }
  }
}
```

### 3.2 Respons Sukses (Tanpa Data — Aksi Mutasi)

```json
{
  "success": true,
  "message": "Kategori berhasil dihapus."
}
```

### 3.3 Respons Error

```json
{
  "success": false,
  "message": "Pesan error yang informatif.",
  "errors": [
    { "field": "expiredDate", "message": "Expired date harus lebih besar dari tanggal penerimaan." }
  ]
}
```

Field `errors` bersifat opsional; muncul hanya pada kasus validasi multi-field (HTTP 400).

### 3.4 HTTP Status Code yang Digunakan

| Kondisi | Status Code |
|---|---|
| Operasi baca/ubah/hapus berhasil | `200 OK` |
| Resource baru berhasil dibuat | `201 Created` |
| Request tidak valid (validasi input gagal) | `400 Bad Request` |
| Token tidak ada, tidak valid, atau kadaluwarsa | `401 Unauthorized` |
| Resource tidak ditemukan | `404 Not Found` |
| Konflik: duplikasi data atau constraint bisnis dilanggar | `409 Conflict` |
| Error internal server | `500 Internal Server Error` |

---

## 4. Autentikasi & Middleware

### 4.1 Alur Token

1. Frontend memanggil `POST /api/auth/login` dengan kredensial.
2. Backend memverifikasi kredensial dan mengembalikan `token` JWT.
3. Frontend menyimpan token di `authStore` (Zustand) dan menyisipkannya pada setiap request berikutnya melalui interceptor Axios (`Authorization: Bearer <token>`).
4. Jika backend mengembalikan HTTP `401`, frontend menghapus token dari store dan me-redirect pengguna ke `/login`.

### 4.2 Payload JWT

```json
{
  "sub": 1,
  "iat": 1700000000,
  "exp": 1700028800
}
```

- `sub`: `id_user` dari tabel `users`.
- `iat`: Waktu token diterbitkan (Unix timestamp).
- `exp`: Waktu token kedaluwarsa (Unix timestamp); masa berlaku dikonfigurasi via `JWT_EXPIRY` (rekomendasi: `8h` sesuai jam operasional klinik).

---

## 5. API-01 — Auth

### 5.1 Login

**Terkait KF:** KF-01

```
POST /api/auth/login
```

**Request Body:**

```json
{
  "username": "admin_farmasi",
  "password": "password123"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `username` | `string` | ✓ | Nama pengguna |
| `password` | `string` | ✓ | Password (plain text; akan diverifikasi terhadap hash di database) |

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Login berhasil.",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "isDefaultPassword": false
  }
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `token` | `string` | JWT token untuk sesi aktif |
| `isDefaultPassword` | `boolean` | `true` jika pengguna belum mengganti password default; frontend wajib redirect ke `/ganti-password` |

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `401` | Kredensial tidak valid | `"Kredensial tidak valid."` |
| `400` | Field wajib tidak diisi | `"Username dan password wajib diisi."` |

---

### 5.2 Ganti Password

**Terkait KF:** KF-01

```
PUT /api/auth/password
```

> Endpoint ini diproteksi middleware JWT. `idUser` diambil dari payload token, bukan dari request body.

**Request Body:**

```json
{
  "passwordBaru": "passwordBaru123"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `passwordBaru` | `string` | ✓ | Password baru yang akan disimpan (min. 8 karakter direkomendasikan) |

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Password berhasil diperbarui."
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `400` | Field `passwordBaru` kosong | `"Password baru wajib diisi."` |
| `401` | Token tidak valid | `"Akses tidak diizinkan."` |

---

## 6. API-02 — Kategori

**Terkait KF:** KF-02

### 6.1 Ambil Daftar Kategori

```
GET /api/kategori
```

**Query Parameters (opsional):**

| Parameter | Tipe | Keterangan |
|---|---|---|
| `search` | `string` | Filter berdasarkan nama kategori (case-insensitive) |
| `page` | `integer` | Halaman (default: 1) |
| `limit` | `integer` | Jumlah per halaman (default: 20) |

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Data kategori berhasil diambil.",
  "data": {
    "items": [
      {
        "idKategori": 1,
        "kodeKategori": "KAT-001",
        "namaKategori": "Skincare",
        "jumlahProdukTerkait": 12
      },
      {
        "idKategori": 2,
        "kodeKategori": "KAT-002",
        "namaKategori": "Injectable",
        "jumlahProdukTerkait": 5
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 5,
      "totalPages": 1
    }
  }
}
```

---

### 6.2 Tambah Kategori

```
POST /api/kategori
```

**Request Body:**

```json
{
  "namaKategori": "Obat"
}
```

**Respons Sukses — `201 Created`:**

```json
{
  "success": true,
  "message": "Kategori berhasil ditambahkan.",
  "data": {
    "idKategori": 3,
    "kodeKategori": "KAT-003",
    "namaKategori": "Obat"
  }
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `400` | Field `namaKategori` kosong | `"Nama kategori wajib diisi."` |
| `409` | Nama kategori sudah ada | `"Nama kategori sudah terdaftar dalam sistem."` |

---

### 6.3 Ubah Kategori

```
PUT /api/kategori/:id
```

**Path Parameter:** `id` — `idKategori` yang akan diubah.

**Request Body:**

```json
{
  "namaKategori": "Skincare & Body Care"
}
```

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Kategori berhasil diperbarui.",
  "data": {
    "idKategori": 1,
    "kodeKategori": "KAT-001",
    "namaKategori": "Skincare & Body Care"
  }
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `400` | Field kosong | `"Nama kategori wajib diisi."` |
| `404` | Kategori tidak ditemukan | `"Kategori tidak ditemukan."` |
| `409` | Nama sudah digunakan kategori lain | `"Nama kategori sudah terdaftar dalam sistem."` |

---

### 6.4 Hapus Kategori

```
DELETE /api/kategori/:id
```

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Kategori berhasil dihapus."
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `404` | Kategori tidak ditemukan | `"Kategori tidak ditemukan."` |
| `409` | Masih memiliki produk terkait | `"Kategori tidak dapat dihapus karena masih memiliki produk terkait."` |

---

## 7. API-03 — Produk

**Terkait KF:** KF-03

### 7.1 Ambil Daftar Produk

```
GET /api/produk
```

**Query Parameters (opsional):**

| Parameter | Tipe | Keterangan |
|---|---|---|
| `idKategori` | `integer` | Filter berdasarkan kategori |
| `polaPenggunaan` | `string` | `FULL_USE` atau `PARTIAL_USE` |
| `search` | `string` | Cari berdasarkan nama produk |
| `page` | `integer` | Halaman (default: 1) |
| `limit` | `integer` | Jumlah per halaman (default: 20) |

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Data produk berhasil diambil.",
  "data": {
    "items": [
      {
        "idProduk": 1,
        "kodeProduk": "PRD-00001",
        "namaProduk": "Botox Bionex 100 IU",
        "kategori": {
          "idKategori": 2,
          "namaKategori": "Injectable"
        },
        "bentukKemasan": "vial",
        "satuanIsi": "IU",
        "isiPerKemasan": 100.00,
        "polaPenggunaan": "PARTIAL_USE",
        "stokKemasan": 5,
        "totalIsiTersedia": 480.00
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 42,
      "totalPages": 3
    }
  }
}
```

---

### 7.2 Ambil Detail Produk

```
GET /api/produk/:id
```

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Data produk berhasil diambil.",
  "data": {
    "idProduk": 1,
    "kodeProduk": "PRD-00001",
    "namaProduk": "Botox Bionex 100 IU",
    "kategori": {
      "idKategori": 2,
      "namaKategori": "Injectable"
    },
    "bentukKemasan": "vial",
    "satuanIsi": "IU",
    "isiPerKemasan": 100.00,
    "polaPenggunaan": "PARTIAL_USE",
    "stokKemasan": 5,
    "totalIsiTersedia": 480.00,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-11-20T08:15:00Z"
  }
}
```

---

### 7.3 Tambah Produk

```
POST /api/produk
```

**Request Body:**

```json
{
  "namaProduk": "Botox Bionex 100 IU",
  "idKategori": 2,
  "bentukKemasan": "vial",
  "satuanIsi": "IU",
  "isiPerKemasan": 100.00,
  "polaPenggunaan": "PARTIAL_USE"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `namaProduk` | `string` | ✓ | Nama produk farmasi |
| `idKategori` | `integer` | ✓ | ID kategori produk (harus ada di tabel `kategori`) |
| `bentukKemasan` | `string` | ✓ | Bentuk kemasan (vial, jar, botol, ampoule, syringe, dll.) |
| `satuanIsi` | `string` | ✓ | Satuan isi (pcs, ml, cc, gram, IU, dll.) |
| `isiPerKemasan` | `number` | ✓ | Kapasitas isi per kemasan; harus > 0 |
| `polaPenggunaan` | `string` | ✓ | `FULL_USE` atau `PARTIAL_USE` |

**Respons Sukses — `201 Created`:**

```json
{
  "success": true,
  "message": "Produk berhasil ditambahkan.",
  "data": {
    "idProduk": 1,
    "kodeProduk": "PRD-00001",
    "namaProduk": "Botox Bionex 100 IU",
    "idKategori": 2,
    "bentukKemasan": "vial",
    "satuanIsi": "IU",
    "isiPerKemasan": 100.00,
    "polaPenggunaan": "PARTIAL_USE"
  }
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `400` | Field wajib tidak diisi atau `isiPerKemasan` ≤ 0 | Pesan validasi per field (lihat § 3.3) |
| `400` | Nilai `polaPenggunaan` tidak valid | `"Pola penggunaan harus FULL_USE atau PARTIAL_USE."` |
| `404` | `idKategori` tidak ditemukan | `"Kategori tidak ditemukan."` |

---

### 7.4 Ubah Produk

```
PUT /api/produk/:id
```

**Request Body:** Sama dengan Tambah Produk (§ 7.3).

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Produk berhasil diperbarui.",
  "data": { ... }
}
```

**Respons Error:** Sama dengan Tambah Produk, ditambah:

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `404` | Produk tidak ditemukan | `"Produk tidak ditemukan."` |

---

### 7.5 Hapus Produk

```
DELETE /api/produk/:id
```

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Data produk berhasil dihapus."
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `404` | Produk tidak ditemukan | `"Produk tidak ditemukan."` |
| `409` | Produk masih memiliki stok aktif | `"Produk tidak dapat dihapus karena masih memiliki stok aktif."` |
| `409` | Produk memiliki riwayat transaksi | `"Produk tidak dapat dihapus karena memiliki riwayat transaksi."` |

---

## 8. API-04 — Stok Masuk

**Terkait KF:** KF-04

### 8.1 Ambil Daftar Stok Masuk

```
GET /api/stok-masuk
```

**Query Parameters (opsional):**

| Parameter | Tipe | Keterangan |
|---|---|---|
| `idProduk` | `integer` | Filter per produk |
| `periodeAwal` | `string` (DATE) | Tanggal penerimaan awal (format: `YYYY-MM-DD`) |
| `periodeAkhir` | `string` (DATE) | Tanggal penerimaan akhir (format: `YYYY-MM-DD`) |
| `page` | `integer` | Halaman (default: 1) |
| `limit` | `integer` | Jumlah per halaman (default: 20) |

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Data stok masuk berhasil diambil.",
  "data": {
    "items": [
      {
        "idStokMasuk": 1,
        "tanggalPenerimaan": "2024-11-01",
        "produk": {
          "idProduk": 1,
          "kodeProduk": "PRD-00001",
          "namaProduk": "Botox Bionex 100 IU",
          "bentukKemasan": "vial",
          "satuanIsi": "IU",
          "isiPerKemasan": 100.00,
          "polaPenggunaan": "PARTIAL_USE"
        },
        "batch": {
          "idBatch": 3,
          "kodeBatch": "BCH-1-20241101-001",
          "expiredDate": "2026-06-30"
        },
        "jumlahKemasan": 5,
        "totalIsiMasuk": 500.00,
        "keterangan": ""
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 30,
      "totalPages": 2
    }
  }
}
```

---

### 8.2 Catat Stok Masuk

```
POST /api/stok-masuk
```

**Request Body:**

```json
{
  "tanggalPenerimaan": "2024-11-01",
  "idProduk": 1,
  "expiredDate": "2026-06-30",
  "jumlahKemasan": 5,
  "keterangan": ""
}
```

| Field | Tipe | Wajib | Validasi Backend |
|---|---|---|---|
| `tanggalPenerimaan` | `string` (DATE) | ✓ | Tidak boleh melebihi tanggal hari ini |
| `idProduk` | `integer` | ✓ | Harus ada di tabel `produk` |
| `expiredDate` | `string` (DATE) | ✓ | Harus lebih besar dari `tanggalPenerimaan` |
| `jumlahKemasan` | `integer` | ✓ | Harus > 0 |
| `keterangan` | `string` | — | Opsional |

**Logika backend:** Backend menentukan apakah akan membuat batch baru atau menambahkan stok ke batch yang sudah ada berdasarkan kombinasi `idProduk` + `expiredDate` (lihat `srs-fr.md` KF-04 BR-04.1, `srs-backend.md` § 7.4).

**Respons Sukses — `201 Created`:**

```json
{
  "success": true,
  "message": "Data stok masuk berhasil disimpan.",
  "data": {
    "idStokMasuk": 10,
    "tanggalPenerimaan": "2024-11-01",
    "produk": {
      "idProduk": 1,
      "namaProduk": "Botox Bionex 100 IU"
    },
    "batch": {
      "idBatch": 3,
      "kodeBatch": "BCH-1-20241101-001",
      "expiredDate": "2026-06-30",
      "statusBatch": "AKTIF",
      "stokKemasanSetelah": 10,
      "totalIsiTersediaSetelah": 980.00
    },
    "jumlahKemasan": 5,
    "totalIsiMasuk": 500.00
  }
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `400` | Field wajib tidak diisi | Pesan validasi per field |
| `400` | `tanggalPenerimaan` > hari ini | `"Tanggal penerimaan tidak boleh melebihi tanggal hari ini."` |
| `400` | `expiredDate` ≤ `tanggalPenerimaan` | `"Tanggal kedaluwarsa harus lebih besar dari tanggal penerimaan."` |
| `400` | `jumlahKemasan` ≤ 0 | `"Jumlah kemasan harus lebih dari 0."` |
| `404` | Produk tidak ditemukan | `"Produk tidak ditemukan."` |

---

## 9. API-05 — Stok Keluar

**Terkait KF:** KF-05, KF-06, KF-07

### 9.1 Preview Batch FEFO

Endpoint ini dipanggil frontend segera setelah user memilih produk pada form stok keluar, untuk menampilkan batch prioritas FEFO secara read-only kepada user sebelum melanjutkan input.

```
GET /api/stok-keluar/preview-batch?idProduk=:idProduk
```

**Query Parameter:**

| Parameter | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `idProduk` | `integer` | ✓ | ID produk yang akan digunakan |

**Respons Sukses — `200 OK`:**

Untuk produk `FULL_USE`:

```json
{
  "success": true,
  "message": "Data preview batch berhasil diambil.",
  "data": {
    "produk": {
      "idProduk": 2,
      "namaProduk": "Sunscreen SPF 50",
      "bentukKemasan": "botol",
      "satuanIsi": "pcs",
      "isiPerKemasan": 1.00,
      "polaPenggunaan": "FULL_USE"
    },
    "batchPrioritas": {
      "idBatch": 5,
      "kodeBatch": "BCH-2-20240901-001",
      "expiredDate": "2025-09-01",
      "stokKemasan": 8,
      "statusBatch": "AKTIF"
    },
    "kemasanTerbuka": null
  }
}
```

Untuk produk `PARTIAL_USE` dengan kemasan terbuka aktif:

```json
{
  "success": true,
  "message": "Data preview batch berhasil diambil.",
  "data": {
    "produk": {
      "idProduk": 1,
      "namaProduk": "Botox Bionex 100 IU",
      "bentukKemasan": "vial",
      "satuanIsi": "IU",
      "isiPerKemasan": 100.00,
      "polaPenggunaan": "PARTIAL_USE"
    },
    "batchPrioritas": {
      "idBatch": 3,
      "kodeBatch": "BCH-1-20241101-001",
      "expiredDate": "2026-06-30",
      "stokKemasan": 4,
      "statusBatch": "AKTIF"
    },
    "kemasanTerbuka": {
      "idKemasanTerbuka": 2,
      "tanggalDibuka": "2024-11-05",
      "bud": "2024-12-03",
      "isiAwal": 100.00,
      "isiTersisa": 50.00,
      "statusBud": "AKTIF"
    }
  }
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `400` | Parameter `idProduk` tidak disertakan | `"Parameter idProduk wajib diisi."` |
| `404` | Produk tidak ditemukan | `"Produk tidak ditemukan."` |
| `404` | Tidak ada batch aktif tersedia | `"Tidak ada stok aktif tersedia untuk produk ini."` |

---

### 9.2 Ambil Daftar Stok Keluar

```
GET /api/stok-keluar
```

**Query Parameters (opsional):**

| Parameter | Tipe | Keterangan |
|---|---|---|
| `idProduk` | `integer` | Filter per produk |
| `periodeAwal` | `string` (DATE) | Tanggal penggunaan awal |
| `periodeAkhir` | `string` (DATE) | Tanggal penggunaan akhir |
| `page` | `integer` | Halaman (default: 1) |
| `limit` | `integer` | Jumlah per halaman (default: 20) |

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Data stok keluar berhasil diambil.",
  "data": {
    "items": [
      {
        "idStokKeluar": 1,
        "tanggalPenggunaan": "2024-11-05",
        "produk": {
          "idProduk": 1,
          "namaProduk": "Botox Bionex 100 IU",
          "polaPenggunaan": "PARTIAL_USE"
        },
        "batch": {
          "idBatch": 3,
          "kodeBatch": "BCH-1-20241101-001",
          "expiredDate": "2026-06-30"
        },
        "jumlahKemasanDipakai": null,
        "jumlahIsiDipakai": 50.00,
        "keterangan": "Pasien: Tindakan Botox Area Dahi"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 55,
      "totalPages": 3
    }
  }
}
```

---

### 9.3 Catat Stok Keluar

```
POST /api/stok-keluar
```

**Request Body untuk `FULL_USE`:**

```json
{
  "tanggalPenggunaan": "2024-11-10",
  "idProduk": 2,
  "jumlahKemasanDipakai": 1,
  "keterangan": "Pasien: Perawatan Facial"
}
```

**Request Body untuk `PARTIAL_USE`:**

```json
{
  "tanggalPenggunaan": "2024-11-10",
  "idProduk": 1,
  "jumlahIsiDipakai": 50.00,
  "keterangan": "Pasien: Tindakan Botox Dahi"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `tanggalPenggunaan` | `string` (DATE) | ✓ | Tanggal produk digunakan |
| `idProduk` | `integer` | ✓ | ID produk yang digunakan |
| `jumlahKemasanDipakai` | `integer` | Kondisional | Diisi untuk produk `FULL_USE`; harus > 0 |
| `jumlahIsiDipakai` | `number` | Kondisional | Diisi untuk produk `PARTIAL_USE`; harus > 0 |
| `keterangan` | `string` | — | Keterangan penggunaan (opsional) |

> Backend menentukan batch yang digunakan secara otomatis menggunakan FEFO, berdasarkan `idProduk`. Field `idBatch` **tidak** dikirim oleh frontend — pemilihan batch sepenuhnya dikelola backend (KF-06, BR-06.2, AC-06.2).

**Respons Sukses — `201 Created`:**

```json
{
  "success": true,
  "message": "Data penggunaan berhasil disimpan.",
  "data": {
    "idStokKeluar": 20,
    "tanggalPenggunaan": "2024-11-10",
    "produk": {
      "idProduk": 1,
      "namaProduk": "Botox Bionex 100 IU"
    },
    "batch": {
      "idBatch": 3,
      "kodeBatch": "BCH-1-20241101-001",
      "expiredDate": "2026-06-30"
    },
    "jumlahIsiDipakai": 50.00,
    "kemasanTerbuka": {
      "idKemasanTerbuka": 2,
      "isiTersisaSetelah": 0.00,
      "statusBud": "KADALUWARSA"
    }
  }
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `400` | Field wajib tidak diisi | Pesan validasi per field |
| `400` | `jumlahKemasanDipakai`/`jumlahIsiDipakai` ≤ 0 | `"Jumlah yang dipakai harus lebih dari 0."` |
| `400` | Stok tidak mencukupi | `"Stok tidak mencukupi untuk transaksi ini."` |
| `404` | Produk tidak ditemukan | `"Produk tidak ditemukan."` |
| `404` | Tidak ada batch aktif | `"Tidak ada stok aktif tersedia untuk produk ini."` |

---

## 10. API-06 — Monitoring Stok

**Terkait KF:** KF-08

### 10.1 Ambil Ringkasan Monitoring Stok

```
GET /api/monitoring
```

Endpoint ini juga digunakan oleh `DashboardPage` untuk menampilkan ringkasan agregat tanpa harus membuat endpoint terpisah (lihat `srs-frontend.md` § 6.3).

**Query Parameters (opsional):**

| Parameter | Tipe | Keterangan |
|---|---|---|
| `idKategori` | `integer` | Filter berdasarkan kategori produk |
| `statusExpired` | `string` | Filter status expired: `AMAN`, `MENDEKATI`, `KADALUWARSA` |
| `statusBud` | `string` | Filter status BUD: `AKTIF`, `KADALUWARSA` |
| `search` | `string` | Cari berdasarkan nama produk |
| `page` | `integer` | Halaman (default: 1) |
| `limit` | `integer` | Jumlah per halaman (default: 20) |

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Data monitoring stok berhasil diambil.",
  "data": {
    "ringkasan": {
      "totalProdukAktif": 42,
      "produkMendekatiExpired": 5,
      "produkKadaluwarsa": 2,
      "kemasanTerbukaAktif": 3
    },
    "items": [
      {
        "idProduk": 1,
        "kodeProduk": "PRD-00001",
        "namaProduk": "Botox Bionex 100 IU",
        "kategori": {
          "idKategori": 2,
          "namaKategori": "Injectable"
        },
        "polaPenggunaan": "PARTIAL_USE",
        "stokKemasan": 4,
        "totalIsiTersedia": 450.00,
        "satuanIsi": "IU",
        "statusExpiredTerdekat": "AMAN",
        "expiredDateTerdekat": "2026-06-30",
        "adaKemasanTerbukaAktif": true
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 42,
      "totalPages": 3
    }
  }
}
```

---

### 10.2 Ambil Detail Produk di Monitoring

```
GET /api/monitoring/:idProduk
```

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Detail monitoring produk berhasil diambil.",
  "data": {
    "produk": {
      "idProduk": 1,
      "kodeProduk": "PRD-00001",
      "namaProduk": "Botox Bionex 100 IU",
      "kategori": { "idKategori": 2, "namaKategori": "Injectable" },
      "bentukKemasan": "vial",
      "satuanIsi": "IU",
      "isiPerKemasan": 100.00,
      "polaPenggunaan": "PARTIAL_USE",
      "totalIsiTersedia": 450.00
    },
    "daftarBatch": [
      {
        "idBatch": 3,
        "kodeBatch": "BCH-1-20241101-001",
        "expiredDate": "2026-06-30",
        "stokKemasan": 4,
        "totalIsiTersedia": 400.00,
        "statusBatch": "AKTIF",
        "statusExpired": "AMAN",
        "kemasanTerbuka": {
          "idKemasanTerbuka": 2,
          "tanggalDibuka": "2024-11-05",
          "bud": "2024-12-03",
          "isiAwal": 100.00,
          "isiTersisa": 50.00,
          "statusBud": "AKTIF"
        }
      },
      {
        "idBatch": 1,
        "kodeBatch": "BCH-1-20240201-001",
        "expiredDate": "2024-08-15",
        "stokKemasan": 0,
        "totalIsiTersedia": 0.00,
        "statusBatch": "KADALUWARSA",
        "statusExpired": "KADALUWARSA",
        "kemasanTerbuka": null
      }
    ]
  }
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `404` | Produk tidak ditemukan | `"Produk tidak ditemukan."` |

---

## 11. API-07 — Stock Opname

**Terkait KF:** KF-09

### 11.1 Ambil Daftar Sesi Opname

```
GET /api/opname
```

**Query Parameters (opsional):**

| Parameter | Tipe | Keterangan |
|---|---|---|
| `page` | `integer` | Halaman (default: 1) |
| `limit` | `integer` | Jumlah per halaman (default: 20) |

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Data stok opname berhasil diambil.",
  "data": {
    "items": [
      {
        "idOpname": 5,
        "tanggalOpname": "2024-11-30",
        "statusOpname": "SELESAI",
        "jumlahItemDiperiksa": 38,
        "jumlahItemSelisih": 2,
        "dibuatOleh": "admin_farmasi"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 12,
      "totalPages": 1
    }
  }
}
```

---

### 11.2 Mulai Sesi Opname Baru

```
POST /api/opname
```

**Request Body:** Tidak diperlukan (sesi dibuat untuk tanggal hari ini; `idUser` diambil dari token).

**Respons Sukses — `201 Created`:**

Backend membuat sesi opname baru dan mengembalikan seluruh item yang perlu diperiksa.

```json
{
  "success": true,
  "message": "Sesi stok opname baru berhasil dimulai.",
  "data": {
    "idOpname": 6,
    "tanggalOpname": "2024-12-01",
    "statusOpname": "SELESAI",
    "daftarItem": [
      {
        "tipeItem": "batch",
        "idBatch": 3,
        "kodeBatch": "BCH-1-20241101-001",
        "produk": {
          "idProduk": 1,
          "namaProduk": "Botox Bionex 100 IU",
          "satuanIsi": "IU"
        },
        "expiredDate": "2026-06-30",
        "stokSistem": 4,
        "stokFisik": null
      },
      {
        "tipeItem": "kemasanTerbuka",
        "idKemasanTerbuka": 2,
        "idBatch": 3,
        "produk": {
          "idProduk": 1,
          "namaProduk": "Botox Bionex 100 IU",
          "satuanIsi": "IU"
        },
        "bud": "2024-12-03",
        "statusBud": "AKTIF",
        "stokSistem": 50.00,
        "stokFisik": null
      }
    ]
  }
}
```

---

### 11.3 Ambil Detail Sesi Opname

```
GET /api/opname/:id
```

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Detail stok opname berhasil diambil.",
  "data": {
    "idOpname": 5,
    "tanggalOpname": "2024-11-30",
    "statusOpname": "SELESAI",
    "dibuatOleh": "admin_farmasi",
    "detailItem": [
      {
        "idDetailOpname": 10,
        "tipeItem": "batch",
        "idBatch": 3,
        "kodeBatch": "BCH-1-20241101-001",
        "produk": {
          "idProduk": 1,
          "namaProduk": "Botox Bionex 100 IU",
          "satuanIsi": "IU"
        },
        "stokSistem": 5,
        "stokFisik": 4,
        "selisih": -1,
        "keterangan": "Satu vial pecah saat penyimpanan"
      }
    ]
  }
}
```

---

### 11.4 Selesaikan Sesi Opname (Input Hasil Fisik & Simpan)

```
PUT /api/opname/:id/selesai
```

**Request Body:**

```json
{
  "items": [
    {
      "tipeItem": "batch",
      "idBatch": 3,
      "stokFisik": 4,
      "keterangan": "Satu vial pecah saat penyimpanan"
    },
    {
      "tipeItem": "kemasanTerbuka",
      "idKemasanTerbuka": 2,
      "stokFisik": 45.00,
      "keterangan": ""
    }
  ]
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `items` | `array` | ✓ | Daftar hasil pemeriksaan fisik |
| `items[].tipeItem` | `string` | ✓ | `"batch"` atau `"kemasanTerbuka"` |
| `items[].idBatch` | `integer` | Kondisional | Wajib jika `tipeItem = "batch"` |
| `items[].idKemasanTerbuka` | `integer` | Kondisional | Wajib jika `tipeItem = "kemasanTerbuka"` |
| `items[].stokFisik` | `number` | ✓ | Hasil hitungan fisik; harus ≥ 0 |
| `items[].keterangan` | `string` | Kondisional | Wajib diisi jika `stokFisik ≠ stokSistem` (BR-09.2) |

**Respons Sukses — `200 OK`:**

```json
{
  "success": true,
  "message": "Stok opname berhasil disimpan.",
  "data": {
    "idOpname": 6,
    "statusOpname": "SELESAI",
    "jumlahItemDiperiksa": 15,
    "jumlahItemSelisih": 1
  }
}
```

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `400` | `stokFisik` < 0 | `"Stok fisik tidak boleh bernilai negatif."` |
| `400` | Keterangan kosong saat ada selisih | `"Keterangan wajib diisi untuk setiap item yang mengalami selisih."` |
| `404` | Sesi opname tidak ditemukan | `"Sesi stok opname tidak ditemukan."` |

---

## 12. API-08 — Laporan Stok

**Terkait KF:** KF-10

### 12.1 Laporan Stok Masuk

```
GET /api/laporan/stok-masuk
```

**Query Parameters:**

| Parameter | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `periodeAwal` | `string` (DATE) | ✓ | Tanggal mulai periode laporan |
| `periodeAkhir` | `string` (DATE) | ✓ | Tanggal akhir periode laporan (harus ≥ `periodeAwal`) |
| `idKategori` | `integer` | — | Filter per kategori (opsional; kosong = semua kategori) |
| `format` | `string` | — | `json` (default), `pdf`, atau `excel` |

**Respons Sukses — `200 OK` (format JSON):**

```json
{
  "success": true,
  "message": "Laporan berhasil dibuat.",
  "data": {
    "parameter": {
      "periodeAwal": "2024-11-01",
      "periodeAkhir": "2024-11-30",
      "kategori": null
    },
    "ringkasan": {
      "totalTransaksi": 25,
      "totalKemasan": 87,
      "totalIsiMasuk": 8200.00
    },
    "items": [
      {
        "idStokMasuk": 1,
        "tanggalPenerimaan": "2024-11-01",
        "produk": {
          "kodeProduk": "PRD-00001",
          "namaProduk": "Botox Bionex 100 IU",
          "kategori": "Injectable",
          "bentukKemasan": "vial",
          "satuanIsi": "IU"
        },
        "batch": {
          "kodeBatch": "BCH-1-20241101-001",
          "expiredDate": "2026-06-30"
        },
        "jumlahKemasan": 5,
        "totalIsiMasuk": 500.00
      }
    ]
  }
}
```

**Respons format PDF/Excel:** Backend mengembalikan file biner dengan header `Content-Disposition: attachment; filename="laporan-stok-masuk-2024-11.pdf"` (atau `.xlsx`). Frontend menerima sebagai Blob dan memicu unduhan.

**Respons Error:**

| HTTP Code | Kondisi | Pesan |
|---|---|---|
| `400` | `periodeAkhir` < `periodeAwal` | `"Tanggal akhir tidak boleh lebih kecil dari tanggal awal."` |
| `400` | Parameter wajib tidak diisi | `"Periode awal dan akhir wajib diisi."` |

---

### 12.2 Laporan Stok Keluar

```
GET /api/laporan/stok-keluar
```

**Query Parameters:** Sama dengan Laporan Stok Masuk (§ 12.1).

**Respons Sukses — `200 OK` (format JSON):**

```json
{
  "success": true,
  "message": "Laporan berhasil dibuat.",
  "data": {
    "parameter": {
      "periodeAwal": "2024-11-01",
      "periodeAkhir": "2024-11-30",
      "kategori": null
    },
    "ringkasan": {
      "totalTransaksi": 60,
      "totalIsiKeluar": 3200.00
    },
    "items": [
      {
        "idStokKeluar": 1,
        "tanggalPenggunaan": "2024-11-05",
        "produk": {
          "kodeProduk": "PRD-00001",
          "namaProduk": "Botox Bionex 100 IU",
          "kategori": "Injectable",
          "satuanIsi": "IU",
          "polaPenggunaan": "PARTIAL_USE"
        },
        "batch": {
          "kodeBatch": "BCH-1-20241101-001",
          "expiredDate": "2026-06-30"
        },
        "jumlahKemasanDipakai": null,
        "jumlahIsiDipakai": 50.00,
        "keterangan": "Tindakan Botox Area Dahi"
      }
    ]
  }
}
```

---

### 12.3 Laporan Sisa Stok

```
GET /api/laporan/sisa-stok
```

**Query Parameters:** Sama dengan Laporan Stok Masuk; `periodeAkhir` digunakan sebagai titik acuan posisi stok yang dilaporkan.

**Respons Sukses — `200 OK` (format JSON):**

```json
{
  "success": true,
  "message": "Laporan berhasil dibuat.",
  "data": {
    "parameter": {
      "periodeAkhir": "2024-11-30",
      "kategori": null
    },
    "items": [
      {
        "produk": {
          "kodeProduk": "PRD-00001",
          "namaProduk": "Botox Bionex 100 IU",
          "kategori": "Injectable",
          "bentukKemasan": "vial",
          "satuanIsi": "IU",
          "polaPenggunaan": "PARTIAL_USE"
        },
        "daftarBatch": [
          {
            "kodeBatch": "BCH-1-20241101-001",
            "expiredDate": "2026-06-30",
            "stokKemasan": 4,
            "totalIsiTersedia": 400.00,
            "statusBatch": "AKTIF",
            "isiTersisaKemasanTerbuka": 50.00
          }
        ],
        "totalStokKemasan": 4,
        "totalIsiTersedia": 450.00
      }
    ]
  }
}
```

---

## 13. Kode Error & Pesan Standar

Tabel berikut merangkum seluruh pesan error yang digunakan oleh API, konsisten dengan pesan yang didefinisikan pada `srs-fr.md` dan `srs-backend.md` § Error Sentinel.

| Kode HTTP | Kode Error Internal | Pesan |
|---|---|---|
| `401` | `ERR_KREDENSIAL_TIDAK_VALID` | `"Kredensial tidak valid."` |
| `401` | `ERR_TOKEN_TIDAK_VALID` | `"Akses tidak diizinkan."` |
| `401` | `ERR_TOKEN_KADALUWARSA` | `"Sesi telah berakhir. Silakan login kembali."` |
| `400` | `ERR_VALIDASI_INPUT` | `"Data yang dikirim tidak valid."` (dengan `errors[]`) |
| `400` | `ERR_TANGGAL_PENERIMAAN` | `"Tanggal penerimaan tidak boleh melebihi tanggal hari ini."` |
| `400` | `ERR_EXPIRED_DATE` | `"Tanggal kedaluwarsa harus lebih besar dari tanggal penerimaan."` |
| `400` | `ERR_STOK_TIDAK_CUKUP` | `"Stok tidak mencukupi untuk transaksi ini."` |
| `400` | `ERR_PARAMETER_LAPORAN` | `"Tanggal akhir tidak boleh lebih kecil dari tanggal awal."` |
| `404` | `ERR_KATEGORI_NOT_FOUND` | `"Kategori tidak ditemukan."` |
| `404` | `ERR_PRODUK_NOT_FOUND` | `"Produk tidak ditemukan."` |
| `404` | `ERR_BATCH_NOT_FOUND` | `"Tidak ada stok aktif tersedia untuk produk ini."` |
| `404` | `ERR_OPNAME_NOT_FOUND` | `"Sesi stok opname tidak ditemukan."` |
| `409` | `ERR_DUPLIKASI_KATEGORI` | `"Nama kategori sudah terdaftar dalam sistem."` |
| `409` | `ERR_KATEGORI_MEMILIKI_PRODUK` | `"Kategori tidak dapat dihapus karena masih memiliki produk terkait."` |
| `409` | `ERR_PRODUK_STOK_AKTIF` | `"Produk tidak dapat dihapus karena masih memiliki stok aktif."` |
| `409` | `ERR_PRODUK_RIWAYAT_TRANSAKSI` | `"Produk tidak dapat dihapus karena memiliki riwayat transaksi."` |
| `500` | `ERR_INTERNAL` | `"Terjadi kesalahan pada server. Silakan coba beberapa saat lagi."` |

---

## 14. Matriks Traceability API

| Grup Endpoint | Endpoint | Method | KF Terkait | Komponen Backend |
|---|---|---|---|---|
| API-01 Auth | `/api/auth/login` | POST | KF-01 | `auth_handler` → `auth_service` |
| API-01 Auth | `/api/auth/password` | PUT | KF-01 | `auth_handler` → `auth_service` |
| API-02 Kategori | `/api/kategori` | GET | KF-02 | `kategori_handler` → `kategori_service` |
| API-02 Kategori | `/api/kategori` | POST | KF-02 | `kategori_handler` → `kategori_service` |
| API-02 Kategori | `/api/kategori/:id` | PUT | KF-02 | `kategori_handler` → `kategori_service` |
| API-02 Kategori | `/api/kategori/:id` | DELETE | KF-02 | `kategori_handler` → `kategori_service` |
| API-03 Produk | `/api/produk` | GET | KF-03 | `produk_handler` → `produk_service` |
| API-03 Produk | `/api/produk/:id` | GET | KF-03 | `produk_handler` → `produk_service` |
| API-03 Produk | `/api/produk` | POST | KF-03 | `produk_handler` → `produk_service` |
| API-03 Produk | `/api/produk/:id` | PUT | KF-03 | `produk_handler` → `produk_service` |
| API-03 Produk | `/api/produk/:id` | DELETE | KF-03 | `produk_handler` → `produk_service` |
| API-04 Stok Masuk | `/api/stok-masuk` | GET | KF-04 | `stok_masuk_handler` → `stok_masuk_service` |
| API-04 Stok Masuk | `/api/stok-masuk` | POST | KF-04 | `stok_masuk_handler` → `stok_masuk_service` (DB Tx) |
| API-05 Stok Keluar | `/api/stok-keluar/preview-batch` | GET | KF-05, KF-06 | `stok_keluar_handler` → `stok_keluar_service` |
| API-05 Stok Keluar | `/api/stok-keluar` | GET | KF-05 | `stok_keluar_handler` → `stok_keluar_service` |
| API-05 Stok Keluar | `/api/stok-keluar` | POST | KF-05, KF-06, KF-07 | `stok_keluar_handler` → `stok_keluar_service` (DB Tx) |
| API-06 Monitoring | `/api/monitoring` | GET | KF-08 | `monitoring_handler` → `monitoring_service` |
| API-06 Monitoring | `/api/monitoring/:idProduk` | GET | KF-08 | `monitoring_handler` → `monitoring_service` |
| API-07 Opname | `/api/opname` | GET | KF-09 | `opname_handler` → `opname_service` |
| API-07 Opname | `/api/opname` | POST | KF-09 | `opname_handler` → `opname_service` |
| API-07 Opname | `/api/opname/:id` | GET | KF-09 | `opname_handler` → `opname_service` |
| API-07 Opname | `/api/opname/:id/selesai` | PUT | KF-09 | `opname_handler` → `opname_service` (DB Tx) |
| API-08 Laporan | `/api/laporan/stok-masuk` | GET | KF-10 | `laporan_handler` → `laporan_service` |
| API-08 Laporan | `/api/laporan/stok-keluar` | GET | KF-10 | `laporan_handler` → `laporan_service` |
| API-08 Laporan | `/api/laporan/sisa-stok` | GET | KF-10 | `laporan_handler` → `laporan_service` |
