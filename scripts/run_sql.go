// scripts/run_sql.go - Helper untuk menjalankan SQL script
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
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run scripts/run_sql.go <sql_file> [env_file]")
	}

	sqlFile := os.Args[1]
	envFile := ".env.production"
	if len(os.Args) >= 3 {
		envFile = os.Args[2]
	}

	// Load environment
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: gagal load %s: %v", envFile, err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL tidak di-set")
	}

	// Read SQL file
	sqlContent, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Gagal membaca file %s: %v", sqlFile, err)
	}

	// Connect to database
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	defer conn.Close(ctx)

	// Execute SQL
	fmt.Printf("🔄 Menjalankan script: %s\n", sqlFile)
	fmt.Printf("🗄️  Database: %s\n\n", maskDatabaseURL(dbURL))

	_, err = conn.Exec(ctx, string(sqlContent))
	if err != nil {
		log.Fatalf("❌ Gagal eksekusi SQL: %v", err)
	}

	fmt.Println("✅ Script berhasil dijalankan!")
}

func maskDatabaseURL(url string) string {
	// Mask password in connection string for display
	if len(url) > 50 {
		return url[:30] + "***" + url[len(url)-20:]
	}
	return "***"
}
