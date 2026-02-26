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

// SystemUseCase handles system business logic
type SystemUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Client
	log         *zap.Logger
}

// NewSystemUseCase creates a new system use case
func NewSystemUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Client,
	log *zap.Logger,
) *SystemUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &SystemUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
		log:         log.Named("system-usecase"),
	}
}

// GetResource retrieves system resources
func (uc *SystemUseCase) GetResource(ctx context.Context, routerID uint) (*dto.SystemResource, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	resourceCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.GetSystemResource(resourceCtx, router)
}

// GetHealth retrieves system health
func (uc *SystemUseCase) GetHealth(ctx context.Context, routerID uint) (*dto.SystemHealth, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.GetSystemHealth(healthCtx, router)
}

// GetIdentity retrieves system identity
func (uc *SystemUseCase) GetIdentity(ctx context.Context, routerID uint) (*dto.SystemIdentity, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	identityCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.GetSystemIdentity(identityCtx, router)
}

// GetRouterBoardInfo retrieves routerboard information
func (uc *SystemUseCase) GetRouterBoardInfo(ctx context.Context, routerID uint) (*dto.RouterBoardInfo, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	rbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.GetRouterBoardInfo(rbCtx, router)
}

// GetClock retrieves system clock
func (uc *SystemUseCase) GetClock(ctx context.Context, routerID uint) (*dto.SystemClock, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	clockCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.GetSystemClock(clockCtx, router)
}

// GetDashboardData retrieves complete dashboard data
func (uc *SystemUseCase) GetDashboardData(ctx context.Context, routerID uint) (*dto.DashboardData, error) {
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

	identityCtx, identityCancel := context.WithTimeout(ctx, 10*time.Second)
	identity, err := uc.mikrotikSvc.GetSystemIdentity(identityCtx, router)
	identityCancel()
	if err == nil {
		dashboard.Identity = identity
	} else {
		uc.log.Warn("GetSystemIdentity failed", zap.Uint("routerID", routerID), zap.Error(err))
	}

	rbCtx, rbCancel := context.WithTimeout(ctx, 10*time.Second)
	routerboard, err := uc.mikrotikSvc.GetRouterBoardInfo(rbCtx, router)
	rbCancel()
	if err == nil {
		dashboard.RouterBoard = routerboard
	} else {
		uc.log.Warn("GetRouterBoardInfo failed", zap.Uint("routerID", routerID), zap.Error(err))
	}

	healthCtx, healthCancel := context.WithTimeout(ctx, 10*time.Second)
	health, err := uc.mikrotikSvc.GetSystemHealth(healthCtx, router)
	healthCancel()
	if err == nil {
		dashboard.Health = health
	} else {
		uc.log.Warn("GetSystemHealth failed", zap.Uint("routerID", routerID), zap.Error(err))
	}

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
