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
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil data produk. Silakan coba lagi.", nil)
	}
	return response.Success(c, http.StatusOK, "Data produk berhasil diambil", data)
}

func (h *ProdukHandler) Create(c *fiber.Ctx) error {
	var req model.CreateProdukRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid. Pastikan data yang dikirim berupa JSON.", nil)
	}

	if req.NamaProduk == "" {
		return response.Error(c, http.StatusBadRequest, "Nama produk wajib diisi.", nil)
	}
	if req.IDKategori == "" {
		return response.Error(c, http.StatusBadRequest, "Kategori wajib dipilih.", nil)
	}
	if req.BentukKemasan == "" {
		return response.Error(c, http.StatusBadRequest, "Bentuk kemasan wajib diisi.", nil)
	}
	if req.SatuanIsi == "" {
		return response.Error(c, http.StatusBadRequest, "Satuan isi wajib diisi.", nil)
	}
	if req.PolaPenggunaan == "" {
		return response.Error(c, http.StatusBadRequest, "Pola penggunaan wajib dipilih (FULL_USE atau PARTIAL_USE).", nil)
	}

	data, err := h.svc.Create(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProdukIsiPerKemasanDiperlukan):
			return response.Error(c, http.StatusBadRequest, "Isi per kemasan wajib diisi dan harus lebih dari 0 untuk produk Partial Use.", nil)
		case errors.Is(err, service.ErrProdukKategoriNotFound):
			return response.Error(c, http.StatusNotFound, "Kategori yang dipilih tidak ditemukan. Pastikan kategori masih aktif.", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal menyimpan produk. Silakan coba lagi.", nil)
		}
	}
	return response.Success(c, http.StatusCreated, "Produk berhasil ditambahkan.", data)
}

func (h *ProdukHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateProdukRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid.", nil)
	}

	data, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProdukNotFound):
			return response.Error(c, http.StatusNotFound, "Produk tidak ditemukan.", nil)
		case errors.Is(err, service.ErrProdukEditLocked):
			return response.Error(c, http.StatusConflict, "Produk tidak dapat diubah karena sudah memiliki riwayat transaksi masuk atau keluar.", nil)
		case errors.Is(err, service.ErrProdukPolaPenggunaanLocked):
			return response.Error(c, http.StatusConflict, "Pola penggunaan tidak dapat diubah karena produk sudah memiliki riwayat transaksi.", nil)
		case errors.Is(err, service.ErrProdukIsiPerKemasanDiperlukan):
			return response.Error(c, http.StatusBadRequest, "Isi per kemasan wajib diisi dan harus lebih dari 0 untuk produk Partial Use.", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal memperbarui produk. Silakan coba lagi.", nil)
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
			return response.Error(c, http.StatusNotFound, "Produk tidak ditemukan.", nil)
		case errors.Is(err, service.ErrProdukStokAktif):
			return response.Error(c, http.StatusConflict, "Produk tidak dapat dihapus karena masih memiliki stok aktif.", nil)
		case errors.Is(err, service.ErrProdukHasTransaksi):
			return response.Error(c, http.StatusConflict, "Produk tidak dapat dihapus karena memiliki riwayat transaksi.", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal menghapus produk. Silakan coba lagi.", nil)
		}
	}
	return response.Success(c, http.StatusOK, "Produk berhasil dihapus.", nil)
}
