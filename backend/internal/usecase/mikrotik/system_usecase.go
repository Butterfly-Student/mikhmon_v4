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
	mikrotikSvc *mikrotik.Manager
	log         *zap.Logger
}

// NewSystemUseCase creates a new system use case
func NewSystemUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
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

// client resolves the *mikrotik.Client for the given routerID, connecting on demand.
func (uc *SystemUseCase) client(ctx context.Context, routerID uint) (*mikrotik.Client, string, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, "", err
	}
	cfg := mikrotik.Config{
		Host:     router.Host,
		Port:     router.Port,
		Username: router.Username,
		Password: router.Password,
		UseTLS:   router.UseSSL,
		Timeout:  time.Duration(router.Timeout) * time.Second,
	}
	c, err := uc.mikrotikSvc.GetOrConnect(ctx, router.Name, cfg)
	if err != nil {
		return nil, router.Name, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}
	return c, router.Name, nil
}

// GetResource retrieves system resources
func (uc *SystemUseCase) GetResource(ctx context.Context, routerID uint) (*dto.SystemResource, error) {
	c, _, err := uc.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	resourceCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetSystemResource(resourceCtx)
}

// GetHealth retrieves system health
func (uc *SystemUseCase) GetHealth(ctx context.Context, routerID uint) (*dto.SystemHealth, error) {
	c, _, err := uc.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetSystemHealth(healthCtx)
}

// GetIdentity retrieves system identity
func (uc *SystemUseCase) GetIdentity(ctx context.Context, routerID uint) (*dto.SystemIdentity, error) {
	c, _, err := uc.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	identityCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetSystemIdentity(identityCtx)
}

// GetRouterBoardInfo retrieves routerboard information
func (uc *SystemUseCase) GetRouterBoardInfo(ctx context.Context, routerID uint) (*dto.RouterBoardInfo, error) {
	c, _, err := uc.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	rbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetRouterBoardInfo(rbCtx)
}

// GetClock retrieves system clock
func (uc *SystemUseCase) GetClock(ctx context.Context, routerID uint) (*dto.SystemClock, error) {
	c, _, err := uc.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	clockCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetSystemClock(clockCtx)
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

	cfg := mikrotik.Config{
		Host:     router.Host,
		Port:     router.Port,
		Username: router.Username,
		Password: router.Password,
		UseTLS:   router.UseSSL,
		Timeout:  time.Duration(router.Timeout) * time.Second,
	}
	c, err := uc.mikrotikSvc.GetOrConnect(ctx, router.Name, cfg)
	if err != nil {
		dashboard.ConnectionError = fmt.Sprintf("Router not connected: %v", err)
		return dashboard, nil
	}

	uc.log.Info("Fetching dashboard data",
		zap.Uint("routerID", routerID),
		zap.String("router", router.Name),
	)

	identityCtx, identityCancel := context.WithTimeout(ctx, 10*time.Second)
	identity, err := c.GetSystemIdentity(identityCtx)
	identityCancel()
	if err == nil {
		dashboard.Identity = identity
	} else {
		uc.log.Warn("GetSystemIdentity failed", zap.Uint("routerID", routerID), zap.Error(err))
	}

	rbCtx, rbCancel := context.WithTimeout(ctx, 10*time.Second)
	routerboard, err := c.GetRouterBoardInfo(rbCtx)
	rbCancel()
	if err == nil {
		dashboard.RouterBoard = routerboard
	} else {
		uc.log.Warn("GetRouterBoardInfo failed", zap.Uint("routerID", routerID), zap.Error(err))
	}

	healthCtx, healthCancel := context.WithTimeout(ctx, 10*time.Second)
	health, err := c.GetSystemHealth(healthCtx)
	healthCancel()
	if err == nil {
		dashboard.Health = health
	} else {
		uc.log.Warn("GetSystemHealth failed", zap.Uint("routerID", routerID), zap.Error(err))
	}

	statsCtx, statsCancel := context.WithTimeout(ctx, 10*time.Second)
	defer statsCancel()

	stats := &dto.HotspotStats{}
	if totalUsers, err := c.GetHotspotUsersCount(statsCtx); err == nil {
		stats.TotalUsers = totalUsers
	} else {
		uc.log.Warn("GetHotspotUsersCount failed", zap.Uint("routerID", routerID), zap.Error(err))
	}
	if activeUsers, err := c.GetHotspotActiveCount(statsCtx); err == nil {
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
