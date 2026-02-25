package mikrotik

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-routeros/routeros/v3/proto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

// PingConfig holds configuration for a ping session
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

// applyDefaults fills in missing values with defaults
func (cfg *PingConfig) applyDefaults() {
	if cfg.Interval <= 0 {
		cfg.Interval = 1 * time.Second
	}
	if cfg.Size <= 0 {
		cfg.Size = 64
	}
	// Count 0 = infinite, negative tidak diperbolehkan
	if cfg.Count < 0 {
		cfg.Count = 0
	}
}

// PingResult represents a single ping result
type PingResult struct {
	Seq       int       `json:"seq"`
	Address   string    `json:"address"`
	TimeMs    float64   `json:"timeMs"`
	TTL       int       `json:"ttl"`
	Size      int       `json:"size"`
	Received  bool      `json:"received"`
	Timestamp time.Time `json:"timestamp"`
}

// StartPingListen starts listening to ping results from MikroTik using
// the RouterOS ListenArgsContext API for real-time streaming.
// RouterOS /ping dengan count=0 menghasilkan stream !re sentences berkelanjutan.
//
// PENTING: Ping menggunakan koneksi DEDICATED (bukan dari pool) karena streaming
// command memblokir koneksi — jika menggunakan pooled connection, health check
// dari goroutine lain akan conflict dengan stream yang sedang berjalan.
func (c *Client) StartPingListen(
	ctx context.Context,
	router *entity.Router,
	cfg PingConfig,
	resultChan chan<- PingResult,
) (func() error, error) {
	// Apply defaults untuk field yang tidak diset
	cfg.applyDefaults()

	// Dial koneksi BARU yang dedicated — tidak dari pool.
	// Ini mencegah health check dari getClient() yang berjalan concurrent
	// menginterferensi stream ping yang sedang aktif di koneksi yang sama.
	client, err := c.dial(router)
	if err != nil {
		return nil, fmt.Errorf("failed to connect for ping: %w", err)
	}

	// Build ping command
	// ListenArgsContext menerima []string (bukan variadic)
	args := []string{
		"/ping",
		fmt.Sprintf("=address=%s", cfg.Address),
		fmt.Sprintf("=interval=%s", formatInterval(cfg.Interval)),
		fmt.Sprintf("=count=%d", cfg.Count),
		fmt.Sprintf("=size=%d", cfg.Size),
	}

	// Start listening menggunakan ListenArgsContext
	listenReply, err := client.ListenArgsContext(ctx, args)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to start ping listen: %w", err)
	}

	seq := 0

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

	// Return cancel function — gunakan Cancel() biasa agar tidak tergantung ctx yang sudah done
	return func() error {
		_, err := listenReply.Cancel()
		return err
	}, nil
}

// formatInterval formats time.Duration ke format yang diterima RouterOS (contoh: "1s", "500ms")
func formatInterval(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%ds", int(d.Seconds()))
}

// parsePingSentence parses a proto.Sentence into PingResult.
// sentence.Map adalah field bertipe map[string]string (bukan method).
func parsePingSentence(sentence *proto.Sentence, seq int, address string) PingResult {
	m := sentence.Map // field, bukan method

	result := PingResult{
		Seq:     seq,
		Address: address,
	}

	// RouterOS mengembalikan "sent" dan "received" sebagai counter string
	// Jika received > 0, berarti ping berhasil diterima
	if received := m["received"]; received != "" && received != "0" {
		result.Received = true
	}

	// Parse size
	if size, err := strconv.Atoi(m["size"]); err == nil {
		result.Size = size
	}

	// Parse TTL
	if ttl, err := strconv.Atoi(m["ttl"]); err == nil {
		result.TTL = ttl
	}

	// Parse time — RouterOS mengembalikan dalam format "Xms" atau angka float
	if timeStr := m["time"]; timeStr != "" {
		// Hapus suffix "ms" jika ada
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
