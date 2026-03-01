package dto


// SystemResource represents MikroTik system resource
// Mapping dari: /system/resource/print
type SystemResource struct {
	Uptime               string  `json:"uptime,omitempty"`
	Version              string  `json:"version,omitempty"`
	BuildTime            string  `json:"buildTime,omitempty"`
	FreeMemory           int64   `json:"freeMemory,omitempty"`
	TotalMemory          int64   `json:"totalMemory,omitempty"`
	FreeHddSpace         int64   `json:"freeHddSpace,omitempty"`
	TotalHddSpace        int64   `json:"totalHddSpace,omitempty"`
	WriteSectSinceReboot int64   `json:"writeSectSinceReboot,omitempty"`
	WriteSectTotal       int64   `json:"writeSectTotal,omitempty"`
	BadBlocks            float64 `json:"badBlocks,omitempty"`
	ArchitectureName     string  `json:"architectureName,omitempty"`
	BoardName            string  `json:"boardName,omitempty"`
	Platform             string  `json:"platform,omitempty"`
	Cpu                  string  `json:"cpu,omitempty"`
	CpuCount             int     `json:"cpuCount,omitempty"`
	CpuFrequency         int     `json:"cpuFrequency,omitempty"`
	CpuLoad              int     `json:"cpuLoad,omitempty"`
}

// SystemHealth represents MikroTik system health
// Mapping dari: /system/health/print
type SystemHealth struct {
	Voltage     string `json:"voltage,omitempty"`
	Temperature string `json:"temperature,omitempty"`
	FanSpeed    string `json:"fanSpeed,omitempty"`
	FanSpeed2   string `json:"fanSpeed2,omitempty"`
	FanSpeed3   string `json:"fanSpeed3,omitempty"`
}

// SystemIdentity represents MikroTik system identity
// Mapping dari: /system/identity/print
type SystemIdentity struct {
	Name string `json:"name,omitempty"`
}

// SystemClock represents MikroTik system clock
// Mapping dari: /system/clock/print
type SystemClock struct {
	Time         string `json:"time,omitempty"`
	Date         string `json:"date,omitempty"`
	TimeZoneName string `json:"timeZoneName,omitempty"`
	TimeZoneAuto string `json:"timeZoneAuto,omitempty"`
	DSTActive    string `json:"dstActive,omitempty"`
}

// RouterBoardInfo represents routerboard information
// Mapping dari: /system/routerboard/print
type RouterBoardInfo struct {
	RouterBoard     string `json:"routerboard,omitempty"`
	Model           string `json:"model,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	FirmwareType    string `json:"firmwareType,omitempty"`
	FactoryFirmware string `json:"factoryFirmware,omitempty"`
	CurrentFirmware string `json:"currentFirmware,omitempty"`
	UpgradeFirmware string `json:"upgradeFirmware,omitempty"`
}

// SystemResourceMonitorStats berisi semua field dari /system/resource/print interval=1s.
// Semua field diperbarui setiap detik secara real-time.
//
// Contoh output MikroTik:
//
//	uptime: 2d1h31m46s
//	version: 6.49.11 (stable)
//	build-time: Dec/08/2023 14:37:03
//	free-memory: 6.2MiB
//	total-memory: 32.0MiB
//	cpu: MIPS 24Kc V7.4
//	cpu-count: 1
//	cpu-frequency: 680MHz
//	cpu-load: 5%
//	free-hdd-space: 46.9MiB
//	total-hdd-space: 63.8MiB
//	write-sect-since-reboot: 11279
//	write-sect-total: 20782897
//	bad-blocks: 0%
//	architecture-name: mipsbe
//	board-name: RB750G
//	platform: MikroTik
type SystemResourceMonitorStats struct {
	Uptime               string  `json:"uptime"`
	Version              string  `json:"version"`
	BuildTime            string  `json:"buildTime"`
	FreeMemory           int64   `json:"freeMemory"`           // bytes
	TotalMemory          int64   `json:"totalMemory"`          // bytes
	CPU                  string  `json:"cpu"`
	CPUCount             int     `json:"cpuCount"`
	CPUFrequency         int     `json:"cpuFrequency"`         // MHz
	CPULoad              int     `json:"cpuLoad"`              // percentage
	FreeHddSpace         int64   `json:"freeHddSpace"`         // bytes
	TotalHddSpace        int64   `json:"totalHddSpace"`        // bytes
	WriteSectSinceReboot int64   `json:"writeSectSinceReboot"`
	WriteSectTotal       int64   `json:"writeSectTotal"`
	BadBlocks            float64 `json:"badBlocks"`            // percentage
	ArchitectureName     string  `json:"architectureName"`
	BoardName            string  `json:"boardName"`
	Platform             string  `json:"platform"`

}
