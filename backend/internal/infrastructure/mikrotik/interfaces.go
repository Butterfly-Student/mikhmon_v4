package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

// GetInterfaces retrieves all network interfaces
// Field dari MikroTik: .id, name, type, actual-mtu, l2mtu, max-l2mtu, mac-address, running, disabled
// Note: rx/tx stats tidak tersedia di /interface/print biasa, gunakan /interface/monitor-traffic
func (c *Client) GetInterfaces(ctx context.Context, router *entity.Router) ([]*dto.Interface, error) {
	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/interface/print")
	if err != nil {
		return nil, err
	}

	interfaces := make([]*dto.Interface, 0, len(reply.Re))
	for _, re := range reply.Re {
		interfaces = append(interfaces, &dto.Interface{
			ID:         re.Map[".id"],
			Name:       re.Map["name"],
			Type:       re.Map["type"],
			MTU:        int(parseInt(re.Map["actual-mtu"])),
			MacAddress: re.Map["mac-address"],
			Running:    parseBool(re.Map["running"]),
			Disabled:   parseBool(re.Map["disabled"]),
			Comment:    re.Map["comment"],
		})
	}

	return interfaces, nil
}


// ==================== Interface Traffic Monitor ====================



// TrafficMonitorStats represents real-time interface traffic statistics
// Mapping dari: /interface/monitor-traffic
// Format: rx-bits-per-second=6.0kbps (dengan unit), rx-packets-per-second=10 (plain number)
type TrafficMonitorStats struct {
	Name               string    `json:"name"`
	RxBitsPerSecond    int64     `json:"rxBitsPerSecond"`    // parsed from rx-bits-per-second (dengan unit)
	TxBitsPerSecond    int64     `json:"txBitsPerSecond"`    // parsed from tx-bits-per-second (dengan unit)
	RxPacketsPerSecond int64     `json:"rxPacketsPerSecond"` // rx-packets-per-second
	TxPacketsPerSecond int64     `json:"txPacketsPerSecond"` // tx-packets-per-second
	FpRxBitsPerSecond  int64     `json:"fpRxBitsPerSecond"`  // fp-rx-bits-per-second (dengan unit)
	FpTxBitsPerSecond  int64     `json:"fpTxBitsPerSecond"`  // fp-tx-bits-per-second (dengan unit)
	FpRxPacketsPerSecond int64   `json:"fpRxPacketsPerSecond"` // fp-rx-packets-per-second
	FpTxPacketsPerSecond int64   `json:"fpTxPacketsPerSecond"` // fp-tx-packets-per-second
	RxDropsPerSecond   int64     `json:"rxDropsPerSecond"`   // rx-drops-per-second
	TxDropsPerSecond   int64     `json:"txDropsPerSecond"`   // tx-drops-per-second
	TxQueueDropsPerSecond int64  `json:"txQueueDropsPerSecond"` // tx-queue-drops-per-second
	RxErrorsPerSecond  int64     `json:"rxErrorsPerSecond"`  // rx-errors-per-second
	TxErrorsPerSecond  int64     `json:"txErrorsPerSecond"`  // tx-errors-per-second
	Timestamp          time.Time `json:"timestamp"`
}

// StartTrafficMonitorListen starts listening to interface traffic from MikroTik using
// the RouterOS ListenArgsContext API for real-time streaming.
// RouterOS /interface/monitor-traffic menghasilkan stream !re sentences berkelanjutan.
//
// PENTING: Traffic monitor menggunakan koneksi DEDICATED (bukan dari pool) karena streaming
// command memblokir koneksi — jika menggunakan pooled connection, health check
// dari goroutine lain akan conflict dengan stream yang sedang berjalan.
func (c *Client) StartTrafficMonitorListen(
	ctx context.Context,
	router *entity.Router,
	Name string,
	resultChan chan<- TrafficMonitorStats,
) (func() error, error) {

	if Name == "" {
		return nil, fmt.Errorf("interface name is required")
	}

	// Dial koneksi BARU yang dedicated — tidak dari pool.
	client, err := c.dial(router)
	if err != nil {
		return nil, fmt.Errorf("failed to connect for traffic monitor: %w", err)
	}

	// Build command: /interface/monitor-traffic
	// ListenArgsContext menerima []string (bukan variadic)
	args := []string{
		"/interface/monitor-traffic",
		fmt.Sprintf("=interface=%s", Name),
	}

	// Start listening menggunakan ListenArgsContext
	listenReply, err := client.ListenArgsContext(ctx, args)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to start traffic monitor listen: %w", err)
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
				result := parseTrafficMonitorSentence(sentence, Name)
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

// parseTrafficMonitorSentence parses a proto.Sentence into TrafficMonitorStats.
// Field rate (bps/kbps/Mbps/Gbps) diparse menggunakan parseRate
func parseTrafficMonitorSentence(sentence *proto.Sentence, name string) TrafficMonitorStats {
	m := sentence.Map

	return TrafficMonitorStats{
		Name:                  name,
		RxBitsPerSecond:       parseRate(m["rx-bits-per-second"]),
		TxBitsPerSecond:       parseRate(m["tx-bits-per-second"]),
		RxPacketsPerSecond:    parseInt(m["rx-packets-per-second"]),
		TxPacketsPerSecond:    parseInt(m["tx-packets-per-second"]),
		FpRxBitsPerSecond:     parseRate(m["fp-rx-bits-per-second"]),
		FpTxBitsPerSecond:     parseRate(m["fp-tx-bits-per-second"]),
		FpRxPacketsPerSecond:  parseInt(m["fp-rx-packets-per-second"]),
		FpTxPacketsPerSecond:  parseInt(m["fp-tx-packets-per-second"]),
		RxDropsPerSecond:      parseInt(m["rx-drops-per-second"]),
		TxDropsPerSecond:      parseInt(m["tx-drops-per-second"]),
		TxQueueDropsPerSecond: parseInt(m["tx-queue-drops-per-second"]),
		RxErrorsPerSecond:     parseInt(m["rx-errors-per-second"]),
		TxErrorsPerSecond:     parseInt(m["tx-errors-per-second"]),
	}
}