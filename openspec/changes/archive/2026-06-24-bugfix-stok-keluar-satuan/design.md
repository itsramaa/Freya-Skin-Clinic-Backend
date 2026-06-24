## Context

Empat bug ditemukan melalui tracing kode pada fitur stok keluar (KF-05), monitoring (KF-08), dan laporan stok (KF-10). Bug bersifat silently-failing — tidak ada panic, hanya response error atau data kosong yang membingungkan user. Semua fix adalah perubahan minimal di layer aplikasi tanpa schema migration.

**Current state:**
- Full use: selalu gagal karena `jumlah_kemasan_dipakai` tidak punya default value di frontend
- Partial buka kemasan baru: gagal jika `isi_per_kemasan` > 1 karena fallback hardcode ke `1.0`
- `satuan_isi` tidak muncul di response stok keluar, laporan stok keluar, laporan sisa stok, monitoring
- `ReduceStok` SQL: subquery ambigu menyebabkan status batch tidak pernah jadi `HABIS`

## Goals / Non-Goals

**Goals:**
- Full use bisa disimpan dengan default `jumlah_kemasan_dipakai = 1`
- Partial use bisa buka kemasan baru setelah kemasan lama habis
- `satuan_isi` muncul di semua response yang relevan
- Status batch `HABIS` ter-set dengan benar saat stok habis

**Non-Goals:**
- Perubahan schema database
- Perubahan contract API (endpoint, method, path)
- Refactor arsitektur service/repository

## Decisions

### Bug 1 — Frontend default value `jumlah_kemasan_dipakai`

**Keputusan:** Tambah `jumlah_kemasan_dipakai: 1` di `defaultValues` react-hook-form.

**Alasan:** Backend sudah benar — validasi `<= 0` adalah guard yang valid. Masalah di source: form tidak menginisialisasi nilai field. Fix di frontend, bukan menghapus validasi backend.

### Bug 2 — Fallback `isi_per_kemasan` ke `1.0`

**Keputusan:** Ganti fallback `1.0` dengan return `ErrIsiPerKemasanTidakDiset` jika `produk.IsiPerKemasan == nil`.

**Alasan:** Produk partial use wajib punya `isi_per_kemasan`. Fallback `1.0` adalah silent assumption yang berbahaya — lebih baik explicit error supaya user tahu produknya belum dikonfigurasi benar. Alternatif (return 0 atau skip) lebih buruk karena menyembunyikan masalah konfigurasi.

### Bug 3 — `satuan_isi` tidak ada di response

**Keputusan:** Tambah `p.satuan_isi` di SELECT query dan Scan di tiga repository, serta tambah field `SatuanIsi string` di tiga struct model.

**Alasan:** Data sudah ada di DB (`produk.satuan_isi`), tinggal di-project. Tidak perlu JOIN tambahan — produk sudah di-JOIN di semua query tersebut.

**Flow yang terdampak:**
```
GET /api/laporan/stok-keluar
  LaporanHandler → LaporanService → LaporanRepository.GetStokKeluar
    SELECT ... p.satuan_isi ...  ← tambah di sini
    Scan → LaporanStokKeluarItem.SatuanIsi  ← tambah field

GET /api/laporan/sisa-stok
  LaporanHandler → LaporanService → LaporanRepository.GetSisaStok
    SELECT ... p.satuan_isi ...  ← tambah di sini
    Scan → LaporanSisaStokItem.SatuanIsi  ← tambah field

GET /api/monitoring
  MonitoringHandler → MonitoringService → MonitoringRepository.FindAllForMonitoring
    SELECT ... p.satuan_isi ...  ← tambah di sini
    Scan → MonitoringProdukItem.SatuanIsi  ← tambah field
```

### Bug 4 — SQL subquery ambigu di `ReduceStok`

**Keputusan:** Fix referensi kolom di subquery dari `id_batch = id` menjadi `id_batch = batch_stok.id` menggunakan alias tabel eksplisit.

**Query sebelum (bug):**
```sql
UPDATE batch_stok
SET status = CASE
    WHEN (stok_kemasan - $1) <= 0 AND NOT EXISTS (
        SELECT 1 FROM kemasan_terbuka
        WHERE id_batch = id  -- BUG: "id" ambigu, resolve ke kemasan_terbuka.id
        AND status_bud = 'AKTIF' AND isi_tersisa > 0
    ) THEN 'HABIS'
    ELSE status
END
WHERE id = $3
```

**Query setelah (fix):**
```sql
UPDATE batch_stok bs
SET status = CASE
    WHEN (bs.stok_kemasan - $1) <= 0 AND NOT EXISTS (
        SELECT 1 FROM kemasan_terbuka kt
        WHERE kt.id_batch = bs.id  -- FIXED: eksplisit referensi alias
        AND kt.status_bud = 'AKTIF' AND kt.isi_tersisa > 0
    ) THEN 'HABIS'
    ELSE bs.status
END
WHERE bs.id = $3
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| Fix SQL `ReduceStok` bisa mengubah behavior batch yang sebelumnya tidak ter-set HABIS | Acceptable — ini adalah perilaku yang memang diharapkan. Data lama yang status-nya salah tidak otomatis terkoreksi, tapi tidak akan memperburuk keadaan |
| Error baru `ErrIsiPerKemasanTidakDiset` untuk produk partial use tanpa `isi_per_kemasan` | Acceptable — produk tersebut memang tidak terkonfigurasi dengan benar. Error lebih baik daripada silent failure |
| Monitoring repository perlu di-trace lebih lanjut untuk lokasi exact Scan | Baca `monitoring_repository.go` sebelum implement untuk pastikan urutan kolom Scan benar |
