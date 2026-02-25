package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/usecase"
	"go.uber.org/zap"
)

// VoucherHandler handles voucher generation
type VoucherHandler struct {
	BaseHandler
	voucherUC *usecase.VoucherUseCase
}

// NewVoucherHandler creates a new voucher handler
func NewVoucherHandler(voucherUC *usecase.VoucherUseCase, log *zap.Logger) *VoucherHandler {
	return &VoucherHandler{
		BaseHandler: BaseHandler{Log: log.Named("voucher")},
		voucherUC:   voucherUC,
	}
}

// GenerateVouchers generates batch vouchers
func (h *VoucherHandler) GenerateVouchers(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	var req dto.VoucherGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.voucherUC.GenerateVouchers(c.Request.Context(), uint(routerID), req)
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Success(c, result)
}

// GetVouchers retrieves vouchers by comment
func (h *VoucherHandler) GetVouchers(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	comment := c.Query("comment")

	vouchers, err := h.voucherUC.GetVouchersByComment(c.Request.Context(), uint(routerID), comment)
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Success(c, vouchers)
}

// DeleteVouchers deletes vouchers by comment
func (h *VoucherHandler) DeleteVouchers(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	comment := c.Query("comment")

	if err := h.voucherUC.DeleteVouchersByComment(c.Request.Context(), uint(routerID), comment); err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.SuccessWithMessage(c, "Vouchers deleted successfully", nil)
}

// RegisterRoutes registers voucher routes
func (h *VoucherHandler) RegisterRoutes(r *gin.RouterGroup) {
	voucher := r.Group("/vouchers")
	{
		voucher.POST("/:router_id/generate", h.GenerateVouchers)
		voucher.GET("/:router_id", h.GetVouchers)
		voucher.DELETE("/:router_id", h.DeleteVouchers)
	}
}
