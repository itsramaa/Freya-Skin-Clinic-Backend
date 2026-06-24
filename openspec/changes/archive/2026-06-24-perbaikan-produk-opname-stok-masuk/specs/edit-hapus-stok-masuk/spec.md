## ADDED Requirements

### Requirement: Edit stok masuk terbatas
Sistem SHALL mengizinkan edit data stok masuk hanya jika batch terkait belum pernah digunakan dalam transaksi stok keluar. Sistem MUST menolak edit jika batch sudah digunakan. Saat edit berhasil, sistem MUST menyesuaikan `batch_stok.stok_kemasan` dan `batch_stok.total_isi_tersedia` sesuai perubahan jumlah kemasan.

#### Scenario: Edit berhasil jika batch belum digunakan
- **WHEN** user mengedit data stok masuk dan batch terkait belum ada di tabel `stok_keluar`
- **THEN** sistem menyimpan perubahan, menyesuaikan stok batch, dan mengembalikan HTTP 200

#### Scenario: Edit ditolak jika batch sudah digunakan
- **WHEN** user mengedit data stok masuk dan batch terkait sudah ada di tabel `stok_keluar`
- **THEN** sistem mengembalikan HTTP 400 dengan pesan "Batch sudah digunakan dalam transaksi stok keluar, tidak dapat diubah"

#### Scenario: Stok batch menyesuaikan saat edit
- **WHEN** user mengubah `jumlah_kemasan` dari 10 menjadi 15 pada stok masuk
- **THEN** `batch_stok.stok_kemasan` bertambah 5 dan `total_isi_tersedia` menyesuaikan

### Requirement: Hapus stok masuk terbatas
Sistem SHALL mengizinkan hapus data stok masuk hanya jika batch terkait belum pernah digunakan dalam transaksi stok keluar. Sistem MUST menolak hapus jika batch sudah digunakan. Saat hapus berhasil, sistem MUST menghapus batch terkait dan menyesuaikan stok produk.

#### Scenario: Hapus berhasil jika batch belum digunakan
- **WHEN** user menghapus data stok masuk dan batch terkait belum ada di tabel `stok_keluar`
- **THEN** sistem menghapus record stok masuk dan batch, mengembalikan HTTP 200

#### Scenario: Hapus ditolak jika batch sudah digunakan
- **WHEN** user menghapus data stok masuk dan batch terkait sudah ada di tabel `stok_keluar`
- **THEN** sistem mengembalikan HTTP 400 dengan pesan "Batch sudah digunakan dalam transaksi stok keluar, tidak dapat dihapus"

#### Scenario: Stok batch terhapus saat hapus stok masuk
- **WHEN** user menghapus stok masuk yang batchnya belum digunakan
- **THEN** `batch_stok` terkait ikut terhapus dan stok produk berkurang sesuai jumlah kemasan yang dihapus
