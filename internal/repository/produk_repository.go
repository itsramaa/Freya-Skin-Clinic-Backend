package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

var ErrProdukNotFound = errors.New("produk tidak ditemukan")

type ProdukRepository interface {
	FindAll(ctx context.Context) ([]model.ProdukResponse, error)
	FindByID(ctx context.Context, id string) (*model.Produk, error)
	Create(ctx context.Context, produk *model.Produk) error
	Update(ctx context.Context, id string, req model.UpdateProdukRequest) error
	Delete(ctx context.Context, id string) error
	HasTransaksi(ctx context.Context, id string) (bool, error)
	CountStokAktif(ctx context.Context, id string) (int, error)
	CountTransaksi(ctx context.Context, id string) (int, error)
}

type produkRepository struct {
	db *pgxpool.Pool
}

func NewProdukRepository(db *pgxpool.Pool) ProdukRepository {
	return &produkRepository{db: db}
}

func (r *produkRepository) FindAll(ctx context.Context) ([]model.ProdukResponse, error) {
	query := `
		SELECT p.id, p.kode_produk, p.nama_produk, p.id_kategori, k.nama_kategori,
		       p.bentuk_kemasan, p.satuan_isi, p.isi_per_kemasan, p.pola_penggunaan,
		       COALESCE(SUM(b.jumlah_kemasan) FILTER (WHERE b.status = 'AKTIF'), 0)::int AS stok_kemasan,
		       COALESCE(SUM(b.sisa_isi) FILTER (WHERE b.status = 'AKTIF'), 0)::float8 AS total_isi_tersedia,
		       EXISTS(SELECT 1 FROM stok_masuk sm WHERE sm.id_produk = p.id) AS has_transaksi
		FROM produk p
		JOIN kategori k ON k.id = p.id_kategori
		LEFT JOIN batch_stok b ON b.id_produk = p.id
		GROUP BY p.id, k.nama_kategori
		ORDER BY p.kode_produk
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.ProdukResponse
	for rows.Next() {
		var p model.ProdukResponse
		if err := rows.Scan(
			&p.ID, &p.KodeProduk, &p.NamaProduk, &p.IDKategori, &p.NamaKategori,
			&p.BentukKemasan, &p.SatuanIsi, &p.IsiPerKemasan, &p.PolaPenggunaan,
			&p.StokKemasan, &p.TotalIsiTersedia, &p.HasTransaksi,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	if result == nil {
		result = []model.ProdukResponse{}
	}
	return result, nil
}

func (r *produkRepository) FindByID(ctx context.Context, id string) (*model.Produk, error) {
	query := `SELECT id, kode_produk, nama_produk, id_kategori, bentuk_kemasan, satuan_isi, isi_per_kemasan, pola_penggunaan, created_at, updated_at FROM produk WHERE id = $1`
	var p model.Produk
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.KodeProduk, &p.NamaProduk, &p.IDKategori,
		&p.BentukKemasan, &p.SatuanIsi, &p.IsiPerKemasan, &p.PolaPenggunaan,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProdukNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *produkRepository) Create(ctx context.Context, produk *model.Produk) error {
	// Generate kode_produk: PRD-{3-char-prefix}-{seq}
	var count int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM produk WHERE id_kategori = $1`, produk.IDKategori).Scan(&count)

	// Ambil 3 huruf pertama nama kategori
	var namaKategori string
	_ = r.db.QueryRow(ctx, `SELECT nama_kategori FROM kategori WHERE id = $1`, produk.IDKategori).Scan(&namaKategori)
	prefix := strings.ToUpper(namaKategori)
	if len(prefix) > 3 {
		prefix = prefix[:3]
	}
	produk.KodeProduk = fmt.Sprintf("PRD-%s-%03d", prefix, count+1)

	query := `
		INSERT INTO produk (kode_produk, nama_produk, id_kategori, bentuk_kemasan, satuan_isi, isi_per_kemasan, pola_penggunaan)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, kode_produk, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		produk.KodeProduk, produk.NamaProduk, produk.IDKategori,
		produk.BentukKemasan, produk.SatuanIsi, produk.IsiPerKemasan, produk.PolaPenggunaan,
	).Scan(&produk.ID, &produk.KodeProduk, &produk.CreatedAt, &produk.UpdatedAt)
}

func (r *produkRepository) Update(ctx context.Context, id string, req model.UpdateProdukRequest) error {
	query := `
		UPDATE produk
		SET nama_produk=$1, id_kategori=$2, bentuk_kemasan=$3, satuan_isi=$4,
		    isi_per_kemasan=$5, pola_penggunaan=$6, updated_at=NOW()
		WHERE id=$7
	`
	_, err := r.db.Exec(ctx, query,
		req.NamaProduk, req.IDKategori, req.BentukKemasan, req.SatuanIsi,
		req.IsiPerKemasan, req.PolaPenggunaan, id,
	)
	return err
}

func (r *produkRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM produk WHERE id = $1`, id)
	return err
}

func (r *produkRepository) HasTransaksi(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM stok_masuk WHERE id_produk = $1)`, id).Scan(&exists)
	return exists, err
}

func (r *produkRepository) CountStokAktif(ctx context.Context, id string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM batch_stok WHERE id_produk = $1 AND status = 'AKTIF'`, id).Scan(&count)
	return count, err
}

func (r *produkRepository) CountTransaksi(ctx context.Context, id string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM stok_masuk WHERE id_produk = $1`, id).Scan(&count)
	return count, err
}
