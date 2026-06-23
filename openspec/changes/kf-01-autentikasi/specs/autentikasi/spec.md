# Autentikasi Specifications

## Requirements

### Requirement: User dapat login dengan kredensial valid

Sistem HARUS menyediakan endpoint untuk autentikasi Admin Farmasi menggunakan username dan password, serta mengembalikan JWT token yang valid.

#### Scenario: Login berhasil dengan kredensial valid
- **WHEN** Admin Farmasi mengirim POST request ke `/api/auth/login` dengan username dan password yang valid
- **THEN** Sistem memvalidasi kredensial, generate JWT token dengan claims (user_id, username), dan mengembalikan response dengan format: `{ token: string, user: { id, username }, is_default_password: boolean }` dengan status 200 OK

#### Scenario: Login gagal dengan username tidak ditemukan
- **WHEN** Admin Farmasi mengirim POST request ke `/api/auth/login` dengan username yang tidak terdaftar
- **THEN** Sistem mengembalikan response error dengan status 401 Unauthorized dan message "Kredensial tidak valid"

#### Scenario: Login gagal dengan password salah
- **WHEN** Admin Farmasi mengirim POST request ke `/api/auth/login` dengan username valid tetapi password salah
- **THEN** Sistem mengembalikan response error dengan status 401 Unauthorized dan message "Kredensial tidak valid"

#### Scenario: Login dengan password default
- **WHEN** Admin Farmasi login dengan kredensial valid dan user memiliki `is_default_password = true`
- **THEN** Sistem mengembalikan response dengan `is_default_password: true` untuk memicu redirect ke halaman ganti password di frontend

### Requirement: User wajib mengganti password saat login pertama

Sistem HARUS menyediakan endpoint untuk mengganti password dan memaksa penggantian password saat login pertama kali dengan password default.

#### Scenario: Ganti password berhasil
- **WHEN** Admin Farmasi mengirim PUT request ke `/api/auth/password` dengan Bearer token valid dan new_password yang memenuhi syarat (minimal 8 karakter)
- **THEN** Sistem melakukan hash password baru, update database dengan password_hash baru, set `is_default_password = false`, dan mengembalikan response success dengan status 200 OK

#### Scenario: Ganti password gagal karena token tidak valid
- **WHEN** Admin Farmasi mengirim PUT request ke `/api/auth/password` tanpa Bearer token atau dengan token yang invalid
- **THEN** Sistem mengembalikan response error dengan status 401 Unauthorized

#### Scenario: Ganti password gagal karena password terlalu pendek
- **WHEN** Admin Farmasi mengirim PUT request ke `/api/auth/password` dengan new_password yang kurang dari 8 karakter
- **THEN** Sistem mengembalikan response error dengan status 400 Bad Request dan message "Password minimal 8 karakter"

#### Scenario: Ganti password gagal karena password kosong
- **WHEN** Admin Farmasi mengirim PUT request ke `/api/auth/password` dengan new_password kosong atau null
- **THEN** Sistem mengembalikan response error dengan status 400 Bad Request dan message "Password tidak boleh kosong"

### Requirement: Middleware JWT melindungi semua endpoint

Sistem HARUS menyediakan middleware yang memvalidasi JWT token pada setiap request ke endpoint yang dilindungi dan menolak request tanpa token yang valid.

#### Scenario: Request dengan token valid
- **WHEN** User mengirim request ke endpoint protected dengan header `Authorization: Bearer <valid_token>`
- **THEN** Middleware memvalidasi token, extract user_id dari claims, set user_id ke context, dan meneruskan request ke handler berikutnya

#### Scenario: Request tanpa Authorization header
- **WHEN** User mengirim request ke endpoint protected tanpa header Authorization
- **THEN** Middleware mengembalikan response error dengan status 401 Unauthorized dan message "Authorization token required"

#### Scenario: Request dengan token format salah
- **WHEN** User mengirim request ke endpoint protected dengan header Authorization yang tidak mengikuti format "Bearer <token>"
- **THEN** Middleware mengembalikan response error dengan status 401 Unauthorized dan message "Invalid authorization header format"

#### Scenario: Request dengan token invalid
- **WHEN** User mengirim request ke endpoint protected dengan token yang signature-nya tidak valid atau sudah di-tamper
- **THEN** Middleware mengembalikan response error dengan status 401 Unauthorized dan message "Invalid token"

#### Scenario: Request dengan token expired
- **WHEN** User mengirim request ke endpoint protected dengan token yang sudah expired
- **THEN** Middleware mengembalikan response error dengan status 401 Unauthorized dan message "Token expired"

### Requirement: Password disimpan dalam bentuk hash

Sistem HARUS menyimpan password dalam bentuk hash menggunakan algoritma bcrypt untuk keamanan.

#### Scenario: Password di-hash saat login validation
- **WHEN** User login dengan password plain text
- **THEN** Sistem membandingkan password dengan password_hash di database menggunakan bcrypt.compare, bukan plain text comparison

#### Scenario: Password di-hash saat ganti password
- **WHEN** User mengganti password dengan new_password plain text
- **THEN** Sistem melakukan hash dengan bcrypt (cost factor minimal 10) sebelum menyimpan ke database

### Requirement: Seed data initial admin user

Sistem HARUS menyediakan seed data untuk user Admin Farmasi pertama kali dengan password default.

#### Scenario: Migration membuat user admin default
- **WHEN** Migration dijalankan untuk pertama kali
- **THEN** Sistem membuat user dengan username "admin" dan password hash dari default password (misalnya "admin123"), dengan `is_default_password = true`
