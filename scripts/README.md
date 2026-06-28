# Database Reset & Seed - Freya Skin Clinic Backend

## Overview
Script untuk cleansing data produksi dan seeding data dummy untuk development/testing.

## Files
- `scripts/reset_and_seed.sql` - SQL script untuk cleansing dan seeding
- `scripts/run_sql.go` - Helper Go untuk menjalankan SQL script
- `scripts/verify_seed.go` - Helper Go untuk verifikasi hasil seeding

## Cara Penggunaan

### 1. Reset & Seed Database Produksi
```bash
# Menggunakan Go helper (recommended untuk Windows)
go run scripts/run_sql.go scripts/reset_and_seed.sql .env.production

# Atau menggunakan Makefile (memerlukan psql di PATH)
make reset-seed-prod
```

### 2. Verifikasi Hasil Seeding
```bash
go run scripts/verify_seed.go .env.production
```

## Data yang Di-seed

### Kategori (5)
1. **Skincare** - Produk perawatan kulit
2. **Injectable** - Botox, filler, mesotherapy
3. **Obat** - Obat tablet dan salep
4. **Threadlift** - PDO threads untuk prosedur
5. **Facial IPL Laser** - Konsumabel untuk laser/IPL

### Produk (14)
- Skincare: 3 produk (PARTIAL_USE)
- Injectable: 3 produk (FULL_USE & PARTIAL_USE)
- Obat: 3 produk (FULL_USE & PARTIAL_USE)
- Threadlift: 2 produk (FULL_USE)
- Facial IPL Laser: 3 produk (PARTIAL_USE & FULL_USE)

### Users
- Semua user di-reset ke **default password**: `admin`
- `is_default_password = true` (memaksa ganti password saat login)
- `session_id = NULL` (invalidate semua sesi aktif)

## Efek Reset Auth

Setelah reset:
1. User login dengan username & password default `admin`
2. Backend return `is_default_password: true`
3. Frontend redirect ke halaman **Ganti Password**
4. User harus ganti password sebelum bisa akses aplikasi

## Catatan Penting

⚠️ **WARNING**: Script ini akan:
- Menghapus SEMUA data transaksi (stok_masuk, stok_keluar, batch_stok, opname, dll)
- Menghapus SEMUA produk dan kategori
- Reset password SEMUA user ke default
- Invalidate SEMUA sesi login aktif

✅ **Safe**: Script ini TIDAK menghapus:
- Tabel users (hanya reset password)
- Schema database
- Migrations history
