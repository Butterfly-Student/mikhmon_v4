//go:build integration

package mikrotik_test

import (
	"context"
	"testing"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Expire Monitor ───────────────────────────────────────────────────────────

func TestIntegration_ExpireMonitor_Ensure(t *testing.T) {
	c := integrationClient(t)

	gen := mikrotik.NewOnLoginGenerator()
	script := gen.GenerateExpireMonitorScript()
	require.NotEmpty(t, script, "expire monitor script must not be empty")

	status, err := c.EnsureExpireMonitor(context.Background(), script)
	require.NoError(t, err)
	assert.Contains(t, []string{"created", "enabled", "existing"}, status)
	t.Logf("EnsureExpireMonitor status: %s", status)
}

func TestIntegration_ExpireMonitor_IdempotentCall(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	gen := mikrotik.NewOnLoginGenerator()
	script := gen.GenerateExpireMonitorScript()

	status1, err := c.EnsureExpireMonitor(ctx, script)
	require.NoError(t, err)
	t.Logf("first call status: %s", status1)

	status2, err := c.EnsureExpireMonitor(ctx, script)
	require.NoError(t, err)
	t.Logf("second call status: %s", status2)

	assert.Equal(t, "existing", status2, "second call should return 'existing'")
}

// TestIntegration_ExpireMonitor_ScriptContent memvalidasi konten script
// yang dihasilkan OnLoginGenerator sebelum dikirim ke router.
func TestIntegration_ExpireMonitor_ScriptContent(t *testing.T) {
	gen := mikrotik.NewOnLoginGenerator()
	script := gen.GenerateExpireMonitorScript()

	t.Logf("expire monitor script length: %d chars", len(script))
	assert.NotEmpty(t, script)

	req := &dto.ProfileRequest{
		Name:         "TestProfile",
		ExpireMode:   "remc",
		Validity:     "30d",
		Price:        5000,
		SellingPrice: 5500,
		LockUser:     "Enable",
		LockServer:   "Disable",
	}

	loginScript := gen.Generate(req)
	assert.NotEmpty(t, loginScript)
	t.Logf("on-login script length: %d chars", len(loginScript))

	parsed := gen.Parse(loginScript)
	assert.Equal(t, "remc", parsed.ExpireMode)
	assert.Equal(t, "30d", parsed.Validity)
	assert.InDelta(t, float64(5000), parsed.Price, 0.01)
	assert.InDelta(t, float64(5500), parsed.SellingPrice, 0.01)
	t.Logf("parsed: expireMode=%s validity=%s price=%.0f sellingPrice=%.0f",
		parsed.ExpireMode, parsed.Validity, parsed.Price, parsed.SellingPrice)
}
