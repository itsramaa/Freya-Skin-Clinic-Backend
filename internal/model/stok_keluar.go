package model

import "time"

type KemasanTerbuka struct {
	ID            string    `db:"id"`
	IDBatch       string    `db:"id_batch"`
	TanggalDibuka time.Time `db:"tanggal_dibuka"`
	BUD           time.Time `db:"bud"`
	IsiAwal       float64   `db:"isi_awal"`
	IsiTersisa    float64   `db:"isi_tersisa"`
	StatusBUD     string    `db:"status_bud"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type StokKeluar struct {
	ID                   string    `db:"id"`
	IDProduk             string    `db:"id_produk"`
	IDBatch              string    `db:"id_batch"`
	IDKemasanTerbuka     *string   `db:"id_kemasan_terbuka"`
	IDUser               string    `db:"id_user"`
	TanggalPenggunaan    time.Time `db:"tanggal_penggunaan"`
	JumlahKemasanDipakai int       `db:"jumlah_kemasan_dipakai"`
	JumlahIsiDipakai     float64   `db:"jumlah_isi_dipakai"`
	Keterangan           string    `db:"keterangan"`
	CreatedAt            time.Time `db:"created_at"`
}

type StokKeluarRequest struct {
	IDProduk             string  `json:"id_produk" validate:"required"`
	TanggalPenggunaan    string  `json:"tanggal_penggunaan" validate:"required"`
	JumlahKemasanDipakai int     `json:"jumlah_kemasan_dipakai"`
	JumlahIsiDipakai     float64 `json:"jumlah_isi_dipakai"`
	Keterangan           string  `json:"keterangan"`
}

type StokKeluarResponse struct {
	ID                   string  `json:"id"`
	IDProduk             string  `json:"id_produk"`
	NamaProduk           string  `json:"nama_produk"`
	KodeBatch            string  `json:"kode_batch"`
	PolaPenggunaan       string  `json:"pola_penggunaan"`
	SatuanIsi            string  `json:"satuan_isi"`
	TanggalPenggunaan    string  `json:"tanggal_penggunaan"`
	JumlahKemasanDipakai int     `json:"jumlah_kemasan_dipakai"`
	JumlahIsiDipakai     float64 `json:"jumlah_isi_dipakai"`
	Keterangan           string  `json:"keterangan"`
	CreatedAt            string  `json:"created_at"`
}

type PreviewBatchResponse struct {
	IDBatch          string              `json:"id_batch"`
	KodeBatch        string              `json:"kode_batch"`
	ExpiredDate      string              `json:"expired_date"`
	StokKemasan      int                 `json:"stok_kemasan"`
	TotalIsiTersedia float64             `json:"total_isi_tersedia"`
	PolaPenggunaan   string              `json:"pola_penggunaan"`
	SatuanIsi        string              `json:"satuan_isi"`
	IsiPerKemasan    *float64            `json:"isi_per_kemasan"`
	KemasanTerbuka   *KemasanTerbukaInfo `json:"kemasan_terbuka"`
}

type KemasanTerbukaInfo struct {
	ID         string  `json:"id"`
	BUD        string  `json:"bud"`
	IsiTersisa float64 `json:"isi_tersisa"`
	StatusBUD  string  `json:"status_bud"`
}
