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