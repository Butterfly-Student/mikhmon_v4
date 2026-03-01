package dto

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
