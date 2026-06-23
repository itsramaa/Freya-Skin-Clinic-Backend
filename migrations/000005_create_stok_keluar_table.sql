-- +goose Up
CREATE TABLE IF NOT EXISTS kemasan_terbuka (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_batch        UUID NOT NULL,
    tanggal_dibuka  DATE NOT NULL,
    bud             DATE NOT NULL,
    isi_awal        DECIMAL(10,3) NOT NULL CHECK (isi_awal > 0),
    isi_tersisa     DECIMAL(10,3) NOT NULL CHECK (isi_tersisa >= 0),
    status_bud      VARCHAR(20) NOT NULL DEFAULT 'AKTIF' CHECK (status_bud IN ('AKTIF','KADALUWARSA')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_kemasan_terbuka_id_batch FOREIGN KEY (id_batch)
        REFERENCES batch_stok(id) ON DELETE RESTRICT,
    CONSTRAINT uq_kemasan_terbuka_batch UNIQUE (id_batch)
);

CREATE INDEX IF NOT EXISTS idx_kemasan_terbuka_id_batch   ON kemasan_terbuka(id_batch);
CREATE INDEX IF NOT EXISTS idx_kemasan_terbuka_status_bud ON kemasan_terbuka(status_bud);

CREATE TABLE IF NOT EXISTS stok_keluar (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_produk              UUID NOT NULL,
    id_batch               UUID NOT NULL,
    id_kemasan_terbuka     UUID,
    id_user                UUID NOT NULL,
    tanggal_penggunaan     DATE NOT NULL,
    jumlah_kemasan_dipakai INTEGER DEFAULT 0,
    jumlah_isi_dipakai     DECIMAL(10,3) DEFAULT 0,
    keterangan             TEXT,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_stok_keluar_id_produk  FOREIGN KEY (id_produk)  REFERENCES produk(id),
    CONSTRAINT fk_stok_keluar_id_batch   FOREIGN KEY (id_batch)   REFERENCES batch_stok(id),
    CONSTRAINT fk_stok_keluar_id_user    FOREIGN KEY (id_user)    REFERENCES users(id),
    CONSTRAINT fk_stok_keluar_id_kemasan FOREIGN KEY (id_kemasan_terbuka) REFERENCES kemasan_terbuka(id)
);

CREATE INDEX IF NOT EXISTS idx_stok_keluar_id_produk ON stok_keluar(id_produk);
CREATE INDEX IF NOT EXISTS idx_stok_keluar_tanggal   ON stok_keluar(tanggal_penggunaan);

-- +goose Down
DROP TABLE IF EXISTS stok_keluar;
DROP TABLE IF EXISTS kemasan_terbuka;
