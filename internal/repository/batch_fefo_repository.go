package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

type BatchFEFORepository interface {
	FindBatchPrioritasFEFO(ctx context.Context, idProduk string) (*model.BatchStok, error)
	FindBatchPartialUseFEFO(ctx context.Context, idProduk string) (*model.BatchStok, error)
	FindAllBatchFEFO(ctx context.Context, idProduk string) ([]model.BatchStok, error)
	ReduceStok(ctx context.Context, id string, kurangiKemasan int, kurangiIsi float64) error
}

type batchFEFORepository struct {
	db *pgxpool.Pool
}

func NewBatchFEFORepository(db *pgxpool.Pool) BatchFEFORepository {
	return &batchFEFORepository{db: db}
}

// FindBatchPrioritasFEFO — untuk Full Use: cari batch AKTIF dengan stok_kemasan > 0, belum expired, expired ASC
func (r *batchFEFORepository) FindBatchPrioritasFEFO(ctx context.Context, idProduk string) (*model.BatchStok, error) {
	query := `
		SELECT id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at
		FROM batch_stok
		WHERE id_produk = $1 AND status = 'AKTIF' AND stok_kemasan > 0 AND expired_date >= CURRENT_DATE
		ORDER BY expired_date ASC
		LIMIT 1
	`
	var b model.BatchStok
	err := r.db.QueryRow(ctx, query, idProduk).Scan(
		&b.ID, &b.IDProduk, &b.KodeBatch, &b.ExpiredDate,
		&b.StokKemasan, &b.TotalIsiTersedia, &b.Status, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// FindBatchPartialUseFEFO — untuk Partial Use:
// Prioritas 1: batch AKTIF yang punya kemasan terbuka AKTIF (expired ASC, belum expired)
// Prioritas 2: batch AKTIF yang stok_kemasan > 0 (expired ASC, belum expired)
func (r *batchFEFORepository) FindBatchPartialUseFEFO(ctx context.Context, idProduk string) (*model.BatchStok, error) {
	// Prioritas 1: ada kemasan terbuka aktif, BUD belum kadaluwarsa, batch belum expired
	query1 := `
		SELECT b.id, b.id_produk, b.kode_batch, b.expired_date, b.stok_kemasan, b.total_isi_tersedia, b.status, b.created_at, b.updated_at
		FROM batch_stok b
		JOIN kemasan_terbuka kt ON kt.id_batch = b.id AND kt.status_bud = 'AKTIF' AND kt.isi_tersisa > 0 AND kt.bud >= CURRENT_DATE
		WHERE b.id_produk = $1 AND b.status = 'AKTIF' AND b.expired_date >= CURRENT_DATE
		ORDER BY b.expired_date ASC
		LIMIT 1
	`
	var b model.BatchStok
	err := r.db.QueryRow(ctx, query1, idProduk).Scan(
		&b.ID, &b.IDProduk, &b.KodeBatch, &b.ExpiredDate,
		&b.StokKemasan, &b.TotalIsiTersedia, &b.Status, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == nil {
		return &b, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Prioritas 2: batch dengan stok kemasan > 0
	return r.FindBatchPrioritasFEFO(ctx, idProduk)
}

// FindAllBatchFEFO — untuk batch splitting Full Use: ambil semua batch AKTIF dengan stok_kemasan > 0, belum expired, expired ASC
func (r *batchFEFORepository) FindAllBatchFEFO(ctx context.Context, idProduk string) ([]model.BatchStok, error) {
	query := `
		SELECT id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at
		FROM batch_stok
		WHERE id_produk = $1 AND status = 'AKTIF' AND stok_kemasan > 0 AND expired_date >= CURRENT_DATE
		ORDER BY expired_date ASC
	`
	rows, err := r.db.Query(ctx, query, idProduk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batches []model.BatchStok
	for rows.Next() {
		var b model.BatchStok
		if err := rows.Scan(
			&b.ID, &b.IDProduk, &b.KodeBatch, &b.ExpiredDate,
			&b.StokKemasan, &b.TotalIsiTersedia, &b.Status, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		batches = append(batches, b)
	}
	return batches, rows.Err()
}

func (r *batchFEFORepository) ReduceStok(ctx context.Context, id string, kurangiKemasan int, kurangiIsi float64) error {
	query := `
		UPDATE batch_stok bs
		SET stok_kemasan = bs.stok_kemasan - $1,
		    total_isi_tersedia = bs.total_isi_tersedia - $2,
		    status = CASE
		        WHEN (bs.stok_kemasan - $1) <= 0 AND NOT EXISTS (
		            SELECT 1 FROM kemasan_terbuka kt
		            WHERE kt.id_batch = bs.id AND kt.status_bud = 'AKTIF' AND kt.isi_tersisa > 0
		        ) THEN 'HABIS'
		        ELSE bs.status
		    END,
		    updated_at = NOW()
		WHERE bs.id = $3
	`
	_, err := r.db.Exec(ctx, query, kurangiKemasan, kurangiIsi, id)
	return err
}
