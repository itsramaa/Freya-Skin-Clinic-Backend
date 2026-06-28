-- Seed Transaksi dengan Berbagai Kondisi
-- Freya Skin Clinic Backend
-- 
-- Kondisi yang di-cover:
-- 1. PARTIAL_USE: Botox (Injectable), Serum (Skincare), IPL Gel
-- 2. FULL_USE: Antibiotic Tablet, PDO Thread, Laser Tip Cover
-- 3. Batch: AKTIF, KADALUWARSA (expired)
-- 4. Kemasan Terbuka: BUD aktif, BUD expired, belum dibuka (no BUD)

BEGIN;

-- Ambil ID produk yang akan digunakan (pastikan sudah ada dari seed sebelumnya)
DO $$
DECLARE
    id_user UUID;
    
    -- PARTIAL_USE products
    id_serum UUID := 'a2222222-2222-2222-2222-222222222222';  -- Hydrating Serum
    id_botox UUID := 'b1111111-1111-1111-1111-111111111111';  -- Botox 50 Unit (akan dijadikan PARTIAL_USE)
    id_ipl_gel UUID := 'e1111111-1111-1111-1111-111111111111'; -- IPL Gel Conductivity
    id_cream UUID := 'a3333333-3333-3333-3333-333333333333';  -- Brightening Cream
    
    -- FULL_USE products
    id_tablet UUID := 'c1111111-1111-1111-1111-111111111111'; -- Antibiotic Tablet
    id_thread UUID := 'd1111111-1111-1111-1111-111111111111'; -- PDO Thread Mono
    id_laser_tip UUID := 'e3333333-3333-3333-3333-333333333333'; -- Laser Tip Cover
    
    -- Batch IDs
    batch_serum_aktif UUID := '10000001-0000-0000-0000-000000000001';
    batch_serum_aktif_2 UUID := '10000002-0000-0000-0000-000000000002';
    batch_serum_expired UUID := '10000003-0000-0000-0000-000000000003';
    batch_botox_aktif UUID := '10000004-0000-0000-0000-000000000004';
    batch_ipl_aktif UUID := '10000005-0000-0000-0000-000000000005';
    batch_ipl_aktif_2 UUID := '10000006-0000-0000-0000-000000000006';
    batch_cream_expired UUID := '10000007-0000-0000-0000-000000000007';
    batch_tablet_aktif UUID := '10000008-0000-0000-0000-000000000008';
    batch_thread_aktif UUID := '10000009-0000-0000-0000-000000000009';
    batch_laser_aktif UUID := '10000010-0000-0000-0000-000000000010';
    
    -- Stok Masuk IDs
    sm1 UUID := '20000001-0000-0000-0000-000000000001';
    sm2 UUID := '20000002-0000-0000-0000-000000000002';
    sm3 UUID := '20000003-0000-0000-0000-000000000003';
    sm4 UUID := '20000004-0000-0000-0000-000000000004';
    sm5 UUID := '20000005-0000-0000-0000-000000000005';
    sm6 UUID := '20000006-0000-0000-0000-000000000006';
    sm7 UUID := '20000007-0000-0000-0000-000000000007';
    sm8 UUID := '20000008-0000-0000-0000-000000000008';
    sm9 UUID := '20000009-0000-0000-0000-000000000009';
    sm10 UUID := '20000010-0000-0000-0000-000000000010';
    
