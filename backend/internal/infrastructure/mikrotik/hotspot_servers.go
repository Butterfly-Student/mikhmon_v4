package mikrotik

import (
	"context"
)

// GetHotspotServers retrieves all hotspot server names.
func (c *Client) GetHotspotServers(ctx context.Context) ([]string, error) {
	reply, err := c.RunContext(ctx, "/ip/hotspot/print")
	if err != nil {
		return nil, err
	}

	servers := []string{"all"}
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			servers = append(servers, name)
		}
	}

	return servers, nil
}
