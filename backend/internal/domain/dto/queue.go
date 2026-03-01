package dto

import "time"

// QueueStats represents real-time queue statistics
// Mapping dari: /queue/simple/print stats (streaming)
// Format MikroTik: field=value1/value2 untuk in/out, total-* tanpa split
type QueueStats struct {
	Name             string `json:"name"`
	BytesIn          int64  `json:"bytesIn"`
	BytesOut         int64  `json:"bytesOut"`
	PacketsIn        int64  `json:"packetsIn"`
	PacketsOut       int64  `json:"packetsOut"`
	QueuedBytesIn    int64  `json:"queuedBytesIn"`    // queued-bytes in
	QueuedBytesOut   int64  `json:"queuedBytesOut"`   // queued-bytes out
	QueuedPacketsIn  int64  `json:"queuedPacketsIn"`  // queued-packets in
	QueuedPacketsOut int64  `json:"queuedPacketsOut"` // queued-packets out
	DroppedIn        int64  `json:"droppedIn"`
	DroppedOut       int64  `json:"droppedOut"`
	RateIn           int64  `json:"rateIn"`        // bits per second
	RateOut          int64  `json:"rateOut"`       // bits per second
	PacketRateIn     int64  `json:"packetRateIn"`  // packets per second
	PacketRateOut    int64  `json:"packetRateOut"` // packets per second
	// Total fields (tanpa split)
	TotalBytes         int64     `json:"totalBytes"`
	TotalPackets       int64     `json:"totalPackets"`
	TotalQueuedBytes   int64     `json:"totalQueuedBytes"`
	TotalQueuedPackets int64     `json:"totalQueuedPackets"`
	TotalDropped       int64     `json:"totalDropped"`
	TotalRate          int64     `json:"totalRate"`       // bits per second
	TotalPacketRate    int64     `json:"totalPacketRate"` // packets per second
	Timestamp          time.Time `json:"timestamp"`
}

// QueueStatsConfig holds configuration for queue stats monitoring
type QueueStatsConfig struct {
	Name string // Queue name (required)
}

// DefaultQueueStatsConfig returns QueueStatsConfig with sensible defaults
func DefaultQueueStatsConfig(name string) QueueStatsConfig {
	return QueueStatsConfig{
		Name: name,
	}
}
