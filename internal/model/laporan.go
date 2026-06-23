package model

import "time"

type LaporanStokMasukItem struct {
	ID                string  `json:"id"`
	TanggalPenerimaan string  `json:"tanggal_penerimaan"`
	NamaProduk        string  `json:"nama_produk"`
	KodeProduk        string  `json:"kode_produk"`
	NamaKategori      string  `json:"nama_kategori"`
	PolaPenggunaan    string  `json:"pola_penggunaan"`
	SatuanIsi         string  `json:"satuan_isi"`
	KodeBatch         string  `json:"kode_batch"`
	ExpiredDate       string  `json:"expired_date"`
	JumlahKemasan     int     `json:"jumlah_kemasan"`
	TotalIsiMasuk     float64 `json:"total_isi_masuk"`
	Keterangan        string  `json:"keterangan"`
}

type LaporanStokKeluarItem struct {
	ID                   string  `json:"id"`
	TanggalPenggunaan    string  `json:"tanggal_penggunaan"`
	NamaProduk           string  `json:"nama_produk"`
	KodeProduk           string  `json:"kode_produk"`
	NamaKategori         string  `json:"nama_kategori"`
	KodeBatch            string  `json:"kode_batch"`
	PolaPenggunaan       string  `json:"pola_penggunaan"`
	JumlahKemasanDipakai int     `json:"jumlah_kemasan_dipakai"`
	JumlahIsiDipakai     float64 `json:"jumlah_isi_dipakai"`
	Keterangan           string  `json:"keterangan"`
}

type LaporanSisaStokItem struct {
	KodeProduk     string  `json:"kode_produk"`
	NamaProduk     string  `json:"nama_produk"`
	NamaKategori   string  `json:"nama_kategori"`
	PolaPenggunaan string  `json:"pola_penggunaan"`
	TotalStok      int     `json:"total_stok"`
	TotalIsi       float64 `json:"total_isi"`
}

type LaporanFilter struct {
	Dari       time.Time
	Sampai     time.Time
	KategoriID string
	ProdukID   string
}
