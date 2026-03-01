package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// parsePPPProfile maps a RouterOS sentence map to a PPPProfile DTO.
func parsePPPProfile(m map[string]string) *dto.PPPProfile {
	return &dto.PPPProfile{
		ID:                m[".id"],
		Name:              m["name"],
		LocalAddress:      m["local-address"],
		RemoteAddress:     m["remote-address"],
		DNSServer:         m["dns-server"],
		SessionTimeout:    m["session-timeout"],
		IdleTimeout:       m["idle-timeout"],
		OnlyOne:           parseBool(m["only-one"]),
		Comment:           m["comment"],
		RateLimit:         m["rate-limit"],
		ParentQueue:       m["parent-queue"],
		QueueType:         m["queue-type"],
		UseCompression:    parseBool(m["use-compression"]),
		UseEncryption:     parseBool(m["use-encryption"]),
		UseMPLS:           parseBool(m["use-mpls"]),
		UseUPnP:           parseBool(m["use-upnp"]),
		Bridge:            m["bridge"],
		AddressList:       m["address-list"],
		InterfaceList:     m["interface-list"],
		OnUp:              m["on-up"],
		OnDown:            m["on-down"],
		ChangeTCPMSS:      parseBool(m["change-tcp-mss"]),
		IncomingFilter:    m["incoming-filter"],
		OutgoingFilter:    m["outgoing-filter"],
		InsertQueueBefore: m["insert-queue-before"],
		WinsServer:        m["wins-server"],
		BridgeHorizon:     m["bridge-horizon"],
		BridgeLearning:    parseBool(m["bridge-learning"]),
		BridgePathCost:    int(parseInt(m["bridge-path-cost"])),
		BridgePortPriority: int(parseInt(m["bridge-port-priority"])),
	}
}

// GetPPPProfiles retrieves all PPP profiles.
func (c *Client) GetPPPProfiles(ctx context.Context) ([]*dto.PPPProfile, error) {
	reply, err := c.RunContext(ctx, "/ppp/profile/print")
	if err != nil {
		return nil, err
	}

	profiles := make([]*dto.PPPProfile, 0, len(reply.Re))
	for _, re := range reply.Re {
		profiles = append(profiles, parsePPPProfile(re.Map))
	}

	return profiles, nil
}

// GetPPPProfileByID retrieves a PPP profile by ID.
func (c *Client) GetPPPProfileByID(ctx context.Context, id string) (*dto.PPPProfile, error) {
	reply, err := c.RunContext(ctx, "/ppp/profile/print", "?.id="+id)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, nil
	}

	return parsePPPProfile(reply.Re[0].Map), nil
}

// GetPPPProfileByName retrieves a PPP profile by name.
func (c *Client) GetPPPProfileByName(ctx context.Context, name string) (*dto.PPPProfile, error) {
	reply, err := c.RunContext(ctx, "/ppp/profile/print", "?name="+name)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, nil
	}

	return parsePPPProfile(reply.Re[0].Map), nil
}

// AddPPPProfile adds a new PPP profile.
func (c *Client) AddPPPProfile(ctx context.Context, profile *dto.PPPProfile) error {
	args := []string{
		"/ppp/profile/add",
		"=name=" + profile.Name,
	}

	args = appendPPPProfileArgs(args, profile)

	_, err := c.RunArgsContext(ctx, args)
	return err
}

// UpdatePPPProfile updates an existing PPP profile.
func (c *Client) UpdatePPPProfile(ctx context.Context, id string, profile *dto.PPPProfile) error {
	args := []string{
		"/ppp/profile/set",
		"=.id=" + id,
	}

	if profile.Name != "" {
		args = append(args, "=name="+profile.Name)
	}
	args = appendPPPProfileArgs(args, profile)

	_, err := c.RunArgsContext(ctx, args)
	return err
}

// appendPPPProfileArgs appends optional PPP profile fields to an args slice.
func appendPPPProfileArgs(args []string, profile *dto.PPPProfile) []string {
	if profile.LocalAddress != "" {
		args = append(args, "=local-address="+profile.LocalAddress)
	}
	if profile.RemoteAddress != "" {
		args = append(args, "=remote-address="+profile.RemoteAddress)
	}
	if profile.DNSServer != "" {
		args = append(args, "=dns-server="+profile.DNSServer)
	}
	if profile.SessionTimeout != "" {
		args = append(args, "=session-timeout="+profile.SessionTimeout)
	}
	if profile.IdleTimeout != "" {
		args = append(args, "=idle-timeout="+profile.IdleTimeout)
	}
	if profile.RateLimit != "" {
		args = append(args, "=rate-limit="+profile.RateLimit)
	}
	if profile.ParentQueue != "" {
		args = append(args, "=parent-queue="+profile.ParentQueue)
	}
	if profile.QueueType != "" {
		args = append(args, "=queue-type="+profile.QueueType)
	}
	if profile.Bridge != "" {
		args = append(args, "=bridge="+profile.Bridge)
	}
	if profile.AddressList != "" {
		args = append(args, "=address-list="+profile.AddressList)
	}
	if profile.InterfaceList != "" {
		args = append(args, "=interface-list="+profile.InterfaceList)
	}
	if profile.OnUp != "" {
		args = append(args, "=on-up="+profile.OnUp)
	}
	if profile.OnDown != "" {
		args = append(args, "=on-down="+profile.OnDown)
	}
	if profile.IncomingFilter != "" {
		args = append(args, "=incoming-filter="+profile.IncomingFilter)
	}
	if profile.OutgoingFilter != "" {
		args = append(args, "=outgoing-filter="+profile.OutgoingFilter)
	}
	if profile.InsertQueueBefore != "" {
		args = append(args, "=insert-queue-before="+profile.InsertQueueBefore)
	}
	if profile.WinsServer != "" {
		args = append(args, "=wins-server="+profile.WinsServer)
	}
	if profile.Comment != "" {
		args = append(args, "=comment="+profile.Comment)
	}
	return args
}

// RemovePPPProfile removes a PPP profile.
func (c *Client) RemovePPPProfile(ctx context.Context, id string) error {
	_, err := c.RunContext(ctx, "/ppp/profile/remove", "=.id="+id)
	return err
}

// DisablePPPProfile disables a PPP profile.
func (c *Client) DisablePPPProfile(ctx context.Context, id string) error {
	_, err := c.RunContext(ctx, "/ppp/profile/disable", "=.id="+id)
	return err
}

// EnablePPPProfile enables a PPP profile.
func (c *Client) EnablePPPProfile(ctx context.Context, id string) error {
	_, err := c.RunContext(ctx, "/ppp/profile/enable", "=.id="+id)
	return err
}
