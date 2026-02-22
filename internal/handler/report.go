package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"mikhmon_v4/config"
	ros "mikhmon_v4/internal/routeros"
)

// ReportHandler serves sales reports and live income data.
type ReportHandler struct {
	store sessions.Store
	pool  *ros.Pool
}

func NewReportHandler(store sessions.Store, pool *ros.Pool) *ReportHandler {
	return &ReportHandler{store: store, pool: pool}
}

// ReportPage renders the sales report page.
func (h *ReportHandler) ReportPage(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.Redirect(http.StatusFound, "/admin/settings")
		return
	}
	c.HTML(http.StatusOK, "web/templates/report.html", gin.H{
		"Title":   "Report — " + router.HotspotName,
		"Session": session,
		"Router":  router,
		"Active":  "report",
	})
}

// GetReports returns sales records filtered by date (source field).
// Query: day=jan/01/2025
func (h *ReportHandler) GetReports(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	day := c.Query("day")
	if day == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "day param required"})
		return
	}

	reports, err := ros.RunArgs(client, "/system/script/print", "?source="+day, "?comment=mikhmon")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse each script name: date-|-time-|-user-|-price-|-address-|-mac-|-validity-|-profile-|-comment
	var parsed []map[string]string
	for _, r := range reports {
		parts := strings.Split(r["name"], "-|-")
		entry := map[string]string{
			"date":     safeIdx(parts, 0),
			"time":     safeIdx(parts, 1),
			"user":     safeIdx(parts, 2),
			"price":    safeIdx(parts, 3),
			"address":  safeIdx(parts, 4),
			"mac":      safeIdx(parts, 5),
			"validity": safeIdx(parts, 6),
			"profile":  safeIdx(parts, 7),
			"comment":  safeIdx(parts, 8),
			"source":   r["source"],
			"owner":    r["owner"],
		}
		parsed = append(parsed, entry)
	}

	if parsed == nil {
		parsed = []map[string]string{}
	}
	c.JSON(http.StatusOK, parsed)
}

// LiveReportPage renders the live income report page.
func (h *ReportHandler) LiveReportPage(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.Redirect(http.StatusFound, "/admin/settings")
		return
	}
	c.HTML(http.StatusOK, "web/templates/live_report.html", gin.H{
		"Title":   "Live Report — " + router.HotspotName,
		"Session": session,
		"Router":  router,
		"Active":  "livereport",
	})
}

// GetLiveReport returns aggregated income data for a given month.
// Query: month=Jan2025
func (h *ReportHandler) GetLiveReport(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	month := c.Query("month")
	if month == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month param required"})
		return
	}

	reports, err := ros.RunArgs(client, "/system/script/print", "?owner="+month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if reports == nil {
		reports = []map[string]string{}
	}
	c.JSON(http.StatusOK, reports)
}

func safeIdx(s []string, i int) string {
	if i < len(s) {
		return s[i]
	}
	return ""
}
