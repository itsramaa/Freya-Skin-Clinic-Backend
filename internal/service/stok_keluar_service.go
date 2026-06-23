package service

import (
	"context"
	"errors"
	"time"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/repository"
)

var (
	ErrStokKurang             = errors.New("Stok tidak mencukupi untuk jumlah yang diminta.")
	ErrTidakAdaBatch          = errors.New("Tidak ada stok aktif untuk produk ini.")
	ErrIsiDipakaiMelebihiSisa = errors.New("Jumlah isi yang dipakai melebihi sisa isi kemasan terbuka.")
)

type StokKeluarService interface {
	Create(ctx context.Context, req model.StokKeluarRequest, userID string) (*model.StokKeluarResponse, error)
	GetAll(ctx context.Context) ([]model.StokKeluarResponse, error)
	GetPreviewBatch(ctx context.Context, idProduk string) (*model.PreviewBatchResponse, error)
}

type stokKeluarService struct {
	stokKeluarRepo   repository.StokKeluarRepository
	batchRepo        repository.BatchRepository
	batchFEFORepo    repository.BatchFEFORepository
	kemasanRepo      repository.KemasanTerbukaRepository
	produkRepo       repository.ProdukRepository
}

func NewStokKeluarService(
	stokKeluarRepo repository.StokKeluarRepository,
	batchRepo repository.BatchRepository,
	batchFEFORepo repository.BatchFEFORepository,
	kemasanRepo repository.KemasanTerbukaRepository,
	produkRepo repository.ProdukRepository,
) StokKeluarService {
	return &stokKeluarService{
		stokKeluarRepo: stokKeluarRepo,
		batchRepo:      batchRepo,
		batchFEFORepo:  batchFEFORepo,
		kemasanRepo:    kemasanRepo,
		produkRepo:     produkRepo,
	}
}

func (s *stokKeluarService) GetAll(ctx context.Context) ([]model.StokKeluarResponse, error) {
	return s.stokKeluarRepo.FindAll(ctx)
}

func (s *stokKeluarService) GetPreviewBatch(ctx context.Context, idProduk string) (*model.PreviewBatchResponse, error) {
	produk, err := s.produkRepo.FindByID(ctx, idProduk)
	if err != nil {
		return nil, ErrProdukKategoriNotFound
	}

	batch, err := s.batchFEFORepo.FindBatchPrioritasFEFO(ctx, idProduk)
	if err != nil {
		return nil, err
	}
	if batch == nil {
		return nil, ErrTidakAdaBatch
	}

	preview := &model.PreviewBatchResponse{
		IDBatch:          batch.ID,
		KodeBatch:        batch.KodeBatch,
		ExpiredDate:      batch.ExpiredDate.Format("2006-01-02"),
		StokKemasan:      batch.StokKemasan,
		TotalIsiTersedia: batch.TotalIsiTersedia,
		PolaPenggunaan:   produk.PolaPenggunaan,
	}

	if produk.PolaPenggunaan == "PARTIAL_USE" {
		kt, err := s.kemasanRepo.FindAktifByBatch(ctx, batch.ID)
		if err != nil {
			return nil, err
		}
		if kt != nil {
			preview.KemasanTerbuka = &model.KemasanTerbukaInfo{
				ID:         kt.ID,
				BUD:        kt.BUD.Format("2006-01-02"),
				IsiTersisa: kt.IsiTersisa,
				StatusBUD:  kt.StatusBUD,
			}
		}
	}

	return preview, nil
}

