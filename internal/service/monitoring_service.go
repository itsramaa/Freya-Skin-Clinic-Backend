package service

import (
	"context"

	"freya-skin-clinic-backend/internal/model"
	"freya-skin-clinic-backend/internal/repository"
)

type MonitoringService interface {
	GetAll(ctx context.Context, filter model.MonitoringFilter) ([]model.MonitoringProdukItem, error)
}

type monitoringService struct {
	monitoringRepo repository.MonitoringRepository
}

func NewMonitoringService(monitoringRepo repository.MonitoringRepository) MonitoringService {
	return &monitoringService{monitoringRepo: monitoringRepo}
}

func (s *monitoringService) GetAll(ctx context.Context, filter model.MonitoringFilter) ([]model.MonitoringProdukItem, error) {
	return s.monitoringRepo.FindAllForMonitoring(ctx, filter)
}
