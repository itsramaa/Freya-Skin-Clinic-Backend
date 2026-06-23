package model

import "time"

type Kategori struct {
	ID           string    `db:"id"`
	KodeKategori string    `db:"kode_kategori"`
	NamaKategori string    `db:"nama_kategori"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type KategoriResponse struct {
	ID           string `json:"id"`
	KodeKategori string `json:"kode_kategori"`
	NamaKategori string `json:"nama_kategori"`
	JumlahProduk int    `json:"jumlah_produk"`
}

type CreateKategoriRequest struct {
	NamaKategori string `json:"nama_kategori" validate:"required,min=1"`
}

type UpdateKategoriRequest struct {
	NamaKategori string `json:"nama_kategori" validate:"required,min=1"`
}
