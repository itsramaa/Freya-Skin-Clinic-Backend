package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/pkg/response"
	"freya-skin-clinic-backend/internal/repository"
	"freya-skin-clinic-backend/internal/service"
)

type ProdukHandler struct {
	svc service.ProdukService
}

func NewProdukHandler(svc service.ProdukService) *ProdukHandler {
	return &ProdukHandler{svc: svc}
}

func (h *ProdukHandler) GetAll(c *fiber.Ctx) error {
	data, err := h.svc.GetAll(c.Context())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil data produk", nil)
	}
	return response.Success(c, http.StatusOK, "Data produk berhasil diambil", data)
}

func (h *ProdukHandler) Create(c *fiber.Ctx) error {
	var req model.CreateProdukRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid", nil)
	}

	if req.NamaProduk == "" || req.IDKategori == "" || req.BentukKemasan == "" ||
		req.SatuanIsi == "" || req.PolaPenggunaan == "" {
		return response.Error(c, http.StatusBadRequest, "Semua field wajib diisi", nil)
	}

	data, err := h.svc.Create(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProdukIsiPerKemasanDiperlukan):
			return response.Error(c, http.StatusBadRequest, err.Error(), nil)
		case errors.Is(err, service.ErrProdukKategoriNotFound):
			return response.Error(c, http.StatusNotFound, err.Error(), nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal menambahkan produk", nil)
		}
	}
	return response.Success(c, http.StatusCreated, "Produk berhasil ditambahkan.", data)
}

func (h *ProdukHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateProdukRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid", nil)
	}

	data, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProdukNotFound):
			return response.Error(c, http.StatusNotFound, "Produk tidak ditemukan", nil)
		case errors.Is(err, service.ErrProdukPolaPenggunaanLocked):
			return response.Error(c, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, service.ErrProdukIsiPerKemasanDiperlukan):
			return response.Error(c, http.StatusBadRequest, err.Error(), nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal memperbarui produk", nil)
		}
	}
	return response.Success(c, http.StatusOK, "Produk berhasil diperbarui.", data)
}

func (h *ProdukHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.svc.Delete(c.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProdukNotFound):
			return response.Error(c, http.StatusNotFound, "Produk tidak ditemukan", nil)
		case errors.Is(err, service.ErrProdukStokAktif):
			return response.Error(c, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, service.ErrProdukHasTransaksi):
			return response.Error(c, http.StatusConflict, err.Error(), nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal menghapus produk", nil)
		}
	}
	return response.Success(c, http.StatusOK, "Data produk berhasil dihapus.", nil)
}
