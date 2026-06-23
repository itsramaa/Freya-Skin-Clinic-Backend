# Project Structure — Freya Skin Clinic Backend

Freya-Skin-Clinic-Backend/
├── cmd/ # Application entry points (main packages)
│ ├── api/ # Main HTTP API server
│ ├── api-validate/ # OpenAPI spec validation CLI
│ ├── hashgen/ # Password hash generator CLI
│ └── migrate/ # Database migration runner CLI
│
├── internal/ # Private application code (Go convention)
│ ├── config/ # App configuration loader (env/yaml)
│ ├── handler/ # HTTP handlers (controllers)
│ ├── service/ # Business logic layer
│ ├── repository/ # Data access layer (database queries)
│ ├── model/ # Domain models & DTOs
│ ├── middleware/ # HTTP middleware (auth, RBAC)
│ ├── router/ # Route registration
│ └── pkg/ # Shared internal packages/utilities
│
├── migrations/ # SQL migration files (sequential numbering)
├── api/ # API documentation & contracts (OpenAPI, Scalar)
├── docs/ # Project documentation (SRS, test plans, QA)
│ └── evidence/ # Test evidence & artifacts
├── openspec/ # OpenSpec workflow documentation
├── .trae/ # Trae IDE configuration
│ ├── rules/ # AI agent steering rules
│ └── skills/ # Installed AI skills
└── tmp/ # Hot-reload build artifacts (gitignored)

```

## Architecture Layers

Proyek ini mengikuti **Clean Architecture** dengan layer dependency sebagai berikut:

```

Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → Database
↓ ↓ ↓
Middleware Model/DTO pkg (utilities)

```

| Layer          | Directory              | Tanggung Jawab                                       |
| -------------- | ---------------------- | ---------------------------------------------------- |
| **Handler**    | `internal/handler/`    | Parse HTTP request, validasi input, return response  |
| **Service**    | `internal/service/`    | Business logic, orkestrasi antar repo, aturan domain |
| **Repository** | `internal/repository/` | Query database, CRUD operations                      |
| **Model**      | `internal/model/`      | Struct domain, DTO, error types                      |
| **Middleware** | `internal/middleware/` | Auth, RBAC, request-scoped logic                     |
| **Pkg**        | `internal/pkg/`        | Utility packages (JWT, validator, response, dll)     |

## Key Conventions

- **Setiap domain** punya file terpisah di setiap layer: `{domain}_handler.go`, `{domain}_service.go`, `{domain}_repo.go`, `{domain}.go` (model).
- **Test file** mengikuti konvensi `{file}_test.go` di samping file sumber.
- **Migration** menggunakan sequential numbering (`000001_`, `000002_`, dst).
- **cmd/** berisi entry points saja — semua logic ada di `internal/`.
- **internal/pkg/** untuk utility yang dipakai lintas domain, **bukan** business logic.
```
