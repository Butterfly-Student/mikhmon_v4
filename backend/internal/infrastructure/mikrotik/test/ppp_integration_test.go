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

// ─── PPP Secrets ─────────────────────────────────────────────────────────────

func TestIntegration_PPPSecrets_List(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	secrets, err := c.GetPPPSecrets(ctx, "")
	require.NoError(t, err)
	assert.NotNil(t, secrets)
	t.Logf("found %d PPP secrets", len(secrets))
	for _, s := range secrets {
		t.Logf("  name=%-20s profile=%-15s service=%s disabled=%v",
			s.Name, s.Profile, s.Service, s.Disabled)
	}
}

func TestIntegration_PPPSecrets_CRUD(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	name := "test-secret-crud"

	// Clean up any pre-existing test secret.
	existing, err := c.GetPPPSecretByName(ctx, name)
	require.NoError(t, err)
	if existing != nil {
		_ = c.RemovePPPSecret(ctx, existing.ID)
	}

	// Add
	secret := &dto.PPPSecret{
		Name:     name,
		Password: "testpass123",
		Profile:  "default",
		Service:  "any",
		Comment:  "integration-test",
	}
	err = c.AddPPPSecret(ctx, secret)
	require.NoError(t, err, "AddPPPSecret should succeed")

	// GetByName
	got, err := c.GetPPPSecretByName(ctx, name)
	require.NoError(t, err)
	require.NotNil(t, got, "GetPPPSecretByName should return the created secret")
	assert.Equal(t, name, got.Name)
	assert.Equal(t, "default", got.Profile)
	t.Logf("created PPP secret id=%s name=%s", got.ID, got.Name)

	// Update
	updated := &dto.PPPSecret{
		Name:    name,
		Profile: "default",
		Comment: "integration-test-updated",
	}
	err = c.UpdatePPPSecret(ctx, got.ID, updated)
	require.NoError(t, err, "UpdatePPPSecret should succeed")

	// GetByID
	got2, err := c.GetPPPSecretByID(ctx, got.ID)
	require.NoError(t, err)
	require.NotNil(t, got2)
	assert.Equal(t, "integration-test-updated", got2.Comment)
	t.Logf("updated PPP secret comment=%s", got2.Comment)

	// Remove
	err = c.RemovePPPSecret(ctx, got.ID)
	require.NoError(t, err, "RemovePPPSecret should succeed")

	// Verify removal
	gone, err := c.GetPPPSecretByName(ctx, name)
	require.NoError(t, err)
	assert.Nil(t, gone, "secret should be gone after removal")
	t.Logf("PPP secret removed successfully")
}

// ─── PPP Secret Disable / Enable ─────────────────────────────────────────────

func TestIntegration_PPPSecrets_DisableEnable(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	name := "test-secret-disable"

	// Clean up any pre-existing test secret.
	existing, err := c.GetPPPSecretByName(ctx, name)
	require.NoError(t, err)
	if existing != nil {
		_ = c.RemovePPPSecret(ctx, existing.ID)
	}

	// Add a secret to operate on.
	err = c.AddPPPSecret(ctx, &dto.PPPSecret{
		Name:    name,
		Profile: "default",
		Comment: "integration-test-disable",
	})
	require.NoError(t, err)

	// Lookup ID by name (RouterOS 6.x does not return ID from /add).
	got, err := c.GetPPPSecretByName(ctx, name)
	require.NoError(t, err)
	require.NotNil(t, got)
	t.Logf("created secret id=%s name=%s disabled=%v", got.ID, got.Name, got.Disabled)

	// Disable
	err = c.DisablePPPSecret(ctx, got.ID)
	require.NoError(t, err, "DisablePPPSecret should succeed")

	dis, err := c.GetPPPSecretByID(ctx, got.ID)
	require.NoError(t, err)
	require.NotNil(t, dis)
	assert.True(t, dis.Disabled, "secret should be disabled")
	t.Logf("disabled secret id=%s disabled=%v", dis.ID, dis.Disabled)

	// Re-enable
	err = c.EnablePPPSecret(ctx, got.ID)
	require.NoError(t, err, "EnablePPPSecret should succeed")

	en, err := c.GetPPPSecretByID(ctx, got.ID)
	require.NoError(t, err)
	require.NotNil(t, en)
	assert.False(t, en.Disabled, "secret should be enabled again")
	t.Logf("re-enabled secret id=%s disabled=%v", en.ID, en.Disabled)

	// Cleanup
	require.NoError(t, c.RemovePPPSecret(ctx, got.ID))
}

// ─── PPP Profiles ─────────────────────────────────────────────────────────────

