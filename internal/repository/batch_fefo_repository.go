package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

// Tambah method FEFO ke BatchRepository
type BatchFEFORepository interface {
	FindBatchPrioritasFEFO(ctx context.Context, idProduk string) (*model.BatchStok, error)
	ReduceStok(ctx context.Context, id string, kurangiKemasan int, kurangiIsi float64) error
}

type batchFEFORepository struct {
	db *pgxpool.Pool
}

func NewBatchFEFORepository(db *pgxpool.Pool) BatchFEFORepository {
	return &batchFEFORepository{db: db}
}

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

func (r *batchFEFORepository) ReduceStok(ctx context.Context, id string, kurangiKemasan int, kurangiIsi float64) error {
	query := `
		UPDATE batch_stok
		SET stok_kemasan = stok_kemasan - $1,
		    total_isi_tersedia = total_isi_tersedia - $2,
		    status = CASE
		        WHEN stok_kemasan - $1 <= 0 THEN 'HABIS'
		        ELSE status
		    END,
		    updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, kurangiKemasan, kurangiIsi, id)
	return err
}
