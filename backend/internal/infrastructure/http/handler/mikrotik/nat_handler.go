package mikrotik

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// NATHandler handles NAT rules management
type NATHandler struct {
	handler.BaseHandler
	natUC *mikrotik.NATUseCase
}

// NewNATHandler creates a new NAT handler
func NewNATHandler(natUC *mikrotik.NATUseCase, log *zap.Logger) *NATHandler {
	return &NATHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("nat")},
		natUC:       natUC,
	}
}

// GetNATRules retrieves firewall NAT rules
func (h *NATHandler) GetNATRules(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	rules, err := h.natUC.GetNATRules(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, rules)
}

// RegisterRoutes registers NAT routes
func (h *NATHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/nat", h.GetNATRules)
}
