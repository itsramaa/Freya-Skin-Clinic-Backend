# SRS — Non-Functional Requirements (NFR)
## Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic

---

## Daftar Isi

1. [Tujuan & Cakupan](#1-tujuan--cakupan)
2. [Daftar Kebutuhan Non-Fungsional](#2-daftar-kebutuhan-non-fungsional)
3. [KNF-01 — Usability](#3-knf-01--usability)
4. [KNF-02 — Security](#4-knf-02--security)
5. [KNF-03 — Portability (Berbasis Web)](#5-knf-03--portability-berbasis-web)
6. [KNF-04 — Reliability & Akurasi Data](#6-knf-04--reliability--akurasi-data)
7. [KNF-05 — Maintainability & Penyimpanan Historis](#7-knf-05--maintainability--penyimpanan-historis)
8. [KNF-06 — Performance](#8-knf-06--performance-tambahan)
9. [KNF-07 — Availability Background Worker](#9-knf-07--availability-background-worker)
10. [Matriks Traceability NFR](#10-matriks-traceability-nfr)

---

## 1. Tujuan & Cakupan

Dokumen ini merinci kebutuhan non-fungsional (KNF) sistem, yaitu atribut kualitas yang tidak berkaitan langsung dengan fungsi bisnis utama, tetapi menentukan bagaimana sistem harus berperilaku dari sisi kemudahan pakai, keamanan, ketersediaan platform, akurasi data, dan kemampuan pemeliharaan. KNF-01 s.d. KNF-05 mengikuti penomoran pada tahap analisis kebutuhan (Bab IV); KNF-06 dan KNF-07 ditambahkan sebagai elaborasi teknis yang relevan dengan pemilihan arsitektur Go + React + PostgreSQL agar dapat diuji secara terukur.

Setiap KNF dirinci dengan **metrik yang dapat diukur** (bukan hanya deskripsi kualitatif) agar dapat dijadikan acuan pengujian non-fungsional.

---

## 2. Daftar Kebutuhan Non-Fungsional

| Kode | Nama | Kategori (ISO 25010) |
|---|---|---|
| KNF-01 | Kemudahan Penggunaan | Usability |
| KNF-02 | Keamanan Akses | Security |
| KNF-03 | Berbasis Web | Portability / Compatibility |
| KNF-04 | Keakuratan dan Integritas Data | Reliability / Functional Suitability |
| KNF-05 | Penyimpanan Data Historis | Maintainability / Reliability |
| KNF-06 | Performance | Performance Efficiency |
| KNF-07 | Availability Background Worker | Reliability |

---

## 3. KNF-01 — Usability

**Deskripsi:** Sistem dirancang dengan antarmuka yang sederhana dan mudah digunakan oleh Admin Farmasi tanpa memerlukan pelatihan khusus.

### 3.1 Kebutuhan Detail
- Seluruh form input menggunakan label berbahasa Indonesia yang sesuai istilah operasional farmasi (mis. "Tanggal Kedaluwarsa", bukan "Expired Date" pada UI, kecuali singkatan baku seperti FEFO/BUD).
- Pola tabel, filter, dan formulir dibuat **konsisten lintas halaman** (Kategori, Produk, Stok Masuk, Stok Keluar) menggunakan komponen UI yang sama, sehingga pengguna yang sudah familier dengan satu modul dapat langsung memahami modul lain.
- Setiap aksi destruktif (hapus) **wajib** melalui dialog konfirmasi sebelum eksekusi (sesuai activity diagram Bab IV).
- Pesan error dan sukses sistem ditampilkan dalam bahasa yang sama dengan pesan yang telah dirumuskan pada `srs-fr.md` (mis. "Kategori berhasil ditambahkan.").

### 3.2 Metrik Pengujian
- AC-NFR01.1: Admin Farmasi baru dapat menyelesaikan transaksi stok masuk dasar dalam ≤ 5 menit tanpa bantuan dokumentasi tertulis (uji usability testing).
- AC-NFR01.2: Tidak ada aksi hapus yang dieksekusi tanpa dialog konfirmasi.

---

## 4. KNF-02 — Security

**Deskripsi:** Sistem menyediakan mekanisme autentikasi (login) untuk memastikan akses hanya diberikan kepada pengguna yang berwenang.

### 4.1 Kebutuhan Detail
- Autentikasi menggunakan **token JWT** yang divalidasi oleh middleware Fiber pada setiap request ke endpoint yang dilindungi (lihat `srs-backend.md` § Middleware Autentikasi).
- Password pengguna disimpan dalam bentuk **hash** (bcrypt/argon2), tidak pernah dalam bentuk plain text, sesuai BR-01.1 pada `srs-fr.md`.
- Token memiliki **masa berlaku (expiry)** dan harus ditolak (HTTP 401) apabila telah kedaluwarsa atau tidak valid.
- Karena sistem hanya memiliki satu peran pengguna (Admin Farmasi), tidak diperlukan mekanisme role-based access control (RBAC) granular; namun seluruh endpoint API (kecuali `/login`) tetap **wajib** diproteksi token.
- Koneksi antara client dan server direkomendasikan menggunakan HTTPS pada lingkungan produksi.

### 4.2 Metrik Pengujian
- AC-NFR02.1: Request ke endpoint terproteksi tanpa token atau dengan token tidak valid selalu mengembalikan HTTP 401.
- AC-NFR02.2: Password tidak pernah muncul dalam bentuk plain text pada log aplikasi maupun respons API.
- AC-NFR02.3: Token kedaluwarsa otomatis ditolak oleh middleware tanpa memerlukan validasi tambahan di level handler.

---

## 5. KNF-03 — Portability (Berbasis Web)

**Deskripsi:** Sistem dirancang berbasis web sehingga dapat diakses melalui browser tanpa instalasi aplikasi tambahan.

### 5.1 Kebutuhan Detail
- Frontend dibangun sebagai **Single Page Application (SPA)** menggunakan React + Vite yang di-build menjadi aset statis (HTML/CSS/JS) dan disajikan melalui browser.
- Sistem harus dapat diakses minimal pada browser berbasis Chromium (Chrome, Edge) versi dua tahun terakhir, mengingat lingkungan operasional klinik umumnya menggunakan PC/laptop standar.
- Tidak ada dependensi pada plugin browser tambahan (Flash, Java Applet, ActiveX, dsb).
- Resolusi minimum yang didukung: 1280×720 (desktop/laptop), karena admin farmasi mengoperasikan sistem dari komputer kerja gudang/farmasi, bukan perangkat mobile.

### 5.2 Metrik Pengujian
- AC-NFR03.1: Seluruh fungsi utama (stok masuk, stok keluar, monitoring, opname, laporan) dapat dijalankan penuh tanpa instalasi software tambahan selain browser.

---

## 6. KNF-04 — Reliability & Akurasi Data

**Deskripsi:** Sistem mampu melakukan perhitungan stok, penerapan FEFO, dan pengelolaan BUD secara otomatis untuk memastikan keakuratan dan konsistensi data.

### 6.1 Kebutuhan Detail
- Seluruh perhitungan stok (penambahan saat stok masuk, pengurangan saat stok keluar, penyesuaian saat opname) **dihitung oleh server**, bukan oleh client, untuk mencegah manipulasi atau kesalahan perhitungan di sisi browser.
- Operasi yang mengubah lebih dari satu tabel sekaligus (mis. stok keluar yang memengaruhi `batch_stok` dan `kemasan_terbuka`) **wajib** dibungkus dalam **transaksi database (DB transaction)** agar tidak terjadi data parsial apabila terjadi kegagalan di tengah proses.
- Pemilihan batch FEFO ditentukan oleh logika backend (`getBatchPrioritasFEFO()`), bukan oleh pilihan manual pengguna, sesuai relasi `<<include>>` pada `srs-fr.md` KF-06.
- Validasi data input (tanggal, jumlah, expired date) dilakukan di sisi backend sebagai validasi akhir, meskipun frontend juga melakukan validasi awal — backend tidak boleh mempercayai validasi client sepenuhnya.

### 6.2 Metrik Pengujian
- AC-NFR04.1: Tidak ditemukan kondisi *race condition* yang menyebabkan stok negatif pada pengujian transaksi stok keluar konkuren terhadap batch yang sama.
- AC-NFR04.2: Apabila salah satu langkah dalam transaksi stok keluar gagal (mis. gagal update `kemasan_terbuka`), seluruh perubahan pada transaksi tersebut di-*rollback*.

---

## 7. KNF-05 — Maintainability & Penyimpanan Historis

**Deskripsi:** Sistem menyimpan seluruh riwayat transaksi secara terpusat untuk mendukung pelacakan data, monitoring, dan kebutuhan pelaporan.

### 7.1 Kebutuhan Detail
- Seluruh transaksi (`stok_masuk`, `stok_keluar`, `stok_opname`, `detail_opname`) disimpan secara permanen pada satu basis data terpusat (PostgreSQL) — menjawab permasalahan P-05 (data tidak lagi tersebar dalam file Excel per periode bulanan).
- Tidak ada mekanisme penghapusan data transaksi historis dari aplikasi (data bersifat *append-only* untuk tabel transaksi); koreksi dilakukan melalui mekanisme stock opname yang tetap mencatat jejak (`detail_opname`), bukan dengan menghapus data lama.
- Struktur kode (backend) mengikuti pemisahan **handler → service → repository** (lihat `srs-backend.md`) agar perubahan logika bisnis di masa depan tidak memengaruhi lapisan lain secara langsung.

### 7.2 Metrik Pengujian
- AC-NFR05.1: Data transaksi bulan-bulan sebelumnya tetap dapat diakses dan dihitung ulang melalui modul Laporan Stok tanpa bergantung pada file eksternal.
- AC-NFR05.2: Setiap koreksi stok melalui Stock Opname menghasilkan baris baru pada `detail_opname`, tidak menimpa (overwrite) data transaksi sebelumnya.

---

## 8. KNF-06 — Performance (Tambahan)

> KNF ini ditambahkan sebagai elaborasi teknis yang relevan dengan pemilihan Go + Fiber serta beban kerja query relasional pada PostgreSQL; tidak terdapat pada tabel KNF asli hasil analisis (Bab IV), namun konsisten dengan kebutuhan KF-08 (Monitoring Stok real-time).

### 8.1 Kebutuhan Detail
- Endpoint Monitoring Stok yang membaca data produk, batch, dan kemasan terbuka secara bersamaan harus memanfaatkan index database (lihat `srs-database.md` § Index) agar waktu respons tetap rendah meski jumlah batch bertambah dari waktu ke waktu.
- Operasi background worker (pengecekan status batch & BUD) **tidak boleh memblokir** thread yang menangani request HTTP, sesuai sifat goroutine pada Go yang berjalan independen dari request-response cycle.

### 8.2 Metrik Pengujian (indikatif, disesuaikan saat UAT)
- AC-NFR06.1: Waktu respons endpoint Monitoring Stok dan Dashboard berada pada rentang yang wajar untuk operasional harian satu farmasi internal (skala data: ratusan produk, ribuan baris transaksi per tahun) — ambang batas presisi ditetapkan pada tahap pengujian non-fungsional (Bab 4.7 Pengujian Sistem), bukan diasumsikan di sini.

---

## 9. KNF-07 — Availability Background Worker

> KNF ini ditambahkan karena kebutuhan KF-06 (FEFO otomatis berbasis status batch) dan KF-07 (BUD otomatis) bergantung pada proses berkala yang berjalan **di luar** alur interaksi pengguna, sebagaimana dijelaskan pada Bab IV bagian arsitektur Go yang menggunakan goroutine.

### 9.1 Kebutuhan Detail
- Background worker pemeriksa status batch (expired) dan status BUD harus berjalan secara periodik selama aplikasi backend aktif, tanpa memerlukan trigger manual dari pengguna.
- Apabila backend di-restart, background worker harus otomatis berjalan kembali bersamaan dengan proses startup aplikasi (tidak memerlukan proses terpisah/cron eksternal pada versi awal sistem).
- Kegagalan satu siklus pemeriksaan worker (mis. galat sementara koneksi database) tidak boleh menghentikan keseluruhan proses aplikasi (worker harus *fault-tolerant*: log error, lanjut ke siklus berikutnya).

### 9.2 Metrik Pengujian
- AC-NFR07.1: Status batch yang melewati `expiredDate` berubah menjadi `KADALUWARSA` dalam satu siklus worker tanpa interaksi user.
- AC-NFR07.2: Status `KemasanTerbuka` yang BUD-nya lewat berubah menjadi `KADALUWARSA` dalam satu siklus worker tanpa interaksi user.
- AC-NFR07.3: Restart aplikasi backend tidak menyebabkan worker berhenti permanen.

---

## 10. Matriks Traceability NFR

| KNF | Kategori | KF/UC Terkait | Komponen Arsitektur Terkait |
|---|---|---|---|
| KNF-01 | Usability | Seluruh KF | Frontend (React, komponen reusable) |
| KNF-02 | Security | KF-01 | Backend (middleware JWT) |
| KNF-03 | Portability | Seluruh KF | Frontend (SPA berbasis browser) |
| KNF-04 | Reliability | KF-04, KF-05, KF-06, KF-09 | Backend (service layer, DB transaction) |
| KNF-05 | Maintainability | KF-04, KF-05, KF-09, KF-10 | Database (PostgreSQL terpusat) |
| KNF-06 | Performance | KF-08 | Database (index), Backend (goroutine) |
| KNF-07 | Reliability | KF-06, KF-07 | Backend (background worker) |
