-- +goose Up
CREATE TABLE IF NOT EXISTS batch_stok (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_produk           UUID NOT NULL,
    kode_batch          VARCHAR(30) UNIQUE NOT NULL,
    expired_date        DATE NOT NULL,
    stok_kemasan        INTEGER NOT NULL DEFAULT 0 CHECK (stok_kemasan >= 0),
    total_isi_tersedia  DECIMAL(12,3) NOT NULL DEFAULT 0 CHECK (total_isi_tersedia >= 0),
    status              VARCHAR(20) NOT NULL DEFAULT 'AKTIF' CHECK (status IN ('AKTIF','HABIS','KADALUWARSA')),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_batch_id_produk FOREIGN KEY (id_produk)
        REFERENCES produk(id) ON DELETE RESTRICT,
    CONSTRAINT uq_batch_produk_expired UNIQUE (id_produk, expired_date)
);

CREATE INDEX IF NOT EXISTS idx_batch_stok_id_produk ON batch_stok(id_produk);
CREATE INDEX IF NOT EXISTS idx_batch_stok_status ON batch_stok(status);
CREATE INDEX IF NOT EXISTS idx_batch_stok_expired_date ON batch_stok(expired_date);

CREATE TABLE IF NOT EXISTS stok_masuk (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_produk           UUID NOT NULL,
    id_batch            UUID NOT NULL,
    id_user             UUID NOT NULL,
    tanggal_penerimaan  DATE NOT NULL,
    jumlah_kemasan      INTEGER NOT NULL CHECK (jumlah_kemasan > 0),
    total_isi_masuk     DECIMAL(12,3) NOT NULL CHECK (total_isi_masuk > 0),
    keterangan          TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_stok_masuk_id_produk FOREIGN KEY (id_produk) REFERENCES produk(id),
    CONSTRAINT fk_stok_masuk_id_batch  FOREIGN KEY (id_batch)  REFERENCES batch_stok(id),
    CONSTRAINT fk_stok_masuk_id_user   FOREIGN KEY (id_user)   REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_stok_masuk_id_produk ON stok_masuk(id_produk);
CREATE INDEX IF NOT EXISTS idx_stok_masuk_tanggal   ON stok_masuk(tanggal_penerimaan);

-- +goose Down
DROP TABLE IF EXISTS stok_masuk;
DROP TABLE IF EXISTS batch_stok;
