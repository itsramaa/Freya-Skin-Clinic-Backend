package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/pkg/response"
	"freya-skin-clinic-backend/internal/service"
)

type StokKeluarHandler struct {
	svc service.StokKeluarService
}

func NewStokKeluarHandler(svc service.StokKeluarService) *StokKeluarHandler {
	return &StokKeluarHandler{svc: svc}
}

func (h *StokKeluarHandler) GetAll(c *fiber.Ctx) error {
	data, err := h.svc.GetAll(c.Context())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil data stok keluar", nil)
	}
	return response.Success(c, http.StatusOK, "Data stok keluar berhasil diambil", data)
}

func (h *StokKeluarHandler) GetPreviewBatch(c *fiber.Ctx) error {
	idProduk := c.Query("produk_id")
	if idProduk == "" {
		return response.Error(c, http.StatusBadRequest, "produk_id diperlukan", nil)
	}

	data, err := h.svc.GetPreviewBatch(c.Context(), idProduk)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTidakAdaBatch):
			return response.Error(c, http.StatusNotFound, err.Error(), nil)
		case errors.Is(err, service.ErrProdukKategoriNotFound):
			return response.Error(c, http.StatusNotFound, "Produk tidak ditemukan", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal mengambil preview batch", nil)
		}
	}
	return response.Success(c, http.StatusOK, "Preview batch berhasil diambil", data)
}

func (h *StokKeluarHandler) Create(c *fiber.Ctx) error {
	var req model.StokKeluarRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid", nil)
	}

	if req.IDProduk == "" || req.TanggalPenggunaan == "" {
		return response.Error(c, http.StatusBadRequest, "Field wajib tidak lengkap", nil)
	}

	userID, _ := c.Locals("user_id").(string)

	data, err := h.svc.Create(c.Context(), req, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStokKurang):
			return response.Error(c, http.StatusBadRequest, err.Error(), nil)
		case errors.Is(err, service.ErrTidakAdaBatch):
			return response.Error(c, http.StatusBadRequest, err.Error(), nil)
		case errors.Is(err, service.ErrIsiDipakaiMelebihiSisa):
			return response.Error(c, http.StatusBadRequest, err.Error(), nil)
		case errors.Is(err, service.ErrProdukKategoriNotFound):
			return response.Error(c, http.StatusNotFound, "Produk tidak ditemukan", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal menyimpan stok keluar", nil)
		}
	}
	return response.Success(c, http.StatusCreated, "Data penggunaan berhasil disimpan.", data)
}
