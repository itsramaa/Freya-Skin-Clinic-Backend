package service

import (
	"context"
	"errors"
	"time"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/repository"
)

var (
	ErrStokMasukExpiredTooEarly = errors.New("Tanggal kedaluwarsa tidak boleh kurang dari atau sama dengan tanggal penerimaan.")
	ErrStokMasukTanggalFuture   = errors.New("Tanggal penerimaan tidak boleh melebihi tanggal hari ini.")
	ErrBatchSudahDigunakan      = errors.New("Batch sudah digunakan dalam transaksi stok keluar, tidak dapat diubah atau dihapus.")
	ErrBatchStokNegatif         = errors.New("Perubahan tidak dapat dilakukan karena stok batch akan menjadi negatif.")
)

type StokMasukService interface {
	Create(ctx context.Context, req model.StokMasukRequest, userID string) (*model.StokMasukResponse, error)
	GetAll(ctx context.Context) ([]model.StokMasukResponse, error)
	Update(ctx context.Context, id string, req model.UpdateStokMasukRequest) error
	Delete(ctx context.Context, id string) error
}

type stokMasukService struct {
	stokMasukRepo repository.StokMasukRepository
	batchRepo     repository.BatchRepository
	produkRepo    repository.ProdukRepository
}

func NewStokMasukService(
	stokMasukRepo repository.StokMasukRepository,
	batchRepo repository.BatchRepository,
	produkRepo repository.ProdukRepository,
) StokMasukService {
	return &stokMasukService{
		stokMasukRepo: stokMasukRepo,
		batchRepo:     batchRepo,
		produkRepo:    produkRepo,
	}
}

func (s *stokMasukService) GetAll(ctx context.Context) ([]model.StokMasukResponse, error) {
	return s.stokMasukRepo.FindAll(ctx)
}

