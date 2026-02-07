package models

type BestSellingProduct struct {
	Name string `json:"name"`
	SellingQuantity int `json:"selling_quantity"`
}

type SalesSummary struct {
	TotalRevenue int `json:"total_revenue"`
	TotalTransaction int `json:"total_transaction"`
	BestSellingProduct *BestSellingProduct `json:"best_selling_product"`
}