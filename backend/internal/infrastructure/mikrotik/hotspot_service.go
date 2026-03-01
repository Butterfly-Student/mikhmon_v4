package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
)

// HotspotService provides high-level hotspot operations across multiple routers.
// It resolves the router's Manager name via routerRepo, then delegates to the
// per-router *Client obtained from the Manager.
type HotspotService struct {
	manager    *Manager
	routerRepo repository.RouterRepository
}

// NewHotspotService creates a new HotspotService.
func NewHotspotService(manager *Manager, routerRepo repository.RouterRepository) *HotspotService {
	return &HotspotService{
		manager:    manager,
		routerRepo: routerRepo,
	}
}

// client resolves the *Client for the given routerID, connecting on demand.
func (s *HotspotService) client(ctx context.Context, routerID uint) (*Client, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}
	cfg := Config{
		Host:     router.Host,
		Port:     router.Port,
		Username: router.Username,
		Password: router.Password,
		UseTLS:   router.UseSSL,
		Timeout:  time.Duration(router.Timeout) * time.Second,
	}
	c, err := s.manager.GetOrConnect(ctx, router.Name, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users from MikroTik: %w", err)
	}
	return c, nil
}

// AddUser adds a hotspot user to the specified router.
func (s *HotspotService) AddUser(ctx context.Context, routerID uint, req dto.CreateUserRequest) (*dto.HotspotUser, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	user := &dto.HotspotUser{
		Name:            req.Name,
		Password:        req.Password,
		Profile:         req.Profile,
		Server:          req.Server,
		MACAddress:      req.MACAddress,
		LimitUptime:     req.LimitUptime,
		LimitBytesTotal: req.LimitBytesTotal,
		Comment:         req.Comment,
		Disabled:        req.Disabled,
	}

	id, err := c.AddHotspotUser(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

// GetUsers retrieves hotspot users from the specified router.
func (s *HotspotService) GetUsers(ctx context.Context, routerID uint, filter dto.UserFilter) ([]dto.HotspotUser, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	var users []*dto.HotspotUser
	if filter.Comment != "" {
		users, err = c.GetHotspotUsersByComment(ctx, filter.Comment)
	} else {
		users, err = c.GetHotspotUsers(ctx, filter.Profile)
	}
	if err != nil {
		return nil, err
	}

	result := make([]dto.HotspotUser, len(users))
	for i, u := range users {
		result[i] = *u
	}
	return result, nil
}

// GetUsersCount retrieves total hotspot users count from the specified router.
func (s *HotspotService) GetUsersCount(ctx context.Context, routerID uint) (int, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return 0, err
	}
	return c.GetHotspotUsersCount(ctx)
}

// RemoveUser removes a hotspot user from the specified router.
func (s *HotspotService) RemoveUser(ctx context.Context, routerID uint, userID string) error {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return err
	}
	return c.RemoveHotspotUser(ctx, userID)
}

// GetAddressPools retrieves address pools from the specified router.
func (s *HotspotService) GetAddressPools(ctx context.Context, routerID uint) ([]string, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return c.GetAddressPools(ctx)
}

// GetParentQueues retrieves parent queues from the specified router.
func (s *HotspotService) GetParentQueues(ctx context.Context, routerID uint) ([]string, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return c.GetAllParentQueues(ctx)
}

// GetAllQueues retrieves all queues from the specified router.
func (s *HotspotService) GetAllQueues(ctx context.Context, routerID uint) ([]string, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return c.GetAllQueues(ctx)
}

// GetServers retrieves hotspot servers from the specified router.
func (s *HotspotService) GetServers(ctx context.Context, routerID uint) ([]string, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return c.GetHotspotServers(ctx)
}

// GetUserByID retrieves a specific user by ID.
func (s *HotspotService) GetUserByID(ctx context.Context, routerID uint, userID string) (*dto.HotspotUser, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return c.GetHotspotUserByID(ctx, userID)
}

// UpdateUser updates a hotspot user.
func (s *HotspotService) UpdateUser(ctx context.Context, routerID uint, userID string, req dto.UpdateUserRequest) error {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return err
	}

	user := &dto.HotspotUser{
		Name:            req.Name,
		Password:        req.Password,
		Profile:         req.Profile,
		Server:          req.Server,
		MACAddress:      req.MACAddress,
		LimitUptime:     req.TimeLimit,
		LimitBytesTotal: ParseDataLimit(req.DataLimit),
		Comment:         req.Comment,
		Disabled:        req.Disabled,
	}

	return c.UpdateHotspotUser(ctx, userID, user)
}

// GetUserProfiles retrieves user profiles from the specified router.
func (s *HotspotService) GetUserProfiles(ctx context.Context, routerID uint) ([]dto.UserProfile, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	profiles, err := c.GetUserProfiles(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserProfile, len(profiles))
	for i, p := range profiles {
		result[i] = *p
	}
	return result, nil
}

