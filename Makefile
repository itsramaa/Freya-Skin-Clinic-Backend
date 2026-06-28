# Freya Skin Clinic Backend — Makefile
# Usage: make <target>

BINARY_NAME=freya-api
BUILD_DIR=./bin
CMD_PATH=./cmd/api
VALIDATOR=./cmd/api-validate
MIGRATE_PATH=./cmd/migrate
MIGRATE_BINARY=$(BUILD_DIR)/migrate
MIGRATE=go run $(MIGRATE_PATH)

.PHONY: all dev build test lint clean \
        migrate-up migrate-up-one migrate-up-to migrate-down migrate-down-to \
        migrate-reset migrate-status migrate-version migrate-create migrate-create-go \
        migrate-fix migrate-validate migrate-seed migrate-build \
        migrate-up-prod migrate-up-one-prod migrate-up-to-prod \
        migrate-down-prod migrate-down-to-prod migrate-reset-prod \
        migrate-status-prod migrate-version-prod \
        migrate-seed-prod \
        docs

all: build

## dev: Run the server with live reload (requires air: go install github.com/air-verse/air@latest)
dev:
	@which air > /dev/null 2>&1 || (echo "air not found. Install: go install github.com/air-verse/air@latest" && exit 1)
	air

## build: Compile the binary to ./bin/sihuni-api
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

## run: Build and run the binary directly
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

## test: Run all unit tests with race detector
test:
	go test -race -v ./...

## test-short: Run tests without verbose output
test-short:
	go test ./...

## lint: Run golangci-lint
lint:
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not found. Install: go install github.com/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

## tidy: Tidy go modules
tidy:
	go mod tidy

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)

## docker-up: Start local PostgreSQL via Docker Compose
docker-up:
	docker compose up -d db

## docker-down: Stop local PostgreSQL
docker-down:
	docker compose down

# ── Migrations (embedded goose via cmd/migrate) ────────────────────────────────

## migrate-build: Build migration binary ke ./bin/migrate
migrate-build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(MIGRATE_BINARY) $(MIGRATE_PATH)
	@echo "Built $(MIGRATE_BINARY)"

## migrate-up: Apply semua pending migrations
migrate-up:
	$(MIGRATE) up

## migrate-up-one: Apply satu migration berikutnya
migrate-up-one:
	$(MIGRATE) up-by-one

## migrate-up-to VERSION=<version>: Apply sampai versi tertentu
migrate-up-to:
	@ [ "$(VERSION)" ] || (echo "Usage: make migrate-up-to VERSION=<number>" && exit 1)
	$(MIGRATE) up-to $(VERSION)

## migrate-down: Rollback migration terakhir
migrate-down:
	$(MIGRATE) down

## migrate-down-to VERSION=<version>: Rollback sampai versi tertentu
migrate-down-to:
	@ [ "$(VERSION)" ] || (echo "Usage: make migrate-down-to VERSION=<number>" && exit 1)
	$(MIGRATE) down-to $(VERSION)

## migrate-reset: Rollback SEMUA migrations (DESTRUCTIVE!)
migrate-reset:
	$(MIGRATE) reset

## migrate-status: Tampilkan status semua migrations
migrate-status:
	$(MIGRATE) status

## migrate-version: Tampilkan versi migration aktif
migrate-version:
	$(MIGRATE) version

## migrate-create NAME=<name>: Buat SQL migration file baru
migrate-create:
	@ [ "$(NAME)" ] || (echo "Usage: make migrate-create NAME=<migration_name>" && exit 1)
	$(MIGRATE) create $(NAME)

## migrate-create-go NAME=<name>: Buat Go migration file baru
migrate-create-go:
	@ [ "$(NAME)" ] || (echo "Usage: make migrate-create-go NAME=<migration_name>" && exit 1)
	$(MIGRATE) create-go $(NAME)

## migrate-fix: Normalisasi urutan file migration
migrate-fix:
	$(MIGRATE) fix

## migrate-validate: Validasi semua migration files
migrate-validate:
	$(MIGRATE) validate

## migrate-seed: Apply migrations (seed user sudah ada di migration 000002)
migrate-seed: migrate-up
	@echo "Migrations applied. Seed user sudah termasuk di migration 000002."

