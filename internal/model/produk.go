package model

import "time"

type Produk struct {
	ID             string    `db:"id"`
	KodeProduk     string    `db:"kode_produk"`
	NamaProduk     string    `db:"nama_produk"`
	IDKategori     string    `db:"id_kategori"`
	NamaKategori   string    `db:"nama_kategori"`
	BentukKemasan  string    `db:"bentuk_kemasan"`
	SatuanIsi      string    `db:"satuan_isi"`
	IsiPerKemasan  *float64  `db:"isi_per_kemasan"`
	PolaPenggunaan string    `db:"pola_penggunaan"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type ProdukResponse struct {
	ID              string   `json:"id"`
	KodeProduk      string   `json:"kode_produk"`
	NamaProduk      string   `json:"nama_produk"`
	IDKategori      string   `json:"id_kategori"`
	NamaKategori    string   `json:"nama_kategori"`
	BentukKemasan   string   `json:"bentuk_kemasan"`
	SatuanIsi       string   `json:"satuan_isi"`
	IsiPerKemasan   *float64 `json:"isi_per_kemasan"`
	PolaPenggunaan  string   `json:"pola_penggunaan"`
	StokKemasan     int      `json:"stok_kemasan"`
	TotalIsiTersedia float64 `json:"total_isi_tersedia"`
	HasTransaksi    bool     `json:"has_transaksi"`
}

type CreateProdukRequest struct {
	NamaProduk     string   `json:"nama_produk" validate:"required"`
	IDKategori     string   `json:"id_kategori" validate:"required"`
	BentukKemasan  string   `json:"bentuk_kemasan" validate:"required"`
	SatuanIsi      string   `json:"satuan_isi" validate:"required"`
	IsiPerKemasan  *float64 `json:"isi_per_kemasan"`
	PolaPenggunaan string   `json:"pola_penggunaan" validate:"required,oneof=FULL_USE PARTIAL_USE"`
}

type UpdateProdukRequest struct {
	NamaProduk     string   `json:"nama_produk" validate:"required"`
	IDKategori     string   `json:"id_kategori" validate:"required"`
	BentukKemasan  string   `json:"bentuk_kemasan" validate:"required"`
	SatuanIsi      string   `json:"satuan_isi" validate:"required"`
	IsiPerKemasan  *float64 `json:"isi_per_kemasan"`
	PolaPenggunaan string   `json:"pola_penggunaan" validate:"required,oneof=FULL_USE PARTIAL_USE"`
}
