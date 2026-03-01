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

// PPPProfileUseCase handles PPP profile business logic
type PPPProfileUseCase struct {
	routerRepo repository.RouterRepository
	mgr        *mikrotik.Manager
	log        *zap.Logger
}

// NewPPPProfileUseCase creates a new PPP profile use case
func NewPPPProfileUseCase(
	routerRepo repository.RouterRepository,
	mgr *mikrotik.Manager,
	log *zap.Logger,
) *PPPProfileUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &PPPProfileUseCase{
		routerRepo: routerRepo,
		mgr:        mgr,
		log:        log.Named("ppp-profile-usecase"),
	}
}

func (uc *PPPProfileUseCase) connect(ctx context.Context, routerID uint) (*mikrotik.Client, error) {
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

// GetProfiles retrieves all PPP profiles.
func (uc *PPPProfileUseCase) GetProfiles(ctx context.Context, routerID uint) ([]*dto.PPPProfile, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetPPPProfiles(tctx)
}

// GetProfileByID retrieves a PPP profile by ID.
func (uc *PPPProfileUseCase) GetProfileByID(ctx context.Context, routerID uint, id string) (*dto.PPPProfile, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetPPPProfileByID(tctx, id)
}

// GetProfileByName retrieves a PPP profile by name.
func (uc *PPPProfileUseCase) GetProfileByName(ctx context.Context, routerID uint, name string) (*dto.PPPProfile, error) {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.GetPPPProfileByName(tctx, name)
}

// AddProfile adds a new PPP profile.
func (uc *PPPProfileUseCase) AddProfile(ctx context.Context, routerID uint, req *dto.PPPProfileRequest) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	profile := &dto.PPPProfile{
		Name:              req.Name,
		LocalAddress:      req.LocalAddress,
		RemoteAddress:     req.RemoteAddress,
		DNSServer:         req.DNSServer,
		SessionTimeout:    req.SessionTimeout,
		IdleTimeout:       req.IdleTimeout,
		OnlyOne:           req.OnlyOne,
		Comment:           req.Comment,
		RateLimit:         req.RateLimit,
		ParentQueue:       req.ParentQueue,
		QueueType:         req.QueueType,
		UseCompression:    req.UseCompression,
		UseEncryption:     req.UseEncryption,
		UseMPLS:           req.UseMPLS,
		UseUPnP:           req.UseUPnP,
		Bridge:            req.Bridge,
		AddressList:       req.AddressList,
		InterfaceList:     req.InterfaceList,
		OnUp:              req.OnUp,
		OnDown:            req.OnDown,
		ChangeTCPMSS:      req.ChangeTCPMSS,
		IncomingFilter:    req.IncomingFilter,
		OutgoingFilter:    req.OutgoingFilter,
		InsertQueueBefore: req.InsertQueueBefore,
		WinsServer:        req.WinsServer,
	}

	return c.AddPPPProfile(tctx, profile)
}

// UpdateProfile updates an existing PPP profile.
func (uc *PPPProfileUseCase) UpdateProfile(ctx context.Context, routerID uint, id string, req *dto.PPPProfileRequest) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	profile := &dto.PPPProfile{
		Name:              req.Name,
		LocalAddress:      req.LocalAddress,
		RemoteAddress:     req.RemoteAddress,
		DNSServer:         req.DNSServer,
		SessionTimeout:    req.SessionTimeout,
		IdleTimeout:       req.IdleTimeout,
		OnlyOne:           req.OnlyOne,
		Comment:           req.Comment,
		RateLimit:         req.RateLimit,
		ParentQueue:       req.ParentQueue,
		QueueType:         req.QueueType,
		UseCompression:    req.UseCompression,
		UseEncryption:     req.UseEncryption,
		UseMPLS:           req.UseMPLS,
		UseUPnP:           req.UseUPnP,
		Bridge:            req.Bridge,
		AddressList:       req.AddressList,
		InterfaceList:     req.InterfaceList,
		OnUp:              req.OnUp,
		OnDown:            req.OnDown,
		ChangeTCPMSS:      req.ChangeTCPMSS,
		IncomingFilter:    req.IncomingFilter,
		OutgoingFilter:    req.OutgoingFilter,
		InsertQueueBefore: req.InsertQueueBefore,
		WinsServer:        req.WinsServer,
	}

	return c.UpdatePPPProfile(tctx, id, profile)
}

// RemoveProfile removes a PPP profile.
func (uc *PPPProfileUseCase) RemoveProfile(ctx context.Context, routerID uint, id string) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.RemovePPPProfile(tctx, id)
}

// DisableProfile disables a PPP profile.
func (uc *PPPProfileUseCase) DisableProfile(ctx context.Context, routerID uint, id string) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.DisablePPPProfile(tctx, id)
}

// EnableProfile enables a PPP profile.
func (uc *PPPProfileUseCase) EnableProfile(ctx context.Context, routerID uint, id string) error {
	c, err := uc.connect(ctx, routerID)
	if err != nil {
		return err
	}

	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return c.EnablePPPProfile(tctx, id)
}
