package services

import (
	"go-kasir-api/models"
	"go-kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetTodaySummary() (*models.SalesSummary, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	return s.repo.GetSalesSummary(startOfDay, endOfDay)
}

func (s *ReportService) GetSummaryByDateRange(startDate, endDate time.Time) (*models.SalesSummary, error) {
	endDateExclusive := endDate.AddDate(0, 0, 1)
	return s.repo.GetSalesSummary(startDate, endDateExclusive)
}