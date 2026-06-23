-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS session_id VARCHAR(36);

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS session_id;
