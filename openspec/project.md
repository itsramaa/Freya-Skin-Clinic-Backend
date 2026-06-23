# Project: Freya Skin Clinic — Backend

## Overview

Sistem Manajemen Stok Farmasi Internal Freya Skin Clinic adalah aplikasi web berbasis SPA (Single Page Application) yang menggantikan proses pengelolaan persediaan semi-manual berbasis catatan fisik dan Microsoft Excel di Farmasi Internal Freya Skin Clinic, Sumedang.

Proyek ini adalah **project baru** — implementasi dimulai dari nol.

---

## Domain & Business Context

- **Klien:** Freya Skin Clinic, Sumedang
- **Pengguna:** Admin Farmasi (single user, single role)
- **Domain:** Manajemen stok farmasi internal klinik kecantikan
- **Kategori produk:** Skincare, Injectable, Obat, Threadlift, Facial IPL Laser
- **Metode stok:** FEFO (First Expired First Out)
- **Fitur khusus:** BUD (Beyond Use Date) untuk kemasan partial use yang sudah dibuka

---

## Problems Solved

| Kode | Masalah Lama |
|------|-------------|
| P-01 | Tidak ada pencatatan expired date & batch; FEFO manual |
| P-02 | Pencatatan full use & partial use serta sisa stok manual |
| P-03 | BUD tidak terdokumentasi |
| P-04 | Monitoring, opname, dan rekap stok manual tanpa histori |
| P-05 | Data stok tersimpan di file terpisah per bulan |

---

## Tech Stack

| Layer | Teknologi |
|-------|-----------|
| Language | Go 1.22+ |
| Web Framework | Fiber v2 |
| Database | PostgreSQL 15+ |
| ORM/Driver | pgx/v5 |
| Auth | JWT (golang-jwt) |
| Password Hash | bcrypt |
| Config | env vars (joho/godotenv) |
| Migration | golang-migrate |
| Validator | go-playground/validator |

---

## Architecture

- **Pattern:** Clean Architecture (Handler → Service → Repository)
- **DB:** PostgreSQL dengan ACID transaction
- **Background Worker:** Goroutine untuk auto-update status batch & BUD
- **Auth:** Middleware JWT; semua endpoint protected kecuali `/api/auth/login`

```
HTTP Request → Fiber Router → Auth Middleware → Handler → Service → Repository → PostgreSQL
                                                            ↓
                                                    Background Worker (goroutine)
```

---

## Project Structure

```
backend/
├── cmd/
│   ├── api/                # Main HTTP server entry point
│   ├── api-validate/       # OpenAPI validation CLI
│   ├── hashgen/            # Password hash generator CLI
│   └── migrate/            # Database migration runner CLI
├── internal/
│   ├── config/             # Env/config loader
│   ├── handler/            # HTTP handlers (controllers)
│   ├── service/            # Business logic layer
│   ├── repository/         # Data access layer
│   ├── model/              # Domain models & DTOs
│   ├── middleware/         # Auth, RBAC middleware
│   ├── router/             # Route registration
│   └── pkg/                # Shared internal packages
├── migrations/             # SQL migration files
├── api/                    # OpenAPI spec & docs
└── docs/                   # SRS & evidence docs
```

---

## Module / Feature List

| Kode | Modul | Handler | Service | Repository |
|------|-------|---------|---------|------------|
| KF-01 | Autentikasi | auth_handler | auth_service | user_repository |
| KF-02 | Kelola Kategori | kategori_handler | kategori_service | kategori_repository |
| KF-03 | Kelola Produk | produk_handler | produk_service | produk_repository |
| KF-04 | Kelola Stok Masuk | stok_masuk_handler | stok_masuk_service | stok_masuk_repository, batch_repository |
| KF-05 | Kelola Stok Keluar | stok_keluar_handler | stok_keluar_service | stok_keluar_repository, batch_repository |
| KF-06 | Penerapan FEFO | (via stok_keluar) | stok_keluar_service | batch_repository |
| KF-07 | Kelola BUD | (via stok_keluar) | stok_keluar_service | kemasan_terbuka_repository |
| KF-08 | Monitoring Stok | monitoring_handler | monitoring_service | batch_repository |
| KF-09 | Stock Opname | opname_handler | opname_service | opname_repository |
| KF-10 | Laporan Stok | laporan_handler | laporan_service | laporan_repository |

---

## Database Tables

| Tabel | Fungsi |
|-------|--------|
| users | Akun admin farmasi |
| kategori | Kategori produk (Skincare, Injectable, dll) |
| produk | Master produk farmasi |
| stok_masuk | Transaksi penerimaan stok |
| batch_stok | Batch stok dengan expired date & status |
| stok_keluar | Transaksi pengeluaran stok |
| kemasan_terbuka | Produk partial use dengan BUD |
| stok_opname | Sesi stock opname |
| detail_opname | Detail per-item opname & selisih |

---

## Key Conventions

- **Bahasa:** Label & pesan error dalam Bahasa Indonesia
- **Commit:** Conventional commits (`feat`, `fix`, `refactor`, dll)
- **Transaction:** Semua operasi multi-tabel dibungkus DB transaction (ACID)
- **Error handling:** Custom error types; respons terstruktur `{ success, message, data, errors }`
- **Migration:** Sequential numbering (`000001_`, `000002_`, dst)
- **Per-domain files:** `{domain}_handler.go`, `{domain}_service.go`, `{domain}_repo.go`

---

## API Contract

- **Base URL:** `/api`
- **Auth header:** `Authorization: Bearer <token>` (JWT)
- **Format:** JSON, ISO 8601 untuk tanggal
- **Standard response:** `{ "success": bool, "message": string, "data": any, "errors": any }`
- **Detail:** lihat `docs/srs-api.md`

---

## Background Workers

| Worker | Trigger | Fungsi |
|--------|---------|--------|
| Batch status monitor | Periodic (goroutine) | Update batch_stok.status ke KADALUWARSA jika expired_date < now |
| BUD monitor | Periodic (goroutine) | Nonaktifkan kemasan_terbuka jika bud < now |

---

## Non-Functional Requirements Summary

| Kode | Atribut | Target |
|------|---------|--------|
| KNF-02 | Security | JWT auth, bcrypt hash, HTTPS di produksi |
| KNF-04 | Reliability | DB transaction untuk semua operasi multi-tabel |
| KNF-05 | Maintainability | Data historis dipertahankan (soft delete / audit trail) |
| KNF-06 | Performance | Response API ≤ 2 detik |
| KNF-07 | Availability | Background worker berjalan terus tanpa block request |

---

## Project Status

- **Phase:** Initial development — belum ada implementasi
- **Started:** 2026-06-22
- **SRS:** `docs/srs-overview.md`, `docs/srs-fr.md`, `docs/srs-nfr.md`, `docs/srs-backend.md`, `docs/srs-api.md`, `docs/srs-database.md`
