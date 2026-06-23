package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/pkg/response"
	"freya-skin-clinic-backend/internal/service"
)

type StokMasukHandler struct {
	svc service.StokMasukService
}

func NewStokMasukHandler(svc service.StokMasukService) *StokMasukHandler {
	return &StokMasukHandler{svc: svc}
}

func (h *StokMasukHandler) GetAll(c *fiber.Ctx) error {
	data, err := h.svc.GetAll(c.Context())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil data stok masuk", nil)
	}
	return response.Success(c, http.StatusOK, "Data stok masuk berhasil diambil", data)
}

func (h *StokMasukHandler) Create(c *fiber.Ctx) error {
	var req model.StokMasukRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid", nil)
	}

	if req.IDProduk == "" || req.TanggalPenerimaan == "" || req.ExpiredDate == "" || req.JumlahKemasan <= 0 {
		return response.Error(c, http.StatusBadRequest, "Field wajib tidak lengkap", nil)
	}

	userID, _ := c.Locals("user_id").(string)

	data, err := h.svc.Create(c.Context(), req, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStokMasukExpiredTooEarly):
			return response.Error(c, http.StatusBadRequest, err.Error(), nil)
		case errors.Is(err, service.ErrStokMasukTanggalFuture):
			return response.Error(c, http.StatusBadRequest, err.Error(), nil)
		case errors.Is(err, service.ErrProdukKategoriNotFound):
			return response.Error(c, http.StatusNotFound, "Produk tidak ditemukan", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal menyimpan stok masuk", nil)
		}
	}
	return response.Success(c, http.StatusCreated, "Data stok masuk berhasil disimpan.", data)
}
