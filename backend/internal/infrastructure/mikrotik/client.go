package mikrotik

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/go-routeros/routeros/v3"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
	"go.uber.org/zap"
)

// Client wraps the MikroTik RouterOS client with a connection pool.
// Koneksi di-cache per router ID dan di-reconnect otomatis jika terputus.
type Client struct {
	mu      sync.Mutex
	clients map[uint]*routeros.Client
	log     *zap.Logger
}

// NewClient creates a new MikroTik client with connection pool
func NewClient(log *zap.Logger) *Client {
	if log == nil {
		log = zap.NewNop()
	}
	return &Client{
		clients: make(map[uint]*routeros.Client),
		log:     log.Named("mikrotik"),
	}
}

// getClient returns a cached connection or dials a new one.
// Health check dilakukan SETELAH melepas mutex agar tidak blocking operasi lain.
func (c *Client) getClient(router *entity.Router) (*routeros.Client, error) {
	// === Phase 1: cek cache (lock singkat) ===
	c.mu.Lock()
	existing, cached := c.clients[router.ID]
	c.mu.Unlock()

	if cached {
		// Health check di luar lock — tidak mengblokir goroutine lain
		healthCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := existing.RunContext(healthCtx, "/system/identity/print")
		if err == nil {
			return existing, nil
		}

		// Koneksi mati; tutup dan hapus dari pool
		c.log.Warn("Cached connection dead, reconnecting",
			zap.Uint("routerID", router.ID),
			zap.String("host", router.Host),
			zap.Error(err),
		)
		existing.Close()

		c.mu.Lock()
		// Double-check: goroutine lain mungkin sudah update saat lock dilepas
		if c.clients[router.ID] == existing {
			delete(c.clients, router.ID)
		}
		c.mu.Unlock()
	}

	// === Phase 2: buat koneksi baru (dengan retry) ===
	client, err := c.dialWithRetry(router)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.clients[router.ID] = client
	c.mu.Unlock()

	return client, nil
}

// dialWithRetry mencoba dial dengan max 2 retry dan exponential backoff.
func (c *Client) dialWithRetry(router *entity.Router) (*routeros.Client, error) {
	const maxRetries = 2
	delays := []time.Duration{500 * time.Millisecond, 1 * time.Second}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			c.log.Info("Retrying MikroTik connection",
				zap.Uint("routerID", router.ID),
				zap.String("host", router.Host),
				zap.Int("attempt", attempt),
				zap.Error(lastErr),
			)
			time.Sleep(delays[attempt-1])
		}

		client, err := c.dial(router)
		if err == nil {
			if attempt > 0 {
				c.log.Info("MikroTik connection succeeded after retry",
					zap.Uint("routerID", router.ID),
					zap.String("host", router.Host),
					zap.Int("attempt", attempt),
				)
			}
			return client, nil
		}
		lastErr = err
	}

	c.log.Error("MikroTik connection failed after retries",
		zap.Uint("routerID", router.ID),
		zap.String("host", router.Host),
		zap.Int("maxRetries", maxRetries),
		zap.Error(lastErr),
	)
	return nil, fmt.Errorf("connection failed after %d attempts: %w", maxRetries+1, lastErr)
}

// dial membuka koneksi TCP baru ke router dan login.
// Timeout dial dinaikkan ke 10 detik untuk menangani jaringan yang lambat/variabel.
func (c *Client) dial(router *entity.Router) (*routeros.Client, error) {
	password := router.Password // TODO: implement decryption
	addr := net.JoinHostPort(router.Host, strconv.Itoa(router.Port))

	dialer := &net.Dialer{Timeout: 10 * time.Second}

	var client *routeros.Client

	if router.UseSSL {
		conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to router %s: %w", router.Host, err)
		}
		if client, err = routeros.NewClient(conn); err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to create MikroTik client: %w", err)
		}
	} else {
		conn, err := dialer.Dial("tcp", addr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to router %s: %w", router.Host, err)
		}
		if client, err = routeros.NewClient(conn); err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to create MikroTik client: %w", err)
		}
	}

	if err := client.Login(router.Username, password); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to login to router %s: %w", router.Host, err)
	}

	c.log.Info("Connected to MikroTik router",
		zap.String("host", router.Host),
		zap.Uint("routerID", router.ID),
		zap.Bool("ssl", router.UseSSL),
	)

	return client, nil
}

// TestConnection tests if a connection can be established to the router.
// Selalu membuat koneksi baru (tidak menggunakan cache).
func (c *Client) TestConnection(ctx context.Context, router *entity.Router) error {
	client, err := c.dial(router)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.RunContext(ctx, "/system/identity/print")
	return err
}

// CloseRouter menutup dan menghapus koneksi untuk router tertentu dari pool
func (c *Client) CloseRouter(routerID uint) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if client, ok := c.clients[routerID]; ok {
		client.Close()
		delete(c.clients, routerID)
	}
}

// CloseAll menutup semua koneksi dalam pool
func (c *Client) CloseAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, client := range c.clients {
		client.Close()
		delete(c.clients, id)
	}
}

// Helper: parse int from RouterOS string
func parseInt(s string) int64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// Helper: parse float from RouterOS string
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// Helper: parse bool from RouterOS string
func parseBool(s string) bool {
	return s == "true" || s == "yes"
}

// Helper: format int to string
func formatInt(n int64) string {
	return strconv.FormatInt(n, 10)
}
