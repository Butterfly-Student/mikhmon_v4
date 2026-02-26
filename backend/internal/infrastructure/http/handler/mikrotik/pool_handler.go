package mikrotik

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// PoolHandler handles IP pool management
type PoolHandler struct {
	handler.BaseHandler
	poolUC *mikrotik.PoolUseCase
}

// NewPoolHandler creates a new pool handler
func NewPoolHandler(poolUC *mikrotik.PoolUseCase, log *zap.Logger) *PoolHandler {
	return &PoolHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("pool")},
		poolUC:      poolUC,
	}
}

// GetAddressPools retrieves address pools
func (h *PoolHandler) GetAddressPools(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	pools, err := h.poolUC.GetAddressPools(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, pools)
}

// RegisterRoutes registers pool routes
func (h *PoolHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/pools", h.GetAddressPools)
}
