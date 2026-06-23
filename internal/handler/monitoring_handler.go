package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/pkg/response"
	"freya-skin-clinic-backend/internal/service"
)

type MonitoringHandler struct {
	svc service.MonitoringService
}

func NewMonitoringHandler(svc service.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{svc: svc}
}

func (h *MonitoringHandler) GetAll(c *fiber.Ctx) error {
	filter := model.MonitoringFilter{
		KategoriID:  c.Query("kategori_id"),
		StatusBatch: c.Query("status_batch"),
		StatusBUD:   c.Query("status_bud"),
		NamaProduk:  c.Query("nama_produk"),
	}

	data, err := h.svc.GetAll(c.Context(), filter)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Gagal mengambil data monitoring", nil)
	}
	return response.Success(c, http.StatusOK, "Data monitoring berhasil diambil", data)
}
