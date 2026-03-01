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

// PPPSecretUseCase handles PPP secret business logic
type PPPSecretUseCase struct {
	routerRepo repository.RouterRepository
	mgr        *mikrotik.Manager
	log        *zap.Logger
}

// NewPPPSecretUseCase creates a new PPP secret use case
func NewPPPSecretUseCase(
	routerRepo repository.RouterRepository,
	mgr *mikrotik.Manager,
	log *zap.Logger,
) *PPPSecretUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &PPPSecretUseCase{
		routerRepo: routerRepo,
		mgr:        mgr,
		log:        log.Named("ppp-secret-usecase"),
	}
}

func (uc *PPPSecretUseCase) connect(ctx context.Context, routerID uint) (*mikrotik.Client, error) {
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

// GetSecrets retrieves PPP secrets, optionally filtered by profile.
func (uc *PPPSecretUseCase) GetSecrets(ctx context.Context, routerID uint, profile string) ([]*dto.PPPSecret, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetPPPSecrets(tctx, profile)
}

// GetSecretByID retrieves a PPP secret by ID.
func (uc *PPPSecretUseCase) GetSecretByID(ctx context.Context, routerID uint, id string) (*dto.PPPSecret, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetPPPSecretByID(tctx, id)
}

// GetSecretByName retrieves a PPP secret by name.
func (uc *PPPSecretUseCase) GetSecretByName(ctx context.Context, routerID uint, name string) (*dto.PPPSecret, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetPPPSecretByName(tctx, name)
}

// AddSecret adds a new PPP secret.
func (uc *PPPSecretUseCase) AddSecret(ctx context.Context, routerID uint, req *dto.PPPSecretRequest) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	secret := &dto.PPPSecret{
		Name:          req.Name,
		Password:      req.Password,
		Profile:       req.Profile,
		Service:       req.Service,
		CallerID:      req.CallerID,
		LocalAddress:  req.LocalAddress,
		RemoteAddress: req.RemoteAddress,
		Routes:        req.Routes,
		Comment:       req.Comment,
		LimitBytesIn:  req.LimitBytesIn,
		LimitBytesOut: req.LimitBytesOut,
	}

	return c.AddPPPSecret(tctx, secret)
}

// UpdateSecret updates an existing PPP secret.
func (uc *PPPSecretUseCase) UpdateSecret(ctx context.Context, routerID uint, id string, req *dto.PPPSecretUpdateRequest) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	secret := &dto.PPPSecret{
		Name:          req.Name,
		Password:      req.Password,
		Profile:       req.Profile,
		Service:       req.Service,
		Disabled:      req.Disabled,
		CallerID:      req.CallerID,
		LocalAddress:  req.LocalAddress,
		RemoteAddress: req.RemoteAddress,
		Routes:        req.Routes,
		Comment:       req.Comment,
		LimitBytesIn:  req.LimitBytesIn,
		LimitBytesOut: req.LimitBytesOut,
	}

	return c.UpdatePPPSecret(tctx, id, secret)
}

// RemoveSecret removes a PPP secret.
func (uc *PPPSecretUseCase) RemoveSecret(ctx context.Context, routerID uint, id string) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.RemovePPPSecret(tctx, id)
}

// DisableSecret disables a PPP secret.
func (uc *PPPSecretUseCase) DisableSecret(ctx context.Context, routerID uint, id string) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.DisablePPPSecret(tctx, id)
}

// EnableSecret enables a PPP secret.
func (uc *PPPSecretUseCase) EnableSecret(ctx context.Context, routerID uint, id string) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.EnablePPPSecret(tctx, id)
}
