// scripts/verify_seed.go - Verifikasi hasil seeding
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	envFile := ".env.production"
	if len(os.Args) >= 2 {
		envFile = os.Args[1]
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: gagal load %s: %v", envFile, err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL tidak di-set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	defer conn.Close(ctx)

	fmt.Println("📊 VERIFIKASI DATA SEEDING")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println()

	// Check kategori
	var countKategori int
	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM kategori").Scan(&countKategori)
	if err != nil {
		log.Fatalf("Error query kategori: %v", err)
	}
	fmt.Printf("✅ Kategori: %d record\n", countKategori)

	// List kategori
	rows, err := conn.Query(ctx, "SELECT kode_kategori, nama_kategori FROM kategori ORDER BY kode_kategori")
	if err != nil {
		log.Fatalf("Error list kategori: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var kode, nama string
		rows.Scan(&kode, &nama)
		fmt.Printf("   - %s: %s\n", kode, nama)
	}
	fmt.Println()

	// Check produk
	var countProduk int
	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM produk").Scan(&countProduk)
	if err != nil {
		log.Fatalf("Error query produk: %v", err)
	}
	fmt.Printf("✅ Produk: %d record\n", countProduk)

	// List produk by kategori
	rows2, err := conn.Query(ctx, `
		SELECT k.nama_kategori, COUNT(p.id)
		FROM kategori k
		LEFT JOIN produk p ON p.id_kategori = k.id
		GROUP BY k.nama_kategori
		ORDER BY k.nama_kategori
	`)
	if err != nil {
		log.Fatalf("Error list produk: %v", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var kategori string
		var count int
		rows2.Scan(&kategori, &count)
		fmt.Printf("   - %s: %d produk\n", kategori, count)
	}
	fmt.Println()

	// Check users dengan default password
	var countUsers int
	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE is_default_password = true").Scan(&countUsers)
	if err != nil {
		log.Fatalf("Error query users: %v", err)
	}
	fmt.Printf("✅ Users dengan default password: %d user\n", countUsers)

	// Check data transaksi (harus kosong)
	var countStokMasuk, countStokKeluar, countOpname int
	conn.QueryRow(ctx, "SELECT COUNT(*) FROM stok_masuk").Scan(&countStokMasuk)
	conn.QueryRow(ctx, "SELECT COUNT(*) FROM stok_keluar").Scan(&countStokKeluar)
	conn.QueryRow(ctx, "SELECT COUNT(*) FROM stok_opname").Scan(&countOpname)

	fmt.Println()
	fmt.Println("🧹 DATA TRANSAKSI (harus kosong):")
	fmt.Printf("   - Stok Masuk: %d\n", countStokMasuk)
	fmt.Printf("   - Stok Keluar: %d\n", countStokKeluar)
	fmt.Printf("   - Stok Opname: %d\n", countOpname)

	fmt.Println()
	fmt.Println("✅ VERIFIKASI SELESAI!")
}
