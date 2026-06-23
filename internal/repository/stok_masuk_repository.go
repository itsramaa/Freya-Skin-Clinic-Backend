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
		       COALESCE(sm.keterangan, ''), sm.created_at
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
			&s.Keterangan, &created,
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
