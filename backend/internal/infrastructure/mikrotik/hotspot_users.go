package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"go.uber.org/zap"
)

// GetHotspotUsers retrieves all hotspot users, optionally filtered by profile
func (c *Client) GetHotspotUsers(ctx context.Context, router *entity.Router, profile string) ([]*dto.HotspotUser, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	args := []string{"/ip/hotspot/user/print"}
	if profile != "" {
		args = append(args, "?profile="+profile)
	}

	reply, err := client.RunContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	users := make([]*dto.HotspotUser, 0, len(reply.Re))
	for _, re := range reply.Re {
		users = append(users, &dto.HotspotUser{
			ID:              re.Map[".id"],
			Server:          re.Map["server"],
			Name:            re.Map["name"],
			Password:        re.Map["password"],
			Profile:         re.Map["profile"],
			MACAddress:      re.Map["mac-address"],
			IPAddress:       re.Map["address"],
			Uptime:          re.Map["uptime"],
			BytesIn:         parseInt(re.Map["bytes-in"]),
			BytesOut:        parseInt(re.Map["bytes-out"]),
			PacketsIn:       parseInt(re.Map["packets-in"]),
			PacketsOut:      parseInt(re.Map["packets-out"]),
			LimitUptime:     re.Map["limit-uptime"],
			LimitBytesIn:    parseInt(re.Map["limit-bytes-in"]),
			LimitBytesOut:   parseInt(re.Map["limit-bytes-out"]),
			LimitBytesTotal: parseInt(re.Map["limit-bytes-total"]),
			Comment:         re.Map["comment"],
			Disabled:        parseBool(re.Map["disabled"]),
			Email:           re.Map["email"],
		})
	}

	return users, nil
}

// GetHotspotUsersByComment retrieves hotspot users by comment (for vouchers)
func (c *Client) GetHotspotUsersByComment(ctx context.Context, router *entity.Router, comment string) ([]*dto.HotspotUser, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	args := []string{"/ip/hotspot/user/print"}
	if comment != "" {
		args = append(args, "?comment="+comment)
	}

	reply, err := client.RunContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	users := make([]*dto.HotspotUser, 0, len(reply.Re))
	for _, re := range reply.Re {
		users = append(users, &dto.HotspotUser{
			ID:              re.Map[".id"],
			Server:          re.Map["server"],
			Name:            re.Map["name"],
			Password:        re.Map["password"],
			Profile:         re.Map["profile"],
			MACAddress:      re.Map["mac-address"],
			IPAddress:       re.Map["address"],
			Uptime:          re.Map["uptime"],
			BytesIn:         parseInt(re.Map["bytes-in"]),
			BytesOut:        parseInt(re.Map["bytes-out"]),
			LimitUptime:     re.Map["limit-uptime"],
			LimitBytesTotal: parseInt(re.Map["limit-bytes-total"]),
			Comment:         re.Map["comment"],
			Disabled:        parseBool(re.Map["disabled"]),
		})
	}

	return users, nil
}

// GetHotspotUserByID retrieves a hotspot user by ID
func (c *Client) GetHotspotUserByID(ctx context.Context, router *entity.Router, id string) (*dto.HotspotUser, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/user/print", "?.id="+id)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, nil
	}

	re := reply.Re[0]
	return &dto.HotspotUser{
		ID:              re.Map[".id"],
		Server:          re.Map["server"],
		Name:            re.Map["name"],
		Password:        re.Map["password"],
		Profile:         re.Map["profile"],
		MACAddress:      re.Map["mac-address"],
		IPAddress:       re.Map["address"],
		Uptime:          re.Map["uptime"],
		BytesIn:         parseInt(re.Map["bytes-in"]),
		BytesOut:        parseInt(re.Map["bytes-out"]),
		LimitUptime:     re.Map["limit-uptime"],
		LimitBytesTotal: parseInt(re.Map["limit-bytes-total"]),
		Comment:         re.Map["comment"],
		Disabled:        parseBool(re.Map["disabled"]),
	}, nil
}

// GetHotspotUserByName retrieves a hotspot user by name
func (c *Client) GetHotspotUserByName(ctx context.Context, router *entity.Router, name string) (*dto.HotspotUser, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/user/print", "?name="+name)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, nil
	}

	re := reply.Re[0]
	return &dto.HotspotUser{
		ID:              re.Map[".id"],
		Server:          re.Map["server"],
		Name:            re.Map["name"],
		Password:        re.Map["password"],
		Profile:         re.Map["profile"],
		MACAddress:      re.Map["mac-address"],
		IPAddress:       re.Map["address"],
		Uptime:          re.Map["uptime"],
		BytesIn:         parseInt(re.Map["bytes-in"]),
		BytesOut:        parseInt(re.Map["bytes-out"]),
		LimitUptime:     re.Map["limit-uptime"],
		LimitBytesTotal: parseInt(re.Map["limit-bytes-total"]),
		Comment:         re.Map["comment"],
		Disabled:        parseBool(re.Map["disabled"]),
	}, nil
}

