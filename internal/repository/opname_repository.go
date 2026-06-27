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
	GetDetailItems(ctx context.Context, idOpname string) ([]model.OpnameItemResponse, error)
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
	// Satu baris per batch:
	// - FULL_USE: stok_sistem = stok_kemasan, tanpa kemasan terbuka
	// - PARTIAL_USE: stok_sistem = stok_kemasan (kolom Stok Fisik opname),
	//   isi_tersisa dari kemasan terbuka aktif jika ada (kolom Sisa Isi Terbuka Fisik opname)
	query := `
		SELECT
			b.id,
			b.kode_batch,
			p.nama_produk,
			b.expired_date,
			p.pola_penggunaan,
			p.satuan_isi,
			kt.id            AS id_kemasan_terbuka,
			kt.isi_tersisa   AS isi_tersisa,
			b.stok_kemasan::DECIMAL AS stok_sistem,
			CASE WHEN p.pola_penggunaan = 'PARTIAL_USE'
			     THEN b.stok_kemasan::DECIMAL
			     ELSE NULL
			END AS stok_kemasan_sistem
		FROM batch_stok b
		JOIN produk p ON p.id = b.id_produk
		LEFT JOIN kemasan_terbuka kt
			ON kt.id_batch = b.id AND kt.status_bud = 'AKTIF' AND kt.isi_tersisa > 0
		WHERE b.status = 'AKTIF'
		  AND (
		    b.stok_kemasan > 0
		    OR (p.pola_penggunaan = 'PARTIAL_USE' AND kt.id IS NOT NULL)
		  )
		ORDER BY p.nama_produk, b.expired_date
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
			&item.StokKemasanSistem,
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

func (r *opnameRepository) GetDetailItems(ctx context.Context, idOpname string) ([]model.OpnameItemResponse, error) {
	// Baca dari detail_opname (histori audit) — dipakai untuk opname SELESAI/DIBATALKAN
	query := `
		SELECT d.id_batch, b.kode_batch, p.nama_produk, b.expired_date,
		       p.pola_penggunaan, p.satuan_isi,
		       d.id_kemasan_terbuka,
		       d.sisa_isi_sistem AS isi_tersisa,
		       d.stok_sistem, d.stok_fisik, d.selisih,
		       d.sisa_isi_sistem, d.sisa_isi_fisik, d.selisih_sisa_isi,
		       COALESCE(d.keterangan, '')
		FROM detail_opname d
		JOIN batch_stok b ON b.id = d.id_batch
		JOIN produk p ON p.id = b.id_produk
		WHERE d.id_opname = $1
		ORDER BY p.nama_produk, b.expired_date
	`
	rows, err := r.db.Query(ctx, query, idOpname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.OpnameItemResponse
	for rows.Next() {
		var item model.OpnameItemResponse
		var exp time.Time
		var stokFisik, selisih float64
		var keterangan string
		if err := rows.Scan(
			&item.IDBatch, &item.KodeBatch, &item.NamaProduk, &exp,
			&item.PolaPenggunaan, &item.SatuanIsi,
			&item.IDKemasanTerbuka,
			&item.IsiTersisa,
			&item.StokSistem, &stokFisik, &selisih,
			&item.SisaIsiSistem, &item.SisaIsiFisik, &item.SelisihSisaIsi,
			&keterangan,
		); err != nil {
			return nil, err
		}
		item.ExpiredDate = exp.Format("2006-01-02")
		item.StokFisik = &stokFisik
		item.Selisih = &selisih
		item.Keterangan = keterangan
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
		// Selalu ambil stok_kemasan dari batch_stok sebagai basis selisih kemasan
		var stokKemasanSistem float64
		err = tx.QueryRow(ctx, `SELECT stok_kemasan FROM batch_stok WHERE id = $1`, d.IDBatch).Scan(&stokKemasanSistem)
		if err != nil {
			return err
		}

		selisihKemasan := d.StokFisik - stokKemasanSistem

		// Untuk kemasan terbuka: cek juga selisih sisa isi
		var isiTersisaSistem float64
		var selisihSisaIsi float64
		if d.IDKemasanTerbuka != nil && d.SisaIsiTerbuka != nil {
			err = tx.QueryRow(ctx, `SELECT isi_tersisa FROM kemasan_terbuka WHERE id = $1`, *d.IDKemasanTerbuka).Scan(&isiTersisaSistem)
			if err != nil {
				return err
			}
			selisihSisaIsi = *d.SisaIsiTerbuka - isiTersisaSistem
		}

		hasSelisih := selisihKemasan != 0 || selisihSisaIsi != 0

		// Validasi keterangan wajib jika ada selisih
		if hasSelisih && (d.Keterangan == nil || *d.Keterangan == "") {
			return ErrKeteranganWajib
		}

		// Simpan ke detail_opname dengan kolom sisa isi untuk audit trail lengkap
		var sisaIsiSistemPtr *float64
		var sisaIsiFisikPtr *float64
		var selisihSisaIsiPtr *float64
		if d.IDKemasanTerbuka != nil && d.SisaIsiTerbuka != nil {
			sisaIsiSistemPtr = &isiTersisaSistem
			sisaIsiFisikPtr = d.SisaIsiTerbuka
			selisihSisaIsiPtr = &selisihSisaIsi
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO detail_opname (id_opname, id_batch, id_kemasan_terbuka, stok_sistem, stok_fisik, selisih, sisa_isi_sistem, sisa_isi_fisik, selisih_sisa_isi, keterangan)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			idOpname, d.IDBatch, d.IDKemasanTerbuka,
			stokKemasanSistem, d.StokFisik, selisihKemasan,
			sisaIsiSistemPtr, sisaIsiFisikPtr, selisihSisaIsiPtr,
			d.Keterangan,
		)
		if err != nil {
			return err
		}

		// Penyesuaian stok
		if d.IDKemasanTerbuka != nil {
			// PARTIAL_USE dengan kemasan terbuka
			if d.SisaIsiTerbuka != nil && selisihSisaIsi != 0 {
				_, err = tx.Exec(ctx, `UPDATE kemasan_terbuka SET isi_tersisa = $1, updated_at = NOW() WHERE id = $2`, *d.SisaIsiTerbuka, *d.IDKemasanTerbuka)
				if err != nil {
					return err
				}
			}
			if selisihKemasan != 0 {
				newStok := int(d.StokFisik)
				_, err = tx.Exec(ctx, `
					UPDATE batch_stok SET
						stok_kemasan = $1,
						status = CASE WHEN $2 = 0 THEN 'HABIS' ELSE status END,
						updated_at = NOW()
					WHERE id = $3`, newStok, newStok, d.IDBatch)
				if err != nil {
					return err
				}
			}
			// Hitung ulang total_isi_tersedia untuk PARTIAL_USE:
			// = (stok_kemasan * isi_per_kemasan) + sisa isi kemasan terbuka aktif
			_, err = tx.Exec(ctx, `
				UPDATE batch_stok SET
					total_isi_tersedia = (
						SELECT (b2.stok_kemasan * COALESCE(p.isi_per_kemasan, 0))
						       + COALESCE((SELECT kt2.isi_tersisa FROM kemasan_terbuka kt2
						                   WHERE kt2.id_batch = b2.id AND kt2.status_bud = 'AKTIF'
						                   LIMIT 1), 0)
						FROM batch_stok b2
						JOIN produk p ON p.id = b2.id_produk
						WHERE b2.id = $1
					),
					updated_at = NOW()
				WHERE id = $1`, d.IDBatch)
			if err != nil {
				return err
			}
		} else if selisihKemasan != 0 {
			// FULL_USE: update stok_kemasan dan total_isi_tersedia = stok_fisik * isi_per_kemasan
			newStok := int(d.StokFisik)
			_, err = tx.Exec(ctx, `
				UPDATE batch_stok SET
					stok_kemasan = $1,
					total_isi_tersedia = $2 * COALESCE((SELECT isi_per_kemasan FROM produk WHERE id = batch_stok.id_produk), 1),
					status = CASE WHEN $3 = 0 THEN 'HABIS' ELSE status END,
					updated_at = NOW()
				WHERE id = $4`, newStok, newStok, newStok, d.IDBatch)
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
