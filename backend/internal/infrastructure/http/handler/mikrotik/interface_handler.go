package mikrotik

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// InterfaceHandler handles interface management
type InterfaceHandler struct {
	handler.BaseHandler
	interfaceUC *mikrotik.InterfaceUseCase
}

// NewInterfaceHandler creates a new interface handler
func NewInterfaceHandler(interfaceUC *mikrotik.InterfaceUseCase, log *zap.Logger) *InterfaceHandler {
	return &InterfaceHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("interface")},
		interfaceUC: interfaceUC,
	}
}

// GetInterfaces gets network interfaces
func (h *InterfaceHandler) GetInterfaces(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	interfaces, err := h.interfaceUC.GetInterfaces(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, interfaces)
}

// GetTraffic gets traffic stats for an interface
func (h *InterfaceHandler) GetTraffic(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	iface := c.Param("name")
	if iface == "" {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, "interface name is required")
		return
	}

	data, err := h.interfaceUC.GetTraffic(c.Request.Context(), uint(routerID), iface)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, data)
}

// RegisterRoutes registers interface routes
func (h *InterfaceHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/interfaces", h.GetInterfaces)
	r.GET("/interfaces/:name/traffic", h.GetTraffic)
}
