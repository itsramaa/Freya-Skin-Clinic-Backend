# KF-01 Autentikasi - Backend

## Summary

Implementasi modul autentikasi backend untuk login Admin Farmasi, validasi JWT token, dan penggantian password paksa.

## Motivation

Sistem memerlukan mekanisme autentikasi yang aman untuk membatasi akses hanya kepada Admin Farmasi yang berwenang. Tanpa autentikasi, seluruh data farmasi klinik dapat diakses oleh siapa saja.

Modul ini merupakan fondasi untuk semua fitur lain (KF-02 s.d. KF-10) yang memerlukan proteksi akses.

## Goals

- Implementasi endpoint login yang memvalidasi username/password dan mengembalikan JWT token
- Implementasi middleware JWT yang memproteksi semua endpoint kecuali `/api/auth/login`
- Implementasi endpoint ganti password yang wajib diakses saat login pertama kali
- Deteksi password default (is_default_password) untuk memaksa penggantian password
- Hash password menggunakan bcrypt untuk keamanan

## Non-Goals

- Multi-role atau RBAC (sistem hanya memiliki 1 role: Admin Farmasi)
- Refresh token mechanism
- Rate limiting pada login endpoint
- Password reset via email (single user, tidak perlu)

## Capabilities

### ADDED Capabilities

**auth-management**: Endpoint untuk login dan ganti password Admin Farmasi
- `POST /api/auth/login` - Autentikasi dengan username/password, return JWT token
- `PUT /api/auth/password` - Ganti password (wajib untuk login pertama)

**auth-middleware**: Middleware untuk proteksi endpoint dengan JWT
- Validasi Bearer token pada setiap request
- Extract user_id dari token claims
- Return 401 jika token invalid/expired

## Technical Approach

### Migration

- Tabel `users` dengan kolom:
  - `id` (PK, UUID)
  - `username` (UNIQUE, VARCHAR 50)
  - `password_hash` (VARCHAR 255)
  - `is_default_password` (BOOLEAN, default true)
  - `created_at`, `updated_at` (TIMESTAMP)

### Implementation

**Repository Layer** (`user_repository.go`):
- `FindByUsername(ctx, username)` - Query user by username
- `UpdatePassword(ctx, userID, passwordHash)` - Update password hash dan set `is_default_password = false`

**Service Layer** (`auth_service.go`):
- `Login(ctx, username, password)` - Validate credentials, generate JWT, return token + is_default_password
- `ChangePassword(ctx, userID, newPassword)` - Hash new password, update DB

**Handler Layer** (`auth_handler.go`):
- `POST /api/auth/login` - Parse request, call service, return token
- `PUT /api/auth/password` - Parse request, extract user_id from context, call service

**Middleware** (`auth_middleware.go`):
- Parse Authorization header (Bearer token)
- Validate JWT signature dan expiry
- Set user_id ke context untuk handler berikutnya
- Return 401 jika invalid

**JWT Utility** (`internal/pkg/jwt/`):
- `GenerateToken(userID, username)` - Create JWT dengan claims
- `ValidateToken(tokenString)` - Parse dan validate token

### Database Schema

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_default_password BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
```

## Risks & Mitigations

**Risk**: JWT token disimpan di frontend memory (Zustand)
**Mitigation**: Token memiliki expiry time yang relatif pendek (e.g., 24 hours) untuk membatasi window of exposure

**Risk**: Password default yang sama untuk semua user
**Mitigation**: Sistem memaksa penggantian password saat login pertama, dan password default hanya diketahui oleh admin yang setup sistem

## Open Questions

- Berapa durasi expiry JWT yang ideal? (rekomendasi: 24 hours untuk single-user internal system)
- Apakah perlu seed data initial user di migration? (e.g., username "admin" dengan password default)
