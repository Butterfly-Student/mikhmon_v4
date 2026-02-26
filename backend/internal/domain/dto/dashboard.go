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

// Interface represents a network interface
// Mapping dari: /interface/print
type Interface struct {
	ID         string `json:".id,omitempty"`
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	MTU        int    `json:"mtu,omitempty"`
	MacAddress string `json:"macAddress,omitempty"`
	Running    bool   `json:"running,omitempty"`
	Disabled   bool   `json:"disabled,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

// TrafficStats represents interface traffic statistics
// Mapping dari: /interface/monitor-traffic
type TrafficStats struct {
	Name                  string `json:"name,omitempty"`
	RxBitsPerSecond       int64  `json:"rxBitsPerSecond,omitempty"`
	TxBitsPerSecond       int64  `json:"txBitsPerSecond,omitempty"`
	RxPacketsPerSecond    int64  `json:"rxPacketsPerSecond,omitempty"`
	TxPacketsPerSecond    int64  `json:"txPacketsPerSecond,omitempty"`
	FpRxBitsPerSecond     int64  `json:"fpRxBitsPerSecond,omitempty"`
	FpTxBitsPerSecond     int64  `json:"fpTxBitsPerSecond,omitempty"`
	FpRxPacketsPerSecond  int64  `json:"fpRxPacketsPerSecond,omitempty"`
	FpTxPacketsPerSecond  int64  `json:"fpTxPacketsPerSecond,omitempty"`
	RxDropsPerSecond      int64  `json:"rxDropsPerSecond,omitempty"`
	TxDropsPerSecond      int64  `json:"txDropsPerSecond,omitempty"`
	TxQueueDropsPerSecond int64  `json:"txQueueDropsPerSecond,omitempty"`
	RxErrorsPerSecond     int64  `json:"rxErrorsPerSecond,omitempty"`
	TxErrorsPerSecond     int64  `json:"txErrorsPerSecond,omitempty"`
}

// NATRule represents a firewall NAT rule
// Mapping from: /ip/firewall/nat/print
type NATRule struct {
	ID              string `json:"id,omitempty"`
	Chain           string `json:"chain,omitempty"`
	Action          string `json:"action,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	SrcAddress      string `json:"srcAddress,omitempty"`
	DstAddress      string `json:"dstAddress,omitempty"`
	SrcPort         string `json:"srcPort,omitempty"`
	DstPort         string `json:"dstPort,omitempty"`
	InInterface     string `json:"inInterface,omitempty"`
	OutInterface    string `json:"outInterface,omitempty"`
	ToAddresses     string `json:"toAddresses,omitempty"`
	ToPorts         string `json:"toPorts,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Dynamic         bool   `json:"dynamic,omitempty"`
	Invalid         bool   `json:"invalid,omitempty"`
	Bytes           int64  `json:"bytes,omitempty"`
	Packets         int64  `json:"packets,omitempty"`
	ConnectionBytes int64  `json:"connectionBytes,omitempty"`
}

// LogEntry represents a log entry
// Mapping dari: /log/print
type LogEntry struct {
	ID      string `json:".id,omitempty"`
	Time    string `json:"time,omitempty"`
	Topics  string `json:"topics,omitempty"`
	Message string `json:"message,omitempty"`
}

// HotspotStats represents hotspot statistics
type HotspotStats struct {
	TotalUsers  int `json:"totalUsers"`
	ActiveUsers int `json:"activeUsers"`
}

// DashboardData represents complete dashboard data
type DashboardData struct {
	RouterID        uint             `json:"routerId"`
	RouterName      string           `json:"routerName,omitempty"`
	SystemTime      *SystemClock     `json:"systemTime,omitempty"`
	Resource        *SystemResource  `json:"resource,omitempty"`
	Health          *SystemHealth    `json:"health,omitempty"`
	Identity        *SystemIdentity  `json:"identity,omitempty"`
	RouterBoard     *RouterBoardInfo `json:"routerBoard,omitempty"`
	Stats           *HotspotStats    `json:"stats,omitempty"`
	Interfaces      []*Interface     `json:"interfaces,omitempty"`
	HotspotLogs     []*LogEntry      `json:"hotspotLogs,omitempty"`
	ConnectionError string           `json:"connectionError,omitempty"`
}

// TrafficRequest represents a traffic monitoring request
type TrafficRequest struct {
	Interface string `json:"interface" validate:"required"`
}
