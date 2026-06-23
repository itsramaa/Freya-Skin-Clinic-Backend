# SRS — Database Specification
## Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic

---

## Daftar Isi

1. [Tujuan & Cakupan](#1-tujuan--cakupan)
2. [Justifikasi Teknologi](#2-justifikasi-teknologi)
3. [Konvensi Penamaan](#3-konvensi-penamaan)
4. [Skema Tabel](#4-skema-tabel)
   - [TBL-01 — users](#41-tbl-01--users)
   - [TBL-02 — kategori](#42-tbl-02--kategori)
   - [TBL-03 — produk](#43-tbl-03--produk)
   - [TBL-04 — stok_masuk](#44-tbl-04--stok_masuk)
   - [TBL-05 — batch_stok](#45-tbl-05--batch_stok)
   - [TBL-06 — stok_keluar](#46-tbl-06--stok_keluar)
   - [TBL-07 — kemasan_terbuka](#47-tbl-07--kemasan_terbuka)
   - [TBL-08 — stok_opname](#48-tbl-08--stok_opname)
   - [TBL-09 — detail_opname](#49-tbl-09--detail_opname)
5. [Entity Relationship Diagram (Deskriptif)](#5-entity-relationship-diagram-deskriptif)
6. [Relasi Antar Tabel](#6-relasi-antar-tabel)
7. [Index](#7-index)
8. [Constraint & Business Rules pada Level Database](#8-constraint--business-rules-pada-level-database)
9. [DDL Lengkap (SQL)](#9-ddl-lengkap-sql)
10. [Catatan Integritas Data](#10-catatan-integritas-data)

---

## 1. Tujuan & Cakupan

Dokumen ini merinci spesifikasi lapisan penyimpanan data Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic, yang dibangun menggunakan **PostgreSQL** sebagai basis data relasional sesuai arsitektur yang ditetapkan pada `srs-overview.md` § 8.3. Dokumen ini mencakup skema seluruh tabel, relasi antar entitas, index, constraint, dan DDL (Data Definition Language) yang digunakan sebagai acuan implementasi dan pengujian integritas data.

---

## 2. Justifikasi Teknologi

Lapisan penyimpanan data menggunakan **PostgreSQL** sebagai basis data relasional untuk menyimpan seluruh data operasional secara persisten, mencakup data kategori, produk, batch stok, transaksi stok masuk dan keluar, kemasan terbuka, serta riwayat stock opname.

Pemilihan PostgreSQL didasarkan pada kebutuhan integritas relasional yang tinggi, mengingat data operasional farmasi memiliki banyak relasi antar entitas yang harus terjaga konsistensinya — data stok keluar harus selalu terhubung ke batch yang valid, data kemasan terbuka harus terikat pada batch yang aktif, dan data detail opname harus terhubung ke sesi opname yang bersangkutan. PostgreSQL mendukung **ACID compliance** penuh, yang relevan langsung dengan KNF-04 (Reliability & Akurasi Data) yang mewajibkan seluruh operasi multi-tabel dibungkus dalam transaksi database agar tidak terjadi data parsial apabila terjadi kegagalan di tengah proses.

Selain itu, PostgreSQL mendukung tipe data `ENUM`, `DECIMAL` dengan presisi tinggi, dan `TIMESTAMP WITH TIME ZONE` yang sesuai dengan kebutuhan pencatatan expired date, BUD, dan stempel waktu transaksi yang akurat pada sistem farmasi ini.

---

## 3. Konvensi Penamaan

| Elemen | Konvensi | Contoh |
|---|---|---|
| Nama tabel | `snake_case`, jamak | `batch_stok`, `stok_masuk` |
| Nama kolom | `snake_case` | `expired_date`, `pola_penggunaan` |
| Primary key | `id_<nama_tabel_singular>` | `id_batch`, `id_produk` |
| Foreign key | `id_<tabel_referensi_singular>` | `id_kategori`, `id_batch` |
| ENUM type | `UPPERCASE` | `AKTIF`, `HABIS`, `KADALUWARSA` |
| Index | `idx_<tabel>_<kolom>` | `idx_batch_stok_id_produk` |
| Constraint FK | `fk_<tabel>_<kolom>` | `fk_produk_id_kategori` |
| Constraint UNIQUE | `uq_<tabel>_<kolom>` | `uq_kategori_nama` |

---

## 4. Skema Tabel

### 4.1 TBL-01 — `users`

Menyimpan data akun pengguna. Sistem hanya memiliki satu peran pengguna (Admin Farmasi), sehingga tidak diperlukan kolom role (B-02).

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id_user` | `SERIAL` | `PRIMARY KEY` | Identitas unik pengguna |
| `username` | `VARCHAR(50)` | `NOT NULL, UNIQUE` | Nama pengguna untuk login |
| `password` | `VARCHAR(255)` | `NOT NULL` | Password dalam bentuk bcrypt hash (KNF-02, BR-01.1) |
| `is_default_password` | `BOOLEAN` | `NOT NULL, DEFAULT TRUE` | Flag password default; `TRUE` = belum diganti (BR-01.2) |
| `created_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembuatan akun |
| `updated_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembaruan terakhir |

**Catatan:** Kolom `password` tidak pernah dikembalikan ke frontend melalui API (KNF-02, AC-NFR02.2).

---

### 4.2 TBL-02 — `kategori`

Menyimpan data master kategori produk farmasi. Kategori digunakan sebagai pengelompokan pada seluruh modul (produk, monitoring, laporan).

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id_kategori` | `SERIAL` | `PRIMARY KEY` | Identitas unik kategori |
| `kode_kategori` | `VARCHAR(20)` | `NOT NULL, UNIQUE` | Kode unik kategori, digenerate otomatis oleh sistem |
| `nama_kategori` | `VARCHAR(100)` | `NOT NULL, UNIQUE` | Nama kategori; validasi duplikasi pada level DB (KF-02) |
| `created_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembuatan |
| `updated_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembaruan terakhir |

**Constraint bisnis:** Penghapusan hanya diperbolehkan jika tidak ada `produk` yang merujuk ke `id_kategori` ini (diperiksa di service layer, bukan ON DELETE CASCADE — sesuai KF-02 A1).

---

### 4.3 TBL-03 — `produk`

Menyimpan data master produk farmasi. Atribut `pola_penggunaan` menentukan mekanisme pengelolaan stok keluar dan BUD di seluruh modul berikutnya (KF-03, KF-05, KF-07).

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id_produk` | `SERIAL` | `PRIMARY KEY` | Identitas unik produk |
| `kode_produk` | `VARCHAR(20)` | `NOT NULL, UNIQUE` | Kode unik produk, digenerate otomatis oleh sistem |
| `nama_produk` | `VARCHAR(150)` | `NOT NULL` | Nama produk farmasi |
| `id_kategori` | `INTEGER` | `NOT NULL, FK → kategori.id_kategori` | Kategori produk (KF-03) |
| `bentuk_kemasan` | `VARCHAR(50)` | `NOT NULL` | Bentuk kemasan (vial, jar, botol, ampoule, syringe, set, pack, box, dll.) |
| `satuan_isi` | `VARCHAR(20)` | `NOT NULL` | Satuan isi kemasan (pcs, ml, cc, gram, IU, dll.) |
| `isi_per_kemasan` | `DECIMAL(10,2)` | `NOT NULL, CHECK > 0` | Kapasitas isi per kemasan sesuai satuan |
| `pola_penggunaan` | `VARCHAR(20)` | `NOT NULL, CHECK IN ('FULL_USE','PARTIAL_USE')` | Pola penggunaan produk; menentukan alur stok keluar dan BUD |
| `created_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembuatan |
| `updated_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembaruan terakhir |

**Catatan kolom `pola_penggunaan`:** Nilai `FULL_USE` berarti produk dihabiskan dalam satu pemakaian; `PARTIAL_USE` berarti kemasan dapat digunakan sebagian dan menyisakan isi yang perlu dilacak via `kemasan_terbuka`.

---

### 4.4 TBL-04 — `stok_masuk`

Menyimpan setiap transaksi penerimaan produk dari supplier. Setiap baris mewakili satu kejadian penerimaan. Tabel ini bersifat *append-only* — tidak ada penghapusan data historis (KNF-05).

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id_stok_masuk` | `SERIAL` | `PRIMARY KEY` | Identitas unik transaksi penerimaan |
| `id_produk` | `INTEGER` | `NOT NULL, FK → produk.id_produk` | Produk yang diterima |
| `id_batch` | `INTEGER` | `NOT NULL, FK → batch_stok.id_batch` | Batch yang terbentuk atau diperbarui dari penerimaan ini |
| `id_user` | `INTEGER` | `NOT NULL, FK → users.id_user` | Admin Farmasi yang mencatat penerimaan |
| `tanggal_penerimaan` | `DATE` | `NOT NULL` | Tanggal fisik penerimaan barang (≤ tanggal saat ini) |
| `jumlah_kemasan` | `INTEGER` | `NOT NULL, CHECK > 0` | Jumlah kemasan yang diterima |
| `total_isi_masuk` | `DECIMAL(12,2)` | `NOT NULL, CHECK > 0` | Total isi = `jumlah_kemasan × isi_per_kemasan` produk |
| `keterangan` | `TEXT` | | Catatan tambahan (opsional) |
| `created_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pencatatan di sistem |

**Logika batch:** Jika produk dan `expired_date` sama dengan batch yang sudah ada, stok ditambahkan ke batch tersebut (relasi `stok_masuk → batch_stok` = N:1 untuk kasus ini). Jika belum ada batch dengan kombinasi tersebut, batch baru dibuat. Logika ini dikelola di service layer (`stok_masuk_service.go`), bukan di level database trigger.

---

### 4.5 TBL-05 — `batch_stok`

Menyimpan setiap batch produk. Batch merupakan unit utama penerapan FEFO — sistem menentukan prioritas penggunaan berdasarkan `expired_date` terkecil di antara batch berstatus `AKTIF` (KF-06).

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id_batch` | `SERIAL` | `PRIMARY KEY` | Identitas unik batch |
| `id_produk` | `INTEGER` | `NOT NULL, FK → produk.id_produk` | Produk yang diwakili batch ini |
| `kode_batch` | `VARCHAR(30)` | `NOT NULL, UNIQUE` | Kode batch unik, digenerate otomatis oleh sistem |
| `expired_date` | `DATE` | `NOT NULL` | Tanggal kedaluwarsa produk pada batch ini |
| `stok_kemasan` | `INTEGER` | `NOT NULL, DEFAULT 0, CHECK >= 0` | Jumlah kemasan utuh yang tersedia |
| `total_isi_tersedia` | `DECIMAL(12,2)` | `NOT NULL, DEFAULT 0, CHECK >= 0` | Total isi tersedia pada batch (termasuk kemasan utuh) |
| `status_batch` | `VARCHAR(20)` | `NOT NULL, DEFAULT 'AKTIF', CHECK IN ('AKTIF','HABIS','KADALUWARSA')` | Status batch; diperbarui otomatis oleh background worker dan service |
| `created_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembuatan batch |
| `updated_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembaruan terakhir |

**Constraint unik komposit:** `UNIQUE(id_produk, expired_date)` — memastikan satu produk dengan expired date yang sama hanya memiliki satu batch, konsisten dengan logika penggabungan stok masuk.

**Status batch:**
- `AKTIF`: `stok_kemasan > 0` dan `expired_date >= CURRENT_DATE`
- `HABIS`: `stok_kemasan = 0` dan `expired_date >= CURRENT_DATE`
- `KADALUWARSA`: `expired_date < CURRENT_DATE` (diperbarui oleh background worker, KNF-07)

---

### 4.6 TBL-06 — `stok_keluar`

Menyimpan setiap transaksi penggunaan produk dalam pelayanan pasien. Tabel ini bersifat *append-only* (KNF-05). Kolom `jumlah_kemasan_dipakai` dan `jumlah_isi_dipakai` bersifat kondisional sesuai `pola_penggunaan` produk.

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id_stok_keluar` | `SERIAL` | `PRIMARY KEY` | Identitas unik transaksi penggunaan |
| `id_produk` | `INTEGER` | `NOT NULL, FK → produk.id_produk` | Produk yang digunakan |
| `id_batch` | `INTEGER` | `NOT NULL, FK → batch_stok.id_batch` | Batch yang digunakan (dipilih otomatis oleh FEFO) |
| `id_user` | `INTEGER` | `NOT NULL, FK → users.id_user` | Admin Farmasi yang mencatat penggunaan |
| `tanggal_penggunaan` | `DATE` | `NOT NULL` | Tanggal produk digunakan dalam pelayanan |
| `jumlah_kemasan_dipakai` | `INTEGER` | `CHECK >= 0` | Jumlah kemasan utuh yang dipakai; diisi untuk produk `FULL_USE` |
| `jumlah_isi_dipakai` | `DECIMAL(10,2)` | `CHECK >= 0` | Jumlah isi yang dipakai; diisi untuk produk `PARTIAL_USE` |
| `keterangan` | `TEXT` | | Catatan penggunaan (opsional) |
| `created_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pencatatan di sistem |

**Catatan:** Tepat salah satu dari `jumlah_kemasan_dipakai` atau `jumlah_isi_dipakai` akan terisi sesuai `pola_penggunaan` produk yang bersangkutan. Validasi ini dikelola di service layer (bukan CHECK constraint database) karena memerlukan lookup ke tabel `produk`.

---

### 4.7 TBL-07 — `kemasan_terbuka`

Menyimpan informasi kemasan produk `PARTIAL_USE` yang telah dibuka dan masih memiliki sisa isi. Hanya satu kemasan terbuka aktif yang diperbolehkan per batch pada satu waktu (BR-07.2). Tabel ini berkaitan langsung dengan KF-07 (Kelola BUD) dan dipantau oleh background worker (KNF-07).

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id_kemasan_terbuka` | `SERIAL` | `PRIMARY KEY` | Identitas unik kemasan terbuka |
| `id_batch` | `INTEGER` | `NOT NULL, FK → batch_stok.id_batch` | Batch asal kemasan terbuka ini |
| `tanggal_dibuka` | `DATE` | `NOT NULL` | Tanggal kemasan pertama kali dibuka |
| `bud` | `DATE` | `NOT NULL` | Beyond Use Date = `tanggal_dibuka + 28 hari` (BR-07.1, B-08) |
| `isi_awal` | `DECIMAL(10,2)` | `NOT NULL, CHECK > 0` | Isi kemasan saat pertama dibuka (= `isi_per_kemasan` produk) |
| `isi_tersisa` | `DECIMAL(10,2)` | `NOT NULL, CHECK >= 0` | Sisa isi yang masih tersedia |
| `status_bud` | `VARCHAR(20)` | `NOT NULL, DEFAULT 'AKTIF', CHECK IN ('AKTIF','KADALUWARSA')` | Status BUD; diperbarui oleh background worker atau saat stok keluar berikutnya |
| `created_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembuatan record |
| `updated_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembaruan terakhir |

**Constraint satu kemasan terbuka aktif per batch:** Dijaga di service layer menggunakan query `FindAktifByBatch()` sebelum membuat kemasan terbuka baru, konsisten dengan relasi `BatchStok → KemasanTerbuka = 1 : 0..1` pada class diagram.

---

### 4.8 TBL-08 — `stok_opname`

Menyimpan header sesi stock opname. Setiap sesi merepresentasikan satu kegiatan pencocokan stok sistem dengan kondisi fisik gudang (KF-09).

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id_opname` | `SERIAL` | `PRIMARY KEY` | Identitas unik sesi opname |
| `id_user` | `INTEGER` | `NOT NULL, FK → users.id_user` | Admin Farmasi yang memulai sesi opname |
| `tanggal_opname` | `DATE` | `NOT NULL` | Tanggal sesi opname dilakukan |
| `status_opname` | `VARCHAR(20)` | `NOT NULL, DEFAULT 'SELESAI', CHECK IN ('SELESAI','DIBATALKAN')` | Status akhir sesi opname |
| `created_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembuatan sesi |
| `updated_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pembaruan terakhir |

---

### 4.9 TBL-09 — `detail_opname`

Menyimpan rincian hasil pemeriksaan setiap item pada suatu sesi opname. Setiap baris mencatat perbandingan stok sistem vs. stok fisik beserta selisih dan keterangannya. Tabel ini adalah implementasi dari pencatatan histori koreksi stok yang menjawab permasalahan P-04 dan memenuhi KNF-05 (append-only, tidak ada overwrite data historis).

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id_detail_opname` | `SERIAL` | `PRIMARY KEY` | Identitas unik detail opname |
| `id_opname` | `INTEGER` | `NOT NULL, FK → stok_opname.id_opname` | Sesi opname induk |
| `id_batch` | `INTEGER` | `NOT NULL, FK → batch_stok.id_batch` | Batch yang diperiksa |
| `id_kemasan_terbuka` | `INTEGER` | `FK → kemasan_terbuka.id_kemasan_terbuka, NULLABLE` | Kemasan terbuka yang diperiksa; `NULL` jika item adalah kemasan utuh (BR-09.3) |
| `stok_sistem` | `DECIMAL(10,2)` | `NOT NULL` | Nilai stok menurut sistem saat opname dilakukan |
| `stok_fisik` | `DECIMAL(10,2)` | `NOT NULL, CHECK >= 0` | Nilai stok hasil perhitungan fisik oleh Admin Farmasi |
| `selisih` | `DECIMAL(10,2)` | `NOT NULL` | `stok_fisik − stok_sistem`; positif = surplus, negatif = kurang (BR-09.1) |
| `keterangan` | `TEXT` | `NOT NULL` jika `selisih ≠ 0` | Penjelasan selisih; wajib diisi apabila ditemukan perbedaan (BR-09.2) |
| `created_at` | `TIMESTAMP WITH TIME ZONE` | `NOT NULL, DEFAULT NOW()` | Stempel waktu pencatatan |

**Catatan `id_kemasan_terbuka`:** Tepat salah satu dari (kemasan utuh via `id_batch`) atau (kemasan terbuka via `id_kemasan_terbuka`) yang relevan per baris. Jika `id_kemasan_terbuka IS NOT NULL`, berarti item yang diperiksa adalah kemasan terbuka; nilai `stok_sistem` adalah `isi_tersisa` pada kemasan tersebut.

---

## 5. Entity Relationship Diagram (Deskriptif)

Sistem memiliki 9 entitas utama dengan relasi sebagai berikut (dideskripsikan sesuai class diagram § 4.3.4 PDF dan dipetakan ke tabel relasional):

```
users (1) ──────────────────────── (N) stok_masuk
users (1) ──────────────────────── (N) stok_keluar
users (1) ──────────────────────── (N) stok_opname

kategori (1) ───────────────────── (N) produk

produk (1) ──────────────────────── (N) stok_masuk
produk (1) ──────────────────────── (N) batch_stok
produk (1) ──────────────────────── (N) stok_keluar

stok_masuk (N) ──────────────────── (1) batch_stok
  └─ (banyak stok_masuk bisa merujuk ke batch yang sama
      jika produk dan expired_date identik)

batch_stok (1) ──────────────────── (N) stok_keluar
batch_stok (1) ──────────────────── (0..1) kemasan_terbuka
batch_stok (1) ──────────────────── (N) detail_opname

kemasan_terbuka (0..1) ─────────── (N) detail_opname

stok_opname (1) ─────────────────── (N) detail_opname
```

---

## 6. Relasi Antar Tabel

| Tabel Anak | Kolom FK | Tabel Induk | Kolom PK | ON DELETE | Keterangan |
|---|---|---|---|---|---|
| `produk` | `id_kategori` | `kategori` | `id_kategori` | `RESTRICT` | Kategori tidak bisa dihapus jika ada produk |
| `stok_masuk` | `id_produk` | `produk` | `id_produk` | `RESTRICT` | Produk tidak bisa dihapus jika ada transaksi |
| `stok_masuk` | `id_batch` | `batch_stok` | `id_batch` | `RESTRICT` | — |
| `stok_masuk` | `id_user` | `users` | `id_user` | `RESTRICT` | — |
| `batch_stok` | `id_produk` | `produk` | `id_produk` | `RESTRICT` | Produk tidak bisa dihapus jika ada batch |
| `stok_keluar` | `id_produk` | `produk` | `id_produk` | `RESTRICT` | — |
| `stok_keluar` | `id_batch` | `batch_stok` | `id_batch` | `RESTRICT` | — |
| `stok_keluar` | `id_user` | `users` | `id_user` | `RESTRICT` | — |
| `kemasan_terbuka` | `id_batch` | `batch_stok` | `id_batch` | `RESTRICT` | — |
| `stok_opname` | `id_user` | `users` | `id_user` | `RESTRICT` | — |
| `detail_opname` | `id_opname` | `stok_opname` | `id_opname` | `CASCADE` | Detail dihapus bersama sesi opname |
| `detail_opname` | `id_batch` | `batch_stok` | `id_batch` | `RESTRICT` | — |
| `detail_opname` | `id_kemasan_terbuka` | `kemasan_terbuka` | `id_kemasan_terbuka` | `SET NULL` | Nullable; kemasan bisa dihapus terlepas dari histori |

**Catatan `ON DELETE RESTRICT`:** Digunakan secara dominan agar data historis selalu memiliki referensi yang valid. Penghapusan data master (kategori, produk) dikendalikan di service layer dengan pengecekan eksplisit sesuai activity diagram (KF-02, KF-03).

---

## 7. Index

Index dirancang berdasarkan pola query yang paling sering digunakan (KNF-06 — Performance), khususnya pada endpoint Monitoring Stok (KF-08) yang membaca data produk, batch, dan kemasan terbuka secara bersamaan, serta pada proses FEFO yang memerlukan pengurutan batch berdasarkan `expired_date`.

| Nama Index | Tabel | Kolom | Tipe | Justifikasi |
|---|---|---|---|---|
| `idx_produk_id_kategori` | `produk` | `id_kategori` | B-tree | Filter produk per kategori (monitoring, laporan) |
| `idx_batch_stok_id_produk` | `batch_stok` | `id_produk` | B-tree | Lookup batch milik produk (FEFO, monitoring) |
| `idx_batch_stok_expired_date` | `batch_stok` | `expired_date` | B-tree | Pengurutan FEFO (ASC), pembaruan status oleh worker |
| `idx_batch_stok_status` | `batch_stok` | `status_batch` | B-tree | Filter batch AKTIF/HABIS/KADALUWARSA |
| `idx_batch_stok_produk_status` | `batch_stok` | `(id_produk, status_batch)` | Komposit | Query FEFO: batch AKTIF milik produk tertentu |
| `idx_stok_masuk_id_produk` | `stok_masuk` | `id_produk` | B-tree | Laporan stok masuk per produk |
| `idx_stok_masuk_tanggal` | `stok_masuk` | `tanggal_penerimaan` | B-tree | Filter laporan berdasarkan periode |
| `idx_stok_keluar_id_produk` | `stok_keluar` | `id_produk` | B-tree | Laporan stok keluar per produk |
| `idx_stok_keluar_tanggal` | `stok_keluar` | `tanggal_penggunaan` | B-tree | Filter laporan berdasarkan periode |
| `idx_kemasan_terbuka_id_batch` | `kemasan_terbuka` | `id_batch` | B-tree | Lookup kemasan terbuka milik batch |
| `idx_kemasan_terbuka_status` | `kemasan_terbuka` | `status_bud` | B-tree | Pembaruan status BUD oleh background worker |
| `idx_kemasan_terbuka_bud` | `kemasan_terbuka` | `bud` | B-tree | Worker: kemasan dengan BUD yang telah terlewati |
| `idx_detail_opname_id_opname` | `detail_opname` | `id_opname` | B-tree | Lookup detail milik sesi opname |

---

## 8. Constraint & Business Rules pada Level Database

Constraint berikut diterapkan langsung pada DDL untuk menjaga integritas data bahkan ketika diakses di luar service layer:

| Tabel | Constraint | Ekspresi | KF/BR Terkait |
|---|---|---|---|
| `users` | `uq_users_username` | `UNIQUE(username)` | KF-01 |
| `kategori` | `uq_kategori_nama` | `UNIQUE(nama_kategori)` | KF-02, BR duplikasi |
| `kategori` | `uq_kategori_kode` | `UNIQUE(kode_kategori)` | KF-02 |
| `produk` | `uq_produk_kode` | `UNIQUE(kode_produk)` | KF-03 |
| `produk` | `chk_produk_isi` | `isi_per_kemasan > 0` | KF-03 |
| `produk` | `chk_produk_pola` | `pola_penggunaan IN ('FULL_USE','PARTIAL_USE')` | KF-03, KF-05 |
| `batch_stok` | `uq_batch_produk_expired` | `UNIQUE(id_produk, expired_date)` | KF-04, FEFO |
| `batch_stok` | `uq_batch_kode` | `UNIQUE(kode_batch)` | KF-04 |
| `batch_stok` | `chk_batch_stok_kemasan` | `stok_kemasan >= 0` | KF-05, KNF-04 |
| `batch_stok` | `chk_batch_total_isi` | `total_isi_tersedia >= 0` | KF-05, KNF-04 |
| `batch_stok` | `chk_batch_status` | `status_batch IN ('AKTIF','HABIS','KADALUWARSA')` | KF-06 |
| `stok_masuk` | `chk_masuk_jumlah` | `jumlah_kemasan > 0` | KF-04 |
| `stok_keluar` | `chk_keluar_kemasan` | `jumlah_kemasan_dipakai >= 0` | KF-05 |
| `stok_keluar` | `chk_keluar_isi` | `jumlah_isi_dipakai >= 0` | KF-05 |
| `kemasan_terbuka` | `chk_kt_isi_awal` | `isi_awal > 0` | KF-07 |
| `kemasan_terbuka` | `chk_kt_isi_tersisa` | `isi_tersisa >= 0` | KF-07, KNF-04 |
| `kemasan_terbuka` | `chk_kt_status` | `status_bud IN ('AKTIF','KADALUWARSA')` | KF-07 |
| `stok_opname` | `chk_opname_status` | `status_opname IN ('SELESAI','DIBATALKAN')` | KF-09 |
| `detail_opname` | `chk_detail_stok_fisik` | `stok_fisik >= 0` | KF-09, BR-09.1 |

---

## 9. DDL Lengkap (SQL)

```sql
-- ============================================================
-- Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic
-- DDL PostgreSQL
-- ============================================================

-- TBL-01: users
CREATE TABLE users (
    id_user             SERIAL PRIMARY KEY,
    username            VARCHAR(50)  NOT NULL,
    password            VARCHAR(255) NOT NULL,
    is_default_password BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_users_username UNIQUE (username)
);

-- TBL-02: kategori
CREATE TABLE kategori (
    id_kategori     SERIAL PRIMARY KEY,
    kode_kategori   VARCHAR(20)  NOT NULL,
    nama_kategori   VARCHAR(100) NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_kategori_kode  UNIQUE (kode_kategori),
    CONSTRAINT uq_kategori_nama  UNIQUE (nama_kategori)
);

-- TBL-03: produk
CREATE TABLE produk (
    id_produk       SERIAL PRIMARY KEY,
    kode_produk     VARCHAR(20)  NOT NULL,
    nama_produk     VARCHAR(150) NOT NULL,
    id_kategori     INTEGER      NOT NULL,
    bentuk_kemasan  VARCHAR(50)  NOT NULL,
    satuan_isi      VARCHAR(20)  NOT NULL,
    isi_per_kemasan DECIMAL(10,2) NOT NULL,
    pola_penggunaan VARCHAR(20)  NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_produk_kode        UNIQUE (kode_produk),
    CONSTRAINT chk_produk_isi        CHECK  (isi_per_kemasan > 0),
    CONSTRAINT chk_produk_pola       CHECK  (pola_penggunaan IN ('FULL_USE', 'PARTIAL_USE')),
    CONSTRAINT fk_produk_id_kategori FOREIGN KEY (id_kategori)
        REFERENCES kategori (id_kategori) ON DELETE RESTRICT
);

CREATE INDEX idx_produk_id_kategori ON produk (id_kategori);

-- TBL-05: batch_stok (dibuat sebelum stok_masuk karena stok_masuk merujuk ke sini)
CREATE TABLE batch_stok (
    id_batch            SERIAL PRIMARY KEY,
    id_produk           INTEGER      NOT NULL,
    kode_batch          VARCHAR(30)  NOT NULL,
    expired_date        DATE         NOT NULL,
    stok_kemasan        INTEGER      NOT NULL DEFAULT 0,
    total_isi_tersedia  DECIMAL(12,2) NOT NULL DEFAULT 0,
    status_batch        VARCHAR(20)  NOT NULL DEFAULT 'AKTIF',
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_batch_kode           UNIQUE (kode_batch),
    CONSTRAINT uq_batch_produk_expired UNIQUE (id_produk, expired_date),
    CONSTRAINT chk_batch_stok_kemasan  CHECK  (stok_kemasan >= 0),
    CONSTRAINT chk_batch_total_isi     CHECK  (total_isi_tersedia >= 0),
    CONSTRAINT chk_batch_status        CHECK  (status_batch IN ('AKTIF', 'HABIS', 'KADALUWARSA')),
    CONSTRAINT fk_batch_id_produk      FOREIGN KEY (id_produk)
        REFERENCES produk (id_produk) ON DELETE RESTRICT
);

CREATE INDEX idx_batch_stok_id_produk    ON batch_stok (id_produk);
CREATE INDEX idx_batch_stok_expired_date ON batch_stok (expired_date);
CREATE INDEX idx_batch_stok_status       ON batch_stok (status_batch);
CREATE INDEX idx_batch_stok_produk_status ON batch_stok (id_produk, status_batch);

-- TBL-04: stok_masuk
CREATE TABLE stok_masuk (
    id_stok_masuk    SERIAL PRIMARY KEY,
    id_produk        INTEGER       NOT NULL,
    id_batch         INTEGER       NOT NULL,
    id_user          INTEGER       NOT NULL,
    tanggal_penerimaan DATE        NOT NULL,
    jumlah_kemasan   INTEGER       NOT NULL,
    total_isi_masuk  DECIMAL(12,2) NOT NULL,
    keterangan       TEXT,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_masuk_jumlah   CHECK  (jumlah_kemasan > 0),
    CONSTRAINT chk_masuk_isi      CHECK  (total_isi_masuk > 0),
    CONSTRAINT fk_masuk_id_produk FOREIGN KEY (id_produk)
        REFERENCES produk (id_produk) ON DELETE RESTRICT,
    CONSTRAINT fk_masuk_id_batch  FOREIGN KEY (id_batch)
        REFERENCES batch_stok (id_batch) ON DELETE RESTRICT,
    CONSTRAINT fk_masuk_id_user   FOREIGN KEY (id_user)
        REFERENCES users (id_user) ON DELETE RESTRICT
);

CREATE INDEX idx_stok_masuk_id_produk ON stok_masuk (id_produk);
CREATE INDEX idx_stok_masuk_tanggal   ON stok_masuk (tanggal_penerimaan);

-- TBL-06: stok_keluar
CREATE TABLE stok_keluar (
    id_stok_keluar         SERIAL PRIMARY KEY,
    id_produk              INTEGER       NOT NULL,
    id_batch               INTEGER       NOT NULL,
    id_user                INTEGER       NOT NULL,
    tanggal_penggunaan     DATE          NOT NULL,
    jumlah_kemasan_dipakai INTEGER       CHECK (jumlah_kemasan_dipakai >= 0),
    jumlah_isi_dipakai     DECIMAL(10,2) CHECK (jumlah_isi_dipakai >= 0),
    keterangan             TEXT,
    created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_keluar_id_produk FOREIGN KEY (id_produk)
        REFERENCES produk (id_produk) ON DELETE RESTRICT,
    CONSTRAINT fk_keluar_id_batch  FOREIGN KEY (id_batch)
        REFERENCES batch_stok (id_batch) ON DELETE RESTRICT,
    CONSTRAINT fk_keluar_id_user   FOREIGN KEY (id_user)
        REFERENCES users (id_user) ON DELETE RESTRICT
);

CREATE INDEX idx_stok_keluar_id_produk ON stok_keluar (id_produk);
CREATE INDEX idx_stok_keluar_tanggal   ON stok_keluar (tanggal_penggunaan);

-- TBL-07: kemasan_terbuka
CREATE TABLE kemasan_terbuka (
    id_kemasan_terbuka SERIAL PRIMARY KEY,
    id_batch           INTEGER       NOT NULL,
    tanggal_dibuka     DATE          NOT NULL,
    bud                DATE          NOT NULL,
    isi_awal           DECIMAL(10,2) NOT NULL,
    isi_tersisa        DECIMAL(10,2) NOT NULL,
    status_bud         VARCHAR(20)   NOT NULL DEFAULT 'AKTIF',
    created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_kt_isi_awal    CHECK  (isi_awal > 0),
    CONSTRAINT chk_kt_isi_tersisa CHECK  (isi_tersisa >= 0),
    CONSTRAINT chk_kt_status      CHECK  (status_bud IN ('AKTIF', 'KADALUWARSA')),
    CONSTRAINT fk_kt_id_batch     FOREIGN KEY (id_batch)
        REFERENCES batch_stok (id_batch) ON DELETE RESTRICT
);

CREATE INDEX idx_kemasan_terbuka_id_batch ON kemasan_terbuka (id_batch);
CREATE INDEX idx_kemasan_terbuka_status   ON kemasan_terbuka (status_bud);
CREATE INDEX idx_kemasan_terbuka_bud      ON kemasan_terbuka (bud);

-- TBL-08: stok_opname
CREATE TABLE stok_opname (
    id_opname      SERIAL PRIMARY KEY,
    id_user        INTEGER     NOT NULL,
    tanggal_opname DATE        NOT NULL,
    status_opname  VARCHAR(20) NOT NULL DEFAULT 'SELESAI',
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_opname_status  CHECK  (status_opname IN ('SELESAI', 'DIBATALKAN')),
    CONSTRAINT fk_opname_id_user  FOREIGN KEY (id_user)
        REFERENCES users (id_user) ON DELETE RESTRICT
);

-- TBL-09: detail_opname
CREATE TABLE detail_opname (
    id_detail_opname   SERIAL PRIMARY KEY,
    id_opname          INTEGER       NOT NULL,
    id_batch           INTEGER       NOT NULL,
    id_kemasan_terbuka INTEGER,
    stok_sistem        DECIMAL(10,2) NOT NULL,
    stok_fisik         DECIMAL(10,2) NOT NULL,
    selisih            DECIMAL(10,2) NOT NULL,
    keterangan         TEXT,
    created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_detail_stok_fisik CHECK  (stok_fisik >= 0),
    CONSTRAINT fk_detail_id_opname   FOREIGN KEY (id_opname)
        REFERENCES stok_opname (id_opname) ON DELETE CASCADE,
    CONSTRAINT fk_detail_id_batch    FOREIGN KEY (id_batch)
        REFERENCES batch_stok (id_batch) ON DELETE RESTRICT,
    CONSTRAINT fk_detail_id_kt       FOREIGN KEY (id_kemasan_terbuka)
        REFERENCES kemasan_terbuka (id_kemasan_terbuka) ON DELETE SET NULL
);

CREATE INDEX idx_detail_opname_id_opname ON detail_opname (id_opname);
```

---

## 10. Catatan Integritas Data

### 10.1 Data Bersifat Append-Only untuk Tabel Transaksi

Sesuai KNF-05 (Penyimpanan Data Historis), tabel `stok_masuk`, `stok_keluar`, dan `detail_opname` **tidak** memiliki mekanisme penghapusan dari aplikasi. Koreksi stok dilakukan melalui sesi stock opname yang mencatat histori selisih baru, bukan menimpa data yang sudah ada.

### 10.2 DB Transaction untuk Operasi Multi-Tabel

Sesuai KNF-04, seluruh operasi yang mengubah lebih dari satu tabel wajib dieksekusi dalam satu transaksi database (`sql.Tx`):

| Operasi | Tabel yang Terpengaruh |
|---|---|
| Simpan Stok Masuk | `stok_masuk` + `batch_stok` (INSERT atau UPDATE) |
| Simpan Stok Keluar (Full Use) | `stok_keluar` + `batch_stok` |
| Simpan Stok Keluar (Partial Use, kemasan baru) | `stok_keluar` + `batch_stok` + `kemasan_terbuka` (INSERT) |
| Simpan Stok Keluar (Partial Use, kemasan terbuka) | `stok_keluar` + `kemasan_terbuka` (UPDATE `isi_tersisa`) |
| Selesaikan Stock Opname | `stok_opname` + `detail_opname` + `batch_stok` dan/atau `kemasan_terbuka` |

### 10.3 Nilai Stok Tidak Pernah Dihitung di Frontend

Sesuai KNF-04 AC-NFR04.1, nilai `stok_kemasan`, `total_isi_tersedia`, dan `isi_tersisa` selalu dihitung dan diperbarui oleh backend saat transaksi berlangsung. Frontend hanya menampilkan nilai yang dikembalikan API.

### 10.4 Pembaruan `updated_at`

Kolom `updated_at` diperbarui oleh aplikasi (bukan database trigger) setiap kali baris diubah, untuk menjaga kompatibilitas lintas database dan menghindari kompleksitas trigger.
