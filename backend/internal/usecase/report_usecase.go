package usecase

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
)

// ReportUseCase handles report business logic
type ReportUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Client
}

// NewReportUseCase creates a new report use case
func NewReportUseCase(routerRepo repository.RouterRepository, mikrotikSvc *mikrotik.Client) *ReportUseCase {
	return &ReportUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
	}
}

// GetSalesReport retrieves sales report from MikroTik
func (uc *ReportUseCase) GetSalesReport(ctx context.Context, routerID uint, owner string, force bool) ([]*dto.SalesReport, error) {
	// Placeholder implementation
	return []*dto.SalesReport{}, nil
}

// GetSalesReportByDay retrieves sales report by day
func (uc *ReportUseCase) GetSalesReportByDay(ctx context.Context, routerID uint, day string, force bool) ([]*dto.SalesReport, error) {
	// Placeholder implementation
	return []*dto.SalesReport{}, nil
}

// CalculateSummary calculates summary from reports
func (uc *ReportUseCase) CalculateSummary(reports []*dto.SalesReport) *dto.ReportSummary {
	if len(reports) == 0 {
		return &dto.ReportSummary{
			TotalVouchers: 0,
			TotalAmount:   0,
			ByProfile:     make(map[string]dto.ProfileSummary),
		}
	}

	var totalAmount float64
	byProfile := make(map[string]dto.ProfileSummary)

	for _, r := range reports {
		totalAmount += r.Price
		ps := byProfile[r.Profile]
		ps.Count++
		ps.Total += r.Price
		byProfile[r.Profile] = ps
	}

	return &dto.ReportSummary{
		TotalVouchers: len(reports),
		TotalAmount:   totalAmount,
		ByProfile:     byProfile,
	}
}
