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

// PPPProfileHandler handles PPP profile endpoints
type PPPProfileHandler struct {
	handler.BaseHandler
	uc *mikrotikUC.PPPProfileUseCase
}

// NewPPPProfileHandler creates a new PPP profile handler
func NewPPPProfileHandler(uc *mikrotikUC.PPPProfileUseCase, log *zap.Logger) *PPPProfileHandler {
	return &PPPProfileHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("ppp-profile")},
		uc:          uc,
	}
}

// GetProfiles handles GET /ppp/profiles
func (h *PPPProfileHandler) GetProfiles(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	profiles, err := h.uc.GetProfiles(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, profiles)
}

// GetProfileByName handles GET /ppp/profiles/by-name/:name
func (h *PPPProfileHandler) GetProfileByName(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	name := c.Param("name")

	profile, err := h.uc.GetProfileByName(c.Request.Context(), uint(routerID), name)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, profile)
}

// GetProfileByID handles GET /ppp/profiles/:id
func (h *PPPProfileHandler) GetProfileByID(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	profile, err := h.uc.GetProfileByID(c.Request.Context(), uint(routerID), id)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, profile)
}

// AddProfile handles POST /ppp/profiles
func (h *PPPProfileHandler) AddProfile(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	var req dto.PPPProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}

	if err := h.uc.AddProfile(c.Request.Context(), uint(routerID), &req); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// UpdateProfile handles PUT /ppp/profiles/:id
func (h *PPPProfileHandler) UpdateProfile(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	var req dto.PPPProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}

	if err := h.uc.UpdateProfile(c.Request.Context(), uint(routerID), id, &req); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// RemoveProfile handles DELETE /ppp/profiles/:id
func (h *PPPProfileHandler) RemoveProfile(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	if err := h.uc.RemoveProfile(c.Request.Context(), uint(routerID), id); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// DisableProfile handles POST /ppp/profiles/:id/disable
func (h *PPPProfileHandler) DisableProfile(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	if err := h.uc.DisableProfile(c.Request.Context(), uint(routerID), id); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// EnableProfile handles POST /ppp/profiles/:id/enable
func (h *PPPProfileHandler) EnableProfile(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	id := c.Param("id")

	if err := h.uc.EnableProfile(c.Request.Context(), uint(routerID), id); err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, nil)
}

// RegisterRoutes registers PPP profile routes
func (h *PPPProfileHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/ppp/profiles", h.GetProfiles)
	r.GET("/ppp/profiles/by-name/:name", h.GetProfileByName)
	r.GET("/ppp/profiles/:id", h.GetProfileByID)
	r.POST("/ppp/profiles", h.AddProfile)
	r.PUT("/ppp/profiles/:id", h.UpdateProfile)
	r.DELETE("/ppp/profiles/:id", h.RemoveProfile)
	r.POST("/ppp/profiles/:id/disable", h.DisableProfile)
	r.POST("/ppp/profiles/:id/enable", h.EnableProfile)
}
