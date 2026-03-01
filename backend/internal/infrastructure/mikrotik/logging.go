package mikrotik

import (
	"context"
	"fmt"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// ─── Snapshot (non-streaming) ─────────────────────────────────────────────────

// GetLogs retrieves a snapshot of current log entries from /log/print.
// If topics is non-empty (e.g. "hotspot,info" or "ppp,pppoe"), only matching
// entries are returned. Pass topics="" to get all logs.
// If limit > 0, at most limit entries are returned.
func (c *Client) GetLogs(ctx context.Context, topics string, limit int) ([]*dto.LogEntry, error) {
	args := []string{"/log/print"}
	if topics != "" {
		args = append(args, fmt.Sprintf("?topics=%s", topics))
	}

	reply, err := c.RunContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	logs := make([]*dto.LogEntry, 0, len(reply.Re))
	for i, re := range reply.Re {
		if limit > 0 && i >= limit {
			break
		}
		logs = append(logs, parseLogEntry(re.Map))
	}

	return logs, nil
}

// GetHotspotLogs retrieves hotspot log entries (best-effort: enables logging first).
func (c *Client) GetHotspotLogs(ctx context.Context, limit int) ([]*dto.LogEntry, error) {
	_ = c.EnableHotspotLogging(ctx)
	return c.GetLogs(ctx, "hotspot,info,debug", limit)
}

// GetPPPLogs retrieves PPP log entries (best-effort: enables logging first).
func (c *Client) GetPPPLogs(ctx context.Context, limit int) ([]*dto.LogEntry, error) {
	_ = c.EnablePPPLogging(ctx)
	return c.GetLogs(ctx, "ppp,pppoe,info", limit)
}

// ─── Logging configuration ────────────────────────────────────────────────────

// EnableHotspotLogging configures hotspot logging if not already configured.
func (c *Client) EnableHotspotLogging(ctx context.Context) error {
	reply, err := c.RunContext(ctx, "/system/logging/print", "?prefix=->")
	if err != nil {
		return err
	}

	if len(reply.Re) > 0 {
		return nil
	}

	_, err = c.RunContext(ctx,
		"/system/logging/add",
		"=action=disk",
		"=prefix=->",
		"=topics=hotspot,info,debug",
	)
	return err
}

// EnablePPPLogging configures PPP logging if not already configured.
func (c *Client) EnablePPPLogging(ctx context.Context) error {
	reply, err := c.RunContext(ctx, "/system/logging/print", "?prefix=ppp->")
	if err != nil {
		return err
	}

	if len(reply.Re) > 0 {
		return nil
	}

	_, err = c.RunContext(ctx,
		"/system/logging/add",
		"=action=disk",
		"=prefix=ppp->",
		"=topics=ppp,pppoe,info",
	)
	return err
}

// ─── Streaming ────────────────────────────────────────────────────────────────

// parseLogEntry maps a RouterOS sentence map to a LogEntry DTO.
func parseLogEntry(m map[string]string) *dto.LogEntry {
	return &dto.LogEntry{
		ID:      m[".id"],
		Time:    m["time"],
		Topics:  m["topics"],
		Message: m["message"],
	}
}

// ListenLogs streams new log entries using /log/print =follow-only=.
// =follow-only= skips the historical backlog — only entries arriving after
// the subscription starts are delivered.
// Pass topics="" to stream all logs, or e.g. "hotspot,info" for filtered streaming.
// The caller controls the lifetime via ctx. resultChan is closed when the goroutine exits.
func (c *Client) ListenLogs(
	ctx context.Context,
	topics string,
	resultChan chan<- *dto.LogEntry,
) (func() error, error) {
	args := []string{"/log/print", "=follow-only="}
	if topics != "" {
		args = append(args, fmt.Sprintf("?topics=%s", topics))
	}

	listenReply, err := c.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start log listen: %w", err)
	}

	go func() {
		defer close(resultChan)

		for {
			select {
			case <-ctx.Done():
				listenReply.Cancel() //nolint:errcheck
				return

			case sentence, ok := <-listenReply.Chan():
				if !ok {
					return
				}

				select {
				case resultChan <- parseLogEntry(sentence.Map):
				case <-ctx.Done():
					listenReply.Cancel() //nolint:errcheck
					return
				}
			}
		}
	}()

	return func() error {
		_, err := listenReply.Cancel()
		return err
	}, nil
}

// ListenAllLogs streams all new log entries (no topic filter).
func (c *Client) ListenAllLogs(
	ctx context.Context,
	resultChan chan<- *dto.LogEntry,
) (func() error, error) {
	return c.ListenLogs(ctx, "", resultChan)
}

// ListenHotspotLogs streams hotspot-related log entries.
func (c *Client) ListenHotspotLogs(
	ctx context.Context,
	resultChan chan<- *dto.LogEntry,
) (func() error, error) {
	return c.ListenLogs(ctx, "hotspot,info", resultChan)
}

// ListenPPPLogs streams PPP-related log entries.
func (c *Client) ListenPPPLogs(
	ctx context.Context,
	resultChan chan<- *dto.LogEntry,
) (func() error, error) {
	return c.ListenLogs(ctx, "ppp,pppoe,info", resultChan)
}
