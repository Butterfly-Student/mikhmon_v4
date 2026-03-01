package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// GetHotspotHosts retrieves hotspot hosts.
func (c *Client) GetHotspotHosts(ctx context.Context) ([]*dto.HotspotHost, error) {
	reply, err := c.RunContext(ctx, "/ip/hotspot/host/print")
	if err != nil {
		return nil, err
	}

	hosts := make([]*dto.HotspotHost, 0, len(reply.Re))
	for _, re := range reply.Re {
		hosts = append(hosts, &dto.HotspotHost{
			ID:           re.Map[".id"],
			MACAddress:   re.Map["mac-address"],
			Address:      re.Map["address"],
			ToAddress:    re.Map["to-address"],
			Server:       re.Map["server"],
			Authorized:   parseBool(re.Map["authorized"]),
			Bypassed:     parseBool(re.Map["bypassed"]),
			Blocked:      parseBool(re.Map["blocked"]),
			FoundBy:      re.Map["found-by"],
			HostDeadTime: re.Map["host-dead-time"],
			Comment:      re.Map["comment"],
		})
	}

	return hosts, nil
}

// RemoveHotspotHost removes a hotspot host.
func (c *Client) RemoveHotspotHost(ctx context.Context, id string) error {
	_, err := c.RunContext(ctx, "/ip/hotspot/host/remove", "=.id="+id)
	return err
}
