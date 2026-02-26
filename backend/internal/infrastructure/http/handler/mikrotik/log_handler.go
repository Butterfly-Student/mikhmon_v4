package mikrotik

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// LogHandler handles log management
type LogHandler struct {
	handler.BaseHandler
	logUC *mikrotik.LogUseCase
}

// NewLogHandler creates a new log handler
func NewLogHandler(logUC *mikrotik.LogUseCase, log *zap.Logger) *LogHandler {
	return &LogHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("log")},
		logUC:       logUC,
	}
}

// GetHotspotLogs gets hotspot logs
func (h *LogHandler) GetHotspotLogs(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	limit := 100
	if raw := c.Query("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	logs, err := h.logUC.GetHotspotLogs(c.Request.Context(), uint(routerID), limit)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, logs)
}

// RegisterRoutes registers log routes
func (h *LogHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/logs", h.GetHotspotLogs)
}
