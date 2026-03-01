//go:build integration

package mikrotik_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Hotspot User Profiles ────────────────────────────────────────────────────

func TestIntegration_HotspotProfiles_GetAll(t *testing.T) {
	c := integrationClient(t)

	profiles, err := c.GetUserProfiles(context.Background())
	require.NoError(t, err)
	t.Logf("total profiles: %d", len(profiles))
	for _, p := range profiles {
		t.Logf("  name=%-20s expireMode=%-6s price=%.0f validity=%s",
			p.Name, p.ExpireMode, p.Price, p.Validity)
	}
}

func TestIntegration_HotspotProfiles_GetByID(t *testing.T) {
	c := integrationClient(t)

	profiles, err := c.GetUserProfiles(context.Background())
	require.NoError(t, err)
	if len(profiles) == 0 {
		t.Skip("no user profiles found on router")
	}

	first := profiles[0]
	profile, err := c.GetUserProfileByID(context.Background(), first.ID)
	require.NoError(t, err)
	assert.Equal(t, first.ID, profile.ID)
	assert.Equal(t, first.Name, profile.Name)
	t.Logf("GetByID ok: name=%s", profile.Name)
}

func TestIntegration_HotspotProfiles_GetByName(t *testing.T) {
	c := integrationClient(t)

	profiles, err := c.GetUserProfiles(context.Background())
	require.NoError(t, err)
	if len(profiles) == 0 {
		t.Skip("no user profiles found on router")
	}

	name := profiles[0].Name
	profile, err := c.GetUserProfileByName(context.Background(), name)
	require.NoError(t, err)
	assert.Equal(t, name, profile.Name)
	t.Logf("GetByName ok: name=%s", profile.Name)
}

func TestIntegration_HotspotProfiles_AddUpdateRemove(t *testing.T) {
	c := integrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testName := fmt.Sprintf("inttest-%d", time.Now().UnixMilli())

	gen := mikrotik.NewOnLoginGenerator()
	script := gen.Generate(&dto.ProfileRequest{
		Name:         testName,
		ExpireMode:   "remc",
		Validity:     "1d",
		Price:        1000,
		SellingPrice: 1100,
		LockUser:     "Disable",
		LockServer:   "Disable",
	})

	newProfile := &dto.UserProfile{
		Name:    testName,
		OnLogin: script,
	}

	// Add — RouterOS tidak selalu mengembalikan ID dari /add,
	// jadi kita lookup by name setelah add.
	_, err := c.AddUserProfile(ctx, newProfile)
	require.NoError(t, err)
	t.Logf("added profile name=%s", testName)

	created, err := c.GetUserProfileByName(ctx, testName)
	require.NoError(t, err)
	require.NotNil(t, created, "profile must exist after creation")
	require.NotEmpty(t, created.ID, "profile must have an ID after creation")
	id := created.ID
	t.Logf("resolved id=%s for name=%s", id, testName)
	assert.Equal(t, testName, created.Name)

	t.Cleanup(func() {
		ctxClean, cancelClean := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelClean()
		_ = c.RemoveUserProfile(ctxClean, id)
	})

	// Update — ubah rate-limit
	updated := &dto.UserProfile{
		Name:      testName,
		RateLimit: "1M/1M",
		OnLogin:   script,
	}
	err = c.UpdateUserProfile(ctx, id, updated)
	require.NoError(t, err)
	t.Logf("updated profile id=%s", id)

	after, err := c.GetUserProfileByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "1M/1M", after.RateLimit)

	// Remove
	err = c.RemoveUserProfile(ctx, id)
	require.NoError(t, err)
	t.Logf("removed profile id=%s", id)

	gone, err := c.GetUserProfileByName(ctx, testName)
	assert.NoError(t, err)
	assert.Nil(t, gone, "profile should not exist after removal")
}