// AddHotspotUser adds a new hotspot user
func (c *Client) AddHotspotUser(ctx context.Context, router *entity.Router, user *dto.HotspotUser) (string, error) {
	client, err := c.getClient(router)
	if err != nil {
		return "", err
	}

	args := []string{
		"/ip/hotspot/user/add",
		"=name=" + user.Name,
		"=profile=" + user.Profile,
		"=disabled=no",
	}

	if user.Server != "" && user.Server != "all" {
		args = append(args, "=server="+user.Server)
	}
	if user.Password != "" {
		args = append(args, "=password="+user.Password)
	}
	if user.MACAddress != "" {
		args = append(args, "=mac-address="+user.MACAddress)
	}
	if user.LimitUptime != "" {
		args = append(args, "=limit-uptime="+user.LimitUptime)
	}
	if user.LimitBytesTotal > 0 {
		args = append(args, "=limit-bytes-total="+formatInt(user.LimitBytesTotal))
	}
	if user.Comment != "" {
		args = append(args, "=comment="+user.Comment)
	}
	if user.Email != "" {
		args = append(args, "=email="+user.Email)
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

// UpdateHotspotUser updates an existing hotspot user
func (c *Client) UpdateHotspotUser(ctx context.Context, router *entity.Router, id string, user *dto.HotspotUser) error {
	client, err := c.getClient(router)
	if err != nil {
		return err
	}

	args := []string{
		"/ip/hotspot/user/set",
		"=.id=" + id,
	}

	if user.Name != "" {
		args = append(args, "=name="+user.Name)
	}
	if user.Password != "" {
		args = append(args, "=password="+user.Password)
	}
	if user.Profile != "" {
		args = append(args, "=profile="+user.Profile)
	}
	if user.Server != "" {
		args = append(args, "=server="+user.Server)
	}
	if user.MACAddress != "" {
		args = append(args, "=mac-address="+user.MACAddress)
	}
	if user.LimitUptime != "" {
		args = append(args, "=limit-uptime="+user.LimitUptime)
	}
	if user.LimitBytesTotal > 0 {
		args = append(args, "=limit-bytes-total="+formatInt(user.LimitBytesTotal))
	}
	if user.Comment != "" {
		args = append(args, "=comment="+user.Comment)
	}
	if user.Email != "" {
		args = append(args, "=email="+user.Email)
	}
	if user.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
	}

	_, err = client.RunContext(ctx, args...)
	return err
}

// RemoveHotspotUser removes a hotspot user
func (c *Client) RemoveHotspotUser(ctx context.Context, router *entity.Router, id string) error {
	client, err := c.getClient(router)
	if err != nil {
		return err
	}

	_, err = client.RunContext(ctx, "/ip/hotspot/user/remove", "=.id="+id)
	return err
}

// RemoveHotspotUsersByComment removes hotspot users by comment (for voucher deletion)
func (c *Client) RemoveHotspotUsersByComment(ctx context.Context, router *entity.Router, comment string) error {
	// First get users with this comment
	users, err := c.GetHotspotUsersByComment(ctx, router, comment)
	if err != nil {
		return err
	}

	// Remove each user
	for _, user := range users {
		if err := c.RemoveHotspotUser(ctx, router, user.ID); err != nil {
			// Log error but continue removing other users
			c.log.Warn("Failed to remove hotspot user",
				zap.String("userName", user.Name),
				zap.String("userID", user.ID),
				zap.Error(err),
			)
		}
	}

	return nil
}

// ResetHotspotUserCounters resets counters for a hotspot user
func (c *Client) ResetHotspotUserCounters(ctx context.Context, router *entity.Router, id string) error {
	client, err := c.getClient(router)
	if err != nil {
		return err
	}

	_, err = client.RunContext(ctx, "/ip/hotspot/user/reset-counters", "=.id="+id)
	return err
}

// GetHotspotUsersCount retrieves the count of hotspot users.
// Count is reduced by 1 to exclude the admin user, matching MIKHMON PHP original behavior.
func (c *Client) GetHotspotUsersCount(ctx context.Context, router *entity.Router) (int, error) {
	client, err := c.getClient(router)
	if err != nil {
		return 0, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/user/print", "=count-only=")
	if err != nil {
		return 0, err
	}

	if len(reply.Re) > 0 {
		// Subtract 1 to exclude the admin user (matching MIKHMON original behavior)
		count := int(parseInt(reply.Re[0].Map["ret"])) - 1
		if count < 0 {
			count = 0
		}
		return count, nil
	}

	return 0, nil
}
