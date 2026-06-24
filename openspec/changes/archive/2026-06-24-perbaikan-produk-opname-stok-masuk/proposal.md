## Why

Terdapat tiga area yang perlu diperbaiki dan disempurnakan untuk memastikan konsistensi desain sistem FEFO + BUD: (1) tampilan kolom "isi per kemasan" di kelola produk salah kalkulasi untuk produk full use, (2) alur stok opname belum mencerminkan proses bisnis yang benar — selisih harus otomatis menjadi dasar stock adjustment, UI opname perlu dipisah tab Full/Partial dalam 1 sesi, dan keterangan wajib jika ada selisih, (3) stok masuk perlu fitur edit dan hapus terbatas — hanya jika batch belum pernah digunakan di stok keluar.

## What Changes

- **[KF-03] Fix tampilan isi per kemasan produk Full Use**: Kolom "Isi Per Kemasan" di halaman kelola produk saat ini mengalikan nilai dengan jumlah kemasan untuk Full Use. Harusnya Full Use ditampilkan sebagai "per pcs" tanpa kalkulasi isi.

- **[KF-09] Revamp alur stok opname**:
  - UI opname dipisah tab **Full Use** dan **Partial Use** dalam 1 sesi yang sama
  - Tab Full Use: input stok fisik kemasan (integer)
  - Tab Partial Use: input stok fisik kemasan + sisa isi kemasan terbuka (float)
  - Sistem otomatis hitung selisih (stok fisik - stok sistem)
  - Jika ada selisih, keterangan WAJIB diisi
  - Saat sesi diselesaikan, sistem otomatis koreksi `batch_stok.stok_kemasan` dan `kemasan_terbuka.isi_tersisa` sesuai stok fisik
  - Hasil opname menjadi acuan stok terbaru — Monitoring membaca hasil akhir ini
  - Stok keluar tidak ada edit/hapus; koreksi salah input dilakukan via stok opname

- **[KF-04] Edit dan hapus stok masuk terbatas**:
  - Edit dan hapus hanya diizinkan jika batch belum digunakan di stok keluar
  - Jika batch sudah digunakan: edit dan hapus ditolak dengan pesan error jelas
  - Saat edit: sistem menyesuaikan `batch_stok.stok_kemasan` dan `total_isi_tersedia`
  - Saat hapus: sistem mengurangi/menghapus batch dan menyesuaikan stok produk

## Capabilities

### New Capabilities
- `edit-hapus-stok-masuk`: Edit dan hapus stok masuk terbatas — hanya jika batch belum digunakan

### Modified Capabilities
- `kelola-produk`: Fix tampilan kolom isi per kemasan untuk produk Full Use
- `stok-opname`: Revamp alur — tab Full/Partial, selisih otomatis jadi stock adjustment, keterangan wajib jika selisih

## Impact

**Backend:**
- `internal/handler/stok_masuk_handler.go` — tambah endpoint PUT dan DELETE stok masuk
- `internal/service/stok_masuk_service.go` — tambah logika update/delete dengan guard batch belum dipakai
- `internal/repository/stok_masuk_repository.go` — tambah Update, Delete, CheckBatchUsed
- `internal/handler/opname_handler.go` — update endpoint SelesaikanOpname untuk handle selisih + keterangan wajib
- `internal/service/opname_service.go` — update logika SelesaikanOpname: hitung selisih, koreksi batch/kemasan, validasi keterangan wajib
- `internal/repository/opname_repository.go` — update SaveDetailAndAdjust: koreksi batch + kemasan terbuka sesuai stok fisik

**Frontend:**
- `src/features/produk/` — fix tampilan kolom isi per kemasan (full use = "per pcs")
- `src/features/stok-masuk/` — tambah tombol edit dan hapus dengan guard
- `src/features/opname/` — revamp UI: tab Full/Partial, field sisa isi kemasan terbuka, keterangan wajib jika selisih

**API (new endpoints):**
- `PUT /api/stok-masuk/:id` — edit stok masuk (terbatas)
- `DELETE /api/stok-masuk/:id` — hapus stok masuk (terbatas)

**Migration:** Tidak diperlukan — semua perubahan di layer aplikasi.

**Modul terdampak:** KF-03, KF-04, KF-09

**Acceptance Criteria:**
- Produk Full Use di halaman kelola produk menampilkan "per pcs" pada kolom isi per kemasan, bukan hasil kalkulasi
- Edit stok masuk berhasil jika batch belum digunakan, gagal dengan error jelas jika sudah digunakan
- Hapus stok masuk berhasil jika batch belum digunakan, gagal dengan error jelas jika sudah digunakan
- Saat edit stok masuk, stok batch dan total stok produk ikut menyesuaikan
- Saat hapus stok masuk, batch dan total stok produk ikut berkurang
- UI stok opname menampilkan dua tab: Full Use dan Partial Use dalam 1 sesi
- Tab Partial Use menampilkan field sisa isi kemasan terbuka
- Jika selisih ≠ 0, field keterangan wajib diisi sebelum sesi bisa diselesaikan
- Setelah sesi opname selesai, `batch_stok.stok_kemasan` dan `kemasan_terbuka.isi_tersisa` ter-update sesuai stok fisik
- Monitoring menampilkan data stok terbaru hasil opname
