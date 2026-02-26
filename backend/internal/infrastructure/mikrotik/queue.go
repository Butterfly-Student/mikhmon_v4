package mikrotik

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
)


// QueueStats represents real-time queue statistics
// Mapping dari: /queue/simple/print stats
// Format MikroTik: field=value1/value2 untuk in/out, total-* tanpa split
type QueueStats struct {
	Name        string    `json:"name"`
	BytesIn     int64     `json:"bytesIn"`
	BytesOut    int64     `json:"bytesOut"`
	PacketsIn   int64     `json:"packetsIn"`
	PacketsOut  int64     `json:"packetsOut"`
	QueuedBytesIn   int64 `json:"queuedBytesIn"`   // queued-bytes in
	QueuedBytesOut  int64 `json:"queuedBytesOut"`  // queued-bytes out
	QueuedPacketsIn int64 `json:"queuedPacketsIn"` // queued-packets in
	QueuedPacketsOut int64 `json:"queuedPacketsOut"` // queued-packets out
	DroppedIn   int64     `json:"droppedIn"`
	DroppedOut  int64     `json:"droppedOut"`
	RateIn      int64     `json:"rateIn"`      // bits per second
	RateOut     int64     `json:"rateOut"`     // bits per second
	PacketRateIn  int64   `json:"packetRateIn"`  // packets per second
	PacketRateOut int64   `json:"packetRateOut"` // packets per second
	// Total fields (tanpa split)
	TotalBytes        int64 `json:"totalBytes"`
	TotalPackets      int64 `json:"totalPackets"`
	TotalQueuedBytes  int64 `json:"totalQueuedBytes"`
	TotalQueuedPackets int64 `json:"totalQueuedPackets"`
	TotalDropped      int64 `json:"totalDropped"`
	TotalRate         int64 `json:"totalRate"`        // bits per second
	TotalPacketRate   int64 `json:"totalPacketRate"`  // packets per second
	Timestamp         time.Time `json:"timestamp"`
}

// GetParentQueues retrieves all tree queue names for parent selection
// MIKHMON original uses /queue/tree/print (not /queue/simple/print)
func (c *Client) GetAllQueues(ctx context.Context, router *entity.Router) ([]string, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/queue/simple/print")
	if err != nil {
		return nil, err
	}

	queues := make([]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			queues = append(queues, name)
		}
	}

	return queues, nil
}

func (c *Client) GetAllParentQueues(ctx context.Context, router *entity.Router) ([]string, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(
		ctx,
		"/queue/simple/print",
		"?dynamic=false",
	)
	if err != nil {
		return nil, err
	}

	queues := make([]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			queues = append(queues, name)
		}
	}

	return queues, nil
}

// ==================== Queue Stats Monitor ====================

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



// StartQueueStatsListen starts listening to queue statistics from MikroTik using
// the RouterOS ListenArgsContext API for real-time streaming.
// RouterOS /queue/simple/print stats menghasilkan stream !re sentences berkelanjutan.
//
// PENTING: Queue stats menggunakan koneksi DEDICATED (bukan dari pool) karena streaming
// command memblokir koneksi — jika menggunakan pooled connection, health check
// dari goroutine lain akan conflict dengan stream yang sedang berjalan.
func (c *Client) StartQueueStatsListen(
	ctx context.Context,
	router *entity.Router,
	cfg QueueStatsConfig,
	resultChan chan<- QueueStats,
) (func() error, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("queue name is required")
	}

	// Dial koneksi BARU yang dedicated — tidak dari pool.
	client, err := c.dial(router)
	if err != nil {
		return nil, fmt.Errorf("failed to connect for queue stats: %w", err)
	}

	// Build command: /queue/simple/print stats
	// ListenArgsContext menerima []string (bukan variadic)
	args := []string{
		"/queue/simple/print",
		"=stats=",
		"=interval=1s",
		fmt.Sprintf("?name=%s", cfg.Name),
	}

	// Start listening menggunakan ListenArgsContext
	listenReply, err := client.ListenArgsContext(ctx, args)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to start queue stats listen: %w", err)
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
				result := parseQueueStatsSentence(sentence, cfg.Name)
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

// splitSlashValue memisahkan value dengan format "in/out" menjadi (in, out)
// Contoh: "11519824317/96401664078" -> (11519824317, 96401664078)
// Jika tidak ada slash, mengembalikan (value, 0)
func splitSlashValue(value string) (int64, int64) {
	if value == "" {
		return 0, 0
	}
	parts := strings.Split(value, "/")
	if len(parts) == 2 {
		return parseInt(parts[0]), parseInt(parts[1])
	}
	return parseInt(value), 0
}

// splitRateValue memisahkan rate dengan unit (bps, kbps, Mbps, Gbps)
// Contoh: "74.3kbps/2.2Mbps" -> (74300, 2200000)
func splitRateValue(value string) (int64, int64) {
	if value == "" {
		return 0, 0
	}
	parts := strings.Split(value, "/")
	if len(parts) == 2 {
		return parseRate(parts[0]), parseRate(parts[1])
	}
	return parseRate(value), 0
}

// parseQueueStatsSentence parses a proto.Sentence into QueueStats.
// Format MikroTik: bytes=11519824317/96401664078 (in/out), total-bytes=0 (tanpa split)
func parseQueueStatsSentence(sentence *proto.Sentence, name string) QueueStats {
	m := sentence.Map

	bytesIn, bytesOut := splitSlashValue(m["bytes"])
	packetsIn, packetsOut := splitSlashValue(m["packets"])
	queuedBytesIn, queuedBytesOut := splitSlashValue(m["queued-bytes"])
	queuedPacketsIn, queuedPacketsOut := splitSlashValue(m["queued-packets"])
	droppedIn, droppedOut := splitSlashValue(m["dropped"])
	rateIn, rateOut := splitRateValue(m["rate"])
	packetRateIn, packetRateOut := splitSlashValue(m["packet-rate"])

	return QueueStats{
		Name:               name,
		BytesIn:            bytesIn,
		BytesOut:           bytesOut,
		PacketsIn:          packetsIn,
		PacketsOut:         packetsOut,
		QueuedBytesIn:      queuedBytesIn,
		QueuedBytesOut:     queuedBytesOut,
		QueuedPacketsIn:    queuedPacketsIn,
		QueuedPacketsOut:   queuedPacketsOut,
		DroppedIn:          droppedIn,
		DroppedOut:         droppedOut,
		RateIn:             rateIn,
		RateOut:            rateOut,
		PacketRateIn:       packetRateIn,
		PacketRateOut:      packetRateOut,
		TotalBytes:         parseInt(m["total-bytes"]),
		TotalPackets:       parseInt(m["total-packets"]),
		TotalQueuedBytes:   parseInt(m["total-queued-bytes"]),
		TotalQueuedPackets: parseInt(m["total-queued-packets"]),
		TotalDropped:       parseInt(m["total-dropped"]),
		TotalRate:          parseRate(m["total-rate"]),
		TotalPacketRate:    parseInt(m["total-packet-rate"]),
	}
}