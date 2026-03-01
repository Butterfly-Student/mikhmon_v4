package mikrotik

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// GetSalesReports retrieves sales reports by owner (month).
// Reports are stored in /system/script with owner field as month identifier.
func (c *Client) GetSalesReports(ctx context.Context, owner string) ([]*dto.SalesReport, error) {
	args := []string{"/system/script/print"}
	if owner != "" {
		args = append(args, "?owner="+owner)
	}

	reply, err := c.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}

	reports := make([]*dto.SalesReport, 0, len(reply.Re))
	for _, re := range reply.Re {
		if report := parseSalesReport(re.Map); report != nil {
			reports = append(reports, report)
		}
	}

	return reports, nil
}

// GetSalesReportsByDay retrieves sales reports by day (source field).
func (c *Client) GetSalesReportsByDay(ctx context.Context, day string) ([]*dto.SalesReport, error) {
	args := []string{"/system/script/print"}
	if day != "" {
		args = append(args, "?source="+day)
	}

	reply, err := c.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}

	reports := make([]*dto.SalesReport, 0, len(reply.Re))
	for _, re := range reply.Re {
		if report := parseSalesReport(re.Map); report != nil {
			reports = append(reports, report)
		}
	}

	return reports, nil
}

// parseSalesReport parses a sales report from a script entry.
// Name format: $date-|-$time-|-$user-|-$price-|-$address-|-$mac-|-$validity-|-$profile-|-$comment
func parseSalesReport(data map[string]string) *dto.SalesReport {
	name := data["name"]
	if !strings.Contains(name, "-|-") {
		return nil
	}

	parts := strings.Split(name, "-|-")
	if len(parts) < 8 {
		return nil
	}

	price, _ := strconv.ParseFloat(parts[3], 64)

	return &dto.SalesReport{
		ID:             data[".id"],
		Name:           name,
		Owner:          data["owner"],
		Source:         data["source"],
		Comment:        data["comment"],
		DontReq:        data["dont-require-permissions"],
		RunCount:       data["run-count"],
		CopyOf:         data["copy-of"],
		Date:           parts[0],
		Time:           parts[1],
		Username:       parts[2],
		Price:          price,
		IPAddress:      parts[4],
		MACAddress:     parts[5],
		Validity:       parts[6],
		Profile:        parts[7],
		VoucherComment: getPart(parts, 8),
	}
}

// getPart safely retrieves an element from a slice.
func getPart(parts []string, index int) string {
	if index < len(parts) {
		return parts[index]
	}
	return ""
}

// AddSalesReport adds a sales report entry.
func (c *Client) AddSalesReport(ctx context.Context, report *dto.SalesReport) error {
	name := fmt.Sprintf("%s-|-%s-|-%s-|-%.0f-|-%s-|-%s-|-%s-|-%s",
		report.Date,
		report.Time,
		report.Username,
		report.Price,
		report.IPAddress,
		report.MACAddress,
		report.Validity,
		report.Profile,
	)

	if report.VoucherComment != "" {
		name = name + "-|-" + report.VoucherComment
	}

	_, err := c.RunContext(ctx,
		"/system/script/add",
		"=name="+name,
		"=owner="+report.Owner,
		"=source="+report.Source,
		"=comment=mikhmon",
	)
	return err
}
