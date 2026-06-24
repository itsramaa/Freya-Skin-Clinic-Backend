## MODIFIED Requirements

### Requirement: Response laporan stok keluar menyertakan satuan isi
Sistem SHALL menyertakan field `satuan_isi` pada setiap item response GET /api/laporan/stok-keluar. Field ini MUST diambil dari `produk.satuan_isi`.

#### Scenario: Satuan isi tersedia di laporan stok keluar
- **WHEN** client melakukan GET /api/laporan/stok-keluar dengan filter periode valid
- **THEN** setiap item response memiliki field `satuan_isi` berisi satuan produk

#### Scenario: Satuan isi digunakan untuk format tampilan partial use
- **WHEN** item laporan stok keluar memiliki `pola_penggunaan = "PARTIAL_USE"`
- **THEN** field `satuan_isi` tersedia sehingga frontend dapat menampilkan "X ml" / "X gram" pada kolom jumlah dipakai

### Requirement: Response laporan sisa stok menyertakan satuan isi
Sistem SHALL menyertakan field `satuan_isi` pada setiap item response GET /api/laporan/sisa-stok. Field ini MUST diambil dari `produk.satuan_isi`.

#### Scenario: Satuan isi tersedia di laporan sisa stok
- **WHEN** client melakukan GET /api/laporan/sisa-stok
- **THEN** setiap item response memiliki field `satuan_isi` berisi satuan produk

#### Scenario: Satuan isi digunakan untuk kolom total isi partial use
- **WHEN** item sisa stok memiliki `pola_penggunaan = "PARTIAL_USE"`
- **THEN** field `satuan_isi` tersedia sehingga frontend dapat menampilkan "X ml" / "X gram" pada kolom total isi tersedia

#### Scenario: Total isi kosong untuk produk full use
- **WHEN** item sisa stok memiliki `pola_penggunaan = "FULL_USE"`
- **THEN** frontend menampilkan "—" pada kolom total isi tersedia (field tetap ada di response, hanya tidak ditampilkan)
