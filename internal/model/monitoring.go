package model

import "time"

// Status indikator batch berdasarkan expired_date
// AMAN: > 30 hari, MENDEKATI: <= 30 hari, KADALUWARSA: expired

type MonitoringBatchItem struct {
	IDBatch          string               `json:"id_batch"`
	KodeBatch        string               `json:"kode_batch"`
	ExpiredDate      string               `json:"expired_date"`
	StokKemasan      int                  `json:"stok_kemasan"`
	TotalIsiTersedia float64              `json:"total_isi_tersedia"`
	StatusBatch      string               `json:"status_batch"`
	IndikatorExpired string               `json:"indikator_expired"` // AMAN, MENDEKATI, KADALUWARSA
	KemasanTerbuka   *MonitoringKemasanTerbuka `json:"kemasan_terbuka"`
}

type MonitoringKemasanTerbuka struct {
	ID         string  `json:"id"`
	BUD        string  `json:"bud"`
	IsiTersisa float64 `json:"isi_tersisa"`
	StatusBUD  string  `json:"status_bud"`
}

type MonitoringProdukItem struct {
	IDProduk       string                `json:"id_produk"`
	KodeProduk     string                `json:"kode_produk"`
	NamaProduk     string                `json:"nama_produk"`
	IDKategori     string                `json:"id_kategori"`
	NamaKategori   string                `json:"nama_kategori"`
	PolaPenggunaan string                `json:"pola_penggunaan"`
	TotalStok      int                   `json:"total_stok"`
	TotalIsi       float64               `json:"total_isi"`
	Batches        []MonitoringBatchItem `json:"batches"`
}

type MonitoringFilter struct {
	KategoriID  string
	StatusBatch string
	StatusBUD   string
	NamaProduk  string
}

// Background worker models
type WorkerJobResult struct {
	UpdatedBatches  int
	UpdatedKemasans int
	Errors          []error
	ExecutedAt      time.Time
}
