## 1. Migration & Database

- [x] 1.1 Buat `migrations/000001_create_users_table.up.sql`
- [x] 1.2 Buat `migrations/000001_create_users_table.down.sql`
- [x] 1.3 Seed data INSERT admin user (password default "sihuni123" bcrypt hash)
- [ ] 1.4 Jalankan `go run cmd/migrate/main.go up` untuk verifikasi migration berjalan tanpa error

## 2. Shared Packages (internal/pkg)

- [x] 2.1 Buat `internal/pkg/response/response.go`
- [x] 2.2 Buat `internal/pkg/hash/hash.go`
- [x] 2.3 Buat `internal/pkg/jwt/jwt.go`

## 3. Config

- [x] 3.1 Buat `internal/config/config.go`

## 4. Model & DTO

- [x] 4.1 Buat `internal/model/user.go`

## 5. Repository Layer

- [x] 5.1 Buat `internal/repository/user_repository.go` — FindByUsername
- [x] 5.2 Tambah method UpdatePassword

## 6. Service Layer

- [x] 6.1 Buat `internal/service/auth_service.go` — Login
- [x] 6.2 Tambah method ChangePassword

## 7. Handler Layer

- [x] 7.1 Buat `internal/handler/auth_handler.go`
- [x] 7.2 Implementasi Login handler
- [x] 7.3 Implementasi ChangePassword handler

## 8. Middleware

- [x] 8.1 Buat `internal/middleware/auth_middleware.go`

## 9. Router & Main

- [x] 9.1 Buat `internal/router/router.go`
- [x] 9.2 Buat `cmd/api/main.go`

## 10. CLI Tools

- [x] 10.1 `cmd/hashgen/main.go` — sudah ada
- [x] 10.2 `cmd/migrate/main.go` — sudah ada

## 11. Verifikasi

- [x] 11.1 Verifikasi `POST /api/auth/login` dengan kredensial valid
- [x] 11.2 Verifikasi `POST /api/auth/login` dengan kredensial salah → 401
- [x] 11.3 Verifikasi `PUT /api/auth/password` dengan token valid → 200
- [x] 11.4 Verifikasi `PUT /api/auth/password` tanpa token → 401
- [x] 11.5 Verifikasi request endpoint lain tanpa token → 401
- [x] 11.6 Verifikasi password tidak muncul plain text di response
- [x] 11.7 Jalankan `go build ./...` — tidak ada compile error
