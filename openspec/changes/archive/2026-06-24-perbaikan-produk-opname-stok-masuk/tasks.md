## 1. Frontend — Kelola Produk

- [x] 1.1 Fix tampilan kolom isi per kemasan di `src/features/produk/pages/ProdukPage.tsx` — Full Use tampilkan "per pcs", Partial Use tampilkan `isi_per_kemasan + satuan_isi`

## 2. Backend — Model

- [x] 2.1 Tambah error var `ErrBatchSudahDigunakan` di `internal/service/stok_masuk_service.go`
- [x] 2.2 Tambah struct `UpdateStokMasukRequest` di `internal/model/stok_masuk.go`
- [x] 2.3 Tambah field `SisaIsiTerbuka *float64` di `DetailOpnameInput` di `internal/model/opname.go` untuk input sisa isi kemasan terbuka saat opname

## 3. Backend — Repository Stok Masuk

- [x] 3.1 Tambah method `CheckBatchUsed(ctx, idBatch) (bool, error)` di `internal/repository/stok_masuk_repository.go` — cek eksistensi di `stok_keluar`
- [x] 3.2 Tambah method `Update(ctx, id, req) error` di `internal/repository/stok_masuk_repository.go` — UPDATE stok_masuk + batch delta
- [x] 3.3 Tambah method `Delete(ctx, id) error` di `internal/repository/stok_masuk_repository.go` — DELETE stok_masuk + DELETE batch_stok
- [x] 3.4 Update interface `StokMasukRepository` dengan method baru

## 4. Backend — Service Stok Masuk

- [x] 4.1 Tambah method `Update(ctx, id, req, userID) error` di `internal/service/stok_masuk_service.go` — guard CheckBatchUsed + hitung delta + update
- [x] 4.2 Tambah method `Delete(ctx, id) error` di `internal/service/stok_masuk_service.go` — guard CheckBatchUsed + delete
- [x] 4.3 Update interface `StokMasukService` dengan method baru

## 5. Backend — Handler Stok Masuk

- [x] 5.1 Tambah handler `Update` di `internal/handler/stok_masuk_handler.go` — parse body, call service, handle ErrBatchSudahDigunakan
- [x] 5.2 Tambah handler `Delete` di `internal/handler/stok_masuk_handler.go` — call service, handle ErrBatchSudahDigunakan
- [x] 5.3 Daftarkan route `PUT /api/stok-masuk/:id` dan `DELETE /api/stok-masuk/:id` di `internal/router/`

## 6. Backend — Repository Opname

- [x] 6.1 Update `SaveDetailAndAdjust` di `internal/repository/opname_repository.go` — gunakan nilai stok fisik langsung (SET = stok_fisik, bukan += selisih)
- [x] 6.2 Update logika UPDATE `kemasan_terbuka` — gunakan field `SisaIsiTerbuka` dari `DetailOpnameInput` untuk partial use

## 7. Backend — Service Opname

- [x] 7.1 Update `SelesaikanOpname` di `internal/service/opname_service.go` — tambah validasi: jika selisih ≠ 0 dan keterangan kosong → return `ErrKeteranganWajib`
- [x] 7.2 Tambah error var `ErrKeteranganWajib` di service opname

## 8. Backend — Handler Opname

- [x] 8.1 Tambah case `ErrKeteranganWajib` di handler opname — return HTTP 400 dengan pesan "Keterangan wajib diisi untuk item yang memiliki selisih"

## 9. Frontend — Stok Masuk

- [x] 9.1 Tambah tombol edit di tabel stok masuk `src/features/stok-masuk/pages/StokMasukPage.tsx`
- [x] 9.2 Tambah tombol hapus di tabel stok masuk dengan konfirmasi dialog
- [x] 9.3 Buat `StokMasukEditModal.tsx` — form edit jumlah kemasan + tanggal + keterangan
- [x] 9.4 Tambah hook `useUpdateStokMasuk` dan `useDeleteStokMasuk` di `src/features/stok-masuk/hooks/useStokMasuk.ts`
- [x] 9.5 Tambah API call `updateStokMasuk` dan `deleteStokMasuk` di `src/features/stok-masuk/api/stokMasukApi.ts`
- [x] 9.6 Handle error "batch sudah digunakan" di UI — tampilkan toast/alert yang jelas

## 10. Frontend — Stok Opname

- [x] 10.1 Revamp `src/features/opname/pages/StokOpnameDetailPage.tsx` — pisahkan tab Full Use dan Partial Use
- [x] 10.2 Tab Full Use: tampilkan batch + field input stok fisik kemasan
- [x] 10.3 Tab Partial Use: tampilkan batch + field input stok fisik kemasan + field input sisa isi kemasan terbuka
- [x] 10.4 Tambah kalkulasi selisih real-time di frontend (stok fisik - stok sistem)
- [x] 10.5 Tampilkan warna berbeda untuk selisih positif (hijau), negatif (merah), nol (abu)
- [x] 10.6 Tampilkan field keterangan per item — required jika selisih ≠ 0
- [x] 10.7 Disable tombol "Selesaikan Opname" jika ada item selisih tanpa keterangan
- [x] 10.8 Update tipe `DetailOpnameInput` di frontend untuk sertakan `sisa_isi_terbuka`

## 11. Verifikasi

- [x] 11.1 Build backend berhasil: `go build ./...`
- [x] 11.2 Test edit stok masuk batch belum digunakan → berhasil
- [x] 11.3 Test edit stok masuk batch sudah digunakan → HTTP 400
- [x] 11.4 Test hapus stok masuk batch belum digunakan → berhasil
- [x] 11.5 Test hapus stok masuk batch sudah digunakan → HTTP 400
- [x] 11.6 Test selesaikan opname dengan selisih tanpa keterangan → HTTP 400
- [x] 11.7 Test selesaikan opname dengan selisih + keterangan → stok ter-update
- [x] 11.8 Verifikasi monitoring menampilkan stok terbaru setelah opname
- [x] 11.9 Verifikasi kolom isi per kemasan produk Full Use menampilkan "per pcs"
