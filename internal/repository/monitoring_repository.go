package repository

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"freya-skin-clinic-backend/internal/model"
)

type MonitoringRepository interface {
	FindAllForMonitoring(ctx context.Context, filter model.MonitoringFilter) ([]model.MonitoringProdukItem, error)
}

type monitoringRepository struct {
	db *pgxpool.Pool
}

func NewMonitoringRepository(db *pgxpool.Pool) MonitoringRepository {
	return &monitoringRepository{db: db}
}

func (r *monitoringRepository) FindAllForMonitoring(ctx context.Context, filter model.MonitoringFilter) ([]model.MonitoringProdukItem, error) {
	// Ambil batch yang masih relevan untuk monitoring:
	// - status AKTIF dengan stok > 0, ATAU
	// - status AKTIF dengan kemasan terbuka masih ada isinya (Partial Use)
	// Tidak bergantung pada worker untuk update status — indikator dihitung real-time dari expired_date
	query := `
		SELECT b.id, b.id_produk, b.kode_batch, b.expired_date,
		       b.stok_kemasan, b.total_isi_tersedia, b.status,
		       p.id, p.kode_produk, p.nama_produk, p.id_kategori, p.pola_penggunaan, p.satuan_isi,
		       k.nama_kategori
		FROM batch_stok b
		JOIN produk p ON p.id = b.id_produk
		JOIN kategori k ON k.id = p.id_kategori
		WHERE b.status = 'AKTIF'
		  AND (
		    b.stok_kemasan > 0
		    OR EXISTS (
		      SELECT 1 FROM kemasan_terbuka kt
		      WHERE kt.id_batch = b.id AND kt.status_bud = 'AKTIF' AND kt.isi_tersisa > 0
		    )
		  )
	`
	args := []interface{}{}
	idx := 1

	if filter.KategoriID != "" {
		query += ` AND p.id_kategori = $` + itoa(idx)
		args = append(args, filter.KategoriID)
		idx++
	}
	if filter.NamaProduk != "" {
		query += ` AND LOWER(p.nama_produk) LIKE $` + itoa(idx)
		args = append(args, "%"+strings.ToLower(filter.NamaProduk)+"%")
		idx++
	}

	query += ` ORDER BY p.kode_produk, b.expired_date ASC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	produkMap := map[string]*model.MonitoringProdukItem{}
	produkOrder := []string{}

	for rows.Next() {
		var (
			bID, bIDProduk, bKodeBatch, bStatus            string
			bExpiredDate                                   time.Time
			bStokKemasan                                   int
			bTotalIsi                                      float64
			pID, pKode, pNama, pIDKategori, pPola, pSatuan string
			kNama                                          string
		)
		if err := rows.Scan(
			&bID, &bIDProduk, &bKodeBatch, &bExpiredDate,
			&bStokKemasan, &bTotalIsi, &bStatus,
			&pID, &pKode, &pNama, &pIDKategori, &pPola, &pSatuan,
			&kNama,
		); err != nil {
			return nil, err
		}

		// Hitung indikator
		now := time.Now().Truncate(24 * time.Hour)
		diff := bExpiredDate.Sub(now).Hours() / 24
		indikator := "AMAN"
		if diff < 0 {
			indikator = "KADALUWARSA"
		} else if diff <= 30 {
			indikator = "MENDEKATI"
		}

		batchItem := model.MonitoringBatchItem{
			IDBatch:          bID,
			KodeBatch:        bKodeBatch,
			ExpiredDate:      bExpiredDate.Format("2006-01-02"),
			StokKemasan:      bStokKemasan,
			TotalIsiTersedia: bTotalIsi,
			StatusBatch:      bStatus,
			IndikatorExpired: indikator,
		}

		if _, exists := produkMap[pID]; !exists {
			produkMap[pID] = &model.MonitoringProdukItem{
				IDProduk:       pID,
				KodeProduk:     pKode,
				NamaProduk:     pNama,
				IDKategori:     pIDKategori,
				NamaKategori:   kNama,
				PolaPenggunaan: pPola,
				SatuanIsi:      pSatuan,
				Batches:        []model.MonitoringBatchItem{},
			}
			produkOrder = append(produkOrder, pID)
		}

		produkMap[pID].Batches = append(produkMap[pID].Batches, batchItem)
		produkMap[pID].TotalStok += bStokKemasan
		produkMap[pID].TotalIsi += bTotalIsi
	}

	// 2. Fetch kemasan terbuka untuk setiap batch (jika PARTIAL_USE)
	today := time.Now().Truncate(24 * time.Hour)
	for _, pID := range produkOrder {
		p := produkMap[pID]
		if p.PolaPenggunaan != "PARTIAL_USE" {
			continue
		}
		for i, b := range p.Batches {
			// Fetch kemasan terbuka tanpa filter status_bud — hitung real-time dari BUD date
			ktQuery := `
				SELECT id, tanggal_dibuka, bud, isi_awal, isi_tersisa, status_bud
				FROM kemasan_terbuka
				WHERE id_batch = $1
				ORDER BY updated_at DESC
				LIMIT 1
			`
			var ktID, ktStatus string
			var ktTanggalDibuka, ktBUD time.Time
			var ktIsiAwal, ktIsi float64
			err := r.db.QueryRow(ctx, ktQuery, b.IDBatch).Scan(&ktID, &ktTanggalDibuka, &ktBUD, &ktIsiAwal, &ktIsi, &ktStatus)
			if err == nil {
				// Hitung status BUD real-time berdasarkan tanggal
				effectiveStatusBUD := ktStatus
				if ktBUD.Before(today) {
					effectiveStatusBUD = "KADALUWARSA"
				}

				// Filter status BUD jika ada — pakai effective status (real-time)
				if filter.StatusBUD != "" && filter.StatusBUD != effectiveStatusBUD {
					continue
				}

				p.Batches[i].KemasanTerbuka = &model.MonitoringKemasanTerbuka{
					ID:            ktID,
					TanggalDibuka: ktTanggalDibuka.Format("2006-01-02"),
					BUD:           ktBUD.Format("2006-01-02"),
					IsiAwal:       ktIsiAwal,
					IsiTersisa:    ktIsi,
					StatusBUD:     effectiveStatusBUD,
				}
			}
		}
	}

	// 3. Post-filter: hapus batch PARTIAL_USE yang tidak punya kemasan terbuka setelah filter StatusBUD
	for _, pID := range produkOrder {
		p := produkMap[pID]
		if p.PolaPenggunaan != "PARTIAL_USE" || filter.StatusBUD == "" {
			continue
		}
		filtered := p.Batches[:0]
		for _, b := range p.Batches {
			if b.KemasanTerbuka != nil {
				filtered = append(filtered, b)
			}
		}
		produkMap[pID].Batches = filtered
	}

	result := make([]model.MonitoringProdukItem, 0, len(produkOrder))
	for _, pID := range produkOrder {
		item := *produkMap[pID]

		// Jika filter StatusBUD aktif, FULL_USE tidak punya kemasan terbuka — skip
		if filter.StatusBUD != "" && item.PolaPenggunaan != "PARTIAL_USE" {
			continue
		}

		// Jika filter StatusBUD aktif dan semua batch sudah tidak punya kemasan terbuka — skip
		if filter.StatusBUD != "" && len(item.Batches) == 0 {
			continue
		}

		// Post-filter IndikatorExpired: filter batch yang tidak sesuai indikator
		if filter.IndikatorExpired != "" {
			filtered := item.Batches[:0]
			for _, b := range item.Batches {
				if b.IndikatorExpired == filter.IndikatorExpired {
					filtered = append(filtered, b)
				}
			}
			item.Batches = filtered
			if len(item.Batches) == 0 {
				continue
			}
		}
		result = append(result, item)
	}
	return result, nil
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
