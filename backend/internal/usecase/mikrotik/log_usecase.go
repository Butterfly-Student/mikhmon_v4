package mikrotik

import (
	"context"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// LogUseCase handles log business logic
type LogUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Client
	log         *zap.Logger
}

// NewLogUseCase creates a new log use case
func NewLogUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Client,
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

// GetHotspotLogs retrieves hotspot logs
func (uc *LogUseCase) GetHotspotLogs(ctx context.Context, routerID uint, limit int) ([]*dto.LogEntry, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	logCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.GetHotspotLogs(logCtx, router, limit)
}
