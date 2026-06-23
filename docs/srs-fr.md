# SRS — Functional Requirements (FR)
## Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic

---

## Daftar Isi

1. [Tujuan & Cakupan](#1-tujuan--cakupan)
2. [Daftar Use Case](#2-daftar-use-case)
3. [Daftar Kebutuhan Fungsional](#3-daftar-kebutuhan-fungsional)
4. [KF-01 — Autentikasi Pengguna](#4-kf-01--autentikasi-pengguna)
5. [KF-02 — Kelola Data Kategori](#5-kf-02--kelola-data-kategori)
6. [KF-03 — Kelola Data Produk](#6-kf-03--kelola-data-produk)
7. [KF-04 — Kelola Stok Masuk](#7-kf-04--kelola-stok-masuk)
8. [KF-05 — Kelola Stok Keluar](#8-kf-05--kelola-stok-keluar)
9. [KF-06 — Penerapan FEFO](#9-kf-06--penerapan-fefo)
10. [KF-07 — Kelola BUD](#10-kf-07--kelola-bud)
11. [KF-08 — Monitoring Stok](#11-kf-08--monitoring-stok)
12. [KF-09 — Stock Opname](#12-kf-09--stock-opname)
13. [KF-10 — Laporan Stok](#13-kf-10--laporan-stok)
14. [Matriks Traceability](#14-matriks-traceability)

---

## 1. Tujuan & Cakupan

Dokumen ini merinci seluruh kebutuhan fungsional (KF) Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic. Setiap KF diturunkan langsung dari use case diagram dan activity diagram pada hasil analisis dan perancangan sistem (Bab IV), serta dirinci hingga level alur, validasi, pesan sistem, dan kondisi akhir (post-condition) yang dapat dijadikan acuan implementasi maupun pengujian (black box testing).

Penomoran KF mengikuti tabel kebutuhan fungsional yang telah ditetapkan pada tahap analisis kebutuhan (KF-01 s.d. KF-10) dan dipertahankan agar konsisten dengan `srs-overview.md`.

---

## 2. Daftar Use Case

> **Catatan koreksi penomoran:** Pada dokumen hasil analisis, kode UC-07 digunakan dua kali (untuk Stok Opname dan Laporan Stok). Pada SRS ini, penomoran dikoreksi menjadi UC-07 (Stok Opname) dan UC-08 (Laporan Stok) agar tidak ambigu. Tidak ada perubahan substansi fungsi.

| Kode | Use Case | Aktor | Relasi |
|---|---|---|---|
| UC-01 | Login | Admin Farmasi | — |
| UC-02 | Kelola Data Kategori | Admin Farmasi | — |
| UC-03 | Kelola Data Produk | Admin Farmasi | — |
| UC-04 | Kelola Stok Masuk | Admin Farmasi | — |
| UC-05 | Kelola Stok Keluar | Admin Farmasi | `<<include>>` Penerapan FEFO; `<<extend>>` Kelola BUD |
| UC-06 | Monitoring Stok | Admin Farmasi | — |
| UC-07 | Stok Opname | Admin Farmasi | — |
| UC-08 | Laporan Stok | Admin Farmasi | — |

---

## 3. Daftar Kebutuhan Fungsional

| Kode | Nama | Use Case Terkait |
|---|---|---|
| KF-01 | Autentikasi pengguna | UC-01 |
| KF-02 | Kelola data kategori | UC-02 |
| KF-03 | Kelola data produk | UC-03 |
| KF-04 | Kelola stok masuk | UC-04 |
| KF-05 | Kelola stok keluar | UC-05 |
| KF-06 | Penerapan FEFO | UC-05 (include) |
| KF-07 | Kelola BUD | UC-05 (extend) |
| KF-08 | Monitoring stok | UC-06 |
| KF-09 | Stock opname | UC-07 |
| KF-10 | Laporan stok | UC-08 |

---

## 4. KF-01 — Autentikasi Pengguna

**Use Case:** UC-01 Login
**Aktor:** Admin Farmasi
**Deskripsi:** Sistem menyediakan mekanisme login untuk membatasi akses hanya kepada pengguna yang memiliki kredensial valid, serta memaksa penggantian password pada login pertama menggunakan password default.

### 4.1 Pre-condition
- Pengguna memiliki akun yang telah terdaftar di basis data (`users`).

### 4.2 Main Flow
1. User membuka sistem melalui browser; sistem menampilkan halaman login.
2. User menginput `username` dan `password`, kemudian menekan tombol **Login**.
3. Sistem memvalidasi kredensial terhadap data akun yang tersimpan di basis data.
4. Jika kredensial valid, sistem memeriksa status password (`is_default_password`):
   - Jika password masih default → sistem menampilkan halaman ganti password (paksa). User menginput password baru, sistem menyimpan password baru (dengan hashing), lalu menampilkan dashboard.
   - Jika password bukan default → sistem langsung menampilkan dashboard.
5. Proses selesai; user dapat mengakses seluruh menu sistem sesuai sesi aktif.

### 4.3 Alternate Flow
- **A1 — Kredensial tidak valid:** Sistem menampilkan pesan **"Kredensial tidak valid."** dan mengembalikan user ke halaman login untuk mencoba kembali.

### 4.4 Business Rules
- BR-01.1: Password disimpan dalam bentuk hash, tidak pernah dalam bentuk plain text.
- BR-01.2: Login pertama kali (password default) **wajib** diikuti penggantian password sebelum dapat mengakses menu lain.
- BR-01.3: Sesi pengguna direpresentasikan oleh token (lihat `srs-backend.md` § Autentikasi) yang memiliki masa berlaku (expiry).

### 4.5 Post-condition
- Sesi aktif terbentuk (token diterbitkan) dan user diarahkan ke dashboard.

### 4.6 Acceptance Criteria
- AC-01.1: Login dengan kredensial salah menampilkan pesan kesalahan dan tidak menerbitkan token.
- AC-01.2: Login pertama dengan password default mengarahkan ke halaman ganti password sebelum dashboard dapat diakses.
- AC-01.3: Login dengan password yang sudah diganti langsung menampilkan dashboard.

---

## 5. KF-02 — Kelola Data Kategori

**Use Case:** UC-02 Kelola Data Kategori
**Aktor:** Admin Farmasi
**Deskripsi:** Sistem mengelola data kategori sebagai data master pengelompokan produk farmasi (CRUD: Tambah, Ubah, Hapus), termasuk validasi duplikasi nama dan validasi ketergantungan produk sebelum penghapusan.

### 5.1 Tampilan Daftar
Daftar kategori menampilkan kolom: **Kode Kategori, Nama Kategori, Jumlah Produk Terkait**, serta kolom **Aksi** (Tambah, Ubah, Hapus).

### 5.2 Sub-fungsi: Tambah Kategori

**Main Flow**
1. User memilih menu **Tambah Kategori** dan mengisi `namaKategori`.
2. Sistem menampilkan notifikasi konfirmasi sebelum data disimpan.
3. Jika user memilih **Batal** → proses dibatalkan, tidak ada perubahan data.
4. Jika user memilih **Lanjut** → sistem memeriksa duplikasi nama kategori.

**Alternate Flow**
- **A1 — Duplikasi ditemukan:** Sistem menampilkan pesan **"Nama kategori sudah terdaftar dalam sistem."** Data tidak disimpan.
- **A2 — Tidak ada duplikasi:** Sistem menyimpan kategori baru dan menampilkan pesan **"Kategori berhasil ditambahkan."**

### 5.3 Sub-fungsi: Ubah Kategori

**Main Flow**
1. User memilih kategori yang akan diubah; sistem menampilkan formulir terisi data kategori tersebut.
2. User mengubah `namaKategori` dan menekan **Simpan**.
3. Sistem menampilkan notifikasi konfirmasi.
4. Jika **Batal** → perubahan tidak disimpan.
5. Jika **Lanjut** → sistem memvalidasi duplikasi nama kategori (terhadap kategori lain, bukan dirinya sendiri).

**Alternate Flow**
- **A1 — Duplikasi ditemukan:** Pesan **"Nama kategori sudah terdaftar dalam sistem."** Perubahan tidak disimpan.
- **A2 — Tidak ada duplikasi:** Sistem menyimpan perubahan dan menampilkan pesan **"Kategori berhasil diperbarui."**

### 5.4 Sub-fungsi: Hapus Kategori

**Main Flow**
1. User memilih kategori yang akan dihapus.
2. Sistem memeriksa jumlah produk yang masih terkait dengan kategori tersebut.

**Alternate Flow**
- **A1 — Masih memiliki produk terkait (jumlah > 0):** Sistem menolak proses dan menampilkan pesan **"Kategori tidak dapat dihapus karena masih memiliki produk terkait."**
- **A2 — Tidak ada produk terkait (jumlah = 0):** Sistem menampilkan notifikasi konfirmasi penghapusan.
  - Jika **Batal** → proses dibatalkan.
  - Jika **Lanjut** → sistem menghapus data kategori dan menampilkan pesan **"Kategori berhasil dihapus."**

### 5.5 Business Rules
- BR-02.1: Nama kategori bersifat unik (case-insensitive direkomendasikan agar "Skincare" dan "skincare" dianggap sama).
- BR-02.2: Kategori tidak dapat dihapus selama masih direferensikan oleh minimal satu produk (integrity guard, lihat `srs-database.md`).

### 5.6 Acceptance Criteria
- AC-02.1: Penambahan kategori dengan nama yang sudah ada ditolak dengan pesan yang sesuai.
- AC-02.2: Penghapusan kategori yang masih memiliki produk terkait ditolak.
- AC-02.3: Operasi tambah/ubah/hapus yang valid memperbarui daftar kategori secara konsisten.

---

## 6. KF-03 — Kelola Data Produk

**Use Case:** UC-03 Kelola Data Produk
**Aktor:** Admin Farmasi
**Deskripsi:** Sistem mengelola data produk sebagai data master acuan seluruh transaksi stok (stok masuk, stok keluar, monitoring, stock opname, laporan). Atribut kunci: kategori, bentuk kemasan, satuan isi, isi per kemasan, dan pola penggunaan (`FULL_USE` / `PARTIAL_USE`).

### 6.1 Tampilan Daftar
Kolom: **Kode Produk, Nama Produk, Kategori, Bentuk Kemasan, Satuan Isi, Isi per Kemasan, Pola Penggunaan, Stok Kemasan, Total Isi Tersedia**, serta **Aksi**.

### 6.2 Sub-fungsi: Tambah Produk

**Main Flow**
1. User mengisi data produk yang diperlukan (`namaProduk`, `idKategori`, `bentukKemasan`, `satuanIsi`, `isiPerKemasan`, `polaPenggunaan`).
2. User menekan **Simpan**; sistem menampilkan notifikasi konfirmasi.
3. Jika **Batal** → proses dibatalkan.
4. Jika **Lanjut** → sistem memvalidasi kelengkapan data, menghasilkan **kode produk otomatis**, menyimpan data, dan menampilkan pesan **"Produk berhasil ditambahkan."**

### 6.3 Sub-fungsi: Ubah Produk

**Main Flow**
1. User memilih produk yang akan diubah; sistem menampilkan formulir terisi.
2. User melakukan perubahan; sistem menampilkan notifikasi konfirmasi.
3. Jika **Batal** → perubahan tidak disimpan.
4. Jika **Lanjut** → sistem memvalidasi data, menyimpan perubahan, dan menampilkan pesan **"Produk berhasil diperbarui."**

> **Catatan desain:** `polaPenggunaan` sebaiknya dikunci (tidak dapat diubah) apabila produk sudah memiliki transaksi stok masuk/keluar, untuk menjaga konsistensi histori batch dan kemasan terbuka yang telah terbentuk berdasarkan pola tersebut.

### 6.4 Sub-fungsi: Hapus Produk

**Main Flow**
1. User memilih produk yang akan dihapus.
2. Sistem memeriksa apakah produk masih memiliki **stok aktif** atau **riwayat transaksi**.

**Alternate Flow**
- **A1 — Masih memiliki stok aktif:** Pesan **"Produk tidak dapat dihapus karena masih memiliki stok aktif."**
- **A2 — Memiliki riwayat transaksi (stok = 0 tetapi pernah bertransaksi):** Pesan **"Produk tidak dapat dihapus karena memiliki riwayat transaksi."**
- **A3 — Tidak memiliki stok aktif maupun riwayat transaksi:** Sistem menampilkan notifikasi konfirmasi.
  - Jika **Batal** → proses dibatalkan.
  - Jika **Lanjut** → sistem menghapus data produk dan menampilkan pesan **"Data produk berhasil dihapus."**

### 6.5 Business Rules
- BR-03.1: Kode produk dihasilkan otomatis oleh sistem (format direkomendasikan: `PRD-{idKategori}-{sequence}`, didetailkan di `srs-backend.md`).
- BR-03.2: `isiPerKemasan` wajib diisi untuk produk `PARTIAL_USE` karena menjadi basis perhitungan sisa isi kemasan terbuka.
- BR-03.3: Pengecekan stok aktif dilakukan terhadap tabel `batch_stok` (status `AKTIF`); pengecekan riwayat transaksi dilakukan terhadap `stok_masuk` dan `stok_keluar`.

### 6.6 Acceptance Criteria
- AC-03.1: Produk dengan data tidak lengkap tidak dapat disimpan.
- AC-03.2: Produk yang memiliki stok aktif tidak dapat dihapus.
- AC-03.3: Kode produk yang dihasilkan bersifat unik di seluruh sistem.

---

## 7. KF-04 — Kelola Stok Masuk

**Use Case:** UC-04 Kelola Stok Masuk
**Aktor:** Admin Farmasi
**Deskripsi:** Sistem mencatat penerimaan produk dari supplier, mendukung penerapan FEFO dengan mencatat tanggal kedaluwarsa, dan menghasilkan kode batch secara otomatis.

### 7.1 Main Flow
1. User mengakses menu Kelola Stok Masuk dan memilih **Tambah Data Penerimaan**.
2. Sistem menampilkan formulir: `tanggalPenerimaan`, `produk`, `expiredDate`, `jumlahKemasan`.
3. Setelah produk dipilih, sistem menampilkan informasi produk otomatis dari data master (kategori, bentuk kemasan, satuan isi, kapasitas isi per kemasan, pola penggunaan).
4. User mengisi seluruh data dan menekan **Simpan**.
5. Sistem menampilkan notifikasi konfirmasi: **"Apakah Anda yakin ingin menyimpan data stok masuk ini?"**
6. Jika **Lanjut**, sistem memvalidasi data input.

### 7.2 Alternate Flow — Validasi Gagal
Sistem menolak penyimpanan dan menampilkan pesan kesalahan apabila:
- Field wajib belum diisi.
- `tanggalPenerimaan` melebihi tanggal saat ini.
- `expiredDate` ≤ `tanggalPenerimaan`.
- `jumlahKemasan` ≤ 0.

### 7.3 Logika Pembentukan Batch (Business Rule Kunci)
- BR-04.1: Jika kombinasi `idProduk` + `expiredDate` **sudah memiliki batch aktif**, maka `jumlahKemasan` dan `totalIsiMasuk` ditambahkan ke batch yang sudah ada (`tambahStok()`), **bukan** membentuk batch baru.
- BR-04.2: Jika kombinasi tersebut belum ada, sistem membentuk batch baru melalui `generateKodeBatch()` dengan status awal `AKTIF`.
- BR-04.3: `totalIsiMasuk = jumlahKemasan × isiPerKemasan` (mengacu pada `hitungTotalIsi()`).
- BR-04.4: Relasi `StokMasuk → BatchStok` bersifat **1 : N** — satu batch dapat terbentuk dari beberapa transaksi penerimaan selama produk dan `expiredDate` sama, demi menghindari duplikasi batch dan tetap konsisten dengan logika FEFO.

### 7.4 Post-condition
- Data `stok_masuk` tersimpan; `batch_stok` terbentuk atau bertambah jumlahnya; total stok produk diperbarui.

### 7.5 Acceptance Criteria
- AC-04.1: Penerimaan dengan expired date yang lebih awal dari tanggal penerimaan ditolak.
- AC-04.2: Dua transaksi penerimaan dengan produk dan expired date identik menghasilkan satu batch dengan stok terakumulasi, bukan dua batch berbeda.
- AC-04.3: Stok produk pada tampilan data master bertambah secara otomatis setelah stok masuk tersimpan.

---

## 8. KF-05 — Kelola Stok Keluar

**Use Case:** UC-05 Kelola Stok Keluar
**Aktor:** Admin Farmasi
**Deskripsi:** Sistem mencatat penggunaan produk berdasarkan pelayanan pasien, menerapkan FEFO secara otomatis (KF-06, include), serta menangani percabangan **Full Use** dan **Partial Use** termasuk pengelolaan BUD (KF-07, extend) untuk Partial Use.

### 8.1 Main Flow
1. User mengakses menu Stok Keluar dan memilih **Tambah Penggunaan**.
2. User menginput `tanggalPenggunaan` dan memilih produk.
3. Sistem menampilkan informasi produk dan menjalankan mekanisme **FEFO** (`getBatchPrioritasFEFO()`) untuk menentukan batch yang akan digunakan.
4. Sistem memeriksa `polaPenggunaan` produk dan bercabang sesuai sub-fungsi berikut.

### 8.2 Sub-fungsi: Full Use (`prosesFullUse()`)
1. User menginput `jumlahKemasanDipakai`.
2. Sistem memvalidasi ketersediaan stok pada batch prioritas.
3. Jika stok mencukupi: sistem mengurangi stok batch (`kurangiStok()`), memperbarui total stok produk, menyimpan riwayat transaksi (`simpanStokKeluar()`), dan menampilkan pesan **"Data penggunaan berhasil disimpan."**

Contoh sesuai hasil analisis: pengeluaran 1 sunscreen → `jumlahKemasanDipakai = 1`, sisa = 0.

### 8.3 Sub-fungsi: Partial Use (`prosesPartialUse()`)
1. Sistem memeriksa keberadaan **kemasan terbuka aktif** (`getKemasanTerbukaAktif()`) pada batch prioritas.
2. **Jika tidak ada kemasan terbuka aktif:**
   - Sistem membuka kemasan baru, menetapkan BUD otomatis 28 hari (`tetapkanBUD()` — lihat KF-07), menghitung sisa isi kemasan, memperbarui stok, dan menyimpan data kemasan terbuka.
3. **Jika ada kemasan terbuka aktif:**
   - Sistem memeriksa status BUD (`cekStatusBUD()`).
   - Jika **BUD masih berlaku**: sistem menggunakan isi kemasan terbuka sesuai `jumlahIsiDipakai` dan memperbarui `isiTersisa` (`kurangiIsi()`).
   - Jika **BUD telah terlewati**: sistem menonaktifkan kemasan terbuka tersebut (`nonaktifkan()`) dan melanjutkan penggunaan sesuai aturan FEFO (membuka kemasan baru pada batch prioritas berikutnya, mengulang langkah 2).
4. Setelah seluruh proses selesai: sistem memperbarui stok produk, menyimpan riwayat penggunaan, dan menampilkan pesan **"Data penggunaan berhasil disimpan."**

Contoh sesuai hasil analisis: penggunaan 50 IU Botox Bionex dari 1 vial berisi 100 IU → `jumlahIsiDipakai = 50 IU`, `isiTersisa = 50 IU`.

### 8.4 Business Rules
- BR-05.1: Pemilihan batch **selalu** mengikuti urutan FEFO; user tidak dapat memilih batch secara manual (lihat KF-06).
- BR-05.2: Pencatatan stok keluar wajib tervalidasi terhadap ketersediaan stok pada batch prioritas sebelum disimpan.
- BR-05.3: Setiap transaksi stok keluar tercatat dengan `idUser` yang melakukan input, untuk keperluan audit trail (lihat KNF-05).

### 8.5 Acceptance Criteria
- AC-05.1: Penggunaan Full Use mengurangi `stokKemasan` batch sesuai jumlah yang diinput.
- AC-05.2: Penggunaan Partial Use pada kemasan yang belum pernah dibuka membentuk satu baris `kemasan_terbuka` baru dengan BUD terhitung otomatis.
- AC-05.3: Penggunaan Partial Use pada kemasan terbuka yang BUD-nya sudah lewat tidak menggunakan isi kemasan tersebut, melainkan membuka kemasan baru.

---

## 9. KF-06 — Penerapan FEFO

**Relasi:** `<<include>>` terhadap KF-05 (Kelola Stok Keluar)
**Deskripsi:** Penerapan FEFO merupakan bagian wajib yang **selalu** dijalankan setiap kali proses Kelola Stok Keluar dieksekusi. Sistem secara otomatis mengidentifikasi dan memprioritaskan batch produk dengan `expiredDate` paling dekat untuk dikeluarkan terlebih dahulu, tanpa bergantung pada keputusan pengguna. Relasi ini merupakan solusi atas permasalahan P-01, di mana penerapan FEFO pada sistem berjalan masih dilakukan secara manual melalui pengecekan fisik.

### 9.1 Algoritma `getBatchPrioritasFEFO()`

```
INPUT: idProduk
1. Ambil seluruh batch dengan idProduk = INPUT DAN statusBatch = 'AKTIF'
2. Urutkan batch berdasarkan expiredDate ASCENDING
3. RETURN batch pertama pada urutan (expiredDate paling dekat)
```

> **Catatan desain (elaborasi rancangan):** Hasil analisis tidak merinci skenario ketika stok pada batch prioritas tidak mencukupi jumlah yang diminta dalam satu transaksi Full Use. Sebagai kelengkapan rancangan, sistem direkomendasikan menerapkan **pemotongan stok lintas-batch** secara berurutan sesuai prioritas FEFO: apabila batch prioritas tidak mencukupi, sisa kebutuhan diambil dari batch berikutnya pada urutan `expiredDate` terdekat, dan transaksi `stok_keluar` dapat menyimpan rincian per batch yang terlibat. Implementasi detail dirinci pada `srs-backend.md`.

### 9.2 Business Rules
- BR-06.1: Batch dengan status `HABIS` atau `KADALUWARSA` tidak masuk dalam kandidat prioritas FEFO.
- BR-06.2: Penerapan FEFO berjalan untuk **seluruh** transaksi stok keluar, baik Full Use maupun Partial Use.

### 9.3 Acceptance Criteria
- AC-06.1: Sistem selalu memilih batch dengan `expiredDate` terdekat di antara batch berstatus `AKTIF`.
- AC-06.2: User tidak diberikan opsi pemilihan batch secara manual pada form stok keluar.

---

## 10. KF-07 — Kelola BUD

**Relasi:** `<<extend>>` terhadap KF-05 (Kelola Stok Keluar)
**Deskripsi:** Kelola BUD memperluas Kelola Stok Keluar secara kondisional. Use case ini hanya diaktifkan apabila produk yang dicatat bersifat **Partial Use** dan kemasannya **baru pertama kali dibuka**; sistem kemudian secara otomatis menetapkan BUD selama 28 hari ke depan. Apabila produk bersifat Full Use, Kelola BUD tidak dijalankan. Relasi ini merupakan solusi atas permasalahan P-03.

### 10.1 Main Flow
1. Kondisi pemicu: `polaPenggunaan = PARTIAL_USE` DAN tidak ditemukan kemasan terbuka aktif pada batch prioritas.
2. Sistem menjalankan `tetapkanBUD()`: `bud = tanggalDibuka + 28 hari`.
3. Sistem menyimpan baris baru pada `kemasan_terbuka` dengan `statusBUD = AKTIF`.

### 10.2 Pemantauan BUD secara berkala (background worker)
Selain dipicu interaktif saat stok keluar, status BUD juga dipantau oleh **background worker** (lihat `srs-backend.md` § Background Worker) yang berjalan periodik untuk memeriksa kemasan terbuka yang BUD-nya telah lewat namun belum digunakan kembali, lalu menjalankan `nonaktifkan()` agar `statusBUD` selalu mencerminkan kondisi terkini tanpa menunggu transaksi stok keluar berikutnya.

### 10.3 Business Rules
- BR-07.1: BUD bersifat fixed 28 hari sejak `tanggalDibuka`, sesuai Batasan Sistem B-08.
- BR-07.2: Satu batch hanya dapat memiliki **satu** kemasan terbuka aktif pada satu waktu (relasi `BatchStok → KemasanTerbuka` bersifat 1 : 0..1).
- BR-07.3: Kemasan terbuka yang `statusBUD = KADALUWARSA` tidak dapat digunakan kembali meskipun masih memiliki `isiTersisa > 0`.

### 10.4 Acceptance Criteria
- AC-07.1: Produk Full Use tidak pernah menghasilkan baris `kemasan_terbuka`.
- AC-07.2: BUD yang ditampilkan ke user selalu sama dengan `tanggalDibuka + 28 hari`.
- AC-07.3: Kemasan terbuka yang BUD-nya lewat berubah status menjadi `KADALUWARSA` baik melalui interaksi user maupun melalui proses background worker.

---

## 11. KF-08 — Monitoring Stok

**Use Case:** UC-06 Monitoring Stok
**Aktor:** Admin Farmasi
**Deskripsi:** Sistem menampilkan kondisi stok secara real-time berdasarkan data transaksi yang tersimpan. User hanya membaca informasi (read-only); seluruh nilai dihitung otomatis oleh sistem.

### 11.1 Main Flow
1. User mengakses menu Monitoring Stok.
2. Sistem membaca seluruh data produk, data batch, dan data kemasan terbuka aktif.
3. Sistem menghitung jumlah stok tersedia, status `expiredDate` setiap batch, dan status BUD setiap kemasan terbuka.
4. Sistem menampilkan ringkasan kondisi stok seluruh produk.
5. User dapat memfilter berdasarkan: **kategori produk, status expired date, status BUD, nama produk**.
6. User dapat memilih satu produk untuk melihat rincian per batch, termasuk status `expiredDate`.
7. Apabila batch yang dipilih memiliki kemasan terbuka, sistem menampilkan detail: `tanggalDibuka`, `bud`, `isiTersisa`, `statusBUD`.

### 11.2 Business Rules
- BR-08.1: Status expired date ditampilkan dalam kategori indikator (direkomendasikan): `AMAN` (> 30 hari), `MENDEKATI` (≤ 30 hari), `KADALUWARSA` (terlewati) — ambang batas dapat dikonfigurasi pada tahap implementasi.
- BR-08.2: Data yang ditampilkan bersifat read-only; tidak ada aksi tambah/ubah/hapus pada modul ini.

### 11.3 Acceptance Criteria
- AC-08.1: Filter kombinasi (kategori + status expired + status BUD) menghasilkan data yang konsisten dengan kondisi aktual basis data.
- AC-08.2: Drill-down ke detail batch menampilkan seluruh batch aktif maupun tidak aktif milik produk terkait.

---

## 12. KF-09 — Stock Opname

**Use Case:** UC-07 Stok Opname
**Aktor:** Admin Farmasi
**Deskripsi:** Sistem mendukung pencocokan data stok sistem dengan kondisi fisik gudang, beserta pencatatan histori selisih sebagai jejak koreksi — menjawab permasalahan P-04.

### 12.1 Main Flow
1. User mengakses menu Stok Opname dan memilih **Mulai Stok Opname Baru**.
2. Sistem membuat sesi opname baru (`mulaiOpname()`) dan menampilkan seluruh data yang perlu diperiksa: stok kemasan utuh per batch, serta data kemasan terbuka aktif untuk produk Partial Use.
3. User melakukan pemeriksaan fisik di gudang dan menginput hasil perhitungan ke sistem.
4. User memilih **Hitung Selisih**; sistem membandingkan data fisik dengan data sistem dan menghitung selisih per item (`hitungSelisih()`).

### 12.2 Alternate Flow — Tidak Ada Selisih
5a. Sistem menampilkan konfirmasi: **"Tidak ditemukan selisih stok. Selesaikan stok opname?"**
6a. Jika **Lanjut**: sistem menyimpan hasil opname dan menampilkan pesan **"Stok opname berhasil diselesaikan."** (`statusOpname = SELESAI`).

### 12.3 Alternate Flow — Ditemukan Selisih
5b. Sistem menampilkan rincian perbedaan yang terdeteksi dan meminta `keterangan` untuk setiap item yang mengalami selisih.
6b. User melengkapi keterangan dan memilih **Simpan**.
7b. Sistem menampilkan konfirmasi: **"Apakah Anda yakin ingin menyimpan hasil stok opname dan melakukan penyesuaian stok?"**
8b. Jika **Batal** → penyimpanan dibatalkan.
9b. Jika **Lanjut** → sistem melakukan penyesuaian stok berdasarkan hasil fisik (`lakukanPenyesuaian()`), menyimpan riwayat koreksi (`detail_opname`), memperbarui total stok, dan menampilkan pesan **"Stok opname berhasil disimpan."**

### 12.4 Business Rules
- BR-09.1: `selisih = stokFisik − stokSistem`.
- BR-09.2: Setiap penyesuaian wajib memiliki `keterangan` (tidak boleh kosong) sebagai bagian dari histori koreksi.
- BR-09.3: `DetailOpname` dapat berelasi dengan `BatchStok` (stok kemasan utuh) **atau** `KemasanTerbuka` (sisa isi), tidak keduanya pada satu baris yang sama — sesuai relasi `KemasanTerbuka → DetailOpname` bersifat `0..1 : N` dengan `idKemasanTerbuka` nullable.

### 12.5 Acceptance Criteria
- AC-09.1: Sesi opname yang dibatalkan tidak mengubah stok aktual (`statusOpname = DIBATALKAN`).
- AC-09.2: Setiap penyesuaian stok tercatat lengkap dengan `stokSistem`, `stokFisik`, `selisih`, dan `keterangan` pada `detail_opname`.
- AC-09.3: Total stok produk setelah opname selalu konsisten dengan akumulasi `stokFisik` hasil opname terakhir.

---

## 13. KF-10 — Laporan Stok

**Use Case:** UC-08 Laporan Stok
**Aktor:** Admin Farmasi
**Deskripsi:** Sistem menghasilkan laporan stok masuk, stok keluar, dan sisa stok secara periodik dalam satu sistem terintegrasi — menjawab permasalahan P-05.

### 13.1 Main Flow
1. User mengakses menu Laporan Stok.
2. Sistem menampilkan formulir parameter: `jenisLaporan`, `periodeWaktu` (tanggal awal–akhir), `kategoriProduk`.
3. User mengisi parameter dan menekan **Tampilkan**.
4. Sistem memvalidasi parameter input.

### 13.2 Alternate Flow — Validasi Gagal
Sistem menampilkan pesan kesalahan dan meminta user memperbaiki parameter apabila:
- Tanggal akhir lebih kecil dari tanggal awal.
- Terdapat parameter wajib yang belum diisi.

### 13.3 Main Flow (lanjutan) — Validasi Berhasil
5. Sistem mengambil data sesuai parameter dari basis data, melakukan perhitungan yang diperlukan, dan menyusun laporan berdasarkan jenis laporan yang dipilih.
6. Sistem menampilkan hasil laporan kepada user.
7. User dapat meninjau laporan atau mengunduh laporan untuk kebutuhan dokumentasi.
8. Sistem menampilkan pesan **"Laporan berhasil dibuat."**

### 13.4 Jenis Laporan
| Jenis Laporan | Isi |
|---|---|
| Laporan Stok Masuk | Rekap penerimaan produk per periode: tanggal, produk, batch, jumlah, expired date |
| Laporan Stok Keluar | Rekap penggunaan produk per periode: tanggal, produk, batch, pola penggunaan, jumlah/isi terpakai |
| Laporan Sisa Stok | Posisi stok terkini per produk dan batch pada akhir periode yang dipilih |

### 13.5 Business Rules
- BR-10.1: Laporan dapat difilter per kategori produk; kosongkan filter berarti seluruh kategori.
- BR-10.2: Format unduhan laporan direkomendasikan PDF dan/atau Excel (rincian teknis pada `srs-backend.md` dan `srs-api.md`).

### 13.6 Acceptance Criteria
- AC-10.1: Parameter tanggal akhir < tanggal awal selalu ditolak sistem.
- AC-10.2: Nilai pada laporan sisa stok konsisten dengan nilai yang ditampilkan pada modul Monitoring Stok untuk periode yang sama.

---

## 14. Matriks Traceability

| KF | UC | Kelas Terkait (Class Diagram) | Permasalahan Sistem Berjalan yang Diselesaikan |
|---|---|---|---|
| KF-01 | UC-01 | User | — |
| KF-02 | UC-02 | Kategori | — |
| KF-03 | UC-03 | Produk | — |
| KF-04 | UC-04 | StokMasuk, BatchStok | P-01, P-05 |
| KF-05 | UC-05 | StokKeluar, BatchStok, KemasanTerbuka | P-02 |
| KF-06 | UC-05 (include) | BatchStok | P-01 |
| KF-07 | UC-05 (extend) | KemasanTerbuka | P-03 |
| KF-08 | UC-06 | Produk, BatchStok, KemasanTerbuka | P-01, P-04 |
| KF-09 | UC-07 | StokOpname, DetailOpname | P-04 |
| KF-10 | UC-08 | StokMasuk, StokKeluar, BatchStok | P-05 |
