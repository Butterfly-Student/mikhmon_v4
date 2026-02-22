package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"mikhmon_v4/config"
	ros "mikhmon_v4/internal/routeros"
)

// TrafficHandler streams real-time interface traffic via Server-Sent Events.
type TrafficHandler struct {
	store sessions.Store
	pool  *ros.Pool
}

func NewTrafficHandler(store sessions.Store, pool *ros.Pool) *TrafficHandler {
	return &TrafficHandler{store: store, pool: pool}
}

// GetTraffic streams tx/rx bits-per-second for the selected interface via SSE.
// Query: iface=<interface-name>
func (h *TrafficHandler) GetTraffic(c *gin.Context) {
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

	iface := c.DefaultQuery("iface", "ether1")

	// Set SSE headers.
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	ctx := c.Request.Context()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			result, err := ros.RunArgs(client,
				"/interface/monitor-traffic",
				"=interface="+iface,
				"=once=",
			)
			if err != nil {
				return false
			}

			tx, rx := "0", "0"
			if len(result) > 0 {
				tx = result[0]["tx-bits-per-second"]
				rx = result[0]["rx-bits-per-second"]
			}

			data := map[string]string{"tx": tx, "rx": rx}
			b, _ := json.Marshal(data)
			fmt.Fprintf(w, "data: %s\n\n", b)
			return true
		}
	})
}
