package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// QueueUseCase handles queue business logic
type QueueUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Manager
	log         *zap.Logger
}

// NewQueueUseCase creates a new queue use case
func NewQueueUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	log *zap.Logger,
) *QueueUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &QueueUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
		log:         log.Named("queue-usecase"),
	}
}

// GetAllQueues retrieves all queues
func (uc *QueueUseCase) GetAllQueues(ctx context.Context, routerID uint) ([]string, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	c, err := connectRouter(ctx, uc.mikrotikSvc, router)
	if err != nil {
		return nil, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}

	queueCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	queues, err := c.GetAllQueues(queueCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get queues from MikroTik: %w", err)
	}
	return queues, nil
}

// GetParentQueues retrieves parent queues
func (uc *QueueUseCase) GetParentQueues(ctx context.Context, routerID uint) ([]string, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	c, err := connectRouter(ctx, uc.mikrotikSvc, router)
	if err != nil {
		return nil, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}

	queueCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	queues, err := c.GetAllParentQueues(queueCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent queues from MikroTik: %w", err)
	}
	return queues, nil
}
