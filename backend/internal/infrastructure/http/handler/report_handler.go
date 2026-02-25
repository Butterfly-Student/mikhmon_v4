package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/usecase"
	"go.uber.org/zap"
)

// ReportHandler handles report requests
type ReportHandler struct {
	BaseHandler
	reportUC *usecase.ReportUseCase
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportUC *usecase.ReportUseCase, log *zap.Logger) *ReportHandler {
	return &ReportHandler{
		BaseHandler: BaseHandler{Log: log.Named("report")},
		reportUC:    reportUC,
	}
}

// GetSalesReport gets sales report
func (h *ReportHandler) GetSalesReport(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	owner := c.Query("owner")

	reports, err := h.reportUC.GetSalesReport(c.Request.Context(), uint(routerID), owner, true)
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Success(c, reports)
}

// GetReportSummary gets summary for dashboard
func (h *ReportHandler) GetReportSummary(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	owner := c.DefaultQuery("owner", "")

	reports, err := h.reportUC.GetSalesReport(c.Request.Context(), uint(routerID), owner, true)
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	summary := h.reportUC.CalculateSummary(reports)
	h.Success(c, summary)
}

// ExportToCSV exports reports to CSV
func (h *ReportHandler) ExportToCSV(c *gin.Context) {
	// TODO: Implement CSV export
	h.Error(c, http.StatusNotImplemented, "CSV export not implemented yet")
}

// RegisterRoutes registers report routes
func (h *ReportHandler) RegisterRoutes(r *gin.RouterGroup) {
	report := r.Group("/reports")
	{
		report.GET("/:router_id/sales", h.GetSalesReport)
		report.GET("/:router_id/summary", h.GetReportSummary)
		report.GET("/:router_id/export", h.ExportToCSV)
	}
}
