package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

type StokMasukRepository interface {
	Create(ctx context.Context, sm *model.StokMasuk) error
	FindAll(ctx context.Context) ([]model.StokMasukResponse, error)
	CheckBatchUsed(ctx context.Context, idBatch string) (bool, error)
	FindByID(ctx context.Context, id string) (*model.StokMasuk, error)
	Update(ctx context.Context, id string, req model.UpdateStokMasukRequest, deltaKemasan int, deltaIsi float64) error
	Delete(ctx context.Context, id string, idBatch string) error
}

type stokMasukRepository struct {
	db *pgxpool.Pool
}

func NewStokMasukRepository(db *pgxpool.Pool) StokMasukRepository {
	return &stokMasukRepository{db: db}
}

func (r *stokMasukRepository) Create(ctx context.Context, sm *model.StokMasuk) error {
	query := `
		INSERT INTO stok_masuk (id_produk, id_batch, id_user, tanggal_penerimaan, jumlah_kemasan, total_isi_masuk, keterangan)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		sm.IDProduk, sm.IDBatch, sm.IDUser,
		sm.TanggalPenerimaan, sm.JumlahKemasan, sm.TotalIsiMasuk, sm.Keterangan,
	).Scan(&sm.ID, &sm.CreatedAt)
}

func (r *stokMasukRepository) FindAll(ctx context.Context) ([]model.StokMasukResponse, error) {
	query := `
		SELECT sm.id, sm.id_produk, p.kode_produk, p.nama_produk, k.nama_kategori,
		       p.pola_penggunaan, p.satuan_isi, p.isi_per_kemasan,
		       b.kode_batch, sm.tanggal_penerimaan, b.expired_date,
		       sm.jumlah_kemasan, sm.total_isi_masuk,
		       COALESCE(sm.keterangan, ''), sm.created_at,
		       EXISTS (SELECT 1 FROM stok_keluar sk WHERE sk.id_batch = sm.id_batch) AS batch_digunakan
		FROM stok_masuk sm
		JOIN produk p ON p.id = sm.id_produk
		JOIN kategori k ON k.id = p.id_kategori
		JOIN batch_stok b ON b.id = sm.id_batch
		ORDER BY sm.created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.StokMasukResponse
	for rows.Next() {
		var s model.StokMasukResponse
		var tgl, exp, created time.Time
		if err := rows.Scan(
			&s.ID, &s.IDProduk, &s.KodeProduk, &s.NamaProduk, &s.NamaKategori,
			&s.PolaPenggunaan, &s.SatuanIsi, &s.IsiPerKemasan,
			&s.KodeBatch, &tgl, &exp,
			&s.JumlahKemasan, &s.TotalIsiMasuk,
			&s.Keterangan, &created, &s.BatchDigunakan,
		); err != nil {
			return nil, err
		}
		s.TanggalPenerimaan = tgl.Format("2006-01-02")
		s.ExpiredDate = exp.Format("2006-01-02")
		s.CreatedAt = created.Format("2006-01-02T15:04:05Z")
		result = append(result, s)
	}
	if result == nil {
		result = []model.StokMasukResponse{}
	}
	return result, nil
}

func (r *stokMasukRepository) CheckBatchUsed(ctx context.Context, idBatch string) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM stok_keluar WHERE id_batch = $1`,
		idBatch,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *stokMasukRepository) FindByID(ctx context.Context, id string) (*model.StokMasuk, error) {
	var sm model.StokMasuk
	err := r.db.QueryRow(ctx,
		`SELECT id, id_produk, id_batch, id_user, tanggal_penerimaan, jumlah_kemasan, total_isi_masuk, keterangan, created_at
		 FROM stok_masuk WHERE id = $1`,
		id,
	).Scan(&sm.ID, &sm.IDProduk, &sm.IDBatch, &sm.IDUser,
		&sm.TanggalPenerimaan, &sm.JumlahKemasan, &sm.TotalIsiMasuk,
		&sm.Keterangan, &sm.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &sm, nil
}

func (r *stokMasukRepository) Update(ctx context.Context, id string, req model.UpdateStokMasukRequest, deltaKemasan int, deltaIsi float64) error {
	tgl, err := time.Parse("2006-01-02", req.TanggalPenerimaan)
	if err != nil {
		return err
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update stok_masuk
	_, err = tx.Exec(ctx,
		`UPDATE stok_masuk SET tanggal_penerimaan=$1, jumlah_kemasan=jumlah_kemasan+$2,
		 total_isi_masuk=total_isi_masuk+$3, keterangan=$4 WHERE id=$5`,
		tgl, deltaKemasan, deltaIsi, req.Keterangan, id,
	)
	if err != nil {
		return err
	}

	// Update batch_stok: stok delta saja (expired_date tidak boleh diubah via edit stok masuk)
	_, err = tx.Exec(ctx,
		`UPDATE batch_stok SET stok_kemasan=stok_kemasan+$1, total_isi_tersedia=total_isi_tersedia+$2, updated_at=NOW()
		 WHERE id=(SELECT id_batch FROM stok_masuk WHERE id=$3)`,
		deltaKemasan, deltaIsi, id,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *stokMasukRepository) Delete(ctx context.Context, id string, idBatch string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM stok_masuk WHERE id = $1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM batch_stok WHERE id = $1`, idBatch)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
