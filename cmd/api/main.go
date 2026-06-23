package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/config"
	"freya-skin-clinic-backend/internal/handler"
	"freya-skin-clinic-backend/internal/repository"
	"freya-skin-clinic-backend/internal/router"
	"freya-skin-clinic-backend/internal/service"
)

func main() {
	// Load config
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL tidak dikonfigurasi")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET tidak dikonfigurasi")
	}

	// Init DB connection pool
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Database tidak dapat dijangkau: %v", err)
	}
	log.Println("Database terhubung")

	// Wire dependencies
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authSvc)

	kategoriRepo := repository.NewKategoriRepository(db)
	kategoriSvc := service.NewKategoriService(kategoriRepo)
	kategoriHandler := handler.NewKategoriHandler(kategoriSvc)

	produkRepo := repository.NewProdukRepository(db)
	produkSvc := service.NewProdukService(produkRepo, kategoriRepo)
	produkHandler := handler.NewProdukHandler(produkSvc)

	batchRepo := repository.NewBatchRepository(db)
	batchFEFORepo := repository.NewBatchFEFORepository(db)
	stokMasukRepo := repository.NewStokMasukRepository(db)
	stokMasukSvc := service.NewStokMasukService(stokMasukRepo, batchRepo, produkRepo)
	stokMasukHandler := handler.NewStokMasukHandler(stokMasukSvc)

	kemasanRepo := repository.NewKemasanTerbukaRepository(db)
	stokKeluarRepo := repository.NewStokKeluarRepository(db)
	stokKeluarSvc := service.NewStokKeluarService(stokKeluarRepo, batchRepo, batchFEFORepo, kemasanRepo, produkRepo)
	stokKeluarHandler := handler.NewStokKeluarHandler(stokKeluarSvc)

	monitoringRepo := repository.NewMonitoringRepository(db)
	monitoringSvc := service.NewMonitoringService(monitoringRepo)
	monitoringHandler := handler.NewMonitoringHandler(monitoringSvc)

	opnameRepo := repository.NewOpnameRepository(db)
	opnameSvc := service.NewOpnameService(opnameRepo)
	opnameHandler := handler.NewOpnameHandler(opnameSvc)

	laporanRepo := repository.NewLaporanRepository(db)
	laporanSvc := service.NewLaporanService(laporanRepo)
	laporanHandler := handler.NewLaporanHandler(laporanSvc)

	// Setup router
	app := router.Setup(cfg, userRepo, authHandler, kategoriHandler, produkHandler, stokMasukHandler, stokKeluarHandler, monitoringHandler, opnameHandler, laporanHandler)

	// Start background worker
	workerSvc := service.NewWorkerService(batchRepo, kemasanRepo)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go workerSvc.Start(workerCtx, 1*time.Hour)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server berjalan di http://localhost%s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
