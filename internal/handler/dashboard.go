package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"mikhmon_v4/config"
	ros "mikhmon_v4/internal/routeros"
)

// DashboardHandler serves system info and connection endpoints.
type DashboardHandler struct {
	store sessions.Store
	pool  *ros.Pool
}

func NewDashboardHandler(store sessions.Store, pool *ros.Pool) *DashboardHandler {
	return &DashboardHandler{store: store, pool: pool}
}

// DashboardPage renders the dashboard HTML page.
func (h *DashboardHandler) DashboardPage(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.Redirect(http.StatusFound, "/admin/settings")
		return
	}
	c.HTML(http.StatusOK, "web/templates/dashboard.html", gin.H{
		"Title":   "Dashboard — " + router.HotspotName,
		"Session": session,
		"Router":  router,
		"Active":  "dashboard",
	})
}

// GetDashboard handles the AJAX dashboard data requests.
// ?page=get_sys_resource | get_hotspotinfo
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
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

	page := c.Query("page")
	switch page {
	case "get_sys_resource":
		systime, _ := ros.RunArgs(client, "/system/clock/print")
		resource, _ := ros.RunArgs(client, "/system/resource/print")
		routerboard, _ := ros.RunArgs(client, "/system/routerboard/print")
		identity, _ := ros.RunArgs(client, "/system/identity/print")
		health, _ := ros.RunArgs(client, "/system/health/print")

		resp := gin.H{
			"systime":   safeFirst(systime),
			"resource":  safeFirst(resource),
			"syshealth": safeFirst(health),
			"model":     safeGet(routerboard, "model"),
			"identity":  safeGet(identity, "name"),
		}
		c.JSON(http.StatusOK, resp)

	case "get_hotspotinfo":
		users, _ := ros.RunArgs(client, "/ip/hotspot/user/print", "count-only=")
		active, _ := ros.RunArgs(client, "/ip/hotspot/active/print", "count-only=")

		userCount := 0
		if len(users) > 0 {
			if n := users[0]["ret"]; n != "" {
				fmt.Sscan(n, &userCount)
				userCount-- // exclude default user
				if userCount < 0 {
					userCount = 0
				}
			}
		}
		activeCount := 0
		if len(active) > 0 {
			fmt.Sscan(active[0]["ret"], &activeCount)
		}

		c.JSON(http.StatusOK, gin.H{
			"hotspot_users":  userCount,
			"hotspot_active": activeCount,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown page"})
	}
}

// TestConnection checks connectivity to the router.
func (h *DashboardHandler) TestConnection(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}

	client, err := h.pool.Get(session, router.Host, router.Username, router.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	identity, err := ros.RunArgs(client, "/system/identity/print")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}

	name := ""
	if len(identity) > 0 {
		name = identity[0]["name"]
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "identity": name})
}

// safeFirst returns the first map in a slice or an empty map.
func safeFirst(rows []map[string]string) map[string]string {
	if len(rows) > 0 {
		return rows[0]
	}
	return map[string]string{}
}

// safeGet returns a key from the first row of a result slice.
func safeGet(rows []map[string]string, key string) string {
	if len(rows) > 0 {
		return rows[0][key]
	}
	return ""
}
