package mikrotik

import (
	"context"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// NATUseCase handles NAT business logic
type NATUseCase struct {
	routerRepo  repository.RouterRepository
	mikrotikSvc *mikrotik.Client
	log         *zap.Logger
}

// NewNATUseCase creates a new NAT use case
func NewNATUseCase(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Client,
	log *zap.Logger,
) *NATUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &NATUseCase{
		routerRepo:  routerRepo,
		mikrotikSvc: mikrotikSvc,
		log:         log.Named("nat-usecase"),
	}
}

// GetNATRules retrieves NAT rules
func (uc *NATUseCase) GetNATRules(ctx context.Context, routerID uint) ([]*dto.NATRule, error) {
	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	natCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return uc.mikrotikSvc.GetNATRules(natCtx, router)
}
