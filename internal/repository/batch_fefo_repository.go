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
	ReduceStok(ctx context.Context, id string, kurangiKemasan int, kurangiIsi float64) error
}

type batchFEFORepository struct {
	db *pgxpool.Pool
}

func NewBatchFEFORepository(db *pgxpool.Pool) BatchFEFORepository {
	return &batchFEFORepository{db: db}
}

// FindBatchPrioritasFEFO — untuk Full Use: cari batch AKTIF dengan stok_kemasan > 0, expired ASC
func (r *batchFEFORepository) FindBatchPrioritasFEFO(ctx context.Context, idProduk string) (*model.BatchStok, error) {
	query := `
		SELECT id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at
		FROM batch_stok
		WHERE id_produk = $1 AND status = 'AKTIF' AND stok_kemasan > 0
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
// Prioritas 1: batch AKTIF yang punya kemasan terbuka AKTIF (expired ASC)
// Prioritas 2: batch AKTIF yang stok_kemasan > 0 (expired ASC)
func (r *batchFEFORepository) FindBatchPartialUseFEFO(ctx context.Context, idProduk string) (*model.BatchStok, error) {
	// Prioritas 1: ada kemasan terbuka aktif
	query1 := `
		SELECT b.id, b.id_produk, b.kode_batch, b.expired_date, b.stok_kemasan, b.total_isi_tersedia, b.status, b.created_at, b.updated_at
		FROM batch_stok b
		JOIN kemasan_terbuka kt ON kt.id_batch = b.id AND kt.status_bud = 'AKTIF' AND kt.isi_tersisa > 0
		WHERE b.id_produk = $1 AND b.status = 'AKTIF'
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

func (r *batchFEFORepository) ReduceStok(ctx context.Context, id string, kurangiKemasan int, kurangiIsi float64) error {
	query := `
		UPDATE batch_stok
		SET stok_kemasan = stok_kemasan - $1,
		    total_isi_tersedia = total_isi_tersedia - $2,
		    status = CASE
		        WHEN (stok_kemasan - $1) <= 0 AND NOT EXISTS (
		            SELECT 1 FROM kemasan_terbuka WHERE id_batch = id AND status_bud = 'AKTIF' AND isi_tersisa > 0
		        ) THEN 'HABIS'
		        ELSE status
		    END,
		    updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, kurangiKemasan, kurangiIsi, id)
	return err
}
