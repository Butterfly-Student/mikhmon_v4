//go:build integration

package mikrotik_test

import (
	"context"
	"testing"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── GetPPPLogs ───────────────────────────────────────────────────────────────

func TestIntegration_GetPPPLogs(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	logs, err := c.GetPPPLogs(ctx, 50)
	require.NoError(t, err)
	assert.NotNil(t, logs)
	t.Logf("found %d PPP log entries", len(logs))
	for i, entry := range logs {
		if i >= 5 {
			t.Logf("  ... and %d more", len(logs)-5)
			break
		}
		t.Logf("  [%s] %s: %s", entry.Time, entry.Topics, entry.Message)
	}
}

// ─── ListenHotspotLogs ────────────────────────────────────────────────────────

func TestIntegration_ListenHotspotLogs(t *testing.T) {
	c := integrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChan := make(chan *dto.LogEntry, 20)
	cancelFn, err := c.ListenHotspotLogs(ctx, resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	var entries int
	timeout := time.After(5 * time.Second)
collect:
	for {
		select {
		case entry, ok := <-resultChan:
			if !ok {
				break collect
			}
			entries++
			t.Logf("  log %d: [%s] %s: %s", entries, entry.Time, entry.Topics, entry.Message)
		case <-timeout:
			break collect
		case <-ctx.Done():
			break collect
		}
	}

	t.Logf("collected %d hotspot log entries in 5s", entries)
	// Log streaming may yield 0 entries if no hotspot activity during test window.
}

// ─── ListenPPPLogs ────────────────────────────────────────────────────────────

func TestIntegration_ListenPPPLogs(t *testing.T) {
	c := integrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChan := make(chan *dto.LogEntry, 20)
	cancelFn, err := c.ListenPPPLogs(ctx, resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	var entries int
	timeout := time.After(5 * time.Second)
collect:
	for {
		select {
		case entry, ok := <-resultChan:
			if !ok {
				break collect
			}
			entries++
			t.Logf("  log %d: [%s] %s: %s", entries, entry.Time, entry.Topics, entry.Message)
		case <-timeout:
			break collect
		case <-ctx.Done():
			break collect
		}
	}

	t.Logf("collected %d PPP log entries in 5s", entries)
	// Log streaming may yield 0 entries if no PPP activity during test window.
}

// ─── ListenLogs generic ───────────────────────────────────────────────────────

func TestIntegration_ListenLogs_Generic(t *testing.T) {
	c := integrationClient(t)

	// =follow-only= streams NEW entries only — no historical backlog.
	// This router is very active so we should still see entries within 5s.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChan := make(chan *dto.LogEntry, 50)
	cancelFn, err := c.ListenLogs(ctx, "", resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	var entries int
collect:
	for {
		select {
		case entry, ok := <-resultChan:
			if !ok {
				break collect
			}
			entries++
			if entries <= 5 {
				t.Logf("  log %d: [%s] %s: %s", entries, entry.Time, entry.Topics, entry.Message)
			}
		case <-ctx.Done():
			break collect
		}
	}

	t.Logf("collected %d generic log entries in 5s", entries)
	// =follow-only= delivers real-time entries only; count depends on router activity.
	// assert >= 1 is reasonable for this router (DHCP/system events fire frequently).
	assert.GreaterOrEqual(t, entries, 1, "should receive at least one generic log entry")
}
