package mikrotik

import (
	"context"
	"fmt"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// parsePPPActive maps a RouterOS sentence map to a PPPActive DTO.
func parsePPPActive(m map[string]string) *dto.PPPActive {
	bytesIn, bytesOut := SplitSlashValue(m["bytes"])
	packetsIn, packetsOut := SplitSlashValue(m["packets"])

	return &dto.PPPActive{
		ID:            m[".id"],
		Name:          m["name"],
		Service:       m["service"],
		CallerID:      m["caller-id"],
		Address:       m["address"],
		Uptime:        m["uptime"],
		SessionID:     m["session-id"],
		Encoding:      m["encoding"],
		BytesIn:       bytesIn,
		BytesOut:      bytesOut,
		PacketsIn:     packetsIn,
		PacketsOut:    packetsOut,
		LimitBytesIn:  parseInt(m["limit-bytes-in"]),
		LimitBytesOut: parseInt(m["limit-bytes-out"]),
	}
}

// GetPPPActive retrieves all active PPP sessions, optionally filtered by service.
func (c *Client) GetPPPActive(ctx context.Context, service string) ([]*dto.PPPActive, error) {
	args := []string{"/ppp/active/print"}
	if service != "" {
		args = append(args, "?service="+service)
	}

	reply, err := c.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}

	active := make([]*dto.PPPActive, 0, len(reply.Re))
	for _, re := range reply.Re {
		active = append(active, parsePPPActive(re.Map))
	}

	return active, nil
}

// GetPPPActiveByID retrieves an active PPP session by ID.
func (c *Client) GetPPPActiveByID(ctx context.Context, id string) (*dto.PPPActive, error) {
	reply, err := c.RunContext(ctx, "/ppp/active/print", "?.id="+id)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, nil
	}

	return parsePPPActive(reply.Re[0].Map), nil
}

// RemovePPPActive removes (disconnects) an active PPP session.
func (c *Client) RemovePPPActive(ctx context.Context, id string) error {
	_, err := c.RunContext(ctx, "/ppp/active/remove", "=.id="+id)
	return err
}

// ListenPPPActive starts listening to active PPP sessions using RouterOS streaming API.
// The caller controls the lifetime via ctx. resultChan is closed when the goroutine exits.
func (c *Client) ListenPPPActive(
	ctx context.Context,
	resultChan chan<- []*dto.PPPActive,
) (func() error, error) {

	args := []string{
		"/ppp/active/print",
		"=follow=",
	}

	listenReply, err := c.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start ppp active listen: %w", err)
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

				// Each sentence in a follow reply may contain one active entry
				entry := parsePPPActive(sentence.Map)
				batch := []*dto.PPPActive{entry}

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

// ListenPPPInactive streams inactive PPP secrets — configured secrets not in the
// active sessions table. Both /ppp/secret and /ppp/active are streamed using
// RouterOS =follow= so updates are purely server-driven.
//
// Each time either stream delivers a new snapshot the diff is recomputed and
// sent to resultChan. resultChan is closed when the goroutine exits.
func (c *Client) ListenPPPInactive(
	ctx context.Context,
	resultChan chan<- []*dto.PPPSecret,
) (func() error, error) {

	args := func(path string) []string {
		return []string{path, "=follow="}
	}

	secretsListen, err := c.ListenArgsContext(ctx, args("/ppp/secret/print"))
	if err != nil {
		return nil, fmt.Errorf("failed to start ppp secrets listen: %w", err)
	}

	activeListen, err := c.ListenArgsContext(ctx, args("/ppp/active/print"))
	if err != nil {
		secretsListen.Cancel() //nolint:errcheck
		return nil, fmt.Errorf("failed to start ppp active listen for inactive: %w", err)
	}

	secretsBatches := listenBatches(ctx, secretsListen.Chan(), batchDebounce)
	activeBatches := listenBatches(ctx, activeListen.Chan(), batchDebounce)

	go func() {
		defer close(resultChan)

		var latestSecrets []*dto.PPPSecret
		var latestActive []*dto.PPPActive

		sendDiff := func() {
			if latestSecrets == nil {
				return
			}
			activeSet := make(map[string]struct{}, len(latestActive))
			for _, a := range latestActive {
				activeSet[a.Name] = struct{}{}
			}
			inactive := make([]*dto.PPPSecret, 0)
			for _, s := range latestSecrets {
				if _, ok := activeSet[s.Name]; !ok {
					inactive = append(inactive, s)
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

			case batch, ok := <-secretsBatches:
				if !ok {
					return
				}
				latestSecrets = make([]*dto.PPPSecret, 0, len(batch))
				for _, s := range batch {
					latestSecrets = append(latestSecrets, parsePPPSecret(s.Map))
				}
				sendDiff()

			case batch, ok := <-activeBatches:
				if !ok {
					return
				}
				latestActive = make([]*dto.PPPActive, 0, len(batch))
				for _, s := range batch {
					latestActive = append(latestActive, parsePPPActive(s.Map))
				}
				sendDiff()
			}
		}
	}()

	return func() error {
		_, err1 := secretsListen.Cancel()
		_, err2 := activeListen.Cancel()
		if err1 != nil {
			return err1
		}
		return err2
	}, nil
}
