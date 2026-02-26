package mikrotik

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// VoucherHandler handles voucher generation
type VoucherHandler struct {
	handler.BaseHandler
	voucherUC *mikrotik.VoucherUseCase
}

// NewVoucherHandler creates a new voucher handler
func NewVoucherHandler(voucherUC *mikrotik.VoucherUseCase, log *zap.Logger) *VoucherHandler {
	return &VoucherHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("voucher")},
		voucherUC:   voucherUC,
	}
}

type CacheVoucherRequest struct {
	User     string `json:"user" binding:"required"`
	GComment string `json:"gcomment"`
	Gencode  string `json:"gencode" binding:"required"`
	Qty      int    `json:"qty"`
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

// CacheVouchers caches generated vouchers using legacy comment format
func (h *VoucherHandler) CacheVouchers(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	var req CacheVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}

	count, comment, users, err := h.voucherUC.CacheGeneratedVouchers(
		c.Request.Context(),
		uint(routerID),
		req.User,
		req.Gencode,
		req.GComment,
	)
	if err != nil {
		h.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.Success(c, gin.H{
		"count":   count,
		"comment": comment,
		"users":   users,
	})
}

// RegisterRoutes registers voucher routes
func (h *VoucherHandler) RegisterRoutes(r *gin.RouterGroup) {
	vouchers := r.Group("/vouchers")
	{
		vouchers.POST("/generate", h.GenerateVouchers)
		vouchers.POST("/cache", h.CacheVouchers)
		vouchers.GET("", h.GetVouchers)
		vouchers.DELETE("", h.DeleteVouchers)
	}
}
