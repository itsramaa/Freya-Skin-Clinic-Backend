## Why

Terdapat beberapa bug yang ditemukan melalui tracing kode pada fitur stok keluar, monitoring, dan laporan stok. Bug-bug ini menyebabkan: (1) stok keluar full use selalu gagal, (2) partial use tidak bisa buka kemasan baru setelah kemasan lama habis, dan (3) field `satuan_isi` tidak dikirim ke frontend pada endpoint stok keluar, laporan stok keluar, laporan sisa stok, dan monitoring.

## What Changes

- **[KF-05] Fix full use selalu gagal**: `jumlah_kemasan_dipakai` default value tidak di-set di frontend form modal, menyebabkan backend menerima nilai `0` dan validasi `<= 0` langsung gagal.
- **[KF-05] Fix buka kemasan baru gagal setelah kemasan habis**: Fallback `isi_per_kemasan` di service menggunakan `1.0` jika produk tidak punya `isi_per_kemasan`, sehingga input jumlah isi > 1 langsung kena `ErrIsiDipakaiMelebihiSisa`.
- **[KF-05/KF-08/KF-10] Fix satuan_isi tidak muncul**: Field `satuan_isi` tidak di-SELECT dan tidak di-Scan di tiga repository: `GetStokKeluar` (laporan), `GetSisaStok` (laporan), dan `FindAllForMonitoring` (monitoring). Struct Go juga belum punya field tersebut di `LaporanStokKeluarItem`, `LaporanSisaStokItem`, dan `MonitoringProdukItem`.
- **[KF-05] Fix SQL bug ReduceStok**: Subquery `NOT EXISTS` di `batch_fefo_repository.go` menggunakan `id_batch = id` yang ambigu — `id` resolve ke kolom tabel itu sendiri, bukan ke parameter batch yang sedang di-update, sehingga status batch tidak pernah ter-set `HABIS` dengan benar.

## Capabilities

### New Capabilities

- none

### Modified Capabilities

- `kelola-stok-keluar`: Perbaikan validasi input full use, fallback isi_per_kemasan, dan SQL ReduceStok
- `monitoring-stok`: Tambah field `satuan_isi` pada response monitoring produk
- `laporan-stok`: Tambah field `satuan_isi` pada response stok keluar dan sisa stok

## Impact

**Backend:**

- `internal/model/laporan.go` — tambah `SatuanIsi` ke `LaporanStokKeluarItem` dan `LaporanSisaStokItem`
- `internal/model/monitoring.go` — tambah `SatuanIsi` ke `MonitoringProdukItem`
- `internal/repository/laporan_repository.go` — tambah `p.satuan_isi` di SELECT + Scan untuk `GetStokKeluar` dan `GetSisaStok`
- `internal/repository/monitoring_repository.go` — tambah `p.satuan_isi` di SELECT + Scan untuk `FindAllForMonitoring`
- `internal/repository/batch_fefo_repository.go` — fix subquery SQL di `ReduceStok`
- `internal/service/stok_keluar_service.go` — fix fallback `isi_per_kemasan` saat buka kemasan baru

**Frontend:**

- `src/features/stok-keluar/components/StokKeluarFormModal.tsx` — tambah `jumlah_kemasan_dipakai: 1` di `defaultValues`

**API:** Tidak ada perubahan endpoint atau contract — hanya menambah field yang seharusnya sudah ada.

**Migration:** Tidak diperlukan — semua perubahan di layer aplikasi.

**Modul terdampak:** KF-05, KF-08, KF-10

**Acceptance Criteria:**

- Full use berhasil disimpan dengan `jumlah_kemasan_dipakai = 1` (default)
- Partial use berhasil buka kemasan baru setelah kemasan lama habis, tanpa error `ErrIsiDipakaiMelebihiSisa`
- Response `/api/stok-keluar` menyertakan field `satuan_isi`
- Response `/api/laporan/stok-keluar` menyertakan field `satuan_isi`
- Response `/api/laporan/sisa-stok` menyertakan field `satuan_isi`
- Response `/api/monitoring` menyertakan field `satuan_isi` per produk
- Batch ter-set status `HABIS` dengan benar setelah stok kemasan habis dan tidak ada kemasan terbuka aktif
