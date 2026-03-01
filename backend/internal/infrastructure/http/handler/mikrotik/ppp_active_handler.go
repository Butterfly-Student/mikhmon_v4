package mikrotik

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	mikrotikUC "github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// PPPActiveHandler handles PPP active session endpoints
type PPPActiveHandler struct {
	handler.BaseHandler
	uc *mikrotikUC.PPPActiveUseCase
}

// NewPPPActiveHandler creates a new PPP active handler
func NewPPPActiveHandler(uc *mikrotikUC.PPPActiveUseCase, log *zap.Logger) *PPPActiveHandler {
	return &PPPActiveHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("ppp-active")},
		uc:          uc,
	}
}

// GetActive handles GET /ppp/active
func (h *PPPActiveHandler) GetActive(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	service := c.Query("service")

	active, err := h.uc.GetActive(c.Request.Context(), uint(routerID), service)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, active)
}

// GetActiveByID handles GET /ppp/active/:id
func (h *PPPActiveHandler) GetActiveByID(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	active, err := h.uc.GetActiveByID(c.Request.Context(), uint(routerID), id)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, active)
}

// RemoveActive handles DELETE /ppp/active/:id
func (h *PPPActiveHandler) RemoveActive(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	if err := h.uc.RemoveActive(c.Request.Context(), uint(routerID), id); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// RegisterRoutes registers PPP active routes
func (h *PPPActiveHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/ppp/active", h.GetActive)
	r.GET("/ppp/active/:id", h.GetActiveByID)
	r.DELETE("/ppp/active/:id", h.RemoveActive)
}
