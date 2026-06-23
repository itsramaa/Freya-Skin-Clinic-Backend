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

type KategoriHandler struct {
	svc service.KategoriService
}

func NewKategoriHandler(svc service.KategoriService) *KategoriHandler {
	return &KategoriHandler{svc: svc}
}

func (h *KategoriHandler) GetAll(c *fiber.Ctx) error {
	data, err := h.svc.GetAll(c.Context())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil data kategori. Silakan coba lagi.", nil)
	}
	return response.Success(c, http.StatusOK, "Data kategori berhasil diambil", data)
}

func (h *KategoriHandler) Create(c *fiber.Ctx) error {
	var req model.CreateKategoriRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid.", nil)
	}
	if req.NamaKategori == "" {
		return response.Error(c, http.StatusBadRequest, "Nama kategori wajib diisi.", nil)
	}

	data, err := h.svc.Create(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrKategoriNamaDuplikat) {
			return response.Error(c, http.StatusConflict, "Nama kategori sudah terdaftar. Gunakan nama yang berbeda.", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "Gagal menyimpan kategori. Silakan coba lagi.", nil)
	}
	return response.Success(c, http.StatusCreated, "Kategori berhasil ditambahkan.", data)
}

func (h *KategoriHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateKategoriRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid.", nil)
	}
	if req.NamaKategori == "" {
		return response.Error(c, http.StatusBadRequest, "Nama kategori wajib diisi.", nil)
	}

	data, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrKategoriNotFound):
			return response.Error(c, http.StatusNotFound, "Kategori tidak ditemukan.", nil)
		case errors.Is(err, service.ErrKategoriNamaDuplikat):
			return response.Error(c, http.StatusConflict, "Nama kategori sudah terdaftar. Gunakan nama yang berbeda.", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal memperbarui kategori. Silakan coba lagi.", nil)
		}
	}
	return response.Success(c, http.StatusOK, "Kategori berhasil diperbarui.", data)
}

func (h *KategoriHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.svc.Delete(c.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrKategoriNotFound):
			return response.Error(c, http.StatusNotFound, "Kategori tidak ditemukan.", nil)
		case errors.Is(err, service.ErrKategoriHasProduk):
			return response.Error(c, http.StatusConflict, "Kategori tidak dapat dihapus karena masih memiliki produk terkait.", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, "Gagal menghapus kategori. Silakan coba lagi.", nil)
		}
	}
	return response.Success(c, http.StatusOK, "Kategori berhasil dihapus.", nil)
}
