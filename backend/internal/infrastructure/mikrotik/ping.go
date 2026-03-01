package mikrotik

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// StartPingListen starts listening to ping results from MikroTik using
// the RouterOS ListenArgsContext API for real-time streaming.
// RouterOS /ping dengan count=0 menghasilkan stream !re sentences berkelanjutan.
//
// Karena Client menggunakan async mode, Listen dan Run dapat berjalan
// bersamaan pada koneksi yang sama tanpa saling memblokir.
func (c *Client) StartPingListen(
	ctx context.Context,
	cfg dto.PingConfig,
	resultChan chan<- dto.PingResult,
) (func() error, error) {
	applyPingDefaults(&cfg)

	args := []string{
		"/ping",
		fmt.Sprintf("=address=%s", cfg.Address),
		fmt.Sprintf("=interval=%s", formatInterval(cfg.Interval)),
		fmt.Sprintf("=count=%d", cfg.Count),
		fmt.Sprintf("=size=%d", cfg.Size),
	}

	listenReply, err := c.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start ping listen: %w", err)
	}

	seq := 0

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

				result := parsePingSentence(sentence, seq, cfg.Address)
				result.Timestamp = time.Now()

				select {
				case resultChan <- result:
					seq++
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

// applyPingDefaults fills in missing values with defaults.
func applyPingDefaults(cfg *dto.PingConfig) {
	if cfg.Interval <= 0 {
		cfg.Interval = 1 * time.Second
	}
	if cfg.Size <= 0 {
		cfg.Size = 64
	}
	if cfg.Count < 0 {
		cfg.Count = 0
	}
}

// formatInterval formats time.Duration to RouterOS format (e.g. "1s", "500ms").
func formatInterval(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%ds", int(d.Seconds()))
}

// parsePingSentence parses a proto.Sentence into dto.PingResult.
func parsePingSentence(sentence *proto.Sentence, seq int, address string) dto.PingResult {
	m := sentence.Map

	result := dto.PingResult{
		Seq:     seq,
		Address: address,
	}

	if received := m["received"]; received != "" && received != "0" {
		result.Received = true
	}

	if size, err := strconv.Atoi(m["size"]); err == nil {
		result.Size = size
	}

	if ttl, err := strconv.Atoi(m["ttl"]); err == nil {
		result.TTL = ttl
	}

	if timeStr := m["time"]; timeStr != "" {
		trimmed := timeStr
		if len(timeStr) > 2 && timeStr[len(timeStr)-2:] == "ms" {
			trimmed = timeStr[:len(timeStr)-2]
		}
		if t, err := strconv.ParseFloat(trimmed, 64); err == nil {
			result.TimeMs = t
		}
	}

	return result
}
