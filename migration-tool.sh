#!/bin/bash
# migration-tool.sh — Script wrapper untuk goose migration
# Usage: ./migration-tool.sh [command] [args]

set -e

# Load environment dari .env jika ada
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Pastikan DATABASE_URL ada
if [ -z "$DATABASE_URL" ]; then
    echo "ERROR: DATABASE_URL tidak di-set"
    echo "Pastikan file .env ada dan berisi DATABASE_URL"
    exit 1
fi

# Jalankan migration tool
go run ./cmd/migrate "$@"
