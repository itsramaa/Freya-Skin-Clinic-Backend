-- Script untuk Cleansing Data Produksi dan Seed Data Dummy
-- Freya Skin Clinic Backend
-- 
-- Cara menjalankan:
-- psql <DATABASE_URL> -f scripts/reset_and_seed.sql
-- atau
-- make reset-seed-prod

-- ============================================================================
-- 1. CLEANSING DATA (DELETE ALL TRANSACTION DATA)
-- ============================================================================

BEGIN;

-- Hapus semua data transaksi (keep referential integrity)
DELETE FROM detail_opname;
DELETE FROM stok_opname;
DELETE FROM stok_keluar;
DELETE FROM kemasan_terbuka;
DELETE FROM stok_masuk;
DELETE FROM batch_stok;
DELETE FROM produk;
DELETE FROM kategori;

-- Reset users dan set semua ke default password
UPDATE users SET 
    password_hash = '$2a$12$deA8Pp1xEpWu1iTjTy6O0.HI/cmdnsvE.2VYEBXSUBSOTYWOkCBuC',
    is_default_password = true,
    session_id = NULL,
    updated_at = NOW();

COMMIT;

-- ============================================================================
-- 2. SEED DATA DUMMY
-- ============================================================================

BEGIN;

-- ──────────────────────────────────────────────────────────────────────────
-- 2.1 KATEGORI (5 kategori - sesuai bisnis klinik kulit)
-- ──────────────────────────────────────────────────────────────────────────

INSERT INTO kategori (id, kode_kategori, nama_kategori) VALUES
('11111111-1111-1111-1111-111111111111', 'KAT-001', 'Skincare'),
('22222222-2222-2222-2222-222222222222', 'KAT-002', 'Injectable'),
('33333333-3333-3333-3333-333333333333', 'KAT-003', 'Obat'),
('44444444-4444-4444-4444-444444444444', 'KAT-004', 'Threadlift'),
('55555555-5555-5555-5555-555555555555', 'KAT-005', 'Facial IPL Laser');

-- ──────────────────────────────────────────────────────────────────────────
-- 2.2 PRODUK (15 produk - sesuai bisnis klinik kulit)
-- ──────────────────────────────────────────────────────────────────────────

INSERT INTO produk (id, kode_produk, nama_produk, id_kategori, bentuk_kemasan, satuan_isi, isi_per_kemasan, pola_penggunaan) VALUES
-- Skincare (PARTIAL_USE)
('a1111111-1111-1111-1111-111111111111', 'SKC-001', 'Acne Gel Treatment', '11111111-1111-1111-1111-111111111111', 'Tube', 'gram', 20.000, 'PARTIAL_USE'),
('a2222222-2222-2222-2222-222222222222', 'SKC-002', 'Hydrating Serum', '11111111-1111-1111-1111-111111111111', 'Botol Dropper', 'ml', 30.000, 'PARTIAL_USE'),
('a3333333-3333-3333-3333-333333333333', 'SKC-003', 'Brightening Cream', '11111111-1111-1111-1111-111111111111', 'Jar', 'gram', 50.000, 'PARTIAL_USE'),

-- Injectable (PARTIAL_USE untuk vial, FULL_USE untuk filler syringe)
('b1111111-1111-1111-1111-111111111111', 'INJ-001', 'Botox 50 Unit', '22222222-2222-2222-2222-222222222222', 'Vial', 'unit', 50.000, 'PARTIAL_USE'),
('b2222222-2222-2222-2222-222222222222', 'INJ-002', 'Hyaluronic Acid Filler 1ml', '22222222-2222-2222-2222-222222222222', 'Syringe', 'ml', 1.000, 'FULL_USE'),
('b3333333-3333-3333-3333-333333333333', 'INJ-003', 'Mesotherapy Cocktail', '22222222-2222-2222-2222-222222222222', 'Vial', 'ml', 5.000, 'PARTIAL_USE'),

-- Obat (FULL_USE untuk tablet, PARTIAL_USE untuk salep)
('c1111111-1111-1111-1111-111111111111', 'OBT-001', 'Antibiotic Tablet', '33333333-3333-3333-3333-333333333333', 'Strip', 'tablet', 10.000, 'FULL_USE'),
('c2222222-2222-2222-2222-222222222222', 'OBT-002', 'Pain Reliever Tablet', '33333333-3333-3333-3333-333333333333', 'Strip', 'tablet', 10.000, 'FULL_USE'),
('c3333333-3333-3333-3333-333333333333', 'OBT-003', 'Antibiotic Ointment', '33333333-3333-3333-3333-333333333333', 'Tube', 'gram', 15.000, 'PARTIAL_USE'),

-- Threadlift (FULL_USE - sekali pakai per prosedur)
('d1111111-1111-1111-1111-111111111111', 'THL-001', 'PDO Thread Mono 29G', '44444444-4444-4444-4444-444444444444', 'Pack', 'pcs', 10.000, 'FULL_USE'),
('d2222222-2222-2222-2222-222222222222', 'THL-002', 'PDO Thread Cog 19G', '44444444-4444-4444-4444-444444444444', 'Pack', 'pcs', 5.000, 'FULL_USE'),

-- Facial IPL Laser (PARTIAL_USE - konsumabel per sesi)
('e1111111-1111-1111-1111-111111111111', 'IPL-001', 'IPL Gel Conductivity', '55555555-5555-5555-5555-555555555555', 'Botol Pump', 'ml', 500.000, 'PARTIAL_USE'),
('e2222222-2222-2222-2222-222222222222', 'IPL-002', 'Cooling Gel Post-Laser', '55555555-5555-5555-5555-555555555555', 'Tube', 'gram', 100.000, 'PARTIAL_USE'),
('e3333333-3333-3333-3333-333333333333', 'IPL-003', 'Laser Tip Cover Disposable', '55555555-5555-5555-5555-555555555555', 'Box', 'pcs', 50.000, 'FULL_USE');

COMMIT;

-- ============================================================================
-- SELESAI
-- ============================================================================

-- Verifikasi hasil
SELECT 'Kategori:' AS tabel, COUNT(*) AS jumlah FROM kategori
UNION ALL
SELECT 'Produk:', COUNT(*) FROM produk
UNION ALL
SELECT 'Users (default password):', COUNT(*) FROM users WHERE is_default_password = true;
