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

// ─── Hotspot Active Listen ────────────────────────────────────────────────────

func TestIntegration_HotspotActive_Listen(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	// Skip if there are no active sessions — RouterOS =follow= sends nothing
	// for an empty table, so we would never receive a batch.
	active, err := c.GetHotspotActive(ctx)
	require.NoError(t, err)
	if len(active) == 0 {
		t.Skip("no active hotspot sessions — skipping listen test")
	}

	listenCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resultChan := make(chan []*dto.HotspotActive, 10)
	cancelFn, err := c.ListenHotspotActive(listenCtx, "1s", resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	var batches int
collect:
	for {
		select {
		case batch, ok := <-resultChan:
			if !ok {
				break collect
			}
			batches++
			t.Logf("  batch %d: %d hotspot active entries", batches, len(batch))
			if batches >= 1 {
				break collect
			}
		case <-listenCtx.Done():
			break collect
		}
	}

	assert.GreaterOrEqual(t, batches, 1, "should receive at least one hotspot active batch")
}

// ─── Hotspot Inactive Listen ──────────────────────────────────────────────────

func TestIntegration_HotspotInactive_Listen(t *testing.T) {
	c := integrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resultChan := make(chan []*dto.HotspotUser, 10)
	cancelFn, err := c.ListenHotspotInactive(ctx, "1s", resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	var results int
	timeout := time.After(3 * time.Second)
collect:
	for {
		select {
		case inactive, ok := <-resultChan:
			if !ok {
				break collect
			}
			results++
			t.Logf("  result %d: %d inactive hotspot users", results, len(inactive))
			if results >= 1 {
				break collect
			}
		case <-timeout:
			break collect
		case <-ctx.Done():
			break collect
		}
	}

	assert.GreaterOrEqual(t, results, 1, "should receive at least one hotspot inactive result")
}
