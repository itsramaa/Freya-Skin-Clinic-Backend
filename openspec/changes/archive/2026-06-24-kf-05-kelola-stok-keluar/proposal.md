## Why

KF-05 adalah modul pencatatan penggunaan stok (stok keluar) yang mencakup dua alur berbeda: Full Use dan Partial Use. Modul ini bergantung pada KF-06 (FEFO) untuk pemilihan batch otomatis dan KF-07 (BUD) untuk pengelolaan kemasan terbuka. Ini adalah modul paling kompleks karena ada percabangan logic yang signifikan.

## What Changes

- Migration tabel `stok_keluar` dan `kemasan_terbuka`
- Repository: create stok keluar, update batch stok, create/update kemasan terbuka
- Service: FEFO batch selection, Full Use flow, Partial Use flow + BUD logic
- Handler: `POST /api/stok-keluar`, `GET /api/stok-keluar`, `GET /api/stok-keluar/preview-batch`
- Background worker update status kemasan_terbuka (di KF-07)

## Capabilities

### New Capabilities

- `stok-keluar-management`: Pencatatan penggunaan stok dengan FEFO otomatis, percabangan Full/Partial Use, dan BUD management
- `preview-batch`: Endpoint preview batch FEFO sebelum submit stok keluar

### Modified Capabilities

*(tidak ada)*

## Impact

**API Endpoints:**
- `GET /api/stok-keluar/preview-batch?produk_id=` → preview batch prioritas FEFO
- `GET /api/stok-keluar` → list riwayat penggunaan
- `POST /api/stok-keluar` → catat penggunaan (FEFO + Full/Partial Use logic)

**Dependencies:** KF-04 (batch_stok), KF-06 (FEFO), KF-07 (BUD)

**Acceptance Criteria:**
- AC-05.1: Full Use mengurangi stok_kemasan batch sesuai jumlah input
- AC-05.2: Partial Use tanpa kemasan terbuka → buat kemasan_terbuka baru + BUD 28 hari
- AC-05.3: Partial Use dengan kemasan terbuka BUD expired → nonaktifkan + buka kemasan baru
- AC-05.4: User tidak dapat memilih batch manual (FEFO otomatis)
