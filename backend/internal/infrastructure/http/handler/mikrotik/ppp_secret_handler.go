package mikrotik

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	mikrotikUC "github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// PPPSecretHandler handles PPP secret endpoints
type PPPSecretHandler struct {
	handler.BaseHandler
	uc *mikrotikUC.PPPSecretUseCase
}

// NewPPPSecretHandler creates a new PPP secret handler
func NewPPPSecretHandler(uc *mikrotikUC.PPPSecretUseCase, log *zap.Logger) *PPPSecretHandler {
	return &PPPSecretHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("ppp-secret")},
		uc:          uc,
	}
}

// GetSecrets handles GET /ppp/secrets
func (h *PPPSecretHandler) GetSecrets(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	profile := c.Query("profile")

	secrets, err := h.uc.GetSecrets(c.Request.Context(), uint(routerID), profile)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, secrets)
}

// GetSecretByName handles GET /ppp/secrets/by-name/:name
func (h *PPPSecretHandler) GetSecretByName(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	name := c.Param("name")

	secret, err := h.uc.GetSecretByName(c.Request.Context(), uint(routerID), name)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, secret)
}

// GetSecretByID handles GET /ppp/secrets/:id
func (h *PPPSecretHandler) GetSecretByID(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	secret, err := h.uc.GetSecretByID(c.Request.Context(), uint(routerID), id)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, secret)
}

// AddSecret handles POST /ppp/secrets
func (h *PPPSecretHandler) AddSecret(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	var req dto.PPPSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}

	if err := h.uc.AddSecret(c.Request.Context(), uint(routerID), &req); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// UpdateSecret handles PUT /ppp/secrets/:id
func (h *PPPSecretHandler) UpdateSecret(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	var req dto.PPPSecretUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}

	if err := h.uc.UpdateSecret(c.Request.Context(), uint(routerID), id, &req); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// RemoveSecret handles DELETE /ppp/secrets/:id
func (h *PPPSecretHandler) RemoveSecret(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	if err := h.uc.RemoveSecret(c.Request.Context(), uint(routerID), id); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// DisableSecret handles POST /ppp/secrets/:id/disable
func (h *PPPSecretHandler) DisableSecret(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	if err := h.uc.DisableSecret(c.Request.Context(), uint(routerID), id); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// EnableSecret handles POST /ppp/secrets/:id/enable
func (h *PPPSecretHandler) EnableSecret(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	if err := h.uc.EnableSecret(c.Request.Context(), uint(routerID), id); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// RegisterRoutes registers PPP secret routes
func (h *PPPSecretHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/ppp/secrets", h.GetSecrets)
	r.GET("/ppp/secrets/by-name/:name", h.GetSecretByName)
	r.GET("/ppp/secrets/:id", h.GetSecretByID)
	r.POST("/ppp/secrets", h.AddSecret)
	r.PUT("/ppp/secrets/:id", h.UpdateSecret)
	r.DELETE("/ppp/secrets/:id", h.RemoveSecret)
	r.POST("/ppp/secrets/:id/disable", h.DisableSecret)
	r.POST("/ppp/secrets/:id/enable", h.EnableSecret)
}
