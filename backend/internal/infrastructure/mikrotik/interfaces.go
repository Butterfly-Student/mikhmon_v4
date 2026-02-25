package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

// GetInterfaces retrieves all network interfaces
func (c *Client) GetInterfaces(ctx context.Context, router *entity.Router) ([]*dto.Interface, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/interface/print")
	if err != nil {
		return nil, err
	}

	interfaces := make([]*dto.Interface, 0, len(reply.Re))
	for _, re := range reply.Re {
		interfaces = append(interfaces, &dto.Interface{
			ID:         re.Map[".id"],
			Name:       re.Map["name"],
			Type:       re.Map["type"],
			MTU:        int(parseInt(re.Map["mtu"])),
			MacAddress: re.Map["mac-address"],
			Running:    parseBool(re.Map["running"]),
			Disabled:   parseBool(re.Map["disabled"]),
			Comment:    re.Map["comment"],
			RxByte:     parseInt(re.Map["rx-byte"]),
			TxByte:     parseInt(re.Map["tx-byte"]),
			RxPacket:   parseInt(re.Map["rx-packet"]),
			TxPacket:   parseInt(re.Map["tx-packet"]),
			RxDrop:     parseInt(re.Map["rx-drop"]),
			TxDrop:     parseInt(re.Map["tx-drop"]),
			RxError:    parseInt(re.Map["rx-error"]),
			TxError:    parseInt(re.Map["tx-error"]),
		})
	}

	return interfaces, nil
}

// MonitorTraffic monitors traffic on a specific interface
func (c *Client) MonitorTraffic(ctx context.Context, router *entity.Router, iface string) (*dto.TrafficStats, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx,
		"/interface/monitor-traffic",
		"=interface="+iface,
		"=once=",
	)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return &dto.TrafficStats{Name: iface}, nil
	}

	re := reply.Re[0]
	return &dto.TrafficStats{
		Name:               iface,
		TxBitsPerSecond:    parseInt(re.Map["tx-bits-per-second"]),
		RxBitsPerSecond:    parseInt(re.Map["rx-bits-per-second"]),
		TxPacketsPerSecond: parseInt(re.Map["tx-packets-per-second"]),
		RxPacketsPerSecond: parseInt(re.Map["rx-packets-per-second"]),
		TxDropped:          parseInt(re.Map["tx-dropped"]),
		RxDropped:          parseInt(re.Map["rx-dropped"]),
		TxErrors:           parseInt(re.Map["tx-errors"]),
		RxErrors:           parseInt(re.Map["rx-errors"]),
	}, nil
}
