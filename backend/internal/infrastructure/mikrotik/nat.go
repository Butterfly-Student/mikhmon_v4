package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// GetNATRules retrieves firewall NAT rules.
func (c *Client) GetNATRules(ctx context.Context) ([]*dto.NATRule, error) {
	reply, err := c.RunContext(ctx, "/ip/firewall/nat/print")
	if err != nil {
		return nil, err
	}

	rules := make([]*dto.NATRule, 0, len(reply.Re))
	for _, re := range reply.Re {
		rules = append(rules, &dto.NATRule{
			ID:              re.Map[".id"],
			Chain:           re.Map["chain"],
			Action:          re.Map["action"],
			Protocol:        re.Map["protocol"],
			SrcAddress:      re.Map["src-address"],
			DstAddress:      re.Map["dst-address"],
			SrcPort:         re.Map["src-port"],
			DstPort:         re.Map["dst-port"],
			InInterface:     re.Map["in-interface"],
			OutInterface:    re.Map["out-interface"],
			ToAddresses:     re.Map["to-addresses"],
			ToPorts:         re.Map["to-ports"],
			Disabled:        parseBool(re.Map["disabled"]),
			Comment:         re.Map["comment"],
			Dynamic:         parseBool(re.Map["dynamic"]),
			Invalid:         parseBool(re.Map["invalid"]),
			Bytes:           parseInt(re.Map["bytes"]),
			Packets:         parseInt(re.Map["packets"]),
			ConnectionBytes: parseInt(re.Map["connection-bytes"]),
		})
	}

	return rules, nil
}
