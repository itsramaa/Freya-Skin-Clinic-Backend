## MODIFIED Requirements

### Requirement: UI stok opname dipisah tab Full Use dan Partial Use
Sistem SHALL menampilkan UI stok opname dalam dua tab terpisah — Full Use dan Partial Use — dalam satu sesi opname yang sama. Tab Full Use MUST menampilkan field input stok fisik kemasan (integer). Tab Partial Use MUST menampilkan field input stok fisik kemasan (integer) DAN sisa isi kemasan terbuka (float) per batch.

#### Scenario: Tab Full Use menampilkan batch produk Full Use
- **WHEN** user membuka sesi opname aktif dan memilih tab Full Use
- **THEN** sistem menampilkan daftar batch dengan `pola_penggunaan = "FULL_USE"` beserta field input stok fisik kemasan

#### Scenario: Tab Partial Use menampilkan batch dan kemasan terbuka
- **WHEN** user membuka sesi opname aktif dan memilih tab Partial Use
- **THEN** sistem menampilkan daftar batch dengan `pola_penggunaan = "PARTIAL_USE"` beserta field input stok fisik kemasan dan field input sisa isi kemasan terbuka

### Requirement: Selisih otomatis dihitung sistem
Sistem SHALL otomatis menghitung selisih antara stok sistem dan stok fisik yang diinput. Selisih = stok fisik - stok sistem. Nilai selisih MUST ditampilkan secara real-time saat user mengisi stok fisik.

#### Scenario: Selisih positif ditampilkan saat stok fisik lebih banyak
- **WHEN** user mengisi stok fisik lebih besar dari stok sistem
- **THEN** sistem menampilkan selisih positif (contoh: +3)

#### Scenario: Selisih negatif ditampilkan saat stok fisik lebih sedikit
- **WHEN** user mengisi stok fisik lebih kecil dari stok sistem
- **THEN** sistem menampilkan selisih negatif (contoh: -2)

#### Scenario: Tidak ada selisih jika stok fisik sama dengan stok sistem
- **WHEN** user mengisi stok fisik sama dengan stok sistem
- **THEN** sistem menampilkan selisih 0 dan field keterangan tidak wajib

### Requirement: Keterangan wajib jika ada selisih
Sistem MUST memvalidasi bahwa keterangan diisi jika selisih ≠ 0 pada item apapun dalam sesi opname. Sesi opname tidak dapat diselesaikan jika ada selisih tanpa keterangan.

#### Scenario: Sesi gagal diselesaikan jika ada selisih tanpa keterangan
- **WHEN** user menekan tombol selesaikan opname dan ada item dengan selisih ≠ 0 dan keterangan kosong
- **THEN** sistem mengembalikan HTTP 400 dengan pesan "Keterangan wajib diisi untuk item yang memiliki selisih"

#### Scenario: Sesi berhasil diselesaikan jika semua selisih ada keterangan
- **WHEN** user menekan tombol selesaikan opname dan semua item dengan selisih ≠ 0 memiliki keterangan
- **THEN** sistem memproses stock adjustment dan mengembalikan HTTP 200

### Requirement: Stock adjustment otomatis saat sesi opname selesai
Sistem SHALL otomatis mengkoreksi stok saat sesi opname diselesaikan. Untuk Full Use: `batch_stok.stok_kemasan` MUST di-update ke nilai stok fisik. Untuk Partial Use: `batch_stok.stok_kemasan` MUST di-update ke nilai stok fisik kemasan DAN `kemasan_terbuka.isi_tersisa` MUST di-update ke nilai sisa isi fisik. Semua operasi MUST dalam satu DB transaction.

#### Scenario: Stok batch Full Use ter-update sesuai stok fisik
- **WHEN** sesi opname diselesaikan dengan stok fisik Full Use berbeda dari stok sistem
- **THEN** `batch_stok.stok_kemasan` ter-update ke nilai stok fisik yang diinput

#### Scenario: Stok batch dan kemasan terbuka Partial Use ter-update
- **WHEN** sesi opname diselesaikan dengan stok fisik Partial Use berbeda dari stok sistem
- **THEN** `batch_stok.stok_kemasan` ter-update ke nilai stok fisik kemasan DAN `kemasan_terbuka.isi_tersisa` ter-update ke nilai sisa isi fisik

#### Scenario: Monitoring menampilkan stok terbaru setelah opname
- **WHEN** sesi opname selesai dan user membuka halaman monitoring
- **THEN** monitoring menampilkan data stok sesuai hasil opname terbaru

### Requirement: Stok keluar tidak dapat diedit atau dihapus
Sistem MUST NOT menyediakan endpoint atau UI untuk edit atau hapus data stok keluar. Koreksi kesalahan input stok keluar MUST dilakukan melalui sesi stok opname.

#### Scenario: Tidak ada tombol edit atau hapus di halaman stok keluar
- **WHEN** user membuka halaman kelola stok keluar
- **THEN** tidak ada tombol edit atau hapus di tabel maupun detail item

#### Scenario: Endpoint edit/hapus stok keluar tidak tersedia
- **WHEN** client mencoba akses PUT atau DELETE /api/stok-keluar/:id
- **THEN** server mengembalikan HTTP 404 atau 405
