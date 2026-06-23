package service

import (
	"context"
	"log"
	"time"

	"freya-skin-clinic-backend/internal/repository"
)

type WorkerService struct {
	batchRepo   repository.BatchRepository
	kemasanRepo repository.KemasanTerbukaRepository
}

func NewWorkerService(batchRepo repository.BatchRepository, kemasanRepo repository.KemasanTerbukaRepository) *WorkerService {
	return &WorkerService{
		batchRepo:   batchRepo,
		kemasanRepo: kemasanRepo,
	}
}

// Start menjalankan background worker setiap interval
func (w *WorkerService) Start(ctx context.Context, interval time.Duration) {
	log.Printf("[worker] Background worker started (interval: %s)", interval)

	// Jalankan sekali saat startup
	w.runOnce(ctx)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.runOnce(ctx)
		case <-ctx.Done():
			log.Println("[worker] Background worker stopped")
			return
		}
	}
}

func (w *WorkerService) runOnce(ctx context.Context) {
	updatedBatches := w.updateExpiredBatches(ctx)
	updatedKemasans := w.updateExpiredBUD(ctx)

	if updatedBatches > 0 || updatedKemasans > 0 {
		log.Printf("[worker] Updated: %d batch KADALUWARSA, %d kemasan terbuka BUD KADALUWARSA",
			updatedBatches, updatedKemasans)
	}
}

// updateExpiredBatches: update status batch yang expired_date < NOW() ke KADALUWARSA
func (w *WorkerService) updateExpiredBatches(ctx context.Context) int {
	batches, err := w.batchRepo.FindExpiredBatches(ctx)
	if err != nil {
		log.Printf("[worker] Error fetching expired batches: %v", err)
		return 0
	}

	count := 0
	for _, batch := range batches {
		if err := w.batchRepo.UpdateStatus(ctx, batch.ID, "KADALUWARSA"); err != nil {
			log.Printf("[worker] Error updating batch %s: %v", batch.ID, err)
			continue
		}
		count++
	}
	return count
}

// updateExpiredBUD: update kemasan terbuka yang bud < NOW() ke KADALUWARSA
func (w *WorkerService) updateExpiredBUD(ctx context.Context) int {
	kemasans, err := w.kemasanRepo.FindExpiredBUD(ctx)
	if err != nil {
		log.Printf("[worker] Error fetching expired BUD: %v", err)
		return 0
	}

	count := 0
	for _, kt := range kemasans {
		if err := w.kemasanRepo.UpdateStatus(ctx, kt.ID, "KADALUWARSA"); err != nil {
			log.Printf("[worker] Error updating kemasan terbuka %s: %v", kt.ID, err)
			continue
		}
		count++
	}
	return count
}
