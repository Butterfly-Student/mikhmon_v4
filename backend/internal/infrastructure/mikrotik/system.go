package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

// GetSystemResource retrieves system resource information
func (c *Client) GetSystemResource(ctx context.Context, router *entity.Router) (*dto.SystemResource, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/system/resource/print")
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return &dto.SystemResource{}, nil
	}

	re := reply.Re[0]
	return &dto.SystemResource{
		Uptime:               re.Map["uptime"],
		Version:              re.Map["version"],
		BuildTime:            re.Map["build-time"],
		FreeMemory:           parseInt(re.Map["free-memory"]),
		TotalMemory:          parseInt(re.Map["total-memory"]),
		FreeHddSpace:         parseInt(re.Map["free-hdd-space"]),
		TotalHddSpace:        parseInt(re.Map["total-hdd-space"]),
		WriteSectSinceReboot: parseInt(re.Map["write-sect-since-reboot"]),
		WriteSectTotal:       parseInt(re.Map["write-sect-total"]),
		BadBlocks:            parseFloat(re.Map["bad-blocks"]),
		ArchitectureName:     re.Map["architecture-name"],
		BoardName:            re.Map["board-name"],
		Cpu:                  re.Map["cpu"],
		CpuCount:             int(parseInt(re.Map["cpu-count"])),
		CpuFrequency:         int(parseInt(re.Map["cpu-frequency"])),
		CpuLoad:              int(parseInt(re.Map["cpu-load"])),
	}, nil
}

// GetSystemHealth retrieves system health information
func (c *Client) GetSystemHealth(ctx context.Context, router *entity.Router) (*dto.SystemHealth, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/system/health/print")
	if err != nil {
		// Some routers might not have health monitoring
		return &dto.SystemHealth{}, nil
	}

	if len(reply.Re) == 0 {
		return &dto.SystemHealth{}, nil
	}

	re := reply.Re[0]
	return &dto.SystemHealth{
		Voltage:     re.Map["voltage"],
		Temperature: re.Map["temperature"],
		FanSpeed:    re.Map["fan-speed"],
		FanSpeed2:   re.Map["fan-speed2"],
		FanSpeed3:   re.Map["fan-speed3"],
	}, nil
}

// GetSystemIdentity retrieves system identity
func (c *Client) GetSystemIdentity(ctx context.Context, router *entity.Router) (*dto.SystemIdentity, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/system/identity/print")
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return &dto.SystemIdentity{}, nil
	}

	return &dto.SystemIdentity{
		Name: reply.Re[0].Map["name"],
	}, nil
}

// GetSystemClock retrieves system clock information
func (c *Client) GetSystemClock(ctx context.Context, router *entity.Router) (*dto.SystemClock, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/system/clock/print")
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return &dto.SystemClock{}, nil
	}

	re := reply.Re[0]
	return &dto.SystemClock{
		Time:         re.Map["time"],
		Date:         re.Map["date"],
		TimeZoneName: re.Map["time-zone-name"],
		TimeZoneAuto: re.Map["time-zone-autodetect"],
		DSTActive:    re.Map["dst-active"],
	}, nil
}

// GetRouterBoardInfo retrieves routerboard information
func (c *Client) GetRouterBoardInfo(ctx context.Context, router *entity.Router) (*dto.RouterBoardInfo, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/system/routerboard/print")
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return &dto.RouterBoardInfo{}, nil
	}

	re := reply.Re[0]
	return &dto.RouterBoardInfo{
		RouterBoard:     re.Map["routerboard"],
		Model:           re.Map["model"],
		SerialNumber:    re.Map["serial-number"],
		FirmwareType:    re.Map["firmware-type"],
		FactoryFirmware: re.Map["factory-firmware"],
		CurrentFirmware: re.Map["current-firmware"],
		UpgradeFirmware: re.Map["upgrade-firmware"],
	}, nil
}
