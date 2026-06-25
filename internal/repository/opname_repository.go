package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

var ErrOpnameNotFound = errors.New("sesi opname tidak ditemukan")
var ErrKeteranganWajib = errors.New("keterangan wajib diisi untuk item yang memiliki selisih")

type OpnameRepository interface {
	Create(ctx context.Context, op *model.StokOpname) error
	FindAll(ctx context.Context) ([]model.StokOpnameResponse, error)
	FindByID(ctx context.Context, id string) (*model.StokOpname, error)
	FindAktif(ctx context.Context) (*model.StokOpname, error)
	UpdateStatus(ctx context.Context, id, status string) error
	GetItemsForOpname(ctx context.Context) ([]model.OpnameItemResponse, error)
	SaveDetailAndAdjust(ctx context.Context, idOpname string, details []model.DetailOpnameInput) error
}

type opnameRepository struct {
	db *pgxpool.Pool
}

func NewOpnameRepository(db *pgxpool.Pool) OpnameRepository {
	return &opnameRepository{db: db}
}

func (r *opnameRepository) Create(ctx context.Context, op *model.StokOpname) error {
	query := `INSERT INTO stok_opname (id_user, tanggal_opname, status) VALUES ($1, $2, 'AKTIF') RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, op.IDUser, op.TanggalOpname).Scan(&op.ID, &op.CreatedAt, &op.UpdatedAt)
}

func (r *opnameRepository) FindAktif(ctx context.Context) (*model.StokOpname, error) {
	query := `SELECT id, id_user, tanggal_opname, status, COALESCE(catatan,''), created_at, updated_at FROM stok_opname WHERE status = 'AKTIF' ORDER BY created_at DESC LIMIT 1`
	var op model.StokOpname
	var tgl, created, updated time.Time
	var catatan string
	err := r.db.QueryRow(ctx, query).Scan(&op.ID, &op.IDUser, &tgl, &op.Status, &catatan, &created, &updated)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	op.TanggalOpname = tgl
	op.CreatedAt = created
	op.UpdatedAt = updated
	return &op, nil
}

func (r *opnameRepository) FindAll(ctx context.Context) ([]model.StokOpnameResponse, error) {
	rows, err := r.db.Query(ctx, `SELECT id, id_user, tanggal_opname, status, COALESCE(catatan,''), created_at FROM stok_opname ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []model.StokOpnameResponse
	for rows.Next() {
		var s model.StokOpnameResponse
		var tgl, created time.Time
		if err := rows.Scan(&s.ID, &s.IDUser, &tgl, &s.Status, &s.Catatan, &created); err != nil {
			return nil, err
		}
		s.TanggalOpname = tgl.Format("2006-01-02")
		s.CreatedAt = created.Format("2006-01-02T15:04:05Z")
		result = append(result, s)
	}
	if result == nil {
		result = []model.StokOpnameResponse{}
	}
	return result, nil
}

func (r *opnameRepository) FindByID(ctx context.Context, id string) (*model.StokOpname, error) {
	var op model.StokOpname
	err := r.db.QueryRow(ctx, `SELECT id, id_user, tanggal_opname, status, COALESCE(catatan,''), created_at, updated_at FROM stok_opname WHERE id = $1`, id).
		Scan(&op.ID, &op.IDUser, &op.TanggalOpname, &op.Status, &op.Catatan, &op.CreatedAt, &op.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOpnameNotFound
		}
		return nil, err
	}
	return &op, nil
}

func (r *opnameRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE stok_opname SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

func (r *opnameRepository) GetItemsForOpname(ctx context.Context) ([]model.OpnameItemResponse, error) {
	query := `
		SELECT b.id, b.kode_batch, p.nama_produk, b.expired_date,
		       p.pola_penggunaan, p.satuan_isi,
		       NULL::UUID AS id_kemasan_terbuka, NULL::DECIMAL AS isi_tersisa,
		       b.stok_kemasan::DECIMAL AS stok_sistem
		FROM batch_stok b
		JOIN produk p ON p.id = b.id_produk
		WHERE b.status = 'AKTIF' AND b.stok_kemasan > 0

		UNION ALL

		SELECT b.id, b.kode_batch, p.nama_produk, b.expired_date,
		       p.pola_penggunaan, p.satuan_isi,
		       kt.id, kt.isi_tersisa,
		       kt.isi_tersisa AS stok_sistem
		FROM kemasan_terbuka kt
		JOIN batch_stok b ON b.id = kt.id_batch
		JOIN produk p ON p.id = b.id_produk
		WHERE kt.status_bud = 'AKTIF' AND kt.isi_tersisa > 0

		ORDER BY nama_produk, expired_date
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.OpnameItemResponse
	for rows.Next() {
		var item model.OpnameItemResponse
		var exp time.Time
		if err := rows.Scan(
			&item.IDBatch, &item.KodeBatch, &item.NamaProduk, &exp,
			&item.PolaPenggunaan, &item.SatuanIsi,
			&item.IDKemasanTerbuka, &item.IsiTersisa, &item.StokSistem,
		); err != nil {
			return nil, err
		}
		item.ExpiredDate = exp.Format("2006-01-02")
		result = append(result, item)
	}
	if result == nil {
		result = []model.OpnameItemResponse{}
	}
	return result, nil
}

func (r *opnameRepository) SaveDetailAndAdjust(ctx context.Context, idOpname string, details []model.DetailOpnameInput) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, d := range details {
		// Ambil stok_sistem dari batch atau kemasan_terbuka
		var stokSistem float64
		if d.IDKemasanTerbuka != nil {
			err = tx.QueryRow(ctx, `SELECT isi_tersisa FROM kemasan_terbuka WHERE id = $1`, *d.IDKemasanTerbuka).Scan(&stokSistem)
		} else {
			err = tx.QueryRow(ctx, `SELECT stok_kemasan FROM batch_stok WHERE id = $1`, d.IDBatch).Scan(&stokSistem)
		}
		if err != nil {
			return err
		}

		selisih := d.StokFisik - stokSistem

		// Validasi keterangan wajib jika ada selisih
		if selisih != 0 && d.Keterangan == "" {
			return ErrKeteranganWajib
		}

		// Insert detail_opname
		_, err = tx.Exec(ctx, `
			INSERT INTO detail_opname (id_opname, id_batch, id_kemasan_terbuka, stok_sistem, stok_fisik, selisih, keterangan)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			idOpname, d.IDBatch, d.IDKemasanTerbuka, stokSistem, d.StokFisik, selisih, d.Keterangan,
		)
		if err != nil {
			return err
		}

		// Penyesuaian stok jika ada selisih
		if selisih != 0 {
			if d.IDKemasanTerbuka != nil {
				// Partial use: gunakan SisaIsiTerbuka jika ada, fallback ke StokFisik
				sisaIsi := d.StokFisik
				if d.SisaIsiTerbuka != nil {
					sisaIsi = *d.SisaIsiTerbuka
				}
				_, err = tx.Exec(ctx, `UPDATE kemasan_terbuka SET isi_tersisa = $1, updated_at = NOW() WHERE id = $2`, sisaIsi, *d.IDKemasanTerbuka)
			} else {
				_, err = tx.Exec(ctx, `UPDATE batch_stok SET stok_kemasan = $1, updated_at = NOW() WHERE id = $2`, int(d.StokFisik), d.IDBatch)
			}
			if err != nil {
				return err
			}
		}
	}

	// Update status opname menjadi SELESAI
	_, err = tx.Exec(ctx, `UPDATE stok_opname SET status = 'SELESAI', updated_at = NOW() WHERE id = $1`, idOpname)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
