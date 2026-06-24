## Context

Tiga area berbeda dalam sistem perlu diperbaiki dan disempurnakan secara bersamaan karena saling berkaitan dalam alur bisnis FEFO + BUD:

1. **Tampilan produk** — kolom isi per kemasan untuk Full Use dikalikan secara salah di frontend
2. **Stok opname** — alur saat ini tidak memisahkan Full/Partial di UI, tidak mewajibkan keterangan jika ada selisih, dan koreksi stok belum konsisten dengan desain database
3. **Stok masuk** — tidak ada fitur edit/hapus sama sekali, padahal ada kebutuhan koreksi data entry sebelum batch digunakan

Ketiga area ini saling terkait: stok opname adalah mekanisme koreksi utama untuk stok keluar yang salah, dan stok masuk perlu guard yang ketat agar tidak merusak konsistensi data jika batch sudah masuk dalam transaksi.

## Goals / Non-Goals

**Goals:**
- Fix kalkulasi tampilan isi per kemasan produk Full Use di frontend
- Revamp UI stok opname: tab Full/Partial dalam 1 sesi, field sisa isi kemasan terbuka di tab Partial, keterangan wajib jika selisih ≠ 0
- Stock adjustment otomatis saat sesi opname selesai: koreksi `batch_stok` dan `kemasan_terbuka` sesuai stok fisik
- Tambah endpoint edit dan hapus stok masuk dengan guard "batch belum digunakan"

**Non-Goals:**
- Edit atau hapus stok keluar (by design — koreksi via opname)
- Perubahan skema database (tidak ada migration baru)
- Multi-role atau permission baru
- Undo/revert sesi opname yang sudah selesai

## Decisions

### D1 — Fix tampilan isi per kemasan Full Use

**Keputusan:** Frontend menampilkan "per pcs" untuk Full Use, tanpa kalkulasi apapun. `isi_per_kemasan` hanya relevan untuk Partial Use.

**Alasan:** Full Use tidak memiliki konsep "isi per kemasan" yang perlu ditampilkan — satu kemasan = satu unit. Menampilkan hasil kalkulasi menyesatkan.

---

### D2 — Tab Full/Partial dalam 1 sesi opname

**Keputusan:** 1 sesi `stok_opname` mencakup semua produk (Full dan Partial). UI dipisah tab untuk kemudahan input, bukan dipisah sesi.

**Alasan:** Opname dilakukan sekaligus untuk seluruh stok. Memisahkan sesi akan mempersulit rekonsiliasi dan audit trail.

**Flow UI:**
```
Sesi Opname Aktif
├── Tab Full Use
│   └── Input: stok_fisik_kemasan (int) per batch
└── Tab Partial Use
    ├── Input: stok_fisik_kemasan (int) per batch
    └── Input: sisa_isi_terbuka (float) per kemasan terbuka
```

---

### D3 — Keterangan wajib jika selisih ≠ 0

**Keputusan:** Validasi di service — jika `selisih != 0` dan `keterangan` kosong, return error. Frontend juga menampilkan field keterangan sebagai required jika selisih terdeteksi.

**Alasan:** Audit trail — setiap penyesuaian stok harus ada justifikasi tercatat.

---

### D4 — Stock adjustment otomatis saat sesi selesai

**Keputusan:** `SaveDetailAndAdjust` di repository langsung update `batch_stok.stok_kemasan` dan `kemasan_terbuka.isi_tersisa` ke nilai stok fisik (bukan ditambah selisih).

**Flow:**
```
SelesaikanOpname (Handler)
  → OpnameService.Selesaikan(details)
      → Validasi: selisih ≠ 0 → keterangan wajib
      → Repository.SaveDetailAndAdjust (dalam 1 DB transaction):
          ├── INSERT detail_opname (stok_sistem, stok_fisik, selisih, keterangan)
          ├── UPDATE batch_stok SET stok_kemasan = stok_fisik WHERE id = id_batch
          ├── UPDATE kemasan_terbuka SET isi_tersisa = sisa_isi_terbuka WHERE id = id_kemasan
          └── UPDATE stok_opname SET status = 'SELESAI'
```

**Alasan:** Set langsung ke stok fisik lebih aman daripada menerapkan selisih — menghindari double-apply jika ada retry.

---

### D5 — Guard edit/hapus stok masuk: cek batch digunakan

**Keputusan:** Cek existensi di tabel `stok_keluar` berdasarkan `id_batch`. Jika ada record, tolak edit/hapus.

**Flow edit:**
```
PUT /api/stok-masuk/:id (Handler)
  → StokMasukService.Update(id, req)
      → CheckBatchUsed(id_batch) → jika dipakai → ErrBatchSudahDigunakan
      → Hitung delta (jumlah_kemasan_baru - jumlah_kemasan_lama)
      → UPDATE stok_masuk
      → UPDATE batch_stok SET stok_kemasan += delta, total_isi_tersedia += delta_isi
```

**Flow hapus:**
```
DELETE /api/stok-masuk/:id (Handler)
  → StokMasukService.Delete(id)
      → CheckBatchUsed(id_batch) → jika dipakai → ErrBatchSudahDigunakan
      → DELETE stok_masuk
      → DELETE batch_stok (jika tidak ada transaksi lain)
```

---

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| Sesi opname selesai, lalu ada stok keluar baru masuk sebelum refresh — data stok bisa inkonsisten sesaat | Acceptable — opname adalah snapshot waktu tertentu, bukan real-time lock |
| Hapus stok masuk menghapus batch yang punya kemasan terbuka aktif | Guard: cek juga `kemasan_terbuka` sebelum hapus batch |
| Edit stok masuk mengubah `jumlah_kemasan` — delta bisa negatif jika stok fisik lebih sedikit dari yang sudah dipakai | Guard: validasi `stok_kemasan - delta >= 0` sebelum apply |
| Frontend tab opname menampilkan data berbeda antara Full dan Partial — perlu endpoint `GET /api/opname/items` yang sudah ada memisahkan keduanya | Gunakan field `pola_penggunaan` di response untuk filter di frontend |

## Open Questions

- Apakah sesi opname yang sudah `SELESAI` bisa dibuka kembali untuk koreksi, atau harus buat sesi baru?
  - **Asumsi saat ini:** Tidak bisa dibuka ulang — buat sesi baru jika perlu koreksi lanjutan.
- Apakah hapus stok masuk juga menghapus batch jika `stok_keluar` kosong tapi `stok_kemasan > 0`?
  - **Asumsi saat ini:** Ya, hapus batch sekaligus jika tidak ada referensi stok keluar.
