package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
)

// HotspotUseCase handles hotspot business logic
type HotspotUseCase struct {
	hotspotSvc *mikrotik.HotspotService
}

// NewHotspotUseCase creates a new hotspot use case
func NewHotspotUseCase(hotspotSvc *mikrotik.HotspotService) *HotspotUseCase {
	return &HotspotUseCase{
		hotspotSvc: hotspotSvc,
	}
}

// GetUsers retrieves hotspot users with timeout
func (uc *HotspotUseCase) GetUsers(ctx context.Context, routerID uint, profile string, force bool) ([]*dto.HotspotUser, error) {
	// Use shorter timeout for MikroTik operations
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := dto.UserFilter{
		Profile: profile,
	}
	users, err := uc.hotspotSvc.GetUsers(ctx, routerID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get users from MikroTik: %w", err)
	}

	// Convert to pointers
	result := make([]*dto.HotspotUser, len(users))
	for i := range users {
		result[i] = &users[i]
	}
	return result, nil
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
		LimitBytesTotal: 0,
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

// GetActive retrieves active hotspot sessions with timeout
func (uc *HotspotUseCase) GetActive(ctx context.Context, routerID uint, force bool) ([]*dto.HotspotActive, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	actives, err := uc.hotspotSvc.GetActiveUsers(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users from MikroTik: %w", err)
	}

	// Convert to pointers
	result := make([]*dto.HotspotActive, len(actives))
	for i := range actives {
		result[i] = &actives[i]
	}
	return result, nil
}

// GetHosts retrieves hotspot hosts with timeout
func (uc *HotspotUseCase) GetHosts(ctx context.Context, routerID uint, force bool) ([]*dto.HotspotHost, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	hosts, err := uc.hotspotSvc.GetHosts(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hosts from MikroTik: %w", err)
	}

	// Convert to pointers
	result := make([]*dto.HotspotHost, len(hosts))
	for i := range hosts {
		result[i] = &hosts[i]
	}
	return result, nil
}

// GetProfiles retrieves user profiles with timeout
func (uc *HotspotUseCase) GetProfiles(ctx context.Context, routerID uint, force bool) ([]*dto.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	profiles, err := uc.hotspotSvc.GetUserProfiles(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles from MikroTik: %w", err)
	}

	// Convert to pointers
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

// GetAddressPools retrieves address pools (placeholder)
func (uc *HotspotUseCase) GetAddressPools(ctx context.Context, routerID uint) ([]string, error) {
	return []string{"none"}, nil
}

// GetParentQueues retrieves parent queues (placeholder)
func (uc *HotspotUseCase) GetParentQueues(ctx context.Context, routerID uint) ([]string, error) {
	return []string{"none"}, nil
}

// GetServers retrieves hotspot servers (placeholder)
func (uc *HotspotUseCase) GetServers(ctx context.Context, routerID uint) ([]string, error) {
	return []string{"all"}, nil
}

// ResetUserCounters resets user counters (placeholder)
func (uc *HotspotUseCase) ResetUserCounters(ctx context.Context, routerID uint, userID string) error {
	return nil
}