func TestIntegration_PPPProfiles_List(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	profiles, err := c.GetPPPProfiles(ctx)
	require.NoError(t, err)
	assert.NotNil(t, profiles)
	t.Logf("found %d PPP profiles", len(profiles))
	for _, p := range profiles {
		t.Logf("  name=%-20s localAddr=%-16s remoteAddr=%s",
			p.Name, p.LocalAddress, p.RemoteAddress)
	}
}

func TestIntegration_PPPProfiles_CRUD(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	name := "test-profile-crud"

	// Clean up any pre-existing test profile.
	existing, err := c.GetPPPProfileByName(ctx, name)
	require.NoError(t, err)
	if existing != nil {
		_ = c.RemovePPPProfile(ctx, existing.ID)
	}

	// Add
	err = c.AddPPPProfile(ctx, &dto.PPPProfile{
		Name:      name,
		RateLimit: "1M/1M",
		Comment:   "integration-test",
	})
	require.NoError(t, err, "AddPPPProfile should succeed")

	// GetByName (RouterOS 6.x does not return ID from /add)
	got, err := c.GetPPPProfileByName(ctx, name)
	require.NoError(t, err)
	require.NotNil(t, got, "GetPPPProfileByName should return the created profile")
	assert.Equal(t, name, got.Name)
	assert.Equal(t, "1M/1M", got.RateLimit)
	t.Logf("created PPP profile id=%s name=%s rateLimit=%s", got.ID, got.Name, got.RateLimit)

	// Update
	err = c.UpdatePPPProfile(ctx, got.ID, &dto.PPPProfile{
		Name:      name,
		RateLimit: "2M/2M",
		Comment:   "integration-test-updated",
	})
	require.NoError(t, err, "UpdatePPPProfile should succeed")

	// GetByID
	got2, err := c.GetPPPProfileByID(ctx, got.ID)
	require.NoError(t, err)
	require.NotNil(t, got2)
	assert.Equal(t, "2M/2M", got2.RateLimit)
	assert.Equal(t, "integration-test-updated", got2.Comment)
	t.Logf("updated PPP profile rateLimit=%s comment=%s", got2.RateLimit, got2.Comment)

	// Remove
	err = c.RemovePPPProfile(ctx, got.ID)
	require.NoError(t, err, "RemovePPPProfile should succeed")

	// Verify removal
	gone, err := c.GetPPPProfileByName(ctx, name)
	require.NoError(t, err)
	assert.Nil(t, gone, "profile should be gone after removal")
	t.Logf("PPP profile removed successfully")
}

// ─── PPP Active ───────────────────────────────────────────────────────────────

func TestIntegration_PPPActive_List(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	active, err := c.GetPPPActive(ctx, "")
	require.NoError(t, err)
	assert.NotNil(t, active)
	t.Logf("found %d active PPP sessions", len(active))
	for _, a := range active {
		t.Logf("  name=%-20s service=%-10s address=%-16s uptime=%s",
			a.Name, a.Service, a.Address, a.Uptime)
	}
}

func TestIntegration_PPPActive_Listen(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	// Skip if no active PPP sessions — RouterOS =follow= sends nothing for an
	// empty table, so we would never receive a batch.
	active, err := c.GetPPPActive(ctx, "")
	require.NoError(t, err)
	if len(active) == 0 {
		t.Skip("no active PPP sessions — skipping listen test")
	}

	listenCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resultChan := make(chan []*dto.PPPActive, 10)
	cancelFn, err := c.ListenPPPActive(listenCtx, "1s", resultChan)
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
			t.Logf("  batch %d: %d entries", batches, len(batch))
			if batches >= 1 {
				break collect
			}
		case <-listenCtx.Done():
			break collect
		}
	}

	assert.GreaterOrEqual(t, batches, 1, "should receive at least one PPP active batch")
}

// ─── PPP Inactive ─────────────────────────────────────────────────────────────

func TestIntegration_PPPInactive_Listen(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	// Skip if no PPP secrets configured — RouterOS =follow= sends nothing for an
	// empty table, so listenBatches would block and we'd receive no results.
	secrets, err := c.GetPPPSecrets(ctx, "")
	require.NoError(t, err)
	if len(secrets) == 0 {
		t.Skip("no PPP secrets configured — skipping inactive listen test")
	}

	listenCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resultChan := make(chan []*dto.PPPSecret, 10)
	cancelFn, err := c.ListenPPPInactive(listenCtx, "1s", resultChan)
	require.NoError(t, err)
	defer cancelFn() //nolint:errcheck

	var results int
collect:
	for {
		select {
		case inactive, ok := <-resultChan:
			if !ok {
				break collect
			}
			results++
			t.Logf("  result %d: %d inactive PPP secrets", results, len(inactive))
			if results >= 1 {
				break collect
			}
		case <-listenCtx.Done():
			break collect
		}
	}

	assert.GreaterOrEqual(t, results, 1, "should receive at least one PPP inactive result")
}
