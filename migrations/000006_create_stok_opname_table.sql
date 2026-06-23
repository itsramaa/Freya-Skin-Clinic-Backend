-- +goose Up
CREATE TABLE IF NOT EXISTS stok_opname (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_user         UUID NOT NULL,
    tanggal_opname  DATE NOT NULL DEFAULT CURRENT_DATE,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF' CHECK (status IN ('AKTIF','SELESAI','DIBATALKAN')),
    catatan         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_opname_id_user FOREIGN KEY (id_user) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS detail_opname (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_opname           UUID NOT NULL,
    id_batch            UUID NOT NULL,
    id_kemasan_terbuka  UUID,
    stok_sistem         DECIMAL(10,3) NOT NULL,
    stok_fisik          DECIMAL(10,3) NOT NULL CHECK (stok_fisik >= 0),
    selisih             DECIMAL(10,3) NOT NULL,
    keterangan          TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_detail_id_opname  FOREIGN KEY (id_opname) REFERENCES stok_opname(id) ON DELETE CASCADE,
    CONSTRAINT fk_detail_id_batch   FOREIGN KEY (id_batch)  REFERENCES batch_stok(id),
    CONSTRAINT fk_detail_id_kemasan FOREIGN KEY (id_kemasan_terbuka) REFERENCES kemasan_terbuka(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_detail_opname_id_opname ON detail_opname(id_opname);
CREATE INDEX IF NOT EXISTS idx_stok_opname_status       ON stok_opname(status);

-- +goose Down
DROP TABLE IF EXISTS detail_opname;
DROP TABLE IF EXISTS stok_opname;
