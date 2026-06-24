## 1. Backend — Model

- [x] 1.1 Tambah field `SatuanIsi string` ke struct `LaporanStokKeluarItem` di `internal/model/laporan.go`
- [x] 1.2 Tambah field `SatuanIsi string` ke struct `LaporanSisaStokItem` di `internal/model/laporan.go`
- [x] 1.3 Tambah field `SatuanIsi string` ke struct `MonitoringProdukItem` di `internal/model/monitoring.go`
- [x] 1.4 Tambah error var `ErrIsiPerKemasanTidakDiset` di `internal/service/stok_keluar_service.go`

## 2. Backend — Repository

- [x] 2.1 Fix `ReduceStok` di `internal/repository/batch_fefo_repository.go` — tambah alias tabel `bs` dan `kt`, ganti `id_batch = id` menjadi `kt.id_batch = bs.id`
- [x] 2.2 Tambah `p.satuan_isi` di SELECT query `GetStokKeluar` di `internal/repository/laporan_repository.go`
- [x] 2.3 Tambah `&item.SatuanIsi` di Scan `GetStokKeluar` di `internal/repository/laporan_repository.go`
- [x] 2.4 Tambah `p.satuan_isi` di SELECT query `GetSisaStok` di `internal/repository/laporan_repository.go`
- [x] 2.5 Tambah `&item.SatuanIsi` di Scan `GetSisaStok` di `internal/repository/laporan_repository.go`
- [x] 2.6 Tambah `p.satuan_isi` di SELECT query `FindAllForMonitoring` di `internal/repository/monitoring_repository.go`
- [x] 2.7 Tambah `&p.SatuanIsi` di Scan `FindAllForMonitoring` di `internal/repository/monitoring_repository.go`

## 3. Backend — Service

- [x] 3.1 Fix fallback `isi_per_kemasan` di `internal/service/stok_keluar_service.go` — ganti `isiPerKemasan := 1.0` dengan check nil dan return `ErrIsiPerKemasanTidakDiset` jika `produk.IsiPerKemasan == nil`

## 4. Backend — Handler

- [x] 4.1 Tambah case `ErrIsiPerKemasanTidakDiset` di `internal/handler/stok_keluar_handler.go` — return HTTP 400 dengan pesan "Produk tidak memiliki konfigurasi isi per kemasan"

## 5. Frontend

- [x] 5.1 Tambah `jumlah_kemasan_dipakai: 1` di `defaultValues` react-hook-form di `src/features/stok-keluar/components/StokKeluarFormModal.tsx`

## 6. Verifikasi

- [x] 6.1 Build backend berhasil tanpa error: `go build ./...`
- [x] 6.2 Test full use: submit stok keluar produk FULL_USE → berhasil tersimpan
- [x] 6.3 Test partial buka kemasan baru: setelah kemasan 1 habis, submit lagi → kemasan baru terbuka
- [x] 6.4 Test response monitoring: field `satuan_isi` muncul di response `/api/monitoring`
- [x] 6.5 Test response laporan stok keluar: field `satuan_isi` muncul di response `/api/laporan/stok-keluar`
- [x] 6.6 Test response laporan sisa stok: field `satuan_isi` muncul di response `/api/laporan/sisa-stok`
- [x] 6.7 Test SQL ReduceStok: batch ter-set HABIS setelah stok kemasan habis dan tidak ada kemasan terbuka aktif
