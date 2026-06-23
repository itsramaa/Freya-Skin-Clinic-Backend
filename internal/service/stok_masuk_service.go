package service

import (
	"context"
	"errors"
	"time"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/repository"
)

var (
	ErrStokMasukExpiredTooEarly    = errors.New("Tanggal kedaluwarsa tidak boleh kurang dari atau sama dengan tanggal penerimaan.")
	ErrStokMasukTanggalFuture      = errors.New("Tanggal penerimaan tidak boleh melebihi tanggal hari ini.")
)

type StokMasukService interface {
	Create(ctx context.Context, req model.StokMasukRequest, userID string) (*model.StokMasukResponse, error)
	GetAll(ctx context.Context) ([]model.StokMasukResponse, error)
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
	// Parse tanggal
	tglPenerimaan, err := time.Parse("2006-01-02", req.TanggalPenerimaan)
	if err != nil {
		return nil, errors.New("Format tanggal penerimaan tidak valid (YYYY-MM-DD)")
	}
	expiredDate, err := time.Parse("2006-01-02", req.ExpiredDate)
	if err != nil {
		return nil, errors.New("Format expired date tidak valid (YYYY-MM-DD)")
	}

	// Validasi: tanggal penerimaan tidak boleh di masa depan
	if tglPenerimaan.After(time.Now().Truncate(24 * time.Hour)) {
		return nil, ErrStokMasukTanggalFuture
	}

	// Validasi: expired date harus setelah tanggal penerimaan
	if !expiredDate.After(tglPenerimaan) {
		return nil, ErrStokMasukExpiredTooEarly
	}

	// Cek produk exist
	produk, err := s.produkRepo.FindByID(ctx, req.IDProduk)
	if err != nil {
		return nil, ErrProdukKategoriNotFound
	}

	// Hitung total isi masuk
	isiPerKemasan := 1.0
	if produk.IsiPerKemasan != nil {
		isiPerKemasan = *produk.IsiPerKemasan
	}
	totalIsiMasuk := float64(req.JumlahKemasan) * isiPerKemasan

	// Cek batch existing (merge jika ada)
	existingBatch, err := s.batchRepo.FindByProdukAndExpired(ctx, req.IDProduk, expiredDate)
	if err != nil {
		return nil, err
	}

	var batchID string
	var kodeBatch string

	if existingBatch != nil {
		// Merge: tambah stok ke batch yang ada
		if err := s.batchRepo.UpdateStok(ctx, existingBatch.ID, req.JumlahKemasan, totalIsiMasuk); err != nil {
			return nil, err
		}
		batchID = existingBatch.ID
		kodeBatch = existingBatch.KodeBatch
	} else {
		// Buat batch baru
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

	// Catat stok masuk
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
		NamaProduk:        produk.NamaProduk,
		KodeBatch:         kodeBatch,
		TanggalPenerimaan: tglPenerimaan.Format("2006-01-02"),
		ExpiredDate:       expiredDate.Format("2006-01-02"),
		JumlahKemasan:     req.JumlahKemasan,
		TotalIsiMasuk:     totalIsiMasuk,
		Keterangan:        req.Keterangan,
		CreatedAt:         sm.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
