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

// ─── Sales Reports ────────────────────────────────────────────────────────────

func TestIntegration_Reports_GetAll(t *testing.T) {
	c := integrationClient(t)

	reports, err := c.GetSalesReports(context.Background(), "")
	require.NoError(t, err)
	assert.NotNil(t, reports)
	t.Logf("total sales reports (no filter): %d", len(reports))
	for _, r := range reports {
		t.Logf("  date=%s time=%s user=%-20s price=%.0f profile=%s",
			r.Date, r.Time, r.Username, r.Price, r.Profile)
	}
}

func TestIntegration_Reports_GetByOwner(t *testing.T) {
	c := integrationClient(t)

	all, err := c.GetSalesReports(context.Background(), "")
	require.NoError(t, err)
	if len(all) == 0 {
		t.Skip("no sales reports found — skipping owner filter test")
	}

	owner := all[0].Owner
	reports, err := c.GetSalesReports(context.Background(), owner)
	require.NoError(t, err)
	t.Logf("reports with owner=%q: %d", owner, len(reports))
	for _, r := range reports {
		assert.Equal(t, owner, r.Owner)
	}
}

func TestIntegration_Reports_GetByDay(t *testing.T) {
	c := integrationClient(t)

	all, err := c.GetSalesReports(context.Background(), "")
	require.NoError(t, err)
	if len(all) == 0 {
		t.Skip("no sales reports found — skipping day filter test")
	}

	day := all[0].Source
	if day == "" {
		t.Skip("first report has no source field — skipping")
	}

	reports, err := c.GetSalesReportsByDay(context.Background(), day)
	require.NoError(t, err)
	t.Logf("reports with day=%q: %d", day, len(reports))
}

func TestIntegration_Reports_AddAndVerify(t *testing.T) {
	c := integrationClient(t)
	ctx := context.Background()

	now := time.Now()
	date := fmt.Sprintf("%s/%02d/%d",
		monthAbbr(int(now.Month())), now.Day(), now.Year())
	timeStr := now.Format("15:04:05")
	owner := fmt.Sprintf("test-%d", now.Unix())

	report := &dto.SalesReport{
		Date:       date,
		Time:       timeStr,
		Username:   "inttest-user",
		Price:      5000,
		IPAddress:  "10.0.0.100",
		MACAddress: "AA:BB:CC:DD:EE:FF",
		Validity:   "1d",
		Profile:    "default",
		Owner:      owner,
		Source:     date,
		Comment:    "mikhmon",
	}

	err := c.AddSalesReport(ctx, report)
	require.NoError(t, err)
	t.Logf("AddSalesReport succeeded for owner=%s", owner)

	reports, err := c.GetSalesReports(ctx, owner)
	require.NoError(t, err)
	require.NotEmpty(t, reports, "report should be retrievable by owner")
	assert.Equal(t, "inttest-user", reports[0].Username)
	assert.InDelta(t, float64(5000), reports[0].Price, 0.01)
	t.Logf("verified: date=%s user=%s price=%.0f", reports[0].Date, reports[0].Username, reports[0].Price)
}

// monthAbbr converts 1-based month number to MikroTik 3-letter abbreviation.
func monthAbbr(m int) string {
	months := [...]string{
		"jan", "feb", "mar", "apr", "may", "jun",
		"jul", "aug", "sep", "oct", "nov", "dec",
	}
	if m < 1 || m > 12 {
		return "jan"
	}
	return months[m-1]
}
