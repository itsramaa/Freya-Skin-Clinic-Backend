-- +goose Up
CREATE TABLE IF NOT EXISTS kategori (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode_kategori VARCHAR(10) UNIQUE NOT NULL,
    nama_kategori VARCHAR(100) UNIQUE NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_kategori_nama ON kategori(LOWER(nama_kategori));

-- +goose Down
DROP TABLE IF EXISTS kategori;