# ── Production Migrations (pakai .env.production) ──────────────────────────────

## migrate-up-prod: Apply semua pending migrations pakai .env.production
migrate-up-prod:
	$(MIGRATE) -env .env.production up

## migrate-up-one-prod: Apply satu migration berikutnya pakai .env.production
migrate-up-one-prod:
	$(MIGRATE) -env .env.production up-by-one

## migrate-up-to-prod VERSION=<version>: Apply sampai versi tertentu pakai .env.production
migrate-up-to-prod:
	@ [ "$(VERSION)" ] || (echo "Usage: make migrate-up-to-prod VERSION=<number>" && exit 1)
	$(MIGRATE) -env .env.production up-to $(VERSION)

## migrate-down-prod: Rollback migration terakhir pakai .env.production
migrate-down-prod:
	$(MIGRATE) -env .env.production down

## migrate-down-to-prod VERSION=<version>: Rollback sampai versi tertentu pakai .env.production
migrate-down-to-prod:
	@ [ "$(VERSION)" ] || (echo "Usage: make migrate-down-to-prod VERSION=<number>" && exit 1)
	$(MIGRATE) -env .env.production down-to $(VERSION)

## migrate-reset-prod: Rollback SEMUA migrations pakai .env.production (DESTRUCTIVE!)
migrate-reset-prod:
	$(MIGRATE) -env .env.production reset

## migrate-status-prod: Status migrations pakai .env.production
migrate-status-prod:
	$(MIGRATE) -env .env.production status

## migrate-version-prod: Versi migration aktif pakai .env.production
migrate-version-prod:
	$(MIGRATE) -env .env.production version

## migrate-seed-prod: Apply migrations pakai .env.production (seed user sudah ada di migration 000002)
migrate-seed-prod: migrate-up-prod
	@echo "Migrations applied. Seed user sudah termasuk di migration 000002."

# ── Database Reset & Seed (Production) ────────────────────────────────────────

## reset-seed: Cleansing data dan seed dummy data (LOCAL - gunakan .env)
reset-seed:
	@echo "⚠️  WARNING: Ini akan menghapus SEMUA data transaksi dan reset auth!"
	@echo "Database: LOCAL (.env)"
	@powershell -Command "$$confirm = Read-Host 'Ketik YES untuk melanjutkan'; if($$confirm -ne 'YES'){Write-Host 'Dibatalkan.'; exit 1}"
	@psql $(shell grep DATABASE_URL .env | cut -d '=' -f2-) -f scripts/reset_and_seed.sql

## reset-seed-prod: Cleansing data dan seed dummy data (PRODUCTION - gunakan .env.production)
reset-seed-prod:
	@echo "⚠️  WARNING: Ini akan menghapus SEMUA data transaksi dan reset auth!"
	@echo "Database: PRODUCTION (.env.production)"
	@powershell -Command "$$confirm = Read-Host 'Ketik YES untuk melanjutkan'; if($$confirm -ne 'YES'){Write-Host 'Dibatalkan.'; exit 1}"
	@psql $(shell grep DATABASE_URL .env.production | cut -d '=' -f2-) -f scripts/reset_and_seed.sql

# ── API Documentation ──────────────────────────────────────────────────────────

## docs: Open Scalar API docs in browser (server must be running)
docs:
	@echo "Scalar API Docs: http://localhost:8080/docs"
	@echo "OpenAPI JSON:    http://localhost:8080/openapi.json"

## docs-validate: Validate openapi.json against Go source
docs-validate:
	go run $(VALIDATOR)

## docs-watch: Watch openapi.json changes and revalidate
docs-watch:
	@echo "Watching openapi.json for changes... (Ctrl+C to stop)"
	@powershell -Command "$$f='api/openapi.json'; $$last=(Get-Item $$f).LastWriteTime; Write-Host 'Watching api/openapi.json...'; while(1){Start-Sleep 1;if((Get-Item $$f).LastWriteTime -ne $$last){$$last=(Get-Item $$f).LastWriteTime; go run ./cmd/api-validate}}"

## help: Show this help
help:
	@grep -E '^## ' Makefile | sed 's/## //'
