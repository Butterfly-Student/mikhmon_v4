package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"mikhmon_v4/config"
	ros "mikhmon_v4/internal/routeros"
)

// SystemHandler handles logs and expire monitor.
type SystemHandler struct {
	store sessions.Store
	pool  *ros.Pool
}

func NewSystemHandler(store sessions.Store, pool *ros.Pool) *SystemHandler {
	return &SystemHandler{store: store, pool: pool}
}

// LogPage renders the system log page.
func (h *SystemHandler) LogPage(c *gin.Context) {
	session := c.Param("session")
	router := config.FindRouter(session)
	if router == nil {
		c.Redirect(http.StatusFound, "/admin/settings")
		return
	}
	c.HTML(http.StatusOK, "web/templates/log.html", gin.H{
		"Title":   "Log — " + router.HotspotName,
		"Session": session,
		"Router":  router,
		"Active":  "log",
	})
}

// GetLogs returns hotspot log entries.
func (h *SystemHandler) GetLogs(c *gin.Context) {
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

	// Ensure logging rule exists.
	loggingRules, _ := ros.RunArgs(client, "/system/logging/print", "?prefix=->")
	hasRule := false
	for _, r := range loggingRules {
		if r["prefix"] == "->" {
			hasRule = true
			break
		}
	}
	if !hasRule {
		ros.RunArgs(client, "/system/logging/add",
			"=action=disk", "=prefix=->", "=topics=hotspot,info,debug")
	}

	logs, err := ros.RunArgs(client, "/log/print", "?topics=hotspot, info, debug")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reverse for latest-first.
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}
	c.JSON(http.StatusOK, logs)
}

// GetExpireMonitor returns expire monitor scheduler entries.
func (h *SystemHandler) GetExpireMonitor(c *gin.Context) {
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

	schedulers, err := ros.RunArgs(client, "/system/scheduler/print", "?comment=mikhmon_expire")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedulers)
}

// SetExpireMonitor creates, enables, or disables the expire monitor scheduler.
func (h *SystemHandler) SetExpireMonitor(c *gin.Context) {
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

	var body struct {
		Action   string `json:"action"`   // "add" | "enable" | "disable" | "remove"
		Interval string `json:"interval"` // e.g. "1m"
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find existing scheduler.
	schedulers, _ := ros.RunArgs(client, "/system/scheduler/print", "?comment=mikhmon_expire")

	switch body.Action {
	case "add":
		interval := body.Interval
		if interval == "" {
			interval = "1m"
		}
		now := time.Now()
		startDate := fmt.Sprintf("%s/%02d/%d",
			strings.ToUpper(now.Format("Jan")), now.Day(), now.Year())

		if len(schedulers) > 0 {
			id := schedulers[0][".id"]
			ros.RunArgs(client, "/system/scheduler/set",
				"=.id="+id,
				"=interval="+interval,
				"=disabled=no",
			)
		} else {
			script := expireMonitorScript()
			ros.RunArgs(client, "/system/scheduler/add",
				"=name=mikhmon_expire_monitor",
				"=start-date="+startDate,
				"=interval="+interval,
				"=comment=mikhmon_expire",
				"=on-event="+script,
			)
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})

	case "enable":
		if len(schedulers) > 0 {
			ros.RunArgs(client, "/system/scheduler/set",
				"=.id="+schedulers[0][".id"], "=disabled=no")
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})

	case "disable":
		if len(schedulers) > 0 {
			ros.RunArgs(client, "/system/scheduler/set",
				"=.id="+schedulers[0][".id"], "=disabled=yes")
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})

	case "remove":
		if len(schedulers) > 0 {
			ros.RunArgs(client, "/system/scheduler/remove",
				"=.id="+schedulers[0][".id"])
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown action"})
	}
}

// expireMonitorScript returns the RouterOS script for the expire monitor.
func expireMonitorScript() string {
	return `:local users [/ip hotspot user find where comment~"N" or comment~"X"];` +
		`:foreach u in=$users do={` +
		`:local name [/ip hotspot user get $u name];` +
		`:local comment [/ip hotspot user get $u comment];` +
		`:local expdate [:pick $comment 0 11];` +
		`:local mode [:pick $comment 12 13];` +
		`:local now [/system clock get date];` +
		`:if ($expdate <= $now) do={` +
		`:if ($mode = "X") do={ /ip hotspot user remove $u };` +
		`:if ($mode = "N") do={ /ip hotspot user set limit-uptime=1s $u };` +
		`}};`
}
