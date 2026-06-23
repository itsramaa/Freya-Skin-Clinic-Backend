package model

import "time"

type StokOpname struct {
	ID            string    `db:"id"`
	IDUser        string    `db:"id_user"`
	TanggalOpname time.Time `db:"tanggal_opname"`
	Status        string    `db:"status"`
	Catatan       string    `db:"catatan"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type DetailOpname struct {
	ID               string  `db:"id"`
	IDOpname         string  `db:"id_opname"`
	IDBatch          string  `db:"id_batch"`
	IDKemasanTerbuka *string `db:"id_kemasan_terbuka"`
	StokSistem       float64 `db:"stok_sistem"`
	StokFisik        float64 `db:"stok_fisik"`
	Selisih          float64 `db:"selisih"`
	Keterangan       string  `db:"keterangan"`
}

type DetailOpnameInput struct {
	IDBatch          string  `json:"id_batch"`
	IDKemasanTerbuka *string `json:"id_kemasan_terbuka"`
	StokFisik        float64 `json:"stok_fisik"`
	Keterangan       string  `json:"keterangan"`
}

type SelesaikanOpnameRequest struct {
	Details []DetailOpnameInput `json:"details" validate:"required"`
}

type OpnameItemResponse struct {
	IDBatch            string   `json:"id_batch"`
	KodeBatch          string   `json:"kode_batch"`
	NamaProduk         string   `json:"nama_produk"`
	ExpiredDate        string   `json:"expired_date"`
	PolaPenggunaan     string   `json:"pola_penggunaan"`
	SatuanIsi          string   `json:"satuan_isi"`
	IDKemasanTerbuka   *string  `json:"id_kemasan_terbuka"`
	IsiTersisa         *float64 `json:"isi_tersisa"`
	StokSistem         float64  `json:"stok_sistem"`
	StokFisik          *float64 `json:"stok_fisik"`
	Selisih            *float64 `json:"selisih"`
	Keterangan         string   `json:"keterangan"`
}

type StokOpnameResponse struct {
	ID            string               `json:"id"`
	IDUser        string               `json:"id_user"`
	TanggalOpname string               `json:"tanggal_opname"`
	Status        string               `json:"status"`
	Catatan       string               `json:"catatan"`
	CreatedAt     string               `json:"created_at"`
	Items         []OpnameItemResponse `json:"items,omitempty"`
}
