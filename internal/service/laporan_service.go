package service

import (
	"context"
	"errors"
	"time"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/repository"
)

type LaporanService interface {
	GetStokMasuk(ctx context.Context, dari, sampai, kategoriID, produkID string) ([]model.LaporanStokMasukItem, error)
	GetStokKeluar(ctx context.Context, dari, sampai, kategoriID, produkID string) ([]model.LaporanStokKeluarItem, error)
	GetSisaStok(ctx context.Context, kategoriID, produkID string) ([]model.LaporanSisaStokItem, error)
}

type laporanService struct {
	repo repository.LaporanRepository
}

func NewLaporanService(repo repository.LaporanRepository) LaporanService {
	return &laporanService{repo: repo}
}

func parseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("tanggal tidak boleh kosong")
	}
	return time.Parse("2006-01-02", s)
}

func (s *laporanService) GetStokMasuk(ctx context.Context, dari, sampai, kategoriID, produkID string) ([]model.LaporanStokMasukItem, error) {
	dariT, err := parseDate(dari)
	if err != nil {
		return nil, errors.New("Format tanggal 'dari' tidak valid")
	}
	sampaiT, err := parseDate(sampai)
	if err != nil {
		return nil, errors.New("Format tanggal 'sampai' tidak valid")
	}
	return s.repo.GetStokMasuk(ctx, model.LaporanFilter{
		Dari: dariT, Sampai: sampaiT, KategoriID: kategoriID, ProdukID: produkID,
	})
}

func (s *laporanService) GetStokKeluar(ctx context.Context, dari, sampai, kategoriID, produkID string) ([]model.LaporanStokKeluarItem, error) {
	dariT, err := parseDate(dari)
	if err != nil {
		return nil, errors.New("Format tanggal 'dari' tidak valid")
	}
	sampaiT, err := parseDate(sampai)
	if err != nil {
		return nil, errors.New("Format tanggal 'sampai' tidak valid")
	}
	return s.repo.GetStokKeluar(ctx, model.LaporanFilter{
		Dari: dariT, Sampai: sampaiT, KategoriID: kategoriID, ProdukID: produkID,
	})
}

func (s *laporanService) GetSisaStok(ctx context.Context, kategoriID, produkID string) ([]model.LaporanSisaStokItem, error) {
	return s.repo.GetSisaStok(ctx, model.LaporanFilter{KategoriID: kategoriID, ProdukID: produkID})
}
