package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

// GetAddressPools retrieves all IP pool names
func (c *Client) GetAddressPools(ctx context.Context, router *entity.Router) ([]string, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/ip/pool/print")
	if err != nil {
		return nil, err
	}

	pools := make([]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			pools = append(pools, name)
		}
	}

	return pools, nil
}

// GetParentQueues retrieves all simple queue names for parent selection
func (c *Client) GetParentQueues(ctx context.Context, router *entity.Router) ([]string, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/queue/simple/print")
	if err != nil {
		return nil, err
	}

	queues := make([]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			queues = append(queues, name)
		}
	}

	return queues, nil
}

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
