package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// DashboardUseCase handles dashboard business logic
type DashboardUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Client
	log         *zap.Logger
}

// NewDashboardUseCase creates a new dashboard use case
func NewDashboardUseCase(routerRepo repository.RouterRepository, mikrotikSvc *mikrotik.Client, log *zap.Logger) *DashboardUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &DashboardUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
		log:         log.Named("dashboard-usecase"),
	}
}

// GetDashboardData retrieves complete dashboard data.
// Timeout dinaikkan ke 10 detik untuk menangani router yang lambat / jaringan variabel.
func (uc *DashboardUseCase) GetDashboardData(ctx context.Context, routerID uint) (*dto.DashboardData, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("router not found: %w", err)
	}

	dashboard := &dto.DashboardData{
		RouterID:   routerID,
		RouterName: router.Name,
	}

	uc.log.Info("Fetching dashboard data",
		zap.Uint("routerID", routerID),
		zap.String("host", router.Host),
	)

	// Get system resource — timeout 10 detik
	resourceCtx, resourceCancel := context.WithTimeout(ctx, 10*time.Second)
	defer resourceCancel()

	resource, err := uc.mikrotikSvc.GetSystemResource(resourceCtx, router)
	if err != nil {
		uc.log.Warn("GetSystemResource failed, returning partial data",
			zap.Uint("routerID", routerID),
			zap.Error(err),
		)
		dashboard.ConnectionError = fmt.Sprintf("Connection error: %v", err)
		return dashboard, nil
	}
	dashboard.Resource = resource

	// Get system identity — separate 10s timeout
	identityCtx, identityCancel := context.WithTimeout(ctx, 10*time.Second)
	identity, err := uc.mikrotikSvc.GetSystemIdentity(identityCtx, router)
	identityCancel()
	if err == nil {
		dashboard.Identity = identity
	} else {
		uc.log.Warn("GetSystemIdentity failed", zap.Uint("routerID", routerID), zap.Error(err))
	}

	// Get routerboard info — separate 10s timeout
	rbCtx, rbCancel := context.WithTimeout(ctx, 10*time.Second)
	routerboard, err := uc.mikrotikSvc.GetRouterBoardInfo(rbCtx, router)
	rbCancel()
	if err == nil {
		dashboard.RouterBoard = routerboard
	} else {
		uc.log.Warn("GetRouterBoardInfo failed", zap.Uint("routerID", routerID), zap.Error(err))
	}

	// Get hotspot stats — separate 10s timeout
	statsCtx, statsCancel := context.WithTimeout(ctx, 10*time.Second)
	defer statsCancel()

	stats := &dto.HotspotStats{}
	if totalUsers, err := uc.mikrotikSvc.GetHotspotUsersCount(statsCtx, router); err == nil {
		stats.TotalUsers = totalUsers
	} else {
		uc.log.Warn("GetHotspotUsersCount failed", zap.Uint("routerID", routerID), zap.Error(err))
	}
	if activeUsers, err := uc.mikrotikSvc.GetHotspotActiveCount(statsCtx, router); err == nil {
		stats.ActiveUsers = activeUsers
	} else {
		uc.log.Warn("GetHotspotActiveCount failed", zap.Uint("routerID", routerID), zap.Error(err))
	}
	dashboard.Stats = stats

	uc.log.Info("Dashboard data fetched successfully",
		zap.Uint("routerID", routerID),
		zap.Int("activeUsers", stats.ActiveUsers),
		zap.Int("totalUsers", stats.TotalUsers),
	)

	return dashboard, nil
}

// GetResource retrieves system resources with a 10s timeout.
func (uc *DashboardUseCase) GetResource(ctx context.Context, routerID uint, force bool) (*dto.SystemResource, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	resourceCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.GetSystemResource(resourceCtx, router)
}

// GetTraffic retrieves traffic stats with a 10s timeout.
func (uc *DashboardUseCase) GetTraffic(ctx context.Context, routerID uint, iface string) (*dto.TrafficStats, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	trafficCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.MonitorTraffic(trafficCtx, router, iface)
}

// GetInterfaces retrieves network interfaces with a 10s timeout.
func (uc *DashboardUseCase) GetInterfaces(ctx context.Context, routerID uint) ([]*dto.Interface, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	interfacesCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.GetInterfaces(interfacesCtx, router)
}
