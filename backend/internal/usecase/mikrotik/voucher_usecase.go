package mikrotik

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/cache"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"go.uber.org/zap"
)

// VoucherUseCase handles voucher business logic
type VoucherUseCase struct {
	routerRepo repository.RouterRepository
	hotspotSvc *mikrotik.HotspotService
	cache      cache.Cache
	log        *zap.Logger
}

// NewVoucherUseCase creates a new voucher use case
func NewVoucherUseCase(
	routerRepo repository.RouterRepository,
	hotspotSvc *mikrotik.HotspotService,
	cache cache.Cache,
	log *zap.Logger,
) *VoucherUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &VoucherUseCase{
		routerRepo: routerRepo,
		hotspotSvc: hotspotSvc,
		cache:      cache,
		log:        log.Named("voucher-usecase"),
	}
}

// GenerateVouchers generates vouchers and adds them directly to MikroTik
func (uc *VoucherUseCase) GenerateVouchers(ctx context.Context, routerID uint, req dto.VoucherGenerateRequest) (*dto.VoucherBatchResult, error) {
	_, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}
	if err := validateVoucherGenerateRequest(req); err != nil {
		return nil, err
	}

	generator := mikrotik.NewVoucherGenerator()
	vouchers := generator.GenerateBatch(&req)

	var createdVouchers []dto.Voucher
	gencode := normalizeGencode(req.Gencode)
	comment := fmt.Sprintf("%s-%s-%s-%s", req.Mode, gencode, time.Now().Format("01.02.06"), strings.TrimSpace(req.Comment))
	dataLimitBytes := mikrotik.ParseDataLimit(req.DataLimit)

	for _, voucher := range vouchers {
		userReq := dto.CreateUserRequest{
			Name:            voucher.Username,
			Password:        voucher.Password,
			Profile:         voucher.Profile,
			Server:          voucher.Server,
			LimitUptime:     req.TimeLimit,
			LimitBytesTotal: dataLimitBytes,
			Comment:         comment,
		}

		user, err := uc.hotspotSvc.AddUser(ctx, routerID, userReq)
		if err != nil {
			continue
		}

		createdVouchers = append(createdVouchers, dto.Voucher{
			Username:  user.Name,
			Password:  user.Password,
			Profile:   user.Profile,
			Server:    user.Server,
			Comment:   user.Comment,
			TimeLimit: req.TimeLimit,
			DataLimit: req.DataLimit,
		})
	}

	if uc.cache != nil {
		uc.cache.Delete(ctx, fmt.Sprintf("users:%d:*", routerID))
	}

	return &dto.VoucherBatchResult{
		Count:    len(createdVouchers),
		Comment:  comment,
		Vouchers: createdVouchers,
	}, nil
}

// GetVouchersByComment retrieves vouchers by comment (for printing)
func (uc *VoucherUseCase) GetVouchersByComment(ctx context.Context, routerID uint, comment string) ([]dto.HotspotUser, error) {
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

	if uc.cache != nil {
		uc.cache.Delete(ctx, fmt.Sprintf("users:%d:*", routerID))
	}

	return nil
}

// CacheGeneratedVouchers replicates legacy cache behavior before printing vouchers
func (uc *VoucherUseCase) CacheGeneratedVouchers(
	ctx context.Context,
	routerID uint,
	user string,
	gencode string,
	gcomment string,
) (int, string, []dto.HotspotUser, error) {
	_, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return 0, "", nil, err
	}

	comment := fmt.Sprintf("%s-%s-%s", user, gencode, time.Now().Format("01.02.06"))
	if strings.TrimSpace(gcomment) != "" {
		comment += "-" + strings.TrimSpace(gcomment)
	}

	users, err := uc.hotspotSvc.GetUsers(ctx, routerID, dto.UserFilter{
		Comment: comment,
	})
	if err != nil {
		return 0, "", nil, err
	}

	filtered := make([]dto.HotspotUser, 0, len(users))
	for _, u := range users {
		if u.Uptime == "" || u.Uptime == "0s" {
			filtered = append(filtered, u)
		}
	}

	return len(filtered), comment, filtered, nil
}

func normalizeGencode(gencode string) string {
	gencode = strings.TrimSpace(gencode)
	if gencode != "" {
		return gencode
	}
	n, err := rand.Int(rand.Reader, big.NewInt(899))
	if err != nil {
		return "101"
	}
	return fmt.Sprintf("%03d", n.Int64()+101)
}

func validateVoucherGenerateRequest(req dto.VoucherGenerateRequest) error {
	if req.Mode == "up" {
		switch req.CharacterSet {
		case "lower", "upper", "upplow", "mix", "mix1", "mix2":
			return nil
		default:
			return errors.New("invalid characterSet for mode=up")
		}
	}

	if req.Mode == "vc" {
		switch req.CharacterSet {
		case "lower1", "upper1", "upplow1", "mix", "mix1", "mix2", "num":
			return nil
		default:
			return errors.New("invalid characterSet for mode=vc")
		}
	}

	return errors.New("invalid mode")
}
