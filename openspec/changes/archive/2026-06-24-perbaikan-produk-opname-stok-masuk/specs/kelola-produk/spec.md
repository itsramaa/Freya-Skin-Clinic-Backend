## MODIFIED Requirements

### Requirement: Tampilkan informasi isi per kemasan produk
Sistem SHALL menampilkan kolom isi per kemasan pada halaman kelola produk. Untuk produk FULL_USE, sistem MUST menampilkan teks "per pcs" tanpa kalkulasi apapun. Untuk produk PARTIAL_USE, sistem MUST menampilkan nilai `isi_per_kemasan` beserta satuannya.

#### Scenario: Produk Full Use menampilkan per pcs
- **WHEN** halaman kelola produk ditampilkan dan produk memiliki `pola_penggunaan = "FULL_USE"`
- **THEN** kolom isi per kemasan menampilkan teks "per pcs" tanpa nilai numerik

#### Scenario: Produk Partial Use menampilkan nilai isi per kemasan
- **WHEN** halaman kelola produk ditampilkan dan produk memiliki `pola_penggunaan = "PARTIAL_USE"`
- **THEN** kolom isi per kemasan menampilkan nilai `isi_per_kemasan` beserta `satuan_isi` (contoh: "20 ml")
