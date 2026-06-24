## MODIFIED Requirements

### Requirement: Response monitoring menyertakan satuan isi produk
Sistem SHALL menyertakan field `satuan_isi` pada setiap item produk di response GET /api/monitoring. Field ini MUST diambil dari `produk.satuan_isi` dan disertakan di level produk (bukan level batch).

#### Scenario: Satuan isi tersedia di response monitoring
- **WHEN** client melakukan GET /api/monitoring
- **THEN** setiap item produk memiliki field `satuan_isi` berisi satuan produk (misal: "ml", "gram", "pcs")

#### Scenario: Satuan isi tersedia untuk produk partial use
- **WHEN** client melakukan GET /api/monitoring dengan filter partial use
- **THEN** setiap item produk partial use memiliki `satuan_isi` yang dapat digunakan frontend untuk menampilkan satuan pada kolom total isi dan sisa kemasan terbuka
