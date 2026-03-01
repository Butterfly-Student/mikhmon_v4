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

// LogUseCase handles log business logic
type LogUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Manager
	log         *zap.Logger
}

// NewLogUseCase creates a new log use case
func NewLogUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	log *zap.Logger,
) *LogUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &LogUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
		log:         log.Named("log-usecase"),
	}
}

// GetLogs retrieves a snapshot of log entries. Pass topics="" for all logs.
func (uc *LogUseCase) GetLogs(ctx context.Context, routerID uint, topics string, limit int) ([]*dto.LogEntry, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	c, err := connectRouter(ctx, uc.mikrotikSvc, router)
	if err != nil {
		return nil, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}

	logCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetLogs(logCtx, topics, limit)
}

// GetPPPLogs retrieves PPP log entries.
func (uc *LogUseCase) GetPPPLogs(ctx context.Context, routerID uint, limit int) ([]*dto.LogEntry, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	c, err := connectRouter(ctx, uc.mikrotikSvc, router)
	if err != nil {
		return nil, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}

	logCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetPPPLogs(logCtx, limit)
}

// GetHotspotLogs retrieves hotspot logs
func (uc *LogUseCase) GetHotspotLogs(ctx context.Context, routerID uint, limit int) ([]*dto.LogEntry, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	c, err := connectRouter(ctx, uc.mikrotikSvc, router)
	if err != nil {
		return nil, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}

	logCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetHotspotLogs(logCtx, limit)
}
