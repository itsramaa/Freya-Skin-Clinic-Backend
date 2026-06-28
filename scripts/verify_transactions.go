// scripts/verify_transactions.go - Verifikasi hasil seeding transaksi
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

	fmt.Println("📊 VERIFIKASI DATA TRANSAKSI")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()

	// Batch Stok
	rows, _ := conn.Query(ctx, `
		SELECT status, COUNT(*) 
		FROM batch_stok 
		GROUP BY status 
		ORDER BY status
	`)
	fmt.Println("📦 BATCH STOK:")
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		icon := "✅"
		if status == "KADALUWARSA" {
			icon = "⚠️"
		}
		fmt.Printf("   %s %s: %d batch\n", icon, status, count)
	}
	rows.Close()

	// Stok Masuk
	var countStokMasuk int
	conn.QueryRow(ctx, "SELECT COUNT(*) FROM stok_masuk").Scan(&countStokMasuk)
	fmt.Printf("\n📥 STOK MASUK: %d transaksi\n", countStokMasuk)

	// Kemasan Terbuka
	rows2, _ := conn.Query(ctx, `
		SELECT status_bud, COUNT(*) 
		FROM kemasan_terbuka 
		GROUP BY status_bud 
		ORDER BY status_bud
	`)
	fmt.Println("\n📂 KEMASAN TERBUKA:")
	for rows2.Next() {
		var status string
		var count int
		rows2.Scan(&status, &count)
		icon := "✅"
		if status == "KADALUWARSA" {
			icon = "⚠️"
		}
		fmt.Printf("   %s %s: %d kemasan\n", icon, status, count)
	}
	rows2.Close()

	// Stok Keluar
	var countStokKeluar int
	conn.QueryRow(ctx, "SELECT COUNT(*) FROM stok_keluar").Scan(&countStokKeluar)
	fmt.Printf("\n📤 STOK KELUAR: %d transaksi\n", countStokKeluar)

	// Detail per produk
	rows3, _ := conn.Query(ctx, `
		SELECT p.nama_produk, p.pola_penggunaan, COUNT(DISTINCT b.id) as batch_count
		FROM produk p
		LEFT JOIN batch_stok b ON b.id_produk = p.id
		GROUP BY p.nama_produk, p.pola_penggunaan
		HAVING COUNT(DISTINCT b.id) > 0
		ORDER BY p.nama_produk
	`)
	fmt.Println("\n🏷️  PRODUK DENGAN TRANSAKSI:")
	for rows3.Next() {
		var nama, pola string
		var batchCount int
		rows3.Scan(&nama, &pola, &batchCount)
		fmt.Printf("   - %s (%s): %d batch\n", nama, pola, batchCount)
	}
	rows3.Close()

	fmt.Println("\n✅ VERIFIKASI SELESAI!")
}
