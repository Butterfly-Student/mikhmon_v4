package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/usecase"
	"go.uber.org/zap"
)

// DashboardHandler handles dashboard requests
type DashboardHandler struct {
	BaseHandler
	dashboardUC *usecase.DashboardUseCase
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardUC *usecase.DashboardUseCase, log *zap.Logger) *DashboardHandler {
	return &DashboardHandler{
		BaseHandler: BaseHandler{Log: log.Named("dashboard")},
		dashboardUC: dashboardUC,
	}
}

// GetDashboardData gets dashboard data
func (h *DashboardHandler) GetDashboardData(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	h.Log.Info("GetDashboardData called", zap.Uint64("routerID", routerID))

	data, err := h.dashboardUC.GetDashboardData(c.Request.Context(), uint(routerID))
	if err != nil {
		h.Log.Error("GetDashboardData failed",
			zap.Uint64("routerID", routerID),
			zap.Error(err),
		)
		code := ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("GetDashboardData success",
		zap.Uint64("routerID", routerID),
		zap.Int("activeUsers", data.Stats.ActiveUsers),
		zap.Int("totalUsers", data.Stats.TotalUsers),
	)
	h.Success(c, data)
}

// GetSystemResources gets system resources
func (h *DashboardHandler) GetSystemResources(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	h.Log.Info("GetSystemResources called", zap.Uint64("routerID", routerID))

	data, err := h.dashboardUC.GetResource(c.Request.Context(), uint(routerID), true)
	if err != nil {
		h.Log.Error("GetSystemResources failed",
			zap.Uint64("routerID", routerID),
			zap.Error(err),
		)
		code := ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("GetSystemResources success",
		zap.Uint64("routerID", routerID),
		zap.Int("cpuLoad", data.CpuLoad),
	)
	h.Success(c, data)
}

// GetRouterStatus gets router status (returns SystemInfo shape for frontend)
func (h *DashboardHandler) GetRouterStatus(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	h.Log.Info("GetRouterStatus called", zap.Uint64("routerID", routerID))

	data, err := h.dashboardUC.GetDashboardData(c.Request.Context(), uint(routerID))
	if err != nil {
		h.Log.Error("GetRouterStatus failed",
			zap.Uint64("routerID", routerID),
			zap.Error(err),
		)
		code := ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	// Return SystemInfo-shaped response to match frontend TypeScript interface
	type SystemInfo struct {
		Uptime    string `json:"uptime"`
		BoardName string `json:"boardName"`
		Model     string `json:"model"`
		Version   string `json:"version"`
	}

	info := SystemInfo{}
	if data.Resource != nil {
		info.Uptime = data.Resource.Uptime
		info.BoardName = data.Resource.BoardName
		info.Version = data.Resource.Version
	}
	if data.RouterBoard != nil {
		info.Model = data.RouterBoard.Model
	}

	h.Log.Info("GetRouterStatus success",
		zap.Uint64("routerID", routerID),
		zap.String("model", info.Model),
		zap.String("version", info.Version),
	)
	h.Success(c, info)
}

// RegisterRoutes registers dashboard routes
func (h *DashboardHandler) RegisterRoutes(r *gin.RouterGroup) {
	dashboard := r.Group("/dashboard")
	{
		dashboard.GET("/:router_id", h.GetDashboardData)
		dashboard.GET("/:router_id/resources", h.GetSystemResources)
		dashboard.GET("/:router_id/status", h.GetRouterStatus)
	}
}

// isMikrotikTimeout checks if the error is a context deadline / timeout error.
func isMikrotikTimeout(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "deadline exceeded") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "timed out")
}
