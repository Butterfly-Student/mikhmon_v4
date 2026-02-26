package mikrotik

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-routeros/routeros/v3/proto"
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


// ==================== System Resource Monitor ====================



// SystemResourceMonitorStats represents real-time system resource statistics
// Mapping dari: /system/resource/monitor
// Output contoh:
//   cpu-used: 4%
//   free-memory: 6756KiB
type SystemResourceMonitorStats struct {
	CPUUsed         int       `json:"cpuUsed"`         // CPU usage percentage
	FreeMemory      int64     `json:"freeMemory"`      // Free memory in bytes
	TotalMemory     int64     `json:"totalMemory"`     // Total memory in bytes
	FreeHddSpace    int64     `json:"freeHddSpace"`    // Free HDD space in bytes
	TotalHddSpace   int64     `json:"totalHddSpace"`   // Total HDD space in bytes
	WriteSectSinceReboot int64 `json:"writeSectSinceReboot"` // Write sectors since reboot
	Uptime          string    `json:"uptime"`          // System uptime
	Timestamp       time.Time `json:"timestamp"`
}

// StartSystemResourceMonitorListen starts listening to system resource statistics from MikroTik using
// the RouterOS ListenArgsContext API for real-time streaming.
// RouterOS /system/resource/monitor menghasilkan stream !re sentences berkelanjutan.
//
// PENTING: Resource monitor menggunakan koneksi DEDICATED (bukan dari pool) karena streaming
// command memblokir koneksi — jika menggunakan pooled connection, health check
// dari goroutine lain akan conflict dengan stream yang sedang berjalan.
func (c *Client) StartSystemResourceMonitorListen(
	ctx context.Context,
	router *entity.Router,
	resultChan chan<- SystemResourceMonitorStats,
) (func() error, error) {

	// Dial koneksi BARU yang dedicated — tidak dari pool.
	client, err := c.dial(router)
	if err != nil {
		return nil, fmt.Errorf("failed to connect for resource monitor: %w", err)
	}

	// Build command: /system/resource/monitor
	// ListenArgsContext menerima []string (bukan variadic)
	args := []string{
		"/system/resource/monitor",
	}

	// Start listening menggunakan ListenArgsContext
	listenReply, err := client.ListenArgsContext(ctx, args)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to start resource monitor listen: %w", err)
	}

	// Process replies in a goroutine
	go func() {
		defer close(resultChan)
		defer client.Close() // Tutup koneksi dedicated ketika selesai

		for {
			select {
			case <-ctx.Done():
				// Context cancelled, cancel the RouterOS command
				listenReply.Cancel()
				return

			case sentence, ok := <-listenReply.Chan():
				if !ok {
					// Channel closed (done or cancelled)
					return
				}

				// Parse the sentence
				result := parseSystemResourceMonitorSentence(sentence)
				result.Timestamp = time.Now()

				select {
				case resultChan <- result:
				case <-ctx.Done():
					listenReply.Cancel()
					return
				}
			}
		}
	}()

	// Return cancel function
	return func() error {
		_, err := listenReply.Cancel()
		return err
	}, nil
}

// parseSystemResourceMonitorSentence parses a proto.Sentence into SystemResourceMonitorStats.
// RouterOS mengembalikan nilai dengan format seperti "4%" untuk cpu-used dan "6756KiB" untuk memory.
func parseSystemResourceMonitorSentence(sentence *proto.Sentence) SystemResourceMonitorStats {
	m := sentence.Map

	return SystemResourceMonitorStats{
		CPUUsed:              parsePercentage(m["cpu-used"]),
		FreeMemory:           parseByteSize(m["free-memory"]),
		TotalMemory:          parseByteSize(m["total-memory"]),
		FreeHddSpace:         parseByteSize(m["free-hdd-space"]),
		TotalHddSpace:        parseByteSize(m["total-hdd-space"]),
		WriteSectSinceReboot: parseInt(m["write-sect-since-reboot"]),
		Uptime:               m["uptime"],
	}
}

// parsePercentage parses percentage string like "4%" to int
func parsePercentage(s string) int {
	if s == "" {
		return 0
	}
	// Remove % suffix
	if len(s) > 1 && s[len(s)-1] == '%' {
		s = s[:len(s)-1]
	}
	i, _ := strconv.Atoi(s)
	return i
}

// parseByteSize parses byte size strings like "6756KiB", "10MiB", "2GiB" to int64
func parseByteSize(s string) int64 {
	if s == "" {
		return 0
	}

	// Map of suffixes to multipliers
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
	}

	// Try to find and remove suffix
	for suffix, multiplier := range multipliers {
		if len(s) > len(suffix) && s[len(s)-len(suffix):] == suffix {
			numStr := s[:len(s)-len(suffix)]
			if val, err := strconv.ParseInt(numStr, 10, 64); err == nil {
				return val * multiplier
			}
			return 0
		}
	}

	// No suffix found, try parsing as plain number
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}
