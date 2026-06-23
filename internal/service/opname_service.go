package service

import (
	"context"
	"errors"
	"time"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/repository"
)

var (
	ErrOpnameAktifSudahAda  = errors.New("Sudah ada sesi opname yang aktif. Selesaikan atau batalkan terlebih dahulu.")
	ErrOpnameKeteranganWajib = errors.New("Keterangan wajib diisi untuk item yang memiliki selisih.")
)

type OpnameService interface {
	MulaiOpname(ctx context.Context, userID string) (*model.StokOpnameResponse, error)
	GetAll(ctx context.Context) ([]model.StokOpnameResponse, error)
	GetDetail(ctx context.Context, id string) (*model.StokOpnameResponse, error)
	SelesaikanOpname(ctx context.Context, id string, req model.SelesaikanOpnameRequest) error
	BatalkanOpname(ctx context.Context, id string) error
}

type opnameService struct {
	repo repository.OpnameRepository
}

func NewOpnameService(repo repository.OpnameRepository) OpnameService {
	return &opnameService{repo: repo}
}

func (s *opnameService) MulaiOpname(ctx context.Context, userID string) (*model.StokOpnameResponse, error) {
	op := &model.StokOpname{
		IDUser:        userID,
		TanggalOpname: time.Now().Truncate(24 * time.Hour),
		Status:        "AKTIF",
	}
	if err := s.repo.Create(ctx, op); err != nil {
		return nil, err
	}

	// Ambil items untuk opname
	items, err := s.repo.GetItemsForOpname(ctx)
	if err != nil {
		return nil, err
	}

	return &model.StokOpnameResponse{
		ID:            op.ID,
		IDUser:        op.IDUser,
		TanggalOpname: op.TanggalOpname.Format("2006-01-02"),
		Status:        op.Status,
		CreatedAt:     op.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Items:         items,
	}, nil
}

func (s *opnameService) GetAll(ctx context.Context) ([]model.StokOpnameResponse, error) {
	return s.repo.FindAll(ctx)
}

func (s *opnameService) GetDetail(ctx context.Context, id string) (*model.StokOpnameResponse, error) {
	op, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetItemsForOpname(ctx)
	if err != nil {
		return nil, err
	}

	return &model.StokOpnameResponse{
		ID:            op.ID,
		IDUser:        op.IDUser,
		TanggalOpname: op.TanggalOpname.Format("2006-01-02"),
		Status:        op.Status,
		Catatan:       op.Catatan,
		CreatedAt:     op.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Items:         items,
	}, nil
}

func (s *opnameService) SelesaikanOpname(ctx context.Context, id string, req model.SelesaikanOpnameRequest) error {
	// Cek exist
	op, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if op.Status != "AKTIF" {
		return errors.New("Sesi opname tidak dalam status AKTIF")
	}

	// Validasi keterangan wajib untuk item dengan selisih
	for _, d := range req.Details {
		if d.Keterangan == "" {
			// Akan dicek di repository saat menghitung selisih
			// Validasi di sini berdasarkan input user: jika selisih != 0 maka keterangan wajib
			// Kita tidak tahu stok_sistem di sini, validasi di repo layer
			_ = d
		}
	}

	return s.repo.SaveDetailAndAdjust(ctx, id, req.Details)
}

func (s *opnameService) BatalkanOpname(ctx context.Context, id string) error {
	op, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if op.Status != "AKTIF" {
		return errors.New("Hanya sesi opname AKTIF yang dapat dibatalkan")
	}
	return s.repo.UpdateStatus(ctx, id, "DIBATALKAN")
}
