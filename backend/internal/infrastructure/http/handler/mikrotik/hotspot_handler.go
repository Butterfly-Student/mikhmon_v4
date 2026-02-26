package mikrotik

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/infrastructure/http/handler"
	"github.com/irhabi89/mikhmon/internal/usecase/mikrotik"
	"go.uber.org/zap"
)

// HotspotHandler handles hotspot management
type HotspotHandler struct {
	handler.BaseHandler
	hotspotUC *mikrotik.HotspotUseCase
}

// NewHotspotHandler creates a new hotspot handler
func NewHotspotHandler(hotspotUC *mikrotik.HotspotUseCase, log *zap.Logger) *HotspotHandler {
	return &HotspotHandler{
		BaseHandler: handler.BaseHandler{Log: log.Named("hotspot")},
		hotspotUC:   hotspotUC,
	}
}

// GetActiveCount returns count of active sessions
func (h *HotspotHandler) GetActiveCount(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	count, err := h.hotspotUC.GetActiveCount(c.Request.Context(), uint(routerID))
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.Success(c, gin.H{"activeUsers": count})
}

// GetProfileByIDs gets a specific profile
func (h *HotspotHandler) GetProfileByIDs(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	profileID := c.Param("id")

	h.Log.Info("GetProfileByIDs called", zap.Uint64("routerID", routerID), zap.String("profileID", profileID))

	profile, err := h.hotspotUC.GetProfileByID(c.Request.Context(), uint(routerID), profileID)
	if err != nil {
		h.Log.Error("GetProfileByIDs failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}
	if profile == nil {
		h.ErrorWithCode(c, http.StatusNotFound, handler.ErrCodeNotFound, "Profile not found")
		return
	}

	h.Success(c, profile)
}

// GetProfileByName gets a profile by name
func (h *HotspotHandler) GetProfileByName(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	profileName := c.Param("name")

	h.Log.Info("GetProfileByName called", zap.Uint64("routerID", routerID), zap.String("profileName", profileName))

	profile, err := h.hotspotUC.GetProfileByName(c.Request.Context(), uint(routerID), profileName)
	if err != nil {
		h.Log.Error("GetProfileByName failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}
	if profile == nil {
		h.ErrorWithCode(c, http.StatusNotFound, handler.ErrCodeNotFound, "Profile not found")
		return
	}

	h.Success(c, profile)
}

// GetUsers lists all hotspot users
func (h *HotspotHandler) GetUsers(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	profile := c.Query("profile")

	h.Log.Info("GetUsers called", zap.Uint64("routerID", routerID), zap.String("profile", profile))

	users, err := h.hotspotUC.GetUsers(c.Request.Context(), uint(routerID), profile)
	if err != nil {
		h.Log.Error("GetUsers failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("GetUsers success", zap.Uint64("routerID", routerID), zap.Int("count", len(users)))
	h.Success(c, users)
}

// GetUser gets a specific user
func (h *HotspotHandler) GetUser(c *gin.Context) {
	routerID, err := strconv.ParseUint(c.Param("router_id"), 10, 32)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, "Invalid router ID")
		return
	}
	userID := c.Param("id")

	h.Log.Info("GetUser called", zap.Uint64("routerID", routerID), zap.String("userID", userID))

	user, err := h.hotspotUC.GetUser(c.Request.Context(), uint(routerID), userID)
	if err != nil {
		h.Log.Error("GetUser failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}
	if user == nil {
		h.ErrorWithCode(c, http.StatusNotFound, handler.ErrCodeNotFound, "User not found")
		return
	}

	h.Log.Info("GetUser success", zap.Uint64("routerID", routerID), zap.String("userID", userID))
	h.Success(c, user)
}

// CreateUser creates a new hotspot user
func (h *HotspotHandler) CreateUser(c *gin.Context) {
	routerID, err := strconv.ParseUint(c.Param("router_id"), 10, 32)
	if err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, "Invalid router ID")
		return
	}

	var req dto.AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}
	req.RouterID = uint(routerID)

	h.Log.Info("CreateUser called", zap.Uint64("routerID", routerID), zap.String("name", req.Name))

	if err := h.hotspotUC.AddUser(c.Request.Context(), &req); err != nil {
		h.Log.Error("CreateUser failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
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
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}

	h.Log.Info("UpdateUser called", zap.Uint64("routerID", routerID), zap.String("userID", userID))

	if err := h.hotspotUC.UpdateUser(c.Request.Context(), uint(routerID), userID, &req); err != nil {
		h.Log.Error("UpdateUser failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	if req.Reset {
		if err := h.hotspotUC.ResetUserCounters(c.Request.Context(), uint(routerID), userID); err != nil {
			h.Log.Warn("ResetUserCounters failed during update", zap.Uint64("routerID", routerID), zap.Error(err))
		}
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
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.SuccessWithMessage(c, "User deleted successfully", nil)
}

// GetUsersCount returns total users count
func (h *HotspotHandler) GetUsersCount(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	profile := c.Query("profile")

	users, err := h.hotspotUC.GetUsers(c.Request.Context(), uint(routerID), profile)
	if err != nil {
		h.Log.Error("GetUsersCount failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Success(c, gin.H{"totalUsers": len(users)})
}

// GetProfiles lists user profiles
func (h *HotspotHandler) GetProfiles(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	h.Log.Info("GetProfiles called", zap.Uint64("routerID", routerID))

	profiles, err := h.hotspotUC.GetProfiles(c.Request.Context(), uint(routerID))
	if err != nil {
		h.Log.Error("GetProfiles failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
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
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}

	h.Log.Info("CreateProfile called", zap.Uint64("routerID", routerID), zap.String("name", req.Name))

	if err := h.hotspotUC.AddProfile(c.Request.Context(), uint(routerID), &req); err != nil {
		h.Log.Error("CreateProfile failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
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
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}

	h.Log.Info("UpdateProfile called", zap.Uint64("routerID", routerID), zap.String("profileID", profileID))

	if err := h.hotspotUC.UpdateProfile(c.Request.Context(), uint(routerID), profileID, &req); err != nil {
		h.Log.Error("UpdateProfile failed", zap.Uint64("routerID", routerID), zap.Error(err))
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
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
		h.ErrorWithCode(c, http.StatusBadGateway, handler.ErrCodeMikrotikConnection, err.Error())
		return
	}

	h.SuccessWithMessage(c, "Profile deleted successfully", nil)
}

// GetActiveUsers lists active sessions
func (h *HotspotHandler) GetActiveUsers(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	h.Log.Info("GetActiveUsers called", zap.Uint64("routerID", routerID))

	users, err := h.hotspotUC.GetActive(c.Request.Context(), uint(routerID))
	if err != nil {
		h.Log.Error("GetActiveUsers failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("GetActiveUsers success", zap.Uint64("routerID", routerID), zap.Int("count", len(users)))
	h.Success(c, users)
}

// DeleteActiveUser removes an active hotspot session
func (h *HotspotHandler) DeleteActiveUser(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	activeID := c.Param("id")

	if err := h.hotspotUC.RemoveActive(c.Request.Context(), uint(routerID), activeID); err != nil {
		h.Log.Error("DeleteActiveUser failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.SuccessWithMessage(c, "Active session removed successfully", nil)
}

// GetHosts lists hotspot hosts
func (h *HotspotHandler) GetHosts(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	h.Log.Info("GetHosts called", zap.Uint64("routerID", routerID))

	hosts, err := h.hotspotUC.GetHosts(c.Request.Context(), uint(routerID))
	if err != nil {
		h.Log.Error("GetHosts failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Log.Info("GetHosts success", zap.Uint64("routerID", routerID), zap.Int("count", len(hosts)))
	h.Success(c, hosts)
}

// DeleteHost removes a hotspot host
func (h *HotspotHandler) DeleteHost(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)
	hostID := c.Param("id")

	if err := h.hotspotUC.RemoveHost(c.Request.Context(), uint(routerID), hostID); err != nil {
		h.Log.Error("DeleteHost failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.SuccessWithMessage(c, "Host removed successfully", nil)
}

// GetHotspotServers lists hotspot server names
func (h *HotspotHandler) GetHotspotServers(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	servers, err := h.hotspotUC.GetServers(c.Request.Context(), uint(routerID))
	if err != nil {
		h.Log.Error("GetHotspotServers failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	h.Success(c, servers)
}

// SetupExpireMonitor creates/enables Mikhmon expire monitor scheduler
func (h *HotspotHandler) SetupExpireMonitor(c *gin.Context) {
	routerID, _ := strconv.ParseUint(c.Param("router_id"), 10, 32)

	var req struct {
		Script string `json:"script"`
		ExpMon string `json:"expmon"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorWithCode(c, http.StatusBadRequest, handler.ErrCodeValidation, err.Error())
		return
	}

	script := req.Script
	if script == "" {
		script = req.ExpMon
	}

	status, err := h.hotspotUC.SetupExpireMonitor(c.Request.Context(), uint(routerID), script)
	if err != nil {
		h.Log.Error("SetupExpireMonitor failed", zap.Uint64("routerID", routerID), zap.Error(err))
		code := handler.ErrCodeMikrotikConnection
		if isMikrotikTimeout(err) {
			code = handler.ErrCodeMikrotikTimeout
		}
		h.ErrorWithCode(c, http.StatusBadGateway, code, err.Error())
		return
	}

	message := "success"
	if status == "existing" {
		message = "Mikhmon-Expire-Monitor"
	}
	h.Success(c, gin.H{"message": message, "status": status})
}

// GetExpireMonitorScript returns default expire monitor script used by backend.
func (h *HotspotHandler) GetExpireMonitorScript(c *gin.Context) {
	h.Success(c, gin.H{
		"script": h.hotspotUC.GetExpireMonitorScript(),
	})
}

// RegisterRoutes registers hotspot routes
func (h *HotspotHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("", h.GetUsers)
		users.GET("/count", h.GetUsersCount)
		users.GET("/:id", h.GetUser)
		users.POST("", h.CreateUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}

	profiles := r.Group("/profiles")
	{
		profiles.GET("", h.GetProfiles)
		profiles.GET("/by-name/:name", h.GetProfileByName)
		profiles.POST("", h.CreateProfile)
		profiles.GET("/:id", h.GetProfileByIDs)
		profiles.PUT("/:id", h.UpdateProfile)
		profiles.DELETE("/:id", h.DeleteProfile)
	}

	active := r.Group("/active")
	{
		active.GET("", h.GetActiveUsers)
		active.GET("/count", h.GetActiveCount)
		active.DELETE("/:id", h.DeleteActiveUser)
	}

	hosts := r.Group("/hosts")
	{
		hosts.GET("", h.GetHosts)
		hosts.DELETE("/:id", h.DeleteHost)
	}

	r.GET("/servers", h.GetHotspotServers)
	r.POST("/expire-monitor", h.SetupExpireMonitor)
	r.GET("/expire-monitor/script", h.GetExpireMonitorScript)
}

// isMikrotikTimeout checks if error is a context deadline / timeout error
func isMikrotikTimeout(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "deadline exceeded") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "timed out")
}
