## ADDED Requirements

### Requirement: Sistem menyediakan endpoint GET list produk

Sistem SHALL menyediakan endpoint untuk mengambil seluruh data produk beserta nama kategori, stok kemasan, dan total isi tersedia.

**Referensi:** KF-03, srs-fr.md § 6.1

#### Scenario: GET list produk berhasil
- **WHEN** Admin Farmasi mengirim GET ke `/api/produk` dengan Bearer token valid
- **THEN** Sistem mengembalikan 200 dengan array produk, masing-masing memiliki field: id, kode_produk, nama_produk, nama_kategori, bentuk_kemasan, satuan_isi, isi_per_kemasan, pola_penggunaan, stok_kemasan, total_isi_tersedia, has_transaksi

#### Scenario: GET list produk tanpa token
- **WHEN** Request ke `/api/produk` tanpa Authorization header
- **THEN** Sistem mengembalikan 401 Unauthorized

### Requirement: Sistem menyediakan endpoint POST tambah produk

Sistem SHALL menyediakan endpoint untuk menambah produk baru dengan validasi kelengkapan data dan auto-generate kode produk.

**Referensi:** KF-03, srs-fr.md § 6.2, BR-03.1, BR-03.2, AC-03.1, AC-03.4

#### Scenario: POST produk baru FULL_USE berhasil
- **WHEN** Admin Farmasi mengirim POST ke `/api/produk` dengan semua field wajib (nama_produk, id_kategori, bentuk_kemasan, satuan_isi, pola_penggunaan=FULL_USE) dan token valid
- **THEN** Sistem menyimpan produk dengan kode_produk auto-generated, mengembalikan 201 Created dengan data produk lengkap

#### Scenario: POST produk PARTIAL_USE tanpa isi_per_kemasan
- **WHEN** Admin Farmasi mengirim POST ke `/api/produk` dengan pola_penggunaan=PARTIAL_USE dan isi_per_kemasan kosong/null
- **THEN** Sistem mengembalikan 400 Bad Request dengan message "Isi per kemasan wajib diisi untuk produk Partial Use" (BR-03.2, AC-03.1)

#### Scenario: POST produk dengan field wajib kosong
- **WHEN** Admin Farmasi mengirim POST ke `/api/produk` tanpa salah satu field wajib
- **THEN** Sistem mengembalikan 400 Bad Request dengan detail field yang kosong (AC-03.1)

#### Scenario: POST produk dengan id_kategori tidak valid
- **WHEN** Admin Farmasi mengirim POST ke `/api/produk` dengan id_kategori yang tidak ada di tabel kategori
- **THEN** Sistem mengembalikan 404 Not Found dengan message "Kategori tidak ditemukan"

### Requirement: Sistem menyediakan endpoint PUT ubah produk

Sistem SHALL menyediakan endpoint untuk mengubah data produk, dengan proteksi perubahan `pola_penggunaan` jika produk sudah memiliki transaksi.

**Referensi:** KF-03, srs-fr.md § 6.3, AC-03.5

#### Scenario: PUT ubah produk berhasil
- **WHEN** Admin Farmasi mengirim PUT ke `/api/produk/:id` dengan data valid dan token valid
- **THEN** Sistem memperbarui data produk dan mengembalikan 200 OK dengan data produk terbaru

#### Scenario: PUT ubah pola_penggunaan saat ada transaksi
- **WHEN** Admin Farmasi mengirim PUT ke `/api/produk/:id` dengan pola_penggunaan berbeda dan produk sudah punya riwayat transaksi
- **THEN** Sistem mengembalikan 409 Conflict dengan message "Pola penggunaan tidak dapat diubah karena produk sudah memiliki transaksi." (AC-03.5)

#### Scenario: PUT produk dengan ID tidak ditemukan
- **WHEN** Admin Farmasi mengirim PUT ke `/api/produk/:id` dengan ID yang tidak ada
- **THEN** Sistem mengembalikan 404 Not Found dengan message "Produk tidak ditemukan"

### Requirement: Sistem menyediakan endpoint DELETE hapus produk

Sistem SHALL menyediakan endpoint untuk menghapus produk dengan pengecekan stok aktif dan riwayat transaksi.

**Referensi:** KF-03, srs-fr.md § 6.4, BR-03.3, AC-03.2, AC-03.3

#### Scenario: DELETE produk berhasil
- **WHEN** Admin Farmasi mengirim DELETE ke `/api/produk/:id` untuk produk tanpa stok aktif dan tanpa riwayat transaksi
- **THEN** Sistem menghapus produk dan mengembalikan 200 OK dengan message "Data produk berhasil dihapus."

#### Scenario: DELETE produk dengan stok aktif
- **WHEN** Admin Farmasi mengirim DELETE ke `/api/produk/:id` untuk produk yang memiliki batch_stok dengan status AKTIF
- **THEN** Sistem mengembalikan 409 Conflict dengan message "Produk tidak dapat dihapus karena masih memiliki stok aktif." (AC-03.2)

#### Scenario: DELETE produk dengan riwayat transaksi
- **WHEN** Admin Farmasi mengirim DELETE ke `/api/produk/:id` untuk produk yang pernah ada di stok_masuk atau stok_keluar (stok sudah 0)
- **THEN** Sistem mengembalikan 409 Conflict dengan message "Produk tidak dapat dihapus karena memiliki riwayat transaksi." (AC-03.3)

### Requirement: Kode produk di-generate otomatis

Sistem SHALL men-generate kode produk secara otomatis saat produk baru dibuat dengan format unik.

**Referensi:** KF-03, BR-03.1, AC-03.4

#### Scenario: Kode produk unik di-generate saat tambah produk
- **WHEN** Produk baru berhasil disimpan
- **THEN** Sistem menetapkan kode_produk yang unik dengan format PRD-{kode_kategori}-{sequence} (contoh: PRD-SKC-001)
