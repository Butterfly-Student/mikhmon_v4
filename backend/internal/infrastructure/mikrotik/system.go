package mikrotik

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// GetSystemResource retrieves system resource information.
func (c *Client) GetSystemResource(ctx context.Context) (*dto.SystemResource, error) {
	reply, err := c.RunContext(ctx, "/system/resource/print")
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
		FreeMemory:           parseByteSize(re.Map["free-memory"]),
		TotalMemory:          parseByteSize(re.Map["total-memory"]),
		FreeHddSpace:         parseByteSize(re.Map["free-hdd-space"]),
		TotalHddSpace:        parseByteSize(re.Map["total-hdd-space"]),
		WriteSectSinceReboot: parseInt(re.Map["write-sect-since-reboot"]),
		WriteSectTotal:       parseInt(re.Map["write-sect-total"]),
		BadBlocks:            parsePercentageFloat(re.Map["bad-blocks"]),
		ArchitectureName:     re.Map["architecture-name"],
		BoardName:            re.Map["board-name"],
		Platform:             re.Map["platform"],
		Cpu:                  re.Map["cpu"],
		CpuCount:             int(parseInt(re.Map["cpu-count"])),
		CpuFrequency:         int(parseInt(re.Map["cpu-frequency"])),
		CpuLoad:              parsePercentage(re.Map["cpu-load"]),
	}, nil
}

// GetSystemHealth retrieves system health information.
func (c *Client) GetSystemHealth(ctx context.Context) (*dto.SystemHealth, error) {
	reply, err := c.RunContext(ctx, "/system/health/print")
	if err != nil {
		// Some routers might not have health monitoring.
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

// GetSystemIdentity retrieves system identity.
func (c *Client) GetSystemIdentity(ctx context.Context) (*dto.SystemIdentity, error) {
	reply, err := c.RunContext(ctx, "/system/identity/print")
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

// GetSystemClock retrieves system clock information.
func (c *Client) GetSystemClock(ctx context.Context) (*dto.SystemClock, error) {
	reply, err := c.RunContext(ctx, "/system/clock/print")
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

// GetRouterBoardInfo retrieves routerboard information.
func (c *Client) GetRouterBoardInfo(ctx context.Context) (*dto.RouterBoardInfo, error) {
	reply, err := c.RunContext(ctx, "/system/routerboard/print")
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

// ==================== System Resource Monitor ====================

// StartSystemResourceMonitorListen starts streaming system resource statistics.
// Menggunakan /system/resource/print interval=1s sehingga semua field diperbarui
// setiap detik secara real-time dalam satu koneksi tanpa fetch statis terpisah.
//
// Karena Client menggunakan async mode, Listen dan Run dapat berjalan
// bersamaan pada koneksi yang sama tanpa saling memblokir.
func (c *Client) StartSystemResourceMonitorListen(
	ctx context.Context,
	resultChan chan<- dto.SystemResourceMonitorStats,
) (func() error, error) {
	listenReply, err := c.ListenArgsContext(ctx, []string{
		"/system/resource/print",
		"=interval=1s",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start resource monitor listen: %w", err)
	}

	go func() {
		defer close(resultChan)

		for {
			select {
			case <-ctx.Done():
				listenReply.Cancel()
				return

			case sentence, ok := <-listenReply.Chan():
				if !ok {
					return
				}

				result := parseSystemResourcePrintSentence(sentence)

				select {
				case resultChan <- result:
				case <-ctx.Done():
					listenReply.Cancel()
					return
				}
			}
		}
	}()

	return func() error {
		_, err := listenReply.Cancel()
		return err
	}, nil
}

// parseSystemResourcePrintSentence parses a proto.Sentence from /system/resource/print
// into dto.SystemResourceMonitorStats. Semua field tersedia karena print interval=1s
// mengembalikan data lengkap setiap iterasi.
func parseSystemResourcePrintSentence(sentence *proto.Sentence) dto.SystemResourceMonitorStats {
	m := sentence.Map

	rawData := make(map[string]string, len(m))
	for k, v := range m {
		rawData[k] = v
	}

	return dto.SystemResourceMonitorStats{
		Uptime:               m["uptime"],
		Version:              m["version"],
		BuildTime:            m["build-time"],
		FreeMemory:           parseByteSize(m["free-memory"]),
		TotalMemory:          parseByteSize(m["total-memory"]),
		CPU:                  m["cpu"],
		CPUCount:             int(parseInt(m["cpu-count"])),
		CPUFrequency:         int(parseByteSize(m["cpu-frequency"])), // e.g. "680MHz" → 680
		CPULoad:              parsePercentage(m["cpu-load"]),
		FreeHddSpace:         parseByteSize(m["free-hdd-space"]),
		TotalHddSpace:        parseByteSize(m["total-hdd-space"]),
		WriteSectSinceReboot: parseInt(m["write-sect-since-reboot"]),
		WriteSectTotal:       parseInt(m["write-sect-total"]),
		BadBlocks:            parsePercentageFloat(m["bad-blocks"]),
		ArchitectureName:     m["architecture-name"],
		BoardName:            m["board-name"],
		Platform:             m["platform"],
	}
}

// parsePercentage parses a percentage string like "4%" to int.
func parsePercentage(s string) int {
	if s == "" {
		return 0
	}
	if len(s) > 1 && s[len(s)-1] == '%' {
		s = s[:len(s)-1]
	}
	i, _ := strconv.Atoi(s)
	return i
}

// parsePercentageFloat parses a percentage string like "0%" or "13.5%" to float64.
func parsePercentageFloat(s string) float64 {
	if s == "" {
		return 0
	}
	if len(s) > 1 && s[len(s)-1] == '%' {
		s = s[:len(s)-1]
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// parseByteSize parses byte size strings like "6.4MiB", "6756KiB", "10MiB", "2GiB" to int64.
func parseByteSize(s string) int64 {
	if s == "" {
		return 0
	}

	multipliers := map[string]int64{
		"KiB": 1024,
		"MiB": 1024 * 1024,
		"GiB": 1024 * 1024 * 1024,
		"TiB": 1024 * 1024 * 1024 * 1024,
		"KB":  1024,
		"MB":  1024 * 1024,
		"GB":  1024 * 1024 * 1024,
		"TB":  1024 * 1024 * 1024 * 1024,
		"B":   1,
		"MHz": 1,
		"GHz": 1000,
	}

	for suffix, multiplier := range multipliers {
		if len(s) > len(suffix) && s[len(s)-len(suffix):] == suffix {
			numStr := s[:len(s)-len(suffix)]
			if val, err := strconv.ParseFloat(numStr, 64); err == nil {
				return int64(val * float64(multiplier))
			}
			return 0
		}
	}

	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}
