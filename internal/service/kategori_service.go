package service

import (
	"context"
	"errors"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/repository"
)

var (
	ErrKategoriNamaDuplikat = errors.New("Nama kategori sudah terdaftar dalam sistem.")
	ErrKategoriHasProduk    = errors.New("Kategori tidak dapat dihapus karena masih memiliki produk terkait.")
)

type KategoriService interface {
	GetAll(ctx context.Context) ([]model.KategoriResponse, error)
	Create(ctx context.Context, req model.CreateKategoriRequest) (*model.KategoriResponse, error)
	Update(ctx context.Context, id string, req model.UpdateKategoriRequest) (*model.KategoriResponse, error)
	Delete(ctx context.Context, id string) error
}

type kategoriService struct {
	repo repository.KategoriRepository
}

func NewKategoriService(repo repository.KategoriRepository) KategoriService {
	return &kategoriService{repo: repo}
}

func (s *kategoriService) GetAll(ctx context.Context) ([]model.KategoriResponse, error) {
	return s.repo.FindAll(ctx)
}

func (s *kategoriService) Create(ctx context.Context, req model.CreateKategoriRequest) (*model.KategoriResponse, error) {
	// Cek duplikasi nama (case-insensitive)
	existing, err := s.repo.FindByNama(ctx, req.NamaKategori)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrKategoriNamaDuplikat
	}

	kategori := &model.Kategori{
		NamaKategori: req.NamaKategori,
	}
	if err := s.repo.Create(ctx, kategori); err != nil {
		return nil, err
	}

	return &model.KategoriResponse{
		ID:           kategori.ID,
		KodeKategori: kategori.KodeKategori,
		NamaKategori: kategori.NamaKategori,
		JumlahProduk: 0,
	}, nil
}

func (s *kategoriService) Update(ctx context.Context, id string, req model.UpdateKategoriRequest) (*model.KategoriResponse, error) {
	// Cek exist
	current, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cek duplikasi nama (exclude self)
	existing, err := s.repo.FindByNama(ctx, req.NamaKategori)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.ID != current.ID {
		return nil, ErrKategoriNamaDuplikat
	}

	return s.repo.Update(ctx, id, req.NamaKategori)
}

func (s *kategoriService) Delete(ctx context.Context, id string) error {
	// Cek exist
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}

	// Cek produk terkait
	count, err := s.repo.CountProdukByKategoriID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrKategoriHasProduk
	}

	return s.repo.Delete(ctx, id)
}
