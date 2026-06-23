package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

var ErrKemasanTerbukaNotFound = errors.New("kemasan terbuka tidak ditemukan")

type KemasanTerbukaRepository interface {
	FindAktifByBatch(ctx context.Context, idBatch string) (*model.KemasanTerbuka, error)
	Create(ctx context.Context, kt *model.KemasanTerbuka) error
	UpdateIsiTersisa(ctx context.Context, id string, isiTersisa float64) error
	UpdateStatus(ctx context.Context, id, status string) error
	FindExpiredBUD(ctx context.Context) ([]model.KemasanTerbuka, error)
}

type kemasanTerbukaRepository struct {
	db *pgxpool.Pool
}

func NewKemasanTerbukaRepository(db *pgxpool.Pool) KemasanTerbukaRepository {
	return &kemasanTerbukaRepository{db: db}
}

func (r *kemasanTerbukaRepository) FindAktifByBatch(ctx context.Context, idBatch string) (*model.KemasanTerbuka, error) {
	query := `SELECT id, id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud, created_at, updated_at
		FROM kemasan_terbuka WHERE id_batch = $1 AND status_bud = 'AKTIF'`
	var kt model.KemasanTerbuka
	err := r.db.QueryRow(ctx, query, idBatch).Scan(
		&kt.ID, &kt.IDBatch, &kt.TanggalDibuka, &kt.BUD,
		&kt.IsiAwal, &kt.IsiTersisa, &kt.StatusBUD, &kt.CreatedAt, &kt.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &kt, nil
}

func (r *kemasanTerbukaRepository) Create(ctx context.Context, kt *model.KemasanTerbuka) error {
	query := `
		INSERT INTO kemasan_terbuka (id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud)
		VALUES ($1, $2, $3, $4, $5, 'AKTIF')
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		kt.IDBatch, kt.TanggalDibuka, kt.BUD, kt.IsiAwal, kt.IsiTersisa,
	).Scan(&kt.ID, &kt.CreatedAt, &kt.UpdatedAt)
}

func (r *kemasanTerbukaRepository) UpdateIsiTersisa(ctx context.Context, id string, isiTersisa float64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE kemasan_terbuka SET isi_tersisa = $1, updated_at = NOW() WHERE id = $2`,
		isiTersisa, id,
	)
	return err
}

func (r *kemasanTerbukaRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE kemasan_terbuka SET status_bud = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	return err
}

func (r *kemasanTerbukaRepository) FindExpiredBUD(ctx context.Context) ([]model.KemasanTerbuka, error) {
	query := `SELECT id, id_batch, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud, created_at, updated_at
		FROM kemasan_terbuka WHERE bud < CURRENT_DATE AND status_bud = 'AKTIF'`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []model.KemasanTerbuka
	for rows.Next() {
		var kt model.KemasanTerbuka
		if err := rows.Scan(&kt.ID, &kt.IDBatch, &kt.TanggalDibuka, &kt.BUD,
			&kt.IsiAwal, &kt.IsiTersisa, &kt.StatusBUD, &kt.CreatedAt, &kt.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, kt)
	}
	return result, nil
}

// StokKeluarRepository
type StokKeluarRepository interface {
	Create(ctx context.Context, sk *model.StokKeluar) error
	FindAll(ctx context.Context) ([]model.StokKeluarResponse, error)
}

type stokKeluarRepository struct {
	db *pgxpool.Pool
}

func NewStokKeluarRepository(db *pgxpool.Pool) StokKeluarRepository {
	return &stokKeluarRepository{db: db}
}

func (r *stokKeluarRepository) Create(ctx context.Context, sk *model.StokKeluar) error {
	query := `
		INSERT INTO stok_keluar (id_produk, id_batch, id_kemasan_terbuka, id_user, tanggal_penggunaan,
		                         jumlah_kemasan_dipakai, jumlah_isi_dipakai, keterangan)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		sk.IDProduk, sk.IDBatch, sk.IDKemasanTerbuka, sk.IDUser,
		sk.TanggalPenggunaan, sk.JumlahKemasanDipakai, sk.JumlahIsiDipakai, sk.Keterangan,
	).Scan(&sk.ID, &sk.CreatedAt)
}

func (r *stokKeluarRepository) FindAll(ctx context.Context) ([]model.StokKeluarResponse, error) {
	query := `
		SELECT sk.id, sk.id_produk, p.nama_produk, b.kode_batch, p.pola_penggunaan, p.satuan_isi,
		       sk.tanggal_penggunaan, sk.jumlah_kemasan_dipakai, sk.jumlah_isi_dipakai,
		       COALESCE(sk.keterangan,''), sk.created_at
		FROM stok_keluar sk
		JOIN produk p ON p.id = sk.id_produk
		JOIN batch_stok b ON b.id = sk.id_batch
		ORDER BY sk.created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []model.StokKeluarResponse
	for rows.Next() {
		var s model.StokKeluarResponse
		var tgl, created time.Time
		if err := rows.Scan(
			&s.ID, &s.IDProduk, &s.NamaProduk, &s.KodeBatch, &s.PolaPenggunaan, &s.SatuanIsi,
			&tgl, &s.JumlahKemasanDipakai, &s.JumlahIsiDipakai, &s.Keterangan, &created,
		); err != nil {
			return nil, err
		}
		s.TanggalPenggunaan = tgl.Format("2006-01-02")
		s.CreatedAt = created.Format("2006-01-02T15:04:05Z")
		result = append(result, s)
	}
	if result == nil {
		result = []model.StokKeluarResponse{}
	}
	return result, nil
}
