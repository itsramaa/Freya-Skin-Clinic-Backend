package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

var ErrKategoriNotFound = errors.New("kategori tidak ditemukan")

type KategoriRepository interface {
	FindAll(ctx context.Context) ([]model.KategoriResponse, error)
	FindByID(ctx context.Context, id string) (*model.Kategori, error)
	FindByNama(ctx context.Context, nama string) (*model.Kategori, error)
	Create(ctx context.Context, kategori *model.Kategori) error
	Update(ctx context.Context, id, namaKategori string) (*model.KategoriResponse, error)
	Delete(ctx context.Context, id string) error
	CountProdukByKategoriID(ctx context.Context, id string) (int, error)
}

type kategoriRepository struct {
	db *pgxpool.Pool
}

func NewKategoriRepository(db *pgxpool.Pool) KategoriRepository {
	return &kategoriRepository{db: db}
}

func (r *kategoriRepository) FindAll(ctx context.Context) ([]model.KategoriResponse, error) {
	query := `
		SELECT k.id, k.kode_kategori, k.nama_kategori,
		       COUNT(p.id) AS jumlah_produk
		FROM kategori k
		LEFT JOIN produk p ON p.id_kategori = k.id
		GROUP BY k.id
		ORDER BY k.kode_kategori
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.KategoriResponse
	for rows.Next() {
		var k model.KategoriResponse
		if err := rows.Scan(&k.ID, &k.KodeKategori, &k.NamaKategori, &k.JumlahProduk); err != nil {
			return nil, err
		}
		result = append(result, k)
	}
	if result == nil {
		result = []model.KategoriResponse{}
	}
	return result, nil
}

func (r *kategoriRepository) FindByID(ctx context.Context, id string) (*model.Kategori, error) {
	query := `SELECT id, kode_kategori, nama_kategori, created_at, updated_at FROM kategori WHERE id = $1`
	var k model.Kategori
	err := r.db.QueryRow(ctx, query, id).Scan(&k.ID, &k.KodeKategori, &k.NamaKategori, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKategoriNotFound
		}
		return nil, err
	}
	return &k, nil
}

func (r *kategoriRepository) FindByNama(ctx context.Context, nama string) (*model.Kategori, error) {
	query := `SELECT id, kode_kategori, nama_kategori, created_at, updated_at FROM kategori WHERE LOWER(nama_kategori) = LOWER($1)`
	var k model.Kategori
	err := r.db.QueryRow(ctx, query, nama).Scan(&k.ID, &k.KodeKategori, &k.NamaKategori, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &k, nil
}

func (r *kategoriRepository) Create(ctx context.Context, kategori *model.Kategori) error {
	// Generate kode_kategori: KTG-XXX
	var count int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM kategori`).Scan(&count)
	kategori.KodeKategori = fmt.Sprintf("KTG-%03d", count+1)

	query := `
		INSERT INTO kategori (kode_kategori, nama_kategori)
		VALUES ($1, $2)
		RETURNING id, kode_kategori, nama_kategori, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, kategori.KodeKategori, kategori.NamaKategori).
		Scan(&kategori.ID, &kategori.KodeKategori, &kategori.NamaKategori, &kategori.CreatedAt, &kategori.UpdatedAt)
}

func (r *kategoriRepository) Update(ctx context.Context, id, namaKategori string) (*model.KategoriResponse, error) {
	query := `
		UPDATE kategori SET nama_kategori = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, kode_kategori, nama_kategori
	`
	var k model.KategoriResponse
	err := r.db.QueryRow(ctx, query, namaKategori, id).Scan(&k.ID, &k.KodeKategori, &k.NamaKategori)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *kategoriRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM kategori WHERE id = $1`, id)
	return err
}

func (r *kategoriRepository) CountProdukByKategoriID(ctx context.Context, id string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM produk WHERE id_kategori = $1`, id).Scan(&count)
	return count, err
}