func (s *stokKeluarService) Create(ctx context.Context, req model.StokKeluarRequest, userID string) (*model.StokKeluarResponse, error) {
	tglPenggunaan, err := time.Parse("2006-01-02", req.TanggalPenggunaan)
	if err != nil {
		return nil, errors.New("Format tanggal tidak valid")
	}

	produk, err := s.produkRepo.FindByID(ctx, req.IDProduk)
	if err != nil {
		return nil, ErrProdukKategoriNotFound
	}

	// FEFO: ambil batch prioritas
	batch, err := s.batchFEFORepo.FindBatchPrioritasFEFO(ctx, req.IDProduk)
	if err != nil {
		return nil, err
	}
	if batch == nil {
		return nil, ErrTidakAdaBatch
	}

	sk := &model.StokKeluar{
		IDProduk:          req.IDProduk,
		IDBatch:           batch.ID,
		IDUser:            userID,
		TanggalPenggunaan: tglPenggunaan,
		Keterangan:        req.Keterangan,
	}

	if produk.PolaPenggunaan == "FULL_USE" {
		// Full Use: kurangi stok kemasan
		if req.JumlahKemasanDipakai <= 0 {
			return nil, errors.New("Jumlah kemasan dipakai harus > 0")
		}
		if batch.StokKemasan < req.JumlahKemasanDipakai {
			return nil, ErrStokKurang
		}

		isiPerKemasan := 1.0
		if produk.IsiPerKemasan != nil {
			isiPerKemasan = *produk.IsiPerKemasan
		}
		totalIsi := float64(req.JumlahKemasanDipakai) * isiPerKemasan

		if err := s.batchFEFORepo.ReduceStok(ctx, batch.ID, req.JumlahKemasanDipakai, totalIsi); err != nil {
			return nil, err
		}

		sk.JumlahKemasanDipakai = req.JumlahKemasanDipakai
		sk.JumlahIsiDipakai = totalIsi

	} else {
		// Partial Use
		if req.JumlahIsiDipakai <= 0 {
			return nil, errors.New("Jumlah isi dipakai harus > 0")
		}

		kt, err := s.kemasanRepo.FindAktifByBatch(ctx, batch.ID)
		if err != nil {
			return nil, err
		}

		if kt != nil && kt.BUD.Before(tglPenggunaan) {
			// BUD expired — nonaktifkan dan buka kemasan baru
			if err := s.kemasanRepo.UpdateStatus(ctx, kt.ID, "KADALUWARSA"); err != nil {
				return nil, err
			}
			kt = nil
		}

		var idKemasan *string

		if kt == nil {
			// Buka kemasan baru
			if batch.StokKemasan < 1 {
				return nil, ErrStokKurang
			}
			isiPerKemasan := 1.0
			if produk.IsiPerKemasan != nil {
				isiPerKemasan = *produk.IsiPerKemasan
			}
			if req.JumlahIsiDipakai > isiPerKemasan {
				return nil, ErrIsiDipakaiMelebihiSisa
			}

			newKT := &model.KemasanTerbuka{
				IDBatch:       batch.ID,
				TanggalDibuka: tglPenggunaan,
				BUD:           tglPenggunaan.AddDate(0, 0, 28),
				IsiAwal:       isiPerKemasan,
				IsiTersisa:    isiPerKemasan - req.JumlahIsiDipakai,
				StatusBUD:     "AKTIF",
			}
			if err := s.kemasanRepo.Create(ctx, newKT); err != nil {
				return nil, err
			}
			idKemasan = &newKT.ID

			// Kurangi 1 kemasan dari batch
			if err := s.batchFEFORepo.ReduceStok(ctx, batch.ID, 1, req.JumlahIsiDipakai); err != nil {
				return nil, err
			}
		} else {
			// Pakai kemasan terbuka yang ada
			if req.JumlahIsiDipakai > kt.IsiTersisa {
				return nil, ErrIsiDipakaiMelebihiSisa
			}
			newSisa := kt.IsiTersisa - req.JumlahIsiDipakai
			if err := s.kemasanRepo.UpdateIsiTersisa(ctx, kt.ID, newSisa); err != nil {
				return nil, err
			}
			idKemasan = &kt.ID

			if err := s.batchFEFORepo.ReduceStok(ctx, batch.ID, 0, req.JumlahIsiDipakai); err != nil {
				return nil, err
			}
		}

		sk.IDKemasanTerbuka = idKemasan
		sk.JumlahIsiDipakai = req.JumlahIsiDipakai
	}

	if err := s.stokKeluarRepo.Create(ctx, sk); err != nil {
		return nil, err
	}

	return &model.StokKeluarResponse{
		ID:                   sk.ID,
		IDProduk:             req.IDProduk,
		NamaProduk:           produk.NamaProduk,
		KodeBatch:            batch.KodeBatch,
		PolaPenggunaan:       produk.PolaPenggunaan,
		TanggalPenggunaan:    tglPenggunaan.Format("2006-01-02"),
		JumlahKemasanDipakai: sk.JumlahKemasanDipakai,
		JumlahIsiDipakai:     sk.JumlahIsiDipakai,
		Keterangan:           req.Keterangan,
		CreatedAt:            sk.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