BEGIN
    -- Ambil user admin untuk transaksi
    SELECT id INTO id_user FROM users WHERE username = 'admin' LIMIT 1;
    
    -- ========================================================================
    -- 1. BATCH STOK dengan berbagai kondisi
    -- ========================================================================
    
    -- Batch Serum AKTIF (masih bisa dipakai)
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_serum_aktif, id_serum, 'SRM-2026-01', '2027-06-30', 5, 150.000, 'AKTIF', NOW(), NOW());
    
    -- Batch Serum AKTIF 2 (batch kedua dengan expired date berbeda)
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_serum_aktif_2, id_serum, 'SRM-2026-02', '2027-07-31', 3, 88.000, 'AKTIF', NOW(), NOW());
    
    -- Batch Serum KADALUWARSA (expired)
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_serum_expired, id_serum, 'SRM-2025-01', '2026-01-31', 2, 40.000, 'KADALUWARSA', NOW(), NOW());
    
    -- Batch Botox AKTIF
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_botox_aktif, id_botox, 'BTX-2026-03', '2027-12-31', 8, 390.000, 'AKTIF', NOW(), NOW());
    
    -- Batch IPL Gel AKTIF
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_ipl_aktif, id_ipl_gel, 'IPL-2026-02', '2027-03-31', 2, 950.000, 'AKTIF', NOW(), NOW());
    
    -- Batch IPL Gel AKTIF 2 (expired date berbeda)
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_ipl_aktif_2, id_ipl_gel, 'IPL-2026-03', '2027-04-30', 1, 480.000, 'AKTIF', NOW(), NOW());
    
    -- Batch Cream KADALUWARSA
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_cream_expired, id_cream, 'CRM-2025-12', '2026-05-31', 1, 20.000, 'KADALUWARSA', NOW(), NOW());
    
    -- Batch Tablet AKTIF (FULL_USE)
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_tablet_aktif, id_tablet, 'TBL-2026-01', '2028-01-31', 15, 150.000, 'AKTIF', NOW(), NOW());
    
    -- Batch Thread AKTIF (FULL_USE)
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_thread_aktif, id_thread, 'THR-2026-04', '2027-06-30', 10, 100.000, 'AKTIF', NOW(), NOW());
    
    -- Batch Laser Tip AKTIF (FULL_USE)
    INSERT INTO batch_stok (id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at)
    VALUES (batch_laser_aktif, id_laser_tip, 'LSR-2026-05', '2028-12-31', 20, 1000.000, 'AKTIF', NOW(), NOW());
    
    -- ========================================================================
    -- 2. STOK MASUK untuk setiap batch
    -- ========================================================================
    
    INSERT INTO stok_masuk (id, id_produk, id_batch, id_user, tanggal_penerimaan, jumlah_kemasan, total_isi_masuk, keterangan, created_at)
    VALUES 
        (sm1, id_serum, batch_serum_aktif, id_user, '2026-06-01', 10, 300.000, 'Stok awal serum aktif batch 1', NOW()),
        (sm2, id_serum, batch_serum_aktif_2, id_user, '2026-06-10', 5, 150.000, 'Stok serum aktif batch 2', NOW()),
        (sm3, id_serum, batch_serum_expired, id_user, '2025-12-15', 5, 150.000, 'Batch lama yang sudah expired', NOW()),
        (sm4, id_botox, batch_botox_aktif, id_user, '2026-06-10', 10, 500.000, 'Botox stok baru', NOW()),
        (sm5, id_ipl_gel, batch_ipl_aktif, id_user, '2026-06-05', 3, 1500.000, 'IPL Gel untuk treatment batch 1', NOW()),
        (sm6, id_ipl_gel, batch_ipl_aktif_2, id_user, '2026-06-15', 2, 1000.000, 'IPL Gel batch 2', NOW()),
        (sm7, id_cream, batch_cream_expired, id_user, '2025-11-20', 3, 150.000, 'Cream batch lama', NOW()),
        (sm8, id_tablet, batch_tablet_aktif, id_user, '2026-06-20', 20, 200.000, 'Antibiotic tablet', NOW()),
        (sm9, id_thread, batch_thread_aktif, id_user, '2026-06-15', 10, 100.000, 'PDO Thread untuk prosedur', NOW()),
        (sm10, id_laser_tip, batch_laser_aktif, id_user, '2026-06-25', 50, 2500.000, 'Laser tip cover disposable', NOW());
    
    -- ========================================================================
    -- 3. KEMASAN TERBUKA (hanya untuk PARTIAL_USE products)
    -- ========================================================================
    
    -- Serum: Kemasan terbuka dengan BUD AKTIF
    INSERT INTO kemasan_terbuka (id, id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud, created_at, updated_at)
    VALUES 
        ('30000001-0000-0000-0000-000000000001', batch_serum_aktif, '2026-06-15', '2026-09-15', 30.000, 25.000, 'AKTIF', NOW(), NOW());
    
    -- Serum: Kemasan terbuka batch 2 dengan BUD AKTIF
    INSERT INTO kemasan_terbuka (id, id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud, created_at, updated_at)
    VALUES
        ('30000002-0000-0000-0000-000000000002', batch_serum_aktif_2, '2026-06-20', '2026-09-20', 30.000, 28.000, 'AKTIF', NOW(), NOW());
    
    -- Serum: Kemasan terbuka dengan BUD KADALUWARSA
    INSERT INTO kemasan_terbuka (id, id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud, created_at, updated_at)
    VALUES 
        ('30000003-0000-0000-0000-000000000003', batch_serum_expired, '2026-03-01', '2026-06-01', 30.000, 15.000, 'KADALUWARSA', NOW(), NOW());
    
    -- Botox: Kemasan terbuka BUD AKTIF
    INSERT INTO kemasan_terbuka (id, id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud, created_at, updated_at)
    VALUES 
        ('30000004-0000-0000-0000-000000000004', batch_botox_aktif, '2026-06-26', '2026-07-03', 50.000, 40.000, 'AKTIF', NOW(), NOW());
    
    -- IPL Gel: Kemasan terbuka BUD AKTIF
    INSERT INTO kemasan_terbuka (id, id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud, created_at, updated_at)
    VALUES 
        ('30000005-0000-0000-0000-000000000005', batch_ipl_aktif, '2026-06-10', '2027-06-10', 500.000, 450.000, 'AKTIF', NOW(), NOW());
    
    -- IPL Gel batch 2: Kemasan terbuka BUD AKTIF
    INSERT INTO kemasan_terbuka (id, id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud, created_at, updated_at)
    VALUES
        ('30000006-0000-0000-0000-000000000006', batch_ipl_aktif_2, '2026-06-20', '2027-06-20', 500.000, 480.000, 'AKTIF', NOW(), NOW());
    
    -- Cream: Kemasan terbuka BUD KADALUWARSA (batch juga expired)
    INSERT INTO kemasan_terbuka (id, id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud, created_at, updated_at)
    VALUES 
        ('30000007-0000-0000-0000-000000000007', batch_cream_expired, '2026-01-15', '2026-04-15', 50.000, 30.000, 'KADALUWARSA', NOW(), NOW());
    
    -- ========================================================================
    -- 4. STOK KELUAR dengan berbagai skenario
    -- ========================================================================
    
    -- PARTIAL_USE: Serum - buka kemasan baru
    INSERT INTO stok_keluar (id, id_produk, id_batch, id_kemasan_terbuka, id_user, tanggal_penggunaan, jumlah_kemasan_dipakai, jumlah_isi_dipakai, keterangan, created_at)
    VALUES 
        ('40000001-0000-0000-0000-000000000001', id_serum, batch_serum_aktif, '30000001-0000-0000-0000-000000000001', id_user, '2026-06-15', 1, 0.000, 'Buka kemasan baru', NOW()),
        ('40000002-0000-0000-0000-000000000002', id_serum, batch_serum_aktif, '30000001-0000-0000-0000-000000000001', id_user, '2026-06-15', 0, 5.000, 'Treatment pasien A', NOW());
    
    -- PARTIAL_USE: Serum - pakai dari kemasan terbuka
    INSERT INTO stok_keluar (id, id_produk, id_batch, id_kemasan_terbuka, id_user, tanggal_penggunaan, jumlah_kemasan_dipakai, jumlah_isi_dipakai, keterangan, created_at)
    VALUES 
        ('40000003-0000-0000-0000-000000000003', id_serum, batch_serum_aktif_2, '30000002-0000-0000-0000-000000000002', id_user, '2026-06-22', 0, 2.000, 'Treatment pasien B', NOW());
    
    -- PARTIAL_USE: Botox - buka dan pakai
    INSERT INTO stok_keluar (id, id_produk, id_batch, id_kemasan_terbuka, id_user, tanggal_penggunaan, jumlah_kemasan_dipakai, jumlah_isi_dipakai, keterangan, created_at)
    VALUES 
        ('40000004-0000-0000-0000-000000000004', id_botox, batch_botox_aktif, '30000004-0000-0000-0000-000000000004', id_user, '2026-06-26', 1, 0.000, 'Buka vial botox baru', NOW()),
        ('40000005-0000-0000-0000-000000000005', id_botox, batch_botox_aktif, '30000004-0000-0000-0000-000000000004', id_user, '2026-06-26', 0, 10.000, 'Injeksi pasien C (forehead)', NOW());
    
    -- PARTIAL_USE: IPL Gel - pakai dari kemasan terbuka
    INSERT INTO stok_keluar (id, id_produk, id_batch, id_kemasan_terbuka, id_user, tanggal_penggunaan, jumlah_kemasan_dipakai, jumlah_isi_dipakai, keterangan, created_at)
    VALUES 
        ('40000006-0000-0000-0000-000000000006', id_ipl_gel, batch_ipl_aktif, '30000005-0000-0000-0000-000000000005', id_user, '2026-06-12', 0, 50.000, 'IPL facial pasien D', NOW()),
        ('40000007-0000-0000-0000-000000000007', id_ipl_gel, batch_ipl_aktif_2, '30000006-0000-0000-0000-000000000006', id_user, '2026-06-22', 0, 20.000, 'Laser treatment pasien E', NOW());
    
    -- FULL_USE: Antibiotic Tablet - langsung habis per strip
    INSERT INTO stok_keluar (id, id_produk, id_batch, id_kemasan_terbuka, id_user, tanggal_penggunaan, jumlah_kemasan_dipakai, jumlah_isi_dipakai, keterangan, created_at)
    VALUES 
        ('40000008-0000-0000-0000-000000000008', id_tablet, batch_tablet_aktif, NULL, id_user, '2026-06-21', 1, 10.000, 'Resep pasien F (1 strip)', NOW()),
        ('40000009-0000-0000-0000-000000000009', id_tablet, batch_tablet_aktif, NULL, id_user, '2026-06-23', 1, 10.000, 'Resep pasien G (1 strip)', NOW());
    
    -- FULL_USE: PDO Thread - langsung habis per prosedur
    INSERT INTO stok_keluar (id, id_produk, id_batch, id_kemasan_terbuka, id_user, tanggal_penggunaan, jumlah_kemasan_dipakai, jumlah_isi_dipakai, keterangan, created_at)
    VALUES 
        ('40000010-0000-0000-0000-000000000010', id_thread, batch_thread_aktif, NULL, id_user, '2026-06-18', 1, 10.000, 'Threadlift prosedur pasien H', NOW());
    
    -- FULL_USE: Laser Tip Cover - disposable
    INSERT INTO stok_keluar (id, id_produk, id_batch, id_kemasan_terbuka, id_user, tanggal_penggunaan, jumlah_kemasan_dipakai, jumlah_isi_dipakai, keterangan, created_at)
    VALUES 
        ('40000011-0000-0000-0000-000000000011', id_laser_tip, batch_laser_aktif, NULL, id_user, '2026-06-26', 5, 5.000, 'Laser treatment 5 pasien', NOW()),
        ('40000012-0000-0000-0000-000000000012', id_laser_tip, batch_laser_aktif, NULL, id_user, '2026-06-27', 3, 3.000, 'IPL treatment 3 pasien', NOW());
    
END $$;

COMMIT;

-- Verifikasi hasil
SELECT 'Batch Stok:' AS tabel, COUNT(*) AS jumlah FROM batch_stok
UNION ALL
SELECT 'Stok Masuk:', COUNT(*) FROM stok_masuk
UNION ALL
SELECT 'Kemasan Terbuka:', COUNT(*) FROM kemasan_terbuka
UNION ALL
SELECT 'Stok Keluar:', COUNT(*) FROM stok_keluar
UNION ALL
SELECT 'Batch AKTIF:', COUNT(*) FROM batch_stok WHERE status = 'AKTIF'
UNION ALL
SELECT 'Batch KADALUWARSA:', COUNT(*) FROM batch_stok WHERE status = 'KADALUWARSA'
UNION ALL
SELECT 'Kemasan BUD AKTIF:', COUNT(*) FROM kemasan_terbuka WHERE status_bud = 'AKTIF'
UNION ALL
SELECT 'Kemasan BUD KADALUWARSA:', COUNT(*) FROM kemasan_terbuka WHERE status_bud = 'KADALUWARSA';
