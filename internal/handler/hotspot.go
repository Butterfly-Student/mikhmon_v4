package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"mikhmon_v4/config"
	ros "mikhmon_v4/internal/routeros"
)

// HotspotHandler manages hotspot users, active sessions, hosts, and servers.
type HotspotHandler struct {
	store sessions.Store
	pool  *ros.Pool
}

func NewHotspotHandler(store sessions.Store, pool *ros.Pool) *HotspotHandler {
	return &HotspotHandler{store: store, pool: pool}
}

// HotspotPage renders the main hotspot management page.
func (h *HotspotHandler) HotspotPage(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.Redirect(http.StatusFound, "/admin/settings")
		return
	}
	c.HTML(http.StatusOK, "web/templates/hotspot.html", gin.H{
		"Title":   "Hotspot — " + router.HotspotName,
		"Session": session,
		"Router":  router,
		"Active":  "hotspot",
	})
}

// getClient is a helper to get the RouterOS client for the current session.
func (h *HotspotHandler) getClient(c *gin.Context) (interface{}, *config.RouterConfig) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return nil, nil
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, nil
	}
	return client, router
}

// GetUsers lists hotspot users, optionally filtered by profile.
func (h *HotspotHandler) GetUsers(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	args := []string{"/ip/hotspot/user/print"}
	if prof := c.Query("prof"); prof != "" && prof != "all" {
		args = append(args, "?profile="+prof)
	}
	users, err := ros.RunArgs(client, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser retrieves a single hotspot user by .id.
func (h *HotspotHandler) GetUser(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users, err := ros.RunArgs(client, "/ip/hotspot/user/print", "?.id="+c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, users[0])
}

// AddUser creates a new hotspot user.
func (h *HotspotHandler) AddUser(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args := []string{"/ip/hotspot/user/add"}
	for k, v := range body {
		args = append(args, "="+k, v)
	}
	result, err := ros.RunArgs(client, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": result})
}

// UpdateUser updates a hotspot user.
func (h *HotspotHandler) UpdateUser(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args := []string{"/ip/hotspot/user/set", "=.id=" + c.Param("id")}
	for k, v := range body {
		args = append(args, "="+k+"="+v)
	}
	_, err = ros.RunArgs(client, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}
	// Reset counters.
	ros.RunArgs(client, "/ip/hotspot/user/reset-counters", "=.id="+c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// RemoveUser deletes a hotspot user.
func (h *HotspotHandler) RemoveUser(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = ros.RunArgs(client, "/ip/hotspot/user/remove", "=.id="+c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetActiveSessions lists active hotspot sessions.
func (h *HotspotHandler) GetActiveSessions(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	active, err := ros.RunArgs(client, "/ip/hotspot/active/print")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, active)
}

// RemoveActiveSession kicks an active session.
func (h *HotspotHandler) RemoveActiveSession(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = ros.RunArgs(client, "/ip/hotspot/active/remove", "=.id="+c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetHosts lists hotspot hosts.
func (h *HotspotHandler) GetHosts(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hosts, err := ros.RunArgs(client, "/ip/hotspot/host/print")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hosts)
}

// RemoveHost deletes a hotspot host entry.
func (h *HotspotHandler) RemoveHost(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = ros.RunArgs(client, "/ip/hotspot/host/remove", "=.id="+c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": gin.H{"error": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// GetHotspotServers lists configured hotspot servers.
func (h *HotspotHandler) GetHotspotServers(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}
	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	servers, err := ros.RunArgs(client, "/ip/hotspot/print")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, servers)
}
