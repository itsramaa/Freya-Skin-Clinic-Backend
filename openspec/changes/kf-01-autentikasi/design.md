## Context

KF-01 adalah modul autentikasi backend untuk Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic. Backend dibangun dengan Go + Fiber mengikuti Clean Architecture (Handler → Service → Repository).

Ini adalah modul pertama yang diimplementasi karena semua endpoint lain (KF-02 s.d. KF-10) bergantung pada middleware JWT yang didefinisikan di sini. Sistem hanya memiliki satu aktor (Admin Farmasi) — tidak ada RBAC granular.

Alur khusus: login pertama dengan password default mengembalikan `is_default_password: true` sehingga frontend dapat memaksa penggantian password sebelum dashboard dapat diakses.

## Goals / Non-Goals

**Goals:**
- Migration tabel `users` dengan seed data admin default
- Repository layer: `FindByUsername`, `UpdatePassword`
- Service layer: `Login` (validate credentials, generate JWT), `ChangePassword` (hash + update DB)
- Handler layer: `POST /api/auth/login`, `PUT /api/auth/password`
- JWT utility package: `GenerateToken`, `ValidateToken`
- Auth middleware: parse Bearer token, validate, inject user_id ke Fiber context
- Router setup: register auth routes (public) dan apply middleware ke protected routes
- Hash utility: bcrypt wrapper untuk hash dan compare password
- Standard response format: `{ success, message, data, errors }`

**Non-Goals:**
- Multi-role / RBAC
- Refresh token / token rotation
- Rate limiting pada login endpoint
- Audit log untuk login attempts

## Decisions

### 1. JWT Library: golang-jwt/jwt

**Keputusan:** Gunakan `github.com/golang-jwt/jwt/v5` untuk generate dan validate JWT.

**Alasan:** Library standar, actively maintained, mendukung claims validation (expiry, issuer). Compatible dengan Fiber middleware pattern.

**Alternatif:** `github.com/lestrrat-go/jwx` — terlalu kompleks untuk use case single-user ini.

### 2. Password Hashing: bcrypt dengan cost factor 12

**Keputusan:** Gunakan `golang.org/x/crypto/bcrypt` dengan cost factor 12.

**Alasan:** bcrypt adalah standar industri untuk password hashing. Cost factor 12 memberikan keseimbangan antara keamanan dan performa — masih < 500ms pada hardware modern, cukup lambat untuk brute-force attack.

**Alternatif:** argon2id — lebih aman secara teoritis tetapi lebih kompleks untuk di-configure dan bukan kebutuhan critical untuk sistem internal single-user.

### 3. JWT Claims: Minimal claims (user_id, username, exp)

**Keputusan:** Hanya simpan `user_id`, `username`, dan `exp` di JWT claims.

**Alasan:** Tidak ada role-based data yang perlu di-encode. Minimal claims mengurangi token size dan attack surface. Expiry 24 jam untuk sistem internal klinik yang beroperasi single-shift.

### 4. JWT Secret: Environment variable

**Keputusan:** JWT signing key dibaca dari environment variable `JWT_SECRET`.

**Alasan:** Secret tidak boleh di-hardcode di source code. Mudah dirotasi tanpa perubahan kode. Mengikuti 12-factor app principles.

### 5. Middleware Placement: Global middleware via Fiber App.Use()

**Keputusan:** Auth middleware diaplikasikan secara global di router level, kecuali route login yang di-whitelist.

**Alasan:** Semua endpoint KF-02 s.d. KF-10 akan protected. Lebih aman daripada opt-in per-route (mencegah kelupaan protect route baru).

**Implementasi:** Router mendaftarkan `/api/auth/login` sebagai public route. Semua route lain di-group dan di-apply auth middleware.

### 6. Error Response: Tidak reveal username exists/not exists

**Keputusan:** Login failure selalu mengembalikan "Kredensial tidak valid" tanpa membedakan apakah username tidak ditemukan atau password salah.

**Alasan:** Mencegah username enumeration attack.

## API Endpoints

