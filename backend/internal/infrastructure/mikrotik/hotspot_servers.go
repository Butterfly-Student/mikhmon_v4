package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
)




// GetHotspotServers retrieves all hotspot server names
func (c *Client) GetHotspotServers(ctx context.Context, router *entity.Router) ([]string, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/print")
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
