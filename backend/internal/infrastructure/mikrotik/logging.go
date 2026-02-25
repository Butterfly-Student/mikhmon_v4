package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

// GetHotspotLogs retrieves hotspot logs
func (c *Client) GetHotspotLogs(ctx context.Context, router *entity.Router, limit int) ([]*dto.LogEntry, error) {
	// First ensure logging is configured
	if err := c.EnableHotspotLogging(ctx, router); err != nil {
		// Log error but continue
	}

	client, err := c.getClient(router)
	if err != nil {
		return nil, err
	}

	reply, err := client.RunContext(ctx, "/log/print", "?topics=hotspot,info,debug")
	if err != nil {
		return nil, err
	}

	logs := make([]*dto.LogEntry, 0, len(reply.Re))
	for i, re := range reply.Re {
		if limit > 0 && i >= limit {
			break
		}
		logs = append(logs, &dto.LogEntry{
			ID:      re.Map[".id"],
			Time:    re.Map["time"],
			Topics:  re.Map["topics"],
			Message: re.Map["message"],
		})
	}

	return logs, nil
}

// EnableHotspotLogging configures hotspot logging if not already configured
func (c *Client) EnableHotspotLogging(ctx context.Context, router *entity.Router) error {
	client, err := c.getClient(router)
	if err != nil {
		return err
	}

	// Check if logging is already configured
	reply, err := client.RunContext(ctx, "/system/logging/print", "?prefix=->")
	if err != nil {
		return err
	}

	// If already configured, return
	if len(reply.Re) > 0 {
		return nil
	}

	// Add logging configuration
	_, err = client.RunContext(ctx,
		"/system/logging/add",
		"=action=disk",
		"=prefix=->",
		"=topics=hotspot,info,debug",
	)
	return err
}
