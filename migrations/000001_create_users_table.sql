-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username            VARCHAR(50) UNIQUE NOT NULL,
    password_hash       VARCHAR(255) NOT NULL,
    is_default_password BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Seed: admin user dengan password default "admin"
INSERT INTO users (username, password_hash, is_default_password)
VALUES (
    'admin',
    '$2a$10$deA8Pp1xEpWu1iTjTy6O0.HI/cmdnsvE.2VYEBXSUBSOTYWOkCBuC',
    true
) ON CONFLICT (username) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS users;
