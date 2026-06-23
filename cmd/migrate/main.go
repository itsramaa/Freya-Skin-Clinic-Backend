// cmd/migrate/main.go — Migration tool untuk SiHuni API
// Penggunaan:
//   go run ./cmd/migrate [command] [flags]
//
// Commands:
//   up              — Apply semua pending migrations
//   up-by-one       — Apply satu migration berikutnya
//   up-to VERSION   — Apply sampai versi tertentu
//   down            — Rollback migration terakhir
//   down-to VERSION — Rollback sampai versi tertentu
//   reset           — Rollback semua migrations (DESTRUCTIVE!)
//   status          — Tampilkan status semua migrations
//   version         — Tampilkan versi migration aktif
//   create NAME     — Buat file migration baru (SQL)
//   create-go NAME  — Buat file migration baru (Go)
//   fix             — Normalisasi urutan file migration
//   validate        — Validasi semua migration files

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "./migrations"

func main() {
	// Load .env jika ada
	_ = godotenv.Load()

	// Parse flags
	dir := flag.String("dir", migrationsDir, "Direktori migration files")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	command := args[0]
	commandArgs := args[1:]

	// Validasi DATABASE_URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable tidak di-set")
	}

	// Handle perintah yang tidak butuh koneksi DB
	if command == "create" || command == "create-go" {
		if len(commandArgs) == 0 {
			log.Fatal("Usage: migrate create <name>")
		}
		migType := "sql"
		if command == "create-go" {
			migType = "go"
		}
		if err := goose.Create(nil, *dir, commandArgs[0], migType); err != nil {
			log.Fatalf("create migration gagal: %v", err)
		}
		return
	}

	// Buka koneksi database
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("gagal membuka koneksi database: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}

	// Set goose dialect
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("gagal set dialect: %v", err)
	}

	// Gunakan sequential numbering (000001, 000002, dst)
	goose.SetSequential(true)

	// Jalankan command
	switch command {
	case "up":
		if err := goose.Up(db, *dir); err != nil {
			log.Fatalf("migrate up gagal: %v", err)
		}

	case "up-by-one":
		if err := goose.UpByOne(db, *dir); err != nil {
			log.Fatalf("migrate up-by-one gagal: %v", err)
		}

	case "up-to":
		if len(commandArgs) == 0 {
			log.Fatal("Usage: migrate up-to <version>")
		}
		var version int64
		fmt.Sscan(commandArgs[0], &version)
		if err := goose.UpTo(db, *dir, version); err != nil {
			log.Fatalf("migrate up-to gagal: %v", err)
		}

	case "down":
		if err := goose.Down(db, *dir); err != nil {
			log.Fatalf("migrate down gagal: %v", err)
		}

	case "down-to":
		if len(commandArgs) == 0 {
			log.Fatal("Usage: migrate down-to <version>")
		}
		var version int64
		fmt.Sscan(commandArgs[0], &version)
		if err := goose.DownTo(db, *dir, version); err != nil {
			log.Fatalf("migrate down-to gagal: %v", err)
		}

	case "reset":
		fmt.Println("WARNING: Ini akan rollback SEMUA migrations!")
		fmt.Print("Ketik 'yes' untuk konfirmasi: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Reset dibatalkan.")
			return
		}
		if err := goose.Reset(db, *dir); err != nil {
			log.Fatalf("migrate reset gagal: %v", err)
		}

	case "status":
		if err := goose.Status(db, *dir); err != nil {
			log.Fatalf("migrate status gagal: %v", err)
		}

	case "version":
		if err := goose.Version(db, *dir); err != nil {
			log.Fatalf("migrate version gagal: %v", err)
		}

	case "fix":
		if err := goose.Fix(*dir); err != nil {
			log.Fatalf("migrate fix gagal: %v", err)
		}

	case "validate":
		files, err := os.ReadDir(*dir)
		if err != nil {
			log.Fatalf("gagal membaca direktori migration: %v", err)
		}
		validCount := 0
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			// Check if it's a migration file (SQL or Go)
			if strings.HasSuffix(file.Name(), ".sql") || strings.HasSuffix(file.Name(), ".go") {
				validCount++
			}
		}
		fmt.Printf("✓ Found %d migration files in %s\n", validCount, *dir)

	case "help", "-h", "--help":
		printHelp()

	default:
		fmt.Fprintf(os.Stderr, "Command tidak dikenal: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`SiHuni Migration Tool (powered by goose)

Usage:
  go run ./cmd/migrate [flags] <command> [args]

Flags:
  -dir string    Direktori migration files (default: ./migrations)

Commands:
  up                Apply semua pending migrations
  up-by-one         Apply satu migration berikutnya
  up-to <version>   Apply sampai versi tertentu
  down              Rollback migration terakhir
  down-to <version> Rollback sampai versi tertentu
  reset             Rollback semua migrations (DESTRUCTIVE!)
  status            Tampilkan status semua migrations
  version           Tampilkan versi migration aktif saat ini
  create <name>     Buat SQL migration file baru
  create-go <name>  Buat Go migration file baru
  fix               Normalisasi urutan file migration
  validate          Validasi semua migration files
  help              Tampilkan bantuan ini

Contoh:
  go run ./cmd/migrate up
  go run ./cmd/migrate status
  go run ./cmd/migrate down
  go run ./cmd/migrate create add_users_table
  go run ./cmd/migrate up-to 3

  # Build binary sekali, pakai berkali-kali:
  go build -o ./bin/migrate ./cmd/migrate
  ./bin/migrate status
  ./bin/migrate up

Environment:
  DATABASE_URL  PostgreSQL connection string (wajib)
                Contoh: postgres://user:pass@localhost:5432/sihuni
`)
}
