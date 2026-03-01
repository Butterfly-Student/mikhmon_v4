package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// PoolUseCase handles pool business logic
type PoolUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Manager
	log         *zap.Logger
}

// NewPoolUseCase creates a new pool use case
func NewPoolUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	log *zap.Logger,
) *PoolUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &PoolUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
		log:         log.Named("pool-usecase"),
	}
}

// GetAddressPools retrieves address pools
func (uc *PoolUseCase) GetAddressPools(ctx context.Context, routerID uint) ([]string, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	c, err := connectRouter(ctx, uc.mikrotikSvc, router)
	if err != nil {
		return nil, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}

	poolCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pools, err := c.GetAddressPools(poolCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get address pools from MikroTik: %w", err)
	}
	return pools, nil
}
