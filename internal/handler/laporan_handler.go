package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"freya-skin-clinic-backend/internal/pkg/response"
	"freya-skin-clinic-backend/internal/service"
)

type LaporanHandler struct {
	svc service.LaporanService
}

func NewLaporanHandler(svc service.LaporanService) *LaporanHandler {
	return &LaporanHandler{svc: svc}
}

func (h *LaporanHandler) GetStokMasuk(c *fiber.Ctx) error {
	dari := c.Query("dari")
	sampai := c.Query("sampai")
	if dari == "" || sampai == "" {
		return response.Error(c, http.StatusBadRequest, "Parameter 'dari' dan 'sampai' wajib diisi dengan format YYYY-MM-DD.", nil)
	}
	data, err := h.svc.GetStokMasuk(c.Context(), dari, sampai, c.Query("kategori_id"), c.Query("produk_id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, http.StatusOK, "Laporan stok masuk berhasil diambil", data)
}

func (h *LaporanHandler) GetStokKeluar(c *fiber.Ctx) error {
	dari := c.Query("dari")
	sampai := c.Query("sampai")
	if dari == "" || sampai == "" {
		return response.Error(c, http.StatusBadRequest, "Parameter 'dari' dan 'sampai' wajib diisi dengan format YYYY-MM-DD.", nil)
	}
	data, err := h.svc.GetStokKeluar(c.Context(), dari, sampai, c.Query("kategori_id"), c.Query("produk_id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, http.StatusOK, "Laporan stok keluar berhasil diambil", data)
}

func (h *LaporanHandler) GetSisaStok(c *fiber.Ctx) error {
	data, err := h.svc.GetSisaStok(c.Context(), c.Query("kategori_id"), c.Query("produk_id"))
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil laporan sisa stok. Silakan coba lagi.", nil)
	}
	return response.Success(c, http.StatusOK, "Laporan sisa stok berhasil diambil", data)
}
