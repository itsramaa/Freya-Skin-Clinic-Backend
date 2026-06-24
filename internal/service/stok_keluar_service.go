package service

import (
	"context"
	"errors"
	"time"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/repository"
)

var (
	ErrStokKurang              = errors.New("Stok tidak mencukupi untuk jumlah yang diminta.")
	ErrTidakAdaBatch           = errors.New("Tidak ada stok aktif untuk produk ini.")
	ErrIsiDipakaiMelebihiSisa  = errors.New("Jumlah isi yang dipakai melebihi sisa isi kemasan terbuka.")
	ErrIsiPerKemasanTidakDiset = errors.New("Produk tidak memiliki konfigurasi isi per kemasan.")
)

type StokKeluarService interface {
	Create(ctx context.Context, req model.StokKeluarRequest, userID string) (*model.StokKeluarResponse, error)
	GetAll(ctx context.Context) ([]model.StokKeluarResponse, error)
	GetPreviewBatch(ctx context.Context, idProduk string) (*model.PreviewBatchResponse, error)
}

type stokKeluarService struct {
	stokKeluarRepo repository.StokKeluarRepository
	batchRepo      repository.BatchRepository
	batchFEFORepo  repository.BatchFEFORepository
	kemasanRepo    repository.KemasanTerbukaRepository
	produkRepo     repository.ProdukRepository
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

	var batch *model.BatchStok
	if produk.PolaPenggunaan == "PARTIAL_USE" {
		batch, err = s.batchFEFORepo.FindBatchPartialUseFEFO(ctx, idProduk)
	} else {
		batch, err = s.batchFEFORepo.FindBatchPrioritasFEFO(ctx, idProduk)
	}
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
		SatuanIsi:        produk.SatuanIsi,
		IsiPerKemasan:    produk.IsiPerKemasan,
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
		return nil, errors.New("Format tanggal tidak valid (YYYY-MM-DD)")
	}

	produk, err := s.produkRepo.FindByID(ctx, req.IDProduk)
	if err != nil {
		return nil, ErrProdukKategoriNotFound
	}

	sk := &model.StokKeluar{
		IDProduk:          req.IDProduk,
		IDUser:            userID,
		TanggalPenggunaan: tglPenggunaan,
		Keterangan:        req.Keterangan,
	}

	if produk.PolaPenggunaan == "FULL_USE" {
		if req.JumlahKemasanDipakai <= 0 {
			return nil, errors.New("Jumlah kemasan dipakai harus lebih dari 0.")
		}

		batch, err := s.batchFEFORepo.FindBatchPrioritasFEFO(ctx, req.IDProduk)
		if err != nil {
			return nil, err
		}
		if batch == nil {
			return nil, ErrTidakAdaBatch
		}
		if batch.StokKemasan < req.JumlahKemasanDipakai {
			return nil, ErrStokKurang
		}

		// Full Use: total isi = jumlah kemasan (tidak dikali isi_per_kemasan)
		totalIsi := float64(req.JumlahKemasanDipakai)

		if err := s.batchFEFORepo.ReduceStok(ctx, batch.ID, req.JumlahKemasanDipakai, totalIsi); err != nil {
			return nil, err
		}

		sk.IDBatch = batch.ID
		sk.JumlahKemasanDipakai = req.JumlahKemasanDipakai
		sk.JumlahIsiDipakai = totalIsi

		if err := s.stokKeluarRepo.Create(ctx, sk); err != nil {
			return nil, err
		}

		return &model.StokKeluarResponse{
			ID:                   sk.ID,
			IDProduk:             req.IDProduk,
			NamaProduk:           produk.NamaProduk,
			KodeBatch:            batch.KodeBatch,
			PolaPenggunaan:       produk.PolaPenggunaan,
			SatuanIsi:            produk.SatuanIsi,
			TanggalPenggunaan:    tglPenggunaan.Format("2006-01-02"),
			JumlahKemasanDipakai: sk.JumlahKemasanDipakai,
			JumlahIsiDipakai:     sk.JumlahIsiDipakai,
			Keterangan:           req.Keterangan,
			CreatedAt:            sk.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}, nil
	}

	// ── PARTIAL USE ──
	if req.JumlahIsiDipakai <= 0 {
		return nil, errors.New("Jumlah isi dipakai harus lebih dari 0.")
	}

	// FEFO Partial: prioritaskan batch dengan kemasan terbuka aktif dulu
	batch, err := s.batchFEFORepo.FindBatchPartialUseFEFO(ctx, req.IDProduk)
	if err != nil {
		return nil, err
	}
	if batch == nil {
		return nil, ErrTidakAdaBatch
	}

	kt, err := s.kemasanRepo.FindAktifByBatch(ctx, batch.ID)
	if err != nil {
		return nil, err
	}

	// BUD expired → nonaktifkan, cari batch baru
	if kt != nil && kt.BUD.Before(tglPenggunaan) {
		if err := s.kemasanRepo.UpdateStatus(ctx, kt.ID, "KADALUWARSA"); err != nil {
			return nil, err
		}
		kt = nil
		// Cari ulang batch yang punya stok kemasan
		batch, err = s.batchFEFORepo.FindBatchPrioritasFEFO(ctx, req.IDProduk)
		if err != nil {
			return nil, err
		}
		if batch == nil {
			return nil, ErrTidakAdaBatch
		}
	}

	var idKemasan *string

	if kt != nil && kt.IsiTersisa > 0 {
		// Pakai kemasan terbuka yang masih ada isinya
		if req.JumlahIsiDipakai > kt.IsiTersisa {
			return nil, ErrIsiDipakaiMelebihiSisa
		}
		newSisa := kt.IsiTersisa - req.JumlahIsiDipakai

		// Update isi tersisa kemasan terbuka
		if err := s.kemasanRepo.UpdateIsiTersisa(ctx, kt.ID, newSisa); err != nil {
			return nil, err
		}
		// Update total_isi_tersedia batch
		if err := s.batchFEFORepo.ReduceStok(ctx, batch.ID, 0, req.JumlahIsiDipakai); err != nil {
			return nil, err
		}
		idKemasan = &kt.ID

		// Jika kemasan terbuka habis, tandai KADALUWARSA (habis terpakai)
		if newSisa == 0 {
			_ = s.kemasanRepo.UpdateStatus(ctx, kt.ID, "KADALUWARSA")
		}
	} else {
		// Tidak ada kemasan terbuka atau habis → buka kemasan baru
		if batch.StokKemasan < 1 {
			return nil, ErrStokKurang
		}
		if produk.IsiPerKemasan == nil {
			return nil, ErrIsiPerKemasanTidakDiset
		}
		isiPerKemasan := *produk.IsiPerKemasan
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

		if err := s.batchFEFORepo.ReduceStok(ctx, batch.ID, 1, req.JumlahIsiDipakai); err != nil {
			return nil, err
		}

		// Jika kemasan langsung habis dalam satu pakai
		if newKT.IsiTersisa == 0 {
			_ = s.kemasanRepo.UpdateStatus(ctx, newKT.ID, "KADALUWARSA")
		}
	}

	sk.IDBatch = batch.ID
	sk.IDKemasanTerbuka = idKemasan
	sk.JumlahIsiDipakai = req.JumlahIsiDipakai

	if err := s.stokKeluarRepo.Create(ctx, sk); err != nil {
		return nil, err
	}

	return &model.StokKeluarResponse{
		ID:                sk.ID,
		IDProduk:          req.IDProduk,
		NamaProduk:        produk.NamaProduk,
		KodeBatch:         batch.KodeBatch,
		PolaPenggunaan:    produk.PolaPenggunaan,
		SatuanIsi:         produk.SatuanIsi,
		TanggalPenggunaan: tglPenggunaan.Format("2006-01-02"),
		JumlahIsiDipakai:  sk.JumlahIsiDipakai,
		Keterangan:        req.Keterangan,
		CreatedAt:         sk.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