| Method | Path | Auth | Request Body | Response |
|--------|------|------|-------------|---------|
| `POST` | `/api/auth/login` | Public | `{ username, password }` | `{ success, data: { token, user: { id, username }, is_default_password } }` |
| `PUT` | `/api/auth/password` | Bearer JWT | `{ password_baru }` | `{ success, message }` |

**Standard Response Format:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": { ... },
  "errors": null
}
```

## Database Schema

```sql
-- Migration: 000001_create_users_table.up.sql
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username    VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_default_password BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);

-- Seed data (dalam migration yang sama atau terpisah)
INSERT INTO users (username, password_hash, is_default_password)
VALUES ('admin', '<bcrypt_hash_of_admin123>', true);
```

## Architecture Flow

```
POST /api/auth/login
  → auth_handler.Login()
      → Validate request body (username, password wajib)
      → auth_service.Login(ctx, username, password)
          → user_repository.FindByUsername(ctx, username)
              → SELECT * FROM users WHERE username = $1
          → bcrypt.CompareHashAndPassword(storedHash, inputPassword)
          → jwt.GenerateToken(userID, username)
          → return { token, user, is_default_password }
      → Return 200 { success: true, data: { token, user, is_default_password } }

PUT /api/auth/password
  → auth_middleware (validate Bearer token, inject user_id to context)
  → auth_handler.ChangePassword()
      → Validate request body (password_baru wajib, min 8 chars)
      → Extract user_id dari Fiber context
      → auth_service.ChangePassword(ctx, userID, passwordBaru)
          → bcrypt.GenerateFromPassword(passwordBaru, cost=12)
          → user_repository.UpdatePassword(ctx, userID, passwordHash)
              → UPDATE users SET password_hash=$1, is_default_password=false, updated_at=NOW() WHERE id=$2
      → Return 200 { success: true, message: "Password berhasil diperbarui" }

Auth Middleware
  → Parse "Authorization: Bearer <token>"
  → jwt.ValidateToken(token)
  → On error → Return 401
  → On success → c.Locals("user_id", claims.UserID) → c.Next()
```

## File Structure

```
internal/
├── handler/
│   └── auth_handler.go         # LoginHandler, ChangePasswordHandler
├── service/
│   └── auth_service.go         # Login(), ChangePassword()
├── repository/
│   └── user_repository.go      # FindByUsername(), UpdatePassword()
├── model/
│   └── user.go                 # User struct, LoginRequest, LoginResponse, ChangePasswordRequest
├── middleware/
│   └── auth_middleware.go      # JWTMiddleware()
├── router/
│   └── router.go               # Route registration
└── pkg/
    ├── jwt/
    │   └── jwt.go              # GenerateToken(), ValidateToken()
    ├── hash/
    │   └── hash.go             # HashPassword(), ComparePassword()
    └── response/
        └── response.go         # Standard response helpers: Success(), Error()

migrations/
├── 000001_create_users_table.up.sql
└── 000001_create_users_table.down.sql
```

## Risks / Trade-offs

- **[Risk] Token tidak dapat di-revoke sebelum expiry** → Mitigasi: Expiry 24 jam adalah trade-off yang diterima untuk sistem internal. Jika user ganti password, token lama masih valid hingga expired. Acceptable untuk sistem single-user internal.

- **[Risk] bcrypt cost=12 membuat login ~200-400ms** → Mitigasi: Acceptable untuk sistem internal dengan satu user. Bukan bottleneck yang perlu dioptimasi.

- **[Risk] Seed password default di migration** → Mitigasi: Gunakan `cmd/hashgen` CLI untuk generate hash yang sudah di-hardcode di migration. Dokumentasikan bahwa password default harus diganti saat deploy pertama.

## Migration Plan

1. Jalankan `migrate up` untuk create tabel users dan insert seed data
2. Jalankan `cmd/hashgen` untuk generate hash password default jika perlu update seed
3. Deploy backend server
4. Verifikasi `POST /api/auth/login` dengan Postman/curl
5. Verifikasi middleware 401 pada endpoint protected

## Open Questions

- Durasi JWT expiry: 24 jam (rekomendasi untuk single-shift internal system)?
- Password default: "admin123" atau dikonfigurasi via env variable `DEFAULT_ADMIN_PASSWORD`?
- Apakah seed data dimasukkan dalam migration atau script terpisah?
