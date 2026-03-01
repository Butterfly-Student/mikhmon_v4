//go:build integration

package mikrotik_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── NAT Rules ────────────────────────────────────────────────────────────────

func TestIntegration_NAT_GetRules(t *testing.T) {
	c := integrationClient(t)

	rules, err := c.GetNATRules(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, rules)
	t.Logf("total NAT rules: %d", len(rules))
	for _, r := range rules {
		t.Logf("  id=%-6s chain=%-10s action=%-15s disabled=%v",
			r.ID, r.Chain, r.Action, r.Disabled)
	}
}

func TestIntegration_NAT_GetRules_HasValidFields(t *testing.T) {
	c := integrationClient(t)

	rules, err := c.GetNATRules(context.Background())
	require.NoError(t, err)

	for i, r := range rules {
		assert.NotEmpty(t, r.ID, "rule %d should have non-empty ID", i)
		assert.NotEmpty(t, r.Chain, "rule %d should have non-empty Chain", i)
		assert.NotEmpty(t, r.Action, "rule %d should have non-empty Action", i)
	}
}
