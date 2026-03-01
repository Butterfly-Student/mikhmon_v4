//go:build integration

package mikrotik_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Hotspot Logging ──────────────────────────────────────────────────────────

func TestIntegration_Logging_EnableHotspotLogging(t *testing.T) {
	c := integrationClient(t)

	err := c.EnableHotspotLogging(context.Background())
	require.NoError(t, err)
	t.Log("EnableHotspotLogging succeeded (or logging already configured)")
}

func TestIntegration_Logging_GetHotspotLogs_All(t *testing.T) {
	c := integrationClient(t)

	logs, err := c.GetHotspotLogs(context.Background(), 0)
	require.NoError(t, err)
	assert.NotNil(t, logs)
	t.Logf("total hotspot log entries (no limit): %d", len(logs))
	for i, l := range logs {
		t.Logf("  [%d] time=%-20s topics=%-30s msg=%s", i+1, l.Time, l.Topics, l.Message)
	}
}

func TestIntegration_Logging_EnablePPPLogging(t *testing.T) {
	c := integrationClient(t)

	err := c.EnablePPPLogging(context.Background())
	require.NoError(t, err)
	t.Log("EnablePPPLogging succeeded (or logging already configured)")
}

func TestIntegration_Logging_GetHotspotLogs_WithLimit(t *testing.T) {
	c := integrationClient(t)

	limit := 5
	logs, err := c.GetHotspotLogs(context.Background(), limit)
	require.NoError(t, err)
	assert.NotNil(t, logs)
	assert.LessOrEqual(t, len(logs), limit, "result must not exceed requested limit")
	t.Logf("hotspot logs (limit=%d): got %d entries", limit, len(logs))
	for _, l := range logs {
		assert.NotEmpty(t, l.ID)
		t.Logf("  time=%s msg=%s", l.Time, l.Message)
	}
}
