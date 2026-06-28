-- Clean all transaction data
-- Run this before seed_transactions.sql

BEGIN;

DELETE FROM stok_keluar;
DELETE FROM kemasan_terbuka;
DELETE FROM stok_masuk;
DELETE FROM batch_stok;

COMMIT;

SELECT 'Transaction data cleaned' AS status;
