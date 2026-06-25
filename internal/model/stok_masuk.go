package model

import "time"

type BatchStok struct {
	ID               string    `db:"id"`
	IDProduk         string    `db:"id_produk"`
	KodeBatch        string    `db:"kode_batch"`
	ExpiredDate      time.Time `db:"expired_date"`
	StokKemasan      int       `db:"stok_kemasan"`
	TotalIsiTersedia float64   `db:"total_isi_tersedia"`
	Status           string    `db:"status"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type StokMasuk struct {
	ID                string    `db:"id"`
	IDProduk          string    `db:"id_produk"`
	IDBatch           string    `db:"id_batch"`
	IDUser            string    `db:"id_user"`
	TanggalPenerimaan time.Time `db:"tanggal_penerimaan"`
	JumlahKemasan     int       `db:"jumlah_kemasan"`
	TotalIsiMasuk     float64   `db:"total_isi_masuk"`
	Keterangan        string    `db:"keterangan"`
	CreatedAt         time.Time `db:"created_at"`
}

type StokMasukRequest struct {
	IDProduk          string `json:"id_produk" validate:"required"`
	TanggalPenerimaan string `json:"tanggal_penerimaan" validate:"required"`
	ExpiredDate       string `json:"expired_date" validate:"required"`
	JumlahKemasan     int    `json:"jumlah_kemasan" validate:"required,min=1"`
	Keterangan        string `json:"keterangan"`
}

type UpdateStokMasukRequest struct {
	TanggalPenerimaan string `json:"tanggal_penerimaan" validate:"required"`
	JumlahKemasan     int    `json:"jumlah_kemasan" validate:"required,min=1"`
	Keterangan        string `json:"keterangan"`
}

type StokMasukResponse struct {
	ID                string   `json:"id"`
	IDProduk          string   `json:"id_produk"`
	KodeProduk        string   `json:"kode_produk"`
	NamaProduk        string   `json:"nama_produk"`
	NamaKategori      string   `json:"nama_kategori"`
	PolaPenggunaan    string   `json:"pola_penggunaan"`
	SatuanIsi         string   `json:"satuan_isi"`
	IsiPerKemasan     *float64 `json:"isi_per_kemasan"`
	KodeBatch         string   `json:"kode_batch"`
	TanggalPenerimaan string   `json:"tanggal_penerimaan"`
	ExpiredDate       string   `json:"expired_date"`
	JumlahKemasan     int      `json:"jumlah_kemasan"`
	TotalIsiMasuk     float64  `json:"total_isi_masuk"`
	Keterangan        string   `json:"keterangan"`
	CreatedAt         string   `json:"created_at"`
	BatchDigunakan    bool     `json:"batch_digunakan"`
}
