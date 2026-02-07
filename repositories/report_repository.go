package repositories

import (
	"database/sql"
	"go-kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetSalesSummary(startDate, endDate time.Time) (*models.SalesSummary, error) {
	var totalRevenue sql.NullInt64
	var totalTransaction int

	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`, startDate, endDate).Scan(&totalRevenue, &totalTransaction)
	if err != nil {
		return nil, err
	}

	var bestSellingProduct *models.BestSellingProduct
	var productName sql.NullString
	var sellingQuantity sql.NullInt64

	err = repo.db.QueryRow(`
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&productName, &sellingQuantity)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if productName.Valid && sellingQuantity.Valid {
		bestSellingProduct = &models.BestSellingProduct{
			Name:           productName.String,
			SellingQuantity: int(sellingQuantity.Int64),
		}
	}

	return &models.SalesSummary{
		TotalRevenue:   int(totalRevenue.Int64),
		TotalTransaction: totalTransaction,
		BestSellingProduct: bestSellingProduct,
	}, nil
}
