package mikrotik

import (
	"context"
)

// GetAddressPools retrieves all IP pool names.
func (c *Client) GetAddressPools(ctx context.Context) ([]string, error) {
	reply, err := c.RunContext(ctx, "/ip/pool/print")
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
