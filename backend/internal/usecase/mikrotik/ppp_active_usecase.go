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

// PPPActiveUseCase handles PPP active session business logic
type PPPActiveUseCase struct {
	routerRepo repository.RouterRepository
	mgr        *mikrotik.Manager
	log        *zap.Logger
}

// NewPPPActiveUseCase creates a new PPP active use case
func NewPPPActiveUseCase(
	routerRepo repository.RouterRepository,
	mgr *mikrotik.Manager,
	log *zap.Logger,
) *PPPActiveUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &PPPActiveUseCase{
		routerRepo: routerRepo,
		mgr:        mgr,
		log:        log.Named("ppp-active-usecase"),
	}
}

func (uc *PPPActiveUseCase) connect(ctx context.Context, routerID uint) (*mikrotik.Client, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}
	c, err := connectRouter(ctx, uc.mgr, router)
	if err != nil {
		return nil, fmt.Errorf("router %q not connected: %w", router.Name, err)
	}
	return c, nil
}

// GetActive retrieves active PPP sessions, optionally filtered by service.
func (uc *PPPActiveUseCase) GetActive(ctx context.Context, routerID uint, service string) ([]*dto.PPPActive, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetPPPActive(tctx, service)
}

// GetActiveByID retrieves an active PPP session by ID.
func (uc *PPPActiveUseCase) GetActiveByID(ctx context.Context, routerID uint, id string) (*dto.PPPActive, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetPPPActiveByID(tctx, id)
}

// RemoveActive disconnects an active PPP session.
func (uc *PPPActiveUseCase) RemoveActive(ctx context.Context, routerID uint, id string) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.RemovePPPActive(tctx, id)
}

// ListenActive starts a streaming subscription to active PPP sessions.
// No timeout is applied — the caller controls lifetime via ctx.
func (uc *PPPActiveUseCase) ListenActive(
	ctx context.Context,
	routerID uint,
	resultChan chan<- []*dto.PPPActive,
) (func() error, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	return c.ListenPPPActive(ctx, resultChan)
}

// ListenInactive starts a streaming subscription to inactive PPP secrets.
// No timeout is applied — the caller controls lifetime via ctx.
func (uc *PPPActiveUseCase) ListenInactive(
	ctx context.Context,
	routerID uint,
	resultChan chan<- []*dto.PPPSecret,
) (func() error, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	return c.ListenPPPInactive(ctx, resultChan)
}
