package dto

import "time"

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

// TrafficStats represents interface traffic statistics snapshot
// Mapping dari: /interface/monitor-traffic (single shot)
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

// TrafficMonitorStats represents real-time interface traffic statistics from streaming
// Mapping dari: /interface/monitor-traffic (streaming / ListenArgsContext)
// Format: rx-bits-per-second=6.0kbps (dengan unit), rx-packets-per-second=10 (plain number)
type TrafficMonitorStats struct {
	Name                  string    `json:"name"`
	RxBitsPerSecond       int64     `json:"rxBitsPerSecond"`       // parsed from rx-bits-per-second (dengan unit)
	TxBitsPerSecond       int64     `json:"txBitsPerSecond"`       // parsed from tx-bits-per-second (dengan unit)
	RxPacketsPerSecond    int64     `json:"rxPacketsPerSecond"`    // rx-packets-per-second
	TxPacketsPerSecond    int64     `json:"txPacketsPerSecond"`    // tx-packets-per-second
	FpRxBitsPerSecond     int64     `json:"fpRxBitsPerSecond"`     // fp-rx-bits-per-second (dengan unit)
	FpTxBitsPerSecond     int64     `json:"fpTxBitsPerSecond"`     // fp-tx-bits-per-second (dengan unit)
	FpRxPacketsPerSecond  int64     `json:"fpRxPacketsPerSecond"`  // fp-rx-packets-per-second
	FpTxPacketsPerSecond  int64     `json:"fpTxPacketsPerSecond"`  // fp-tx-packets-per-second
	RxDropsPerSecond      int64     `json:"rxDropsPerSecond"`      // rx-drops-per-second
	TxDropsPerSecond      int64     `json:"txDropsPerSecond"`      // tx-drops-per-second
	TxQueueDropsPerSecond int64     `json:"txQueueDropsPerSecond"` // tx-queue-drops-per-second
	RxErrorsPerSecond     int64     `json:"rxErrorsPerSecond"`     // rx-errors-per-second
	TxErrorsPerSecond     int64     `json:"txErrorsPerSecond"`     // tx-errors-per-second
	Timestamp             time.Time `json:"timestamp"`
}
