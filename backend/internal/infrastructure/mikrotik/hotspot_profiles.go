package mikrotik

import (
	"context"
	"strings"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

// GetUserProfiles retrieves all hotspot user profiles
func (c *Client) GetUserProfiles(ctx context.Context, router *entity.Router) ([]*dto.UserProfile, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/user/profile/print")
	if err != nil {
		return nil, err
	}

	generator := NewOnLoginGenerator()
	profiles := make([]*dto.UserProfile, 0, len(reply.Re))

	for _, re := range reply.Re {
		profile := &dto.UserProfile{
			ID:                re.Map[".id"],
			Name:              re.Map["name"],
			AddressPool:       re.Map["address-pool"],
			SharedUsers:       int(parseInt(re.Map["shared-users"])),
			RateLimit:         re.Map["rate-limit"],
			ParentQueue:       re.Map["parent-queue"],
			StatusAutorefresh: re.Map["status-autorefresh"],
			OnLogin:           re.Map["on-login"],
			OnLogout:          re.Map["on-logout"],
			OnUp:              re.Map["on-up"],
			OnDown:            re.Map["on-down"],
			TransparentProxy:  parseBool(re.Map["transparent-proxy"]),
			OpenStatusPage:    re.Map["open-status-page"],
			Advertise:         parseBool(re.Map["advertise"]),
			AdvertiseInterval: re.Map["advertise-interval"],
			AdvertiseTimeout:  re.Map["advertise-timeout"],
			AdvertiseURL:      re.Map["advertise-url"],
		}

		// Parse Mikhmon metadata from on-login script
		if profile.OnLogin != "" {
			parsed := generator.Parse(profile.OnLogin)
			profile.ExpireMode = parsed.ExpireMode
			profile.Validity = parsed.Validity
			profile.Price = parsed.Price
			profile.SellingPrice = parsed.SellingPrice
			profile.LockUser = parsed.LockUser
			profile.LockServer = parsed.LockServer
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// GetUserProfileByID retrieves a user profile by ID
func (c *Client) GetUserProfileByID(ctx context.Context, router *entity.Router, id string) (*dto.UserProfile, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/user/profile/print", "?.id="+id)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, nil
	}

	re := reply.Re[0]
	generator := NewOnLoginGenerator()

	profile := &dto.UserProfile{
		ID:                re.Map[".id"],
		Name:              re.Map["name"],
		AddressPool:       re.Map["address-pool"],
		SharedUsers:       int(parseInt(re.Map["shared-users"])),
		RateLimit:         re.Map["rate-limit"],
		ParentQueue:       re.Map["parent-queue"],
		StatusAutorefresh: re.Map["status-autorefresh"],
		OnLogin:           re.Map["on-login"],
		OnLogout:          re.Map["on-logout"],
		OnUp:              re.Map["on-up"],
		OnDown:            re.Map["on-down"],
		TransparentProxy:  parseBool(re.Map["transparent-proxy"]),
		OpenStatusPage:    re.Map["open-status-page"],
	}

	// Parse Mikhmon metadata
	if profile.OnLogin != "" {
		parsed := generator.Parse(profile.OnLogin)
		profile.ExpireMode = parsed.ExpireMode
		profile.Validity = parsed.Validity
		profile.Price = parsed.Price
		profile.SellingPrice = parsed.SellingPrice
		profile.LockUser = parsed.LockUser
		profile.LockServer = parsed.LockServer
	}

	return profile, nil
}

// GetUserProfileByName retrieves a user profile by name
func (c *Client) GetUserProfileByName(ctx context.Context, router *entity.Router, name string) (*dto.UserProfile, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/user/profile/print", "?name="+name)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, nil
	}

	re := reply.Re[0]
	generator := NewOnLoginGenerator()

	profile := &dto.UserProfile{
		ID:                re.Map[".id"],
		Name:              re.Map["name"],
		AddressPool:       re.Map["address-pool"],
		SharedUsers:       int(parseInt(re.Map["shared-users"])),
		RateLimit:         re.Map["rate-limit"],
		ParentQueue:       re.Map["parent-queue"],
		StatusAutorefresh: re.Map["status-autorefresh"],
		OnLogin:           re.Map["on-login"],
	}

	// Parse Mikhmon metadata
	if profile.OnLogin != "" {
		parsed := generator.Parse(profile.OnLogin)
		profile.ExpireMode = parsed.ExpireMode
		profile.Validity = parsed.Validity
		profile.Price = parsed.Price
		profile.SellingPrice = parsed.SellingPrice
		profile.LockUser = parsed.LockUser
		profile.LockServer = parsed.LockServer
	}

	return profile, nil
}

// AddUserProfile adds a new hotspot user profile with on-login script
func (c *Client) AddUserProfile(ctx context.Context, router *entity.Router, profile *dto.UserProfile) (string, error) {
	client, err := c.getClient(router)
	if err != nil {
		return "", err
	}

	// Generate on-login script if Mikhmon fields are provided
	var onLoginScript string
	if profile.ExpireMode != "" {
		generator := NewOnLoginGenerator()
		req := &dto.ProfileRequest{
			Name:         profile.Name,
			ExpireMode:   profile.ExpireMode,
			Validity:     profile.Validity,
			Price:        profile.Price,
			SellingPrice: profile.SellingPrice,
			LockUser:     profile.LockUser,
			LockServer:   profile.LockServer,
		}
		onLoginScript = generator.Generate(req)
	}

	args := []string{
		"/ip/hotspot/user/profile/add",
		"=name=" + profile.Name,
		"=status-autorefresh=1m",
	}

	if profile.AddressPool != "" && profile.AddressPool != "none" {
		args = append(args, "=address-pool="+profile.AddressPool)
	}
	if profile.SharedUsers > 0 {
		args = append(args, "=shared-users="+formatInt(int64(profile.SharedUsers)))
	}
	if profile.RateLimit != "" {
		args = append(args, "=rate-limit="+profile.RateLimit)
	}
	if profile.ParentQueue != "" && profile.ParentQueue != "none" {
		args = append(args, "=parent-queue="+profile.ParentQueue)
	}
	if onLoginScript != "" {
		// Replace semicolons with newlines for RouterOS script format
		script := strings.ReplaceAll(onLoginScript, "; ", "\n")
		args = append(args, "=on-login="+script)
	}

	reply, err := client.RunContext(ctx, args...)
	if err != nil {
		return "", err
	}

	if len(reply.Re) > 0 {
		return reply.Re[0].Map["ret"], nil
	}

	return "", nil
}

// UpdateUserProfile updates an existing user profile
func (c *Client) UpdateUserProfile(ctx context.Context, router *entity.Router, id string, profile *dto.UserProfile) error {
	client, err := c.getClient(router)
	if err != nil {
		return err
	}

	// Generate on-login script if Mikhmon fields are provided
	var onLoginScript string
	if profile.ExpireMode != "" {
		generator := NewOnLoginGenerator()
		req := &dto.ProfileRequest{
			Name:         profile.Name,
			ExpireMode:   profile.ExpireMode,
			Validity:     profile.Validity,
			Price:        profile.Price,
			SellingPrice: profile.SellingPrice,
			LockUser:     profile.LockUser,
			LockServer:   profile.LockServer,
		}
		onLoginScript = generator.Generate(req)
	}

	args := []string{
		"/ip/hotspot/user/profile/set",
		"=.id=" + id,
	}

	if profile.Name != "" {
		args = append(args, "=name="+profile.Name)
	}
	if profile.AddressPool != "" {
		args = append(args, "=address-pool="+profile.AddressPool)
	}
	if profile.SharedUsers > 0 {
		args = append(args, "=shared-users="+formatInt(int64(profile.SharedUsers)))
	}
	if profile.RateLimit != "" {
		args = append(args, "=rate-limit="+profile.RateLimit)
	}
	if profile.ParentQueue != "" {
		args = append(args, "=parent-queue="+profile.ParentQueue)
	}
	if onLoginScript != "" {
		// Replace semicolons with newlines for RouterOS script format
		script := strings.ReplaceAll(onLoginScript, "; ", "\n")
		args = append(args, "=on-login="+script)
	}

	_, err = client.RunContext(ctx, args...)
	return err
}

// RemoveUserProfile removes a user profile
func (c *Client) RemoveUserProfile(ctx context.Context, router *entity.Router, id string) error {
	client, err := c.getClient(router)
	if err != nil {
		return err
	}

	_, err = client.RunContext(ctx, "/ip/hotspot/user/profile/remove", "=.id="+id)
	return err
}
