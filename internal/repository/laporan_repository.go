package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

type LaporanRepository interface {
	GetStokMasuk(ctx context.Context, filter model.LaporanFilter) ([]model.LaporanStokMasukItem, error)
	GetStokKeluar(ctx context.Context, filter model.LaporanFilter) ([]model.LaporanStokKeluarItem, error)
	GetSisaStok(ctx context.Context, filter model.LaporanFilter) ([]model.LaporanSisaStokItem, error)
}

type laporanRepository struct {
	db *pgxpool.Pool
}

func NewLaporanRepository(db *pgxpool.Pool) LaporanRepository {
	return &laporanRepository{db: db}
}

func (r *laporanRepository) GetStokMasuk(ctx context.Context, filter model.LaporanFilter) ([]model.LaporanStokMasukItem, error) {
	query := `
		SELECT sm.id, sm.tanggal_penerimaan, p.nama_produk, p.kode_produk, k.nama_kategori,
		       p.pola_penggunaan, p.satuan_isi,
		       b.kode_batch, b.expired_date, sm.jumlah_kemasan, sm.total_isi_masuk,
		       COALESCE(sm.keterangan,'')
		FROM stok_masuk sm
		JOIN produk p ON p.id = sm.id_produk
		JOIN kategori k ON k.id = p.id_kategori
		JOIN batch_stok b ON b.id = sm.id_batch
		WHERE sm.tanggal_penerimaan BETWEEN $1 AND $2
	`
	args := []interface{}{filter.Dari, filter.Sampai}
	idx := 3

	if filter.KategoriID != "" {
		query += ` AND p.id_kategori = $` + strconv.Itoa(idx)
		args = append(args, filter.KategoriID)
		idx++
	}
	if filter.ProdukID != "" {
		query += ` AND p.id = $` + strconv.Itoa(idx)
		args = append(args, filter.ProdukID)
	}
	query += ` ORDER BY sm.tanggal_penerimaan DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.LaporanStokMasukItem
	for rows.Next() {
		var item model.LaporanStokMasukItem
		var tgl, exp time.Time
		if err := rows.Scan(&item.ID, &tgl, &item.NamaProduk, &item.KodeProduk, &item.NamaKategori,
			&item.PolaPenggunaan, &item.SatuanIsi,
			&item.KodeBatch, &exp, &item.JumlahKemasan, &item.TotalIsiMasuk, &item.Keterangan); err != nil {
			return nil, err
		}
		item.TanggalPenerimaan = tgl.Format("2006-01-02")
		item.ExpiredDate = exp.Format("2006-01-02")
		result = append(result, item)
	}
	if result == nil {
		result = []model.LaporanStokMasukItem{}
	}
	return result, nil
}

func (r *laporanRepository) GetStokKeluar(ctx context.Context, filter model.LaporanFilter) ([]model.LaporanStokKeluarItem, error) {
	query := `
		SELECT sk.id, sk.tanggal_penggunaan, p.nama_produk, p.kode_produk, k.nama_kategori,
		       b.kode_batch, p.pola_penggunaan, p.satuan_isi, sk.jumlah_kemasan_dipakai, sk.jumlah_isi_dipakai,
		       COALESCE(sk.keterangan,'')
		FROM stok_keluar sk
		JOIN produk p ON p.id = sk.id_produk
		JOIN kategori k ON k.id = p.id_kategori
		JOIN batch_stok b ON b.id = sk.id_batch
		WHERE sk.tanggal_penggunaan BETWEEN $1 AND $2
	`
	args := []interface{}{filter.Dari, filter.Sampai}
	idx := 3

	if filter.KategoriID != "" {
		query += ` AND p.id_kategori = $` + strconv.Itoa(idx)
		args = append(args, filter.KategoriID)
		idx++
	}
	if filter.ProdukID != "" {
		query += ` AND p.id = $` + strconv.Itoa(idx)
		args = append(args, filter.ProdukID)
	}
	query += ` ORDER BY sk.tanggal_penggunaan DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.LaporanStokKeluarItem
	for rows.Next() {
		var item model.LaporanStokKeluarItem
		var tgl time.Time
		if err := rows.Scan(&item.ID, &tgl, &item.NamaProduk, &item.KodeProduk, &item.NamaKategori,
			&item.KodeBatch, &item.PolaPenggunaan, &item.SatuanIsi, &item.JumlahKemasanDipakai, &item.JumlahIsiDipakai, &item.Keterangan); err != nil {
			return nil, err
		}
		item.TanggalPenggunaan = tgl.Format("2006-01-02")
		result = append(result, item)
	}
	if result == nil {
		result = []model.LaporanStokKeluarItem{}
	}
	return result, nil
}

func (r *laporanRepository) GetSisaStok(ctx context.Context, filter model.LaporanFilter) ([]model.LaporanSisaStokItem, error) {
	query := `
		SELECT p.kode_produk, p.nama_produk, k.nama_kategori, p.pola_penggunaan, p.satuan_isi,
		       COALESCE(SUM(b.stok_kemasan) FILTER (WHERE b.status = 'AKTIF'), 0)::int AS total_stok,
		       COALESCE(SUM(b.total_isi_tersedia) FILTER (WHERE b.status = 'AKTIF'), 0)::float8 AS total_isi
		FROM produk p
		JOIN kategori k ON k.id = p.id_kategori
		LEFT JOIN batch_stok b ON b.id_produk = p.id
		WHERE 1=1
	`
	args := []interface{}{}
	idx := 1

	if filter.KategoriID != "" {
		query += ` AND p.id_kategori = $` + strconv.Itoa(idx)
		args = append(args, filter.KategoriID)
		idx++
	}
	if filter.ProdukID != "" {
		query += ` AND p.id = $` + strconv.Itoa(idx)
		args = append(args, filter.ProdukID)
		idx++
	}
	if !filter.Dari.IsZero() {
		query += ` AND (b.expired_date >= $` + strconv.Itoa(idx) + ` OR b.id IS NULL)`
		args = append(args, filter.Dari)
		idx++
	}
	if !filter.Sampai.IsZero() {
		query += ` AND (b.expired_date <= $` + strconv.Itoa(idx) + ` OR b.id IS NULL)`
		args = append(args, filter.Sampai)
		idx++
	}
	query += ` GROUP BY p.id, k.nama_kategori ORDER BY p.kode_produk`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.LaporanSisaStokItem
	for rows.Next() {
		var item model.LaporanSisaStokItem
		if err := rows.Scan(&item.KodeProduk, &item.NamaProduk, &item.NamaKategori,
			&item.PolaPenggunaan, &item.SatuanIsi, &item.TotalStok, &item.TotalIsi); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if result == nil {
		result = []model.LaporanSisaStokItem{}
	}
	return result, nil
}
