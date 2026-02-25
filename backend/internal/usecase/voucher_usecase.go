package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/cache"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
)

// VoucherUseCase handles voucher business logic
type VoucherUseCase struct {
	routerRepo  repository.RouterRepository
	hotspotSvc  *mikrotik.HotspotService
	cache       cache.Cache
}

// NewVoucherUseCase creates a new voucher use case
func NewVoucherUseCase(
	routerRepo repository.RouterRepository,
	hotspotSvc *mikrotik.HotspotService,
	cache cache.Cache,
) *VoucherUseCase {
	return &VoucherUseCase{
		routerRepo:  routerRepo,
		hotspotSvc:  hotspotSvc,
		cache:       cache,
	}
}

// GenerateVouchers generates vouchers and adds them directly to MikroTik
func (uc *VoucherUseCase) GenerateVouchers(ctx context.Context, routerID uint, req dto.VoucherGenerateRequest) (*dto.VoucherBatchResult, error) {
	// Verify router exists
	_, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	// Generate voucher codes
	generator := mikrotik.NewVoucherGenerator()
	vouchers := generator.GenerateBatch(&req)

	// Add each voucher to MikroTik
	var createdVouchers []dto.Voucher
	comment := fmt.Sprintf("vc-%s-%s", time.Now().Format("01.02.06"), req.Comment)
	
	for _, voucher := range vouchers {
		userReq := dto.CreateUserRequest{
			Name:       voucher.Username,
			Password:   voucher.Password,
			Profile:    voucher.Profile,
			Server:     voucher.Server,
			Comment:    comment,
		}

		user, err := uc.hotspotSvc.AddUser(ctx, routerID, userReq)
		if err != nil {
			// Log error but continue with next voucher
			continue
		}
		
		createdVouchers = append(createdVouchers, dto.Voucher{
			Username: user.Name,
			Password: user.Password,
			Profile:  user.Profile,
		})
	}

	// Invalidate cache
	uc.cache.Delete(ctx, fmt.Sprintf("users:%d:*", routerID))

	return &dto.VoucherBatchResult{
		Count:    len(createdVouchers),
		Comment:  comment,
		Vouchers: createdVouchers,
	}, nil
}

// GetVouchersByComment retrieves vouchers by comment (for printing)
func (uc *VoucherUseCase) GetVouchersByComment(ctx context.Context, routerID uint, comment string) ([]dto.HotspotUser, error) {
	// Verify router exists
	_, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	return uc.hotspotSvc.GetUsers(ctx, routerID, dto.UserFilter{
		Comment: comment,
	})
}

// DeleteVouchersByComment deletes vouchers by comment
func (uc *VoucherUseCase) DeleteVouchersByComment(ctx context.Context, routerID uint, comment string) error {
	// Verify router exists
	_, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}

	users, err := uc.hotspotSvc.GetUsers(ctx, routerID, dto.UserFilter{
		Comment: comment,
	})
	if err != nil {
		return err
	}

	for _, user := range users {
		if err := uc.hotspotSvc.RemoveUser(ctx, routerID, user.ID); err != nil {
			continue
		}
	}

	// Invalidate cache
	uc.cache.Delete(ctx, fmt.Sprintf("users:%d:*", routerID))

	return nil
}
