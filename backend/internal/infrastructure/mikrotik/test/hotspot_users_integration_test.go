//go:build integration

package mikrotik_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Hotspot Users ────────────────────────────────────────────────────────────

func TestIntegration_HotspotUsers_GetAll(t *testing.T) {
	c := integrationClient(t)

	users, err := c.GetHotspotUsers(context.Background(), "")
	require.NoError(t, err)
	t.Logf("total users: %d", len(users))
	for _, u := range users {
		t.Logf("  name=%-20s profile=%-12s disabled=%v", u.Name, u.Profile, u.Disabled)
	}
}

func TestIntegration_HotspotUsers_GetByProfile(t *testing.T) {
	c := integrationClient(t)

	// Ambil semua user lalu cari profile pertama yang tidak kosong.
	all, err := c.GetHotspotUsers(context.Background(), "")
	require.NoError(t, err)
	if len(all) == 0 {
		t.Skip("no hotspot users found on router")
	}

	var profile string
	for _, u := range all {
		if u.Profile != "" {
			profile = u.Profile
			break
		}
	}
	if profile == "" {
		t.Skip("no user with non-empty profile found — skipping profile filter test")
	}

	users, err := c.GetHotspotUsers(context.Background(), profile)
	require.NoError(t, err)
	t.Logf("users with profile=%q: %d", profile, len(users))
	for _, u := range users {
		assert.Equal(t, profile, u.Profile)
	}
}

func TestIntegration_HotspotUsers_GetCount(t *testing.T) {
	c := integrationClient(t)

	count, err := c.GetHotspotUsersCount(context.Background())
	require.NoError(t, err)
	t.Logf("hotspot users count: %d", count)
	assert.GreaterOrEqual(t, count, 0)
}

func TestIntegration_HotspotUsers_GetByID(t *testing.T) {
	c := integrationClient(t)

	users, err := c.GetHotspotUsers(context.Background(), "")
	require.NoError(t, err)
	if len(users) == 0 {
		t.Skip("no hotspot users found on router")
	}

	first := users[0]
	user, err := c.GetHotspotUserByID(context.Background(), first.ID)
	require.NoError(t, err)
	assert.Equal(t, first.ID, user.ID)
	assert.Equal(t, first.Name, user.Name)
	t.Logf("GetByID ok: name=%s", user.Name)
}

func TestIntegration_HotspotUsers_GetByName(t *testing.T) {
	c := integrationClient(t)

	users, err := c.GetHotspotUsers(context.Background(), "")
	require.NoError(t, err)
	if len(users) == 0 {
		t.Skip("no hotspot users found on router")
	}

	name := users[0].Name
	user, err := c.GetHotspotUserByName(context.Background(), name)
	require.NoError(t, err)
	assert.Equal(t, name, user.Name)
	t.Logf("GetByName ok: name=%s", user.Name)
}

func TestIntegration_HotspotUsers_GetByComment(t *testing.T) {
	c := integrationClient(t)

	users, err := c.GetHotspotUsersByComment(context.Background(), "test-integration")
	require.NoError(t, err)
	t.Logf("users with comment='test-integration': %d", len(users))
}

func TestIntegration_HotspotUsers_AddUpdateRemove(t *testing.T) {
	c := integrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testName := fmt.Sprintf("inttest-%d", time.Now().UnixMilli())
	testComment := "mikhmon-integration-test"

	newUser := &dto.HotspotUser{
		Name:     testName,
		Password: "testpass",
		Profile:  "default",
		Comment:  testComment,
	}

	// Add — RouterOS tidak selalu mengembalikan ID dari /add,
	// jadi kita lookup by name setelah add.
	_, err := c.AddHotspotUser(ctx, newUser)
	require.NoError(t, err)
	t.Logf("added user name=%s", testName)

	created, err := c.GetHotspotUserByName(ctx, testName)
	require.NoError(t, err)
	require.NotNil(t, created, "user must exist after creation")
	require.NotEmpty(t, created.ID, "user must have an ID after creation")
	id := created.ID
	t.Logf("resolved id=%s for name=%s", id, testName)
	assert.Equal(t, testName, created.Name)
	assert.Equal(t, "default", created.Profile)

	t.Cleanup(func() {
		ctxClean, cancelClean := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelClean()
		_ = c.RemoveHotspotUser(ctxClean, id)
	})

	// Update
	updated := &dto.HotspotUser{
		Name:    testName,
		Profile: "default",
		Comment: testComment + "-updated",
	}
	err = c.UpdateHotspotUser(ctx, id, updated)
	require.NoError(t, err)
	t.Logf("updated user id=%s", id)

	after, err := c.GetHotspotUserByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, testComment+"-updated", after.Comment)

	// Remove
	err = c.RemoveHotspotUser(ctx, id)
	require.NoError(t, err)
	t.Logf("removed user id=%s", id)

	gone, err := c.GetHotspotUserByName(ctx, testName)
	assert.NoError(t, err)
	assert.Nil(t, gone, "user should not exist after removal")
}

func TestIntegration_HotspotUsers_RemoveByComment(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	comment := fmt.Sprintf("mikhmon-rmbycomment-%d", time.Now().UnixMilli())

	for i := 0; i < 2; i++ {
		u := &dto.HotspotUser{
			Name:     fmt.Sprintf("inttest-rbc%d-%d", i, time.Now().UnixNano()),
			Password: "x",
			Profile:  "default",
			Comment:  comment,
		}
		_, err := c.AddHotspotUser(ctx, u)
		require.NoError(t, err)
	}

	err := c.RemoveHotspotUsersByComment(ctx, comment)
	require.NoError(t, err)
	t.Logf("RemoveByComment comment=%s", comment)

	remaining, err := c.GetHotspotUsersByComment(ctx, comment)
	require.NoError(t, err)
	assert.Empty(t, remaining, "no users should remain after RemoveByComment")
}

func TestIntegration_HotspotUsers_ResetCounters(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	users, err := c.GetHotspotUsers(ctx, "")
	require.NoError(t, err)
	if len(users) == 0 {
		t.Skip("no hotspot users found")
	}

	err = c.ResetHotspotUserCounters(ctx, users[0].ID)
	require.NoError(t, err)
	t.Logf("reset counters for user=%s", users[0].Name)
}
