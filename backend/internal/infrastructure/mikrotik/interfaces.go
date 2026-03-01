package mikrotik

import (
	"context"
	"fmt"
	"time"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// GetInterfaces retrieves all network interfaces.
// Field dari MikroTik: .id, name, type, actual-mtu, l2mtu, max-l2mtu, mac-address, running, disabled
// Note: rx/tx stats tidak tersedia di /interface/print biasa, gunakan /interface/monitor-traffic
func (c *Client) GetInterfaces(ctx context.Context) ([]*dto.Interface, error) {
	reply, err := c.RunContext(ctx, "/interface/print")
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

// StartTrafficMonitorListen starts listening to interface traffic from MikroTik using
// the RouterOS ListenArgsContext API for real-time streaming.
// RouterOS /interface/monitor-traffic menghasilkan stream !re sentences berkelanjutan.
//
// Karena Client menggunakan async mode, Listen dan Run dapat berjalan
// bersamaan pada koneksi yang sama tanpa saling memblokir.
func (c *Client) StartTrafficMonitorListen(
	ctx context.Context,
	name string,
	resultChan chan<- dto.TrafficMonitorStats,
) (func() error, error) {
	if name == "" {
		return nil, fmt.Errorf("interface name is required")
	}

	args := []string{
		"/interface/monitor-traffic",
		fmt.Sprintf("=interface=%s", name),
	}

	listenReply, err := c.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start traffic monitor listen: %w", err)
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

				result := parseTrafficMonitorSentence(sentence, name)
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

// parseTrafficMonitorSentence parses a proto.Sentence into dto.TrafficMonitorStats.
func parseTrafficMonitorSentence(sentence *proto.Sentence, name string) dto.TrafficMonitorStats {
	m := sentence.Map

	return dto.TrafficMonitorStats{
		Name:                  name,
		RxBitsPerSecond:       ParseRate(m["rx-bits-per-second"]),
		TxBitsPerSecond:       ParseRate(m["tx-bits-per-second"]),
		RxPacketsPerSecond:    parseInt(m["rx-packets-per-second"]),
		TxPacketsPerSecond:    parseInt(m["tx-packets-per-second"]),
		FpRxBitsPerSecond:     ParseRate(m["fp-rx-bits-per-second"]),
		FpTxBitsPerSecond:     ParseRate(m["fp-tx-bits-per-second"]),
		FpRxPacketsPerSecond:  parseInt(m["fp-rx-packets-per-second"]),
		FpTxPacketsPerSecond:  parseInt(m["fp-tx-packets-per-second"]),
		RxDropsPerSecond:      parseInt(m["rx-drops-per-second"]),
		TxDropsPerSecond:      parseInt(m["tx-drops-per-second"]),
		TxQueueDropsPerSecond: parseInt(m["tx-queue-drops-per-second"]),
		RxErrorsPerSecond:     parseInt(m["rx-errors-per-second"]),
		TxErrorsPerSecond:     parseInt(m["tx-errors-per-second"]),
	}
}
