package dto

import "time"

// PingConfig holds configuration for a ping session via MikroTik /ping command
type PingConfig struct {
	Address  string        // Target IP/hostname (required)
	Interval time.Duration // Ping interval (default: 1s)
	Count    int           // Number of pings: 0 = infinite (default: 0)
	Size     int           // Packet size in bytes (default: 64)
}

// DefaultPingConfig returns PingConfig with sensible defaults
func DefaultPingConfig(address string) PingConfig {
	return PingConfig{
		Address:  address,
		Interval: 1 * time.Second,
		Count:    0,
		Size:     64,
	}
}

// PingResult represents a single ping result from MikroTik /ping streaming
type PingResult struct {
	Seq       int       `json:"seq"`
	Address   string    `json:"address"`
	TimeMs    float64   `json:"timeMs"`
	TTL       int       `json:"ttl"`
	Size      int       `json:"size"`
	Received  bool      `json:"received"`
	Timestamp time.Time `json:"timestamp"`
}
