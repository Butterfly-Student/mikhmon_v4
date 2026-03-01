//go:build integration

package mikrotik_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Hotspot Active Sessions ──────────────────────────────────────────────────

func TestIntegration_HotspotActive_RemoveSession(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	sessions, err := c.GetHotspotActive(ctx)
	require.NoError(t, err)
	assert.NotNil(t, sessions)

	if len(sessions) == 0 {
		t.Skip("no active hotspot sessions — skipping remove test")
	}

	target := sessions[0]
	t.Logf("removing active session id=%s user=%s address=%s",
		target.ID, target.User, target.Address)

	err = c.RemoveHotspotActive(ctx, target.ID)
	require.NoError(t, err)
	t.Logf("RemoveHotspotActive succeeded for id=%s", target.ID)
}
