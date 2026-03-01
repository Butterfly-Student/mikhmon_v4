package mikrotik

import (
	"context"
	"fmt"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)


// GetHotspotActive retrieves active hotspot sessions.
func (c *Client) GetHotspotActive(ctx context.Context) ([]*dto.HotspotActive, error) {
	reply, err := c.RunContext(ctx, "/ip/hotspot/active/print")
	if err != nil {
		return nil, err
	}

	active := make([]*dto.HotspotActive, 0, len(reply.Re))
	for _, re := range reply.Re {
		active = append(active, &dto.HotspotActive{
			ID:               re.Map[".id"],
			Server:           re.Map["server"],
			User:             re.Map["user"],
			Address:          re.Map["address"],
			MACAddress:       re.Map["mac-address"],
			LoginBy:          re.Map["login-by"],
			Uptime:           re.Map["uptime"],
			SessionTimeLeft:  re.Map["session-time-left"],
			IdleTime:         re.Map["idle-time"],
			IdleTimeout:      re.Map["idle-timeout"],
			KeepaliveTimeout: re.Map["keepalive-timeout"],
			BytesIn:          parseInt(re.Map["bytes-in"]),
			BytesOut:         parseInt(re.Map["bytes-out"]),
			PacketsIn:        parseInt(re.Map["packets-in"]),
			PacketsOut:       parseInt(re.Map["packets-out"]),
			Radius:           parseBool(re.Map["radius"]),
		})
	}

	return active, nil
}

// GetHotspotActiveCount retrieves the count of active hotspot sessions.
func (c *Client) GetHotspotActiveCount(ctx context.Context) (int, error) {
	reply, err := c.RunContext(ctx, "/ip/hotspot/active/print", "=count-only=")
	if err != nil {
		return 0, err
	}

	if len(reply.Re) > 0 {
		return int(parseInt(reply.Re[0].Map["ret"])), nil
	}

	return 0, nil
}

// RemoveHotspotActive removes an active hotspot session.
func (c *Client) RemoveHotspotActive(ctx context.Context, id string) error {
	_, err := c.RunContext(ctx, "/ip/hotspot/active/remove", "=.id="+id)
	return err
}

// parseHotspotActive maps a RouterOS sentence map to a HotspotActive DTO.
func parseHotspotActive(m map[string]string) *dto.HotspotActive {
	return &dto.HotspotActive{
		ID:               m[".id"],
		Server:           m["server"],
		User:             m["user"],
		Address:          m["address"],
		MACAddress:       m["mac-address"],
		LoginBy:          m["login-by"],
		Uptime:           m["uptime"],
		SessionTimeLeft:  m["session-time-left"],
		IdleTime:         m["idle-time"],
		IdleTimeout:      m["idle-timeout"],
		KeepaliveTimeout: m["keepalive-timeout"],
		BytesIn:          parseInt(m["bytes-in"]),
		BytesOut:         parseInt(m["bytes-out"]),
		PacketsIn:        parseInt(m["packets-in"]),
		PacketsOut:       parseInt(m["packets-out"]),
		Radius:           parseBool(m["radius"]),
	}
}

// ListenHotspotActive starts listening to active hotspot sessions using the RouterOS follow API.
// The caller controls the lifetime via ctx. resultChan is closed when the goroutine exits.
func (c *Client) ListenHotspotActive(
	ctx context.Context,
	resultChan chan<- []*dto.HotspotActive,
) (func() error, error) {
	args := []string{
		"/ip/hotspot/active/print",
		"=follow=",
	}

	listenReply, err := c.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start hotspot active listen: %w", err)
	}

	go func() {
		defer close(resultChan)

		for {
			select {
			case <-ctx.Done():
				listenReply.Cancel()
				return

			case sentence, ok := <-listenReply.Chan():
				if !ok {
					return
				}

				entry := parseHotspotActive(sentence.Map)
				batch := []*dto.HotspotActive{entry}

				select {
				case resultChan <- batch:
				case <-ctx.Done():
					listenReply.Cancel()
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

// ListenHotspotInactive streams inactive hotspot users — configured users not in
// the active sessions table. Both /ip/hotspot/user and /ip/hotspot/active are
// streamed using RouterOS =follow=  so updates are server-driven.
//
// Each time either stream delivers a new snapshot the diff is recomputed and
// sent to resultChan. resultChan is closed when the goroutine exits.
func (c *Client) ListenHotspotInactive(
	ctx context.Context,
	resultChan chan<- []*dto.HotspotUser,
) (func() error, error) {

	args := func(path string) []string {
		return []string{path, "=follow="}
	}

	usersListen, err := c.ListenArgsContext(ctx, args("/ip/hotspot/user/print"))
	if err != nil {
		return nil, fmt.Errorf("failed to start hotspot users listen: %w", err)
	}

	activeListen, err := c.ListenArgsContext(ctx, args("/ip/hotspot/active/print"))
	if err != nil {
		usersListen.Cancel() //nolint:errcheck
		return nil, fmt.Errorf("failed to start hotspot active listen for inactive: %w", err)
	}

	usersBatches := listenBatches(ctx, usersListen.Chan(), batchDebounce)
	activeBatches := listenBatches(ctx, activeListen.Chan(), batchDebounce)

	go func() {
		defer close(resultChan)

		var latestUsers []*dto.HotspotUser
		var latestActive []*dto.HotspotActive

		sendDiff := func() {
			if latestUsers == nil {
				return
			}
			activeSet := make(map[string]struct{}, len(latestActive))
			for _, a := range latestActive {
				activeSet[a.User] = struct{}{}
			}
			inactive := make([]*dto.HotspotUser, 0)
			for _, u := range latestUsers {
				if _, ok := activeSet[u.Name]; !ok {
					inactive = append(inactive, u)
				}
			}
			select {
			case resultChan <- inactive:
			case <-ctx.Done():
			}
		}

		for {
			select {
			case <-ctx.Done():
				return

			case batch, ok := <-usersBatches:
				if !ok {
					return
				}
				latestUsers = make([]*dto.HotspotUser, 0, len(batch))
				for _, s := range batch {
					latestUsers = append(latestUsers, &dto.HotspotUser{
						ID:              s.Map[".id"],
						Server:          s.Map["server"],
						Name:            s.Map["name"],
						Password:        s.Map["password"],
						Profile:         s.Map["profile"],
						MACAddress:      s.Map["mac-address"],
						IPAddress:       s.Map["address"],
						Uptime:          s.Map["uptime"],
						BytesIn:         parseInt(s.Map["bytes-in"]),
						BytesOut:        parseInt(s.Map["bytes-out"]),
						LimitUptime:     s.Map["limit-uptime"],
						LimitBytesTotal: parseInt(s.Map["limit-bytes-total"]),
						Comment:         s.Map["comment"],
						Disabled:        parseBool(s.Map["disabled"]),
					})
				}
				sendDiff()

			case batch, ok := <-activeBatches:
				if !ok {
					return
				}
				latestActive = make([]*dto.HotspotActive, 0, len(batch))
				for _, s := range batch {
					latestActive = append(latestActive, parseHotspotActive(s.Map))
				}
				sendDiff()
			}
		}
	}()

	return func() error {
		_, err1 := usersListen.Cancel()
		_, err2 := activeListen.Cancel()
		if err1 != nil {
			return err1
		}
		return err2
	}, nil
}
