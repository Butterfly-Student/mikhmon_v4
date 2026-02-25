package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

// GetHotspotActive retrieves active hotspot sessions
func (c *Client) GetHotspotActive(ctx context.Context, router *entity.Router) ([]*dto.HotspotActive, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/active/print")
	if err != nil {
		return nil, err
	}

	active := make([]*dto.HotspotActive, 0, len(reply.Re))
	for _, re := range reply.Re {
		active = append(active, &dto.HotspotActive{
			ID:               re.Map[".id"],
			Server:           re.Map["server"],
			User:             re.Map["user"],
			Address:          re.Map["address"],
			MACAddress:       re.Map["mac-address"],
			LoginBy:          re.Map["login-by"],
			Uptime:           re.Map["uptime"],
			SessionTimeLeft:  re.Map["session-time-left"],
			IdleTime:         re.Map["idle-time"],
			IdleTimeout:      re.Map["idle-timeout"],
			KeepaliveTimeout: re.Map["keepalive-timeout"],
			BytesIn:          parseInt(re.Map["bytes-in"]),
			BytesOut:         parseInt(re.Map["bytes-out"]),
			PacketsIn:        parseInt(re.Map["packets-in"]),
			PacketsOut:       parseInt(re.Map["packets-out"]),
			Radius:           parseBool(re.Map["radius"]),
		})
	}

	return active, nil
}

// GetHotspotActiveCount retrieves the count of active hotspot sessions
func (c *Client) GetHotspotActiveCount(ctx context.Context, router *entity.Router) (int, error) {
	client, err := c.getClient(router)
	if err != nil {
		return 0, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/active/print", "=count-only=")
	if err != nil {
		return 0, err
	}

	if len(reply.Re) > 0 {
		return int(parseInt(reply.Re[0].Map["ret"])), nil
	}

	return 0, nil
}

// RemoveHotspotActive removes an active hotspot session
func (c *Client) RemoveHotspotActive(ctx context.Context, router *entity.Router, id string) error {
	client, err := c.getClient(router)
	if err != nil {
		return err
	}

	_, err = client.RunContext(ctx, "/ip/hotspot/active/remove", "=.id="+id)
	return err
}

// GetHotspotHosts retrieves hotspot hosts
func (c *Client) GetHotspotHosts(ctx context.Context, router *entity.Router) ([]*dto.HotspotHost, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/ip/hotspot/host/print")
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

// RemoveHotspotHost removes a hotspot host
func (c *Client) RemoveHotspotHost(ctx context.Context, router *entity.Router, id string) error {
	client, err := c.getClient(router)
	if err != nil {
		return err
	}

	_, err = client.RunContext(ctx, "/ip/hotspot/host/remove", "=.id="+id)
	return err
}