func (s *stokMasukService) Create(ctx context.Context, req model.StokMasukRequest, userID string) (*model.StokMasukResponse, error) {
	tglPenerimaan, err := time.Parse("2006-01-02", req.TanggalPenerimaan)
	if err != nil {
		return nil, errors.New("Format tanggal penerimaan tidak valid (YYYY-MM-DD)")
	}
	expiredDate, err := time.Parse("2006-01-02", req.ExpiredDate)
	if err != nil {
		return nil, errors.New("Format expired date tidak valid (YYYY-MM-DD)")
	}

	if tglPenerimaan.After(time.Now().Truncate(24 * time.Hour)) {
		return nil, ErrStokMasukTanggalFuture
	}
	if !expiredDate.After(tglPenerimaan) {
		return nil, ErrStokMasukExpiredTooEarly
	}

	produk, err := s.produkRepo.FindByID(ctx, req.IDProduk)
	if err != nil {
		return nil, ErrProdukKategoriNotFound
	}

	// Hitung total_isi_masuk berdasarkan pola:
	// - Full Use: total = jumlah_kemasan (unit kemasan, tidak dikali isi)
	// - Partial Use: total = jumlah_kemasan * isi_per_kemasan
	var totalIsiMasuk float64
	if produk.PolaPenggunaan == "PARTIAL_USE" && produk.IsiPerKemasan != nil {
		totalIsiMasuk = float64(req.JumlahKemasan) * *produk.IsiPerKemasan
	} else {
		totalIsiMasuk = float64(req.JumlahKemasan)
	}

	existingBatch, err := s.batchRepo.FindByProdukAndExpired(ctx, req.IDProduk, expiredDate)
	if err != nil {
		return nil, err
	}

	var batchID, kodeBatch string

	if existingBatch != nil {
		if err := s.batchRepo.UpdateStok(ctx, existingBatch.ID, req.JumlahKemasan, totalIsiMasuk); err != nil {
			return nil, err
		}
		batchID = existingBatch.ID
		kodeBatch = existingBatch.KodeBatch
	} else {
		newBatch := &model.BatchStok{
			IDProduk:         req.IDProduk,
			ExpiredDate:      expiredDate,
			StokKemasan:      req.JumlahKemasan,
			TotalIsiTersedia: totalIsiMasuk,
			Status:           "AKTIF",
		}
		if err := s.batchRepo.Create(ctx, newBatch); err != nil {
			return nil, err
		}
		batchID = newBatch.ID
		kodeBatch = newBatch.KodeBatch
	}

	sm := &model.StokMasuk{
		IDProduk:          req.IDProduk,
		IDBatch:           batchID,
		IDUser:            userID,
		TanggalPenerimaan: tglPenerimaan,
		JumlahKemasan:     req.JumlahKemasan,
		TotalIsiMasuk:     totalIsiMasuk,
		Keterangan:        req.Keterangan,
	}
	if err := s.stokMasukRepo.Create(ctx, sm); err != nil {
		return nil, err
	}

	return &model.StokMasukResponse{
		ID:                sm.ID,
		IDProduk:          req.IDProduk,
		KodeProduk:        produk.KodeProduk,
		NamaProduk:        produk.NamaProduk,
		NamaKategori:      produk.NamaKategori,
		PolaPenggunaan:    produk.PolaPenggunaan,
		SatuanIsi:         produk.SatuanIsi,
		IsiPerKemasan:     produk.IsiPerKemasan,
		KodeBatch:         kodeBatch,
		TanggalPenerimaan: tglPenerimaan.Format("2006-01-02"),
		ExpiredDate:       expiredDate.Format("2006-01-02"),
		JumlahKemasan:     req.JumlahKemasan,
		TotalIsiMasuk:     totalIsiMasuk,
		Keterangan:        req.Keterangan,
		CreatedAt:         sm.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *stokMasukService) Update(ctx context.Context, id string, req model.UpdateStokMasukRequest) error {
	existing, err := s.stokMasukRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("Data stok masuk tidak ditemukan.")
	}

	if req.JumlahKemasan <= 0 {
		return errors.New("Jumlah kemasan harus lebih dari 0.")
	}

	tglPenerimaan, err := time.Parse("2006-01-02", req.TanggalPenerimaan)
	if err != nil {
		return errors.New("Format tanggal penerimaan tidak valid (YYYY-MM-DD)")
	}
	if tglPenerimaan.After(time.Now().Truncate(24 * time.Hour)) {
		return ErrStokMasukTanggalFuture
	}

	deltaKemasan := req.JumlahKemasan - existing.JumlahKemasan

	// Jika delta negatif (pengurangan), cek apakah stok batch akan menjadi negatif.
	// Ini mencegah perubahan yang merusak integritas stok tanpa memblokir
	// entry stok masuk baru yang berbagi batch dengan transaksi lama.
	if deltaKemasan < 0 {
		batch, err := s.batchRepo.FindByID(ctx, existing.IDBatch)
		if err != nil {
			return err
		}
		if batch.StokKemasan+deltaKemasan < 0 {
			return ErrBatchStokNegatif
		}
	}

	produk, err := s.produkRepo.FindByID(ctx, existing.IDProduk)
	if err != nil {
		return ErrProdukKategoriNotFound
	}

	var deltaIsi float64
	if produk.PolaPenggunaan == "PARTIAL_USE" && produk.IsiPerKemasan != nil {
		deltaIsi = float64(deltaKemasan) * *produk.IsiPerKemasan
	} else {
		deltaIsi = float64(deltaKemasan)
	}

	return s.stokMasukRepo.Update(ctx, id, req, deltaKemasan, deltaIsi)
}

func (s *stokMasukService) Delete(ctx context.Context, id string) error {
	existing, err := s.stokMasukRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("Data stok masuk tidak ditemukan.")
	}

	used, err := s.stokMasukRepo.CheckBatchUsed(ctx, existing.IDBatch)
	if err != nil {
		return err
	}
	if used {
		return ErrBatchSudahDigunakan
	}

	return s.stokMasukRepo.Delete(ctx, id, existing.IDBatch)
}
