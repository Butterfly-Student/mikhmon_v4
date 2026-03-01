package mikrotik

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// GetAllQueues retrieves all simple queue names.
func (c *Client) GetAllQueues(ctx context.Context) ([]string, error) {
	reply, err := c.RunContext(ctx, "/queue/simple/print")
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

// GetAllParentQueues retrieves non-dynamic simple queue names (for parent selection).
func (c *Client) GetAllParentQueues(ctx context.Context) ([]string, error) {
	reply, err := c.RunContext(ctx, "/queue/simple/print", "?dynamic=false")
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

// StartQueueStatsListen starts listening to queue statistics from MikroTik using
// the RouterOS ListenArgsContext API for real-time streaming.
//
// Karena Client menggunakan async mode, Listen dan Run dapat berjalan
// bersamaan pada koneksi yang sama tanpa saling memblokir.
func (c *Client) StartQueueStatsListen(
	ctx context.Context,
	cfg dto.QueueStatsConfig,
	resultChan chan<- dto.QueueStats,
) (func() error, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("queue name is required")
	}

	args := []string{
		"/queue/simple/print",
		"=stats=",
		"=interval=1s",
		fmt.Sprintf("?name=%s", cfg.Name),
	}

	listenReply, err := c.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start queue stats listen: %w", err)
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

				result := ParseQueueStatsSentence(sentence, cfg.Name)
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

	return func() error {
		_, err := listenReply.Cancel()
		return err
	}, nil
}

// SplitSlashValue splits a "in/out" value into (in, out).
// Example: "11519824317/96401664078" -> (11519824317, 96401664078)
func SplitSlashValue(value string) (int64, int64) {
	if value == "" {
		return 0, 0
	}
	parts := strings.Split(value, "/")
	if len(parts) == 2 {
		return parseInt(parts[0]), parseInt(parts[1])
	}
	return parseInt(value), 0
}

// SplitRateValue splits a rate with unit (bps, kbps, Mbps, Gbps).
// Example: "74.3kbps/2.2Mbps" -> (74300, 2200000)
func SplitRateValue(value string) (int64, int64) {
	if value == "" {
		return 0, 0
	}
	parts := strings.Split(value, "/")
	if len(parts) == 2 {
		return ParseRate(parts[0]), ParseRate(parts[1])
	}
	return ParseRate(value), 0
}

// ParseQueueStatsSentence parses a proto.Sentence into dto.QueueStats.
func ParseQueueStatsSentence(sentence *proto.Sentence, name string) dto.QueueStats {
	m := sentence.Map

	bytesIn, bytesOut := SplitSlashValue(m["bytes"])
	packetsIn, packetsOut := SplitSlashValue(m["packets"])
	queuedBytesIn, queuedBytesOut := SplitSlashValue(m["queued-bytes"])
	queuedPacketsIn, queuedPacketsOut := SplitSlashValue(m["queued-packets"])
	droppedIn, droppedOut := SplitSlashValue(m["dropped"])
	rateIn, rateOut := SplitRateValue(m["rate"])
	packetRateIn, packetRateOut := SplitSlashValue(m["packet-rate"])

	return dto.QueueStats{
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
		TotalRate:          ParseRate(m["total-rate"]),
		TotalPacketRate:    parseInt(m["total-packet-rate"]),
	}
}
