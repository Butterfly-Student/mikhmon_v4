package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
)

// HotspotService provides high-level hotspot operations
type HotspotService struct {
	client     *Client
	routerRepo repository.RouterRepository
}

// NewHotspotService creates a new hotspot service
func NewHotspotService(client *Client, routerRepo repository.RouterRepository) *HotspotService {
	return &HotspotService{
		client:     client,
		routerRepo: routerRepo,
	}
}

// AddUser adds a hotspot user to the specified router
func (s *HotspotService) AddUser(ctx context.Context, routerID uint, req dto.CreateUserRequest) (*dto.HotspotUser, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
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

	id, err := s.client.AddHotspotUser(ctx, router, user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

// GetUsers retrieves hotspot users from the specified router
func (s *HotspotService) GetUsers(ctx context.Context, routerID uint, filter dto.UserFilter) ([]dto.HotspotUser, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	var users []*dto.HotspotUser
	if filter.Comment != "" {
		users, err = s.client.GetHotspotUsersByComment(ctx, router, filter.Comment)
	} else {
		users, err = s.client.GetHotspotUsers(ctx, router, filter.Profile)
	}
	if err != nil {
		return nil, err
	}

	// Convert to slice of values
	result := make([]dto.HotspotUser, len(users))
	for i, u := range users {
		result[i] = *u
	}
	return result, nil
}

// GetUsersCount retrieves total hotspot users count from the specified router
func (s *HotspotService) GetUsersCount(ctx context.Context, routerID uint) (int, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return 0, err
	}

	return s.client.GetHotspotUsersCount(ctx, router)
}

// RemoveUser removes a hotspot user from the specified router
func (s *HotspotService) RemoveUser(ctx context.Context, routerID uint, userID string) error {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}

	return s.client.RemoveHotspotUser(ctx, router, userID)
}

// GetAddressPools retrieves address pools from the specified router
func (s *HotspotService) GetAddressPools(ctx context.Context, routerID uint) ([]string, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	return s.client.GetAddressPools(ctx, router)
}

// GetParentQueues retrieves parent queues from the specified router
func (s *HotspotService) GetParentQueues(ctx context.Context, routerID uint) ([]string, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	return s.client.GetAllParentQueues(ctx, router)
}

// GetParentQueues retrieves parent queues from the specified router
func (s *HotspotService) GetAllQueues(ctx context.Context, routerID uint) ([]string, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	return s.client.GetAllQueues(ctx, router)
}

// GetServers retrieves hotspot servers from the specified router
func (s *HotspotService) GetServers(ctx context.Context, routerID uint) ([]string, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	return s.client.GetHotspotServers(ctx, router)
}

// GetUserByID retrieves a specific user by ID
func (s *HotspotService) GetUserByID(ctx context.Context, routerID uint, userID string) (*dto.HotspotUser, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	return s.client.GetHotspotUserByID(ctx, router, userID)
}

// UpdateUser updates a hotspot user
func (s *HotspotService) UpdateUser(ctx context.Context, routerID uint, userID string, req dto.UpdateUserRequest) error {
	router, err := s.routerRepo.GetByID(ctx, routerID)
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

	return s.client.UpdateHotspotUser(ctx, router, userID, user)
}

// GetUserProfiles retrieves user profiles from the specified router
func (s *HotspotService) GetUserProfiles(ctx context.Context, routerID uint) ([]dto.UserProfile, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	profiles, err := s.client.GetUserProfiles(ctx, router)
	if err != nil {
		return nil, err
	}

	// Convert to slice of values
	result := make([]dto.UserProfile, len(profiles))
	for i, p := range profiles {
		result[i] = *p
	}
	return result, nil
}

// AddUserProfile adds a user profile to the specified router
func (s *HotspotService) AddUserProfile(ctx context.Context, routerID uint, req dto.ProfileRequest) (*dto.UserProfile, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
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

	id, err := s.client.AddUserProfile(ctx, router, profile)
	if err != nil {
		return nil, err
	}

	profile.ID = id
	return profile, nil
}

// UpdateUserProfile updates a user profile
func (s *HotspotService) UpdateUserProfile(ctx context.Context, routerID uint, profileID string, req dto.ProfileUpdateRequest) (*dto.UserProfile, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
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

	if err := s.client.UpdateUserProfile(ctx, router, profileID, profile); err != nil {
		return nil, err
	}

	profile.ID = profileID
	return profile, nil
}

// DeleteUserProfile deletes a user profile
func (s *HotspotService) DeleteUserProfile(ctx context.Context, routerID uint, profileID string) error {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}

	return s.client.RemoveUserProfile(ctx, router, profileID)
}

// GetActiveUsers retrieves active hotspot sessions
func (s *HotspotService) GetActiveUsers(ctx context.Context, routerID uint) ([]dto.HotspotActive, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	actives, err := s.client.GetHotspotActive(ctx, router)
	if err != nil {
		return nil, err
	}

	// Convert to slice of values
	result := make([]dto.HotspotActive, len(actives))
	for i, a := range actives {
		result[i] = *a
	}
	return result, nil
}

// GetHosts retrieves hotspot hosts
func (s *HotspotService) GetHosts(ctx context.Context, routerID uint) ([]dto.HotspotHost, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, err
	}

	hosts, err := s.client.GetHotspotHosts(ctx, router)
	if err != nil {
		return nil, err
	}

	// Convert to slice of values
	result := make([]dto.HotspotHost, len(hosts))
	for i, h := range hosts {
		result[i] = *h
	}
	return result, nil
}

// RemoveActiveUser removes an active hotspot session
func (s *HotspotService) RemoveActiveUser(ctx context.Context, routerID uint, activeID string) error {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}

	return s.client.RemoveHotspotActive(ctx, router, activeID)
}

// RemoveHost removes a hotspot host entry
func (s *HotspotService) RemoveHost(ctx context.Context, routerID uint, hostID string) error {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}

	return s.client.RemoveHotspotHost(ctx, router, hostID)
}

// ResetUserCounters resets a hotspot user's counters
func (s *HotspotService) ResetUserCounters(ctx context.Context, routerID uint, userID string) error {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return err
	}

	return s.client.ResetHotspotUserCounters(ctx, router, userID)
}

// SetupExpireMonitor ensures expire monitor scheduler exists and enabled
func (s *HotspotService) SetupExpireMonitor(ctx context.Context, routerID uint, script string) (string, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return "", err
	}

	return s.client.EnsureExpireMonitor(ctx, router, script)
}

// UserFilter for filtering users
type UserFilter struct {
	Profile string
	Comment string
}


