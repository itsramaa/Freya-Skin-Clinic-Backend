package service

import (
	"context"
	"errors"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/repository"
)

var (
	ErrProdukDataTidakLengkap      = errors.New("Data produk tidak lengkap")
	ErrProdukIsiPerKemasanDiperlukan = errors.New("Isi per kemasan wajib diisi untuk produk Partial Use")
	ErrProdukKategoriNotFound      = errors.New("Kategori tidak ditemukan")
	ErrProdukStokAktif             = errors.New("Produk tidak dapat dihapus karena masih memiliki stok aktif.")
	ErrProdukHasTransaksi          = errors.New("Produk tidak dapat dihapus karena memiliki riwayat transaksi.")
	ErrProdukPolaPenggunaanLocked  = errors.New("Pola penggunaan tidak dapat diubah karena produk sudah memiliki transaksi.")
)

type ProdukService interface {
	GetAll(ctx context.Context) ([]model.ProdukResponse, error)
	Create(ctx context.Context, req model.CreateProdukRequest) (*model.ProdukResponse, error)
	Update(ctx context.Context, id string, req model.UpdateProdukRequest) (*model.ProdukResponse, error)
	Delete(ctx context.Context, id string) error
}

type produkService struct {
	produkRepo   repository.ProdukRepository
	kategoriRepo repository.KategoriRepository
}

func NewProdukService(produkRepo repository.ProdukRepository, kategoriRepo repository.KategoriRepository) ProdukService {
	return &produkService{produkRepo: produkRepo, kategoriRepo: kategoriRepo}
}

func (s *produkService) GetAll(ctx context.Context) ([]model.ProdukResponse, error) {
	return s.produkRepo.FindAll(ctx)
}

func (s *produkService) Create(ctx context.Context, req model.CreateProdukRequest) (*model.ProdukResponse, error) {
	// Validasi isi_per_kemasan untuk PARTIAL_USE
	if req.PolaPenggunaan == "PARTIAL_USE" && (req.IsiPerKemasan == nil || *req.IsiPerKemasan <= 0) {
		return nil, ErrProdukIsiPerKemasanDiperlukan
	}

	// Validasi kategori exists
	if _, err := s.kategoriRepo.FindByID(ctx, req.IDKategori); err != nil {
		return nil, ErrProdukKategoriNotFound
	}

	produk := &model.Produk{
		NamaProduk:     req.NamaProduk,
		IDKategori:     req.IDKategori,
		BentukKemasan:  req.BentukKemasan,
		SatuanIsi:      req.SatuanIsi,
		IsiPerKemasan:  req.IsiPerKemasan,
		PolaPenggunaan: req.PolaPenggunaan,
	}

	if err := s.produkRepo.Create(ctx, produk); err != nil {
		return nil, err
	}

	// Fetch full response
	list, err := s.produkRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range list {
		if p.ID == produk.ID {
			return &p, nil
		}
	}
	return nil, nil
}

func (s *produkService) Update(ctx context.Context, id string, req model.UpdateProdukRequest) (*model.ProdukResponse, error) {
	// Cek exist
	current, err := s.produkRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Lock pola_penggunaan jika ada transaksi
	if req.PolaPenggunaan != current.PolaPenggunaan {
		hasTransaksi, err := s.produkRepo.HasTransaksi(ctx, id)
		if err != nil {
			return nil, err
		}
		if hasTransaksi {
			return nil, ErrProdukPolaPenggunaanLocked
		}
	}

	// Validasi isi_per_kemasan
	if req.PolaPenggunaan == "PARTIAL_USE" && (req.IsiPerKemasan == nil || *req.IsiPerKemasan <= 0) {
		return nil, ErrProdukIsiPerKemasanDiperlukan
	}

	if err := s.produkRepo.Update(ctx, id, req); err != nil {
		return nil, err
	}

	// Fetch full response
	list, err := s.produkRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range list {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, nil
}

func (s *produkService) Delete(ctx context.Context, id string) error {
	if _, err := s.produkRepo.FindByID(ctx, id); err != nil {
		return err
	}

	countStok, err := s.produkRepo.CountStokAktif(ctx, id)
	if err != nil {
		return err
	}
	if countStok > 0 {
		return ErrProdukStokAktif
	}

	countTransaksi, err := s.produkRepo.CountTransaksi(ctx, id)
	if err != nil {
		return err
	}
	if countTransaksi > 0 {
		return ErrProdukHasTransaksi
	}

	return s.produkRepo.Delete(ctx, id)
}
