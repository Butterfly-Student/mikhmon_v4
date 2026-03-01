//go:build integration

package mikrotik_test

import (
	"context"
	"testing"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// ─── Konfigurasi koneksi ──────────────────────────────────────────────────────
// Ubah nilai di bawah ini sesuai router MikroTik Anda sebelum menjalankan test.
const (
	testHost     = "192.168.233.1"
	testPort     = 8728
	testUsername = "admin"
	testPassword = "r00t"
)

// ─── Helpers ──────────────────────────────────────────────────────────────────

func integrationClient(t *testing.T) *mikrotik.Client {
	t.Helper()

	if testPassword == "" {
		t.Skip("skipping integration test: testPassword belum diisi di integration_test.go")
	}

	cfg := mikrotik.Config{
		Host:     testHost,
		Port:     testPort,
		Username: testUsername,
		Password: testPassword,
		Timeout:  10 * time.Second,
	}

	logger, _ := zap.NewDevelopment()

	const (
		attemptTimeout = 8 * time.Second
		retryInterval  = 5 * time.Second
		totalTimeout   = 30 * time.Second
	)

	deadline := time.Now().Add(totalTimeout)
	var lastErr error
	for attempt := 1; time.Now().Before(deadline); attempt++ {
		c := mikrotik.NewClient(cfg, logger)
		ctx, cancel := context.WithTimeout(context.Background(), attemptTimeout)
		err := c.Connect(ctx)
		cancel()
		if err == nil {
			t.Logf("connected to %s after %d attempt(s)", cfg.Host, attempt)
			t.Cleanup(func() { c.Close() })
			return c
		}
		lastErr = err
		t.Logf("attempt %d failed: %v — retry in %s", attempt, err, retryInterval)
		time.Sleep(retryInterval)
	}
	t.Fatalf("could not connect to %s: %v", testHost, lastErr)
	return nil
}

// firstRunningInterface returns the name of the first running non-disabled interface.
func firstRunningInterface(t *testing.T, c *mikrotik.Client) string {
	t.Helper()
	ifaces, err := c.GetInterfaces(context.Background())
	require.NoError(t, err)
	for _, iface := range ifaces {
		if iface.Running && !iface.Disabled {
			return iface.Name
		}
	}
	t.Skip("no running interface found on router")
	return ""
}

// firstQueueName returns the name of the first simple queue, or skips if none.
func firstQueueName(t *testing.T, c *mikrotik.Client) string {
	t.Helper()
	queues, err := c.GetAllQueues(context.Background())
	require.NoError(t, err)
	if len(queues) == 0 {
		t.Skip("no simple queue found on router")
	}
	return queues[0]
}

// ─── Connection & Async mode ──────────────────────────────────────────────────

func TestIntegration_Connect(t *testing.T) {
	c := integrationClient(t)
	assert.NotNil(t, c)
}

func TestIntegration_Async_IsAsync(t *testing.T) {
	c := integrationClient(t)
	assert.True(t, c.IsAsync(), "connection must be in async mode after Connect()")
}

func TestIntegration_RunContext_WithTimeout(t *testing.T) {
	c := integrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reply, err := c.RunContext(ctx, "/system/identity/print")
	require.NoError(t, err)
	assert.NotEmpty(t, reply.Re)
}

func TestIntegration_RunMany_Concurrent(t *testing.T) {
	c := integrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	commands := [][]string{
		{"/system/resource/print"},
		{"/system/identity/print"},
		{"/ip/pool/print"},
		{"/interface/print"},
		{"/ip/firewall/nat/print"},
	}

	start := time.Now()
	replies, errs := c.RunMany(ctx, commands)
	elapsed := time.Since(start)

	for i, err := range errs {
		require.NoError(t, err, "command %d (%s) failed", i, commands[i][0])
	}
	assert.NotEmpty(t, replies[0].Re, "resource/print must return data")
	assert.NotEmpty(t, replies[1].Re, "identity/print must return data")
	assert.NotEmpty(t, replies[3].Re, "interface/print must return data")
	t.Logf("5 concurrent commands finished in %s via single async conn", elapsed)
}

// ─── System ───────────────────────────────────────────────────────────────────

func TestIntegration_System_GetResource(t *testing.T) {
	c := integrationClient(t)

	res, err := c.GetSystemResource(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, res.Version)
	assert.NotEmpty(t, res.BoardName)
	assert.Greater(t, res.TotalMemory, int64(0))
	t.Logf("RouterOS %s on %s | uptime: %s | CPU: %d%% | RAM: %d/%d MB",
		res.Version, res.BoardName, res.Uptime,
		res.CpuLoad, res.FreeMemory/1024/1024, res.TotalMemory/1024/1024)
}

func TestIntegration_System_GetIdentity(t *testing.T) {
	c := integrationClient(t)

	id, err := c.GetSystemIdentity(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, id.Name)
	t.Logf("router identity: %s", id.Name)
}

func TestIntegration_System_GetHealth(t *testing.T) {
	c := integrationClient(t)

	health, err := c.GetSystemHealth(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, health)
	t.Logf("health: voltage=%s temperature=%s", health.Voltage, health.Temperature)
}

func TestIntegration_System_GetClock(t *testing.T) {
	c := integrationClient(t)

	clock, err := c.GetSystemClock(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, clock.Time)
	t.Logf("clock: %s %s (%s)", clock.Date, clock.Time, clock.TimeZoneName)
}

func TestIntegration_System_GetRouterBoardInfo(t *testing.T) {
	c := integrationClient(t)

	rb, err := c.GetRouterBoardInfo(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, rb)
	t.Logf("routerboard: model=%s serial=%s firmware=%s", rb.Model, rb.SerialNumber, rb.CurrentFirmware)
}

// ─── Interface ────────────────────────────────────────────────────────────────

func TestIntegration_Interface_GetAll(t *testing.T) {
	c := integrationClient(t)

	ifaces, err := c.GetInterfaces(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, ifaces)
	t.Logf("found %d interfaces:", len(ifaces))
	for _, i := range ifaces {
		t.Logf("  %-20s type=%-12s running=%-5v disabled=%v", i.Name, i.Type, i.Running, i.Disabled)
	}
}

// ─── IP Pool ──────────────────────────────────────────────────────────────────

func TestIntegration_Pool_GetAddressPools(t *testing.T) {
	c := integrationClient(t)

	pools, err := c.GetAddressPools(context.Background())
	require.NoError(t, err)
	t.Logf("found %d address pools:", len(pools))
	for _, p := range pools {
		t.Logf("  %s", p)
	}
}

// ─── Queue ────────────────────────────────────────────────────────────────────

func TestIntegration_Queue_GetAllQueues(t *testing.T) {
	c := integrationClient(t)

	queues, err := c.GetAllQueues(context.Background())
	require.NoError(t, err)
	t.Logf("found %d simple queues:", len(queues))
	for _, q := range queues {
		t.Logf("  %s", q)
	}
}

func TestIntegration_Queue_GetAllParentQueues(t *testing.T) {
	c := integrationClient(t)

	queues, err := c.GetAllParentQueues(context.Background())
	require.NoError(t, err)
	t.Logf("found %d parent queues:", len(queues))
}

// ─── Hotspot ──────────────────────────────────────────────────────────────────

func TestIntegration_Hotspot_GetServers(t *testing.T) {
	c := integrationClient(t)

	servers, err := c.GetHotspotServers(context.Background())
	require.NoError(t, err)
	t.Logf("found %d hotspot servers: %v", len(servers), servers)
}

func TestIntegration_Hotspot_GetActive(t *testing.T) {
	c := integrationClient(t)

	active, err := c.GetHotspotActive(context.Background())
	require.NoError(t, err)
	t.Logf("found %d active hotspot sessions", len(active))
	for _, a := range active {
		t.Logf("  user=%-20s address=%-16s uptime=%s", a.User, a.Address, a.Uptime)
	}
}

func TestIntegration_Hotspot_GetActiveCount(t *testing.T) {
	c := integrationClient(t)

	count, err := c.GetHotspotActiveCount(context.Background())
	require.NoError(t, err)
	t.Logf("active hotspot sessions: %d", count)
}

// ─── Listener: Traffic Monitor ────────────────────────────────────────────────

func TestIntegration_Listener_TrafficMonitor_Once(t *testing.T) {
	c := integrationClient(t)
	ifaceName := firstRunningInterface(t, c)
	t.Logf("monitoring interface: %s", ifaceName)

	// =once= meminta RouterOS mengirim satu sampel lalu menutup stream.
	lr, err := c.ListenArgs([]string{
		"/interface/monitor-traffic",
		"=interface=" + ifaceName,
		"=once=",
	})
	require.NoError(t, err)

	received := 0
	timeout := time.After(8 * time.Second)
loop:
	for {
		select {
		case sen, ok := <-lr.Chan():
			if !ok {
				break loop
			}
			t.Logf("  sample %d: rx=%s bps  tx=%s bps",
				received+1, sen.Map["rx-bits-per-second"], sen.Map["tx-bits-per-second"])
			received++
		case <-timeout:
			break loop
		}
	}

	_, err = lr.Cancel()
	assert.NoError(t, err)
	assert.Greater(t, received, 0, "must receive at least one traffic sample")
}

func TestIntegration_Listener_TrafficMonitor_Continuous(t *testing.T) {
	c := integrationClient(t)
	ifaceName := firstRunningInterface(t, c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resultChan := make(chan dto.TrafficMonitorStats, 10)
	cancelFn, err := c.StartTrafficMonitorListen(ctx, ifaceName, resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	received := 0
	timeout := time.After(8 * time.Second)
collect:
	for {
		select {
		case stat, ok := <-resultChan:
			if !ok {
				break collect
			}
			t.Logf("  stat %d: iface=%s rx=%d bps tx=%d bps",
				received+1, stat.Name, stat.RxBitsPerSecond, stat.TxBitsPerSecond)
			assert.Equal(t, ifaceName, stat.Name)
			assert.False(t, stat.Timestamp.IsZero())
			received++
			if received >= 3 {
				break collect
			}
		case <-timeout:
			break collect
		case <-ctx.Done():
			break collect
		}
	}

	assert.Greater(t, received, 0, "must receive at least one traffic stat")
}

// ─── Listener: Queue Stats ────────────────────────────────────────────────────

func TestIntegration_Listener_QueueStats(t *testing.T) {
	c := integrationClient(t)
	queueName := firstQueueName(t, c)
	t.Logf("monitoring queue: %s", queueName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dto.DefaultQueueStatsConfig(queueName)
	resultChan := make(chan dto.QueueStats, 10)
	cancelFn, err := c.StartQueueStatsListen(ctx, cfg, resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	received := 0
	timeout := time.After(8 * time.Second)
collect:
	for {
		select {
		case stat, ok := <-resultChan:
			if !ok {
				break collect
			}
			t.Logf("  stat %d: queue=%s rateIn=%d rateOut=%d",
				received+1, stat.Name, stat.RateIn, stat.RateOut)
			assert.Equal(t, queueName, stat.Name)
			received++
			if received >= 3 {
				break collect
			}
		case <-timeout:
			break collect
		case <-ctx.Done():
			break collect
		}
	}

	assert.Greater(t, received, 0, "must receive at least one queue stat")
}

// ─── Listener: Resource Monitor ──────────────────────────────────────────────

func TestIntegration_Listener_ResourceMonitor(t *testing.T) {
	c := integrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resultChan := make(chan dto.SystemResourceMonitorStats, 10)
	cancelFn, err := c.StartSystemResourceMonitorListen(ctx, resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	received := 0
	timeout := time.After(8 * time.Second)
collect:
	for {
		select {
		case stat, ok := <-resultChan:
			if !ok {
				break collect
			}
			t.Logf("  stat %d: cpu=%d%% freeMem=%.1fMiB uptime=%s", received+1, stat.CPULoad, float64(stat.FreeMemory)/1024/1024, stat.Uptime)
			assert.NotEmpty(t, stat.Uptime)
			received++
			if received >= 3 {
				break collect
			}
		case <-timeout:
			break collect
		case <-ctx.Done():
			break collect
		}
	}

	assert.Greater(t, received, 0, "must receive at least one resource stat")
}

// ─── Listener: Ping ───────────────────────────────────────────────────────────

func TestIntegration_Listener_Ping(t *testing.T) {
	c := integrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := dto.DefaultPingConfig("8.8.8.8")
	resultChan := make(chan dto.PingResult, 10)
	cancelFn, err := c.StartPingListen(ctx, cfg, resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	received := 0
	timeout := time.After(8 * time.Second)
collect:
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				break collect
			}
			t.Logf("  ping %d: addr=%s time=%.1fms received=%v",
				result.Seq+1, result.Address, result.TimeMs, result.Received)
			received++
			if received >= 3 {
				break collect
			}
		case <-timeout:
			break collect
		case <-ctx.Done():
			break collect
		}
	}

	assert.Greater(t, received, 0, "must receive at least one ping result")
}

// ─── Concurrent: Run + Listen bersamaan ──────────────────────────────────────

func TestIntegration_RunAndListenConcurrent(t *testing.T) {
	c := integrationClient(t)
	ifaceName := firstRunningInterface(t, c)

	// Mulai streaming traffic
	lr, err := c.ListenArgs([]string{
		"/interface/monitor-traffic",
		"=interface=" + ifaceName,
	})
	require.NoError(t, err)
	defer lr.Cancel() //nolint:errcheck

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Jalankan beberapa command secara concurrent di atas koneksi yang sama
	replies, errs := c.RunMany(ctx, [][]string{
		{"/system/resource/print"},
		{"/system/identity/print"},
		{"/ip/pool/print"},
	})
	for i, err := range errs {
		require.NoError(t, err, "RunMany command %d should succeed while listener active", i)
	}
	assert.NotEmpty(t, replies[0].Re)
	assert.NotEmpty(t, replies[1].Re)
	t.Log("RunMany succeeded while ListenArgs streaming — async confirmed")

	// Pastikan traffic sample masih diterima
	select {
	case sen, ok := <-lr.Chan():
		require.True(t, ok)
		t.Logf("concurrent traffic sample: rx=%s tx=%s",
			sen.Map["rx-bits-per-second"], sen.Map["tx-bits-per-second"])
	case <-time.After(5 * time.Second):
		t.Fatal("no traffic sample received during concurrent test")
	}
}
