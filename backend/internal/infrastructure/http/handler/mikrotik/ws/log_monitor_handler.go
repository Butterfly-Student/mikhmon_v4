package ws

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/repository"
	"github.com/irhabi89/mikhmon/internal/infrastructure/mikrotik"
	"github.com/irhabi89/mikhmon/pkg/pubsub"
	"go.uber.org/zap"
)

// LogMonitorHandler handles WebSocket log streaming for all, hotspot, or PPP logs.
// Features:
// - Unified format for history and real-time entries
// - Batched real-time delivery (reduces WebSocket message count)
// - Sequence ID for ordering
// - Circular buffer management (max 1000 entries)
type LogMonitorHandler struct {
	upgrader      websocket.Upgrader
	routerRepo    repository.RouterRepository
	mikrotikSvc   *mikrotik.Manager
	ps            *pubsub.PubSub
	topics        string
	internalWSKey string
	log           *zap.Logger
}

func newLogMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	ps *pubsub.PubSub,
	topics string,
	internalWSKey string,
	log *zap.Logger,
	name string,
) *LogMonitorHandler {
	return &LogMonitorHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  1024,
			WriteBufferSize: 4096,
		},
		routerRepo:    routerRepo,
		mikrotikSvc:   mikrotikSvc,
		ps:            ps,
		topics:        topics,
		internalWSKey: internalWSKey,
		log:           log.Named(name),
	}
}

// NewLogMonitorHandler streams all logs (topics="").
func NewLogMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	ps *pubsub.PubSub,
	internalWSKey string,
	log *zap.Logger,
) *LogMonitorHandler {
	return newLogMonitorHandler(routerRepo, mikrotikSvc, ps, "", internalWSKey, log, "ws-log-monitor")
}

// NewHotspotLogMonitorHandler streams hotspot logs (topics="hotspot,info").
func NewHotspotLogMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	ps *pubsub.PubSub,
	internalWSKey string,
	log *zap.Logger,
) *LogMonitorHandler {
	return newLogMonitorHandler(routerRepo, mikrotikSvc, ps, "hotspot,info", internalWSKey, log, "ws-hotspot-log-monitor")
}

// NewPPPLogMonitorHandler streams PPP logs (topics="ppp,pppoe,info").
func NewPPPLogMonitorHandler(
	routerRepo repository.RouterRepository,
	mikrotikSvc *mikrotik.Manager,
	ps *pubsub.PubSub,
	internalWSKey string,
	log *zap.Logger,
) *LogMonitorHandler {
	return newLogMonitorHandler(routerRepo, mikrotikSvc, ps, "ppp,pppoe,info", internalWSKey, log, "ws-ppp-log-monitor")
}

// logEntryWithSeq wraps log entry with sequence ID
type logEntryWithSeq struct {
	Seq     int64       `json:"seq"`
	Time    int64       `json:"time_ms"` // Unix timestamp in milliseconds
	Entry   *dto.LogEntry `json:"entry"`
}