// AddUserProfile adds a user profile to the specified router.
func (s *HotspotService) AddUserProfile(ctx context.Context, routerID uint, req dto.ProfileRequest) (*dto.UserProfile, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	profile := &dto.UserProfile{
		Name:         req.Name,
		AddressPool:  req.AddressPool,
		SharedUsers:  req.SharedUsers,
		RateLimit:    req.RateLimit,
		ParentQueue:  req.ParentQueue,
		ExpireMode:   req.ExpireMode,
		Validity:     req.Validity,
		Price:        req.Price,
		SellingPrice: req.SellingPrice,
		LockUser:     req.LockUser,
		LockServer:   req.LockServer,
	}

	id, err := c.AddUserProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	profile.ID = id
	return profile, nil
}

// UpdateUserProfile updates a user profile.
func (s *HotspotService) UpdateUserProfile(ctx context.Context, routerID uint, profileID string, req dto.ProfileUpdateRequest) (*dto.UserProfile, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	profile := &dto.UserProfile{
		Name:         req.Name,
		AddressPool:  req.AddressPool,
		SharedUsers:  req.SharedUsers,
		RateLimit:    req.RateLimit,
		ParentQueue:  req.ParentQueue,
		ExpireMode:   req.ExpireMode,
		Validity:     req.Validity,
		Price:        req.Price,
		SellingPrice: req.SellingPrice,
		LockUser:     req.LockUser,
		LockServer:   req.LockServer,
	}

	if err := c.UpdateUserProfile(ctx, profileID, profile); err != nil {
		return nil, err
	}

	profile.ID = profileID
	return profile, nil
}

// DeleteUserProfile deletes a user profile.
func (s *HotspotService) DeleteUserProfile(ctx context.Context, routerID uint, profileID string) error {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return err
	}
	return c.RemoveUserProfile(ctx, profileID)
}

// GetActiveUsers retrieves active hotspot sessions.
func (s *HotspotService) GetActiveUsers(ctx context.Context, routerID uint) ([]dto.HotspotActive, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	actives, err := c.GetHotspotActive(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.HotspotActive, len(actives))
	for i, a := range actives {
		result[i] = *a
	}
	return result, nil
}

// GetHosts retrieves hotspot hosts.
func (s *HotspotService) GetHosts(ctx context.Context, routerID uint) ([]dto.HotspotHost, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}

	hosts, err := c.GetHotspotHosts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.HotspotHost, len(hosts))
	for i, h := range hosts {
		result[i] = *h
	}
	return result, nil
}

// RemoveActiveUser removes an active hotspot session.
func (s *HotspotService) RemoveActiveUser(ctx context.Context, routerID uint, activeID string) error {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return err
	}
	return c.RemoveHotspotActive(ctx, activeID)
}

// RemoveHost removes a hotspot host entry.
func (s *HotspotService) RemoveHost(ctx context.Context, routerID uint, hostID string) error {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return err
	}
	return c.RemoveHotspotHost(ctx, hostID)
}

// ResetUserCounters resets a hotspot user's counters.
func (s *HotspotService) ResetUserCounters(ctx context.Context, routerID uint, userID string) error {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return err
	}
	return c.ResetHotspotUserCounters(ctx, userID)
}

// SetupExpireMonitor ensures expire monitor scheduler exists and is enabled.
func (s *HotspotService) SetupExpireMonitor(ctx context.Context, routerID uint, script string) (string, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return "", err
	}
	return c.EnsureExpireMonitor(ctx, script)
}

// GetActiveCount retrieves the count of active hotspot sessions.
func (s *HotspotService) GetActiveCount(ctx context.Context, routerID uint) (int, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return 0, err
	}
	return c.GetHotspotActiveCount(ctx)
}

// GetProfileByID retrieves a user profile by ID.
func (s *HotspotService) GetProfileByID(ctx context.Context, routerID uint, profileID string) (*dto.UserProfile, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return c.GetUserProfileByID(ctx, profileID)
}

// GetProfileByName retrieves a user profile by name.
func (s *HotspotService) GetProfileByName(ctx context.Context, routerID uint, profileName string) (*dto.UserProfile, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return c.GetUserProfileByName(ctx, profileName)
}

// ListenActive starts a streaming subscription to active hotspot sessions.
// No timeout is applied — the caller controls lifetime via ctx.
func (s *HotspotService) ListenActive(
	ctx context.Context,
	routerID uint,
	resultChan chan<- []*dto.HotspotActive,
) (func() error, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return c.ListenHotspotActive(ctx, resultChan)
}

// ListenInactive starts a streaming subscription to inactive hotspot users.
// No timeout is applied — the caller controls lifetime via ctx.
func (s *HotspotService) ListenInactive(
	ctx context.Context,
	routerID uint,
	resultChan chan<- []*dto.HotspotUser,
) (func() error, error) {
	c, err := s.client(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return c.ListenHotspotInactive(ctx, resultChan)
}

// UserFilter for filtering users.
type UserFilter struct {
	Profile string
	Comment string
}
