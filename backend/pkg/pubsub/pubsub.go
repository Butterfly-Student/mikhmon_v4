package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// PubSub wraps a Redis client to provide publish/subscribe functionality
// for broadcasting streaming RouterOS data to multiple consumers.
type PubSub struct {
	client *redis.Client
}

// New creates a new PubSub backed by the given Redis client.
func New(client *redis.Client) *PubSub {
	return &PubSub{client: client}
}

// Publish JSON-encodes data and publishes it to the given Redis channel.
func (ps *PubSub) Publish(ctx context.Context, channel string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("pubsub publish marshal: %w", err)
	}
	return ps.client.Publish(ctx, channel, b).Err()
}

// Subscribe returns a channel of raw JSON message bytes from a Redis channel.
// The returned cancel func unsubscribes and stops the goroutine.
func (ps *PubSub) Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error) {
	sub := ps.client.Subscribe(ctx, channel)

	// Verify subscription before returning.
	if _, err := sub.Receive(ctx); err != nil {
		_ = sub.Close()
		return nil, nil, fmt.Errorf("pubsub subscribe: %w", err)
	}

	out := make(chan []byte, 64)

	go func() {
		defer close(out)
		ch := sub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				select {
				case out <- []byte(msg.Payload):
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	cancel := func() {
		_ = sub.Close()
	}

	return out, cancel, nil
}

// ─── Channel name helpers ─────────────────────────────────────────────────────

// HotspotActiveChannel returns the Redis channel name for hotspot active streaming.
func HotspotActiveChannel(routerID uint) string {
	return fmt.Sprintf("mikhmon:stream:hotspot:active:%d", routerID)
}

// HotspotInactiveChannel returns the Redis channel name for hotspot inactive streaming.
func HotspotInactiveChannel(routerID uint) string {
	return fmt.Sprintf("mikhmon:stream:hotspot:inactive:%d", routerID)
}

// PPPActiveChannel returns the Redis channel name for PPP active streaming.
func PPPActiveChannel(routerID uint) string {
	return fmt.Sprintf("mikhmon:stream:ppp:active:%d", routerID)
}

// PPPInactiveChannel returns the Redis channel name for PPP inactive streaming.
func PPPInactiveChannel(routerID uint) string {
	return fmt.Sprintf("mikhmon:stream:ppp:inactive:%d", routerID)
}

// LogsChannel returns the Redis channel name for log streaming.
// topics is included in the channel name so different topic filters use separate channels.
func LogsChannel(routerID uint, topics string) string {
	return fmt.Sprintf("mikhmon:stream:logs:%d:%s", routerID, topics)
}