// Handle handles WebSocket connections for log streaming.
func (h *LogMonitorHandler) Handle(c *gin.Context) {
	if c.Query("key") != h.internalWSKey {
		h.log.Warn("WebSocket auth failed", zap.String("key", c.Query("key")))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	routerID, err := strconv.ParseUint(c.Param("router_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid router_id"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("WebSocket upgrade failed", zap.Uint64("routerID", routerID), zap.Error(err))
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Keepalive: ping every 20s, expect pong within 60s.
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
					cancel()
					return
				}
			}
		}
	}()

	router, err := h.routerRepo.GetByID(ctx, uint(routerID))
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "router not found"}) //nolint:errcheck
		return
	}

	cfg := mikrotik.Config{
		Host:     router.Host,
		Port:     router.Port,
		Username: router.Username,
		Password: router.Password,
		UseTLS:   router.UseSSL,
		Timeout:  time.Duration(router.Timeout) * time.Second,
	}
	client, err := h.mikrotikSvc.GetOrConnect(ctx, router.Name, cfg)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "router not connected: " + err.Error()}) //nolint:errcheck
		return
	}

	// Sequence counter for ordering
	var seqCounter int64 = 0
	var seqMu sync.Mutex
	getNextSeq := func() int64 {
		seqMu.Lock()
		defer seqMu.Unlock()
		seqCounter++
		return seqCounter
	}

	// 1. Send history snapshot (up to 1000 entries) with sequence IDs.
	history, err := client.GetLogs(ctx, h.topics, 1000)
	if err != nil {
		h.log.Warn("Failed to fetch log history", zap.Uint64("routerID", routerID), zap.Error(err))
		// Continue even if history fetch fails
		history = []*dto.LogEntry{}
	}

	// Convert history to unified format with sequence
	historyWithSeq := make([]logEntryWithSeq, len(history))
	for i, entry := range history {
		historyWithSeq[i] = logEntryWithSeq{
			Seq:   getNextSeq(),
			Time:  parseLogTime(entry.Time),
			Entry: entry,
		}
	}

	// Send initial state
	if err := conn.WriteJSON(gin.H{
		"type":     "init",
		"data":     historyWithSeq,
		"meta": gin.H{
			"count":    len(historyWithSeq),
			"topics":   h.topics,
			"maxSize":  1000,
			"routerID": routerID,
		},
	}); err != nil {
		h.log.Warn("Failed to write init", zap.Error(err))
		return
	}

	// 2. Start follow-only stream with batching.
	resultChan := make(chan *dto.LogEntry, 128)
	cancelFn, err := client.ListenLogs(ctx, h.topics, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": "failed to start log stream: " + err.Error()}) //nolint:errcheck
		return
	}
	defer cancelFn() //nolint:errcheck

	// Batching goroutine - sends batches every 100ms or when buffer is full
	batchChan := make(chan []logEntryWithSeq, 10)
	go func() {
		defer close(batchChan)
		
		batch := make([]logEntryWithSeq, 0, 50)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		flush := func() {
			if len(batch) > 0 {
				batchCopy := make([]logEntryWithSeq, len(batch))
				copy(batchCopy, batch)
				select {
				case batchChan <- batchCopy:
					batch = batch[:0] // reset
				case <-ctx.Done():
				}
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case entry, ok := <-resultChan:
				if !ok {
					flush()
					return
				}
				
				seq := getNextSeq()
				wrapped := logEntryWithSeq{
					Seq:   seq,
					Time:  time.Now().UnixMilli(),
					Entry: entry,
				}
				
				batch = append(batch, wrapped)
				
				// Flush if batch is full
				if len(batch) >= 50 {
					flush()
				}
				
				// Publish to PubSub if available
				if h.ps != nil {
					ch := pubsub.LogsChannel(uint(routerID), h.topics)
					if err := h.ps.Publish(ctx, ch, entry); err != nil {
						h.log.Debug("pubsub publish failed", zap.Error(err))
					}
				}
				
			case <-ticker.C:
				flush()
			}
		}
	}()

	// Forward goroutine - sends batched updates to client
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case batch, ok := <-batchChan:
				if !ok {
					// Stream ended
					conn.WriteJSON(gin.H{"type": "status", "status": "ended"}) //nolint:errcheck
					cancel()
					return
				}
				
				if err := conn.WriteJSON(gin.H{
					"type": "update",
					"data": batch,
					"meta": gin.H{
						"batchSize": len(batch),
						"totalSeq":  seqCounter,
					},
				}); err != nil {
					h.log.Warn("Failed to write log batch", zap.Error(err))
					cancel()
					return
				}
			}
		}
	}()

	h.log.Info("Log monitor started", 
		zap.Uint64("routerID", routerID), 
		zap.String("topics", h.topics),
		zap.Int("historyCount", len(historyWithSeq)),
	)

	// Read loop: detect client disconnect.
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.log.Info("Log monitor WS closed", zap.Uint64("routerID", routerID), zap.Error(err))
			return
		}
	}
}

// parseLogTime converts RouterOS time string to Unix milliseconds.
// RouterOS format: "mar/01 13:19:12" or similar
func parseLogTime(timeStr string) int64 {
	// Try multiple formats
	formats := []string{
		"jan/02 15:04:05",
		"feb/02 15:04:05",
		"mar/02 15:04:05",
		"apr/02 15:04:05",
		"may/02 15:04:05",
		"jun/02 15:04:05",
		"jul/02 15:04:05",
		"aug/02 15:04:05",
		"sep/02 15:04:05",
		"oct/02 15:04:05",
		"nov/02 15:04:05",
		"dec/02 15:04:05",
	}
	
	// Normalize to lowercase for parsing
	timeStrLower := ""
	if len(timeStr) >= 3 {
		timeStrLower = timeStr[:3] + timeStr[3:]
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, timeStrLower); err == nil {
			// Assume current year since RouterOS doesn't include year
			now := time.Now()
			t = t.AddDate(now.Year(), 0, 0)
			// Adjust if the result is in the future
			if t.After(now) {
				t = t.AddDate(-1, 0, 0)
			}
			return t.UnixMilli()
		}
	}
	
	// Fallback to current time if parsing fails
	return time.Now().UnixMilli()
}
