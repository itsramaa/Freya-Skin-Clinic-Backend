package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"freya-skin-clinic-backend/internal/config"
	"freya-skin-clinic-backend/internal/handler"
	"freya-skin-clinic-backend/internal/middleware"
)

func Setup(cfg *config.Config, authHandler *handler.AuthHandler, kategoriHandler *handler.KategoriHandler, produkHandler *handler.ProdukHandler, stokMasukHandler *handler.StokMasukHandler, stokKeluarHandler *handler.StokKeluarHandler, monitoringHandler *handler.MonitoringHandler, opnameHandler *handler.OpnameHandler, laporanHandler *handler.LaporanHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://freya-skin-clinic.vercel.app,http://localhost:5173,http://localhost:8080",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	api := app.Group("/api")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// Protected routes
	protected := api.Group("/", middleware.JWTMiddleware(cfg))
	protected.Put("/auth/password", authHandler.ChangePassword)

	// Kategori routes
	kategori := protected.Group("/kategori")
	kategori.Get("/", kategoriHandler.GetAll)
	kategori.Post("/", kategoriHandler.Create)
	kategori.Put("/:id", kategoriHandler.Update)
	kategori.Delete("/:id", kategoriHandler.Delete)

	// Produk routes
	produk := protected.Group("/produk")
	produk.Get("/", produkHandler.GetAll)
	produk.Post("/", produkHandler.Create)
	produk.Put("/:id", produkHandler.Update)
	produk.Delete("/:id", produkHandler.Delete)

	// Stok Masuk routes
	stokMasuk := protected.Group("/stok-masuk")
	stokMasuk.Get("/", stokMasukHandler.GetAll)
	stokMasuk.Post("/", stokMasukHandler.Create)

	// Stok Keluar routes
	stokKeluar := protected.Group("/stok-keluar")
	stokKeluar.Get("/", stokKeluarHandler.GetAll)
	stokKeluar.Get("/preview-batch", stokKeluarHandler.GetPreviewBatch)
	stokKeluar.Post("/", stokKeluarHandler.Create)

	// Monitoring routes
	protected.Get("/monitoring", monitoringHandler.GetAll)

	// Opname routes
	opname := protected.Group("/opname")
	opname.Get("/", opnameHandler.GetAll)
	opname.Post("/", opnameHandler.MulaiOpname)
	opname.Get("/:id", opnameHandler.GetDetail)
	opname.Post("/:id/selesaikan", opnameHandler.SelesaikanOpname)
	opname.Post("/:id/batalkan", opnameHandler.BatalkanOpname)

	// Laporan routes
	laporan := protected.Group("/laporan")
	laporan.Get("/stok-masuk", laporanHandler.GetStokMasuk)
	laporan.Get("/stok-keluar", laporanHandler.GetStokKeluar)
	laporan.Get("/sisa-stok", laporanHandler.GetSisaStok)

	return app
}
