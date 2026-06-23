## ADDED Requirements

### Requirement: Sistem menyediakan endpoint GET list kategori

Sistem SHALL menyediakan endpoint untuk mengambil seluruh data kategori beserta jumlah produk terkait.

**Referensi:** KF-02, srs-fr.md § 5.1, AC-02.3

#### Scenario: GET list kategori berhasil
- **WHEN** Admin Farmasi mengirim GET request ke `/api/kategori` dengan Bearer token valid
- **THEN** Sistem mengembalikan response 200 dengan array kategori, masing-masing memiliki field: id, kode_kategori, nama_kategori, jumlah_produk, created_at, updated_at

#### Scenario: GET list kategori tanpa token
- **WHEN** Request ke `/api/kategori` tanpa Authorization header
- **THEN** Sistem mengembalikan 401 Unauthorized

### Requirement: Sistem menyediakan endpoint POST tambah kategori

Sistem SHALL menyediakan endpoint untuk menambah kategori baru dengan validasi nama unik (case-insensitive).

**Referensi:** KF-02, srs-fr.md § 5.2, BR-02.1, AC-02.1

#### Scenario: POST kategori baru berhasil
- **WHEN** Admin Farmasi mengirim POST ke `/api/kategori` dengan body `{ "nama_kategori": "Skincare" }` dan token valid
- **THEN** Sistem menyimpan kategori baru dengan kode_kategori yang di-generate otomatis, mengembalikan 201 Created dengan data kategori yang baru dibuat

#### Scenario: POST kategori dengan nama duplikat (case-insensitive)
- **WHEN** Admin Farmasi mengirim POST ke `/api/kategori` dengan nama yang sudah ada (misal "skincare" saat "Skincare" sudah ada)
- **THEN** Sistem mengembalikan 409 Conflict dengan message "Nama kategori sudah terdaftar dalam sistem."

#### Scenario: POST kategori dengan nama kosong
- **WHEN** Admin Farmasi mengirim POST ke `/api/kategori` dengan body `{ "nama_kategori": "" }` atau tanpa field nama_kategori
- **THEN** Sistem mengembalikan 400 Bad Request dengan message "Nama kategori wajib diisi"

### Requirement: Sistem menyediakan endpoint PUT ubah kategori

Sistem SHALL menyediakan endpoint untuk mengubah nama kategori berdasarkan ID dengan validasi duplikasi terhadap kategori lain.

**Referensi:** KF-02, srs-fr.md § 5.3, BR-02.1

#### Scenario: PUT ubah kategori berhasil
- **WHEN** Admin Farmasi mengirim PUT ke `/api/kategori/:id` dengan nama baru yang belum dipakai kategori lain
- **THEN** Sistem memperbarui nama kategori dan mengembalikan 200 OK dengan data kategori yang diperbarui

#### Scenario: PUT ubah kategori dengan nama duplikat kategori lain
- **WHEN** Admin Farmasi mengirim PUT ke `/api/kategori/:id` dengan nama yang sudah dipakai kategori lain
- **THEN** Sistem mengembalikan 409 Conflict dengan message "Nama kategori sudah terdaftar dalam sistem."

#### Scenario: PUT ubah kategori dengan ID tidak ditemukan
- **WHEN** Admin Farmasi mengirim PUT ke `/api/kategori/:id` dengan ID yang tidak ada di database
- **THEN** Sistem mengembalikan 404 Not Found dengan message "Kategori tidak ditemukan"

### Requirement: Sistem menyediakan endpoint DELETE hapus kategori

Sistem SHALL menyediakan endpoint untuk menghapus kategori dengan pengecekan integritas referensial — kategori tidak dapat dihapus jika masih ada produk terkait.

**Referensi:** KF-02, srs-fr.md § 5.4, BR-02.2, AC-02.2

#### Scenario: DELETE kategori berhasil
- **WHEN** Admin Farmasi mengirim DELETE ke `/api/kategori/:id` untuk kategori yang tidak memiliki produk terkait
- **THEN** Sistem menghapus kategori dan mengembalikan 200 OK dengan message "Kategori berhasil dihapus."

#### Scenario: DELETE kategori yang masih memiliki produk terkait
- **WHEN** Admin Farmasi mengirim DELETE ke `/api/kategori/:id` untuk kategori yang masih direferensikan oleh minimal satu produk
- **THEN** Sistem mengembalikan 409 Conflict dengan message "Kategori tidak dapat dihapus karena masih memiliki produk terkait."

#### Scenario: DELETE kategori dengan ID tidak ditemukan
- **WHEN** Admin Farmasi mengirim DELETE ke `/api/kategori/:id` dengan ID yang tidak ada
- **THEN** Sistem mengembalikan 404 Not Found dengan message "Kategori tidak ditemukan"

### Requirement: Kode kategori di-generate otomatis oleh sistem

Sistem SHALL men-generate kode kategori secara otomatis saat kategori baru dibuat.

**Referensi:** KF-02, srs-fr.md § 5.1

#### Scenario: Kode kategori unik di-generate saat tambah kategori
- **WHEN** Kategori baru berhasil disimpan
- **THEN** Sistem menetapkan kode_kategori yang unik dan konsisten (format: KTG-XXX dengan auto-increment sequence)
