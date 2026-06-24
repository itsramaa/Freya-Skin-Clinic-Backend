## MODIFIED Requirements

### Requirement: Catat penggunaan stok full use
Sistem SHALL memproses penggunaan stok produk bertipe FULL_USE. Frontend MUST mengirimkan `jumlah_kemasan_dipakai` dengan nilai default minimal 1. Backend MUST menolak request jika `jumlah_kemasan_dipakai <= 0`.

#### Scenario: Full use berhasil dengan nilai default
- **WHEN** user submit form stok keluar untuk produk FULL_USE tanpa mengubah field jumlah
- **THEN** sistem menggunakan nilai default `jumlah_kemasan_dipakai = 1` dan menyimpan data berhasil

#### Scenario: Full use gagal jika jumlah 0
- **WHEN** request dikirim dengan `jumlah_kemasan_dipakai = 0`
- **THEN** sistem mengembalikan HTTP 400 dengan pesan error validasi

### Requirement: Catat penggunaan stok partial use — buka kemasan baru
Sistem SHALL membuka kemasan baru saat tidak ada kemasan terbuka aktif. Sistem MUST menggunakan nilai `isi_per_kemasan` dari data produk. Sistem MUST mengembalikan error eksplisit jika produk partial use tidak memiliki `isi_per_kemasan` yang dikonfigurasi.

#### Scenario: Buka kemasan baru berhasil
- **WHEN** tidak ada kemasan terbuka aktif dan produk memiliki `isi_per_kemasan` valid
- **THEN** sistem membuat kemasan terbuka baru dengan `isi_awal = isi_per_kemasan`, mengurangi `stok_kemasan` sebesar 1, dan menyimpan stok keluar

#### Scenario: Buka kemasan baru gagal jika isi_per_kemasan tidak dikonfigurasi
- **WHEN** produk partial use tidak memiliki `isi_per_kemasan` (null)
- **THEN** sistem mengembalikan HTTP 400 dengan pesan "Produk tidak memiliki konfigurasi isi per kemasan"

#### Scenario: Buka kemasan baru gagal jika jumlah isi melebihi isi_per_kemasan
- **WHEN** `jumlah_isi_dipakai` lebih besar dari `isi_per_kemasan` produk
- **THEN** sistem mengembalikan HTTP 400 dengan pesan "Jumlah isi yang dipakai melebihi sisa isi kemasan terbuka"

### Requirement: Response stok keluar menyertakan satuan isi
Sistem SHALL menyertakan field `satuan_isi` pada setiap item response GET /api/stok-keluar.

#### Scenario: Satuan isi tersedia di response
- **WHEN** client melakukan GET /api/stok-keluar
- **THEN** setiap item response memiliki field `satuan_isi` berisi satuan produk (misal: "ml", "gram", "pcs")

### Requirement: Status batch ter-set HABIS dengan benar
Sistem SHALL mengubah status batch menjadi `HABIS` saat `stok_kemasan <= 0` dan tidak ada kemasan terbuka aktif. Query UPDATE MUST menggunakan alias tabel eksplisit pada subquery untuk menghindari ambiguitas referensi kolom.

#### Scenario: Batch jadi HABIS setelah stok habis
- **WHEN** `stok_kemasan` batch berkurang hingga 0 dan tidak ada kemasan terbuka dengan `status_bud = 'AKTIF'`
- **THEN** `batch_stok.status` ter-update menjadi `'HABIS'`

#### Scenario: Batch tetap AKTIF jika masih ada kemasan terbuka
- **WHEN** `stok_kemasan` batch berkurang hingga 0 tapi masih ada kemasan terbuka aktif dengan `isi_tersisa > 0`
- **THEN** `batch_stok.status` tetap `'AKTIF'`
