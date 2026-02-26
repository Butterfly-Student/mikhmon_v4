package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// ReportUseCase handles report business logic
type ReportUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Client
	log         *zap.Logger
}

// NewReportUseCase creates a new report use case
func NewReportUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Client,
	log *zap.Logger,
) *ReportUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &ReportUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
		log:         log.Named("report-usecase"),
	}
}

// GetSalesReport retrieves sales report from MikroTik
func (uc *ReportUseCase) GetSalesReport(ctx context.Context, routerID uint, owner string) ([]*dto.SalesReport, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	reportCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	reports, err := uc.mikrotikSvc.GetSalesReports(reportCtx, router, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales report from MikroTik: %w", err)
	}

	return reports, nil
}

// GetSalesReportByDay retrieves sales report by day
func (uc *ReportUseCase) GetSalesReportByDay(ctx context.Context, routerID uint, day string) ([]*dto.SalesReport, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	reportCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	reports, err := uc.mikrotikSvc.GetSalesReportsByDay(reportCtx, router, day)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily sales report from MikroTik: %w", err)
	}

	return reports, nil
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
