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

type OpnameHandler struct {
	svc service.OpnameService
}

func NewOpnameHandler(svc service.OpnameService) *OpnameHandler {
	return &OpnameHandler{svc: svc}
}

func (h *OpnameHandler) GetAll(c *fiber.Ctx) error {
	data, err := h.svc.GetAll(c.Context())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil daftar sesi opname. Silakan coba lagi.", nil)
	}
	return response.Success(c, http.StatusOK, "Data opname berhasil diambil", data)
}

func (h *OpnameHandler) MulaiOpname(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	data, err := h.svc.MulaiOpname(c.Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Gagal memulai sesi opname. Silakan coba lagi.", nil)
	}
	return response.Success(c, http.StatusCreated, "Sesi opname berhasil dimulai.", data)
}

func (h *OpnameHandler) GetDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	data, err := h.svc.GetDetail(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOpnameNotFound) {
			return response.Error(c, http.StatusNotFound, "Sesi opname tidak ditemukan.", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil detail opname. Silakan coba lagi.", nil)
	}
	return response.Success(c, http.StatusOK, "Detail opname berhasil diambil", data)
}

func (h *OpnameHandler) SelesaikanOpname(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.SelesaikanOpnameRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Format request tidak valid.", nil)
	}
	if len(req.Details) == 0 {
		return response.Error(c, http.StatusBadRequest, "Detail opname wajib diisi. Harap isi stok fisik untuk setiap item.", nil)
	}

	err := h.svc.SelesaikanOpname(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrOpnameNotFound):
			return response.Error(c, http.StatusNotFound, "Sesi opname tidak ditemukan.", nil)
		case errors.Is(err, service.ErrOpnameKeteranganWajib),
			errors.Is(err, repository.ErrKeteranganWajib):
			return response.Error(c, http.StatusBadRequest, "Keterangan wajib diisi untuk setiap item yang memiliki selisih stok.", nil)
		default:
			return response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		}
	}
	return response.Success(c, http.StatusOK, "Stok opname berhasil disimpan.", nil)
}

func (h *OpnameHandler) BatalkanOpname(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.svc.BatalkanOpname(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrOpnameNotFound) {
			return response.Error(c, http.StatusNotFound, "Sesi opname tidak ditemukan.", nil)
		}
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}
	return response.Success(c, http.StatusOK, "Sesi opname berhasil dibatalkan.", nil)
}
