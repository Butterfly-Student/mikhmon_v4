//go:build integration

package mikrotik_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Hotspot Hosts ────────────────────────────────────────────────────────────

func TestIntegration_HotspotHosts_GetAll(t *testing.T) {
	c := integrationClient(t)

	hosts, err := c.GetHotspotHosts(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, hosts)
	t.Logf("total hotspot hosts: %d", len(hosts))
	for _, h := range hosts {
		t.Logf("  mac=%-18s address=%-16s authorized=%v bypassed=%v",
			h.MACAddress, h.Address, h.Authorized, h.Bypassed)
	}
}

func TestIntegration_HotspotHosts_Remove(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	hosts, err := c.GetHotspotHosts(ctx)
	require.NoError(t, err)
	if len(hosts) == 0 {
		t.Skip("no hotspot hosts found — skipping remove test")
	}

	target := hosts[0]
	t.Logf("removing host id=%s mac=%s address=%s", target.ID, target.MACAddress, target.Address)

	err = c.RemoveHotspotHost(ctx, target.ID)
	require.NoError(t, err)
	t.Logf("RemoveHotspotHost succeeded for id=%s", target.ID)
}
