# SRS — Overview
## Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic

---

## Daftar Isi

1. [Tujuan Dokumen](#1-tujuan-dokumen)
2. [Lingkup Sistem](#2-lingkup-sistem)
3. [Latar Belakang](#3-latar-belakang)
4. [Permasalahan Sistem Berjalan](#4-permasalahan-sistem-berjalan)
5. [Tujuan Sistem](#5-tujuan-sistem)
6. [Batasan Sistem](#6-batasan-sistem)
7. [Stakeholder](#7-stakeholder)
8. [Arsitektur Sistem](#8-arsitektur-sistem)
9. [Technology Stack](#9-technology-stack)
10. [Struktur Dokumen SRS](#10-struktur-dokumen-srs)
11. [Konvensi & Glosarium](#11-konvensi--glosarium)

---

## 1. Tujuan Dokumen

Dokumen ini merupakan Software Requirements Specification (SRS) untuk Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic. SRS berfungsi sebagai acuan tunggal (*single source of truth*) yang mendefinisikan seluruh kebutuhan fungsional, kebutuhan non-fungsional, arsitektur sistem, kontrak API, skema basis data, serta spesifikasi antarmuka yang harus dipenuhi selama proses perancangan dan implementasi sistem.

Dokumen SRS ini terdiri dari tujuh bagian terpisah:

| File | Cakupan |
|---|---|
| `srs-overview.md` | Gambaran umum, konteks, arsitektur, dan glosarium |
| `srs-fr.md` | Kebutuhan fungsional lengkap per use case |
| `srs-nfr.md` | Kebutuhan non-fungsional dan atribut kualitas sistem |
| `srs-frontend.md` | Spesifikasi lapisan antarmuka (React + Vite) |
| `srs-backend.md` | Spesifikasi lapisan logika bisnis (Go + Fiber) |
| `srs-database.md` | Spesifikasi skema basis data (PostgreSQL) |
| `srs-api.md` | Kontrak REST API lengkap |

---

## 2. Lingkup Sistem

Sistem yang didefinisikan dalam dokumen ini adalah **Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic**, yaitu aplikasi berbasis web yang dirancang untuk menggantikan proses pengelolaan persediaan semi-manual berbasis catatan fisik dan Microsoft Excel di Farmasi Internal Freya Skin Clinic, Sumedang.

**Termasuk dalam lingkup sistem:**
- Autentikasi pengguna tunggal (Admin Farmasi)
- Pengelolaan data master kategori dan produk farmasi
- Pencatatan stok masuk beserta informasi batch dan expired date
- Pencatatan stok keluar dengan mekanisme full use dan partial use
- Penerapan otomatis metode First Expired First Out (FEFO)
- Pengelolaan Beyond Use Date (BUD) untuk produk partial use
- Monitoring stok secara real-time
- Stock opname dengan pencatatan histori selisih
- Laporan stok masuk, keluar, dan persediaan terkini secara periodik

**Tidak termasuk dalam lingkup sistem:**
- Integrasi dengan sistem kasir, rekam medis, atau sistem lain di Freya Skin Clinic
- Pengadaan barang ke supplier secara otomatis
- Pencatatan penyusutan stok
- Audit persediaan eksternal
- Manajemen multi-pengguna atau multi-role

---

## 3. Latar Belakang

Freya Skin Clinic adalah klinik kecantikan yang berlokasi di Sumedang dan beroperasi setiap Senin–Sabtu pukul 10.00–19.00 WIB. Layanan yang tersedia mencakup skin analysis, perawatan laser, tindakan injeksi, serta perawatan kulit wajah dan tubuh lainnya.

Dalam operasionalnya, klinik ini memiliki Farmasi Internal yang mengelola persediaan produk dengan keragaman tinggi, meliputi lima kategori utama: **Skincare, Injectable, Obat, Threadlift, dan Facial IPL Laser**. Produk dikelola dalam berbagai bentuk kemasan (vial, jar, botol, ampoule, syringe, set, pack, box) dengan satuan yang beragam (pcs, cc/ml, gram, IU).

Sistem pencatatan yang berjalan saat ini menggunakan kombinasi **catatan fisik (buku tulis)** dan **Microsoft Excel** dengan satu file per periode bulanan. Pendekatan ini menimbulkan sejumlah keterbatasan operasional yang memerlukan solusi berbasis sistem informasi yang terintegrasi.

---

## 4. Permasalahan Sistem Berjalan

Berdasarkan hasil analisis sistem berjalan, ditemukan lima permasalahan utama:

| Kode | Permasalahan | Dampak |
|---|---|---|
| P-01 | Tidak adanya pencatatan expired date dan batch secara sistematis; FEFO belum terintegrasi | Pengendalian kedaluwarsa bergantung pada pengecekan fisik; risiko produk kadaluwarsa terpakai |
| P-02 | Pencatatan full use & partial use serta perhitungan sisa stok masih manual | Risiko kesalahan pencatatan dan ketidakakuratan data persediaan |
| P-03 | BUD (Beyond Use Date) tidak terdokumentasi dalam sistem | Tidak ada kontrol masa penggunaan produk setelah kemasan dibuka |
| P-04 | Monitoring, stock opname, dan rekap stok manual tanpa pencatatan histori selisih | Proses tidak efisien; tidak ada jejak koreksi data |
| P-05 | Data stok tersimpan dalam file terpisah per periode bulanan | Data tidak terintegrasi; sulit tracking histori persediaan lintas periode |

---

## 5. Tujuan Sistem

Sistem dirancang untuk:

1. Menggantikan pencatatan manual menjadi sistem terintegrasi berbasis web yang dapat diakses melalui browser
2. Menyediakan pencatatan stok masuk yang mencakup informasi batch dan expired date secara sistematis
3. Mendukung pencatatan stok keluar dengan membedakan pola full use dan partial use secara otomatis
4. Mengimplementasikan metode FEFO secara otomatis berdasarkan data expired date yang tercatat di sistem
5. Mencatat dan memantau BUD setiap kemasan produk partial use yang dibuka
6. Menyediakan monitoring stok real-time berdasarkan kategori, batch, dan status kedaluwarsa
7. Mendukung stock opname dengan pencatatan histori koreksi yang terstruktur
8. Menghasilkan laporan stok secara periodik dalam satu sistem yang terintegrasi

---

## 6. Batasan Sistem

| No | Batasan |
|---|---|
| B-01 | Sistem hanya digunakan oleh Farmasi Internal Freya Skin Clinic; tidak mencakup klinik atau farmasi lain |
| B-02 | Sistem memiliki satu peran pengguna yaitu Admin Farmasi |
| B-03 | Metode pengembangan menggunakan Waterfall (Pressman) dengan pemodelan UML |
| B-04 | Proses pengkodean tidak dibahas secara rinci; fokus pada perancangan logika dan arsitektur sistem |
| B-05 | Sistem tidak mencakup fitur pengadaan otomatis ke supplier |
| B-06 | Sistem tidak terintegrasi dengan sistem kasir, rekam medis, atau sistem lain di klinik |
| B-07 | Sistem tidak mencakup pencatatan penyusutan stok dan audit persediaan langsung |
| B-08 | BUD ditetapkan secara fixed selama 28 hari setelah pembukaan kemasan |

---

## 7. Stakeholder

| Stakeholder | Peran | Kepentingan dalam Sistem |
|---|---|---|
| **Admin Farmasi** | Pengguna langsung (aktor utama) | Melakukan seluruh operasi: pencatatan stok masuk/keluar, monitoring, stock opname, laporan |
| **Pengelola Freya Skin Clinic** | Pemilik proses bisnis | Memastikan sistem mendukung operasional farmasi secara efisien dan akurat |
| **Tim Pengembang** | Perancang dan implementator | Membangun sistem sesuai spesifikasi yang telah ditetapkan |
| **Peneliti / Mahasiswa** | Penyusun skripsi | Merancang, mengimplementasikan, dan mengevaluasi sistem sebagai objek penelitian |

---

## 8. Arsitektur Sistem

Sistem dikembangkan menggunakan **arsitektur client-server berbasis web** dengan pemisahan yang tegas antara lapisan antarmuka, lapisan logika bisnis, dan lapisan penyimpanan data. Pemilihan setiap komponen teknologi didasarkan pada karakteristik spesifik kebutuhan yang telah diidentifikasi pada tahap analisis, bukan pada popularitas teknologi secara umum.

### 8.1 Lapisan Antarmuka (Frontend)

Lapisan antarmuka dibangun menggunakan **React dengan Vite** sebagai build tool. React dipilih karena memungkinkan pembangunan antarmuka berbasis komponen yang dapat digunakan ulang, sementara Vite menyediakan proses build dan development server yang ringan sehingga iterasi pengembangan antarmuka yang melibatkan banyak komponen serupa dapat dilakukan secara efisien.

Pendekatan berbasis komponen dipilih karena sistem memiliki beberapa tampilan yang secara struktural serupa namun berbeda konten — halaman kelola kategori, kelola produk, stok masuk, stok keluar, dan laporan memiliki pola tabel, filter, dan formulir yang identik sehingga komponen yang sama dapat digunakan lintas halaman tanpa duplikasi kode.

React juga mendukung pembaruan tampilan secara selektif tanpa memuat ulang seluruh halaman, sehingga perubahan data seperti status batch, status BUD, dan indikator stok kritis pada halaman monitoring dapat diperbarui secara dinamis ketika pengguna melakukan interaksi, tanpa mengganggu tampilan halaman secara keseluruhan.

### 8.2 Lapisan Logika Bisnis (Backend)

Lapisan logika bisnis dibangun menggunakan **Go dengan Fiber** sebagai web framework. Fiber menyediakan routing dan middleware untuk penanganan request API, termasuk middleware autentikasi berbasis token yang mendukung pembatasan akses sesuai hak pengguna yang telah ditentukan.

Pemilihan Go didasarkan pada kebutuhan sistem akan proses yang berjalan di luar alur utama penanganan permintaan pengguna, yaitu:
- **Pemantauan status batch** — pengecekan batch yang telah melewati expired date untuk diperbarui statusnya menjadi `KADALUWARSA` secara otomatis
- **Pemantauan status BUD** — pengecekan kemasan terbuka yang telah melewati batas BUD (28 hari) untuk dinonaktifkan secara otomatis

Go menyediakan model konkurensi berbasis **goroutine dan channel** yang memungkinkan kedua proses pemantauan tersebut berjalan sebagai background worker secara terpisah dari proses penanganan permintaan utama, sehingga validasi data dan respons terhadap pengguna tidak terhambat oleh proses pemantauan yang berjalan di latar belakang.

### 8.3 Lapisan Penyimpanan Data

Lapisan penyimpanan data menggunakan **PostgreSQL** sebagai basis data relasional untuk menyimpan seluruh data operasional secara persisten, mencakup data kategori, produk, batch stok, transaksi stok masuk dan keluar, kemasan terbuka, serta riwayat stock opname.

Pemilihan PostgreSQL didasarkan pada kebutuhan integritas relasional yang tinggi, mengingat data operasional farmasi memiliki banyak relasi antar entitas yang harus terjaga konsistensinya — data stok keluar harus selalu terhubung ke batch yang valid, data kemasan terbuka harus terikat pada batch yang aktif, dan data detail opname harus terhubung ke sesi opname yang bersangkutan.

### 8.4 Diagram Arsitektur

```
┌─────────────────────────────────────────────────────┐
│                   CLIENT (Browser)                   │
│                                                      │
│   ┌─────────────────────────────────────────────┐   │
│   │     React Application (Vite Build)          │   │
│   │  Pages · Components · State · API Client    │   │
│   └─────────────────────────────────────────────┘   │
└─────────────────────┬───────────────────────────────┘
                      │ HTTPS / REST API (JSON)
┌─────────────────────▼───────────────────────────────┐
│                  SERVER (Backend)                    │
│                                                      │
│   ┌─────────────────────────────────────────────┐   │
│   │    Go Application (Fiber Framework)         │   │
│   │                                             │   │
│   │  ┌──────────────┐   ┌─────────────────┐    │   │
│   │  │  HTTP Handler │   │  Background     │    │   │
│   │  │  (Router +   │   │  Workers        │    │   │
│   │  │  Middleware) │   │  (Goroutines)   │    │   │
│   │  └──────┬───────┘   │  - Batch Status │    │   │
│   │         │           │  - BUD Expiry   │    │   │
│   │  ┌──────▼───────┐   └────────┬────────┘    │   │
│   │  │ Service Layer│            │             │   │
│   │  │ (Biz Logic)  │◄───────────┘             │   │
│   │  └──────┬───────┘                          │   │
│   │         │                                  │   │
│   │  ┌──────▼───────┐                          │   │
│   │  │ Repository   │                          │   │
│   │  │ Layer (DB)   │                          │   │
│   │  └──────┬───────┘                          │   │
│   └─────────┼───────────────────────────────────┘   │
└─────────────┼───────────────────────────────────────┘
              │ SQL
┌─────────────▼───────────────────────────────────────┐
│              PostgreSQL Database                     │
│  users · kategori · produk · stok_masuk             │
│  batch_stok · stok_keluar · kemasan_terbuka         │
│  stok_opname · detail_opname                        │
└─────────────────────────────────────────────────────┘
```

---

## 9. Technology Stack

| Layer | Teknologi | Versi (rekomendasi) | Justifikasi |
|---|---|---|---|
| Frontend | React | ^18.x | Komponen reusable, pembaruan UI selektif |
| Frontend Build Tool | Vite | ^5.x | Build ringan, dev server cepat |
| Frontend Routing | React Router | ^6.x | Client-side routing SPA |
| Frontend State | Zustand / React Query | latest | State management ringan + data fetching |
| Frontend UI | Tailwind CSS | ^3.x | Utility-first, konsistensi styling |
| Backend Language | Go | ^1.22 | Goroutine untuk background worker |
| Backend Framework | Fiber | ^2.x | HTTP routing, middleware, performance |
| Backend ORM/Query | sqlx / pgx | latest | Typed query ke PostgreSQL |
| Database | PostgreSQL | ^16.x | Integritas relasional, ACID compliance |
| Authentication | JWT (jose/golang-jwt) | latest | Stateless token-based auth |
| Containerization | Docker + Docker Compose | latest | Portabilitas deployment |

---

## 10. Struktur Dokumen SRS

```
srs-freya-farmasi/
├── srs-overview.md      ← Dokumen ini
├── srs-fr.md            ← Functional Requirements (KF-01 s/d KF-10)
├── srs-nfr.md           ← Non-Functional Requirements (KNF-01 s/d KNF-05+)
├── srs-frontend.md      ← Spesifikasi Frontend (React + Vite)
├── srs-backend.md       ← Spesifikasi Backend (Go + Fiber)
├── srs-database.md      ← Spesifikasi Database (PostgreSQL)
└── srs-api.md           ← Kontrak REST API
```

---

## 11. Konvensi & Glosarium

### 11.1 Konvensi Kode

| Prefix | Domain |
|---|---|
| `KF-xx` | Kebutuhan Fungsional |
| `KNF-xx` | Kebutuhan Non-Fungsional |
| `UC-xx` | Use Case |
| `P-xx` | Permasalahan sistem berjalan |
| `API-xx` | Endpoint API |
| `TBL-xx` | Tabel database |

### 11.2 Glosarium

| Istilah | Definisi |
|---|---|
| **Admin Farmasi** | Satu-satunya aktor dalam sistem; pengguna yang bertanggung jawab atas seluruh pengelolaan stok farmasi |
| **Batch** | Kelompok produk dengan expired date yang sama yang diterima dalam satu atau lebih transaksi penerimaan |
| **BUD (Beyond Use Date)** | Batas waktu penggunaan produk setelah kemasan pertama kali dibuka, ditetapkan 28 hari sejak tanggal pembukaan |
| **FEFO (First Expired First Out)** | Metode pengeluaran stok yang memprioritaskan produk dengan expired date paling dekat |
| **Full Use** | Pola penggunaan produk yang dihabiskan seluruhnya dalam satu kali pemakaian |
| **Partial Use** | Pola penggunaan produk yang hanya digunakan sebagian sehingga menyisakan isi untuk pemakaian berikutnya |
| **Kemasan Terbuka** | Kemasan produk partial use yang telah dibuka dan masih memiliki sisa isi yang dapat digunakan kembali |
| **Stock Opname** | Kegiatan pencocokan data stok dalam sistem dengan kondisi fisik aktual di gudang |
| **Stok Masuk** | Transaksi penerimaan produk dari supplier yang menghasilkan penambahan stok |
| **Stok Keluar** | Transaksi penggunaan produk dalam pelayanan pasien yang menghasilkan pengurangan stok |
| **Status Batch** | Kondisi batch: `AKTIF` (stok tersedia), `HABIS` (stok = 0), `KADALUWARSA` (expired date terlewati) |
| **Status BUD** | Kondisi kemasan terbuka: `AKTIF` (BUD belum terlewati), `KADALUWARSA` (BUD telah terlewati) |
