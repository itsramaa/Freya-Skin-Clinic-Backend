package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

var ErrBatchNotFound = errors.New("batch tidak ditemukan")

type BatchRepository interface {
	FindByProdukAndExpired(ctx context.Context, idProduk string, expiredDate time.Time) (*model.BatchStok, error)
	Create(ctx context.Context, batch *model.BatchStok) error
	UpdateStok(ctx context.Context, id string, tambahKemasan int, tambahIsi float64) error
	FindByID(ctx context.Context, id string) (*model.BatchStok, error)
	UpdateStatus(ctx context.Context, id, status string) error
	FindExpiredBatches(ctx context.Context) ([]model.BatchStok, error)
}

type batchRepository struct {
	db *pgxpool.Pool
}

func NewBatchRepository(db *pgxpool.Pool) BatchRepository {
	return &batchRepository{db: db}
}

func (r *batchRepository) FindByProdukAndExpired(ctx context.Context, idProduk string, expiredDate time.Time) (*model.BatchStok, error) {
	query := `SELECT id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at FROM batch_stok WHERE id_produk = $1 AND expired_date = $2`
	var b model.BatchStok
	err := r.db.QueryRow(ctx, query, idProduk, expiredDate).Scan(
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

func (r *batchRepository) Create(ctx context.Context, batch *model.BatchStok) error {
	// Generate kode_batch: BCH-YYYYMM-{seq}
	var count int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM batch_stok`).Scan(&count)
	batch.KodeBatch = fmt.Sprintf("BCH-%s-%04d", batch.ExpiredDate.Format("200601"), count+1)

	query := `
		INSERT INTO batch_stok (id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status)
		VALUES ($1, $2, $3, $4, $5, 'AKTIF')
		RETURNING id, kode_batch, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		batch.IDProduk, batch.KodeBatch, batch.ExpiredDate, batch.StokKemasan, batch.TotalIsiTersedia,
	).Scan(&batch.ID, &batch.KodeBatch, &batch.CreatedAt, &batch.UpdatedAt)
}

func (r *batchRepository) UpdateStok(ctx context.Context, id string, tambahKemasan int, tambahIsi float64) error {
	query := `
		UPDATE batch_stok
		SET stok_kemasan = stok_kemasan + $1,
		    total_isi_tersedia = total_isi_tersedia + $2,
		    updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, tambahKemasan, tambahIsi, id)
	return err
}

func (r *batchRepository) FindByID(ctx context.Context, id string) (*model.BatchStok, error) {
	query := `SELECT id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at FROM batch_stok WHERE id = $1`
	var b model.BatchStok
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.IDProduk, &b.KodeBatch, &b.ExpiredDate,
		&b.StokKemasan, &b.TotalIsiTersedia, &b.Status, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBatchNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *batchRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE batch_stok SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *batchRepository) FindExpiredBatches(ctx context.Context) ([]model.BatchStok, error) {
	query := `SELECT id, id_produk, kode_batch, expired_date, stok_kemasan, total_isi_tersedia, status, created_at, updated_at FROM batch_stok WHERE expired_date < CURRENT_DATE AND status = 'AKTIF'`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []model.BatchStok
	for rows.Next() {
		var b model.BatchStok
		if err := rows.Scan(&b.ID, &b.IDProduk, &b.KodeBatch, &b.ExpiredDate, &b.StokKemasan, &b.TotalIsiTersedia, &b.Status, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}
