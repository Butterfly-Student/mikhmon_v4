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

// HotspotUseCase handles hotspot business logic
type HotspotUseCase struct {
	routerRepo  repository.RouterRepository
	hotspotSvc  *mikrotik.HotspotService
	mikrotikSvc *mikrotik.Client
	log         *zap.Logger
}

// NewHotspotUseCase creates a new hotspot use case
func NewHotspotUseCase(
	routerRepo repository.RouterRepository,
	hotspotSvc *mikrotik.HotspotService,
	mikrotikSvc *mikrotik.Client,
	log *zap.Logger,
) *HotspotUseCase {
	if log == nil {
		log = zap.NewNop()
	}
	return &HotspotUseCase{
		routerRepo:  routerRepo,
		hotspotSvc:  hotspotSvc,
		mikrotikSvc: mikrotikSvc,
		log:         log.Named("hotspot-usecase"),
	}
}

// GetUsers retrieves hotspot users with timeout
func (uc *HotspotUseCase) GetUsers(ctx context.Context, routerID uint, profile string) ([]*dto.HotspotUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := dto.UserFilter{
		Profile: profile,
	}
	users, err := uc.hotspotSvc.GetUsers(ctx, routerID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get users from MikroTik: %w", err)
	}

	result := make([]*dto.HotspotUser, len(users))
	for i := range users {
		result[i] = &users[i]
	}
	return result, nil
}

// GetUser retrieves a specific hotspot user by ID with timeout
func (uc *HotspotUseCase) GetUser(ctx context.Context, routerID uint, userID string) (*dto.HotspotUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := uc.hotspotSvc.GetUserByID(ctx, routerID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from MikroTik: %w", err)
	}
	return user, nil
}

// AddUser adds a new hotspot user
func (uc *HotspotUseCase) AddUser(ctx context.Context, req *dto.AddUserRequest) error {
	createReq := dto.CreateUserRequest{
		Server:          req.Server,
		Name:            req.Name,
		Password:        req.Password,
		Profile:         req.Profile,
		MACAddress:      req.MACAddress,
		LimitUptime:     req.TimeLimit,
		LimitBytesTotal: mikrotik.ParseDataLimit(req.DataLimit),
		Comment:         req.Comment,
	}
	_, err := uc.hotspotSvc.AddUser(ctx, req.RouterID, createReq)
	return err
}

// UpdateUser updates a hotspot user
func (uc *HotspotUseCase) UpdateUser(ctx context.Context, routerID uint, userID string, req *dto.UpdateUserRequest) error {
	return uc.hotspotSvc.UpdateUser(ctx, routerID, userID, *req)
}

// RemoveUser removes a hotspot user
func (uc *HotspotUseCase) RemoveUser(ctx context.Context, routerID uint, userID string) error {
	return uc.hotspotSvc.RemoveUser(ctx, routerID, userID)
}

// ResetUserCounters resets user counters
func (uc *HotspotUseCase) ResetUserCounters(ctx context.Context, routerID uint, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return uc.hotspotSvc.ResetUserCounters(ctx, routerID, userID)
}

// GetProfiles retrieves user profiles with timeout
func (uc *HotspotUseCase) GetProfiles(ctx context.Context, routerID uint) ([]*dto.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	profiles, err := uc.hotspotSvc.GetUserProfiles(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles from MikroTik: %w", err)
	}

	result := make([]*dto.UserProfile, len(profiles))
	for i := range profiles {
		result[i] = &profiles[i]
	}
	return result, nil
}

// AddProfile adds a user profile
func (uc *HotspotUseCase) AddProfile(ctx context.Context, routerID uint, req *dto.ProfileRequest) error {
	_, err := uc.hotspotSvc.AddUserProfile(ctx, routerID, *req)
	return err
}

// UpdateProfile updates a user profile
func (uc *HotspotUseCase) UpdateProfile(ctx context.Context, routerID uint, profileID string, req *dto.ProfileUpdateRequest) error {
	_, err := uc.hotspotSvc.UpdateUserProfile(ctx, routerID, profileID, *req)
	return err
}

// RemoveProfile removes a user profile
func (uc *HotspotUseCase) RemoveProfile(ctx context.Context, routerID uint, profileID string) error {
	return uc.hotspotSvc.DeleteUserProfile(ctx, routerID, profileID)
}

// GetActive retrieves active hotspot sessions with timeout
func (uc *HotspotUseCase) GetActive(ctx context.Context, routerID uint) ([]*dto.HotspotActive, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	actives, err := uc.hotspotSvc.GetActiveUsers(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users from MikroTik: %w", err)
	}

	result := make([]*dto.HotspotActive, len(actives))
	for i := range actives {
		result[i] = &actives[i]
	}
	return result, nil
}

// RemoveActive removes an active hotspot session
func (uc *HotspotUseCase) RemoveActive(ctx context.Context, routerID uint, activeID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return uc.hotspotSvc.RemoveActiveUser(ctx, routerID, activeID)
}

// GetHosts retrieves hotspot hosts with timeout
func (uc *HotspotUseCase) GetHosts(ctx context.Context, routerID uint) ([]*dto.HotspotHost, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	hosts, err := uc.hotspotSvc.GetHosts(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hosts from MikroTik: %w", err)
	}

	result := make([]*dto.HotspotHost, len(hosts))
	for i := range hosts {
		result[i] = &hosts[i]
	}
	return result, nil
}

// RemoveHost removes a hotspot host entry
func (uc *HotspotUseCase) RemoveHost(ctx context.Context, routerID uint, hostID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return uc.hotspotSvc.RemoveHost(ctx, routerID, hostID)
}

// GetServers retrieves hotspot servers
func (uc *HotspotUseCase) GetServers(ctx context.Context, routerID uint) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	servers, err := uc.hotspotSvc.GetServers(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hotspot servers from MikroTik: %w", err)
	}
	return servers, nil
}

// GetActiveCount retrieves count of active hotspot sessions
func (uc *HotspotUseCase) GetActiveCount(ctx context.Context, routerID uint) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return 0, err
	}

	count, err := uc.mikrotikSvc.GetHotspotActiveCount(ctx, router)
	if err != nil {
		return 0, fmt.Errorf("failed to get active users count from MikroTik: %w", err)
	}
	return count, nil
}

// GetProfileByID retrieves a profile by ID
func (uc *HotspotUseCase) GetProfileByID(ctx context.Context, routerID uint, profileID string) (*dto.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	profile, err := uc.mikrotikSvc.GetUserProfileByID(ctx, router, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile from MikroTik: %w", err)
	}
	return profile, nil
}

// GetProfileByName retrieves a profile by name
func (uc *HotspotUseCase) GetProfileByName(ctx context.Context, routerID uint, profileName string) (*dto.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	router, err := uc.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	profile, err := uc.mikrotikSvc.GetUserProfileByName(ctx, router, profileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile from MikroTik: %w", err)
	}
	return profile, nil
}

// SetupExpireMonitor sets up or enables Mikhmon expire monitor scheduler
func (uc *HotspotUseCase) SetupExpireMonitor(ctx context.Context, routerID uint, script string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	status, err := uc.hotspotSvc.SetupExpireMonitor(ctx, routerID, script)
	if err != nil {
		return "", fmt.Errorf("failed to setup expire monitor: %w", err)
	}
	return status, nil
}

// GetExpireMonitorScript returns default expire monitor script
func (uc *HotspotUseCase) GetExpireMonitorScript() string {
	return mikrotik.NewOnLoginGenerator().GenerateExpireMonitorScript()
}
