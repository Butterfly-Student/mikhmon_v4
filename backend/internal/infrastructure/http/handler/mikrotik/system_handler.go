package mikrotik

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// SystemHandler handles system management
type SystemHandler struct {
	handler.BaseHandler
	systemUC *mikrotik.SystemUseCase
}

// NewSystemHandler creates a new system handler
func NewSystemHandler(systemUC *mikrotik.SystemUseCase, log *zap.Logger) *SystemHandler {
	return &SystemHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("system")},
		systemUC:    systemUC,
	}
}

// GetResources gets system resources
func (h *SystemHandler) GetResources(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	resource, err := h.systemUC.GetResource(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, resource)
}

// GetHealth gets system health
func (h *SystemHandler) GetHealth(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	health, err := h.systemUC.GetHealth(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, health)
}

// GetIdentity gets system identity
func (h *SystemHandler) GetIdentity(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	identity, err := h.systemUC.GetIdentity(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, identity)
}

// GetRouterBoard gets routerboard info
func (h *SystemHandler) GetRouterBoard(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	rb, err := h.systemUC.GetRouterBoardInfo(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, rb)
}

// GetClock gets system clock
func (h *SystemHandler) GetClock(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	clock, err := h.systemUC.GetClock(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, clock)
}

// GetDashboardData gets complete dashboard data
func (h *SystemHandler) GetDashboardData(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	data, err := h.systemUC.GetDashboardData(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, data)
}

// GetRouterStatus returns SystemInfo-shaped response (for frontend compatibility)
func (h *SystemHandler) GetRouterStatus(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	data, err := h.systemUC.GetDashboardData(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

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

	h.Success(c, info)
}

// RegisterRoutes registers system routes
func (h *SystemHandler) RegisterRoutes(r *gin.RouterGroup) {
	system := r.Group("/system")
	{
		system.GET("/resources", h.GetResources)
		system.GET("/health", h.GetHealth)
		system.GET("/identity", h.GetIdentity)
		system.GET("/routerboard", h.GetRouterBoard)
		system.GET("/clock", h.GetClock)
	}

	r.GET("/dashboard", h.GetDashboardData)
	r.GET("/status", h.GetRouterStatus)
}
