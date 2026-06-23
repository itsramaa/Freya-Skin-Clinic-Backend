-- +goose Up
CREATE TYPE pola_penggunaan_enum AS ENUM ('FULL_USE', 'PARTIAL_USE');

CREATE TABLE IF NOT EXISTS produk (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode_produk      VARCHAR(20) UNIQUE NOT NULL,
    nama_produk      VARCHAR(200) NOT NULL,
    id_kategori      UUID NOT NULL,
    bentuk_kemasan   VARCHAR(50) NOT NULL,
    satuan_isi       VARCHAR(20) NOT NULL,
    isi_per_kemasan  DECIMAL(10,3),
    pola_penggunaan  pola_penggunaan_enum NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_produk_id_kategori FOREIGN KEY (id_kategori)
        REFERENCES kategori(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_produk_id_kategori ON produk(id_kategori);
CREATE INDEX IF NOT EXISTS idx_produk_nama ON produk(nama_produk);

-- +goose Down
DROP TABLE IF EXISTS produk;
DROP TYPE IF EXISTS pola_penggunaan_enum;
