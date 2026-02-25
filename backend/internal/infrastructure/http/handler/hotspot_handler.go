package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/usecase"
	"go.uber.org/zap"
)

// HotspotHandler handles hotspot management
type HotspotHandler struct {
	BaseHandler
	hotspotUC *usecase.HotspotUseCase
}

// NewHotspotHandler creates a new hotspot handler
func NewHotspotHandler(hotspotUC *usecase.HotspotUseCase, log *zap.Logger) *HotspotHandler {
	return &HotspotHandler{
		BaseHandler: BaseHandler{Log: log.Named("hotspot")},
		hotspotUC:   hotspotUC,
	}
}

// GetUsers lists all hotspot users
func (h *HotspotHandler) GetUsers(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	profile := c.Query("profile")

	h.Log.Info("GetUsers called", zap.Uint64("routerID", routerID), zap.String("profile", profile))

	users, err := h.hotspotUC.GetUsers(c.Request.Context(), uint(routerID), profile, true)
	if err != nil {
		h.Log.Error("GetUsers failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("GetUsers success", zap.Uint64("routerID", routerID), zap.Int("count", len(users)))
	h.Success(c, users)
}

// GetUser gets a specific user
func (h *HotspotHandler) GetUser(c *gin.Context) {
	h.ErrorWithCode(c, http.StatusNotImplemented, ErrCodeInternal, "Not implemented")
}

// CreateUser creates a new hotspot user
func (h *HotspotHandler) CreateUser(c *gin.Context) {
	routerID, err := strconv.ParseUint(c.Param("router_id"), 10, 32)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, ErrCodeValidation, "Invalid router ID")
		return
	}

	var req dto.AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}
	req.RouterID = uint(routerID)

	h.Log.Info("CreateUser called", zap.Uint64("routerID", routerID), zap.String("name", req.Name))

	if err := h.hotspotUC.AddUser(c.Request.Context(), &req); err != nil {
		h.Log.Error("CreateUser failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("CreateUser success", zap.Uint64("routerID", routerID), zap.String("name", req.Name))
	h.Created(c, gin.H{"message": "User created successfully"})
}

// UpdateUser updates a hotspot user
func (h *HotspotHandler) UpdateUser(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	userID := c.Param("id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	h.Log.Info("UpdateUser called", zap.Uint64("routerID", routerID), zap.String("userID", userID))

	if err := h.hotspotUC.UpdateUser(c.Request.Context(), uint(routerID), userID, &req); err != nil {
		h.Log.Error("UpdateUser failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, gin.H{"message": "User updated successfully"})
}

// DeleteUser deletes a hotspot user
func (h *HotspotHandler) DeleteUser(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	userID := c.Param("id")

	h.Log.Info("DeleteUser called", zap.Uint64("routerID", routerID), zap.String("userID", userID))

	if err := h.hotspotUC.RemoveUser(c.Request.Context(), uint(routerID), userID); err != nil {
		h.Log.Error("DeleteUser failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.SuccessWithMessage(c, "User deleted successfully", nil)
}

// GetProfiles lists user profiles
func (h *HotspotHandler) GetProfiles(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	h.Log.Info("GetProfiles called", zap.Uint64("routerID", routerID))

	profiles, err := h.hotspotUC.GetProfiles(c.Request.Context(), uint(routerID), true)
	if err != nil {
		h.Log.Error("GetProfiles failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("GetProfiles success", zap.Uint64("routerID", routerID), zap.Int("count", len(profiles)))
	h.Success(c, profiles)
}

// CreateProfile creates a user profile
func (h *HotspotHandler) CreateProfile(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	var req dto.ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	h.Log.Info("CreateProfile called", zap.Uint64("routerID", routerID), zap.String("name", req.Name))

	if err := h.hotspotUC.AddProfile(c.Request.Context(), uint(routerID), &req); err != nil {
		h.Log.Error("CreateProfile failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Created(c, gin.H{"message": "Profile created successfully"})
}

// UpdateProfile updates a user profile
func (h *HotspotHandler) UpdateProfile(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	profileID := c.Param("id")

	var req dto.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, ErrCodeValidation, err.Error())
		return
	}

	h.Log.Info("UpdateProfile called", zap.Uint64("routerID", routerID), zap.String("profileID", profileID))

	if err := h.hotspotUC.UpdateProfile(c.Request.Context(), uint(routerID), profileID, &req); err != nil {
		h.Log.Error("UpdateProfile failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, gin.H{"message": "Profile updated successfully"})
}

// DeleteProfile deletes a user profile
func (h *HotspotHandler) DeleteProfile(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	profileID := c.Param("id")

	h.Log.Info("DeleteProfile called", zap.Uint64("routerID", routerID), zap.String("profileID", profileID))

	if err := h.hotspotUC.RemoveProfile(c.Request.Context(), uint(routerID), profileID); err != nil {
		h.Log.Error("DeleteProfile failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.SuccessWithMessage(c, "Profile deleted successfully", nil)
}

// GetActiveUsers lists active sessions
func (h *HotspotHandler) GetActiveUsers(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	h.Log.Info("GetActiveUsers called", zap.Uint64("routerID", routerID))

	users, err := h.hotspotUC.GetActive(c.Request.Context(), uint(routerID), true)
	if err != nil {
		h.Log.Error("GetActiveUsers failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("GetActiveUsers success", zap.Uint64("routerID", routerID), zap.Int("count", len(users)))
	h.Success(c, users)
}

// GetHosts lists hotspot hosts
func (h *HotspotHandler) GetHosts(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	h.Log.Info("GetHosts called", zap.Uint64("routerID", routerID))

	hosts, err := h.hotspotUC.GetHosts(c.Request.Context(), uint(routerID), true)
	if err != nil {
		h.Log.Error("GetHosts failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("GetHosts success", zap.Uint64("routerID", routerID), zap.Int("count", len(hosts)))
	h.Success(c, hosts)
}

// RegisterRoutes registers hotspot routes
func (h *HotspotHandler) RegisterRoutes(r *gin.RouterGroup) {
	hotspot := r.Group("/hotspot")
	{
		// Users
		hotspot.GET("/:router_id/users", h.GetUsers)
		hotspot.GET("/:router_id/users/:id", h.GetUser)
		hotspot.POST("/:router_id/users", h.CreateUser)
		hotspot.PUT("/:router_id/users/:id", h.UpdateUser)
		hotspot.DELETE("/:router_id/users/:id", h.DeleteUser)

		// Profiles
		hotspot.GET("/:router_id/profiles", h.GetProfiles)
		hotspot.POST("/:router_id/profiles", h.CreateProfile)
		hotspot.PUT("/:router_id/profiles/:id", h.UpdateProfile)
		hotspot.DELETE("/:router_id/profiles/:id", h.DeleteProfile)

		// Active
		hotspot.GET("/:router_id/active", h.GetActiveUsers)

		// Hosts
		hotspot.GET("/:router_id/hosts", h.GetHosts)
	}
}
