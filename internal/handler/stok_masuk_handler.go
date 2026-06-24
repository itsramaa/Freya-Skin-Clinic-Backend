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
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil data stok masuk. Silakan coba lagi.", nil)
	}
	return response.Success(c, http.StatusOK, "Data stok masuk berhasil diambil", data)
}

func (h *StokMasukHandler) Create(c *fiber.Ctx) error {
	var req model.StokMasukRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid.", nil)
	}
	if req.IDProduk == "" {
		return response.Error(c, http.StatusBadRequest, "Produk wajib dipilih.", nil)
	}
	if req.TanggalPenerimaan == "" {
		return response.Error(c, http.StatusBadRequest, "Tanggal penerimaan wajib diisi.", nil)
	}
	if req.ExpiredDate == "" {
		return response.Error(c, http.StatusBadRequest, "Tanggal kedaluwarsa wajib diisi.", nil)
	}
	if req.JumlahKemasan <= 0 {
		return response.Error(c, http.StatusBadRequest, "Jumlah kemasan harus lebih dari 0.", nil)
	}

	userID, _ := c.Locals("user_id").(string)

	data, err := h.svc.Create(c.Context(), req, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStokMasukExpiredTooEarly):
			return response.Error(c, http.StatusBadRequest, "Tanggal kedaluwarsa harus setelah tanggal penerimaan.", nil)
		case errors.Is(err, service.ErrStokMasukTanggalFuture):
			return response.Error(c, http.StatusBadRequest, "Tanggal penerimaan tidak boleh melebihi tanggal hari ini.", nil)
		case errors.Is(err, service.ErrProdukKategoriNotFound):
			return response.Error(c, http.StatusNotFound, "Produk tidak ditemukan.", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal menyimpan data stok masuk. Silakan coba lagi.", nil)
		}
	}
	return response.Success(c, http.StatusCreated, "Data stok masuk berhasil disimpan.", data)
}

func (h *StokMasukHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateStokMasukRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid.", nil)
	}
	if req.JumlahKemasan <= 0 {
		return response.Error(c, http.StatusBadRequest, "Jumlah kemasan harus lebih dari 0.", nil)
	}

	err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBatchSudahDigunakan):
			return response.Error(c, http.StatusBadRequest, "Batch sudah digunakan dalam transaksi stok keluar, tidak dapat diubah.", nil)
		case errors.Is(err, service.ErrProdukKategoriNotFound):
			return response.Error(c, http.StatusNotFound, "Produk tidak ditemukan.", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal mengubah data stok masuk. Silakan coba lagi.", nil)
		}
	}
	return response.Success(c, http.StatusOK, "Data stok masuk berhasil diubah.", nil)
}

func (h *StokMasukHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.svc.Delete(c.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBatchSudahDigunakan):
			return response.Error(c, http.StatusBadRequest, "Batch sudah digunakan dalam transaksi stok keluar, tidak dapat dihapus.", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal menghapus data stok masuk. Silakan coba lagi.", nil)
		}
	}
	return response.Success(c, http.StatusOK, "Data stok masuk berhasil dihapus.", nil)
}
