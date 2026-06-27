-- +goose Up
ALTER TABLE detail_opname
    ADD COLUMN IF NOT EXISTS sisa_isi_sistem   DECIMAL(10,3),
    ADD COLUMN IF NOT EXISTS sisa_isi_fisik    DECIMAL(10,3),
    ADD COLUMN IF NOT EXISTS selisih_sisa_isi  DECIMAL(10,3);

-- +goose Down
ALTER TABLE detail_opname
    DROP COLUMN IF EXISTS sisa_isi_sistem,
    DROP COLUMN IF EXISTS sisa_isi_fisik,
    DROP COLUMN IF EXISTS selisih_sisa_isi;
