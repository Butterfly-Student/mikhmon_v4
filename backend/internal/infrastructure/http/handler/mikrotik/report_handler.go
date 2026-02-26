package mikrotik

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// ReportHandler handles report requests
type ReportHandler struct {
	handler.BaseHandler
	reportUC *mikrotik.ReportUseCase
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportUC *mikrotik.ReportUseCase, log *zap.Logger) *ReportHandler {
	return &ReportHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("report")},
		reportUC:    reportUC,
	}
}

// GetSalesReport gets sales report
func (h *ReportHandler) GetSalesReport(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	owner, day := resolveReportQuery(c.Query("owner"), c.Query("month"), c.Query("year"), c.Query("date"))

	var (
		reports []*dto.SalesReport
		err     error
	)
	if day != "" {
		reports, err = h.getSalesByDay(c, uint(routerID), day)
	} else {
		reports, err = h.getSalesByOwner(c, uint(routerID), owner)
	}
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Success(c, reports)
}

// GetReportSummary gets summary for dashboard
func (h *ReportHandler) GetReportSummary(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	owner, day := resolveReportQuery(c.Query("owner"), c.Query("month"), c.Query("year"), c.Query("date"))

	var (
		reports []*dto.SalesReport
		err     error
	)
	if day != "" {
		reports, err = h.getSalesByDay(c, uint(routerID), day)
	} else {
		reports, err = h.getSalesByOwner(c, uint(routerID), owner)
	}
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	summary := h.reportUC.CalculateSummary(reports)
	avg := 0.0
	if summary.TotalVouchers > 0 {
		avg = summary.TotalAmount / float64(summary.TotalVouchers)
	}
	h.Success(c, gin.H{
		"totalSales":        summary.TotalAmount,
		"totalTransactions": summary.TotalVouchers,
		"averageTicket":     avg,
		"byProfile":         summary.ByProfile,
	})
}

// ExportToCSV exports reports to CSV
func (h *ReportHandler) ExportToCSV(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	owner, day := resolveReportQuery(c.Query("owner"), c.Query("month"), c.Query("year"), c.Query("date"))

	var (
		reports []*dto.SalesReport
		err     error
	)
	if day != "" {
		reports, err = h.getSalesByDay(c, uint(routerID), day)
	} else {
		reports, err = h.getSalesByOwner(c, uint(routerID), owner)
	}
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{
		"date", "time", "username", "price", "ipAddress", "macAddress", "validity", "profile", "comment",
	})
	for _, r := range reports {
		_ = w.Write([]string{
			r.Date,
			r.Time,
			r.Username,
			fmt.Sprintf("%.0f", r.Price),
			r.IPAddress,
			r.MACAddress,
			r.Validity,
			r.Profile,
			r.VoucherComment,
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	filename := "sales-report.csv"
	if owner != "" {
		filename = fmt.Sprintf("sales-report-%s.csv", owner)
	}
	if day != "" {
		filename = fmt.Sprintf("sales-report-%s.csv", strings.ReplaceAll(day, "/", "-"))
	}
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.String(http.StatusOK, buf.String())
}

// RegisterRoutes registers report routes
func (h *ReportHandler) RegisterRoutes(r *gin.RouterGroup) {
	reports := r.Group("/reports")
	{
		reports.GET("/sales", h.GetSalesReport)
		reports.GET("/summary", h.GetReportSummary)
		reports.GET("/export", h.ExportToCSV)
	}
}

func (h *ReportHandler) getSalesByOwner(c *gin.Context, routerID uint, owner string) ([]*dto.SalesReport, error) {
	return h.reportUC.GetSalesReport(c.Request.Context(), routerID, owner)
}

func (h *ReportHandler) getSalesByDay(c *gin.Context, routerID uint, day string) ([]*dto.SalesReport, error) {
	return h.reportUC.GetSalesReportByDay(c.Request.Context(), routerID, day)
}

func resolveReportQuery(owner, month, year, date string) (resolvedOwner, resolvedDay string) {
	if day := normalizeDay(date); day != "" {
		return "", day
	}

	if owner != "" {
		return strings.ToLower(owner), ""
	}

	if month != "" && year != "" {
		if m, err := strconv.Atoi(month); err == nil && m >= 1 && m <= 12 {
			return fmt.Sprintf("%s%s", monthAbbr(time.Month(m)), year), ""
		}
	}

	return "", ""
}

func normalizeDay(value string) string {
	if value == "" {
		return ""
	}
	if strings.Count(value, "/") == 2 && len(value) >= 8 && len(value) <= 12 {
		return strings.ToLower(value)
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s/%d/%d", monthAbbr(t.Month()), t.Day(), t.Year())
}

func monthAbbr(m time.Month) string {
	switch m {
	case time.January:
		return "jan"
	case time.February:
		return "feb"
	case time.March:
		return "mar"
	case time.April:
		return "apr"
	case time.May:
		return "may"
	case time.June:
		return "jun"
	case time.July:
		return "jul"
	case time.August:
		return "aug"
	case time.September:
		return "sep"
	case time.October:
		return "oct"
	case time.November:
		return "nov"
	case time.December:
		return "dec"
	default:
		return "jan"
	}
}
